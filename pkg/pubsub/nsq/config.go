// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package nsq

import (
	gonsq "github.com/nsqio/go-nsq"
)

// Config holds the basic information for working with NSQ.
type Config struct {
	NsqdAddr string `mapstructure:"nsqd_addr"`
	// the data source name (DSN) for connecting to the broker server.
	NsqLookupds []string `mapstructure:"nsqlookupds"`
	// Topic name to consume from or to produce to.
	Topic       string `mapstructure:"topic"`
	Channel     string `mapstructure:"channel"`
	Concurrency int    `mapstructure:"concurrency"`
	// NSQ Config for more control over the underlying nsq client.
	Config *gonsq.Config
}
