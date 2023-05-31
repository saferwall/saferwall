// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package micro

// Progress of a file scan.
const (
	Queued     = iota + 1
	Processing
	Finished
)
