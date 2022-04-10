// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package bitdefender

import (
	"regexp"
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	bdscan = "/opt/BitDefender-scanner/bin/bdscan"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// ProgramVersion returns Bitdefender Anti-Virus version.
func ProgramVersion() (string, error) {

	out, err := utils.ExecCmd(bdscan, "--version")
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
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	res.Out, err = utils.ExecCmd(bdscan, "--action=ignore", filePath)
	// --action=[disinfect|quarantine|delete|ignore]

	// Exit status codes:
	// 254: Your license has expired.
	if err != nil && err.Error() != "exit status 1" {
		return res, err
	}

	// BitDefender Antivirus Scanner for Unices v7.141118 Linux-amd64
	// Copyright (C) 1996-2014 BitDefender. All rights reserved.

	// Infected file action: ignore
	// Suspected file action: ignore
	// Loading plugins, please wait
	// Plugins loaded.

	// /samples/locky  infected: Trojan.GenericKD.3048400

	lines := strings.Split(res.Out, "\n")
	if len(lines) == 0 {
		return res, nil
	}

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
