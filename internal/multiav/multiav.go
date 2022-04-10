// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package multiav

import "errors"

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// Result represents detection results.
type Result struct {
	// Infected is true when the file has been detected as so by the antivirus.
	Infected bool `json:"infected"`
	// The detection name.
	Output string `json:"output"`
	// Out represent the std out from the av scanner during the cmd line scan.
	Out string `json:"-"`
}

var (
	ErrParseDetection = errors.New("failed to parse detection name")
)
