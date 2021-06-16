// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"encoding/json"
	"errors"
	"path"
	"runtime/debug"

	"github.com/minio/minio-go/v7"
	nsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/pkg/utils"
	"go.uber.org/zap"
)

// Status represents a file's status during scanning.
type Status int

// File scan progress status.
const (
	Queued Status = iota
	Processing
	Finished
)

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
type MessageHandler struct {
	cfg         *Config
	minioClient *minio.Client
	logger      *zap.SugaredLogger
	authToken   string
}

// processMessage is the core processing function when
// handling NSQ messages.
// processMessage will setup contextual logging and run
// the different scanners.
func (h *MessageHandler) processMessage(b []byte) error {
	// read the sample hash
	sha256 := string(b)
	// setup contextual logging using the hash as context
	h.logger.Debug("start scanning ...")
	h.logger = h.logger.With(
		"sha256", sha256,
	)
	return h.scanFile(sha256)
}

// scanFile does the actual file scanning
func (h *MessageHandler) scanFile(sha256 string) error {

	// Handle unexpected panics.
	defer func() {
		if r := recover(); r != nil {
			h.logger.DPanic("panic occured in file scan: %v", debug.Stack())
		}
	}()

	// Create a new file instance.
	f := File{SHA256: sha256}

	// Download the sample first before updating queue status.
	filePath := path.Join(h.cfg.DownloadDir, f.SHA256)
	b, err := h.fetchSample(filePath, &f)
	if err != nil {
		h.logger.Errorf("failed to download sample from s3: %v", err)
		return err
	}
	// Set the file status to `processing`.
	f.Status = int(Processing)
	err = h.updateMsgProgress(&f)
	if err != nil {
		h.logger.Errorf("failed to update message status: %v", err)
		return err
	}

	// Scan the file.
	err = f.Scan(sha256, filePath, b, h.logger, h.cfg)
	if err != nil {
		h.logger.Errorf("failed to scan the file: %v", err)
		return err
	}

	// Set the file status to `finished`.
	f.Status = int(Finished)
	err = h.updateMsgProgress(&f)
	if err != nil {
		h.logger.Errorf("failed to update message status: %v", err)
		return err
	}

	// Delete the file from the network share.
	if utils.Exists(filePath) {
		if err = utils.DeleteFile(filePath); err != nil {
			h.logger.Errorf("failed to delete file path %s", filePath)
		}
	}

	return nil
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

	return h.processMessage(m.Body)
}

func (h *MessageHandler) updateMsgProgress(f *File) error {

	// Marshell results.
	buff, err := json.Marshal(f)
	if err != nil {
		return err
	}
	// Update document.
	err = updateDocument(f.SHA256, h.authToken, h.cfg, buff)
	if err == errHTTPStatusUnauthorized {
		// Get a new fresh jwt token.
		h.authToken, err = fetchAuthToken(h.cfg)
		if err != nil {
			return err
		}
		err = updateDocument(f.SHA256, h.authToken, h.cfg, buff)
	}
	return err
}

// fetch sample from the object storage.
func (h *MessageHandler) fetchSample(filePath string, f *File) ([]byte, error) {
	bucketName := h.cfg.Minio.Spacename
	return utils.Download(h.minioClient, bucketName, filePath, f.SHA256)
}
