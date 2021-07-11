// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/pe"
	"github.com/saferwall/saferwall/pkg/log"
	"github.com/saferwall/saferwall/pkg/pubsub"
	"github.com/saferwall/saferwall/pkg/pubsub/nsq"
	micro "github.com/saferwall/saferwall/services"
	"github.com/saferwall/saferwall/services/config"
)

// Config represents our application config.
type Config struct {
	LogLevel   string             `mapstructure:"log_level"`
	StorageDir string             `mapstructure:"storage_dir"`
	Producer   config.ProducerCfg `mapstructure:"producer"`
	Consumer   config.ConsumerCfg `mapstructure:"consumer"`
}

// Service represents the PE scan service. It adheres to the nsq.Handler
// interface. This allows us to define our own custom handlers for our messages.
// Think of these handlers much like you would an http handler.
type Service struct {
	cfg    Config
	logger log.Logger
	pub    pubsub.Publisher
	sub    pubsub.Subscriber
}

// New create a new PE scanner service.
func New(cfg Config, logger log.Logger) (Service, error) {
	var err error
	s := Service{}
	s.sub, err = nsq.NewSubscriber(
		cfg.Consumer.Topic,
		cfg.Consumer.Channel,
		cfg.Consumer.Lookupds,
		cfg.Consumer.Concurrency,
		&s,
	)
	if err != nil {
		return Service{}, err
	}

	s.pub, err = nsq.NewPublisher(cfg.Producer.Nsqd)
	if err != nil {
		return Service{}, err
	}

	s.cfg = cfg
	s.logger = logger
	return s, nil

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
		// returning an error results in the message being re-enqueued
		// a REQ is sent to nsqd
		return errors.New("body is blank re-enqueue message")
	}

	// deserialize the message.
	msg := micro.Message{}
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		s.logger.Error("failed to unmarshal msg")
		return err
	}

	sha256 := msg.Sha256
	s.logger = s.logger.With(context.TODO(), "sha256", sha256)

	s.logger.Infof("processing %s", sha256)
	filepath := filepath.Join(s.cfg.StorageDir, sha256)
	result, err := parse(filepath)
	if err != nil {
		s.logger.Error("failed to process file ...")
		return err
	}
	s.logger.Info("file parse success")

	msg.Body, err = json.Marshal(result)
	if err != nil {
		s.logger.Error("failed to process file ...")
		return err
	}

	data, err := json.Marshal(msg)
	if err != nil {
		s.logger.Error("failed to process file ...")
		return err
	}

	s.pub.Publish(context.TODO(), s.cfg.Producer.Topic, data)

	// Returning nil signals to the consumer that the message has
	// been handled with success. A FIN is sent to nsqd.
	return nil
}

func parse(filePath string) (*pe.File, error) {

	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	// Open the file and prepare it to be parsed.
	opts := pe.Options{SectionEntropy: true}
	f, err := pe.New(filePath, &opts)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Parse the PE.
	err = f.Parse()
	return f, err
}
