// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avast

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cmd          = "/bin/scan"
	avastService = "/etc/init.d/avast"
	licenseFile  = "/etc/avast/license.avastlic"
	vpsUpdate    = "/var/lib/avast/Setup/avast.vpsupdate"
)

// Result represents detection results.
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// GetVPSVersion returns VPS version.
func GetVPSVersion() (string, error) {

	// Run the scanner to grab the version
	out, err := utils.ExecCommand(cmd, "-V")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// GetProgramVersion returns program version.
func GetProgramVersion() (string, error) {

	// Run the scanner to grab the version
	out, err := utils.ExecCommand(cmd, "-v")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// ScanFilePath scans from a filepath.
func ScanFilePath(filepath string) (Result, error) {

	res := Result{}
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return res, err
	}

	// Execute the scanner with the given file path
	//  -a         Print all scanned files/URLs, not only infected.
	//  -b         Report decompression bombs as infections.
	//  -f         Scan full files.
	//  -u         Report potentionally unwanted programs (PUP).
	out, err := utils.ExecCommand(cmd, "-abfu", filepath)

	// Exit status:
	//  0 - no infections were found
	//  1 - some infected file was found
	//  2 - an error occurred
	if err != nil && err.Error() != "exit status 1" {
		return res, err
	}

	// Check if the file is infected
	if strings.Contains(out, "[OK]") {
		res.Infected = false
		return res, nil
	}

	// Sanitize the detection output
	det := strings.Split(out, "\t")
	res.Output = strings.TrimSpace(det[1])
	res.Infected = true
	return res, nil
}

// ScanURL scans a given URL
func ScanURL(url string) (string, error) {

	// Execute the scanner with the given URL
	out, err := utils.ExecCommand(cmd, "-U", url)

	// 	Exit status:
	// 0 - no infections were found
	// 1 - some infected file was found
	// 2 - an error occurred
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}

	// Check if we got a clean URL
	if out == "" {
		return "[OK]", nil
	}

	// Sanitize the output and return
	str := strings.Split(out, "\t")
	result := strings.TrimSpace(str[1])
	return result, nil
}

// UpdateVPS performs a VPS update.
func UpdateVPS() error {
	_, err := utils.ExecCommand(vpsUpdate)
	return err
}

// IsLicenseExpired returns true if license was expired
func IsLicenseExpired() (bool, error) {

	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		return true, errors.New("License not found")
	}
	  
	out, err := utils.ExecCommand(avastService, "status")
	if err != nil {
		return true, err
	}

	if strings.Contains(out, "License expired") {
		return true, nil
	}
	return false, nil
}

// ActivateLicense activate the license.
func ActivateLicense(r io.Reader) error {
	// Write the license file to disk
	_, err := utils.WriteBytesFile(licenseFile, r)
	if err != nil {
		return err
	}

	// Change the owner of the license file to `avast` user
	err = utils.ChownFileUsername(licenseFile, "avast")
	if err != nil {
		return err
	}

	// Restart the daemon to apply the license
	err = RestartService()
	if err != nil {
		return err
	}

	isExpired, err := IsLicenseExpired()
	if err != nil {
		return err
	}

	if isExpired {
		return errors.New("License is expired ")
	}

	return nil
}

// RestartService re-starts the Avast service.
func RestartService() error {
	// check if service is running
	_, err := utils.ExecCommand(avastService, "status")
	if err != nil && err.Error() != "exit status 3" {
		return err
	}

	// exit code 3 means program is not running
	if err.Error() == "exit status 3" {
		_, err = utils.ExecCommand(avastService, "start")
	} else {
		_, err = utils.ExecCommand(avastService, "restart")
	}

	return err
}

// ScanFileBinary receives a binary files, write it to disk then scan it
func ScanFileBinary(r io.Reader) (Result, error) {
	// Write the license file to disk
	_, err := utils.WriteBytesFile("sample", r)
	if err != nil {
		return Result{Output: ""}, err
	}

	return ScanFilePath("sample")
}
