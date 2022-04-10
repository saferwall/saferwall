// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"sync"

	"github.com/digitalocean/go-libvirt"
	gonsq "github.com/nsqio/go-nsq"
	agent "github.com/saferwall/saferwall/internal/agent"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	"github.com/saferwall/saferwall/internal/utils"
	"github.com/saferwall/saferwall/internal/vmmanager"
	"github.com/saferwall/saferwall/services/config"
	pb "github.com/saferwall/saferwall/services/proto"
	"google.golang.org/protobuf/proto"
)

// Config represents our application config.
type Config struct {
	LogLevel     string             `mapstructure:"log_level"`
	SharedVolume string             `mapstructure:"shared_volume"`
	Producer     config.ProducerCfg `mapstructure:"producer"`
	Consumer     config.ConsumerCfg `mapstructure:"consumer"`
	Agent        AgentCfg           `mapstructure:"agent"`
	virtMgr      VirtManagerCfg     `mapstructure:"virt_manager"`
	snapshotName string             `mapstructure:"snapshot_name"`
}

// AgentCfg represents the guest agent config.
type AgentCfg struct {
	// Destinary directory inside the guest where the agent is deployed.
	AgentDestDir string `mapstructure:"dest_dir"`
	// The sandbox binary components.
	PackageName string `mapstructure:"package_name"`
}

// VirtManagerCfg represents the virtualization manager config.
// For now, only libvirt server.
type VirtManagerCfg struct {
	Network string `mapstructure:"network"`
	Address string `mapstructure:"address"`
	port    string `mapstructure:"port"`
	user    string `mapstructure:"user"`
}

// VM represents a virtual machine config.
type VM struct {
	// ID identify uniquely the VM.
	ID int32
	// Name of the VM.
	Name string
	// IP address of the VM.
	IP string
	// Snapshots list names.
	Snapshots []string
	// InUse represents the availability of the VM.
	InUse bool
	// Pointer to the domain object.
	Dom *libvirt.Domain
}

// Service represents the sandbox scan service. It adheres to the nsq.Handler
// interface. This allows us to define our own custom handlers for our messages.
// Think of these handlers much like you would an http handler.
type Service struct {
	cfg     Config
	mu      sync.Mutex
	logger  log.Logger
	pub     pubsub.Publisher
	sub     pubsub.Subscriber
	vms     []VM
	vmm     vmmanager.VMManager
	sandbox []byte
}

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// New create a new sandbox service.
func New(cfg Config, logger log.Logger) (*Service, error) {
	var err error
	s := Service{}

	// retrieve the list of active VMs.
	conn, err := vmmanager.New(cfg.virtMgr.Network, cfg.virtMgr.Address,
		cfg.virtMgr.port, cfg.virtMgr.user)
	if err != nil {
		return nil, err
	}
	dd, err := conn.Domains()
	if err != nil {
		return nil, err
	}

	var vms []VM
	for _, d := range dd {
		vms = append(vms, VM{
			ID:        d.Dom.ID,
			Name:      d.Dom.Name,
			IP:        d.IP,
			Snapshots: d.Snapshots,
		})
	}

	// the number of concurrent workers have to match the number of
	// available virtual machines.
	s.sub, err = nsq.NewSubscriber(
		cfg.Consumer.Topic,
		cfg.Consumer.Channel,
		cfg.Consumer.Lookupds,
		len(vms),
		&s,
	)
	if err != nil {
		return nil, err
	}

	s.pub, err = nsq.NewPublisher(cfg.Producer.Nsqd)
	if err != nil {
		return nil, err
	}

	// download the sandbox release package.
	zipPackageData, err := utils.ReadAll(s.cfg.Agent.PackageName)
	if err != nil {
		return nil, err
	}

	s.sandbox = zipPackageData
	s.vms = vms
	s.cfg = cfg
	s.logger = logger
	s.vmm = conn
	return &s, nil

}

// Start kicks in the service to start consuming events.
func (s *Service) Start() error {
	s.logger.Infof("start consuming from topic: %s ...", s.cfg.Consumer.Topic)
	s.sub.Start()

	return nil
}

// HandleMessage is the only requirement needed to fulfill the nsq.Handler.
func (s *Service) HandleMessage(m *gonsq.Message) error {
	if len(m.Body) == 0 {
		return errors.New("body is blank re-enqueue message")
	}

	fileScanCfg := config.FileScanCfg{}
	ctx := context.Background()

	// Deserialize the msg sent from the web apis.
	err := json.Unmarshal(m.Body, &fileScanCfg)
	if err != nil {
		s.logger.Errorf("failed unmarshalling json messge body: %v", err)
		return err
	}

	sha256 := fileScanCfg.SHA256
	logger := s.logger.With(ctx, "sha256", sha256)
	logger.Info("start processing")

	// Find a free virtual machine to process this job.
	vm := s.findFreeVM()
	if vm == nil {
		logger.Infof("no VM currently available, call 911")
		return errors.New("failed to find a free VM")
	}

	logger = s.logger.With(ctx, "VM", vm.Name)
	logger.Infof("VM %s was selected", vm.Name)

	// Revert the VM to a clean state.
	err = s.vmm.Revert(*vm.Dom, s.cfg.snapshotName)
	if err != nil {
		logger.Errorf("failed to revert the VM: %v", err)
	}

	// Perform the actual detonation.
	res, err := s.detonate(logger, vm, sha256, fileScanCfg.Dynamic)
	if err != nil {
		logger.Errorf("failed to detonation the sample: %v", err)
		s.freeVM(vm)
		return err
	}

	// Make sure to free the VM for next job.
	s.freeVM(vm)

	payloads := []*pb.Message_Payload{
		{Module: "sandbox", Body: toJSON(res)},
	}

	msg := &pb.Message{Sha256: sha256, Payload: payloads}
	peMsg, err := proto.Marshal(msg)
	if err != nil {
		logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	err = s.pub.Publish(ctx, s.cfg.Producer.Topic, peMsg)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	return nil
}

func (s *Service) detonate(logger log.Logger, vm *VM,
	sha256 string, cfg config.DynFileScanCfg) (agent.FileScanResult, error) {

	ctx := context.Background()

	// Establish a gRPC connection to the agent server running
	// inside the guest.
	client, err := agent.New(vm.IP)
	if err != nil {
		logger.Errorf("failed to establish connection to server: %v", err)
		return agent.FileScanResult{}, nil
	}

	// Deploy the sandbox component files inside the guest.
	ver, err := client.Deploy(ctx, s.cfg.Agent.AgentDestDir, s.sandbox)
	if err != nil {
		return agent.FileScanResult{}, nil
	}
	logger.Infof("sandbox version %s has been deployed", ver)

	src := filepath.Join(s.cfg.SharedVolume, sha256)
	sampleContent, err := utils.ReadAll(src)
	if err != nil {
		return agent.FileScanResult{}, nil
	}

	// Analyze the sample. This call will block until results
	// are ready.
	sandboxCfg := toJSON(cfg)
	res, err := client.Analyze(ctx, sandboxCfg, sampleContent)
	if err != nil {
		return agent.FileScanResult{}, nil
	}

	return res, nil

}

// findFreeVM iterates over the list of available VM and find
// one which is currently not in use.
func (s *Service) findFreeVM() *VM {
	var freeVM *VM
	s.mu.Lock()
	for _, vm := range s.vms {
		if !vm.InUse {
			vm.InUse = true
			freeVM = &vm
			break
		}
	}
	s.mu.Unlock()
	return freeVM
}

// freeVM makes the VM free for consumption.
func (s *Service) freeVM(vm *VM) {
	s.mu.Lock()
	vm.InUse = false
	s.mu.Unlock()
}
