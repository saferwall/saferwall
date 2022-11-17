// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	"github.com/saferwall/saferwall/internal/utils"
	"github.com/saferwall/saferwall/internal/vmmanager"
	"github.com/saferwall/saferwall/services"
	"github.com/saferwall/saferwall/services/config"
	pb "github.com/saferwall/saferwall/services/proto"
	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
)

var (
	errNotEnoughResources = errors.New("failed to find a free VM")
)

// Config represents our application config.
type Config struct {
	LogLevel     string             `mapstructure:"log_level"`
	SharedVolume string             `mapstructure:"shared_volume"`
	Agent        AgentCfg           `mapstructure:"agent"`
	VirtMgr      VirtManagerCfg     `mapstructure:"virt_manager"`
	Producer     config.ProducerCfg `mapstructure:"producer"`
	Consumer     config.ConsumerCfg `mapstructure:"consumer"`
}

// AgentCfg represents the guest agent config.
type AgentCfg struct {
	// Destination directory inside the guest where the agent is deployed.
	AgentDestDir string `mapstructure:"dest_dir"`
	// The sandbox binary components.
	PackageName string `mapstructure:"package_name"`
}

// VirtManagerCfg represents the virtualization manager config.
// For now, only libvirt server.
type VirtManagerCfg struct {
	Network      string `mapstructure:"network"`
	Address      string `mapstructure:"address"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	SSHKeyPath   string `mapstructure:"ssh_key_path"`
	SnapshotName string `mapstructure:"snapshot_name"`
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

// generateGuid returns a unique ID to identify a document.
func generateGuid() string {
	id := uuid.New()
	return id.String()
}

// New create a new sandbox service.
func New(cfg Config, logger log.Logger) (*Service, error) {
	var err error
	s := Service{}

	// retrieve the list of active VMs.
	conn, err := vmmanager.New(cfg.VirtMgr.Network, cfg.VirtMgr.Address,
		cfg.VirtMgr.Port, cfg.VirtMgr.User, cfg.VirtMgr.SSHKeyPath)
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
			InUse:     false,
			IsHealthy: true,
			Dom:       d.Dom,
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
	zipPackageData, err := utils.ReadAll(cfg.Agent.PackageName)
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

	ctx := context.Background()
	detonationID := generateGuid()

	// Deserialize the msg sent from the web apis.
	fileScanCfg := config.FileScanCfg{}
	err := json.Unmarshal(m.Body, &fileScanCfg)
	if err != nil {
		s.logger.Errorf("failed un-marshalling json message body: %v", err)
		return err
	}

	sha256 := fileScanCfg.SHA256
	logger := s.logger.With(ctx, "sha256", sha256, "guid", detonationID)
	logger.Info("start processing")

	// Update the state of the job to processing
	status := make(map[string]interface{})
	status["sha256"] = sha256
	status["timestamp"] = time.Now().Unix()

	// Type is only to state to that the document we are storing in the DB is of
	// type `detonate`.
	status["type"] = "dynamic-scan"
	status["status"] = micro.Processing

	payloads := []*pb.Message_Payload{
		{Key: detonationID, Body: toJSON(status), Kind: pb.Message_DBCREATE},
	}

	msg := &pb.Message{Sha256: sha256, Payload: payloads}
	out, err := proto.Marshal(msg)
	if err != nil {
		logger.Errorf("failed to marshal message: %v", err)
		return err
	}
	err = s.pub.Publish(ctx, s.cfg.Producer.Topic, out)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	// Generate a unique ID for this detonation.
	detRes := DetonationResult{}
	detRes.ScanCfg = fileScanCfg.DynFileScanCfg

	// Find a free VM to process this job.
	vm := s.findFreeVM(fileScanCfg.OS)
	if vm == nil {
		logger.Infof("no VM currently available, call 911")
		return errNotEnoughResources
	}
	logger.Infof("VM [%s] with ID: %d was selected", vm.Name, vm.ID)
	logger = logger.With(ctx, "VM", vm.Name)

	// Perform the actual detonation.
	res, errDetonation := s.detonate(logger, vm, fileScanCfg, &detRes)
	if errDetonation != nil {
		logger.Errorf("detonation failed with: %v", errDetonation)
	} else {
		logger.Infof("detonation succeeded")
	}

	// Reverting the VM to a clean state at the end of the analysis
	// is safer than during the start of the analysis, as we instantly
	// stop the malware from running further.
	err = s.vmm.Revert(*vm.Dom, s.cfg.VirtMgr.SnapshotName)
	if err != nil {
		logger.Errorf("failed to revert the VM: %v", err)

		// mark the VM as non healthy so we can repair it.
		logger.Infof("marking the VM as stale")
		s.markStale(vm)

	} else {
		// Free the VM for next job now, then continue on processing
		// sandbox results.
		logger.Infof("freeing the VM")
		s.freeVM(vm)
	}

	// If something went wrong during detonation, we still want to
	// upload the results back to the backend.
	detRes.APITrace = res.TraceLog
	detRes.AgentLog = res.AgentLog
	detRes.SandboxLog = res.SandboxLog
	payloads = []*pb.Message_Payload{
		{Key: detonationID, Body: toJSON(detRes), Kind: pb.Message_DBUPDATE},
	}

	msg = &pb.Message{Sha256: sha256, Payload: payloads}
	out, err = proto.Marshal(msg)
	if err != nil {
		logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	err = s.pub.Publish(ctx, s.cfg.Producer.Topic, out)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	return nil
}
