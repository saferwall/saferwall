// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/saferwall/saferwall/internal/constants"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/utils"
	"github.com/stevedonovan/luar"
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

	// Parse the API Trace JSON file.
	JSONAPITrace, err := utils.ReadAll(filepath.Join(*flagBhvReportPath, APITraceFilename))
	if err != nil {
		return err
	}

	// Init the lua state.
	L := luar.Init()
	defer L.Close()

	// Append the lua dependencies CPATH.
	L.DoString("package.cpath = package.cpath .. ';./lib/?.so'")

	// Execute our behavior lua code.
	err = L.DoFile("behavior.lua")
	if err != nil {
		logger.Errorf("DoFile failed with: %v", err)
		return err
	}

	// Run the rule matching.
	eval := luar.NewLuaObjectFromName(L, "Eval")
	defer eval.Close()

	// Using `Call` we would get a generic `[]interface{}`, which is awkward to
	// work with. But the return type can be specified:
	results := make([]interface{}, 1)
	err = eval.Call(&results, string(JSONAPITrace))
	if err != nil {
		logger.Errorf("eval.Call failed with: %v", err)
		return err
	}

	logger.Info(results[0])

	return nil
}
