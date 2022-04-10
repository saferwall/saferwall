// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/saferwall/saferwall/internal/log"

	"github.com/gabriel-vasile/mimetype"
	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	s "github.com/saferwall/saferwall/internal/storage"
	"github.com/saferwall/saferwall/services/config"
)

// Config represents our application config.
type Config struct {
	// Log level. Defaults to info.
	LogLevel     string             `mapstructure:"log_level"`
	SharedVolume string             `mapstructure:"shared_volume"`
	Producer     config.ProducerCfg `mapstructure:"producer"`
	Consumer     config.ConsumerCfg `mapstructure:"consumer"`
	Storage      config.StorageCfg  `mapstructure:"storage"`
}

// Service represents the PE scan service. It adheres to the nsq.Handler
// interface. This allows us to define our own custom handlers for our messages.
// Think of these handlers much like you would an http handler.
type Service struct {
	cfg     Config
	logger  log.Logger
	pub     pubsub.Publisher
	sub     pubsub.Subscriber
	storage s.Storage
}

// Progress of a file scan.
const (
	queued     = iota
	processing = iota
	finished   = iota
)

// New create a new PE scanner service.
func New(cfg Config, logger log.Logger) (Service, error) {

	svc := Service{}
	var err error

	svc.sub, err = nsq.NewSubscriber(
		cfg.Consumer.Topic,
		cfg.Consumer.Channel,
		cfg.Consumer.Lookupds,
		cfg.Consumer.Concurrency,
		&svc,
	)
	if err != nil {
		return Service{}, err
	}

	svc.pub, err = nsq.NewPublisher(cfg.Producer.Nsqd)
	if err != nil {
		return Service{}, err
	}

	opts := s.Options{}
	switch cfg.Storage.DeploymentKind {
	case "aws":
		opts.Region = cfg.Storage.S3.Region
		opts.AccessKey = cfg.Storage.S3.AccessKey
		opts.SecretKey = cfg.Storage.S3.SecretKey
	case "minio":
		opts.Region = cfg.Storage.Minio.Region
		opts.AccessKey = cfg.Storage.Minio.AccessKey
		opts.SecretKey = cfg.Storage.Minio.SecretKey
		opts.MinioEndpoint = cfg.Storage.Minio.Endpoint
	case "local":
		opts.LocalRootDir = cfg.Storage.Local.RootDir
	}

	sto, err := s.New(cfg.Storage.DeploymentKind, opts)
	if err != nil {
		return Service{}, err
	}

	svc.cfg = cfg
	svc.logger = logger
	svc.storage = sto
	return svc, nil
}

// Start kicks in the service to start consuming events.
func (s *Service) Start() error {
	s.logger.Infof("start consuming from topic: %s ...", s.cfg.Consumer.Topic)
	s.sub.Start()

	return nil
}

// HandleMessage is the only requirement needed to fulfill the nsq.Handler.
func (s *Service) HandleMessage(m *gonsq.Message) error {
	if len(m.Body) == 0 {
		return errors.New("body is blank re-enqueue message")
	}

	fileScanCfg := config.FileScanCfg{}
	ctx := context.Background()

	// Deserialize the msg sent from the web apis.
	err := json.Unmarshal(m.Body, &fileScanCfg)
	if err != nil {
		s.logger.Errorf("failed unmarshalling json messge body: %v", err)
		return err
	}

	sha256 := fileScanCfg.SHA256
	logger := s.logger.With(ctx, "sha256", sha256)

	logger.Info("start processing")

	// Download the file from object storage and place it in a directory
	// shared between all microservices.
	filePath := filepath.Join(s.cfg.SharedVolume, sha256)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("failed creating file: %v", err)
		return err
	}

	// Create a context with a timeout that will abort the download if it takes
	// more than the passed in timeout.
	downloadCtx, cancelFn := context.WithTimeout(
		context.Background(), time.Duration(time.Second*30))
	defer cancelFn()

	if err := s.storage.Download(
		downloadCtx, s.cfg.Storage.Bucket, sha256, file); err != nil {
		logger.Errorf("failed downloading file: %v", err)
		return err
	}
	file.Close()

	logger.Debugf("file downloaded to %s", filePath)

	// always run the multi-av scanner and the metadata
	// extractor no matter what the file format is.
	err = s.pub.Publish(ctx, "topic-multiav", []byte(sha256))
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	logger.Debug("published messaged to topic-multiav")

	err = s.pub.Publish(ctx, "topic-meta", []byte(sha256))
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	logger.Debug("published messaged to topic-meta")

	// depending on what the file format is,
	// we produce events to different consumers.
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		logger.Errorf("failed to detect mimetype: %v", err)
		return err
	}

	logger.Debugf("file type is: %s", mtype.String())

	switch mtype.String() {
	case "application/vnd.microsoft.portable-executable":
		if err = s.pub.Publish(ctx, "topic-pe", []byte(sha256)); err != nil {
			logger.Errorf("failed to publish message: %v", err)
			return err
		}
		logger.Debug("published messaged to topic-pe")

		if err = s.pub.Publish(ctx, "topic-sandbox", m.Body); err != nil {
			logger.Errorf("failed to publish message: %v", err)
			return err
		}
		logger.Debug("published messaged to topic-sandbox")
	}

	err = s.pub.Publish(ctx, "topic-postprocessor", []byte(sha256))
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	logger.Debug("published messaged to topic-postprocessor")
	return nil
}
