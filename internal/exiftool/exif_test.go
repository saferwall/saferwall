// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package exiftool

import (
	"testing"
)

var magictests = []struct {
	in  string
	out string
}{
	{"../../testdata/putty.exe", "Win32 EXE"},
	{"../../testdata/ls", "ELF shared library"},
}

var failuretests = []struct {
	in  string
	out string
}{
	{"../../testdata/notfound", ""},
}

func TestExiftoolScan(t *testing.T) {
	t.Run("TestConsistentOutput", func(t *testing.T) {
		for _, tt := range magictests {
			filePath := tt.in
			got, err := Scan(filePath)
			if err != nil {
				t.Errorf("TestMagicScan(%s) got %v, want %v",
					tt.in, err, tt.in)
			}
			if got["FileType"] != tt.out {
				t.Errorf("TestMagicScan(%s) got %v, want %v",
					tt.in, got["FileType"], tt.out)
			}
		}
	})
	t.Run("TestExpectedFailure", func(t *testing.T) {
		for _, tt := range failuretests {
			filePath := tt.in
			got, err := Scan(filePath)
			if got != nil || err == nil {
				t.Errorf("TestMagicScan(%s) expected error , got %v",
					tt.in,
					err)
			}
		}
	})
}
