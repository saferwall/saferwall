// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"path"

	nsq "github.com/bitly/go-nsq"
	minio "github.com/minio/minio-go"
	"github.com/saferwall/saferwall/pkg/utils"
	"github.com/saferwall/saferwall/web/app/common/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

const (
	addr = "127.0.0.1:4150"
)

var (
	// StoragePath is where we save the samples
	StoragePath string

	// MaxFileSize allowed
	MaxFileSize int64

	// NsqProducer holds an instance of NSQ producer.
	NsqProducer *nsq.Producer

	// DOClient represents an instance of Object Space API client.
	DOClient *minio.Client

	// UserSchema represent a user
	UserSchema *gojsonschema.Schema

	// FileSchema represent a user
	FileSchema *gojsonschema.Schema

	// SamplesSpaceBucket contains the space name of bucket to save samples.
	SamplesSpaceBucket string
)

// loadConfig loads our configration.
func loadConfig() {
	viper.SetConfigName("app")    // no need to include file extension
	viper.AddConfigPath("config") // set the path of your config file

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

// createNSQProducer creates a new NSQ producer.
func createNSQProducer() *nsq.Producer {

	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Create a NewProducer with the name of our topic, the channel, and our config
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}

	log.Info("Got a new NSQ publisher instance")
	return p
}

// initLogging initialize our logging system.
func initLogging() {
	// Add the calling method as a field.
	log.SetReportCaller(false)
}

// loadSchemas will load schemas at server startup.
func loadSchemas() {

	dir, err := utils.Getwd()
	if err != nil {
		log.Error("Failed to GetWd, err: ", err)
	}

	jsonPath := path.Join(dir, "app", "schema", "user.json")
	source := fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader := gojsonschema.NewReferenceLoader(source)
	UserSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Error("Error while loading user schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "file.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	FileSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Error("Error while loading file schema: ", err)
	}
}

// initDOClient returns a client for DigitalOcean Spaces.
func initDOClient() *minio.Client {
	accessKey := viper.GetString("do.accesskey")
	secKey := viper.GetString("do.seckey")
	endpoint := viper.GetString("do.endpoint")
	SamplesSpaceBucket = viper.GetString("do.spacename")
	ssl := true

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		log.Error(err)
	}
	return client
}

// Init will create some directories
func Init() {

	// Init our logger
	initLogging()

	// Load the configuration file
	loadConfig()

	// Load schemas
	loadSchemas()

	// Get an instance of NSQ
	NsqProducer = createNSQProducer()

	DOClient = initDOClient()

	// Connect to the database
	db.Connect()

	StoragePath = viper.GetString("storage.tmp_samples")
	MaxFileSize = int64(viper.GetInt("storage.max_file_size"))
	os.MkdirAll(StoragePath, os.ModePerm)

}
