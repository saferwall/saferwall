// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package magic

import (
	"testing"
)

func TestMagicScan(t *testing.T) {
	var magictests = []struct {
		in  string
		out string
	}{
		{
			"../../testdata/putty.exe",
			"PE32 executable (GUI) Intel 80386, for MS Windows",
		},
		{
			"",
			"",
		},
	}
	for _, tt := range magictests {

		filePath := tt.in
		got, err := Scan(filePath)
		if err != nil {
			t.Errorf("TestMagicScan(%s) got %v, want %v",
				tt.in, err, tt.in)
		}
		if got != tt.out {
			t.Errorf("TestMagicScan(%s) got %v, want %v",
				tt.in, got, tt.out)
		}

	}
}
