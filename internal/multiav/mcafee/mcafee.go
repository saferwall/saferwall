// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package mcafee

import (
	"regexp"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	cmd = "/opt/mcafee/uvscan"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

// Version represents all mcafee components' versions
type Version struct {
	AVEngineVersion string `json:"scancl_version"`
	VDFVersion      string `json:"vdf_version"`
	ProgramVersion  string `json:"program_version"`
}

// GetVersion get Anti-Virus scanner version.
func GetVersion() (Version, error) {

	out, err := utils.ExecCmd(cmd, "--version")

	// McAfee VirusScan Command Line for Linux64 Version: 6.0.4.564
	// Copyright (C) 2013 McAfee, Inc.

	// AV Engine version: 5600.1067 for Linux64.
	// Dat set version: 9118 created Dec 26 2018
	// Scanning for 668680 viruses, trojans and variants.

	if err != nil {
		return Version{}, err
	}

	v := Version{}
	re := regexp.MustCompile("Linux64 Version: ([\\d\\.]+)[\\s\\S]+Engine version: ([\\d\\.]+)[\\s\\S]+set version: ([\\d\\.]+)")
	l := re.FindStringSubmatch(out)
	if len(l) > 2 {
		v.ProgramVersion, v.AVEngineVersion, v.VDFVersion = l[1], l[2], l[3]
	}

	return v, nil
}

// ScanFile scans a given file.
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// Execute the scanner with the given file path
	// --ANALYZE                  : Turn on heuristic analysis for programs and macros
	// --ASCII                    : Display filenames as ASCII text.
	// --MANALYZE                 : Turn on macro heuristics.
	// --MACRO-HEURISTICS         : Turn on macro heuristics.
	// --UNZIP                    : Scan inside archive files, such as those saved in ZIP, LHA, PKarc, ARJ, TAR, CHM, and RAR.

	// /opt/mcafee/uvscan --ANALYZE --ASCII --MANALYZE --MACRO-HEURISTICS --UNZIP sample
	res.Out, err = utils.ExecCmd(cmd, "--ANALYZE", "--ASCII",
		"--MANALYZE", "--MACRO-HEURISTICS", "--UNZIP", filePath)

	// Exit codes:
	//  0 The scanner found no viruses or other potentially unwanted software, and returned no errors.
	//  2 Integrity check on DAT file failed.
	//  6 A general problem occurred.
	//  8 The scanner was unable to find a DAT file.
	//  10 A virus was found in memory.
	//  12 The scanner tried to clean a file, the attempt failed, and the file is still infected.
	//  13 The scanner found one or more viruses or hostile objects — such as a Trojan-horse program, joke program, or test file.
	//  15 The scanner’s self-check failed; the scanner may be infected or damaged.
	//  19 The scanner succeeded in cleaning all infected files.
	//  20 Scanning was prevented because of the /FREQUENCY option.
	//  21 Computer requires a reboot to clean the infection.

	if err != nil && err.Error() != "exit status 13" {
		return res, err
	}

	// McAfee VirusScan Command Line for Linux64 Version: 6.0.4.564
	// Copyright (C) 2013 McAfee, Inc.
	// (408) 988-3832 EVALUATION COPY - December 27 2018

	// AV Engine version: 5600.1067 for Linux64.
	// Dat set version: 9118 created Dec 26 2018
	// Scanning for 668680 viruses, trojans and variants.

	// /home/ubuntu/malware ... Found the RDN/Generic.tfr trojan !!!

	// Time: 00:00.00

	// Grab the detection result
	re := regexp.MustCompile(`Found the (.*) trojan`)
	l := re.FindStringSubmatch(res.Out)
	if len(l) > 0 {
		res.Output = l[1]
		res.Infected = true
	}
	return res, nil
}
