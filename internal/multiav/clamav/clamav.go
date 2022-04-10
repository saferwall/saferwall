// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package clamav

import (
	"context"
	"fmt"
	"strings"
	"time"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	clamdscan = "/usr/bin/clamdscan"
	clamd     = "clamd"
	timeout   = 60
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// ScanFile performs antivirus scan.
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// Execute the scanner with the given file path
	// --no-summary   Disable summary at end of scanning
	res.Out, err = utils.ExecCmd(clamdscan, "--no-summary", filePath)

	// clamscan return values (documented from man clamscan)
	//   0 : No virus found.
	//   1 : Virus(es) found.
	//   2 : Some error(s) occurred.

	//  40: Unknown option passed.
	//  50: Database initialization error.
	//  52: Not supported file type.
	//  53: Can't open directory.
	//  54: Can't open file. (ofm)
	//  55: Error reading file. (ofm)
	//  56: Can't stat input file / directory.
	//  57: Can't get absolute path name of current working directory.
	//  58: I/O error, please check your file system.
	//  62: Can't initialize logger.
	//  63: Can't create temporary files/directories (check permissions).
	//  64: Can't write to temporary directory (please specify another one).
	//  70: Can't allocate memory (calloc).
	//  71: Can't allocate memory (malloc).
	if err != nil && err.Error() != "exit status 1" {
		return res, err
	}

	// samples/locky: Win.Malware.Locky-5540 FOUND
	// samples/putty: OK
	if strings.HasSuffix(res.Out, "OK\n") {
		return res, nil
	}

	if !strings.HasSuffix(res.Out, "FOUND\n") {
		return res, nil
	}

	// Extract detection name if infected
	parts := strings.Split(res.Out, ": ")
	det := parts[len(parts)-1]
	res.Output = strings.TrimSuffix(det, " FOUND\n")
	res.Infected = true
	return res, nil
}

// Version returns program version.
func Version() (string, error) {

	// Execute the scanner with the given file path
	out, err := utils.ExecCmd(clamdscan, "--version")
	if err != nil {
		return "", err
	}

	// Extract the version
	// ClamAV 0.100.2/25284/Wed Jan  9 18:42:45 2019
	ver := strings.Split(out, "/")[0]
	ver = strings.Split(ver, " ")[1]
	return ver, nil
}

// StartDaemon starts the clamd daemon.
func StartDaemon() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	out, err := utils.ExecCmdWithContext(ctx, clamd)
	if err != nil {
		return fmt.Errorf("failed to start daemon, err: %v, out:%s", err, out)
	}
	return nil
}
