// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package drweb

import (
	"regexp"
	"testing"

	multiav "github.com/saferwall/saferwall/internal/multiav"
)

type filePathTest struct {
	filepath string
	want     multiav.Result
}

var filepathScanTest = []filePathTest{
	{"../../../testdata/765c3a580f885f5e4e4f98a709e9f0ce",
		multiav.Result{Infected: true, Output: "Trojan.Siggen2.24456"}},
}

func TestVersion(t *testing.T) {
	ver, err := Version()
	if err != nil {
		t.Fatalf("TestGetVersion failed, got: %s", err)
	}

	re := regexp.MustCompile(`\d{1}\.\d{2}\.\d{2}\.\d{1,5}`)
	l := re.FindStringSubmatch(ver)
	if len(l) == 0 {
		t.Errorf(`Core engine version was incorrect, got: %s,\
		 want something similar to: 7.00.47.04280`, ver)
	}
}

func TestScanFile(t *testing.T) {
	s := Scanner{}
	for _, tt := range filepathScanTest {
		t.Run(tt.filepath, func(t *testing.T) {
			got, err := s.ScanFile(tt.filepath)
			if err != nil {
				t.Fatalf("TestScanFile(%s) failed, err: %s",
					tt.filepath, err)
			}
			if got != tt.want {
				t.Errorf("TestScanFile(%s) got %v, want %v",
					tt.filepath, got, tt.want)
			}
		})
	}
}
