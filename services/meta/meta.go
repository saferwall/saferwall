// Copyright 2018 Saferwall. All rights reserved.
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
	"github.com/saferwall/saferwall/internal/exiftool"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/magic"
	"github.com/saferwall/saferwall/internal/packer"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	"github.com/saferwall/saferwall/internal/trid"
	"github.com/saferwall/saferwall/internal/utils"
	bs "github.com/saferwall/saferwall/pkg/bytestats"
	"github.com/saferwall/saferwall/pkg/crypto"
	str "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/services/config"
	pb "github.com/saferwall/saferwall/services/proto"
	"google.golang.org/protobuf/proto"
)

const (
	maxStrLength = 8
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
		return errors.New("body is blank re-enqueue message")
	}

	sha256 := string(m.Body)
	ctx := context.Background()
	logger := s.logger.With(ctx, "sha256", sha256)

	logger.Info("start processing")

	filePath := filepath.Join(s.cfg.SharedVolume, sha256)
	data, err := utils.ReadAll(filePath)
	if err != nil {
		logger.Errorf("failed to read file path: %s, err: %v",
			filePath, err)
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

	// Determine file format.
	var fileFormat string
	var fileExt string
	for k, v := range typeMap {
		if strings.Contains(magicRes, k) {
			fileFormat = v
			break
		}
	}
	if len(fileFormat) == 0 {
		fileFormat = "unknown"
	}
	logger.Debugf("file format is: %s", fileFormat)

	// Determine file extension.
	if fileFormat != "unknown" {
		fileExt = guessFileExtension(data, magicRes, fileFormat, tridRes)
	} else {
		fileExt = "unknown"
	}
	logger.Debugf("file extension is: %s", fileExt)

	// Extract strings.
	asciiStrings := str.GetASCIIStrings(&data, maxStrLength)
	wideStrings := str.GetUnicodeStrings(&data, maxStrLength)
	asmStrings := str.GetAsmStrings(&data)
	stringRes := map[string]interface{}{
		"ascii": utils.UniqueSlice(asciiStrings),
		"wide":  utils.UniqueSlice(wideStrings),
		"asm":   utils.UniqueSlice(asmStrings),
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
		{Key: sha256, Path: "crc32", Kind: pb.Message_DBUPDATE, Body: toJSON(r.CRC32)},
		{Key: sha256, Path: "md5", Kind: pb.Message_DBUPDATE, Body: toJSON(r.MD5)},
		{Key: sha256, Path: "sha1", Kind: pb.Message_DBUPDATE, Body: toJSON(r.SHA1)},
		{Key: sha256, Path: "sha256", Kind: pb.Message_DBUPDATE, Body: toJSON(r.SHA256)},
		{Key: sha256, Path: "sha256", Kind: pb.Message_DBUPDATE, Body: toJSON(r.SHA256)},
		{Key: sha256, Path: "sha512", Kind: pb.Message_DBUPDATE, Body: toJSON(r.SHA512)},
		{Key: sha256, Path: "ssdeep", Kind: pb.Message_DBUPDATE, Body: toJSON(r.SSDeep)},
		{Key: sha256, Path: "size", Kind: pb.Message_DBUPDATE, Body: toJSON(int64(len(data)))},
		{Key: sha256, Path: "exif", Kind: pb.Message_DBUPDATE, Body: toJSON(exif)},
		{Key: sha256, Path: "trid", Kind: pb.Message_DBUPDATE, Body: toJSON(tridRes)},
		{Key: sha256, Path: "magic", Kind: pb.Message_DBUPDATE, Body: toJSON(magicRes)},
		{Key: sha256, Path: "packer", Kind: pb.Message_DBUPDATE, Body: toJSON(packerRes)},
		{Key: sha256, Path: "tags.packer", Kind: pb.Message_DBUPDATE, Body: toJSON(tags)},
		{Key: sha256, Path: "strings", Kind: pb.Message_DBUPDATE, Body: toJSON(stringRes)},
		{Key: sha256, Path: "histogram", Kind: pb.Message_DBUPDATE, Body: toJSON(bs.ByteHistogram(data))},
		{Key: sha256, Path: "byte_entropy", Kind: pb.Message_DBUPDATE, Body: toJSON(bs.ByteEntropyHistogram(data))},
		{Key: sha256, Path: "file_format", Kind: pb.Message_DBUPDATE, Body: toJSON(fileFormat)},
	}

	if fileExt != "unknown" {
		payloads = append(payloads, &pb.Message_Payload{
			Key:  sha256,
			Path: "file_extension",
			Kind: pb.Message_DBUPDATE,
			Body: toJSON(fileExt)})
	}

	msg := &pb.Message{Sha256: sha256, Payload: payloads}
	out, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("failed to pb marshal: %v", err)
		return err
	}

	err = s.pub.Publish(ctx, s.cfg.Producer.Topic, out)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	return nil
}
