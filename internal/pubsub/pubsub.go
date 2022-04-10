// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pubsub

import (
	"golang.org/x/net/context"
)

// Publisher is a generic interface to encapsulate how we want our publishers
// to behave.
type Publisher interface {
	// Publish will publish a message with context.
	Publish(context.Context, string, []byte) error
}

// Subscriber is a generic interface to encapsulate how we want our subscribers
// to behave. For now the system will auto stop if it encounters any errors. If
// a user encounters a closed channel, they should check the Err() method to see
// what happened.
type Subscriber interface {
	// Start will return a channel of raw messages.
	Start() error
	// Err will contain any errors returned from the consumer connection.
	Err() error
	// Stop will initiate a graceful shutdown of the subscriber connection.
	Stop() error
}
