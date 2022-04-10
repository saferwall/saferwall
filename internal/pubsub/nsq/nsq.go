// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package nsq

import (
	"os"
	"time"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/saferwall/saferwall/internal/pubsub"
	"golang.org/x/net/context"
)

// NoopNSQLogger allows us to pipe NSQ logs to dev/null. The default NSQ logger
// is great for debugging, but did not fit our normally well structured JSON
// logs. Luckily NSQ provides a simple interface for injecting your own
// logger.
type NoopNSQLogger struct{}

// Output allows us to implement the nsq.Logger interface.
func (l *NoopNSQLogger) Output(calldepth int, s string) error {
	return nil
}

// publisher is a publisher that provides an implementation for NSQ.
type publisher struct {
	producer *gonsq.Producer
}

// NewPublisher will initiate a new nsq producer.
func NewPublisher(addr string) (pubsub.Publisher,
	error) {
	p, err := gonsq.NewProducer(addr, gonsq.NewConfig())
	return &publisher{p}, err
}

// Publish will marshal the message to json and produce it to the NSQ topic.
func (p *publisher) Publish(ctx context.Context, topic string, msg []byte) error {
	return p.producer.Publish(topic, msg)
}

// Stop will close the pub connection.
func (p *publisher) Stop() error {
	p.producer.Stop()
	return nil
}

// subscriber is a subscriber that provides an implementation for NSQ.
type subscriber struct {
	concurrency int
	nerr        error
	nsqlookupds []string
	handler     gonsq.Handler
	stop        chan chan error
	cons        *gonsq.Consumer
}

// NewSubscriber will initiate an nsq consumer.
func NewSubscriber(topic, channel string, nsqlookupds []string, concurrency int,
	h gonsq.Handler) (pubsub.Subscriber, error) {

	// Create a new nsq config.
	nsqConfig := gonsq.NewConfig()

	// Maximum number of times this consumer will attempt to process a message
	// before giving up.
	nsqConfig.MaxAttempts = 2

	// Maximum number of messages to allow in flight (concurrency knob).
	nsqConfig.MaxInFlight = 10

	// The server-side message timeout for messages delivered to this client.
	nsqConfig.MsgTimeout = time.Duration(2 * time.Minute)

	cons, err := gonsq.NewConsumer(topic, channel, nsqConfig)

	s := subscriber{
		cons:        cons,
		concurrency: concurrency,
		stop:        make(chan chan error, 1),
		handler:     h,
		nsqlookupds: nsqlookupds,
	}

	return &s, err
}

// Start will start consuming message on the NSQ topic. If it encounters any
// issues, it will populate the Err() error and close the returned channel.
func (s *subscriber) Start() error {

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
	if err := s.cons.ConnectToNSQLookupds(s.nsqlookupds); err != nil {
		return err
	}

	// Let's allow our queues to drain properly during shutdown.
	// We'll create a channel to listen for SIGINT (Ctrl+C) to signal
	// to our application to gracefully shutdown.
	shutdown := make(chan os.Signal, 2)

	for {
		select {
		case <-s.cons.StopChan:
			// consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			s.Stop()
		}
	}
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
