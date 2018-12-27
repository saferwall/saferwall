// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avira

import (
	"regexp"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cmd = "/opt/avira/scancl"
)

// Version represents all avira components' versions
type Version struct {
	ScanCLVersion string `json:"scancl_version"`
	CoreVersion   string `json:"core_version"`
	VDFVersion    string `json:"vdf_version"`
}

// GetVersion get Anti-Virus scanner version
func GetVersion() (Version, error) {

	versionOut, err := utils.ExecCommand(cmd, "--version")

	// Avira / Linux Version 1.9.161.2
	// Copyright (c) 2010 by Avira GmbH
	// All rights reserved.
	
	// operating system:   Linux
	// architecture:       ia32
	// system date:        Dec 27 2018
	// scancl Version:     1.9.161.2 
	// core Version:       1.9.2.0 
	// VDF Version:        7.15.16.96 
	
	if err != nil {
		return Version{}, err
	}

	v := Version{}
	lines := strings.Split(versionOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "scancl Version:") {
			v.ScanCLVersion = strings.TrimSpace(strings.TrimPrefix(line, "scancl Version:"))
		} else if strings.Contains(line, "core Version:") {
			v.CoreVersion = strings.TrimSpace(strings.TrimPrefix(line, "core Version:"))
		} else if strings.Contains(line, "VDF Version:") {
			v.VDFVersion = strings.TrimSpace(strings.TrimPrefix(line, "VDF Version:"))
		}
	}

	return v, nil
}

// ScanFile scans a given file
func ScanFile(filepath string) (string, error) {

	// Execute the scanner with the given file path
	// --nombr ................  do not check any master boot records
	// --nostats ..............  do not display scan statistics
	// --quarantine ...........  set the quarantine directory

	scanclOut, err := utils.ExecCommand(cmd, "--nombr", "--nostats",
		"--quarantine=/tmp", filepath)

	// From Avira documentation
	// 0 Normal program termination, no detection, no error
	// 1 Found concerning file or boot sector
	// 2 A signature was found in memory
	// 3 Suspicious file found
	// 100 Avira has only displayed the help text
	// 101 A macro was found in a document file
	// 20? Program aborted with one of the following error codes:
	// 203 Invalid option
	// 204 Invalid (nonexistent) directory given in the command line
	// 205 The log file could not be created
	// 210 Avira could not find a necessary library file
	// 211 Program aborted, because the self-check failed
	// 212 The virus definition files could not be read
	// 213 An error occurred during initialization (engine and VDF
	// versions incompatible)
	// 214 No valid license found
	// 215 ScanCL self-test failed
	// 216 File access denied (no permissions)
	// 217 Directory access denied (no permissions)
	if err != nil && err.Error() != "exit status 1" && err.Error() != "exit status 2" &&
		err.Error() != "exit status 3" && err.Error() != "exit status 101" {
		return "", err
	}

	// ./scan./scancl --quarantine=/tmp  --nostats /home/ubuntu/malware/locky
	// Avira / Linux Version 1.9.161.2
	// Copyright (c) 2010 by Avira GmbH
	// All rights reserved.

	// engine set:         8.3.52.150
	// VDF Version:        7.15.16.96

	// key file:           /opt/avira/hbedv.key
	// registered user:    Free
	// serial number:      0000149996
	// key expires:        Dec 31 2999

	// Scan start time: Thu Dec 27 11:13:42 2018
	// Command line: ./scancl --quarantine=/tmp --nostats /home/noteworthy/malware/locky

	// auto excluding /sys from scanning (is a special fs)
	// auto excluding /proc from scanning (is a special fs)
	// configuration file: /opt/avira/scancl.conf
	// ALERT: [TR/Agent.53465] /home/noteworthy/malware/locky <<< Is the Trojan horse TR/Agent.53465

	// Grab the detection result
	detection := ""
	re := regexp.MustCompile(" ALERT: \\[(.*)\\] ")
	l := re.FindStringSubmatch(scanclOut)
	if len(l) > 0 {
		detection = l[1]
	}
	return detection, nil
}
