// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package fsecure

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
	{
		"../../test/testdata/765c3a580f885f5e4e4f98a709e9f0ce",
		multiav.Result{
			Infected: true, Output: "Trojan-Downloader:W32/Renos.GTR"},
	},
}

func TestProgramVersion(t *testing.T) {
	ver, err := GetVersion()
	if err != nil {
		t.Fatalf("TestProgramVersion failed, got: %s", err)
	}

	engineRegex := regexp.MustCompile(`\d{1}\.\d{1,2} build \d{1,2}`)
	l := engineRegex.FindStringSubmatch(ver.AquariusEngineVersion)
	if len(l) == 0 {
		t.Fatalf(`Aquarius Engine Version was incorrect, got: %s,\
		 want something similar to: 1.0 build 8`, ver)
	}
	l = engineRegex.FindStringSubmatch(ver.HydraEngineVersion)
	if len(l) == 0 {
		t.Fatalf(`Hydra Engine Version was incorrect, got: %s,\
		 want something similar to: 5.22 build 28`, ver)
	}
	l = engineRegex.FindStringSubmatch(ver.FSecureVersion)
	if len(l) == 0 {
		t.Fatalf(`FSecure Version was incorrect, got: %s,\
		 want something similar to: 11.10 build 68`, ver)
	}

	dbRegex := regexp.MustCompile(`[\d-_]{10,}`)
	l = dbRegex.FindStringSubmatch(ver.AquariusDatabaseVersion)
	if len(l) == 0 {
		t.Fatalf(`Aquarius Database Version was incorrect, got: %s,\
		 want something similar to: 2018-12-21_08`, ver)
	}
	l = dbRegex.FindStringSubmatch(ver.HydraDatabaseVersion)
	if len(l) == 0 {
		t.Fatalf(`Hydra Database Version was incorrect, got: %s,\
		 want something similar to: 2018-12-21_02`, ver)
	}
	l = dbRegex.FindStringSubmatch(ver.DatabaseVersion)
	if len(l) == 0 {
		t.Fatalf(`Database Version was incorrect, got: %s,\
		 want something similar to: 2018-12-21_08`, ver)
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
