// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trendmicro

import (
	"path/filepath"
	"runtime"
	"testing"

	multiav "github.com/saferwall/saferwall/internal/multiav"
)

type filePathTest struct {
	filepath string
	want     multiav.Result
}

var filepathScanTest = []filePathTest{
	{
		absPath("../../test/testdata/765c3a580f885f5e4e4f98a709e9f0ce"),
		multiav.Result{Infected: true, Output: "TROJ_KRYPTK.SMCA"},
	},
	{
		absPath("../../test/testdata/2762e6c9679a8174ba9a60981f63a78f"),
		multiav.Result{Infected: true, Output: "HTML_IFRAME.LCA"},
	},
}

func absPath(testfile string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), testfile)
	return dir
}

func TestStartDaemon(t *testing.T) {
	err := StartDaemon()
	if err != nil {
		t.Fatalf("TestScanFile failed, err: %s", err)
	}
}

func TestScanFile(t *testing.T) {
	s := Scanner{}
	for _, tt := range filepathScanTest {
		t.Run(tt.filepath, func(t *testing.T) {

			got, err := s.ScanFile(tt.filepath)
			if err != nil {
				t.Fatalf("TestScanFile(%s) failed, err: %v",
					tt.filepath, err)
			}
			if got.Output != tt.want.Output {
				t.Errorf("TestScanFile(%s) got %v, want %v",
					tt.filepath, got, tt.want)
			}
		})
	}
}
