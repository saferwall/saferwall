// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"testing"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	sdkDir = "C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0"
)

var rePrototypetests = []struct {
	in  string
	out int
}{
	{sdkDir + "\\um\\fileapi.h", 96},
	{sdkDir + "\\um\\processthreadsapi.h", 93},
}

var reStructtests = []struct {
	in  string
	out int
}{
	{sdkDir + "\\shared\\bcrypt.h", 40},
	{sdkDir + "\\um\\debugapi.h", 0},
	{sdkDir + "\\um\\fileapi.h", 5},
	{sdkDir + "\\um\\libloaderapi.h", 3},
	{sdkDir + "\\um\\memoryapi.h", 2},
	{sdkDir + "\\um\\processthreadsapi.h", 10},
	{sdkDir + "\\um\\sysinfoapi.h", 2},
	{sdkDir + "\\um\\wininet.h", 49},
	{sdkDir + "\\um\\tlhelp32.h", 7},
}

var parseStructTests = []struct {
	path string
	in   string
	out  int
}{
	{sdkDir + "\\um\\processthreadsapi.h", "PROCESS_INFORMATION", 4},
}

func TestGetAPIPrototypes(t *testing.T) {
	for _, tt := range rePrototypetests {
		t.Run(tt.in, func(t *testing.T) {
			data, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetAPIPrototypes(%s) failed, got: %s", tt.in, err)
			}

			r := regexp.MustCompile(RegAPIs)
			matches := r.FindAllString(string(data), -1)
			got := len(matches)
			if got != tt.out {
				t.Errorf("TestGetAPIPrototypes(%s) got %v, want %v", tt.in, got, tt.out)
			}
		})
	}
}

func TestGetStructs(t *testing.T) {
	for _, tt := range reStructtests {
		t.Run(tt.in, func(t *testing.T) {
			data, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetStructs(%s) failed, got: %s", tt.in, err)
			}

			matches, _ := getAllStructs(data)
			got := len(matches)
			if got != tt.out {
				t.Errorf("TestGetStructs(%s) got %v, want %v", tt.in, got, tt.out)
			}
		})
	}
}

func TestParseStruct(t *testing.T) {
	for _, tt := range parseStructTests {
		t.Run(tt.in, func(t *testing.T) {

			data, err := utils.ReadAll(tt.path)
			if err != nil {
				t.Errorf("ReadAll(%s) failed with : %s", tt.path, err)
			}

			_, matches := getAllStructs(data)
			for _, structObj := range matches {
				if structObj.Name == tt.in {
					got := len(structObj.Members)
					if got != tt.out {
						t.Errorf("TestParseStruct(%s) got %v, want %v", tt.in, got, tt.out)
					} else {
						for _, member := range structObj.Members {
							if member.Name == "" || member.Type == "" {
								t.Errorf("TestParseStruct(%s) empty members", tt.in)
							}
						}
					}

				}
			}
		})
	}
}
