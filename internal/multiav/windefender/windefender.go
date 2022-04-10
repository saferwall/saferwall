// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package windefender

import (
	"os"
	"path"
	"strings"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

// Our consts
const (
	loadlibraryPath = "/opt/windows-defender/"
	mpclient        = "./mpclient"
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
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// get current working directory
	dir, err := utils.Getwd()
	if err != nil {
		return res, err
	}

	// mpclient requires us to run from loadlibrary folder or it fails
	if err := os.Chdir(loadlibraryPath); err != nil {
		return res, err
	}
	defer os.Chdir(dir)

	// Execute the scanner with the given file path
	// main(): usage: ./mpclient [filenames...]
	res.Out, err = utils.ExecCmd(mpclient, filePath)
	if err != nil {
		return res, err
	}

	// main(): Scanning /samples/locky...
	// EngineScanCallback(): Scanning input
	// EngineScanCallback(): Threat Ransom:Win32/Locky.A identified.
	lines := strings.Split(res.Out, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "EngineScanCallback(): Threat ") {
			continue
		}

		detection := strings.TrimPrefix(line, "EngineScanCallback(): Threat ")
		res.Output = strings.TrimSuffix(detection, " identified.")
		res.Infected = true
		break
	}

	return res, nil
}
