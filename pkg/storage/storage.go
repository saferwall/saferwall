// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package storage

import (
	"errors"
	"io"

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

// Options for object storage.
type Options struct {
	S3Region  string
	S3AccKey  string
	S3SecKey  string
	LocalRootDir string
}

func New(kind string, opts Options) (Storage, error) {

	switch kind {
	case "aws":
		return s3.New(opts.S3Region, opts.S3AccKey, opts.S3SecKey)
	case "local":
		return local.New(opts.LocalRootDir)
	}

	return nil, errDeploymentNotFound
}
