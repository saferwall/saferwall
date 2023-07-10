// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

type ProcessTree struct {
	// Process image's path.
	ImagePath string `json:"path"`
	// Process identifier.
	ProcessID string `json:"pid"`
	// The name of the process.
	ProcessName string `json:"proc_name"`
	// The file type: doc, exe, etc.
	FileType string `json:"file_type"`
	// Detection contains the family name of the malware if it is malicious,
	// or clean otherwise.
	Detection string `json:"detection"`
	// The children of the process.
	Children *ProcessTree `json:"children"`
}
