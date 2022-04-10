// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avira

import (
	"os"
	"regexp"
	"testing"

	multiav "github.com/saferwall/saferwall/internal/multiav"
)

type filePathTest struct {
	filepath string
	want     multiav.Result
}

var filepathScanTest = []filePathTest{
	{"../../test/testdata/765c3a580f885f5e4e4f98a709e9f0ce",
		multiav.Result{Infected: true, Output: "TR/Crypt.XPACK.Gen2"}},
}

func TestVersion(t *testing.T) {
	ver, err := GetVersion()
	if err != nil {
		t.Fatalf("TestVersion failed, got: %s", err)
	}

	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	l := re.FindStringSubmatch(ver.ScanCLVersion)
	if len(l) == 0 {
		t.Errorf(`ScanCL version was incorrect, got: %s,\
		 want something similar to: 1.9.161.2`, ver)
	}

	l = re.FindStringSubmatch(ver.CoreVersion)
	if len(l) == 0 {
		t.Errorf(`Core version was incorrect, got: %s,\
		 want something similar to: 1.9.2.0`, ver)
	}

	l = re.FindStringSubmatch(ver.VDFVersion)
	if len(l) == 0 {
		t.Errorf(`VDF version was incorrect, got: %s,\
		 want something similar to: 7.15.16.96`, ver)
	}
}

func TestScanFile(t *testing.T) {
	s := Scanner{}
	for _, tt := range filepathScanTest {
		t.Run(tt.filepath, func(t *testing.T) {
			got, err := s.ScanFile(tt.filepath)
			if err != nil {
				t.Fatalf("TestScanFile(%s) failed, err: %s, got: %v",
					tt.filepath, err, got)
			}
			if got.Output != tt.want.Output || got.Infected != tt.want.Infected {
				t.Errorf("TestScanFile(%s) got %v, want %v",
					tt.filepath, got, tt.want)
			}
		})
	}
}

func TestLicenceStatus(t *testing.T) {
	_, err := LicenseStatus()
	if err != nil {
		t.Errorf("TestLicenceStatus failed, err: %s", err)
	}
}

func TestNoLicenseFound(t *testing.T) {

	// Deliberately removing the license file
	err := os.Remove(LicenseKeyPath)
	if err != nil {
		t.Errorf("TestNoLicenseFound failed, err: %s", err)
	}

	_, err = LicenseStatus()
	if err != ErrNoLicenseFound {
		t.Errorf("TestNoLicenseFound failed, err: %s", err)
	}
}
