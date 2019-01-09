// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avast

import (
	"regexp"
	"testing"
)


type filePathTest struct {
	filepath string
	want Result
}

var filepathScanTest = []filePathTest{
	{"../../../test/multiav/eicar.com", Result{Infected: true, Output: "EICAR Test-NOT virus!!!"}},
}

func TestGetVPSVersion(t *testing.T) {
	version, err := GetVPSVersion()
	if err != nil {
		t.Fatalf("TestGetVPSVersion failed, got: %s", err)

	}

	re := regexp.MustCompile(`\d{8}`)
	l := re.FindStringSubmatch(version)
	if len(l) == 0 {
		t.Fatalf("VPS version was incorrect, got: %s, want something similar to: 19010602", version)
	}
}

func TestGetProgramVersion(t *testing.T) {
	version, err := GetProgramVersion()
	if err != nil {
		t.Fatalf("TestGetProgramVersion failed, got: %s", err)
	}

	re := regexp.MustCompile(`\d{1}\.\d{1}\.\d{1}`)
	l := re.FindStringSubmatch(version)
	if len(l) == 0 {
		t.Errorf("Program version was incorrect, got: %s, want something similar to: 2.2.0", version)
	}
}

func TestScanFilePath(t *testing.T) {
	for _, tt := range filepathScanTest {
		t.Run(tt.filepath, func(t *testing.T) {
			got, err := ScanFilePath(tt.filepath)
			if err != nil {
				t.Fatalf("TestScanFilePath(%s) failed, err: %s", tt.filepath, err)
			}
			if got != tt.want {
				t.Fatalf("TestScanFilePath(%s) got %v, want %v", tt.filepath, got, tt.want)
			}
		})
	}
}
