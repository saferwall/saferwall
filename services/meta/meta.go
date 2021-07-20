// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package meta

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/log"
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/packer"
	str "github.com/saferwall/saferwall/pkg/strings"
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

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
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
		logger.Errorf("failed to read file path: %s, err: %v", filePath, err)
		return err
	}

	// Get crypto hashes.
	r := crypto.HashBytes(data)

	// Get exif metadata.
	exif, err := exiftool.Scan(filePath)
	if err != nil {
		logger.Errorf("exiftool scan failed with: %v", err)
	}

	// Get TriD file identifier results.
	tridRes, err := trid.Scan(filePath)
	if err != nil {
		logger.Errorf("trid scan failed with: %v", err)
	}

	// Get lib magic scan results.
	magicRes, err := magic.Scan(filePath)
	if err != nil {
		logger.Errorf("magic scan failed with: %v", err)
	}

	// Retrieve packer/crypter scan results.
	packerRes, err := packer.Scan(filePath)
	if err != nil {
		logger.Errorf("packer scan failed with: %v", err)
	}

	// Determine type.
	var format string
	for k, v := range typeMap {
		if strings.HasPrefix(magicRes, k) {
			format = v
			break
		}
	}

	// Extract strings.
	n := 8
	asciiStrings := str.GetASCIIStrings(data, n)
	wideStrings := str.GetUnicodeStrings(data, n)
	asmStrings := str.GetAsmStrings(data)
	string_res := map[string]interface{} {
		"ascii":  utils.UniqueSlice(asciiStrings),
		"wide":  utils.UniqueSlice(wideStrings),
		"asm":  utils.UniqueSlice(asmStrings),
	}

	logger.Info("file metadata extraction success")

	var tags []string
	for _, out := range packerRes {
		if strings.Contains(out, "packer") ||
			strings.Contains(out, "protector") ||
			strings.Contains(out, "compiler") ||
			strings.Contains(out, "installer") ||
			strings.Contains(out, "library") {
			for sig, tag := range sigMap {
				if strings.Contains(out, sig) {
					tags = append(tags, tag)
				}
			}
		}
	}

	logger.Info("tags extraction success")

	payloads := []*pb.Message_Payload{
		{Module: "crc32", Body: toJSON(r.CRC32)},
		{Module: "md5", Body: toJSON(r.MD5)},
		{Module: "sha1", Body: toJSON(r.SHA1)},
		{Module: "sha256", Body: toJSON(r.SHA256)},
		{Module: "sha256", Body: toJSON(r.SHA256)},
		{Module: "sha512", Body: toJSON(r.SHA512)},
		{Module: "ssdeep", Body: toJSON(r.SSDeep)},
		{Module: "size", Body: toJSON(int64(len(data)))},
		{Module: "exif", Body: toJSON(exif)},
		{Module: "trid", Body: toJSON(tridRes)},
		{Module: "magic", Body: toJSON(magicRes)},
		{Module: "packer", Body: toJSON(packerRes)},
		{Module: "tags.packer", Body: toJSON(tags)},
		{Module: "strings", Body: toJSON(string_res)},
		{Module: "fileformat", Body: toJSON(format)},
	}

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
