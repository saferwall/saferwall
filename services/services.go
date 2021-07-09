// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package micro

// Message represents the message format shared across the different services/
type Message struct {
	Sha256 string
	Body   []byte
}
