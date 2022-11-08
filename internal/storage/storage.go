// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package storage

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/saferwall/saferwall/internal/storage/local"
	"github.com/saferwall/saferwall/internal/storage/minio"
	"github.com/saferwall/saferwall/internal/storage/s3"
)

var (
	errDeploymentNotFound = errors.New("deployment not found")
	timeout               = time.Duration(time.Second * 5)
)

// Storage abstract uploading and download files from different
// object storage solutions.
type Storage interface {
	// Upload uploads a file to an object storage.
	Upload(ctx context.Context, bucket, key string, file io.Reader) error
	// Download downloads a file from a remote object storage location.
	Download(ctx context.Context, bucket, key string, file io.Writer) error
	// MakeBucket creates a new bucket.
	MakeBucket(ctx context.Context, bucket, location string) error
	// Exists checks whether an object exists.
	Exists(ctx context.Context, bucket, key string) (bool, error)
}

// Options for object storage.
type Options struct {
	Region        string
	AccessKey     string
	SecretKey     string
	MinioEndpoint string
	LocalRootDir  string
	Bucket        string
}

func New(kind string, opts Options) (Storage, error) {

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	// Ensure the context is canceled to prevent leaking.
	// See context package for more information, https://golang.org/pkg/context/
	if cancelFn != nil {
		defer cancelFn()
	}

	switch kind {
	case "aws":
		svc, err := s3.New(opts.Region, opts.AccessKey, opts.SecretKey)
		if err != nil {
			return nil, err
		}
		err = svc.MakeBucket(ctx, opts.Bucket, opts.Region)
		if err != nil {
			return nil, err
		}
		return svc, nil

	case "minio":
		svc, err := minio.New(opts.MinioEndpoint, opts.AccessKey, opts.SecretKey)
		if err != nil {
			return nil, err
		}
		err = svc.MakeBucket(ctx, opts.Bucket, opts.Region)
		if err != nil {
			return nil, err
		}
		return svc, nil
	case "local":
		svc, err := local.New(opts.LocalRootDir)
		if err != nil {
			return nil, err
		}
		err = svc.MakeBucket(ctx, opts.Bucket, "")
		if err != nil {
			return nil, err
		}
		return svc, nil
	}

	return nil, errDeploymentNotFound
}
