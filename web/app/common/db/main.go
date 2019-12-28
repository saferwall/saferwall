// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gocb "gopkg.in/couchbase/gocb.v1"
)

var (
	// UsersBucket controlls users models
	UsersBucket *gocb.Bucket

	// FilesBucket controlls files models
	FilesBucket *gocb.Bucket
)

// Connect to couchbase server
func Connect() {

	/* Init our cluster */
	connectStr := viper.GetString("db.server")
	cluster, err := gocb.Connect(connectStr)
	if err != nil {
		log.Error(err)
	}

	/* Authenticate cluster */
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: username,
		Password: password,
	})

	/* Open the `users` bucket */
	bucketUsers, err := cluster.OpenBucket("users", "")
	if err != nil {
		log.Fatal(err)
	}
	UsersBucket = bucketUsers

	/* Open the `files` bucket */
	bucketFiles, err := cluster.OpenBucket("files", "")
	if err != nil {
		log.Fatal(err)
	}
	FilesBucket = bucketFiles

	/* Create primary indexs */
	err = UsersBucket.Manager("", "").CreatePrimaryIndex("", true, false)
	if err != nil {
		log.Errorf("Failed to create an index over users bucket, reason: %s", err.Error())
	}
	err = FilesBucket.Manager("", "").CreatePrimaryIndex("", true, false)
	if err != nil {
		log.Errorf("Failed to create an index over files bucket, reason: %s", err.Error())
	}

	log.Infoln("Connected to couchbase")
}
