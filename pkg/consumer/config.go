// Copyright 2021` Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

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

// ConsumerCfg represents the consumer config.
type ConsumerCfg struct {
	LogLevel string `mapstructure:"log_level"`
}

// MlCfg represents the consumer config.
type MlCfg struct {
	Address string `mapstructure:"ml"`
}

// NSQCfg represents NSQ config.
type NSQCfg struct {
	Lookupds []string `mapstructure:"lookupds"`
}

// MinioCfg represents minio config.
type MinioCfg struct {
	Endpoint  string `mapstructure:"endpoint"`
	SecKey    string `mapstructure:"seckey"`
	AccessKey string `mapstructure:"accesskey"`
	Spacename string `mapstructure:"spacename"`
	Ssl       bool   `mapstructure:"ssl"`
}

// MultiAvCfg represents the multi AV config.
type MultiAvCfg struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"addr"`
}

// Config represents our application config.
type Config struct {
	Backend  BackendCfg            `mapstructure:"backend"`
	Consumer ConsumerCfg           `mapstructure:"consumer"`
	Ml       MlCfg                 `mapstructure:"ml"`
	Nsq      NSQCfg                `mapstructure:"nsq"`
	Minio    MinioCfg              `mapstructure:"minio"`
	MultiAV  map[string]MultiAvCfg `mapstructure:"multiav"`
}

// loadConfig init our configration.
func loadConfig() (Config, error) {

	c := Config{}

	// Set the path of your config file.
	viper.AddConfigPath("configs")
	viper.AddConfigPath("../configs")

	// Load the config type depending on env variable.
	var name string
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		name = "saferwall.dev"
	case "prod":
		name = "saferwall.prod"
	case "test":
		name = "saferwall.test"
	default:
		name = "saferwall.dev"
	}

	viper.SetConfigName(name)
	err := viper.ReadInConfig()
	if err != nil {
		return c, err
	}

	err = viper.Unmarshal(&c)
	return c, err
}
