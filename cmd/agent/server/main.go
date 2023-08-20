// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

//go:build windows

// Package server implements a server for AgentServer service.
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"html/template"
	"io"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"time"

	wapi "github.com/iamacarpet/go-win64api"
	pb "github.com/saferwall/saferwall/internal/agent/proto"
	"github.com/saferwall/saferwall/internal/archiver"
	"github.com/saferwall/saferwall/internal/config"
	"github.com/saferwall/saferwall/internal/constants"
	"github.com/saferwall/saferwall/internal/hasher"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/random"
	"github.com/saferwall/saferwall/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	// grpc library default is 1GB.
	maxMsgSize = 1024 * 1024 * 1024

	// default file scan timeout in seconds.
	defaultFileScanTimeout = 30

	// Hides the window and activates another window.
	SW_HIDE = 0
)

var flagConfig = flag.String("config", "./../../configs/server",
	"path to the config file")

// Config represents our server config.
type Config struct {
	// gRPC server address.
	Address string `mapstructure:"address"`
	// Log level. Defaults to info.
	LogLevel string `mapstructure:"log_level"`
	// Log file where to write logs.
	LogFile string `mapstructure:"log_file"`
	// Hide console window upon startup.
	HideConsoleWindow bool `mapstructure:"hide_console_window"`
	// EnglishWords points to a text file containing a list of english words.
	EnglishWords string `mapstructure:"english_words"`
	// The sandbox config file name to write inside the guest.
	SandboxConfig string `mapstructure:"config_file_name"`
	// File name of the go template file that we use to generate
	// dynamically the sandbox config file.
	TemplateFilename string `mapstructure:"template_file_name"`
	// The file name of the user mode controller component of the sandbox.
	ControllerFilename string `mapstructure:"controller_name"`
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

// Ping probes if the server is healthy and running saferwall analysis VM,
// some information about the guest are returned like OS name, ...
func (s *server) Ping(ctx context.Context, in *emptypb.Empty) (
	*pb.PingReply, error) {

	s.logger.Infof("received a ping request")

	_, os, _, _, _, err := wapi.GetSystemProfile()
	if err != nil {
		s.logger.Error("getting system profile info failed, reason: :%v", err)
		return nil, err
	}

	pingReply := pb.PingReply{Version: constants.Version}

	// Get system information.
	sysInfo, err := json.Marshal(os)
	if err != nil {
		s.logger.Error("marshalling system profile info failed, reason: :%v", err)
		return nil, err
	}

	pingReply.Sysinfo = sysInfo

	return &pingReply, nil
}

// Deploy the sandbox application in the guest.
func (s *server) Deploy(ctx context.Context, in *pb.DeployRequest) (
	*pb.DeployReply, error) {

	s.logger.Infof("received request to deploy package in dest: %s", in.Path)
	s.agentPath = in.Path

	if err := archiver.Unarchive(in.Package, s.agentPath); err != nil {
		s.logger.Error("failed to unarchive package, reason: :%v", err)
		return nil, err
	}

	verFile := filepath.Join(s.agentPath, "VERSION")
	ver, err := utils.ReadAll(verFile)
	if err != nil {
		s.logger.Error("reading sandbox version file failed, reason: :%v", err)
		return nil, err
	}

	return &pb.DeployReply{Version: string(ver)}, nil
}

// Analyze performs the binary analysis.
func (s *server) Analyze(ctx context.Context, in *pb.AnalyzeFileRequest) (
	*pb.AnalyzeFileReply, error) {

	sha256 := s.hasher.Hash(in.Binary)
	logger := s.logger.With(ctx, "sha256", sha256)

	logger.Info("start processing")

	// The config comes from the client as a JSON file.
	// This is explicitly left as map<string>interface{}
	// as we don't want to cause any un-marshalling issues
	// when the scan config changes. All code in the server
	// side has to be carefully written as updating the VMs
	// is expensive.
	var scanCfg map[string]interface{}
	err := json.Unmarshal(in.Config, &scanCfg)
	if err != nil {
		logger.Errorf("failed to unmarshal json config: %v", err)
		return nil, err
	}
	logger.Infof("scan config: %v", scanCfg)

	// Generates the sandbox config from the values from the client.
	tomlConfig, err := s.genSandboxConfig(scanCfg)
	if err != nil {
		return nil, err
	}
	logger.Infof("generated scan config: %v", scanCfg)

	// Write the sandbox TOML config to disk.
	configPath := filepath.Join(s.agentPath, s.cfg.SandboxConfig)
	_, err = utils.WriteBytesFile(configPath, tomlConfig)
	if err != nil {
		logger.Errorf("failed to write config file: %v", err)
		return nil, err
	}

	// Drop the sample to disk.
	sampleData := bytes.NewBuffer(in.Binary)
	samplePath := scanCfg["dest_path"].(string)
	_, err = utils.WriteBytesFile(samplePath, sampleData)
	if err != nil {
		logger.Error("failed to write sample to disk, reason: %v", err)
		return nil, err
	}

	// Add 10 seconds to the timeout to account for bootstrapping the
	// sample execution: loading driver, etc.
	timeout := time.Duration(scanCfg["timeout"].(float64)+10) * time.Second
	deadline := time.Now().Add(timeout)
	logger.Infof("timeout for process to return back: %v", timeout)

	// Create a new context and add a timeout to it.
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	controllerPath := filepath.Join(s.agentPath, s.cfg.ControllerFilename)
	out, err := utils.ExecCmdWithContext(ctx, controllerPath, "-c", configPath)

	// We want to check first the context error to see if the timeout
	// was executed. In any case, we try to collect if they are any
	// artifacts that was created during the analysis.
	if ctx.Err() == context.DeadlineExceeded {
		logger.Errorf("deadline exceeded: %v, output: %s", err, out)
	} else {
		if err != nil {
			logger.Errorf("controller failed: %v, output: %s", err, out)
		} else {
			logger.Info("controller success")
		}
	}

	// At this stage, the controller has terminated and all artifacts are
	// available on disk for collection.
	apiTrace, err := utils.ReadAll(filepath.Join(s.agentPath, "apilog.jsonl"))
	if err != nil {
		logger.Error("failed to read api trace log, reason: %v", err)
	} else {
		logger.Infof("APIs trace logs size is: %d bytes", len(apiTrace))
	}

	// Collect screenshots.
	screenshots := []*pb.AnalyzeFileReply_Screenshot{}
	screenshotsPath := filepath.Join(s.agentPath, "screenshots")
	screenShotId := int32(0)
	err = filepath.Walk(screenshotsPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Errorf("walking screenshots directory failed: %v", err)
				return err
			}
			if !info.IsDir() {
				content, e := utils.ReadAll(path)
				if e != nil {
					logger.Errorf("failed reading screenshot: %s, err: %v", path, err)
				} else {
					screenshots = append(screenshots,
						&pb.AnalyzeFileReply_Screenshot{
							Id: screenShotId, Content: content})
					screenShotId++
				}
			}

			return nil
		})
	if err != nil {
		s.logger.Error("failed to collect screenshots, reason: %v", err)
	} else {
		s.logger.Infof("screenshot collection terminated: %d screenshots acquired", len(screenshots))
	}

	// Collect API buffers.
	apiBuffers := []*pb.AnalyzeFileReply_APIBuffer{}
	apiBuffersPath := filepath.Join(s.agentPath, "api-buffers")
	err = filepath.Walk(apiBuffersPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Errorf("walking api buffers directory failed: %v", err)
				return err
			}

			if !info.IsDir() {
				content, e := utils.ReadAll(path)
				if e != nil {
					logger.Errorf("failed reading api buffer: %s, err: %v", path, err)
				} else {
					apiBuffers = append(apiBuffers,
						&pb.AnalyzeFileReply_APIBuffer{
							Name: info.Name(), Content: content})
				}
			}
			return nil
		})
	if err != nil {
		logger.Error("failed to collect api buffers, reason: %v", err)
	} else {
		logger.Infof("api buffers collection terminated: %d api buffers acquired",
			len(apiBuffers))
	}

	// Collect artifacts.
	artifacts := []*pb.AnalyzeFileReply_Artifact{}
	artifactsPath := filepath.Join(s.agentPath, "artifacts")
	err = filepath.Walk(artifactsPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Errorf("walking artifacts directory failed: %v", err)
				return err
			}

			if !info.IsDir() {
				content, e := utils.ReadAll(path)
				if e != nil {
					logger.Errorf("failed reading artifact: %s, err: %v", path, err)
				} else {
					artifacts = append(artifacts,
						&pb.AnalyzeFileReply_Artifact{Name: info.Name(), Content: content})
				}
			}
			return nil
		})
	if err != nil {
		logger.Error("failed to collect artifacts, reason: %v", err)
	} else {
		logger.Infof("artifacts collection terminated: %d artifact acquired", len(artifacts))
	}

	// Collect the controller logs.
	controllerLog, err := utils.ReadAll(filepath.Join(s.agentPath, "logs", "controller.log"))
	if err != nil {
		logger.Error("failed to read controller log, reason: %v", err)
	} else {
		logger.Infof("controller logs size is: %d bytes", len(controllerLog))
	}

	// Collect the process tree data.
	procTreeLog, err := utils.ReadAll(filepath.Join(s.agentPath, "proc_tree.jsonl"))
	if err != nil {
		logger.Error("failed to read process tree data, reason: %v", err)
	} else {
		logger.Infof("process tree data size is: %d bytes", len(procTreeLog))
	}

	// Keep the agent log the last thing to read.
	agentLog, err := utils.ReadAll(filepath.Join(s.agentPath, s.cfg.LogFile))
	if err != nil {
		logger.Error("failed to read agent log, reason: %v", err)
	} else {
		logger.Infof("agent logs size is: %d bytes", len(agentLog))
	}

	return &pb.AnalyzeFileReply{
		APITrace:      apiTrace,
		APIBuffers:    apiBuffers,
		Screenshots:   screenshots,
		Artifacts:     artifacts,
		ServerLog:     agentLog,
		ControllerLog: controllerLog,
		ProcessTree:   procTreeLog,
	}, nil
}

// Build sandbox config by taking the values from from the `FileScanCfg` and
// generates a TOML config on the fly using go templates.
func (s *server) genSandboxConfig(scanCfg map[string]interface{}) (
	io.Reader, error) {

	_, ok := scanCfg["timeout"]
	if !ok || scanCfg["timeout"] == 0 {
		scanCfg["timeout"] = defaultFileScanTimeout
	}

	_, ok = scanCfg["dest_path"]
	if !ok || scanCfg["dest_path"] == "" {
		randomFilename := s.randomizer.Random()
		scanCfg["dest_path"] = "%USERPROFILE%//Downloads//" + randomFilename + ".exe"
	}

	// For path expansion to work in Windows, we need to replace the
	// `%` with `$`.
	scanCfg["dest_path"] = utils.Resolve(scanCfg["dest_path"].(string))

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

// Hide the current window console.
func hideConsoleWindow(logger log.Logger) error {
	getConsoleWindow := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	hWindow, _, _ := getConsoleWindow.Call()

	// NULL if there is no such associated console.
	if hWindow == 0 {
		return errors.New("no associated console with the current process")
	}

	showWindow := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	st, _, _ := showWindow.Call(hWindow, SW_HIDE)

	if st != 0 {
		// If the window was previously visible, the return value is nonzero.
		logger.Info("window has be hidden")
	} else {
		// 	If the window was previously hidden, the return value is zero.
		logger.Info("window was previously hidden")
	}

	return nil
}

// DefaultServerOpts returns the set of default grpc ServerOption's that Agent
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

	logger := log.New().With(context.TODO(), "version", constants.Version)

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

	// update the logger according to the config.
	logger = log.NewCustomWithFile(c.LogLevel, c.LogFile).With(context.TODO(), "version", constants.Version)

	// the console window should be hidden in prod.
	if c.HideConsoleWindow {
		logger.Info("hiding console window")
		err := hideConsoleWindow(logger)
		if err != nil {
			return err
		}
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
