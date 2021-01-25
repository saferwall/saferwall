// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package drweb

import (
	"regexp"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cmd      = "/opt/drweb.com/bin/drweb-ctl"
	regexStr = "infected with (.*)"

	// ConfigD represents the configuration daemon
	ConfigD = "/opt/drweb.com/bin/drweb-configd"
)

// Result represents detection results
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// Version represents all DrWeb components' versions
type Version struct {
	CoreEngineVersion string `json:"core_engine_version"`
}

// GetVersion returns ScanCL, Core and VDF versions
func GetVersion() (Version, error) {

	out, err := utils.ExecCommand(cmd, "baseinfo")

	// Core engine version: 7.00.47.04280
	// Virus database timestamp: 2020-Aug-11 18:40:16
	// Virus database fingerprint: D2EFA560783BC31243E97B3B73766C18
	// Virus databases loaded: 202
	// Virus records: 9118543
	// Anti-spam core is not loaded
	// Last successful update: 2020-Aug-11 20:42:37
	// Next scheduled update: 2020-Aug-11 21:12:37
	if err != nil {
		return Version{}, err
	}

	ver := Version{}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Core engine version:") {
			ver.CoreEngineVersion = strings.TrimSpace(strings.TrimPrefix(line, "Core engine version:"))
			break
		}
	}

	return ver, nil
}

// ScanFile scans a given file
func ScanFile(filepath string) (Result, error) {

	// # /opt/drweb.com/bin/drweb-ctl scan --help
	// Scan file or directory
	// Usage: drweb-ctl scan <path_to_scan> [options]
	// Available options:
	//   -a [ --Autonomous ]               start autonomous component set
	//   --Report arg (=BRIEF)             report type BRIEF, DEBUG or JSON
	//   --ScanTimeout arg (=0)            scan timeout (in ms), 0 means no timeout
	//   --PackerMaxLevel arg (=8)         limit packer nesting level
	//   --ArchiveMaxLevel arg (=8)        limit archive (like zip) nesting level
	//   --MailMaxLevel arg (=8)           limit mail (like pst, tbb) nesting level
	//   --ContainerMaxLevel arg (=8)      limit container (like html) nesting level
	//   --MaxCompressionRatio arg (=3000) limit compression ratio (must be >= 2)
	//   --HeuristicAnalysis arg (=ON)     use heuristic analysis ON, OFF
	//   --Exclude arg                     exclude specified paths from scan
	// 									(wildcards are allowed)
	//   --OnKnownVirus arg (=REPORT)      action REPORT, CURE, QUARANTINE, DELETE
	//   --OnIncurable arg (=REPORT)       action REPORT, QUARANTINE, DELETE
	//   --OnSuspicious arg (=REPORT)      action REPORT, QUARANTINE, DELETE
	//   --OnAdware arg (=REPORT)          action REPORT, QUARANTINE, DELETE
	//   --OnDialers arg (=REPORT)         action REPORT, QUARANTINE, DELETE
	//   --OnJokes arg (=REPORT)           action REPORT, QUARANTINE, DELETE
	//   --OnRiskware arg (=REPORT)        action REPORT, QUARANTINE, DELETE
	//   --OnHacktools arg (=REPORT)       action REPORT, QUARANTINE, DELETE
	//   --Stdin                           read '\n'-separated paths from stdin
	//   --Stdin0                          read '\0'-separated paths from stdin
	//   -d [ --Debug ]                    extended diagnostic output
	out, err := utils.ExecCommand(cmd, "scan", filepath)

	// # /opt/drweb.com/bin/drweb-ctl scan /eicar
	// /eicar - infected with EICAR Test File (NOT a Virus!)
	// Scanned objects: 1, scan errors: 0, threats found: 1, threats neutralized: 0.
	// Scanned 0.07 KB in 0.08 s with speed 0.80 KB/s.
	//
	// 	List of return codes :
	// 		1: Error on monitor channel
	// 		2: Operation is already in progress
	// 		3: Operation is in pending state
	// 	    4: Interrupted by user
	// 	    5: Operation canceled
	// 	    6: IPC connection terminated
	// 	    7: Invalid IPC message size
	// 	    8: Invalid IPC message format
	// 	    9: Not ready
	// 	   10: The component is not installed
	// 	   11: Unexpected IPC message
	// 	   12: IPC protocol violation
	// 	   13: Subsystem state is unknown
	// 	   20: Path must be absolute
	// 	   21: Not enough memory
	// 	   22: IO error
	// 	   23: No such file or directory
	// 	   24: Permission denied
	// 	   25: Not a directory
	// 	   26: Data file corrupted
	// 	   27: File already exists
	// 	   28: Read-only file system
	// 	   29: Network error
	// 	   30: Not a drive
	// 	   31: Unexpected EOF

	if err != nil {
		return Result{}, err
	}

	// Grab the detection result
	re := regexp.MustCompile(regexStr)
	l := re.FindStringSubmatch(out)
	res := Result{}
	if len(l) > 0 {
		res.Output = l[1]
		res.Infected = true
	}
	return res, nil
}
