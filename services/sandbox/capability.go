// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

// Capability represents any capability found in executable files.
// An example of a capability is: Exfiltration over C2 server.
type Capability struct {
	// Process identifier responsible for generating the capability..
	ProcessID string `json:"pid"`
	// Description describes in a few words the capability.
	Description string `json:"description"`
	// The severity of the capability: informative, suspicious, malicious, etc.
	Severity string `json:"severity"`
	// Category of the capability: Persistence, anti-analysis, etc.
	Category string `json:"category"`
	// The module that generated the capability: yara, behavior, etc.
	Module string `json:"module"`
}
