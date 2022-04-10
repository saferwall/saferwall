// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package postprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/blackironj/periodic"
	"github.com/djherbis/times"
	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/pubsub"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
	"github.com/saferwall/saferwall/internal/db"
	"github.com/saferwall/saferwall/internal/ml"
	"github.com/saferwall/saferwall/services/config"
	pb "github.com/saferwall/saferwall/services/proto"
	"google.golang.org/protobuf/proto"
)

// Config represents our application config.
type Config struct {
	LogLevel     string             `mapstructure:"log_level"`
	MLAddress    string             `mapstructure:"ml_address"`
	SharedVolume string             `mapstructure:"shared_volume"`
	DB           db.Config          `mapstructure:"db"`
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
	db     db.DB
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

	s.db, err = db.Open(cfg.DB.Server, cfg.DB.Username,
		cfg.DB.Password, cfg.DB.BucketName)
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

	// start a background job that deletes samples from
	// the nfs share that has not been accessed since
	// more than a day.
	scheduler := periodic.NewScheduler()
	cleanupStorage, _ := periodic.NewTask(deleteOldSamples, s.cfg.SharedVolume)
	scheduler.RegisterTask("cleanupStorageTask", time.Hour*24, cleanupStorage)
	scheduler.Run()

	//Stop tasks before program is shutting down
	defer func() {
		scheduler.Stop()
		s.logger.Info("every task is stopped")
	}()
	return nil
}

func toJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func deleteOldSamples(logger log.Logger, root string) {
	var files []string

	logger.Info("run cleanup samples from storage task")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		logger.Errorf("error while walk storage dir: %s, reason: %v", root, err)
		return
	}

	yesterday := time.Now().AddDate(0, 0, -1)

	for _, file := range files {
		t, err := times.Stat(file)
		if err != nil {
			logger.Errorf("error while time stating file: %s, reason: %v", file, err)
			continue
		}
		if t.AccessTime().Before(yesterday) {
			err = os.Remove(file)
			if err != nil {
				logger.Errorf("error while deleting file: %s, reason: %v", file, err)
				continue
			}

			logger.Infof("file: %s has been cleaned up from the storage", file)
		}
	}
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

	// wait until all microservices finishes processing.
	sleeprange := [6]time.Duration{6, 5, 4, 3, 2, 1}
	for _, v := range sleeprange {
		logger.Debugf("Iteratation: %d", v)
		time.Sleep(v * time.Second)
		var multiav map[string]interface{}
		err := s.db.Lookup(ctx, sha256, "multiav.last_scan", &multiav)
		if err != nil {
			logger.Errorf("failed to read document: %v", err)
		}
		logger.Debugf("finish av scanners: %d", len(multiav))
		if len(multiav) == 14 {
			break
		}
	}

	var file map[string]interface{}
	err := s.db.Get(ctx, sha256, &file)
	if err != nil {
		logger.Errorf("failed to read document: %v", err)
		return err
	}

	// set the file analysis status to `finished` and
	// set the `last_scan` to now.
	payloads := []*pb.Message_Payload{
		{Module: "status", Body: toJSON(2)},
		{Module: "last_scanned", Body: toJSON(time.Now().Unix())},
	}

	// if the file format is PE, run the ML classifier.
	if file["fileformat"] == "pe" {
		if _, ok := file["pe"]; ok {
			res, err := ml.PEClassPrediction(s.cfg.MLAddress, toJSON(file))
			if err != nil {
				logger.Errorf("failed to get ml classification results: %v", err)
			} else {
				payloads = append(payloads, &pb.Message_Payload{
					Module: "ml.pe", Body: toJSON(res)})
			}
		}
	}

	// update multiav last_scan and first_scan if needed.
	if _, ok := file["multiav"]; ok {
		logger.Debugf("multiav res: %v", file["multiav"])
		multiav := file["multiav"].(map[string]interface{})
		if _, ok := multiav["first_scan"]; !ok {
			payloads = append(payloads, &pb.Message_Payload{
				Module: "multiav.first_scan",
				Body:   toJSON(multiav["last_scan"])})
		} else if len(multiav["first_scan"].(map[string]interface{})) == 0 {
			payloads = append(payloads, &pb.Message_Payload{
				Module: "multiav.first_scan",
				Body:   toJSON(multiav["last_scan"])})
		} else {
			logger.Debugf("multiav first_scan already set to: %v",
				multiav["first_scan"])
		}
	}

	// serialize the message using protobuf.
	msg := &pb.Message{Sha256: sha256, Payload: payloads}
	out, err := proto.Marshal(msg)
	if err != nil {
		logger.Error("failed to pb marshal: %v", err)
		return err
	}

	// finally, produce the message to the right queue.
	err = s.pub.Publish(ctx, s.cfg.Producer.Topic, out)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
		return err
	}

	return nil
}
