// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"path/filepath"
	"strings"
)

// ProcessTree represents an array of processes, each process contains a process
// ID that can helps us track the parent/children relationship.
type ProcessTree []Process

// Process represents a process object within the detonation context.
// This structure help us build the process tree.
type Process struct {
	// Process image's path.
	ImagePath string `json:"path"`
	// Process identifier.
	PID string `json:"pid"`
	// The parent process ID.
	ParentPID string `json:"parent_pid"`
	// The relationship between this process and its parent.
	ParentLink string `json:"parent_link"`
	// The name of the process.
	ProcessName string `json:"proc_name"`
	// The file type: doc, exe, etc.
	FileType string `json:"file_type"`
	// Detection contains the family name of the malware if it is malicious,
	// or clean otherwise.
	Detection string `json:"detection"`
}

func enrichProcTree(procTree ProcessTree) ProcessTree {
	var enrichedProcTree ProcessTree
	for _, p := range procTree {
		p.ProcessName = filepath.Base(p.ImagePath)
		fileExt := filepath.Ext(p.ImagePath)
		p.FileType = strings.Replace(fileExt, ".", "", -1)
		enrichedProcTree = append(enrichedProcTree, p)
	}
	return enrichedProcTree
}
