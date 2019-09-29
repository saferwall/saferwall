// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	nsq "github.com/bitly/go-nsq"
	"github.com/minio/minio-go/v6"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/packer"
	s "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/pkg/trid"
	"github.com/saferwall/saferwall/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path"
	"syscall"
	"github.com/spf13/viper"

)

var (
	minioClient     *minio.Client
	backendEndpoint string
)

type stringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}
type result struct {
	Crc32   string                 `json:"crc32,omitempty"`
	Md5     string                 `json:"md5,omitempty"`
	Sha1    string                 `json:"sha1,omitempty"`
	Sha256  string                 `json:"sha256,omitempty"`
	Sha512  string                 `json:"sha512,omitempty"`
	Ssdeep  string                 `json:"ssdeep,omitempty"`
	Exif    map[string]string      `json:"exif,omitempty"`
	TriD    []string               `json:"trid,omitempty"`
	Packer  []string               `json:"packer,omitempty"`
	Magic   string                 `json:"magic,omitempty"`
	Strings []stringStruct         `json:"strings,omitempty"`
	MultiAV map[string]interface{} `json:"multiav,omitempty"`
	Status  int                    `json:"status,omitempty"`
}

// File scan progress status.
const (
	queued     = iota
	processing = iota
	finished   = iota
)

// NoopNSQLogger allows us to pipe NSQ logs to dev/null
// The default NSQ logger is great for debugging, but did
// not fit our normally well structured JSON logs. Luckily
// NSQ provides a simple interface for injecting your own
// logger.
type NoopNSQLogger struct{}

// Output allows us to implement the nsq.Logger interface
func (l *NoopNSQLogger) Output(calldepth int, s string) error {
	log.Info(s)
	return nil
}

// MessageHandler adheres to the nsq.Handler interface.
// This allows us to define our own custome handlers for
// our messages. Think of these handlers much like you would
// an http handler.
type MessageHandler struct{}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface. This where you'll write your message
// handling logic.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// returning an error results in the message being re-enqueued
		// a REQ is sent to nsqd
		return errors.New("body is blank re-enqueue message")
	}

	sha256 := string(m.Body)
	log.Infof("Processing %s", sha256)

	// set the file status to `processing`
	res := result{}
	res.Status = processing

	// Marshell results
	buff, err := json.Marshal(res)
	if err != nil {
		log.Error("Failed to get object, err: ", err)
		return err
	}

	// Update document
	updateDocument(sha256, buff)

	// Download the sample
	bucketName := viper.GetString("do.spacename")
	b, err := utils.Download(minioClient, bucketName, sha256)
	if err != nil {
		log.Errorf("Failed to download file %s", sha256)
		return err
	}

	filePath := path.Join("/samples", sha256)

	// Run crypto pkg
	r := crypto.HashBytes(b)
	res.Crc32 = r.Crc32
	res.Md5 = r.Md5
	res.Sha1 = r.Sha1
	res.Sha256 = r.Sha256
	res.Sha512 = r.Sha512
	res.Ssdeep = r.Ssdeep
	log.Infof("HashBytes success %s", sha256)

	// Run exiftool pkg
	res.Exif, err = exiftool.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with exiftool, err: ", err)
	}
	log.Infof("exiftool success %s", sha256)

	// Run TRiD pkg
	res.TriD, err = trid.Scan(filePath)
	if err != nil {
		log.Error("Faileds to scan file with trid, err: ", err)
	}
	log.Infof("trid success %s", sha256)

	// Run Magic Pkg
	res.Magic, err = magic.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with magic, err: ", err)
	}
	log.Infof("magic extraction success %s", sha256)

	// Run Die Pkg
	res.Packer, err = packer.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with packer, err: ", err)
	}
	log.Infof("packer extraction success %s", sha256)

	// Run strings pkg
	n := 10
	asciiStrings := s.GetASCIIStrings(b, n)
	wideStrings := s.GetUnicodeStrings(b, n)
	asmStrings := s.GetAsmStrings(b)

	// Remove duplicates
	uniqueASCII := utils.UniqueSlice(asciiStrings)
	uniqueWide := utils.UniqueSlice(wideStrings)
	uniqueAsm := utils.UniqueSlice(asmStrings)

	var strResults []stringStruct
	for _, str := range uniqueASCII {
		strResults = append(strResults, stringStruct{"ascii", str})
	}

	for _, str := range uniqueWide {
		strResults = append(strResults, stringStruct{"wide", str})
	}

	for _, str := range uniqueAsm {
		strResults = append(strResults, stringStruct{"asm", str})
	}
	res.Strings = strResults
	log.Infof("strings success %s", sha256)

	// multiav scanning
	res.MultiAV = multiAvScan(filePath)

	// analysis finished
	res.Status = finished

	// Marshell results
	buff, err = json.Marshal(res)
	if err != nil {
		log.Error("Failed to get object, err: ", err)
		return err
	}

	// Update document
	updateDocument(sha256, buff)

	// Returning nil signals to the consumer that the message has
	// been handled with success. A FIN is sent to nsqd
	return nil
}


func main() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Load consumer config
	err := loadConfig()
	if err != nil {
		log.Fatal("Error parsing the config")
	}

	// Get an minio client instance
	accessKey := viper.GetString("do.accesskey")
	secKey := viper.GetString("do.seckey")
	endpoint := viper.GetString("do.endpoint")
	ssl := viper.GetBool("do.ssl")
	minioClient, err = minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		log.Fatal("Failed to connect to get minio client instance")
	}

	// Set backend API address
	backendEndpoint = viper.GetString("backend.address") + "/v1/files/"

	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Create a NewConsumer with the name of our topic, the channel, and our config
	consumer, err := nsq.NewConsumer("scan", "file", config)
	if err != nil {
		log.Errorln("Could not create consumer")
	}

	// Set the number of messages that can be in flight at any given time
	// you'll want to set this number as the default is only 1. This can be
	// a major concurrency knob that can change the performance of your application.
	consumer.ChangeMaxInFlight(200)

	// Here we set the logger to our NoopNSQLogger to quiet down the default logs.
	// At Reverb we use a custom structured logging format so we'll take the logging
	// from here.
	consumer.SetLogger(
		&NoopNSQLogger{},
		nsq.LogLevelError,
	)

	// Injects our handler into the consumer. You'll define one handler
	// per consumer, but you can have as many concurrently running handlers
	// as specified by the second argument. If your MaxInFlight is less
	// than your number of concurrent handlers you'll  starve your workers
	// as there will never be enough in flight messages for your worker pool
	consumer.AddConcurrentHandlers(
		&MessageHandler{},
		1,
	)

	// Our consumer will discover where topics are located by our three
	// nsqlookupd instances The application will periodically poll
	// these nqslookupd instances to discover new nodes or drop unhealthy
	// producers.
	nsqlds := viper.GetStringSlice("nsq.lookupd")
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		log.Fatal(err)
	}

	log.Infoln("Connected to nsqlookupd")

	// Let's allow our queues to drain properly during shutdown.
	// We'll create a channel to listen for SIGINT (Ctrl+C) to signal
	// to our application to gracefully shutdown.
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)

	// This is our main loop. It will continue to read off of our nsq
	// channel until either the consumer dies or our application is signaled
	// to stop.
	for {
		select {
		case <-consumer.StopChan:
			return // uh oh consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			consumer.Stop()
		}
	}
}
