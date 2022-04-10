// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Service provides abstraction to cloud object storage.
type Service struct {
	// s3 service.
	s3svc *awss3.S3
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

	// S3 service client the Upload manager will use.
	s3Svc := awss3.New(sess)

	// Create an uploader with S3 client and custom options
	uploader := s3manager.NewUploaderWithClient(s3Svc, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5MB per part
		u.LeavePartsOnError = true   // Don't delete the parts if the upload fails.
	})

	// Create a downloader with S3 client and custom options
	downloader := s3manager.NewDownloaderWithClient(s3Svc,
		func(u *s3manager.Downloader) {
			u.PartSize = 5 * 1024 * 1024 // 5MB per part
		})

	return Service{s3Svc, uploader, downloader}, nil
}

// Upload upload an object to s3.
func (s Service) Upload(ctx context.Context, bucket, key string,
	file io.Reader) error {

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
func (s Service) Download(ctx context.Context, bucket, key string,
	file io.Writer) error {

	// Download input parameters.
	input := &awss3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Perform the download.
	_, err := s.downloader.DownloadWithContext(ctx, FakeWriterAt{file}, input)

	return err
}

// ListObjects returns the list of objects in a bucket.
func (s Service) ListObjects(ctx context.Context, bucket string,
	fn func(*awss3.ListObjectsV2Output, bool) bool) error {

	input := &awss3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	err := s.s3svc.ListObjectsV2PagesWithContext(ctx, input, fn)
	if err != nil {
		return err
	}
	return nil
}

// Exists checks if an object key exists without returning the object itself.
func (s Service) Exists(ctx context.Context, bucket string, key string) (bool, error) {

	input := &awss3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3svc.HeadObjectWithContext(ctx, input)
	if err != nil {
		return false, err
	}
	return true, nil
}
