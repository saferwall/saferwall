// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package clamav

import (
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	clamdscan = "/usr/bin/clamdscan"
)

// ScanFile performs antivirus scan
func ScanFile(filePath string) (string, error) {

	// Execute the scanner with the given file path
	// --no-summary   Disable summary at end of scanning
	clamscanOut, err := utils.ExecCommand(clamdscan, "--no-summary", filePath)

	// clamscan return values (documented from man clamscan)
	//  0 : No virus found.
	//  1 : Virus(es) found.
	// 40: Unknown option passed.
	// 50: Database initialization error.
	// 52: Not supported file type.
	// 53: Can't open directory.
	// 54: Can't open file. (ofm)
	// 55: Error reading file. (ofm)
	// 56: Can't stat input file / directory.
	// 57: Can't get absolute path name of current working directory.
	// 58: I/O error, please check your file system.
	// 62: Can't initialize logger.
	// 63: Can't create temporary files/directories (check permissions).
	// 64: Can't write to temporary directory (please specify another one).
	// 70: Can't allocate memory (calloc).
	// 71: Can't allocate memory (malloc).
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}

	// samples/locky: Win.Malware.Locky-5540 FOUND
	// samples/putty: OK
	infected := false
	if strings.HasSuffix(clamscanOut, "OK\n") {
		infected = false
	} else if strings.HasSuffix(clamscanOut, "FOUND\n") {
		infected = true
	}

	// Extract detection name if infected
	detection := ""
	if infected {
		parts := strings.Split(clamscanOut, ": ")
		detection = parts[len(parts)-1]
		detection = strings.TrimSuffix(detection, " FOUND\n")
	}

	return detection, nil
}

// GetVersion returns program version
func GetVersion() (string, error) {

	// Execute the scanner with the given file path
	versionOut, err := utils.ExecCommand(clamdscan, "--version")
	if err != nil {
		return "", err
	}

	// Extract the version
	version := strings.Split(versionOut, "/")[0]
	version = strings.Split(version, " ")[1]
	return version, nil
}
