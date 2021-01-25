// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package eset

import (
	"regexp"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cls = "/opt/eset/efs/sbin/cls/cls"
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
	// --no-quarantine           do not copy detected files to Quarantine

	out, err := utils.ExecCommand(cls, "--unsafe", "--unwanted",
		"--clean-mode=NONE", "--no-quarantine", filePath)

	// Exit codes:
	//  0    no threat found
	//  1    threat found and cleaned
	//  10   some files could not be scanned (may be threats)
	//  50   threat found
	//  100  error
	if err != nil && err.Error() != "exit status 1" &&
		err.Error() != "exit status 50" {
		return Result{}, err
	}

	// Scan started at:   Tue Jan  1 01:32:57 2019
	// name="/samples/Wauchos.exe", result="Win32/TrojanDownloader.Wauchos.A trojan", action="retained", info=""
	re := regexp.MustCompile(`result="([\s\w/.]+)"`)
	l := re.FindStringSubmatch(out)
	if len(l) < 1 {
		return Result{}, nil
	}

	// Clean up detection name
	det := l[1]
	det = strings.TrimPrefix(det, "a variant of")
	det = strings.TrimSuffix(det, "potentially unwanted application")
	det = strings.TrimSuffix(det, "potentially unsafe application")
	det = strings.TrimSuffix(det, " trojan")
	det = strings.TrimSuffix(det, " Constructor")
	det = strings.TrimSuffix(det, " worm")
	det = strings.TrimSpace(det)

	res := Result{}
	res.Infected = true
	res.Output = det
	return res, nil
}

// GetProgramVersion returns program version.
func GetProgramVersion() (string, error) {

	// Execute the scanner with the given file path
	out, err := utils.ExecCommand(cls, "--version")
	if err != nil {
		return "", err
	}

	// Extract the version
	ver := strings.Split(out, " ")[2]
	ver = strings.TrimSuffix(ver, "\n")
	return ver, nil
}
