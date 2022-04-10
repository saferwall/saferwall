// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package comodo

import (
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	cmdscan = "/opt/COMODO/cmdscan"
	cavver  = "/opt/COMODO/cavver.dat"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// ProgramVersion returns COMODO Anti-Virus version.
func ProgramVersion() (string, error) {

	// Read the content of the file to grab the version
	version, err := utils.ReadAll(cavver)
	if err != nil {
		return "", err
	}

	return string(version), nil
}

// ScanFile a file with COMODO scanner.
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	res.Out, err = utils.ExecCmd(cmdscan, "-v", "-s", filePath)
	// -s: scan a file or directory
	// -v: verbose mode, display more detailed output
	if err != nil {
		return res, err
	}

	// -----== Scan Start ==-----
	// /samples/virut ---> Found Virus, Malware Name is Packed.Win32.MUPX.Gen
	// -----== Scan End ==-----
	// Number of Scanned Files: 1
	// Number of Found Viruses: 1
	lines := strings.Split(res.Out, "\n")
	if len(lines) < 2 {
		return res, multiav.ErrParseDetection
	}

	// Check if file is infected.
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
