// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package behavior

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type Scanner struct {
	L *lua.State
}

// Rule describes a behavior rule.
type Rule struct {
	// Description describes the purpose of the rule.
	Description string `json:"description"`
	// ID uniquely identify the rule.
	ID string `json:"id"`
	// Category indicates the category of the behavior rules.
	// examples include: anti-analysis, ransomware, ..
	Category string `json:"category"`
	// Severity indicates how confident the rule is to classify
	// the threat as malicious.
	Severity string `json:"severity"`
}

// Event represents a system event: a registry, network or file event.
type Event struct {
	// Process identifier responsible for generating the event.
	ProcessID string `json:"proc_id"`
	// Type of the system event.
	Type string `json:"type"`
	// Path of the system event. For instance, when the event is of type:
	// `registry`, the path represents the registry key being used. For a
	// `network` event type, the path is the IP or domain used.
	Path string `json:"path"`
	// Th operation requested over the above `Path` field. This field means
	// different things according to the type of the system event.
	// - For file system events: can be either: create, read, write, delete, rename, ..
	// - For registry events: can be either: create, rename, set, delete.
	// - For network events: this represents the protocol of the communication, can
	// be either HTTP, HTTPS, FTP, FTP
	Operation string `json:"operation"`
}

// ScanResult represents the behavior rules scan results.
type ScanResult struct {
	Rules  []Rule  `json:"matches"`
	Events []Event `json:"events"`
}

const (
	behaviorLuaFile = "behavior.lua"
)

func New(behaviorRules string) (Scanner, error) {

	L := luar.Init()

	// Append the lua dependencies CPATH.
	luaCode := fmt.Sprintf("package.cpath = package.cpath .. ';./%s/?.so'",
		behaviorRules)
	err := L.DoString(luaCode)
	if err != nil {
		return Scanner{}, err
	}

	// Execute lua file.
	luaFilePath := filepath.Join(behaviorRules, behaviorLuaFile)
	err = L.DoFile(luaFilePath)
	if err != nil {
		return Scanner{}, err
	}

	return Scanner{L: L}, nil

}

// Scan a behavior report and extract system events and matching rules.
func (s Scanner) Scan(apiTrace []byte) (ScanResult, error) {

	// Run the rule matching.
	eval := luar.NewLuaObjectFromName(s.L, "Eval")
	defer eval.Close()

	// Using `Call` we would get a generic `[]interface{}`, which is awkward to
	// work with. But the return type can be specified:
	results := make([]interface{}, 1)
	err := eval.Call(&results, string(apiTrace))
	if err != nil {
		return ScanResult{}, err
	}

	v, err := json.Marshal(results[0])
	if err != nil {
		return ScanResult{}, err
	}

	scanResult := ScanResult{}

	err = json.Unmarshal(v, &scanResult)
	if err != nil {
		return ScanResult{}, err
	}

	return scanResult, nil

}

// Close the lua state object.
func (s Scanner) Close() {
	if s.L != nil {
		s.L.Close()
	}
}
