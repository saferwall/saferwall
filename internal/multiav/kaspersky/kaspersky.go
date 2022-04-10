// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package kaspersky

import (
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	kesl         = "/opt/kaspersky/kesl/bin/kesl-control"
	keslLauncher = "/opt/kaspersky/kesl/libexec/kesl_launcher.sh"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// Version represents database components' versions.
type Version struct {
	CurrentAVDatabasesDate    string `json:"current_av_db_ate"`
	LastAVDatabasesUpdateDate string `json:"last_av_db_update_date"`
	CurrentAVDatabasesState   string `json:"current_av_db_state"`
	CurrentAVDatabasesRecords string `json:"current_av_db_records"`
}

// GetProgramVersion returns Kaspersky Anti-Virus for Linux File Server version.
func GetProgramVersion() (string, error) {

	out, err := utils.ExecCmd("sudo", kesl, "-S", "--app-info")
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

	databaseOut, err := utils.ExecCmd("sudo", kesl, "--get-stat", "Update")

	ver := Version{}
	if err != nil {
		return ver, nil
	}

	lines := strings.Split(databaseOut, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Current AV databases date") {
			ver.CurrentAVDatabasesDate = strings.TrimSpace(
				strings.TrimPrefix(line, "Current AV databases date:"))
		} else if strings.Contains(line, "Last AV databases update date") {
			ver.LastAVDatabasesUpdateDate = strings.TrimSpace(
				strings.TrimPrefix(line, "Last AV databases update date:"))
		} else if strings.Contains(line, "Current AV databases state") {
			ver.CurrentAVDatabasesState = strings.TrimSpace(
				strings.TrimPrefix(line, "Current AV databases state:"))
		} else if strings.Contains(line, "Current AV databases records") {
			ver.CurrentAVDatabasesRecords = strings.TrimSpace(
				strings.TrimPrefix(line, "Current AV databases records:"))
		}
	}
	return ver, nil
}

// ScanFile a file with Kaspersky scanner
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// Return codes
	// 0 – command / task completed successfully.
	// 1 – general error in command arguments.
	// 2 – error in passed application settings.
	// 64 – Kaspersky Endpoint Security is not running.
	// 65 - Protection is disabled.
	// 66 – anti-virus databases have not been downloaded (used only for the command kesl-control --app-info).
	// 67 – activation 2.0 ended with an error due to network problems.
	// 68 – the command cannot be executed because the application is running under a policy.
	// 128 – unknown error.
	// 65 – all other errors.

	// Run now
	res.Out, err = utils.ExecCmd("sudo", kesl, "--scan-file",
		filePath, "--action", "Skip")
	// root@404e0cc38216:/# /opt/kaspersky/kesl/bin/kesl-control --scan-file eicar.com.txt --action Skip
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
	if !strings.Contains(res.Out, "Total detected objects              : 1") {
		return res, nil
	}

	// Grab detection name with a separate cmd
	// sudo /opt/kaspersky/kesl/bin/kesl-control -E --query "EventType=='ThreatDetected'"
	out, err := utils.ExecCmd("sudo", kesl, "-E", "--query",
		"EventType=='ThreatDetected'")
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
	res.Out += ", query res: " + out
	if err != nil {
		return res, err
	}

	// so hackish, there is no easy way to grab detection name
	// no way to clean all these events as it was in previous version
	// so pretty hardcoded for now
	lines := strings.Split(res.Out, "\n\n")
	if len(lines) > 0 {
		index := len(lines) - 1
		lines = strings.Split(lines[index], "\n")
		if len(lines) > 8 {
			res.Output = strings.TrimSpace(strings.Split(lines[9], "=")[1])
			res.Infected = true
		}
	}

	return res, nil
}

// GetLicenseInfos queries license infos
func GetLicenseInfos() (string, error) {

	out, err := utils.ExecCmd("sudo", kesl, "-L", "--query")
	// Active key information:
	// Expiration date                      : 2019-07-13
	// Days remaining until expiration      : 0
	// Protection                           : No protection
	// Updates                              : No updates
	// Key status                           : Expired
	// License type                         : XYZ
	// Usage restriction                    : 1
	// Application name                     : Kaspersky Endpoint Security 10 SP1 MR1 for Linux
	// Active key                           : XYZ
	// Activation date                      : 2019-06-12
	if err != nil {
		return "", err
	}

	return out, err
}

// StartDaemon starts the kaspersky daemon.
func StartDaemon() error {
	err := utils.ExecCmdBackground("sudo", keslLauncher, "-n", "-D")
	return err
}
