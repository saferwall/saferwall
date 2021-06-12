// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.
package consumer

import (
	"testing"
)

func TestConsumer(t *testing.T) {
	t.Run("TestNew", func(t *testing.T) {
		consumerConfig, err := LoadConfig()
		if err != nil {
			t.Fatal("ConsumerTest failed with error :", err)
		}
		c, err := New(consumerConfig)
		c.Start()
		defer c.Stop()
		if err != nil {
			t.Fatal("ConsumerTest failed with error", err)
		}
	})
}
