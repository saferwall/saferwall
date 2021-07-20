// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/pe"
	bs "github.com/saferwall/saferwall/pkg/bytestats"
	"github.com/saferwall/saferwall/pkg/log"
	"github.com/saferwall/saferwall/pkg/pubsub"
	"github.com/saferwall/saferwall/pkg/pubsub/nsq"
	"github.com/saferwall/saferwall/pkg/utils"
	"github.com/saferwall/saferwall/services/config"
	pb "github.com/saferwall/saferwall/services/proto"
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

// Result represents the scan results produced by the PE svc. It embeds the
// histogram and the byte entropy.
type Result struct {
	*pe.File
	Histogram   []int `json:"histogram,omitempty"`
	ByteEntropy []int `json:"byte_entropy,omitempty"`
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

	sha256 := string(m.Body)
	logger := s.logger.With(context.TODO(), "sha256", sha256)

	logger.Infof("processing %s", sha256)

	filePath := filepath.Join(s.cfg.SharedVolume, sha256)
	peres, err := parse(filePath)
	if err != nil {
		logger.Error("failed to parse pe file: %v", err)
	}

	var tags []string
	if peres.IsEXE() {
		tags = append(tags, "exe")
	} else if peres.IsDriver() {
		tags = append(tags, "sys")
	} else if peres.IsDLL() {
		tags = append(tags, "dll")
	}

	// Extract Byte Histogram and byte entropy.
	b, err := utils.ReadAll(filePath)
	if err != nil {
		logger.Errorf("failed to read file path: %s, err: %v", filePath, err)
		return err
	}

	res := Result{
		File:        peres,
		Histogram:   bs.ByteHistogram(b),
		ByteEntropy: bs.ByteEntropyHistogram(b),
	}

	bodyPE, err := json.Marshal(res)
	if err != nil {
		logger.Error("failed to marshal json: %v", err)
		return err
	}
	bodyTags, err := json.Marshal(tags)
	if err != nil {
		logger.Error("failed to marshal json: %v", err)
		return err
	}

	payloads := []*pb.Message_Payload{}
	payloads = append(payloads, &pb.Message_Payload{
		Module: "pe",
		Body:   bodyPE,
	})
	payloads = append(payloads, &pb.Message_Payload{
		Module: "tags.pe",
		Body:   bodyTags,
	})

	msg := &pb.Message{Sha256: sha256, Payload: payloads}
	out, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("failed to pb marshal: %v", err)
		return err
	}

	err = s.pub.Publish(context.TODO(), s.cfg.Producer.Topic, out)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

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
