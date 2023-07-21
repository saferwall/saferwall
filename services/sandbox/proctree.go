// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

// ProcessTree represents an array of processes, each process contains a process
// ID that can helps us track the parent/children relationship.
type ProcessTree []Process

// Process represents a process object within the detonation context.
// This structure help us build the process tree.
type Process struct {
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
	// The parent process ID.
	ParentPID string `json:"parent_pid"`
	// True indicates that this is the root node of the tree.
	IsRoot bool `json:"is_root"`
}
