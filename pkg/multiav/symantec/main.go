// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package symantec

import (
	"os"
	"strings"
	"time"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cmd     = "/opt/Symantec/symantec_antivirus/sav"
	logsDir = "/var/symantec/sep/Logs/"
)

// Result represents detection results.
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// GetProgramVersion returns Symantec Program version
func GetProgramVersion() (string, error) {

	// Run the scanner to grab the version
	out, err := utils.ExecCommand(cmd, "info", "-p")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out, "\n"), nil
}

// ScanFile scans a given file
func ScanFile(filepath string) (Result, error) {

	res := Result{}

	// Symantec does not really provide a simple way to grab the results
	// We need to read it from a log file
	currentTime := time.Now()
	logfile := logsDir + currentTime.Format("01022006") + ".log"

	// Cleanup
	os.RemoveAll(logsDir)
	os.MkdirAll(logsDir, 0777)
	err := utils.CreateFile(logfile)
	if err != nil {
		return res, err
	}

	// Execute the scanner with the given file path
	_, err = utils.ExecCommand("sudo", cmd, "manualscan", "--clscan", filepath)
	if err != nil {
		return res, err
	}

	data, err := utils.ReadAll(logfile)
	if err != nil {
		return res, err
	}

	// Check if infected.
	savOut := string(data)
	if strings.Contains(savOut, "Scan Complete:  Threats: 0") {
		return res, nil
	}

	// 3100010A1C11,3,2,1,ubuntu,root,,,,,,,16777216,"Scan started on selected drives and folders and all extensions.",1546367310,,0,,,,,0,,,,,,,,,,,,,,,,00:50:56:f9:3d:02,14.2.770.0000,,,,,,,,,,,,,,,,0,,,,,,,,,
	// 3100010A1C11,5,1,1,ubuntu,root,Trojan.Zbot!gen30,/sample/malware.exe,4,4,4,256,33570852,"",1546367310,,0,,0,42393,0,0,0,,,,20190101.002,197726,6,,0,,,,,,,00:50:56:f9:3d:02,14.2.770.0000,,,,,,,,,,,,,,,,0,,,0,,502		318464	2	48d04c7fd164aaf97037fe6c9abdbd290e9fd888152f6854a822138e37ee7413,,,,1
	// 3100010A1C12,2,2,1,ubuntu,root,,,,,,,16777216,"Scan Complete:  Threats: 1   Scanned: 1   Files/Folders/Drives Omitted: 0",1546367310,,0,1:1:1:0,,,,0,,,,,,,,,,,,,,,,00:50:56:f9:3d:02,14.2.770.0000,,,,,,,,,,,,,,,,
	lines := strings.Split(savOut, "\n")
	lines = strings.Split(lines[1], ",")
	res.Output = lines[6]
	res.Infected = true

	return res, nil
}
