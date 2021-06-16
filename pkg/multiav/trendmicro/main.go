// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trendmicro

import (
	"io/ioutil"
	"path"
	"regexp"
	"time"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	splxmain = "/opt/TrendMicro/SProtectLinux/SPLX.vsapiapp/splxmain"
	logDir   = "/var/log/TrendMicro/SProtectLinux/"
	tmsplx   = "/opt/TrendMicro/SProtectLinux/tmsplx.xml"
)

// Result represents detection results
type Result struct {
	Infected bool   `json:"infected"`
	Output   string `json:"output"`
}

// Version represents all avira components' versions
type Version struct {
	EngineVersion         string `json:"engine_version"`
	PatternVersion        string `json:"pattern_version"`
	SpywarePatternVersion string `json:"spyware_pattern_version"`
}

// GetVersion returns Engine, Virus and Patterns versions
func GetVersion() (Version, error) {

	ver := Version{}

	// The information regarding the versions is found on:
	// /opt/TrendMicro/SProtectLinux/tmsplx.xml
	out, err := utils.ReadAll(tmsplx)
	if err != nil {
		return Version{}, err
	}
	data := string(out)

	// Interesting keys:
	// <P Name="EngineVersion" Value="12.000.1008"/>
	// <P Name="PatternVersion" Value="16.259.00"/>
	// <P Name="SpywarePatternVersion" Value="2.337.00"/>
	engineReg := regexp.MustCompile(`<P Name="EngineVersion" Value="(.*)"/>`)
	l := engineReg.FindStringSubmatch(data)
	if len(l) > 0 {
		ver.EngineVersion = l[1]
	}

	patternReg := regexp.MustCompile(`<P Name="PatternVersion" Value="(.*)"/>`)
	l = patternReg.FindStringSubmatch(data)
	if len(l) > 0 {
		ver.PatternVersion = l[1]
	}

	spywareReg := regexp.MustCompile(`<P Name="SpywarePatternVersion" Value="(.*)"/>`)
	l = spywareReg.FindStringSubmatch(data)
	if len(l) > 0 {
		ver.SpywarePatternVersion = l[1]
	}

	return ver, nil
}

// ScanFile scans a given file.
func ScanFile(filepath string) (Result, error) {

	// TrendMicro does not really provide a simple way to grab the results
	// We need to read it from a log file. We start by deleting all files
	// from previous scan results.
	err := utils.DeleteDirContent(logDir)
	if err != nil {
		return Result{}, err
	}

	// splxmain does not seem to be able to scan a file directly,
	// it only take a directory as argument.
	// So we create a tmp dir and we make a copy of the file.
	tempDir, err := ioutil.TempDir("/tmp", "trendmicro")
	if err != nil {
		return Result{}, err
	}

	filename := path.Base(filepath)
	filePathCopy := path.Join(tempDir, filename)
	err = utils.CopyFile(filepath, filePathCopy)
	if err != nil {
		return Result{}, err
	}

	// Execute the scanner with the given file path.
	_, err = utils.ExecCommand(splxmain, "-m", path.Dir(filePathCopy))
	if err != nil {
		return Result{}, err
	}

	// The logs are found in files that have the following pattern:
	// Virus.20201001.0001 or Spyware.20201001.0001.
	currentTime := time.Now()
	todayDate := currentTime.Format("20060102")
	virusLog := path.Join(logDir, "Virus."+todayDate+".0001")

	// The content of the Virus.XYZ and Spyware.XYZ log looks like:
	// ...
	// first_action=0
	// virus_name=Ransom.Win32.GANDCRAB.SMLV2.hp
	// function_code=12
	// ...
	re := regexp.MustCompile("virus_name=(.*)")

	// Read the Virus.XYZ log
	res := Result{}
	data, virusErr := utils.ReadAll(virusLog)
	if virusErr == nil && len(data) > 0 {
		l := re.FindStringSubmatch(string(data))
		if len(l) > 0 {
			res.Output = l[1]
			res.Infected = true
			return res, nil
		}
	}

	// Might be a spyware.
	spywareLog := path.Join(logDir, "Spyware."+todayDate+".0001")
	data, spywareErr := utils.ReadAll(spywareLog)
	if spywareErr == nil && len(data) > 0 {
		l := re.FindStringSubmatch(string(data))
		if len(l) > 0 {
			res.Output = l[1]
			res.Infected = true
			return res, nil
		}
	}

	if virusErr != nil {
		return Result{}, virusErr
	}
	if spywareErr != nil {
		return Result{}, spywareErr
	}

	return Result{}, nil
}
