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

// Scan a behavior report and extract matching rules.
func (s Scanner) Scan(apiTrace []byte) ([]Rule, error) {

	// Run the rule matching.
	eval := luar.NewLuaObjectFromName(s.L, "Eval")
	defer eval.Close()

	// Using `Call` we would get a generic `[]interface{}`, which is awkward to
	// work with. But the return type can be specified:
	results := make([]interface{}, 1)
	err := eval.Call(&results, string(apiTrace))
	if err != nil {
		return nil, err
	}

	v, err := json.Marshal(results[0])
	if err != nil {
		return nil, err
	}

	rules := make([]Rule, 0)
	err = json.Unmarshal(v, &rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

// Close the lua state object.
func (s Scanner) Close() {
	if s.L != nil {
		s.L.Close()
	}
}
