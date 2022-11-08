// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package local

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Service provides abstraction to cloud object storage.
type Service struct {
	// Root directory in the local file system.
	root string
}

// New generates new object storage service.
func New(root string) (Service, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := os.MkdirAll(root, os.ModePerm); err != nil {
			return Service{}, err
		}
	}
	return Service{root}, nil
}

// Upload upload an object to s3.
func (s Service) Upload(ctx context.Context, bucket, key string,
	file io.Reader) error {

	dest := filepath.Join(s.root, bucket, key)

	// Create new file.
	new, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer new.Close()

	// Perform the copy.
	if _, err := io.Copy(new, file); err != nil {
		return err
	}

	return nil
}

// Download downloads an object from the local file system.
func (s Service) Download(ctx context.Context, bucket, key string,
	dst io.Writer) error {

	// Create new file.
	name := filepath.Join(s.root, bucket, key)
	src, err := os.Open(name)
	if err != nil {
		return err
	}
	defer src.Close()

	// Perform the copy.
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// MakeBucket creates a new folder in the local file system that acts like
// a bucket or a container in a object storage.
func (s Service) MakeBucket(ctx context.Context, bucketName, location string) error {
	dest := filepath.Join(s.root, bucketName)
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// Exists checks whether a file exists in disk.
func (s Service) Exists(ctx context.Context, bucketName, key string) (bool, error) {
	name := filepath.Join(s.root, bucketName, key)
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
