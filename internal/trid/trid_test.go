// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trid

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func getAbsoluteFilePath(testfile string) string {
	_, p, _, _ := runtime.Caller(0)
	return path.Join(filepath.Dir(p), "..", "..", testfile)
}

var tridTests = []struct {
	in  string
	out string
}{
	{getAbsoluteFilePath("testdata/putty.exe"),
		"(.EXE) Win64 Executable (generic)",
	},
}

func TestScan(t *testing.T) {
	for _, tt := range tridTests {
		t.Run(tt.in, func(t *testing.T) {
			filePath := tt.in
			got, err := Scan(filePath)
			if err != nil {
				t.Errorf("TestScan(%s) got %v, want %v", tt.in, err, tt.in)
			}
			if len(got) <= 0 || !strings.Contains(got[0], tt.out) {
				t.Errorf("TestScan(%s) got %v, want %v", tt.in, got, tt.out)
			}
		})
	}
}
