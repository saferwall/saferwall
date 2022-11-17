// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package aggregator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	store "github.com/saferwall/saferwall/internal/db"
	"github.com/saferwall/saferwall/internal/log"
	pb "github.com/saferwall/saferwall/services/proto"
	"google.golang.org/protobuf/proto"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	s "github.com/saferwall/saferwall/internal/storage"
	"github.com/saferwall/saferwall/services/config"
)

// Config represents our application config.
type Config struct {
	LogLevel string             `mapstructure:"log_level"`
	Consumer config.ConsumerCfg `mapstructure:"consumer"`
	DB       store.Config       `mapstructure:"db"`
	Storage  config.StorageCfg  `mapstructure:"storage"`
}

// Service represents the PE scan service. It adheres to the nsq.Handler
// interface. This allows us to define our own custom handlers for our messages.
// Think of these handlers much like you would an http handler.
type Service struct {
	cfg     Config
	logger  log.Logger
	sub     pubsub.Subscriber
	db      store.DB
	storage s.Storage
}

// New create a new aggregator scanner service.
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

	svc.db, err = store.Open(cfg.DB.Server, cfg.DB.Username,
		cfg.DB.Password, cfg.DB.BucketName)
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

	opts.Bucket = cfg.Storage.Bucket

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

	msg := &pb.Message{}
	err := proto.Unmarshal(m.Body, msg)
	if err != nil {
		s.logger.Errorf("failed to unmarshal msg: %v", err)
		return err
	}

	sha256 := msg.Sha256
	ctx := context.Background()

	logger := s.logger.With(ctx, "sha256", sha256)

	for _, payload := range msg.Payload {
		key := payload.Key
		path := payload.Path

		switch payload.Kind {
		case pb.Message_DBUPDATE:
			var jsonPayload interface{}
			err = json.Unmarshal(payload.Body, &jsonPayload)
			if err != nil {
				logger.Errorf("failed to unmarshal json payload: %v", err)
				continue
			}
			logger.Debugf("payload is %v", jsonPayload)
			err = s.db.Update(ctx, key, path, jsonPayload)
			if err != nil {
				logger.Errorf("failed to update db: %v", err)
			}
		case pb.Message_DBCREATE:
			var jsonPayload interface{}
			err = json.Unmarshal(payload.Body, &jsonPayload)
			if err != nil {
				logger.Errorf("failed to unmarshal json payload: %v", err)
				continue
			}
			logger.Debugf("payload is %v", jsonPayload)
			err = s.db.Create(ctx, key, jsonPayload)
			if err != nil {
				logger.Errorf("failed to create document key: %s in db: %v", key, err)
			}
		case pb.Message_UPLOAD:
			obj := bytes.NewReader(payload.Body)
			err = s.storage.Upload(ctx, s.cfg.Storage.Bucket, key, obj)
			if err != nil {
				logger.Errorf("failed to upload object %s err: %v", key, err)
			}
		}

	}

	return nil
}
