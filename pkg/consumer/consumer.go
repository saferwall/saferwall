// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nsqio/go-nsq"
)

var (
	defaultNSQMaxAttempt  uint16 = 2
	defaultNSQMaxInFlight        = 1
	defaultNSQMaxTimeout         = time.Duration(2 * time.Minute)
	scanTopic                    = "scan"
	scanChannel                  = "file"
)

// Consumer type wraps an NSQ consumer and handles
// the service runtime.
type Consumer struct {
	cfg       Config
	c         *nsq.Consumer
	handler   MessageHandler
	auth      func(cfg *Config) (string, error)
	authToken string
}

// New creates a new consumer instance.
func New(cfg Config) (*Consumer, error) {
	nsqConfig := NewNSQConfig(defaultNSQMaxAttempt, defaultNSQMaxInFlight, defaultNSQMaxTimeout)
	cons, err := nsq.NewConsumer(scanTopic, scanChannel, nsqConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		cfg:       cfg,
		c:         cons,
		handler:   MessageHandler{},
		auth:      fetchAuthToken,
		authToken: "",
	}, nil
}

// Start will start the consumer service and try
// to fetch an authentification token and setup
// logging.
// If both are successful it will initaite a connection
// to the S3 bucket through minio and start the underlying
// NSQ Consumer which does the queuing and processing.
// Start only fails if one of the requirements are unfullfilled
// - Authentification
// - Logging
// - S3 connection
// Otherwise the service simulates the ListenAndServe function
// and runs in a loop unless a shutdown message is received.
func (c *Consumer) Start() error {
	var err error

	if err != nil {
		log.Fatal(err)
	}
	// Setup logging.
	setupLogging(&c.cfg)
	fmt.Println(c.cfg)
	// When running in Headless mode we only want to run the consumer
	// without interacting with the backend API.
	if !c.cfg.Headless {
		// Setup API Authentification
		c.authToken, err = fetchAuthToken(&c.cfg)
		if err != nil {
			log.Fatalf("failed to get auth token: %v", err)
			return err
		}

	}
	// Setup MinioClient
	minioClient, err := NewMinioClient(c.cfg)
	if err != nil {
		return err
	}
	// Create our message handler structure.
	messageHandler := MessageHandler{
		cfg:         &c.cfg,
		minioClient: minioClient,
		authToken:   c.authToken,
	}
	c.handler = messageHandler

	// Here we set the logger to our NoopNSQLogger to quiet down the default logs.
	// At Reverb we use a custom structured logging format so we'll take the
	// logging from here.
	c.c.SetLogger(
		&NoopNSQLogger{},
		nsq.LogLevelError,
	)

	// Injects our handler into the consumer. You'll define one handler
	// per consumer, but you can have as many concurrently running handlers
	// as specified by the second argument. If your MaxInFlight is less
	// than your number of concurrent handlers you'll starve your workers
	// as there will never be enough in flight messages for your worker pool
	c.c.AddConcurrentHandlers(
		&c.handler,
		1,
	)

	// Our consumer will discover where topics are located by our three
	// nsqlookupd instances The application will periodically poll
	// these nqslookupd instances to discover new nodes or drop unhealthy
	// producers.
	nsqlds := c.cfg.NSQ.Lookupds
	c.c.ConnectToNSQLookupds(nsqlds)
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
		case <-c.c.StopChan:
			return nil // uh oh consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			c.Stop()
		}
	}

}

// Stop the underlying NSQ service gracefully and clean up.
func (c *Consumer) Stop() {
	c.c.Stop()
}

// NewNSQConfig creates a new NSQ Config.
func NewNSQConfig(maxAttempets uint16, maxInFlight int, timeout time.Duration) *nsq.Config {
	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Maximum number of times this consumer will attempt to process a message
	// before giving up.
	config.MaxAttempts = maxAttempets

	// Maximum number of messages to allow in flight (concurrency knob).
	config.MaxInFlight = maxInFlight

	// The server-side message timeout for messages delivered to this client.
	config.MsgTimeout = timeout
	return config
}

// NewMinioClient creates a new minio client instance.
func NewMinioClient(cfg Config) (*minio.Client, error) {
	// Get an minio client instance.
	accessKeyID := cfg.Minio.AccessKey
	secretAccessKey := cfg.Minio.SecKey
	endpoint := cfg.Minio.Endpoint
	useSSL := cfg.Minio.SSL
	opts := minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	}
	minioClient, err := minio.New(endpoint, &opts)
	if err != nil {
		log.Fatalf("failed to connect to object storage: %v", err)
	}
	return minioClient, err
}
