package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadObject(bucket, region, key, filename string) error {
	if bucket == "" || region == "" || key == "" || filename == "" {
		return errors.New("bucket, key and file name required")
	}

	file, err := os.Open(filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(key),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	if err != nil {
		errMsg := fmt.Sprintf("Unable to upload [%q] %q to %q, %v",
			key, filename, bucket, err)
		return errors.New(errMsg)
	}

	fmt.Printf("Successfully uploaded [%q] %q to %q\n", key, filename, bucket)

	return nil
}

func listBucket(region string) {
	// Initialize a session in us-west-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	check(err)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func listObject(bucket, region string, verbose bool) []string {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	check(err)

	// Create S3 service client
	svc := s3.New(sess)

	// Get the list of items
	var objKeys []string
	query := &s3.ListObjectsV2Input{
		Bucket: &bucket,
	}

	for {
		resp, err := svc.ListObjectsV2(query)
		if err != nil {
			exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
		}

		for _, item := range resp.Contents {
			if verbose {
				fmt.Println("Name:         ", *item.Key)
				fmt.Println("Last modified:", *item.LastModified)
				fmt.Println("Size:         ", *item.Size)
				fmt.Println("Storage class:", *item.StorageClass)
				fmt.Println("")
			}

			objKeys = append(objKeys, *item.Key)
		}

		// Fetch the next chunk.
		query.ContinuationToken = resp.NextContinuationToken

		if *resp.IsTruncated == false {
			break
		}
	}

	fmt.Println("Found", len(objKeys), "items in bucket", bucket)
	fmt.Println("")
	return objKeys
}

func isFileFoundInObjStorage(sha256 string) bool {

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create S3 service client
	svc := s3.New(sess)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(sha256)})
	if err != nil {
		log.Printf("svc.HeadObject failed with: %v", err)
		return false
	}

	return true
}

func downloadObject(bucket, region, key string, w io.WriterAt) error {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(w,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	return err
}
