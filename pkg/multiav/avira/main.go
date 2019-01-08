// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avira

import (
	"errors"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/saferwall/saferwall/pkg/utils"
)

var (
	// ErrNoLicenseFound is returned when no license is found
	ErrNoLicenseFound = errors.New("No License Found")

	// ErrInvalidLicense is returned when invalid license is found
	ErrInvalidLicense = errors.New("Invalid License")
)

const (
	cmd = "/opt/avira/scancl"

	// LicenseKeyPath points to the location of the license.
	LicenseKeyPath = "/opt/avira/hbedv.key"
)

// Result represents detection results
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// Version represents all avira components' versions
type Version struct {
	ScanCLVersion string `json:"scancl_version"`
	CoreVersion   string `json:"core_version"`
	VDFVersion    string `json:"vdf_version"`
}

// GetVersion returns ScanCL, Core and VDF versions
func GetVersion() (Version, error) {

	out, err := utils.ExecCommand(cmd, "--version")

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

	ver := Version{}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "scancl Version:") {
			ver.ScanCLVersion = strings.TrimSpace(strings.TrimPrefix(line, "scancl Version:"))
		} else if strings.Contains(line, "core Version:") {
			ver.CoreVersion = strings.TrimSpace(strings.TrimPrefix(line, "core Version:"))
		} else if strings.Contains(line, "VDF Version:") {
			ver.VDFVersion = strings.TrimSpace(strings.TrimPrefix(line, "VDF Version:"))
		}
	}

	return ver, nil
}

// ScanFile scans a given file
func ScanFile(filepath string) (Result, error) {

	// Execute the scanner with the given file path
	// --nombr ................  do not check any master boot records
	// --nostats ..............  do not display scan statistics
	// --quarantine ...........  set the quarantine directory

	out, err := utils.ExecCommand(cmd, "--nombr", "--nostats",
		"--quarantine=/tmp", filepath)

	// 	List of return codes :
	// 		0: Normal program termination, nothing found, no error
	// 		1: Found concerning file(s) or boot sector(s)
	// 		2: A signature was found in memory
	// 		3: Suspicious file(s) found
	// 	  100: Avira only has displayed this help text
	// 	  101: A macro was found in a document file

	//  Abort error codes:
	// 	  203: Invalid option
	// 	  204: Invalid (nonexistent) directory given at command line
	// 	  205: The log file could not be created
	// 	  210: Avira could not find a necessary library file
	// 	  211: Program aborted because the self check failed
	// 	  212: Virus definition file(s) could not be read
	// 	  213: An error occurred during initialization (engine and VDF versions incompatible)
	// 	  214: No valid license found
	// 	  215: Scancl self test failed
	// 	  216: File access denied (no permissions)
	// 	  217: Directory access denied (no permissions)

	res := Result{}
	if err != nil && err.Error() != "exit status 1" && err.Error() != "exit status 2" &&
		err.Error() != "exit status 3" && err.Error() != "exit status 101" {
		return res, err
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
	re := regexp.MustCompile(" ALERT: \\[(.*)\\] ")
	l := re.FindStringSubmatch(out)
	if len(l) > 0 {
		res.Output = l[1]
		res.Infected = true
	}
	return res, nil
}

// IsLicenseExpired returns true if license was expired
func IsLicenseExpired() (bool, error) {
	out, err := utils.ExecCommand(cmd, "-v")
	if err != nil {
		return true, err
	}

	if strings.Contains(out, "No license found") {
		return true, ErrNoLicenseFound
	} else if strings.Contains(out, "invalid license") {
		return true, ErrInvalidLicense
	} else if strings.Contains(out, "key expires") {
		// key file:           /opt/avira/hbedv.key
		// registered user:    Free
		// serial number:      0000149996
		// key expires:        Dec 31 2018

		re := regexp.MustCompile(`key expires:        ([\w\s]+)\n\n`)
		l := re.FindStringSubmatch(out)
		if len(l) > 0 {
			expiresAt, err := time.Parse("Jan 02 2006", l[1])
			if err != nil {
				return true, err
			}
			now := time.Now()
			diff := expiresAt.Sub(now)
			if diff < 0 {
				return true, nil
			}
		}
	}

	return false, nil
}

// ActivateLicense activate the license.
func ActivateLicense(r io.Reader) error {
	// Write the license file to disk
	_, err := utils.WriteBytesFile(LicenseKeyPath, r)
	if err != nil {
		return err
	}

	isExpired, err := IsLicenseExpired()
	if err != nil {
		return err
	}

	if isExpired {
		return errors.New("License was expird")
	}

	return nil
}
