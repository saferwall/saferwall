// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package db

import (
	gocb "github.com/couchbase/gocb/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	// UsersCollection reference the users collection
	UsersCollection *gocb.Collection

	// FilesCollection reference the files collection
	FilesCollection *gocb.Collection

	// UsersBucket represents the users bucket.
	UsersBucket *gocb.Bucket

	// FilesBucket represents the files bucket.
	FilesBucket *gocb.Bucket

	// Cluster is our cluster
	Cluster *gocb.Cluster
)

// Connect to couchbase server
func Connect() {

	/* setup logger */
	// gocb.SetLogger(gocb.DefaultStdioLogger())

	/* Authenticate cluster */
	cbUsername := viper.GetString("db.username")
	cbPassword := viper.GetString("db.password")
	opts := gocb.ClusterOptions{
		Username: cbUsername,
		Password: cbPassword,
	}

	// Init our cluster.
	server := viper.GetString("db.server")
	cluster, err := gocb.Connect(server, opts)
	if err != nil {
		log.Fatal(err)
	}

	// Get a bucket reference over users.
	UsersBucket = cluster.Bucket("users")
	FilesBucket = cluster.Bucket("files")

	// We wait until the bucket is definitely connected and setup.
	err = UsersBucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Errorf("Failed to WaitUntilReady for users bucket, reason: %v", err)
	}
	err = FilesBucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Errorf("Failed to WaitUntilReady for files bucket, reason: %v", err)
	}

	// Get a collection reference.
	UsersCollection = UsersBucket.DefaultCollection()
	FilesCollection = FilesBucket.DefaultCollection()

	Cluster = cluster

	// Create primary indexs.
	mgr := cluster.QueryIndexes()
	err = mgr.CreatePrimaryIndex("users", &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})
	if err != nil {
		log.Errorf("Failed to create an index over users bucket, reason: %v", err)
	}

	err = mgr.CreatePrimaryIndex("files", &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})
	if err != nil {
		log.Errorf("Failed to create an index over files bucket, reason: %v", err)
	}

	// Create secondary indexes.
	err = mgr.CreateIndex("users", "idx_username", []string{"username"},
		&gocb.CreateQueryIndexOptions{
			IgnoreIfExists: true})
	if err != nil {
		log.Errorf("Failed to create secondary index (idx_username) over users bucket, reason: %v", err)
	}
	err = mgr.CreateIndex("users", "idx_email", []string{"email"},
		&gocb.CreateQueryIndexOptions{
			IgnoreIfExists: true})
	if err != nil {
		log.Errorf("Failed to create secondary index (idx_email) over users bucket, reason: %v", err)
	}
	err = mgr.CreateIndex("files", "idx_sha256", []string{"sha256"},
		&gocb.CreateQueryIndexOptions{
			IgnoreIfExists: true})
	if err != nil {
		log.Errorf("Failed to create secondary index (idx_email) over users bucket, reason: %v", err)
	}

	err = mgr.CreateIndex("files", "idx_status", []string{"status"},
		&gocb.CreateQueryIndexOptions{
			IgnoreIfExists: true})
	if err != nil {
		log.Errorf("Failed to create secondary index (idx_status) over users bucket, reason: %v", err)
	}

	log.Infoln("Connected to couchbase")
}
