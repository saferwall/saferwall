// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"testing"
	"log"

	"github.com/dlclark/regexp2"


	"github.com/saferwall/saferwall/pkg/utils"
)

var (
	RegStructs2 = `typedef [\w() ]*struct [\w]+[\n\s]+{(.|\n)+?} (?!DUMMYSTRUCTNAME|DUMMYUNIONNAME)[\w, *]+;`
)

var rePrototypetests = []struct {
	in  string
	out int
}{
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\fileapi.h", 94},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\processthreadsapi.h", 85},
}

var reStructtests = []struct {
	in  string
	out int
}{
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\shared\\bcrypt.h", 27},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\debugapi.h", 0},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\fileapi.h", 5},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\libloaderapi.h", 3},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\memoryapi.h", 2},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\processthreadsapi.h", 10},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\sysinfoapi.h", 2},
	{"C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0\\um\\tlhelp32.h", 7},
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

			r := regexp2.MustCompile(RegStructs2, 0)
			matches := regexp2FindAllString(r, string(data))
			log.Print(matches)
			got := len(matches)
			if got != tt.out {
				t.Errorf("TestGetStructs(%s) got %v, want %v", tt.in, got, tt.out)
			}
		})
	}
}
