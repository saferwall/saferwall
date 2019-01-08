// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avira

import (
	"os"
	"regexp"
	"testing"

	"github.com/saferwall/saferwall/pkg/multiav/avira"
)

type filePathTest struct {
	filepath string
	want     Result
}

var filepathScanTest = []filePathTest{
	{"../../../test/multiav/eicar.com", Result{Infected: true, Output: "Eicar-Test-Signature"}},
}

func TestGetVersion(t *testing.T) {
	ver, err := GetVersion()
	if err != nil {
		t.Errorf("TestGetVersion failed, got: %s", err)

	}

	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	l := re.FindStringSubmatch(ver.ScanCLVersion)
	if len(l) == 0 {
		t.Errorf("ScanCL version was incorrect, got: %s, want something similar to: 1.9.161.2", ver)
	}

	l = re.FindStringSubmatch(ver.CoreVersion)
	if len(l) == 0 {
		t.Errorf("Core version was incorrect, got: %s, want something similar to: 1.9.2.0", ver)
	}

	l = re.FindStringSubmatch(ver.VDFVersion)
	if len(l) == 0 {
		t.Errorf("VDF version was incorrect, got: %s, want something similar to: 7.15.16.96", ver)
	}
}
func TestScanFile(t *testing.T) {
	for _, tt := range filepathScanTest {
		t.Run(tt.filepath, func(t *testing.T) {
			got, err := ScanFile(tt.filepath)
			if err != nil {
				t.Errorf("TestScanFile(%s) failed, err: %s", tt.filepath, err)
			}
			if got != tt.want {
				t.Errorf("TestScanFile(%s) got %v, want %v", tt.filepath, got, tt.want)
			}
		})
	}
}

func TestIsLicenseExpired(t *testing.T) {
	isExpired, err := IsLicenseExpired()
	if err != nil {
		t.Errorf("TestIsLicenseExpired failed, err: %s", err)
	}
	if isExpired {
		t.Errorf("TestIsLicenseExpired failed, license expired")
	}
}

func TestIsLicenseExpiredNoLicenseFound(t *testing.T) {
	// Deliberately removing the license file
	err := os.Remove(avira.LicenseKeyPath)
	if err != nil {
		t.Errorf("TestIsLicenseExpiredNoLicenseFound failed, err: %s", err)
	}

	
	_, err := IsLicenseExpired()
	if err != avira.ErrNoLicenseFound {
		t.Errorf("TestIsLicenseExpiredNoLicenseFound failed, err: %s", err)

	}
}
