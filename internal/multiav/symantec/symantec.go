// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package symantec

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	cmd      = "/opt/Symantec/symantec_antivirus/sav"
	logsDir  = "/var/symantec/sep/Logs/"
	symcfgd  = "/opt/Symantec/symantec_antivirus/symcfgd"
	rtvscand = "/opt/Symantec/symantec_antivirus/rtvscand"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// GetProgramVersion returns Symantec Program version
func GetProgramVersion() (string, error) {

	// Run the scanner to grab the version
	out, err := utils.ExecCmd(cmd, "info", "-p")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out, "\n"), nil
}

// ScanFile scans a given file
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// Symantec does not really provide a simple way to grab the results
	// We need to read it from a log file
	currentTime := time.Now()
	logfile := logsDir + currentTime.Format("01022006") + ".log"

	// Cleanup the previous log results.
	err = os.RemoveAll(logsDir)
	if err != nil {
		return res, err
	}

	// Create a directory to store the new logs.
	err = os.MkdirAll(logsDir, 0777)
	if err != nil {
		return res, err
	}

	err = utils.CreateFile(logfile)
	if err != nil {
		return res, err
	}

	// Execute the scanner with the given file path
	_, err = utils.ExecCmd("sudo", cmd, "manualscan", "--clscan", filePath)
	if err != nil {
		return res, err
	}

	data, err := utils.ReadAll(logfile)
	if err != nil {
		return res, err
	}

	// Check if infected.
	res.Out = string(data)
	if strings.Contains(res.Out, "Scan Complete:  Threats: 0") {
		return res, nil
	}

	// 3100010A1C11,3,2,1,ubuntu,root,,,,,,,16777216,"Scan started on selected drives and folders and all extensions.",1546367310,,0,,,,,0,,,,,,,,,,,,,,,,00:50:56:f9:3d:02,14.2.770.0000,,,,,,,,,,,,,,,,0,,,,,,,,,
	// 3100010A1C11,5,1,1,ubuntu,root,Trojan.Zbot!gen30,/sample/malware.exe,4,4,4,256,33570852,"",1546367310,,0,,0,42393,0,0,0,,,,20190101.002,197726,6,,0,,,,,,,00:50:56:f9:3d:02,14.2.770.0000,,,,,,,,,,,,,,,,0,,,0,,502		318464	2	48d04c7fd164aaf97037fe6c9abdbd290e9fd888152f6854a822138e37ee7413,,,,1
	// 3100010A1C12,2,2,1,ubuntu,root,,,,,,,16777216,"Scan Complete:  Threats: 1   Scanned: 1   Files/Folders/Drives Omitted: 0",1546367310,,0,1:1:1:0,,,,0,,,,,,,,,,,,,,,,00:50:56:f9:3d:02,14.2.770.0000,,,,,,,,,,,,,,,,
	lines := strings.Split(res.Out, "\n")
	if len(lines) < 2 {
		return res, multiav.ErrParseDetection
	}

	lines = strings.Split(lines[1], ",")
	if len(lines) < 7 {
		return res, multiav.ErrParseDetection
	}

	res.Infected = true
	res.Output = lines[6]
	return res, nil
}

// StartDaemon starts the symantec daemon.
func StartDaemon() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out, err := utils.ExecCmdWithContext(ctx, "sudo", symcfgd, "-x")
	if err != nil {
		return fmt.Errorf("failed to start symcfgd daemon, err: %v, out:%s",
			err, out)
	}
	out, err = utils.ExecCmdWithContext(ctx, "sudo", rtvscand, "-x")
	if err != nil {
		return fmt.Errorf("failed to start rtvscand daemon, err: %v, out:%s",
			err, out)
	}

	return nil
}
