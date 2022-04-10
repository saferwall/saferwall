// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package multiav

import (
	"path/filepath"
	"strconv"

	"context"
	"encoding/json"
	"errors"

	"github.com/saferwall/saferwall/internal/utils"
	"google.golang.org/protobuf/proto"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	"github.com/saferwall/saferwall/pkg/avlabel"

	m "github.com/saferwall/saferwall/internal/multiav"
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

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// Start kicks in the service to start consuming events.
func (s *Service) Start() error {
	s.logger.Infof("start consuming from topic: %s ...", s.cfg.Consumer.Topic)
	return s.sub.Start()
}

// HandleMessage is the only requirement needed to fulfill the nsq.Handler.
func (s *Service) HandleMessage(m *gonsq.Message) error {
	if len(m.Body) == 0 {
		return errors.New("body is blank re-enqueue message")
	}

	sha256 := string(m.Body)
	ctx := context.Background()
	logger := s.logger.With(ctx, "sha256", sha256)

	logger.Info("start scanning")

	// some multiav scanners lock the file in the nfs share
	// and prevent the other scanners from accessing the file,
	// we copy the file locally to the container to avoid this issue.
	src := filepath.Join(s.cfg.SharedVolume, sha256)
	dest := filepath.Join("/tmp", sha256)
	if err := utils.CopyFile(src, dest); err != nil {
		logger.Errorf("failed to copy file, reason: %v", err)
		return err
	}

	r, err := s.av.ScanFile(dest)
	if err != nil {
		logger.Errorf("failed to scan file, reason: %v", err)
	}

	if utils.Exists(dest) {
		if err = utils.DeleteFile(dest); err != nil {
			logger.Errorf("Failed to delete file path %s.", dest)
		}
	}

	logger.Debugf("finished scanning: output: %s, infected:%v, out: %s",
		r.Output, r.Infected, r.Out)

	result := ScanResult{
		Infected: r.Infected,
		Output:   r.Output,
		Update:   s.updatedAt}

	payloads := []*pb.Message_Payload{
		{
			Module: "multiav.last_scan." + s.cfg.EngineName,
			Body:   toJSON(result),
		},
	}

	// get the multiav tag.
	if r.Infected && len(r.Output) > 0 {
		parsedDet := avlabel.Parse(s.cfg.EngineName, r.Output)
		if len(parsedDet.Family) > 0 {
			payloads = append(payloads, &pb.Message_Payload{
				Module: "tags." + s.cfg.EngineName,
				Body:   toJSON(parsedDet.Family),
			})
		}
	}

	msg := &pb.Message{Sha256: sha256, Payload: payloads}
	out, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("failed to pb marshal: %v", err)
		return err
	}

	s.pub.Publish(ctx, s.cfg.Producer.Topic, out)

	return nil
}
