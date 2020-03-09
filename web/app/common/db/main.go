// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package db

import (
	gocb "github.com/couchbase/gocb/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			username,
			password,
		},
	}

	/* Init our cluster */
	server := viper.GetString("db.server")
	cluster, err := gocb.Connect(server, opts)
	if err != nil {
		log.Fatal(err)
	}

	// get a bucket reference over users
	UsersBucket = cluster.Bucket("users", &gocb.BucketOptions{})
	FilesBucket = cluster.Bucket("files", &gocb.BucketOptions{})

	// get a collection reference
	UsersCollection = UsersBucket.DefaultCollection()
	FilesCollection = FilesBucket.DefaultCollection()

	Cluster = cluster

	/* Create primary indexs */
	qm, err := cluster.QueryIndexes()

	err = qm.CreatePrimaryIndex("users", &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})
	if err != nil {
		log.Errorf("Failed to create an index over users bucket, reason: %s", err.Error())
	}

	err = qm.CreatePrimaryIndex("files", &gocb.CreatePrimaryQueryIndexOptions{
		IgnoreIfExists: true,
	})
	if err != nil {
		log.Errorf("Failed to create an index over files bucket, reason: %s", err.Error())
	}

	log.Infoln("Connected to couchbase")
}
