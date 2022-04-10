// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package consumer implements the NSQ worker logic.
package config

import (
	"github.com/spf13/viper"
)

// Load returns an application configuration which is populated
// from the given configuration file.
func Load(path, env string, c interface{}) error {

	// Adding our TOML config file.
	viper.AddConfigPath(path)

	// Load the config type depending on env variable.
	var name string
	switch env {
	case "local":
		name = "local"
	case "dev":
		name = "dev"
	case "prod":
		name = "prod"
	default:
		name = "local"
	}

	// Set the config name to choose from the config path
	// Extension not needed.
	viper.SetConfigName(name)

	// Load the configuration from disk.
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Unmarshals the config into our interface.
	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	return err
}
