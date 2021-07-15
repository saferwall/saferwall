// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package meta

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/log"
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/packer"
	"github.com/saferwall/saferwall/pkg/pubsub"
	"github.com/saferwall/saferwall/pkg/pubsub/nsq"
	"github.com/saferwall/saferwall/pkg/trid"
	"github.com/saferwall/saferwall/pkg/utils"
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

// Metadata represents meta data extarcted from a file.
type Metadata struct {
	MD5    string            `json:"md5,omitempty"`
	SHA1   string            `json:"sha1,omitempty"`
	SHA256 string            `json:"sha256,omitempty"`
	SHA512 string            `json:"sha512,omitempty"`
	SSDeep string            `json:"ssdeep,omitempty"`
	CRC32  string            `json:"crc32,omitempty"`
	Magic  string            `json:"magic,omitempty"`
	Size   int64             `json:"size,omitempty"`
	Exif   map[string]string `json:"exif,omitempty"`
	TriD   []string          `json:"trid,omitempty"`
	Packer []string          `json:"packer,omitempty"`
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
	data, err := utils.ReadAll(filePath)
	if err != nil {
		logger.Errorf("failed to read file, err: %v", err)
		return err
	}

	md := Metadata{}

	// Get crypto hashes.
	r := crypto.HashBytes(data)
	md.CRC32 = r.CRC32
	md.MD5 = r.MD5
	md.SHA1 = r.SHA1
	md.SHA256 = r.SHA256
	md.SHA512 = r.SHA512
	md.SSDeep = r.SSDeep

	// Get exif metadata.
	if md.Exif, err = exiftool.Scan(filePath); err != nil {
		logger.Errorf("exiftool scan failed with: %v", err)
	}

	// Get TriD file identifier results.
	if md.TriD, err = trid.Scan(filePath); err != nil {
		logger.Errorf("trid scan failed with: %v", err)
	}

	// Get lib magic scan results.
	if md.Magic, err = magic.Scan(filePath); err != nil {
		logger.Errorf("magic scan failed with: %v", err)
	}
	// Retrieve packer/crypter scan results.
	if md.Packer, err = packer.Scan(filePath); err != nil {
		logger.Errorf("packer scan failed with: %v", err)
	}

	logger.Info("file metadata extraction success")

	msg := &pb.Message{}
	hdr := &pb.Message_Header{Module: "meta", Sha256: sha256}

	msg.Body, err = json.Marshal(md)
	if err != nil {
		logger.Error("failed to process file ...")
		return err
	}

	msg.Header = hdr
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
