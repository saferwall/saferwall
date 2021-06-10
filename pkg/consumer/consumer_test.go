// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.
package consumer

import (
	"log"
	"testing"
)

func TestConsumer(t *testing.T) {
	t.Run("TestNew", func(t *testing.T) {
		c, err := New()
		defer c.Stop()
		if err != nil {
			log.Fatal(err)
		}
	})
}
