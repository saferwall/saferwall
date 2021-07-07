// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package storage

import (
	"errors"
	"io"

	"github.com/saferwall/saferwall/pkg/config"
	"github.com/saferwall/saferwall/pkg/storage/local"
	"github.com/saferwall/saferwall/pkg/storage/s3"
)

var (
	errDeploymentNotFound = errors.New("deployment not found")
)

// Storage abstract accessing object storage from various cloud locations.
type Storage interface {
	// Download uploads a file to an object storage.
	Download(bucket, key string, file io.Writer, timeout int) error
	// Upload uploads a file to an object storage.
	Upload(bucket, key string, file io.Reader, timeout int) error
}

func New(deployment string, cfg config.StorageCfg) (Storage, error) {

	switch deployment {
	case "aws":
		return s3.New(cfg.S3.Region, cfg.S3.AccessKey, cfg.S3.SecretKey)
	case "local":
		return local.New(cfg.Local.RootDir)
	}

	return nil, errDeploymentNotFound
}
