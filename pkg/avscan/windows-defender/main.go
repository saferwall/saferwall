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
	loadlibraryPath = "/usr/local/loadlibrary"
	mpclient        = "./mpclient"
	mpenginedll     = "/engine/mpengine.dll"
)

// GetVersion returns update version
func GetVersion() (string, error) {
	mpenginedll := path.Join(loadlibraryPath, mpenginedll)
	versionOut, err := utils.ExecCommand("exiftool", "-ProductVersion", mpenginedll)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Split(versionOut, ":")[1]), nil
}

// ScanFile a file with Windows Defender scanner
func ScanFile(filePath string) (string, error) {

	// get current working directory
	dir, err := utils.Getwd()
	if err != nil {
		return "", err
	}

	// mpclient requires us to run from loadlibrary folder or it fails
	if err := os.Chdir(loadlibraryPath); err != nil {
		return "", err
	}
	defer os.Chdir(dir)

	// Execute the scanner with the given file path
	// main(): usage: ./mpclient [filenames...]
	mpclientOut, err := utils.ExecCommand(mpclient, filePath)
	if err != nil {
		return "", err
	}

	// main(): Scanning /samples/locky...
	// EngineScanCallback(): Scanning input
	// EngineScanCallback(): Threat Ransom:Win32/Locky.A identified.
	detection := ""
	lines := strings.Split(mpclientOut, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "EngineScanCallback(): Threat ") {
			continue
		}

		detection = strings.TrimPrefix(line, "EngineScanCallback(): Threat ")
		detection = strings.TrimSuffix(detection, " identified.")
		break
	}

	return detection, nil
}
