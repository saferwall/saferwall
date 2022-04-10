// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package fsecure

import (
	"regexp"
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	fsav = "/opt/f-secure/fsav/bin/fsav"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

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
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	res.Out, err = utils.ExecCmd(fsav, "--virus-action1=report",
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
		return res, err
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

	reg := regexp.MustCompile(` \(.*\)`)
	lines := strings.Split(res.Out, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Infected: ") {
			parts := strings.Split(line, "Infected: ")
			detection := parts[len(parts)-1]
			if strings.Contains(detection, " [Aquarius]") {
				res.Output = strings.TrimSuffix(detection, " [Aquarius]")
				res.Infected = true
				continue
			} else if strings.Contains(detection, " [FSE]") {
				res.Output = strings.TrimSuffix(detection, " [FSE]")
				res.Infected = true
				continue
			}
		}
		if strings.Contains(line, "Riskware: ") {
			parts := strings.Split(line, "Riskware: ")
			detection := parts[len(parts)-1]
			if strings.Contains(detection, " [Aquarius]") {
				detection = strings.TrimSuffix(detection, " [Aquarius]")
				res.Output = reg.ReplaceAllString(detection, "")
				res.Infected = true
				continue
			} else if strings.Contains(detection, " [FSE]") {
				detection = strings.TrimSuffix(detection, " [FSE]")
				res.Output = reg.ReplaceAllString(detection, "")
				res.Infected = true
				continue
			}
		}
	}

	return res, nil
}

// GetVersion get Anti-Virus scanner version.
func GetVersion() (Version, error) {

	fsavOut, err := utils.ExecCmd(fsav, "--version")

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

	ver := Version{}
	if err != nil {
		return ver, err
	}

	lines := strings.Split(fsavOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Database version: ") {
			ver.DatabaseVersion = strings.TrimPrefix(line, "Database version: ")
		} else if strings.Contains(line, "F-Secure Linux Security version ") {
			ver.FSecureVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "F-Secure Linux Security version "))
		} else if strings.Contains(line, "Hydra engine version") {
			ver.HydraEngineVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Hydra engine version "))
		} else if strings.Contains(line, "Hydra database version") {
			ver.HydraDatabaseVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Hydra database version "))
		} else if strings.Contains(line, "Aquarius engine version") {
			ver.AquariusEngineVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Aquarius engine version "))
		} else if strings.Contains(line, "Aquarius database version") {
			ver.AquariusDatabaseVersion = strings.TrimSpace(
				strings.TrimPrefix(line, "\tF-Secure Corporation Aquarius database version "))
		}
	}

	return ver, nil
}
