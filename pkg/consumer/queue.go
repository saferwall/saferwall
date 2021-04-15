// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"encoding/json"
	"errors"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	nsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/pkg/utils"
	log "github.com/sirupsen/logrus"
)

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
type MessageHandler struct {
	cfg         *Config
	minioClient *minio.Client
	authToken   string
}

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

	// Always include sha256 in our context logger.
	ctxLogger := log.WithFields(log.Fields{"sha256": sha256})
	ctxLogger.Info("start scanning ...")

	// Create a new file instance.
	f := File{Sha256: sha256}

	// Set the file status to `processing`.
	f.Status = processing
	err := h.updateMsgProgress(&f)
	if err != nil {
		ctxLogger.Errorf("failed to update message status: %v", err)
		return err
	}

	// Download the sample.
	filePath := path.Join("/samples", f.Sha256)
	b, err := h.downloadSample(filePath, &f)
	if err != nil {
		ctxLogger.Errorf("failed to download sample from s3: %v", err)
		return err
	}

	// Scan the file.
	err = f.Scan(sha256, filePath, b, ctxLogger, h.cfg)
	if err != nil {
		ctxLogger.Errorf("failed to scan the file: %v", err)
		return err
	}

	// Set the file status to `finished`.
	f.Status = finished
	err = h.updateMsgProgress(&f)
	if err != nil {
		ctxLogger.Errorf("failed to update message status: %v", err)
		return err
	}

	// Delete the file from the network share.
	if utils.Exists(filePath) {
		if err = utils.DeleteFile(filePath); err != nil {
			log.Errorf("failed to delete file path %s", filePath)
		}
	}

	// Returning nil signals to the consumer that the message has
	// been handled with success. A FIN is sent to nsqd.
	return nil
}

func (h *MessageHandler) updateMsgProgress(f *File) error {

	// Marshell results.
	buff, err := json.Marshal(f)
	if err != nil {
		return err
	}

	// Update document.
	err = updateDocument(f.Sha256, h.authToken, h.cfg, buff)
	if err == errHTTPStatusUnauthorized {

		// Get a new fresh jwt token.
		h.authToken, err = getAuthToken(h.cfg)
		if err != nil {
			return err
		}
		err = updateDocument(f.Sha256, h.authToken, h.cfg, buff)
	}
	return err
}

// Download sample from the object storage.
func (h *MessageHandler) downloadSample(filePath string, f *File) ([]byte, error) {
	bucketName := h.cfg.Minio.Spacename
	return utils.Download(h.minioClient, bucketName, filePath, f.Sha256)
}

// New creates a new consumer.
func New() (*nsq.Consumer, error) {

	// Load the config.
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load the config: %v", err)
	}

	// Setup logging.
	setupLogging(&cfg)

	// Authenticate to the web apis.
	authToken, err := getAuthToken(&cfg)
	if err != nil {
		log.Fatalf("failed to get auth token: %v", err)
	}

	// Get an minio client instance.
	accessKeyID := cfg.Minio.AccessKey
	secretAccessKey := cfg.Minio.SecKey
	endpoint := cfg.Minio.Endpoint
	useSSL := cfg.Minio.Ssl
	opts := minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	}
	minioClient, err := minio.New(endpoint, &opts)
	if err != nil {
		log.Fatalf("failed to connect to object storage: %v", err)
	}

	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Maximum number of times this consumer will attempt to process a message
	// before giving up.
	config.MaxAttempts = 2

	// Maximum number of messages to allow in flight (concurrency knob).
	config.MaxInFlight = 1

	// The server-side message timeout for messages delivered to this client.
	config.MsgTimeout = time.Duration(2 * time.Minute)

	// Create a NewConsumer with the name of our topic, the channel, and our config.
	consumer, err := nsq.NewConsumer("scan", "file", config)
	if err != nil {
		log.Fatalf("failed to create new consumer: %v", err)
	}

	// Create our message handler structure.
	messageHandler := MessageHandler{
		cfg:         &cfg,
		minioClient: minioClient,
		authToken:   authToken,
	}

	// Here we set the logger to our NoopNSQLogger to quiet down the default logs.
	// At Reverb we use a custom structured logging format so we'll take the
	// logging from here.
	consumer.SetLogger(
		&NoopNSQLogger{},
		nsq.LogLevelError,
	)

	// Injects our handler into the consumer. You'll define one handler
	// per consumer, but you can have as many concurrently running handlers
	// as specified by the second argument. If your MaxInFlight is less
	// than your number of concurrent handlers you'll starve your workers
	// as there will never be enough in flight messages for your worker pool
	consumer.AddConcurrentHandlers(
		&messageHandler,
		1,
	)

	// Our consumer will discover where topics are located by our three
	// nsqlookupd instances The application will periodically poll
	// these nqslookupd instances to discover new nodes or drop unhealthy
	// producers.
	nsqlds := cfg.Nsq.Lookupds
	consumer.ConnectToNSQLookupds(nsqlds)

	return consumer, err

}
