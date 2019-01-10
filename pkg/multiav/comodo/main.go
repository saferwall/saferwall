// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package comodo

import (
	"errors"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

// Our consts
const (
	cmdscan = "/opt/COMODO/cmdscan"
	cavver  = "/opt/COMODO/cavver.dat"
)

// Result represents detection results
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// GetProgramVersion returns COMODO Anti-Virus version
func GetProgramVersion() (string, error) {

	// Read the content of the file to grab the version
	version, err := utils.ReadAll(cavver)
	if err != nil {
		return "", err
	}

	return string(version), nil
}

// ScanFile a file with COMODO scanner
func ScanFile(filePath string) (Result, error) {

	// Run now
	cmdscanOut, err := utils.ExecCommand(cmdscan, "-v", "-s", filePath)
	// -s: scan a file or directory
	// -v: verbose mode, display more detailed output
	lines := strings.Split(cmdscanOut, "\n")

	res := Result{}
	if err != nil {
		return res, err
	}

	// -----== Scan Start ==-----
	// /home/noteworthy/malware/virut ---> Found Virus, Malware Name is Packed.Win32.MUPX.Gen
	// -----== Scan End ==-----
	// Number of Scanned Files: 1
	// Number of Found Viruses: 1

	lines = strings.Split(cmdscanOut, "\n")
	if len(lines) == 0 {
		return res, errors.New("we got an empty output")
	}

	// Check if file is infected
	if strings.HasSuffix(lines[1], "---> Not Virus") {
		return res, nil
	}

	// Grab detection name
	// Viruses found: 1
	// Virus name:       Trojan-Ransom.Win32.Locky.d
	// Infected objects: 1
	detection := strings.Split(lines[1], "Malware Name is ")
	res.Output = detection[len(detection)-1]
	res.Infected = true
	return res, nil
}
