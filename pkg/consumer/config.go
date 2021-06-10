// Copyright 2021` Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package consumer implements the NSQ worker logic.
package consumer

import (
	"os"

	"github.com/spf13/viper"
)

// BackendCfg represents the backend config.
type BackendCfg struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"admin_user"`
	Password string `mapstructure:"admin_pwd"`
}

// MLCfg represents the consumer config.
type MLCfg struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address"`
}

// NSQCfg represents NSQ config.
type NSQCfg struct {
	Lookupds []string `mapstructure:"lookupds"`
}

// MinioCfg represents minio config.
type MinioCfg struct {
	Enabled   bool   `mapstructure:"enabled"`
	SSL       bool   `mapstructure:"ssl"`
	Endpoint  string `mapstructure:"endpoint"`
	SecKey    string `mapstructure:"seckey"`
	AccessKey string `mapstructure:"accesskey"`
	Spacename string `mapstructure:"spacename"`
}

// AVVendorCfg represents an AV vendor config.
type AVVendorCfg struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address"`
}

// MultiAVCfg represents the multi AV config.
type MultiAVCfg struct {
	Vendors  map[string]AVVendorCfg `mapstructure:"vendors"`
	Enabled bool `mapstructure:"enabled"`
}

// Config represents our application config.
type Config struct {
	LogLevel    string     `mapstructure:"log_level"`
	DownloadDir string     `mapstructure:"download_dir"`
	Headless    bool       `mapstructure:"headless"`
	Backend     BackendCfg `mapstructure:"backend"`
	ML          MLCfg      `mapstructure:"ml"`
	NSQ         NSQCfg     `mapstructure:"nsq"`
	Minio       MinioCfg   `mapstructure:"minio"`
	MultiAV     MultiAVCfg `mapstructure:"multiav"`
}

// loadConfig init our configration.
func loadConfig() (Config, error) {

	c := Config{}
	// Set the path of your config file.
	// In prod, we drop the configs in the same dir as the compiled bin.
	viper.AddConfigPath("configs")
	// In dev environement, the configs is found insie cmd/ dir.
	viper.AddConfigPath("../configs")

	// Load the config type depending on env variable.
	var name string
	env := os.Getenv("SFW_CONSUMER")
	switch env {
	case "dev":
		name = "dev"
	case "prod":
		name = "prod"
	case "test":
		name = "test"
	default:
		name = "dev"
	}

	viper.SetConfigName(name)
	err := viper.ReadInConfig()
	if err != nil {
		return c, err
	}

	err = viper.Unmarshal(&c)
	return c, err
}
