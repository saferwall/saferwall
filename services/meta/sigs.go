// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package meta

// Compilers, Installers, Packers names as seen by DiE (Detect It Easy)
// This map a signature name substring into a tag.
var sigMap = map[string]string{
	"Nullsoft Scriptable Install System": "nsis",
	"Inno Setup":                         "innosetup",
	"UPX":                                "upx",
	"FSG":                                "fsg",
	"ASPack":                             "aspack",
	"RLPack":                             "rlpack",
	"ASProtect":                          "asprotect",
	"ACProtect":                          "acprotect",
	"PECompact":                          "pecompact",
	"PECrypt32":                          "pecrypt32",
	"PE-Armor":                           "pe-armor",
	"Petite":                             "petite",
	"PELock":                             "pelock",
	"tElock":                             "telock",
	"EXECryptor":                         "execryptor",
	"Obsidium":                           "obsidium",
	"VMProtect":                          "vmprotect",
	"Themida/Winlicense":                 "themida-winlicense",
	"MoleBox":                            "molebox",
	"ENIGMA":                             "enigma",
	"MPRESS":                             "mpress",
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
	"Yoda's Crypter":                     "yodascrypter",
	"Delphi":                             "delphi",
	"AutoIt":                             "autoit",
	"StarForce":                          "starforce",
	"eXPressor":                          "expressor",
	"sfx: Microsoft Cabinet":             "sfx-cab",
	"Smart Assembly":                     "smart-assembly",
	".NET Reactor":                       "dotnet-reactor",
	"Confuser":                           "confuser",
	"Dotfuscator":                        "dotfuscator",
}

var typeMap = map[string]string{
	"PE32":                    "pe",
	"MS-DOS":                  "msdos",
	"XML":                     "xml",
	"HTML":                    "html",
	"ELF":                     "elf",
	"PDF":                     "pdf",
	"Macromedia Flash":        "swf",
	"Zip archive data":        "zip",
	"Java archive data (JAR)": "jar",
	"PEG image data":          "jpeg",
	"PNG image data":          "png",
	"SVG Scalable Vector":     "svg",
}
