// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sophos

import (
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	savscan = "/opt/sophos/bin/savscan"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// Version represents all sophos components' versions.
type Version struct {
	ProductVersion       string `json:"product_version"`
	EngineVersion        string `json:"engine_version"`
	VirusDataVersion     string `json:"virus_data_version"`
	UserInterfaceVersion string `json:"user_interface_version"`
}

// ScanFile a file with Sophos scanner
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	//  Scan parameters
	// -f  		: Full Scan
	// -c  		: Ask for confirmation before disinfection/deletion
	// -b  		: Sound bell on virus detection
	// -ss 		: Don't display anything except on error or virus
	// archive  : All of the above
	// loopback : Scan inside loopback-type files
	// mime 	: Scan files encoded in MIME format
	// oe   	: Scan Microsoft Outlook Express mailbox files (requires -mime)
	// tnef 	: Scan inside TNEF files
	// pua : Scan for adware/potentially unwanted applications (PUAs).

	res.Out, err = utils.ExecCmd(savscan, "-f", "-nc", "-nb", "-ss",
		"-archive", "-loopback", "-mime", "-oe", "-tnef", "-pua", filePath)
	if err != nil && err.Error() != "exit status 3" {
		return res, err
	}

	lines := strings.Split(res.Out, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, ">>> Virus ") {
			detection := strings.Split(line, " ")[2]
			res.Output = detection[1 : len(detection)-1]
			res.Infected = true
			break
		}
	}

	return res, nil
}

// GetVersion returns Sophos components' version
func GetVersion() (Version, error) {

	versionOut, err := utils.ExecCmd(savscan, "--version")

	// SAVScan virus detection utility
	// Copyright (c) 1989-2018 Sophos Limited. All rights reserved.

	// System time 19:28:51, System date 24 December 2018

	// Product version           : 5.53.0
	// Engine version            : 3.74.2
	// Virus data version        : 5.55
	// User interface version    : 2.03.074
	// Platform                  : Linux/AMD64
	// Released                  : 18 September 2018
	// Total viruses (with IDEs) : 25676226

	if err != nil {
		return Version{}, err
	}

	v := Version{}
	lines := strings.Split(versionOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Product version") {
			v.ProductVersion = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Engine version") {
			v.EngineVersion = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Virus data version") {
			v.VirusDataVersion = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "User interface version") {
			v.UserInterfaceVersion = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return v, nil
}
