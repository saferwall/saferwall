// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package storage

import (
	"context"
	"errors"
	"io"

	"github.com/saferwall/saferwall/pkg/storage/local"
	"github.com/saferwall/saferwall/pkg/storage/minio"
	"github.com/saferwall/saferwall/pkg/storage/s3"
)

var (
	errDeploymentNotFound = errors.New("deployment not found")
)

// Storage abstract uploading and download files from different
// object storage solutions.
type Storage interface {
	// Upload uploads a file to an object storage.
	Upload(ctx context.Context, bucket, key string, file io.Reader) error
	// Download downloads a file from a remote object storage location.
	Download(ctx context.Context, bucket, key string, file io.Writer) error
}

// Options for object storage.
type Options struct {
	Region        string
	AccessKey     string
	SecretKey     string
	MinioEndpoint string
	LocalRootDir  string
}

func New(kind string, opts Options) (Storage, error) {
	switch kind {
	case "aws":
		return s3.New(opts.Region, opts.AccessKey, opts.SecretKey)
	case "minio":
		return minio.New(opts.MinioEndpoint, opts.AccessKey, opts.SecretKey)
	case "local":
		return local.New(opts.LocalRootDir)
	}

	return nil, errDeploymentNotFound
}
