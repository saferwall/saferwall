// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/saferwall/saferwall/pkg/log"
	micro "github.com/saferwall/saferwall/services"

	"github.com/gabriel-vasile/mimetype"
	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/pkg/pubsub"
	"github.com/saferwall/saferwall/pkg/pubsub/nsq"
	s "github.com/saferwall/saferwall/pkg/storage"
	"github.com/saferwall/saferwall/services/config"
)

// Config represents our application config.
type Config struct {
	LogLevel   string             `mapstructure:"log_level"`
	Deployment string             `mapstructure:"deployment"`
	Producer   config.ProducerCfg `mapstructure:"producer"`
	Consumer   config.ConsumerCfg `mapstructure:"consumer"`
	Storage    config.StorageCfg  `mapstructure:"storage"`
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
	switch cfg.Deployment {
	case "aws":
		opts.S3AccKey = cfg.Storage.S3.AccessKey
		opts.S3SecKey = cfg.Storage.S3.SecretKey
		opts.S3Region = cfg.Storage.S3.Region
	case "local":
		opts.LocalRootDir = cfg.Storage.Local.RootDir
	}

	sto, err := s.New(cfg.Deployment, opts)
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
	s.logger.Debug("Starting the service ...")
	s.sub.Start()

	return nil
}

// HandleMessage is the only requirement needed to fulfill the nsq.Handler.
func (s *Service) HandleMessage(m *gonsq.Message) error {
	if len(m.Body) == 0 {
		// returning an error results in the message being re-enqueued
		// a REQ is sent to nsqd
		return errors.New("body is blank re-enqueue message")
	}

	sha256 := string(m.Body)
	s.logger = s.logger.With(context.Background(), "sha256", sha256)

	// Download the file from object storage and place it in a directory
	// shared between all microservices.
	filePath := filepath.Join(s.cfg.Storage.SharedVolume, sha256)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	if err := s.storage.Download(s.cfg.Storage.Bucket, sha256, file,
		s.cfg.Storage.Timeout); err != nil {
		return err
	}

	// Depending on what the file format is, we produce events to different
	// consumers.
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return err
	}
	msg, err := json.Marshal(micro.Message{Sha256: sha256, Body: nil})
	if err != nil {
		return err
	}

	switch mtype.String() {
	case "application/vnd.microsoft.portable-executable":
		err = s.pub.Publish(context.Background(), "topic-pe", msg)
		if err != nil {
			return err
		}
	case "elf":
		err = s.pub.Publish(context.TODO(), "topic-elf", m.Body)
		if err != nil {
			return err
		}
	case "mach-o":
		err = s.pub.Publish(context.TODO(), "topic-mach-o", m.Body)
		if err != nil {
			return err
		}
	case "pdf":
		err = s.pub.Publish(context.TODO(), "topic-pdf", m.Body)
		if err != nil {
			return err
		}
	}

	// we always scan the file no matter which format it is with the multi-av
	// scanner.
	s.pub.Publish(context.TODO(), "topic-multiav", m.Body)

	// Returning nil signals to the consumer that the message has
	// been handled with success. A FIN is sent to nsqd.
	return nil
}
