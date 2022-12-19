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
	"github.com/saferwall/saferwall/internal/multiav/trendmicro"
	"github.com/saferwall/saferwall/services/multiav"
)

var flagConfig = flag.String(
	"config", "./../../../../configs/services/multiav/trendmicro",
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

	c := multiav.Config{}

	env := os.Getenv("SAFERWALL_DEPLOYMENT_KIND")

	logger.Infof("loading %s configuration from %s", env, *flagConfig)

	err := config.Load(*flagConfig, env, &c)
	if err != nil {
		return err
	}

	logger.Infof("Starting trendmicro daemon")
	if err = trendmicro.StartDaemon(); err != nil {
		return err
	}

	scanner := trendmicro.Scanner{}
	logger = log.NewCustom(c.LogLevel).With(context.TODO(), "version", constants.Version)
	s, err := multiav.New(c, logger, scanner)
	if err != nil {
		return err
	}

	s.Start()
	return nil
}
