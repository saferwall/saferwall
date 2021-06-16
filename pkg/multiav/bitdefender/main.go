// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package bitdefender

import (
	"fmt"
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

// GetProgramVersion returns Bitdefender Anti-Virus version.
func GetProgramVersion() (string, error) {

	//  Run now
	out, err := utils.ExecCommand(bdscan, "--version")
	if err != nil {
		return "", err
	}
	// BitDefender Antivirus Scanner for Unices v7.141118 Linux-amd64
	// Copyright (C) 1996-2014 BitDefender. All rights reserved.

	ver := ""
	lines := strings.Split(out, "\n")
	if len(lines) > 0 {
		re := regexp.MustCompile(`v\d\.\d{6}`)
		match := re.FindStringSubmatch(lines[0])
		ver = match[0]
	}
	return ver, nil

}

// ScanFile a file with Bitdefender scanner.
func ScanFile(filePath string) (Result, error) {

	//  Run now
	out, err := utils.ExecCommand(bdscan, "--action=ignore", filePath)
	// --action=[disinfect|quarantine|delete|ignore]

	// Exit status codes:
	// 254: Your license has expired.
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

	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		errUnexpectedOutput := fmt.Errorf("unexpected output: %s", out)
		return Result{}, errUnexpectedOutput
	}

	//  Grab detection name
	res := Result{}
	for _, line := range lines {
		if strings.Contains(line, "infected: ") {
			parts := strings.Split(line, "infected: ")
			res.Output = parts[len(parts)-1]
			res.Infected = true
			break
		}
	}

	return res, nil
}
