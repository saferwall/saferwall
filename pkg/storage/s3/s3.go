// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package s3

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Service provides abstraction to cloud object storage.
type Service struct {
	// s3 uploader.
	uploader *s3manager.Uploader
	// s3 downloader.
	downloader *s3manager.Downloader
}

// FakeWriterAt represents a struct that provides the method WriteAt so it will
// satisfy interface io.WriterAt. It will ignore offset and therefore works
// like just io.Writer. AWS SDK is Using io.WriterAt because of concurrent
// download, so it can write at offset position (e.g. in middle of file).
// By disabling concurrent download we can safely ignore the offset argument
// because it will be downloaded sequentially.
type FakeWriterAt struct {
	w io.Writer
}

func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	return fw.w.Write(p)
}

// New generates new s3 object storage service.
func New(region, accessKey, secretKey string) (Service, error) {

	// The session the S3 Uploader will use.
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region),
		Credentials: creds})
	if err != nil {
		return Service{}, nil
	}

	// S3 service client the Upload/Download manager will use.
	s3Svc := awss3.New(sess)

	// Create an uploader with S3 client and custom options.
	uploader := s3manager.NewUploaderWithClient(s3Svc,
		func(u *s3manager.Uploader) {
			u.PartSize = 5 * 1024 * 1024 // 5MB per part
			u.LeavePartsOnError = true   // Don't delete the parts if the upload fails.
		})

	// Create a downloader with S3 client and custom options
	downloader := s3manager.NewDownloaderWithClient(s3Svc,
		func(u *s3manager.Downloader) {
			u.PartSize = 5 * 1024 * 1024 // 5MB per part
		})

	return Service{uploader, downloader}, nil
}

// Upload uploads an object to s3.
func (s Service) Upload(bucket, key string, file io.Reader, timeout int) error {

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, time.Duration(timeout))
	}

	// Ensure the context is canceled to prevent leaking.
	// See context package for more information, https://golang.org/pkg/context/
	if cancelFn != nil {
		defer cancelFn()
	}

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	}

	// Perform an upload.
	_, err := s.uploader.UploadWithContext(ctx, upParams)

	return err
}

// Download downloads an object from s3.
func (s Service) Download(bucket, key string, file io.Writer, timeout int) error {

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, time.Duration(timeout))
	}

	// Ensure the context is canceled to prevent leaking.
	// See context package for more information, https://golang.org/pkg/context/
	if cancelFn != nil {
		defer cancelFn()
	}

	// Download input parameters.
	input := &awss3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Perform the download.
	_, err := s.downloader.DownloadWithContext(ctx, FakeWriterAt{file}, input)

	return err
}
