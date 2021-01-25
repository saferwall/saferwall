// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trendmicro

import (
	"path"
	"testing"

	"github.com/saferwall/saferwall/pkg/utils"
)

type filePathTest struct {
	filepath string
	want     Result
}

func getAbsoluteFilePath(testfile string) string {
	return path.Join(utils.GetRootProjectDir(), testfile)
}

var filepathScanTest = []filePathTest{
	{getAbsoluteFilePath("test/multiav/clean/eicar.com"),
	 Result{Infected: true, Output: "Eicar_test_file"}},
}

func TestScanFilePath(t *testing.T) {
	for _, tt := range filepathScanTest {
		t.Run(tt.filepath, func(t *testing.T) {
			got, err := ScanFile(tt.filepath)
			if err != nil {
				t.Fatalf("TestScanFilePath(%s) failed, err: %s",
				 tt.filepath, err)
			}
			if got != tt.want {
				t.Errorf("TestScanFilePath(%s) got %v, want %v",
				 tt.filepath, got, tt.want)
			}
		})
	}
}
