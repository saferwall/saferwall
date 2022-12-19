// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"

	"github.com/saferwall/saferwall/internal/config"
	"github.com/saferwall/saferwall/internal/constants"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/services/aggregator"
)

var flagConfig = flag.String(
	"config", "./../../../configs/services/aggregator",
	"path to the config file")

func main() {

	flag.Parse()

	// Create root logger tagged with server version.
	logger := log.New().With(context.TODO(), "version", constants.Version)
	if err := run(logger); err != nil {
		logger.Errorf("failed to run the server: %s", err)
		os.Exit(-1)
	}
}

func run(logger log.Logger) error {

	c := aggregator.Config{}

	env := os.Getenv("SAFERWALL_DEPLOYMENT_KIND")

	logger.Infof("loading %s configuration from %s", env, *flagConfig)

	err := config.Load(*flagConfig, env, &c)
	if err != nil {
		return err
	}

	logger = log.NewCustom(c.LogLevel).With(context.TODO(), "version", constants.Version)
	s, err := aggregator.New(c, logger)
	if err != nil {
		return err
	}

	s.Start()
	return nil
}
