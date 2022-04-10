// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package aggregator

import (
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
	"github.com/saferwall/saferwall/services/config"
)

// Config represents our application config.
type Config struct {
	LogLevel string             `mapstructure:"log_level"`
	Consumer config.ConsumerCfg `mapstructure:"consumer"`
	DB       store.Config       `mapstructure:"db"`
}

// Service represents the PE scan service. It adheres to the nsq.Handler
// interface. This allows us to define our own custom handlers for our messages.
// Think of these handlers much like you would an http handler.
type Service struct {
	cfg    Config
	logger log.Logger
	sub    pubsub.Subscriber
	db     store.DB
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

	svc.db, err = store.Open(cfg.DB.Server, cfg.DB.Username,
		cfg.DB.Password, cfg.DB.BucketName)
	if err != nil {
		return Service{}, err
	}

	svc.cfg = cfg
	svc.logger = logger
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

	for _, payload := range msg.Payload {
		path := payload.Module
		logger := s.logger.With(ctx, "sha256", sha256, "module", path)

		var jsonPayload interface{}
		err = json.Unmarshal(payload.Body, &jsonPayload)
		if err != nil {
			logger.Errorf("failed to unmarshal json payload: %v", err)
			continue
		}

		logger.Debugf("payload is %v", jsonPayload)
		err = s.db.Update(ctx, sha256, path, jsonPayload)
		if err != nil {
			logger.Errorf("failed to update db: %v", err)
		}
	}

	return nil
}
