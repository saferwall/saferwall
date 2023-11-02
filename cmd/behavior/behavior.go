// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/saferwall/saferwall/internal/constants"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/utils"
	"github.com/saferwall/saferwall/services/sandbox"
	lua "github.com/yuin/gopher-lua"
)

const (
	APITraceFilename = "api_trace.json"
)

var flagBhvReportPath = flag.String("bhv-path", "",
	"File path of the behavior report scan results")

func main() {

	flag.Parse()

	// Create root logger tagged with server version.
	logger := log.New().With(context.TODO(), "version", constants.Version)
	if err := run(logger); err != nil {
		logger.Errorf("failed to run the server: %s", err)
		os.Exit(-1)
	}
}

func run(logger log.Logger) error {

	// Extract the GUID and the SHA256 from the file path.
	parts := strings.Split(*flagBhvReportPath, string(os.PathSeparator))
	guid := parts[len(parts)-2]
	sha256 := parts[len(parts)-3]
	logger.Infof("processing behavior report for %s : %s", sha256, guid)

	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("load.lua"); err != nil {
		panic(err)
	}

	// Parse the API Trace JSON file.
	JSONAPITrace, err := utils.ReadAll(filepath.Join(*flagBhvReportPath, APITraceFilename))
	if err != nil {
		return err
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("Eval"),
		NRet:    1,
		Protect: true,
	}, lua.LString(JSONAPITrace)); err != nil {
		panic(err)
	}
	ret := L.Get(-1) // returned value
	L.Pop(1)         // remove received value
	fmt.Print(ret)

	var w32APIs []sandbox.Win32API
	err = json.Unmarshal(JSONAPITrace, &w32APIs)
	if err != nil {
		return err
	}

	return nil
}
