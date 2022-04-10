// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avast

import (
	"regexp"
	"testing"

	multiav "github.com/saferwall/saferwall/internal/multiav"
)

type filescanTest struct {
	filepath string
	want     multiav.Result
}

var filescanTests = []filescanTest{
	{"../../test/testdata/765c3a580f885f5e4e4f98a709e9f0ce",
		multiav.Result{Infected: true, Output: "Win32:CodecPack-P [Trj]"}},
}

func TestVPSVersion(t *testing.T) {
	version, err := VPSVersion()
	if err != nil {
		t.Fatalf("TestVPSVersion failed, got: %s", err)

	}

	re := regexp.MustCompile(`\d{8}`)
	l := re.FindStringSubmatch(version)
	if len(l) == 0 {
		t.Fatalf(`VPS version was incorrect, got: %s,\
		 want something similar to: 19010602`, version)
	}
}

func TestProgramVersion(t *testing.T) {
	version, err := ProgramVersion()
	if err != nil {
		t.Fatalf("TestProgramVersion failed, got: %s", err)
	}

	re := regexp.MustCompile(`\d{1}\.\d{1}\.\d{1}`)
	l := re.FindStringSubmatch(version)
	if len(l) == 0 {
		t.Fatalf(`Program version was incorrect, got: %s, \
		want something similar to: 2.2.0`, version)
	}
}

func TestScanFile(t *testing.T) {
	s := Scanner{}
	for _, tt := range filescanTests {
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
