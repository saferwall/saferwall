// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package config

// ProducerCfg represents the producer config.
type ProducerCfg struct {
	Nsqd  string `mapstructure:"nsqd"`
	Topic string `mapstructure:"topic"`
}

// ConsumerCfg represents the consumer config.
type ConsumerCfg struct {
	Lookupds    []string `mapstructure:"lookupds"`
	Topic       string   `mapstructure:"topic"`
	Channel     string   `mapstructure:"channel"`
	Concurrency int      `mapstructure:"concurrency"`
}

// AWSS3Cfg represents AWS S3 credentials.
type AWSS3Cfg struct {
	Region    string `mapstructure:"region"`
	SecretKey string `mapstructure:"secret_key"`
	AccessKey string `mapstructure:"access_key"`
}

// MinioCfg represents Minio credentials.
type MinioCfg struct {
	Endpoint  string `mapstructure:"endpoint"`
	Region    string `mapstructure:"region"`
	SecretKey string `mapstructure:"secret_key"`
	AccessKey string `mapstructure:"access_key"`
}

// LocalFsCfg represents local file system storage data.
type LocalFsCfg struct {
	RootDir string `mapstructure:"root_dir"`
}

// StorageCfg represents the object storage config.
type StorageCfg struct {
	// Deployment kind, possible values: aws, gcp, azure, local.
	DeploymentKind string     `mapstructure:"deployment_kind"`
	Bucket         string     `mapstructure:"bucket"`
	S3             AWSS3Cfg   `mapstructure:"s3"`
	Minio          MinioCfg   `mapstructure:"minio"`
	Local          LocalFsCfg `mapstructure:"local"`
}

// FileScanCfg represents a file scanning config. This map to a 1:1 mapping between
// the config stored in the main saferwall repo.
type FileScanCfg struct {
	// SHA256 hash of the file.
	SHA256 string `json:"sha256,omitempty"`
	// Config used during dynamic file scan.
	DynFileScanCfg `json:"dynamic,omitempty"`
}

// DynFileScanCfg represents the config used to detonate a file.
type DynFileScanCfg struct {
	// Destination path where the sample will be located in the VM.
	DestPath string `json:"dest_path"`
	// Arguments used to run the sample.
	Arguments string `json:"args"`
	// Timeout in seconds for how long to keep the VM running.
	Timeout int `json:"timeout"`
	// Country to route traffic through.
	Country string `json:"country"`
	// Operating System used to run the sample.
	OS string `json:"os"`
}
