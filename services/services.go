// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package micro

// FileScanProgressType represents the file scan progress type.
type FileScanProgressType uint8

// Progress of a file scan.
const (
	FileScanProgressQueued     FileScanProgressType = 1
	FileScanProgressProcessing FileScanProgressType = 2
	FileScanProgressFinished   FileScanProgressType = 3
)

// DocMetadata stores metadata information for saved documents in the DB.
type DocMetadata struct {
	CreatedAt   int64 `json:"created_at,omitempty"`
	LastUpdated int64 `json:"last_updated,omitempty"`
	Version     int   `json:"version,omitempty"`
}
