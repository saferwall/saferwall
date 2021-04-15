// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/saferwall/saferwall/pkg/consumer"
)

func main() {

	// Create a new consumer.
	c, err := consumer.New()
	if err != nil {
		log.Fatal(err)
	}

	// Let's allow our queues to drain properly during shutdown.
	// We'll create a channel to listen for SIGINT (Ctrl+C) to signal
	// to our application to gracefully shutdown.
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)

	// This is our main loop. It will continue to read off of our nsq
	// channel until either the consumer dies or our application is signaled
	// to stop.
	for {
		select {
		case <-c.StopChan:
			return // uh oh consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			c.Stop()
		}
	}
}
