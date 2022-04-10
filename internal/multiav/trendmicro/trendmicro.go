// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trendmicro

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"

	multiav "github.com/saferwall/saferwall/internal/multiav"
	"github.com/saferwall/saferwall/internal/utils"
)

const (
	splxmain = "/opt/TrendMicro/SProtectLinux/SPLX.vsapiapp/splxmain"
	logDir   = "/var/log/TrendMicro/SProtectLinux/"
	tmsplx   = "/opt/TrendMicro/SProtectLinux/tmsplx.xml"
	splx     = "/etc/init.d/splx"
)

// Scanner represents an empty struct that can be used to a method received.
type Scanner struct{}

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
func (Scanner) ScanFile(filePath string) (multiav.Result, error) {

	var err error
	res := multiav.Result{}

	// TrendMicro does not really provide a simple way to grab the results
	// We need to read it from a log file. We start by deleting all files
	// from previous scan results.
	err = utils.DeleteDirContent(logDir)
	if err != nil {
		return res, err
	}

	// splxmain does not seem to be able to scan a file directly,
	// it only take a directory as argument.
	// So we create a tmp dir and we make a copy of the file.
	tempDir, err := ioutil.TempDir("/tmp", "trendmicro")
	if err != nil {
		return res, err
	}

	filename := path.Base(filePath)
	filePathCopy := path.Join(tempDir, filename)
	err = utils.CopyFile(filePath, filePathCopy)
	if err != nil {
		return res, err
	}

	// Execute the scanner with the given file path.
	out, err := utils.ExecCmd("sudo", splxmain, "-m", path.Dir(filePathCopy))
	fmt.Printf("out: %s", out)
	if err != nil {
		return res, err
	}

	// Wait until the scan is finished.
	currentTime := time.Now()
	todayDate := currentTime.Format("20060102")
	scanLog := path.Join(logDir, "Scan."+todayDate+".0001")
	for i := 0; i < 5; i++ {

		time.Sleep(time.Second * 1)

		_, err = utils.ExecCmd("sudo", "chmod", "-R", "0777", logDir)
		if err != nil {
			return res, err
		}
		file, err := os.Open(scanLog)
		if err != nil {
			continue
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			continue
		}
		if fi.Size() > 0 {
			break
		}
	}

	// The logs are found in files that have the following pattern:
	// Virus.20201001.0001 or Spyware.20201001.0001.
	virusLog := path.Join(logDir, "Virus."+todayDate+".0001")

	// The content of the Virus.XYZ and Spyware.XYZ log looks like:
	// ...
	// first_action=0
	// virus_name=Ransom.Win32.GANDCRAB.SMLV2.hp
	// function_code=12
	// ...
	re := regexp.MustCompile("virus_name=(.*)")

	// Read the Virus.XYZ log
	data, virusErr := utils.ReadAll(virusLog)
	res.Out = "VirusOut: " + string(data)
	if virusErr == nil && len(res.Out) > 0 {
		l := re.FindStringSubmatch(res.Out)
		if len(l) > 0 {
			res.Output = l[1]
			res.Infected = true
			return res, nil
		}
	}

	// Might be a spyware.
	spywareLog := path.Join(logDir, "Spyware."+todayDate+".0001")
	data, spywareErr := utils.ReadAll(spywareLog)
	res.Out += "SpywareOut: " + string(data)
	if spywareErr == nil && len(data) > 0 {
		l := re.FindStringSubmatch(res.Out)
		if len(l) > 0 {
			res.Output = l[1]
			res.Infected = true
			return res, nil
		}
	}

	if virusErr != nil {
		return res, virusErr
	}
	if spywareErr != nil {
		return res, spywareErr
	}

	return res, nil
}

// StartDaemon starts the trendmicro daemon.
func StartDaemon() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out, err := utils.ExecCmdWithContext(ctx, "sudo", splx, "restart")
	if err != nil {
		return fmt.Errorf("failed to start daemon, err: %v, out:%s",
			err, out)
	}

	return nil
}
