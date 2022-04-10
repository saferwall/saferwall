// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package hasher

import (
	"encoding/hex"
	"hash"
)

// Hasher is an interface to abstract hash calculation.
type Hasher interface {
	Hash(b []byte) string
}

// Service represents the password reset token management service.
type Service struct {
	h hash.Hash
}

// New initializes the token generation service.
func New(h hash.Hash) Service {
	return Service{h}
}

// Hash hashes a stream of bytes using sha2 algorihtm.
func (s Service) Hash(b []byte) string {
	s.h.Reset()
	s.h.Write(b)
	return hex.EncodeToString(s.h.Sum(nil))
}
