// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"github.com/saferwall/saferwall/internal/behavior"
)

// EventType is the type of the system event. A type can be either:
// `registry`, `network` or `file`.
type EventType string

// Event represents a system event: a registry, network or file event.
type Event struct {
	// Process identifier responsible for generating the event.
	ProcessID string `json:"pid"`
	// Type of the system event.
	Type EventType `json:"type"`
	// Path of the system event. For instance, when the event is of type:
	// `registry`, the path represents the registry key being used. For a
	// `network` event type, the path is the IP or domain used.
	Path string `json:"path"`
	// Th operation requested over the above `Path` field. This field means
	// different things according to the type of the system event.
	// - For file system events: can be either: create, read, write, delete, rename, ..
	// - For registry events: can be either: create, rename, set, delete.
	// - For network events: this represents the protocol of the communication, can
	// be either HTTP, HTTPS, FTP, FTP
	Operation string `json:"op"`
}

func generateEvents(events []behavior.Event) []Event {
	curatedEvents := make([]Event, 0)
	for _, evt := range events {
		curatedEvents = append(curatedEvents, Event{
			ProcessID: evt.ProcessID,
			Type:      EventType(evt.Type),
			Path:      evt.Path,
			Operation: evt.Operation,
		})
	}

	return curatedEvents
}
