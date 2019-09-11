
// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v
// license that can be found in the LICENSE file.

package utils


import (
	"path"
	"github.com/minio/minio-go/v6"
)


// Download downloads an object from a bucket.
func (client *minio.Client) Download (bucketName, objectName string) ([]byte, error) {
	filePath := path.Join("/"+ bucketName, objectName)
	err := client.FGetObject(bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return []byte, err
	}

	data, err := ReadAll(filePath)
	return data, err
}

