// Copyright 2018 Saferwall. All rights reserved.
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
	{"../../test/multiav/clean/putty.exe", "PE32"},
}

func TestExiftoolScan(t *testing.T) {
	for _, tt := range magictests {
		t.Run(tt.in, func(t *testing.T) {
			filePath := tt.in
			got, err := Scan(filePath)
			if err != nil {
				t.Errorf("TestMagicScan(%s) got %v, want %v", tt.in, err, tt.in)
			}
			if got["PeType"] != tt.out {
				t.Errorf("TestMagicScan(%s) got %v, want %v", tt.in, got, tt.out)
			}
		})
	}
}
