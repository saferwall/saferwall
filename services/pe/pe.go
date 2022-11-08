// Copyright 2018 Saferwall. All rights reserved.
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
		file.Close()
		return nil
	}

	// Extract PE related tags and file extension.
	var tags []string
	var ext string
	if file.IsEXE() {
		tags = append(tags, "exe")
		ext = "exe"
	} else if file.IsDLL() {
		tags = append(tags, "sys")
		ext = "sys"
	} else if file.IsDriver() {
		tags = append(tags, "dll")
		ext = "dll"
	}

	payloads := []*pb.Message_Payload{
		{Key: sha256, Path: "pe", Kind: pb.Message_DBUPDATE, Body: curate(file)},
		{Key: sha256, Path: "tags.pe", Kind: pb.Message_DBUPDATE, Body: toJSON(tags)},
		{Key: sha256, Path: "file_extension", Kind: pb.Message_DBUPDATE, Body: toJSON(ext)},
	}

	file.Close()

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

	// Parse the PE.
	err = f.Parse()
	return f, err
}

// curate the PE scan results by dropping all directories
// that are not available, also create a `meta` element
// that stores all the available PE fields. This make it
// easier for the UI to query only the needed part of the file.
func curate(file *pe.File) []byte {

	m := make(map[string]interface{})
	fields := make([]string, 0)

	if file.HasDOSHdr {
		m["dos_header"] = file.DOSHeader
		fields = append(fields, "dos_header")
	}

	if file.HasRichHdr {
		m["rich_header"] = file.RichHeader
		fields = append(fields, "rich_header")
	}

	if file.HasCOFF {
		m["coff"] = file.COFF
		fields = append(fields, "coff")
	}

	if file.HasNTHdr {
		m["nt_header"] = file.NtHeader
		fields = append(fields, "nt_header")
	}

	if file.HasSections {
		m["sections"] = file.Sections
		fields = append(fields, "sections")
	}

	if file.HasExport {
		m["export"] = file.Export
		fields = append(fields, "export")
	}

	if file.HasImport {
		m["import"] = file.Imports
		fields = append(fields, "import")
	}

	if file.HasResource {
		m["resource"] = file.Resources
		fields = append(fields, "resource")
	}

	if file.HasException {
		m["exception"] = file.Exceptions
		fields = append(fields, "exception")
	}

	if file.HasReloc {
		m["reloc"] = file.Relocations
		fields = append(fields, "reloc")
	}

	if file.HasDebug {
		m["debug"] = file.Debugs
		fields = append(fields, "debug")
	}

	if file.HasGlobalPtr {
		m["global_ptr"] = file.GlobalPtr
		fields = append(fields, "global_ptr")
	}

	if file.HasTLS {
		m["tls"] = file.TLS
		fields = append(fields, "tls")
	}

	if file.HasLoadCFG {
		m["load_config"] = file.LoadConfig
		fields = append(fields, "load_config")
	}

	if file.HasBoundImp {
		m["bound_import"] = file.BoundImports
		fields = append(fields, "bound_import")
	}

	if file.HasIAT {
		m["iat"] = file.IAT
		fields = append(fields, "iat")
	}

	if file.HasDelayImp {
		m["delay_import"] = file.DelayImports
		fields = append(fields, "delay_import")
	}

	if file.HasCLR {
		m["clr"] = file.CLR
		fields = append(fields, "clr")
	}

	if file.HasSecurity {
		m["security"] = file.Certificates
		if file.IsSigned {
			if file.Certificates.Verified {
				m["signature"] = "Signed file, valid signature"
			} else {
				m["signature"] = "Signed file, invalid signature"
			}
		}
		fields = append(fields, "security")
	} else {
		m["signature"] = "File is not signed"
	}

	m["meta"] = fields
	return toJSON(m)
}
