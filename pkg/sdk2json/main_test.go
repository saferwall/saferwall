// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	sdkDir = "C:\\Program Files (x86)\\Windows Kits\\10\\Include\\10.0.19041.0"
)

var rePrototypeTests = []struct {
	in  string
	out int
}{
	{sdkDir + "\\um\\fileapi.h", 96},
	{sdkDir + "\\um\\processthreadsapi.h", 93},
}

var reTypedefsTests = []struct {
	in  string
	out int
}{
	
	{sdkDir + "\\shared\\ntdef.h", 186},
	{sdkDir + "\\shared\\minwindef.h", 42},
	{sdkDir + "\\shared\\basetsd.h", 42},
	{sdkDir + "\\um\\winnt.h", 270},
	{sdkDir + "\\um\\minwinbase.h", 12},
}

var reStructTests = []struct {
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
	out  Struct
}{
	{sdkDir + "\\um\\processthreadsapi.h", "PROCESS_INFORMATION", Struct{
		Name:             "PROCESS_INFORMATION",
		TypedefName:      "_PROCESS_INFORMATION",
		PointerAlias:     "PPROCESS_INFORMATION",
		LongPointerAlias: "LPPROCESS_INFORMATION",
		Members: []StructMember{
			{
				Name: "hProcess",
				Type: "HANDLE",
			},
			{
				Name: "hThread",
				Type: "HANDLE",
			},
			{
				Name: "dwProcessId",
				Type: "DWORD",
			},
			{
				Name: "dwThreadId",
				Type: "DWORD",
			},
		}},
	},
}

func TestGetAPIPrototypes(t *testing.T) {
	for _, tt := range rePrototypeTests {
		t.Run(tt.in, func(t *testing.T) {
			data, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("ReadAll(%s) failed, got: %s", tt.in, err)
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


func TestParseTypedefs(t *testing.T) {
	for _, tt := range reTypedefsTests {
		t.Run(tt.in, func(t *testing.T) {
			data, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("ReadAll(%s) failed, got: %s", tt.in, err)
			}
			parseTypedefs(data)
			got := len(typedefs)
			if got != tt.out {
				t.Errorf("TestParseTypedefs(%s) got %v, want %v", tt.in, got, tt.out)
			}
			for k := range typedefs {
				delete(typedefs, k)
			}

		})
	}
}

func TestGetStructs(t *testing.T) {
	for _, tt := range reStructTests {
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
			for _, got := range matches {
				if got.Name == tt.in {
					if !reflect.DeepEqual(got, tt.out) {
						t.Errorf("TestParseStruct(%s) got %v, want %v", tt.in, got, tt.out)
					}

				}
			}
		})
	}
}
