// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package nsq

import (
	"encoding/json"
	"errors"
	"os"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/pkg/pubsub"
	"golang.org/x/net/context"
)

// publisher is a publisher that provides an implementation for NSQ.
type publisher struct {
	producer *gonsq.Producer
	topic    string
}

// NoopNSQLogger allows us to pipe NSQ logs to dev/null
// The default NSQ logger is great for debugging, but did
// not fit our normally well structured JSON logs. Luckily
// NSQ provides a simple interface for injecting your own
// logger.
type NoopNSQLogger struct{}

// Output allows us to implement the nsq.Logger interface
func (l *NoopNSQLogger) Output(calldepth int, s string) error {
	return nil
}

// NewPublisher will initiate a new nsq producer.
func NewPublisher(cfg *Config) (pubsub.Publisher, error) {
	var err error
	p := &publisher{}

	if len(cfg.Topic) == 0 {
		return p, errors.New("topic name is required")
	}
	nsqConfig := cfg.Config
	if nsqConfig == nil {
		nsqConfig = gonsq.NewConfig()
	}
	p.topic = cfg.Topic
	p.producer, err = gonsq.NewProducer(cfg.NsqdAddr, nsqConfig)
	return p, err
}

// Publish will marshal the message to json and produce it to the NSQ topic.
func (p *publisher) Publish(ctx context.Context, key string, m []byte) error {
	msg, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return p.producer.Publish(p.topic, msg)
}

// Stop will close the pub connection.
func (p *publisher) Stop() error {
	p.producer.Stop()
	return nil
}

// subscriber is a subscriber that provides an implementation for NSQ.
type subscriber struct {
	cons        *gonsq.Consumer
	nsqlookupds []string
	topic       string
	channel     string
	concurrency int
	handler     gonsq.Handler
	stop        chan chan error
	nerr        error
}

// NewSubscriber will initiate an nsq consumer.
func NewSubscriber(cfg *Config) (pubsub.Subscriber, error) {

	var err error
	s := &subscriber{
		topic: cfg.Topic,
		stop:  make(chan chan error, 1),
	}

	if len(cfg.NsqLookupds) == 0 {
		return s, errors.New("at least 1 broker host is required")
	}
	if len(cfg.Topic) == 0 {
		return s, errors.New("topic name is required")
	}

	nsqConfig := cfg.Config
	if nsqConfig == nil {
		nsqConfig = gonsq.NewConfig()
	}
	s.topic = cfg.Topic
	s.channel = cfg.Channel
	s.nsqlookupds = cfg.NsqLookupds
	s.cons, err = gonsq.NewConsumer(cfg.Topic, cfg.Channel, nsqConfig)
	return s, err
}

// Start will start consuming message on the NSQ topic. If it encounters any
// issues, it will populate the Err() error and close the returned channel.
func (s *subscriber) Start() <-chan pubsub.SubscriberMessage {

	// Here we set the logger to our NoopNSQLogger to quiet down the default logs.
	// At Reverb we use a custom structured logging format so we'll take the
	// logging from here.
	s.cons.SetLogger(&NoopNSQLogger{}, gonsq.LogLevelError)

	// Injects our handler into the consumer. You'll define one handler
	// per consumer, but you can have as many concurrently running handlers
	// as specified by the second argument. If your MaxInFlight is less
	// than your number of concurrent handlers you'll starve your workers
	// as there will never be enough in flight messages for your worker pool
	s.cons.AddConcurrentHandlers(s.handler, s.concurrency)

	// Our consumer will discover where topics are located by our three
	// nsqlookupd instances The application will periodically poll
	// these nqslookupd instances to discover new nodes or drop unhealthy
	// producers.
	s.cons.ConnectToNSQLookupds(s.nsqlookupds)

	// Let's allow our queues to drain properly during shutdown.
	// We'll create a channel to listen for SIGINT (Ctrl+C) to signal
	// to our application to gracefully shutdown.
	shutdown := make(chan os.Signal, 2)

	output := make(chan pubsub.SubscriberMessage)
	go func(s *subscriber, output chan pubsub.SubscriberMessage) {
		defer close(output)

		// channel until either the consumer dies or our application is signaled
		// to stop.
		for {
			select {
			case <-s.cons.StopChan:
				// consumer disconnected. Time to quit.
				return
			case <-shutdown:
				// Synchronously drain the queue before falling out of main
				s.Stop()
				return
			}
		}
	}(s, output)

	return output
}

// Stop will block until the consumer has stopped consuming messages
// and return any errors seen on consumer close.
func (s *subscriber) Stop() error {
	s.cons.Stop()
	return nil
}

// Err will contain any  errors that occurred during
// consumption. This method should be checked after
// a user encounters a closed channel.
func (s *subscriber) Err() error {
	return s.nerr
}
