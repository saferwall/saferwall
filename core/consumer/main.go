// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	nsq "github.com/bitly/go-nsq"
	minio "github.com/minio/minio-go"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/magic"
	s "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/pkg/trid"
	"github.com/saferwall/saferwall/pkg/utils"
	"github.com/saferwall/saferwall/pkg/utils/do"
	log "github.com/sirupsen/logrus"
)

const (
	addr     = "127.0.0.1:4161"
	endpoint = "http://192.168.99.100:30081/v1/files/"
)

var (
	client *minio.Client
)

type stringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}
type result struct {
	Crc32   string            `json:"crc32"`
	Md5     string            `json:"md5"`
	Sha1    string            `json:"sha1"`
	Sha256  string            `json:"sha256"`
	Sha512  string            `json:"sha512"`
	Ssdeep  string            `json:"ssdeep"`
	Exif    map[string]string `json:"exif"`
	TriD    []string          `json:"trid"`
	Magic   string            `json:"magic"`
	Strings []stringStruct    `json:"strings"`
}

// NoopNSQLogger allows us to pipe NSQ logs to dev/null
// The default NSQ logger is great for debugging, but did
// not fit our normally well structured JSON logs. Luckily
// NSQ provides a simple interface for injecting your own
// logger.
type NoopNSQLogger struct{}

// Output allows us to implement the nsq.Logger interface
func (l *NoopNSQLogger) Output(calldepth int, s string) error {
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
	filePath := path.Join("/tmp", sha256)

	// Download the sample
	err := client.FGetObject("samples", sha256, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Error("Failed to get object, err: ", err)
		return err
	}

	b, err := utils.ReadAll(filePath)
	if err != nil {
		log.Error("Failed to read file, err: ", err)
		return err
	}

	// Run crypto pkg
	r := crypto.HashBytes(b)
	res := result{
		Crc32:  r.Crc32,
		Md5:    r.Md5,
		Sha1:   r.Sha1,
		Sha256: r.Sha256,
		Sha512: r.Sha512,
		Ssdeep: r.Ssdeep,
	}

	// Run exiftool pkg
	res.Exif, err = exiftool.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with exiftool, err: ", err)
		return err
	}

	// Run TRiD pkg
	res.TriD, err = trid.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with trid, err: ", err)
		return err
	}

	// Run Magic Pkg
	res.Magic, err = magic.GetMimeType(b)
	if err != nil {
		log.Error("Failed to scan file with magic, err: ", err)
		return err
	}

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

	buff, err := json.Marshal(res)
	if err != nil {
		log.Error("Failed to get object, err: ", err)
		return err
	}

	client := &http.Client{}
	client.Timeout = time.Second * 15
	url := endpoint + sha256
	body := bytes.NewBuffer(buff)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Errorf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	log.Infof("Response status code: %d, text: %s", resp.StatusCode, string(d))

	// Returning nil signals to the consumer that the message has
	// been handled with success. A FIN is sent to nsqd
	return nil
}

func main() {

	log.Info("0.0.5")

	client = do.GetClient()

	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Create a NewConsumer with the name of our topic, the channel, and our config
	consumer, err := nsq.NewConsumer("scan", "file", config)
	if err != nil {
		log.Fatal("Could not create consumer")
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
		20,
	)

	// Our consumer will discover where topics are located by our three
	// nsqlookupd instances The application will periodically poll
	// these nqslookupd instances to discover new nodes or drop unhealthy
	// producers.
	// nsqlds := []string{
	// 	"nsqlookupd-0.nsqlookupd.default.svc.cluster.local:4161",
	// 	"nsqlookupd-1.nsqlookupd.default.svc.cluster.local:4161",
	// 	"nsqlookupd-2.nsqlookupd.default.svc.cluster.local:4161"}
	if err := consumer.ConnectToNSQLookupd("nsqlookupd:4161"); err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to nsqlookupd")

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
