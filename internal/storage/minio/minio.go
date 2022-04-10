// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package minio

import (
	"bytes"
	"context"
	"io"

	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Service provides abstraction to cloud object storage.
type Service struct {
	// s3 client.
	client *miniogo.Client
}

func New(endpoint, accessKey, secretKey string) (Service, error) {

	// New returns an Amazon S3 compatible client object.
	// API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3Client, err := miniogo.New(endpoint, &miniogo.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return Service{}, nil
	}

	return Service{s3Client}, nil
}

// Upload uploads an object to s3.
func (s Service) Upload(ctx context.Context, bucket, key string,
	file io.Reader) error {

	// Get the size.
	buf := &bytes.Buffer{}
	size, err := io.Copy(buf, file)
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(ctx, bucket, key, buf, size,
		miniogo.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Download(ctx context.Context, bucket, key string,
	file io.Writer) error {

	reader, err := s.client.GetObject(ctx, bucket, key, miniogo.GetObjectOptions{})
	if err != nil {
		return err
	}

	stat, err := reader.Stat()
	if err != nil {
		return err
	}

	_, err = io.CopyN(file, reader, stat.Size)
	if err != nil {
		return err
	}
	return nil
}
