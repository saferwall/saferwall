// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package drweb

import (
	"testing"
	"regexp"
)

type filePathTest struct {
	filepath string
	want     Result
}


var filepathScanTest = []filePathTest{
	{"../../../test/multiav/clean/eicar.com", Result{Infected: true, Output: "EICAR Test File (NOT a Virus!)"}},
}

func TestGetVersion(t *testing.T) {
	ver, err := GetVersion()
	if err != nil {
		t.Fatalf("TestGetVersion failed, got: %s", err)
	}
	
	re := regexp.MustCompile(`\d{1}\.\d{2}\.\d{2}\.\d{1,5}`)
	l := re.FindStringSubmatch(ver.CoreEngineVersion)
	if len(l) == 0 {
		t.Errorf("Core enfine version was incorrect, got: %s, want something similar to: 7.00.47.04280", ver)
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
