// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package multiav

import (
	"path/filepath"
	"strconv"

	"context"
	"encoding/json"
	"errors"

	"github.com/saferwall/saferwall/pkg/utils"
	"google.golang.org/protobuf/proto"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/pkg/log"
	"github.com/saferwall/saferwall/pkg/pubsub"
	"github.com/saferwall/saferwall/pkg/pubsub/nsq"

	m "github.com/saferwall/multiav/pkg"
	"github.com/saferwall/saferwall/services/config"
	pb "github.com/saferwall/saferwall/services/proto"
)

const (
	// Path to the file which holds the last time we updated the AV engine
	// database.
	dbUpdateDateFilePath = "/av_db_update_date.txt"
)

type Scanner interface {
	ScanFile(string) (m.Result, error)
}

// ScanResult av result
type ScanResult struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
	Update   int64  `json:"update"`
}

// Config represents our application config.
type Config struct {
	LogLevel     string             `mapstructure:"log_level"`
	SharedVolume string             `mapstructure:"shared_volume"`
	EngineName   string             `mapstructure:"engine_name"`
	Producer     config.ProducerCfg `mapstructure:"producer"`
	Consumer     config.ConsumerCfg `mapstructure:"consumer"`
}

// Service represents the Avast scan service.
type Service struct {
	cfg       Config
	logger    log.Logger
	av        Scanner
	pub       pubsub.Publisher
	sub       pubsub.Subscriber
	updatedAt int64
}

// updateDate returns a unix timestamp of the date when the database engine was
// updated.
func updateDate() (int64, error) {
	data, err := utils.ReadAll(dbUpdateDateFilePath)
	if err != nil {
		return 0, err
	}
	updateDate, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return int64(updateDate), nil
}

// New create a new PE scanner service.
func New(cfg Config, logger log.Logger, av Scanner) (Service, error) {
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
		return s, err
	}

	s.pub, err = nsq.NewPublisher(cfg.Producer.Nsqd)
	if err != nil {
		return s, err
	}

	s.updatedAt, err = updateDate()
	if err != nil {
		return s, err
	}

	s.cfg = cfg
	s.logger = logger
	s.av = av
	return s, nil

}

// Start kicks in the service to start consuming events.
func (s *Service) Start() error {
	s.logger.Infof("start consuming from topic: %s ...", s.cfg.Consumer.Topic)

	return s.sub.Start()
}

// HandleMessage is the only requirement needed to fulfill the nsq.Handler.
func (s *Service) HandleMessage(m *gonsq.Message) error {
	if len(m.Body) == 0 {
		// returning an error results in the message being re-enqueued
		// a REQ is sent to nsqd
		return errors.New("body is blank re-enqueue message")
	}

	sha256 := string(m.Body)
	ctx := context.Background()
	s.logger = s.logger.With(ctx, "sha256", sha256)

	s.logger.Info("start scanning")
	filepath := filepath.Join(s.cfg.SharedVolume, sha256)
	r, err := s.av.ScanFile(filepath)
	if err != nil {
		s.logger.Errorf("failed to scan file, reason: %v", err)
		return err
	}

	result := ScanResult{
		Infected: r.Infected,
		Output:   r.Output,
		Update:   s.updatedAt}

	msg := &pb.Message{}
	hdr := &pb.Message_Header{
		Module: "multiav." + s.cfg.EngineName,
		Sha256: sha256}

	msg.Body, err = json.Marshal(result)
	if err != nil {
		s.logger.Errorf("failed to marshal message body, reason: %v", err)
		return err
	}

	msg.Header = hdr
	out, err := proto.Marshal(msg)
	if err != nil {
		s.logger.Error("failed to pb marshal: %v", err)
		return err
	}

	s.pub.Publish(ctx, s.cfg.Producer.Topic, out)

	// Returning nil signals to the consumer that the message has
	// been handled with success. A FIN is sent to nsqd.
	return nil
}
