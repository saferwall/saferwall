// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"time"

	client "github.com/saferwall/saferwall/internal/agent"
	"github.com/saferwall/saferwall/internal/config"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/utils"
)

// Config represents our client config.
type Config struct {
	// gRPC server address.
	Address string `mapstructure:"address"`
	// Location inside the guest where to deploy the package.zip
	AgentDestPath string `mapstructure:"agent_dest_path"`
	// Local packages.zip location.
	AgentSrcPath string `mapstructure:"agent_src_path"`
	// Sample to analyze file path.
	SampleSrcPath string `mapstructure:"sample_src_path"`
	// Maximum timeout in seconds for the client to wait for the server
	// to deply back during an alaysis request before it hang up.
	Timeout int `mapstructure:"timeout"`
}

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./../../configs/client",
	"path to the config file")

func main() {

	flag.Parse()

	logger := log.New().With(context.TODO(), "version", Version)

	if err := run(logger); err != nil {
		logger.Errorf("failed to run the server: %s", err)
		os.Exit(-1)
	}
}

func run(logger log.Logger) error {

	c := Config{}

	env := os.Getenv("SAFERWALL_DEPLOYMENT_KIND")

	logger.Infof("loading %s configuration from %s", env, *flagConfig)

	err := config.Load(*flagConfig, env, &c)
	if err != nil {
		return err
	}

	zipPackage, err := utils.ReadAll(c.AgentSrcPath)
	if err != nil {
		logger.Errorf("could not greet: %v", err)
		return err
	}

	logger.Infof("creating new client conn to: %s", c.Address)
	ac, err := client.New(c.Address)
	if err != nil {
		logger.Errorf("could not create grpc client: %v", err)
		return err
	}

	// Deploy saferwall package request.
	ctx, cancelDeploy := context.WithTimeout(context.Background(), time.Second)
	defer cancelDeploy()
	ver, err := ac.Deploy(ctx, c.AgentDestPath, zipPackage)
	if err != nil {
		logger.Errorf("could not deploy sandbox: %v", err)
		return err
	}
	logger.Infof("Version: %s", ver)

	// Analyze sample binary request.
	sampleData, err := utils.ReadAll(c.SampleSrcPath)
	if err != nil {
		logger.Errorf("could not read src sample: %v", err)
		return err
	}

	fileScanCfg, _ := json.Marshal(map[string]interface{}{
		"timeout": 5,
	})

	// Setting a hard timeout on the client side in case
	// the server never replied back.
	timeout := time.Duration(c.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ac.Analyze(ctx, fileScanCfg, sampleData)
	return nil
}
