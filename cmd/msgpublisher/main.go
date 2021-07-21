package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	cfg "github.com/saferwall/saferwall/pkg/config"
	"github.com/saferwall/saferwall/pkg/log"
	"github.com/saferwall/saferwall/pkg/pubsub/nsq"
	"github.com/saferwall/saferwall/services/pe"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

// Config represents our application config.
type Config struct {
	LogLevel          string `mapstructure:"log_level"`
	SharedVolume      string `mapstructure:"shared_volume"`
	Nsqd              string `mapstructure:"nsqd"`
	OrchestratorTopic string `mapstructure:"orchestrator_topic"`
	MLTopic           string `mapstructure:"ml_topic"`
}

func main() {
	flagSvcName := flag.String("service", "", "Service name to write to. (Required)")
	flagSHA256 := flag.String("sha256", "", "Hash of the file to scan. (Required)")
	flagConfig := flag.String("config", "./../../configs/msgpublisher",
		"path to the config file")
	flag.Parse()

	if *flagSvcName == "" || *flagSHA256 == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	logger := log.New().With(context.TODO(), "version", Version, "sha256",
		*flagSHA256)

	if err := run(logger, *flagConfig, *flagSHA256, *flagSvcName); err != nil {
		logger.Errorf("failed to run the server: %s", err)
		os.Exit(-1)
	}
}

func run(logger log.Logger, configFile, sha256, svc string) error {

	c := Config{}
	var topic string
	ctx := context.Background()
	env := os.Getenv("SAFERWALL_DEPLOYMENT_KIND")

	logger.Infof("loading %s configuration from %s", env, configFile)

	err := cfg.Load(configFile, env, &c)
	if err != nil {
		return err
	}

	pub, err := nsq.NewPublisher(c.Nsqd)
	if err != nil {
		return err
	}

	filePath := filepath.Join(c.SharedVolume, sha256)

	switch svc {
	case "orchestrator":
		msg := []byte(sha256)
		topic = c.OrchestratorTopic
		pub.Publish(ctx, topic, msg)
	case "ml":
		msg, err := pe.Scan(filePath, sha256)
		if err != nil {
			return err
		}
		pub.Publish(ctx, c.MLTopic, msg)
	}

	logger.Infof("Message has been produced to topic: %s", topic)

	return nil
}
