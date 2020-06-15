// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	"github.com/saferwall/saferwall/pkg/avlabel"
)

// Compilers, Installers, Packers names as seen by DiE (Detect It Easy)
// This map a signature name substring into a tag.
var sigMap = map[string]string{
	"Nullsoft Scriptable Install System": "nsis",
	"Inno Setup":                         "innosetup",
	"UPX":                                "upx",
	"FSG":                                "fsg",
	"ASPack":                             "aspack",
	"ASProtect":                          "asprotect",
	"ACProtect":                          "acprotect",
	"PECompact":                          "pecompact",
	"PE-Armor":                           "pe-armor",
	"Petite":                             "petite",
	"tElock":                             "telock",
	"EXECryptor":                         "execryptor",
	"Obsidium":                           "obsidium",
	"VMProtect":                          "vmprotect",
	"Themida/Winlicense":                 "themida-winlicense",
	"MoleBox":                            "molebox",
	"ENIGMA":                             "enigma",
	"Armadillo":                          "armadillo",
	"gcc":                                "gcc",
	"MinGW":                              "mingw",
	"Microsoft Visual C/C++":             "vc",
	"Microsoft Visual Basic":             "vb",
	"Borland C++":                        "borland-c++",
	"MASM":                               "masm",
	"FASM":                               "fasm",
	".NET":                               "dotnet",
	"MFC":                                "mfc",
	"Delphi":                             "delphi",
	"AutoIt":                             "autoit",
	"sfx: Microsoft Cabinet":             "sfx-cab",
	"Smart Assembly":                     "smart-assembly",
	".NET Reactor":                       "dotnet-reactor",
	"Confuser":                           "confuser",
	"Dotfuscator":                        "dotfuscator",
}

func (f *result) getTags() error {

	var tags []string
	f.Tags = map[string]interface{}{}

	// File format tags
	switch f.Type {
	case "pe":
		if f.PE.IsEXE() {
			tags = append(tags, "exe")
		} else if f.PE.IsDriver() {
			tags = append(tags, "sys")
		} else if f.PE.IsDLL() {
			tags = append(tags, "dll")
		}
	case "elf":
		tags = append(tags, "elf")
	case "xml":
		tags = append(tags, "xml")
	case "html":
		tags = append(tags, "html")
	case "swf":
		tags = append(tags, "swf")
	}

	f.Tags[f.Type] = tags

	// Packers/Protector/Compilers/Installers tags
	tags = nil
	for _, out := range f.Packer {
		if strings.Contains(out, "packer") ||
			strings.Contains(out, "protector") ||
			strings.Contains(out, "compiler") ||
			strings.Contains(out, "installer") ||
			strings.Contains(out, "library") {
			for sig, tag := range sigMap {
				if strings.Contains(out, sig) {
					tags = append(tags, tag)
				}
			}
		}
	}
	if tags != nil {
		f.Tags["packer"] = tags
	}

	// Multi AV tags
	avMap := f.MultiAV["last_scan"].(map[string]interface{})
	for engine, res := range avMap {
		result := res.(multiav.ScanResult)
		if result.Infected {
			switch engine {
			case "eset":
				parsedDetection := avlabel.ParseEset(result.Output)
				if len(parsedDetection) > 0 {
					f.Tags["eset"] = []string{parsedDetection["Family"]}
				}
			case "windefender":
				parsedDetection := avlabel.ParseWindefender(result.Output)
				if len(parsedDetection) > 0 {
					f.Tags["eset"] = []string{parsedDetection["Family"]}
				}
			case "avira":
				parsedDetection := avlabel.ParseAvira(result.Output)
				if len(parsedDetection) > 0 {
					f.Tags["avira"] = []string{parsedDetection["Family"]}
				}
			}
		}
	}

	return nil
}
