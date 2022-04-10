// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package server implements a server for AgentServer service.
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	pb "github.com/saferwall/saferwall/internal/agent/proto"
	"github.com/saferwall/saferwall/internal/archiver"
	"github.com/saferwall/saferwall/internal/config"
	"github.com/saferwall/saferwall/internal/hasher"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/random"
	"github.com/saferwall/saferwall/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

const (
	// grpc library default is 64MB.
	maxMsgSize = 1024 * 1024 * 64
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./../../configs/server",
	"path to the config file")

// Config represents our server config.
type Config struct {
	// gRPC server address.
	Address string `mapstructure:"address"`
	// EnglishWords points to a text file contaning a list of english words.
	EnglishWords string `mapstructure:"english_words"`
	// The sandbox config file name to write inside the guest.
	SandboxConfig string `mapstructure:"config_file_name"`
	// File name of the go template file that we use to generate
	// dynamically the sandbox config file.
	TemplateFilename string `mapstructure:"template_file_name"`
	// The file name of the user mode controller component of the sandbox.
	ControllerFilename string `mapstructure:"controller_name"`
}

// ScanConfig represents the dynamic malware analysis configuration.
type ScanConfig struct {
	// Destination path where the sample will be located in the VM.
	SampleDestPath string `json:"sample_dest_path"`
	// Arguments used to run the sample.
	Arguments string `json:"arguments"`
	// Timeout in seconds for how long to keep the VM running.
	Timeout int `json:"timeout"`
}

// server is used to implement agent.GreeterServer.
type server struct {
	pb.UnimplementedAgentServer
	logger     log.Logger
	hasher     hasher.Hasher
	randomizer random.Ramdomizer
	cfg        Config
	agentPath  string
}

// Deploy deloys the sandbox application in the guest.
func (s *server) Deploy(ctx context.Context, in *pb.DeployRequest) (
	*pb.DeployReply, error) {

	s.logger.Infof("Received request to deploy package in dest: %s", in.Path)
	s.agentPath = in.Path

	if err := archiver.Unarchive(in.Package, s.agentPath); err != nil {
		s.logger.Error("failed to unarchive package, reason: :%v", err)
		return nil, err
	}

	verfile := filepath.Join(s.agentPath, "VERSION.txt")
	ver, err := utils.ReadAll(verfile)

	return &pb.DeployReply{Version: string(ver)}, err
}

// Analyze performs the binary analysis.
func (s *server) Analyze(ctx context.Context, in *pb.AnalyzeFileRequest) (
	*pb.AnalyzeFileReply, error) {

	sha256 := s.hasher.Hash(in.Binary)
	logger := s.logger.With(ctx, "sha256", sha256)

	logger.Info("start processing")

	// The config comes from the client as a JSON file.
	var scanCfg ScanConfig
	err := json.Unmarshal(in.Config, &scanCfg)
	if err != nil {
		s.logger.Errorf("failed to unmarshal json config: %v", err)
		return nil, err
	}

	// Generates the sandbox config from the values from the client.
	tomlConfig, err := s.genSandboxConfig(&scanCfg)
	if err != nil {
		return nil, err
	}

	// Write the sandbox TOML config to disk.
	configPath := filepath.Join(s.agentPath, s.cfg.SandboxConfig)
	_, err = utils.WriteBytesFile(configPath, tomlConfig)
	if err != nil {
		s.logger.Errorf("failed to write config file: %v", err)
		return nil, err
	}

	// Drop the sample to disk.
	sampleData := bytes.NewBuffer(in.Binary)
	_, err = utils.WriteBytesFile(scanCfg.SampleDestPath, sampleData)
	if err != nil {
		s.logger.Error("failed to write sample to disk, reason: :%v", err)
		return nil, err
	}

	// Create a new context and add a timeout to it.
	// It is possible that the controller crashes and never
	// returns back, this timeout helps remediating this issue.
	timeout := (time.Duration(scanCfg.Timeout + 5)) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	controllerPath := filepath.Join(s.agentPath, s.cfg.ControllerFilename)
	out, err := utils.ExecCmdWithContext(ctx, controllerPath, "-c", configPath)

	// We want to check first the context error to see if the timeout
	// was executed. In any case, we try to collect if they are any
	// artifacts that was created during the analysis.
	if ctx.Err() == context.DeadlineExceeded {
		s.logger.Errorf("deadline exceeded: %v, output: %s", err, out)
	} else {
		if err != nil {
			s.logger.Errorf("controller failed: %v, output: %s", err, out)
		} else {
			s.logger.Info("controller success")

		}
	}

	// At this stage, the controller has terminated and all artifacts are
	// available on disk for collection.
	apiTrace, err := utils.ReadAll(filepath.Join(s.agentPath, "apilog.jsonl"))
	if err != nil {
		s.logger.Error("failed to read api trace log, reason: :%v", err)
		return nil, err
	}

	return &pb.AnalyzeFileReply{
		Apitrace:           apiTrace,
		Screenshots:        nil,
		CollectedArtifacts: nil,
		Memdumps:           nil,
		Serverlog:          nil,
		Controllerlog:      nil,
	}, nil
}

// Build sandbox config by taking the values from from the
// `FileScanCfg` and generates a TOML config on the fly using
// go templates.
func (s *server) genSandboxConfig(scanCfg *ScanConfig) (
	io.Reader, error) {

	if scanCfg.Timeout == 0 {
		scanCfg.Timeout = 60
	}
	if scanCfg.SampleDestPath == "" {
		randomFilename := s.randomizer.Random()
		scanCfg.SampleDestPath = "%USERPROFILE%//Downloads//" + randomFilename + ".exe"
	}

	// For path expansion to work in Windows, we need to replace the
	// `%` with `$`.
	scanCfg.SampleDestPath = utils.Resolve(scanCfg.SampleDestPath)

	configTemplate := filepath.Join(s.agentPath, s.cfg.TemplateFilename)
	tpl, err := template.ParseFiles(configTemplate)
	if err != nil {
		s.logger.Errorf("failed to parse template file: %v", err)
		return nil, err
	}

	tomlConfig := new(bytes.Buffer)
	if err = tpl.Execute(tomlConfig, scanCfg); err != nil {
		s.logger.Errorf("failed to execute template: %v", err)
		return nil, err
	}

	return tomlConfig, nil
}

// DefaultServerOpts returns the set of default grpc ServerOption's that Tiller
// requires.
func DefaultServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	}
}

// NewServer creates a new grpc server.
func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(DefaultServerOpts(), opts...)...)
}

// Serve registers reflection service on gRPC server and start serving.
func Serve(s *grpc.Server, lis net.Listener) error {
	reflection.Register(s)
	return s.Serve(lis)
}

func main() {

	flag.Parse()

	logger := log.New().With(context.TODO(), "version", Version)

	if err := run(logger, *flagConfig); err != nil {
		logger.Errorf("failed to run the server: %s", err)
		os.Exit(-1)
	}

}

func run(logger log.Logger, configFile string) error {

	c := Config{}

	env := os.Getenv("SAFERWALL_DEPLOYMENT_KIND")

	logger.Infof("loading %s configuration from %s", env, *flagConfig)

	err := config.Load(configFile, env, &c)
	if err != nil {
		return err
	}

	// create a hasher.
	hashsvc := hasher.New(sha256.New())

	// create a string randomizer.
	randomSvc, err := random.New(c.EnglishWords)
	if err != nil {
		return err
	}

	// create a tcp listener.
	lis, err := net.Listen("tcp", c.Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object and attach it to the agent service.
	s := NewServer()
	pb.RegisterAgentServer(s, &server{
		logger:     logger,
		hasher:     hashsvc,
		randomizer: randomSvc,
		cfg:        c,
	})

	// register reflection service on gRPC server and serve.
	logger.Infof("starting server on: %s", c.Address)
	err = Serve(s, lis)
	if err != nil {
		grpclog.Fatalf("failed to serve: %v", err)
	}

	return nil
}
