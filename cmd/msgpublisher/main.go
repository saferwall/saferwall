package main

import (
	"context"
	"flag"
	"os"

	cfg "github.com/saferwall/saferwall/internal/config"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/pubsub/nsq"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

// Config represents our application config.
type Config struct {
	LogLevel           string `mapstructure:"log_level"`
	SharedVolume       string `mapstructure:"shared_volume"`
	Nsqd               string `mapstructure:"nsqd"`
	OrchestratorTopic  string `mapstructure:"orchestrator_topic"`
	MetaTopic          string `mapstructure:"meta_topic"`
	PostProcessorTopic string `mapstructure:"postprocessor_topic"`
	PETopic            string `mapstructure:"pe_topic"`
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

	switch svc {
	case "orchestrator":
		msg := []byte(sha256)
		topic = c.OrchestratorTopic
		pub.Publish(ctx, topic, msg)
	case "meta":
		msg := []byte(sha256)
		topic = c.MetaTopic
		pub.Publish(ctx, topic, msg)
	case "pe":
		msg := []byte(sha256)
		topic = c.PETopic
		pub.Publish(ctx, topic, msg)
	case "postprocessor":
		msg := []byte(sha256)
		topic = c.PostProcessorTopic
		pub.Publish(ctx, topic, msg)
	}

	logger.Infof("Message has been produced to topic: %s", topic)

	return nil
}
