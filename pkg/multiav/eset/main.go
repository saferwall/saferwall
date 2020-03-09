// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package eset

import (
	"regexp"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cmd = "/opt/eset/esets/sbin/esets_scan"
)

// Result represents detection results.
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// ScanFile performs antivirus scan.
func ScanFile(filePath string) (Result, error) {

	// Execute the scanner with the given file path
	// --unsafe                  scan for potentially unsafe applications
	// --unwanted                scan for potentially unwanted applications
	// --clean-mode=MODE         use cleaning MODE for infected objects.
	// 							 Available options: none, standard (default),
	// 							 strict, rigorous, delete
	out, err := utils.ExecCommand(cmd, "--unsafe", "--unwanted",
		"--clean-mode=NONE", filePath)

	// Exit codes:
	//  0    no threat found
	//  1    threat found and cleaned
	//  10   some files could not be scanned (may be threats)
	//  50   threat found
	//  100  error

	res := Result{}
	if err != nil && err.Error() != "exit status 1" && err.Error() != "exit status 50" {
		return res, err
	}

	// Scan started at:   Tue Jan  1 01:32:57 2019
	// name="/samples/aa.exx", threat="a variant of Win32/Injector.BIIZ trojan", action="", info=""

	re := regexp.MustCompile(`threat="([\s\w/.]+)"`)
	l := re.FindStringSubmatch(out)
	if len(l) < 1 {
		return res, nil
	}

	// Clean up detection name
	det := l[1]
	det = strings.TrimPrefix(det, "a variant of ")
	det = strings.TrimSuffix(det, "  potentially unwanted application")
	det = strings.TrimSuffix(det, "  potentially unsafe application")

	res.Infected = true
	res.Output = det
	return res, nil
}

// GetProgramVersion returns program version.
func GetProgramVersion() (string, error) {

	// Execute the scanner with the given file path
	out, err := utils.ExecCommand(cmd, "--version")
	if err != nil {
		return "", err
	}

	// Extract the version
	ver := strings.Split(out, " ")[2]
	ver = strings.TrimSuffix(ver, "\n")
	return ver, nil
}
