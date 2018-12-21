// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avast

import (
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cmd = "/bin/scan"
)

// GetVPSVersion returns Avast VPS version
func GetVPSVersion() (string, error) {

	// Run the scanner to grab the version
	vpsOut, err := utils.ExecCommand(cmd, "-V")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(vpsOut), nil
}

// GetProgramVersion returns Avast Program version
func GetProgramVersion() (string, error) {

	// Run the scanner to grab the version
	versionOut, err := utils.ExecCommand(cmd, "-v")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(versionOut), nil
}

// ScanFile scans a given file
func ScanFile(filepath string) (string, error) {

	// Execute the scanner with the given file path
	// -a         Print all scanned files/URLs, not only infected.
	// -b         Report decompression bombs as infections.
	// -f         Scan full files.
	// -u         Report potentionally unwanted programs (PUP).

	avastOut, err := utils.ExecCommand(cmd, "-abfu", filepath)

	// From Avast linux technical documentation:
	// The exit status is 0 if no infected files are found and 1
	// otherwise. If an error occurred, the exit status is 2.
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}

	// Sanitize the output and return
	str := strings.Split(avastOut, "\t")
	result := strings.TrimSpace(str[1])
	return result, nil
}

// ScanURL scans a given URL
func ScanURL(url string) (string, error) {

	// Execute the scanner with the given URL
	avastOut, err := utils.ExecCommand(cmd, "-U", url)

	// From Avast linux technical documentation:
	// The exit status is 0 if no infected files are found and 1
	// otherwise. If an error occurred, the exit status is 2.
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}

	// Check if we got a clean URL
	if avastOut == "" {
		return "[OK]", nil
	}

	// Sanitize the output and return
	str := strings.Split(avastOut, "\t")
	result := strings.TrimSpace(str[1])
	return result, nil
}

// IsInfected returns true if file is not clean
func IsInfected(result string) bool {
	if !strings.Contains(result, "[OK]") {
		return false
	} else {
		return true
	}
}
