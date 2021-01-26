// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package goyara

import (
	"errors"
	"fmt"
	"os"

	"github.com/hillu/go-yara"
)

// Rule represents a Yara rule.
type Rule struct {
	Namespace string
	Filename  string
}

// Load and compile yara rules
func Load(rules []Rule) (*yara.Rules, error) {

	if len(rules) == 0 {
		return nil, errors.New("no rules specified")
	}

	c, err := yara.NewCompiler()
	if err != nil {
		msg := fmt.Sprintf("failed to initialize yara compiler: %s", err)
		return nil, errors.New(msg)
	}

	for _, rule := range rules {
		f, err := os.Open(rule.Filename)
		if err != nil {
			msg := fmt.Sprintf("could not open rule file %s: %s",
				rule.Filename, err)
			return nil, errors.New(msg)
		}
		err = c.AddFile(f, rule.Namespace)
		f.Close()
		if err != nil {
			msg := fmt.Sprintf("could not parse rule file %s: %s",
				rule.Filename, err)
			return nil, errors.New(msg)
		}
	}

	r, err := c.GetRules()
	if err != nil {
		msg := fmt.Sprintf("failed to compile rules: %s", err)
		return nil, errors.New(msg)
	}

	return r, nil
}

// ScanFile performs a scan over a file path
func ScanFile(r *yara.Rules, filepath string) ([]yara.MatchRule, error) {
	m, err := r.ScanFile(filepath, 0, 0)
	return m, err
}
