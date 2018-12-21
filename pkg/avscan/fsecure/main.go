// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package fsecure

import (
	"regexp"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

// Our consts
const (
	fsav = "/opt/f-secure/fsav/bin/fsav"
)

// Detection represents detection results
type Detection struct {
	Infected bool   `json:"infected"`
	Version  string `json:"version"`
	FSE      string `json:"fse"`
	Aquarius string `json:"aquarius"`
}

// Version represents all fsecure components' versions
type Version struct {
	FSecureVersion          string `json:"fsecure_version"`
	DatabaseVersion         string `json:"database_version"`
	HydraEngineVersion      string `json:"hydra_engine_version"`
	HydraDatabaseVersion    string `json:"hydra_db_version"`
	AquariusEngineVersion   string `json:"aquarius_engine_version"`
	AquariusDatabaseVersion string `json:"aquarius_db_version"`
}

// ScanFile a file with FSecure scanner
func ScanFile(filePath string) (Detection, error) {

	// Run now
	fsavOut, err := utils.ExecCommand(fsav, "--virus-action1=report",
		"--suspected-action1=report", filePath)

	// FSAV has the following exit codes:
	// 0 - Normal exit; no viruses or suspicious files found.
	// 1 - Fatal error; unrecoverable error. (Usually a missing or corrupted file.)
	// 3 - A boot virus or file virus found.
	// 4 - Riskware (potential spyware) found.
	// 6 - At least one virus was removed and no infected files left.
	// 7 - Out of memory.
	// 8 - Suspicious files found; these are not necessarily infected by a virus.
	// 9 - Scan error, at least one file scan failed.
	// 130 - Program was terminated by pressing CTRL-C, or by a sigterm or suspend eve

	// The exit status is 0 if no infected files are found and 1
	// otherwise. If an error occurred, the exit status is 2.
	infectedExitCodes := []string{"exit status 3", "exit status 4",
		"exit status 6", "exit status 8"}
	if err != nil && !utils.StringInSlice(err.Error(), infectedExitCodes) {
		return Detection{}, err
	}

	// Parse fsav output
	// F-Secure Anti-Virus CLI version 1.0  build 0069

	// Scan started at Fri Dec 21 03:40:19 2018
	// Database version: 2018-12-21_05

	// Sality: Infected: Win32.Sality.M [Aquarius]
	// loadmoney: Riskware: Gen:Application.LoadMoney.1 (6, 1, 1) [Aquarius]
	// Desktop_update.exe.cld: Infected: Packed:MSIL/SmartIL.A [FSE]

	// Scan ended at Fri Dec 21 03:40:19 2018
	// 1 file scanned
	// 1 file infected

	d := Detection{}
	reg := regexp.MustCompile(` \(.*\)`)
	lines := strings.Split(fsavOut, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Infected: ") {
			parts := strings.Split(line, "Infected: ")
			detection := parts[len(parts)-1]
			if strings.Contains(detection, " [Aquarius]") {
				d.Aquarius = strings.TrimSuffix(detection, " [Aquarius]")
				d.Infected = true
				continue
			} else if strings.Contains(detection, " [FSE]") {
				d.FSE = strings.TrimSuffix(detection, " [FSE]")
				d.Infected = true
				continue
			}
		}

		if strings.Contains(line, "Riskware: ") {
			parts := strings.Split(line, "Riskware: ")
			detection := parts[len(parts)-1]
			if strings.Contains(detection, " [Aquarius]") {
				detection = strings.TrimSuffix(detection, " [Aquarius]")
				d.Aquarius = reg.ReplaceAllString(detection, "")
				d.Infected = true
				continue
			} else if strings.Contains(detection, " [FSE]") {
				detection = strings.TrimSuffix(detection, " [FSE]")
				d.FSE = reg.ReplaceAllString(detection, "")
				d.Infected = true
				continue
			}
		}

		if strings.Contains(line, "Database version: ") {
			d.Version = strings.TrimPrefix(line, "Database version: ")
		}
	}

	return d, nil
}

// GetVersion get Anti-Virus scanner version
func GetVersion() (Version, error) {

	fsavOut, err := utils.ExecCommand(fsav, "--version")

	// F-Secure Linux Security version 11.10 build 68

	// F-Secure Anti-Virus CLI Command line client version:
	// F-Secure Anti-Virus CLI version 1.0  build 0069

	// F-Secure Anti-Virus CLI Daemon version:
	// F-Secure Anti-Virus Daemon version 1.0  build 0161

	// Database version: 2018-12-21_08

	// Scanner Engine versions:
	// 	F-Secure Corporation Hydra engine version 5.22 build 28
	// 	F-Secure Corporation Hydra database version 2018-12-21_02

	// 	F-Secure Corporation Aquarius engine version 1.0 build 8
	// 	F-Secure Corporation Aquarius database version 2018-12-21_08

	if err != nil {
		return Version{}, err
	}

	r := Version{}
	lines := strings.Split(fsavOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Database version: ") {
			r.DatabaseVersion = strings.TrimPrefix(line, "Database version: ")
		} else if strings.Contains(line, "F-Secure Linux Security version ") {
			r.FSecureVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "F-Secure Linux Security version "))
		} else if strings.Contains(line, "Hydra engine version") {
			r.HydraEngineVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Hydra engine version "))
		} else if strings.Contains(line, "Hydra database version") {
			r.HydraDatabaseVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Hydra database version "))
		} else if strings.Contains(line, "Aquarius engine version") {
			r.AquariusEngineVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Aquarius engine version "))
		} else if strings.Contains(line, "Aquarius database version") {
			r.AquariusDatabaseVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Aquarius database version "))
		}
	}

	return r, nil
}
