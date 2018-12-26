// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package bitdefender

import (
	"errors"
	"regexp"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

// Our consts
const (
	bdscan = "/opt/BitDefender-scanner/bin/bdscan"
)

// Result represents detection results
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// GetProgramVersion returns Bitdefender Anti-Virus version
func GetProgramVersion() (string, error) {

	//  Run now
	versionOut, err := utils.ExecCommand(bdscan, "--version")
	if err != nil {
		return "", err
	}
	// BitDefender Antivirus Scanner for Unices v7.141118 Linux-amd64
	// Copyright (C) 1996-2014 BitDefender. All rights reserved.

	version := ""
	lines := strings.Split(versionOut, "\n")
	if len(lines) > 0 {
		re := regexp.MustCompile("\\w+ v([\\d.]+) .*")
		version = re.FindStringSubmatch(lines[0])[1]
	}
	return version, nil

}

// ScanFile a file with COMODO scanner
func ScanFile(filePath string) (Result, error) {

	//  Run now
	bdscanOut, err := utils.ExecCommand(bdscan, "--action=ignore", filePath)
	// --action=[disinfect|quarantine|delete|ignore]

	if err != nil && err.Error() != "exit status 1" {
		return Result{}, err
	}

	// BitDefender Antivirus Scanner for Unices v7.141118 Linux-amd64
	// Copyright (C) 1996-2014 BitDefender. All rights reserved.

	// Infected file action: ignore
	// Suspected file action: ignore
	// Loading plugins, please wait
	// Plugins loaded.

	// /home/linux/malware/locky  infected: Trojan.GenericKD.3048400

	lines := strings.Split(bdscanOut, "\n")
	if len(lines) == 0 {
		return Result{}, errors.New("we got an empty output")
	}

	//  Grab detection name
	r := Result{}
	for _, line := range lines {
		if strings.Contains(line, "infected: ") {
			parts := strings.Split(line, "infected: ")
			r.Output = parts[len(parts)-1]
			r.Infected = true
			break
		}
	}

	return r, nil
}
