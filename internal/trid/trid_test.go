// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trid

import (
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func getAbsoluteFilePath(testfile string) string {
	_, p, _, _ := runtime.Caller(0)
	return path.Join(filepath.Dir(p), "..", "..", testfile)
}

var tridtests = []struct {
	in  string
	out []string
}{
	{getAbsoluteFilePath("testdata/putty.exe"),
		[]string{
			"42.7% (.EXE) Win32 Executable (generic) (4505/5/1)",
			"19.2% (.EXE) OS/2 Executable (generic) (2029/13)",
			"19.0% (.EXE) Generic Win/DOS Executable (2002/3)",
			"18.9% (.EXE) DOS Executable Generic (2000/1)",
		},
	},
}

func TestScan(t *testing.T) {
	for _, tt := range tridtests {
		t.Run(tt.in, func(t *testing.T) {
			filePath := tt.in
			got, err := Scan(filePath)
			if err != nil {
				t.Errorf("TestScan(%s) got %v, want %v", tt.in, err, tt.in)
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("TestScan(%s) got %v, want %v", tt.in, got, tt.out)
			}
		})
	}
}
