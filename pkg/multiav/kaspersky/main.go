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
	kesl = "/opt/kaspersky/kesl/bin/kesl-control"
)

// Result represents detection results
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// Version represents database components' versions.
type Version struct {
	CurrentAVDatabasesDate    string `json:"current_av_db_ate"`
	LastAVDatabasesUpdateDate string `json:"last_av_db_update_date"`
	CurrentAVDatabasesState   string `json:"current_av_db_state"`
	CurrentAVDatabasesRecords string `json:"current_av_db_records"`
}

// GetProgramVersion returns Kaspersky Anti-Virus for Linux File Server version.
func GetProgramVersion() (string, error) {

	// Run kesl to grab the version
	out, err := utils.ExecCommand(kesl, "-S", "--app-info")
	if err != nil {
		return "", err
	}

	version := ""
	lines := strings.Split(out, "\n")
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
	databaseOut, err := utils.ExecCommand(kesl, "--get-stat", "Update")

	ver := Version{}
	if err != nil {
		return ver, nil
	}

	lines := strings.Split(databaseOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Current AV databases date") {
			ver.CurrentAVDatabasesDate = strings.TrimSpace(strings.TrimPrefix(line, "Current AV databases date:"))
		} else if strings.Contains(line, "Last AV databases update date") {
			ver.LastAVDatabasesUpdateDate = strings.TrimSpace(strings.TrimPrefix(line, "Last AV databases update date:"))
		} else if strings.Contains(line, "Current AV databases state") {
			ver.CurrentAVDatabasesState = strings.TrimSpace(strings.TrimPrefix(line, "Current AV databases state:"))
		} else if strings.Contains(line, "Current AV databases records") {
			ver.CurrentAVDatabasesRecords = strings.TrimSpace(strings.TrimPrefix(line, "Current AV databases records:"))
		}
	}
	return ver, nil
}

// ScanFile a file with Kaspersky scanner
func ScanFile(filePath string) (Result, error) {

	// Clean the states
	res := Result{}

	// Run now
	out, err := utils.ExecCommand(kesl, "--scan-file", filePath, "--action", "Skip")
	// root@404e0cc38216:/# /opt/kaspersky/kesl/bin/kesl-control --scan-file eicar.com.txt --action SKip
	// Scanned objects                     : 1
	// Total detected objects              : 1
	// Infected objects and other objects  : 1
	// Disinfected objects                 : 0
	// Moved to Storage                    : 0
	// Removed objects                     : 0
	// Not disinfected objects             : 1
	// Scan errors                         : 0
	// Password-protected objects          : 0
	// Skipped                             : 0

	if err != nil {
		return res, err
	}

	// Check if infected
	if !strings.Contains(out, "Total detected objects              : 1") {
		return res, nil
	}

	// Grab detection name with a separate cmd
	kavOut, err := utils.ExecCommand(kesl, "-E", "--query", "\"EventType=='ThreatDetected'\"")
	// EventType=ThreatDetected
	// EventId=2544
	// Date=2019-06-11 22:12:16
	// DangerLevel=Critical
	// FileName=/eicar
	// ObjectName=File
	// TaskName=Scan_File_ca3f0bc2-ce71-4d4a-bdc1-c8ae502566d0
	// RuntimeTaskId=4
	// TaskId=100
	// DetectName=EICAR-Test-File
	// TaskType=ODS
	// FileOwner=root
	// FileOwnerId=0
	// DetectCertainty=Sure
	// DetectType=Virware
	// DetectSource=Local
	// ObjectId=1
	// AccessUser=root
	// AccessUserId=0
	if err != nil {
		return res, err
	}

	lines := strings.Split(kavOut, "\n")
	if len(lines) > 0 {
		res.Output = strings.TrimSpace(strings.Split(lines[9], "=")[1])
		res.Infected = true
	}

	return res, nil
}
