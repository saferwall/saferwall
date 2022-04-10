// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avira

import (
	"errors"
	"io"
	"regexp"
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

var (
	// ErrNoLicenseFound is returned when no license is found
	ErrNoLicenseFound = errors.New("no license found")

	// ErrInvalidLicense is returned when invalid license is found
	ErrInvalidLicense = errors.New("invalid license")

	// ErrExpiredLicense is returned when license is expired
	ErrExpiredLicense = errors.New("license expired")

	// ErrLicenseUnknownError is returned when unknown error occurred
	ErrLicenseUnknownError = errors.New("license parsing failed")
)

const (
	cmd = "/opt/avira/scancl"

	// LicenseKeyPath points to the location of the license.
	LicenseKeyPath = "/opt/avira/hbedv.key"
)

var (
	re = regexp.MustCompile(`ALERT: \[(.*)\] `)
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// Version represents all avira components' versions
type Version struct {
	ScanCLVersion string `json:"scancl_version"`
	CoreVersion   string `json:"core_version"`
	VDFVersion    string `json:"vdf_version"`
}

// GetVersion returns ScanCL, Core and VDF versions
func GetVersion() (Version, error) {

	out, err := utils.ExecCmd(cmd, "--version")

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
			ver.ScanCLVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "scancl Version:"))
		} else if strings.Contains(line, "core Version:") {
			ver.CoreVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "core Version:"))
		} else if strings.Contains(line, "VDF Version:") {
			ver.VDFVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "VDF Version:"))
		}
	}

	return ver, nil
}

// ScanFile scans a given file
func (Scanner) ScanFile(filepath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// Execute the scanner with the given file path
	// --nombr ................  do not check any master boot records
	// --nostats ..............  do not display scan statistics
	// --quarantine ...........  set the quarantine directory

	res.Out, err = utils.ExecCmd(cmd, "--nombr", "--nostats",
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

	if err != nil && err.Error() != "exit status 1" &&
		err.Error() != "exit status 2" && err.Error() != "exit status 3" &&
		err.Error() != "exit status 101" {
		return res, err
	}

	// ./scan./scancl --quarantine=/tmp  --nostats /samples/locky
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
	// Command line: ./scancl --quarantine=/tmp --nostats /samples/locky

	// auto excluding /sys from scanning (is a special fs)
	// auto excluding /proc from scanning (is a special fs)
	// configuration file: /opt/avira/scancl.conf
	// ALERT: [TR/Agent.53465] /samples/locky <<< Is the Trojan horse TR/Agent.53465

	// Grab the detection result
	l := re.FindStringSubmatch(res.Out)
	if len(l) > 0 {
		res.Output = l[1]
		res.Infected = true
	}
	return res, nil
}

// LicenseStatus checks the validity of the license.
func LicenseStatus() (string, error) {
	out, err := utils.ExecCmd(cmd, "-v")
	// exit code 214 means No valid license found.
	if err != nil {
		if err.Error() == "exit status 214" {
			return "", ErrNoLicenseFound
		}

		return "", err
	}

	if strings.Contains(out, "invalid license") {
		return "", ErrInvalidLicense
	} else if strings.Contains(out, "This key has expired") {
		return "", ErrExpiredLicense
	} else if strings.Contains(out, "key expires:") {
		// key file:           /opt/avira/hbedv.key
		// registered user:    Free
		// serial number:      0000149996
		// key expires:        Dec 31 2018

		re := regexp.MustCompile(`key expires:        ([\w\s]+)\n\n`)
		l := re.FindStringSubmatch(out)
		if len(l) > 0 {
			return l[1], nil
		}
	}

	return "", ErrLicenseUnknownError
}

// ActivateLicense activate the license.
func ActivateLicense(r io.Reader) (string, error) {
	// Write the license file to disk
	_, err := utils.WriteBytesFile(LicenseKeyPath, r)
	if err != nil {
		return "", err
	}

	expireAt, err := LicenseStatus()
	if err != nil {
		return "", err
	}

	return expireAt, nil
}
