// Copyright 2022 Saferwall. All rights reserved.
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
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	"github.com/saferwall/saferwall/services/config"
	pb "github.com/saferwall/saferwall/services/proto"
	"google.golang.org/protobuf/proto"
)

// Config represents our application config.
type Config struct {
	LogLevel     string             `mapstructure:"log_level"`
	SharedVolume string             `mapstructure:"shared_volume"`
	Producer     config.ProducerCfg `mapstructure:"producer"`
	Consumer     config.ConsumerCfg `mapstructure:"consumer"`
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

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
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
		return errors.New("body is blank re-enqueue message")
	}

	ctx := context.Background()
	sha256 := string(m.Body)
	logger := s.logger.With(ctx, "sha256", sha256)

	logger.Info("start processing")

	src := filepath.Join(s.cfg.SharedVolume, sha256)
	file, err := parse(src)
	if err != nil {
		logger.Errorf("pe parsing failed: %v", err)
	}

	var tags []string
	if file.IsEXE() {
		tags = append(tags, "exe")
	} else if file.IsDriver() {
		tags = append(tags, "sys")
	} else if file.IsDLL() {
		tags = append(tags, "dll")
	}

	payloads := []*pb.Message_Payload{
		{Module: "pe", Body: toJSON(file)},
		{Module: "tags.pe", Body: toJSON(tags)},
	}

	msg := &pb.Message{Sha256: sha256, Payload: payloads}
	peMsg, err := proto.Marshal(msg)
	if err != nil {
		logger.Errorf("failed to marshal message: %v", err)
		return err
	}

	err = s.pub.Publish(ctx, s.cfg.Producer.Topic, peMsg)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	return nil
}

func parse(src string) (*pe.File, error) {

	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	// Open the file and prepare it to be parsed.
	opts := pe.Options{SectionEntropy: true}
	f, err := pe.New(src, &opts)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Parse the PE.
	err = f.Parse()
	return f, err
}
