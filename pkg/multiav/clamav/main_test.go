// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package clamav

import (
	"regexp"
	"testing"
)


type filePathTest struct {
	filepath string
	want Result
}

var filepathScanTest = []filePathTest{
	{"../../../test/multiav/eicar.com", Result{Infected: true, Output: "Eicar-Test-Signature"}},
}

func TestGetProgramVersion(t *testing.T) {
	version, err := GetVersion()
	if err != nil {
		t.Fatalf("TestGetProgramVersion failed, got: %s", err)
	}

	re := regexp.MustCompile(`\d{1}\.\d{3}\.\d{1}`)
	l := re.FindStringSubmatch(version)
	if len(l) == 0 {
		t.Fatalf("Program version was incorrect, got: %s, want something similar to: 0.100.2", version)
	}
}

func TestScanFilePath(t *testing.T) {
	for _, tt := range filepathScanTest {
		t.Run(tt.filepath, func(t *testing.T) {
			got, err := ScanFile(tt.filepath)
			if err != nil {
				t.Fatalf("TestScanFilePath(%s) failed, err: %s", tt.filepath, err)
			}
			if got != tt.want {
				t.Errorf("TestScanFilePath(%s) got %v, want %v", tt.filepath, got, tt.want)
			}
		})
	}
}
