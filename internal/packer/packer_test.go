// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package packer

import (
	"reflect"
	"testing"
)

var magictests = []struct {
	in  string
	out []string
}{
	{"../../testdata/putty.exe",
		[]string{
			"PE: compiler: Microsoft Visual C/C++(2015 v.14.0)[-]",
			"PE: linker: unknown(14.0)[EXE32,signed]",
		},
	},
}

func TestPackerScan(t *testing.T) {
	for _, tt := range magictests {
		t.Run(tt.in, func(t *testing.T) {
			filePath := tt.in
			got, err := Scan(filePath)
			if err != nil {
				t.Errorf("TestPackerScan(%s) got %v, want %v",
					tt.in, err, tt.in)
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("TestPackerScan(%s) got %v, want %v",
					tt.in, got, tt.out)
			}
		})
	}
}
