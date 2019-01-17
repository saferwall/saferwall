// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"path"

	nsq "github.com/bitly/go-nsq"
	"github.com/saferwall/saferwall/web/app/common/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

const (
	addr = "127.0.0.1:4161"
)

var (
	// StoragePath is where we save the samples
	StoragePath string

	// MaxFileSize allowed
	MaxFileSize int64

	// NsqProducer holds an instance of NSQ producer.
	NsqProducer *nsq.Producer

	// UserSchemaLoader represent a user
	UserSchemaLoader gojsonschema.Schema
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

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out := path.Join(dir, "app", "schema", "user.json")
	out = fmt.Sprintf("file:///%s", out)
	schemaLoader := gojsonschema.NewReferenceLoader(out)
	userSchema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Error("Error while loading user schema: ", err)
	}
	UserSchemaLoader = *userSchema
}

// Init will create some directories
func Init() {

	// Init our logger
	initLogging()

	// Load the configuration file
	loadConfig()

	// Connect to the database
	db.Connect()

	// Load schemas
	loadSchemas()

	// Get an instance of NSQ
	NsqProducer = createNSQProducer()

	StoragePath = viper.GetString("storage.tmp_samples")
	MaxFileSize = int64(viper.GetInt("storage.max_file_size"))
	os.MkdirAll(StoragePath, os.ModePerm)

}
