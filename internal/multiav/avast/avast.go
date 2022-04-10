// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avast

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	cmd         = "scan"
	avastDaemon = "/usr/bin/avast"
	licenseFile = "/etc/avast/license.avastlic"
	vpsUpdate   = "/var/lib/avast/Setup/avast.vpsupdate"
	scanTimeout = 10 * time.Second
	tmpFilename = "tmpFile"
)

var (
	errLicenseExpired = errors.New("license was expired ")
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// VPSVersion returns VPS version.
func VPSVersion() (string, error) {

	// Run the scanner to grab the version
	out, err := utils.ExecCmd(cmd, "-V")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// ProgramVersion returns the program version.
func ProgramVersion() (string, error) {

	// Run the scanner to grab the version
	out, err := utils.ExecCmd(cmd, "-v")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// ScanFile scans from a filepath.
func (Scanner) ScanFile(filepath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return res, err
	}

	// Create a new context and add a timeout to it.
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(scanTimeout))
	defer cancel()

	// Execute the scanner with the given file path
	//  -a         Print all scanned files/URLs, not only infected.
	//  -b         Report decompression bombs as infections.
	//  -f         Scan full files.
	//  -u         Report potentionally unwanted programs (PUP).
	res.Out, err = utils.ExecCmdWithContext(ctx, cmd, "-abfu", filepath)

	// Exit status:
	//  0 - no infections were found
	//  1 - some infected file was found
	//  2 - an error occurred
	if err != nil && err.Error() != "exit status 1" {
		return res, err
	}

	// Check if the file is clean.
	if strings.Contains(res.Out, "[OK]") {
		return res, nil
	}

	// Sanitize the detection output
	det := strings.Split(res.Out, "\t")
	if len(det) < 2 {
		return res, multiav.ErrParseDetection
	}

	res.Infected = true
	res.Output = strings.TrimSpace(det[1])
	return res, nil
}

// ScanReader scans from a reader.
func (avast Scanner) ScanReader(r io.Reader) (multiav.Result, error) {
	_, err := utils.WriteBytesFile(tmpFilename, r)
	if err != nil {
		return multiav.Result{Output: ""}, err
	}

	return avast.ScanFile(tmpFilename)
}

// ScanURL scans a given URL.
func (Scanner) ScanURL(url string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// Execute the scanner with the given URL
	res.Out, err = utils.ExecCmd(cmd, "-U", url)

	// Exit status:
	//  0 - no infections were found
	//  1 - some infected file was found
	//  2 - an error occurred
	if err != nil && err.Error() != "exit status 1" {
		return res, err
	}

	// Check if we got a clean URL.
	if res.Out == "" {
		return multiav.Result{Infected: false}, nil
	}

	// Sanitize the output.
	str := strings.Split(res.Out, "\t")
	if len(str) < 2 {
		return res, multiav.ErrParseDetection
	}
	return multiav.Result{
		Output: strings.TrimSpace(str[1]), Infected: true}, nil
}

// UpdateVPS performs a VPS update.
func UpdateVPS() error {
	_, err := utils.ExecCmd(vpsUpdate)
	return err
}

// IsLicenseExpired returns true if license was expired.
func IsLicenseExpired() (bool, error) {

	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		return true, errors.New("license not found")
	}

	out, err := utils.ExecCmd(avastDaemon, "status")
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

	// Change the owner of the license file to `avast` user.
	err = utils.ChownFileUsername(licenseFile, "avast")
	if err != nil {
		return err
	}

	// Restart the service to apply the license.
	err = RestartService()
	if err != nil {
		return err
	}

	isExpired, err := IsLicenseExpired()
	if err != nil {
		return err
	}

	if isExpired {
		return errLicenseExpired
	}

	return nil
}

// RestartService re-starts the Avast service.
func RestartService() error {
	// check if service is running
	_, err := utils.ExecCmd(avastDaemon, "status")
	if err != nil && err.Error() != "exit status 3" {
		return err
	}

	// exit code 3 means program is not running
	if err.Error() == "exit status 3" {
		_, err = utils.ExecCmd(avastDaemon, "start")
	} else {
		_, err = utils.ExecCmd(avastDaemon, "restart")
	}

	return err
}

// StartDaemon starts the Avast daemon.
func StartDaemon() error {
	err := utils.ExecCmdBackground("sudo", avastDaemon, "-n", "-D")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = utils.ExecCmdBackground("sudo", avastDaemon, "-n", "-D")
	return err
}
