// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"strings"
)

// Compilers, Installers, Packers names as seen by DiE (Detect It Easy)
const (
	SigNSIS      = "Nullsoft Scriptable Install System"
	SigInnoSetup = "Inno Setup"
	SigUPX       = "UPX"
	SigASPack    = "ASPack"
	SigASProtect = "ASProtect"
	SigPECompact = "PECompact"
	SigEnigma    = "ENIGMA"
	SigGCC       = "gcc"
	SigMSVC      = "Microsoft Visual C/C++"
	SigMSVB      = "Microsoft Visual Basic"
	SigDotNet    = ".NET"
	SigMFC       = "MFC"
	SigDelphi    = "Delphi"
	SigAutoIT    = "AutoIt"
	SigSFXCab    = "sfx: Microsoft Cabinet"
)

// GetTags get file tags.
func (f *result) GetTags() error {

	var tags []string
	packer := f.Packer[0]

	if strings.Contains(packer, SigNSIS) {
		tags = append(tags, "nsis")
	} else if strings.Contains(packer, SigInnoSetup) {
		tags = append(tags, "innosetup")
	} else if strings.Contains(packer, SigUPX) {
		tags = append(tags, "upx")
	} else if strings.Contains(packer, SigUPX) {
		tags = append(tags, "upx")
	} else if strings.Contains(packer, SigASPack) {
		tags = append(tags, "aspack")
	} else if strings.Contains(packer, SigASProtect) {
		tags = append(tags, "asprotect")
	} else if strings.Contains(packer, SigPECompact) {
		tags = append(tags, "pecompact")
	} else if strings.Contains(packer, SigPECompact) {
		tags = append(tags, "enigma")
	} else if strings.Contains(packer, SigMSVC) {
		tags = append(tags, "vc")
	} else if strings.Contains(packer, SigMSVB) {
		tags = append(tags, "vb")
	} else if strings.Contains(packer, SigDotNet) {
		tags = append(tags, "dotnet")
	} else if strings.Contains(packer, SigMFC) {
		tags = append(tags, "mfc")
	} else if strings.Contains(packer, SigDelphi) {
		tags = append(tags, "delphi")
	} else if strings.Contains(packer, SigAutoIT) {
		tags = append(tags, "autoit")
	} else if strings.Contains(packer, SigSFXCab) {
		tags = append(tags, "sfx-cab")
	} else if strings.Contains(packer, SigGCC) {
		tags = append(tags, "gcc")
	}

	log.Println(tags)
	f.Tags = tags
	return nil
}
