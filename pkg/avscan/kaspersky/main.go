// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package kaspersky

import (
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

// Our consts
const (
	kav4fs = "/opt/kaspersky/kav4fs/bin/kav4fs-control"
)

// Detection represents detection results
type Detection struct {
	Infected bool   `json:"infected"`
	Result   string `json:"version"`
}

// Version represents database components' versions
type Version struct {
	CurrentAVDatabasesDate    string `json:"current_av_db_ate"`
	LastAVDatabasesUpdateDate string `json:"last_av_db_update_date"`
	CurrentAVDatabasesState   string `json:"current_av_db_state"`
	CurrentAVDatabasesRecords string `json:"current_av_db_records"`
}

// GetProgramVersion returns Kaspersky Anti-Virus for Linux File Server version
func GetProgramVersion() (string, error) {

	// Run kav4s to grab the version
	versionOut, err := utils.ExecCommand(kav4fs, "-S", "--app-info")
	if err != nil {
		return "", err
	}

	version := ""
	lines := strings.Split(versionOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Version:") {
			version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
			break
		}
	}

	return version, nil
}

// GetDatabaseVersion returns AV database update version
func GetDatabaseVersion() (Version, error) {

	// Run kav4s to grab the database update version
	databaseOut, err := utils.ExecCommand(kav4fs, "--get-stat", "Update")
	if err != nil {
		return Version{}, nil
	}

	v := Version{}
	lines := strings.Split(databaseOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Current AV databases date") {
			v.CurrentAVDatabasesDate = strings.TrimSpace(strings.TrimPrefix(line, "Current AV databases date:"))
		} else if strings.Contains(line, "Last AV databases update date") {
			v.LastAVDatabasesUpdateDate = strings.TrimSpace(strings.TrimPrefix(line, "Last AV databases update date:"))
		} else if strings.Contains(line, "Current AV databases state") {
			v.CurrentAVDatabasesState = strings.TrimSpace(strings.TrimPrefix(line, "Current AV databases state:"))
		} else if strings.Contains(line, "Current AV databases records") {
			v.CurrentAVDatabasesRecords = strings.TrimSpace(strings.TrimPrefix(line, "Current AV databases records:"))
		}
	}
	return v, nil
}

// ScanFile a file with Kaspersky scanner
func ScanFile(filePath string) (Detection, error) {

	// Run now
	kav4fsOut, err := utils.ExecCommand(kav4fs, "--scan-file", filePath)
	// /opt/kaspersky/kav4fs/bin/kav4fs-control --scan-file locky
	// Objects scanned:     1
	// Threats found:       1
	// Riskware found:      0
	// Infected:            1
	// Suspicious:          0
	// Cured:               0
	// Moved to quarantine: 0
	// Removed:             0
	// Not cured:           0
	// Scan errors:         0
	// Password protected:  0
	if err != nil {
		return Detection{}, err
	}

	// Check if file is infected
	infected := false
	if strings.Contains(kav4fsOut, "Threats found:       1") {
		infected = true
	}

	// If not infected, return immediately
	if !infected {
		return Detection{}, nil
	}

	// Grab detection name with a separate cmd
	kavOut, err := utils.ExecCommand(kav4fs, "--top-viruses", "1")
	// Viruses found: 1
	// Virus name:       Trojan-Ransom.Win32.Locky.d
	// Infected objects: 1
	if err != nil {
		return Detection{}, err
	}

	d := Detection{Infected: true}
	lines := strings.Split(kavOut, "\n")
	if len(lines) > 0 {
		d.Result = strings.TrimSpace(strings.Split(lines[1], ":")[1])
	}

	// Clean the states
	utils.ExecCommand(kav4fs, "--clean-stat")
	return d, nil
}
