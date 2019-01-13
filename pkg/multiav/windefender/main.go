// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package windefender

import (
	"os"
	"path"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

// Our consts
const (
	loadlibraryPath = "/opt/windowsdefender/"
	mpclient        = "./mpclient"
	mpenginedll     = "/engine/mpengine.dll"
)

// Result represents detection results.
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// GetVersion returns update version.
func GetVersion() (string, error) {
	mpenginedll := path.Join(loadlibraryPath, mpenginedll)
	versionOut, err := utils.ExecCommand("exiftool", "-ProductVersion", mpenginedll)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(versionOut, ":")[1]), nil
}

// ScanFile a file with Windows Defender scanner.
func ScanFile(filePath string) (Result, error) {

	res := Result{}

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
	mpclientOut, err := utils.ExecCommand(mpclient, filePath)
	if err != nil {
		return res, err
	}

	// main(): Scanning /samples/locky...
	// EngineScanCallback(): Scanning input
	// EngineScanCallback(): Threat Ransom:Win32/Locky.A identified.
	lines := strings.Split(mpclientOut, "\n")
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
