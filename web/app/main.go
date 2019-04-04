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
		log.Fatal(err)
	}

	log.Info("Config was loaded")
}

// createNSQProducer creates a new NSQ producer.
func createNSQProducer() *nsq.Producer {

	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Create a NewProducer with the name of our topic, the channel, and our config
	addr := viper.GetString("nsq.addr")
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		log.Fatal(err)
	}
	if p.Ping() != nil {
		log.Fatal(err)
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
		log.Fatal("Failed to GetWd, err: ", err)
	}

	jsonPath := path.Join(dir, "app", "schema", "user.json")
	source := fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader := gojsonschema.NewReferenceLoader(source)
	UserSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatal("Error while loading user schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "file.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	FileSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatal("Error while loading file schema: ", err)
	}

	log.Info("Schemas were loaded")
}

// initDOClient returns a client for DigitalOcean Spaces.
func initDOClient() *minio.Client {
	accessKey := viper.GetString("do.accesskey")
	secKey := viper.GetString("do.seckey")
	endpoint := viper.GetString("do.endpoint")
	SamplesSpaceBucket = viper.GetString("do.spacename")
	ssl := viper.GetBool("do.ssl")

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Got DO instance")
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

	// Connect to the database
	db.Connect()

	// Get a DO instance
	DOClient = initDOClient()

	StoragePath = viper.GetString("storage.tmp_samples")
	MaxFileSize = int64(viper.GetInt("storage.max_file_size"))
	os.MkdirAll(StoragePath, os.ModePerm)

}
