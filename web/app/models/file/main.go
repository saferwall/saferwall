// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package file

import (
	"fmt"
	"log"
	"time"

	"github.com/saferwall/saferwall/web/app/common/db"
)

// AV vendor
type AV struct {
	Vendor string `json:"vendor,omitempty"`
}

// File represent a sample
type File struct {
	Md5       string    `json:"md5,omitempty"`
	Sha1      string    `json:"sha1,omitempty"`
	Sha256    string    `json:"sha256,omitempty"`
	Sha512    string    `json:"sha512,omitempty"`
	Ssdeep    string    `json:"ssdeep,omitempty"`
	Crc32     string    `json:"crc32,omitempty"`
	Magic     string    `json:"magic,omitempty"`
	Size      int       `json:"size,omitempty"`
	FirstSeen time.Time `json:"first_seen,omitempty"`
}

// Create creates a new file
func (file *File) Create(documentID string) error {
	_, error := db.FilesBucket.Upsert(documentID, file, 0)
	if error != nil {
		log.Fatal(error)
		return error
	}

	return nil
}

// GetFileBySHA256 return user document
func GetFileBySHA256(sha256 string) (File, error) {

	// get our file
	file := File{}
	cas, err := db.FilesBucket.Get(sha256, &file)
	if err != nil {
		fmt.Println(err, cas)
		return file, err
	}

	return file, err
}

// DeleteAllFiles will empty files bucket
func DeleteAllFiles() {
	// Keep in mind that you must have flushing enabled in the buckets configuration.
	db.FilesBucket.Manager("", "").Flush()
}
