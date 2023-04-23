// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package windefender

import (
	"context"
	"os"
	"path"
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	loadlibraryPath = "/opt/windows-defender/"
	mploader        = "./mploader.exe"
	mpenginedll     = "/engine/mpengine.dll"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// GetVersion returns update version.
func GetVersion() (string, error) {
	mpenginedll := path.Join(loadlibraryPath, mpenginedll)
	out, err := utils.ExecCmd("exiftool", "-ProductVersion",
		mpenginedll)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(out, ":")[1]), nil
}

// ScanFile a file with Windows Defender scanner.
func (Scanner) ScanFile(filepath string, opts multiav.Options) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// get current working directory
	dir, err := utils.Getwd()
	if err != nil {
		return res, err
	}

	// mploader requires us to run from loadlibrary folder or it fails
	if err := os.Chdir(loadlibraryPath); err != nil {
		return res, err
	}
	defer os.Chdir(dir)

	if opts.ScanTimeout == 0 {
		opts.ScanTimeout = multiav.DefaultScanTimeout
	}

	// Create a new context and add a timeout to it.
	ctx, cancel := context.WithTimeout(
		context.Background(), opts.ScanTimeout)
	defer cancel()

	// Execute the scanner with the given file path
	// main(): usage: ./mploader -f <filePath> -u
	res.Out, err = utils.ExecCmdWithContext(ctx, "wine", mploader, "-f", filepath, "-u")
	if err != nil {
		return res, err
	}

	// main(): Scanning /samples/locky...
	// ....
	// Engine Boot Success!
	// 0024:fixme:ntdll:EtwEventActivityIdControl 0x2, 0032F9DC: stub
	// Scan Start /eicar
	// 0024:fixme:wintrust:CryptCATAdminAcquireContext2 0032C530 (null) L"SHA1" 00000000 0 stub
	// 0024:fixme:wintrust:CryptCATAdminAcquireContext2 0032C530 (null) L"SHA256" 00000000 0 stub
	// 0024:fixme:bcrypt:BCryptGenRandom ignoring selected algorithm
	// Threat Virus:DOS/EICAR_Test_File identified.
	lines := strings.Split(res.Out, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "Threat ") {
			continue
		}
		if strings.Contains(line, "No Threat identified ") {
			continue
		}

		detection := strings.TrimPrefix(line, "Threat ")
		res.Output = strings.TrimSuffix(detection, " identified.\r")
		res.Infected = true
		break
	}

	return res, nil
}
