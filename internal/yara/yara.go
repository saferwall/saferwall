// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package goyara

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yara "github.com/hillu/go-yara/v4"
)

// Rule represents a Yara rule.
type Rule struct {
	Namespace string
	Filename  string
}

type Scanner struct {
	scanner *yara.Scanner
}

func New(rulesPath string) (Scanner, error) {
	rules := []Rule{}

	files, err := os.ReadDir(rulesPath)
	if err != nil {
		return Scanner{}, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yar") || strings.HasSuffix(file.Name(), ".yara") {
			filePath := filepath.Join(rulesPath, file.Name())
			rules = append(rules, Rule{Namespace: file.Name(), Filename: filePath})
		}
	}

	r, err := Load(rules)
	if err != nil {
		return Scanner{}, err
	}

	s, err := yara.NewScanner(r)
	if err != nil {
		return Scanner{}, err
	}

	scanner := Scanner{scanner: s}
	return scanner, nil

}

func NewFromRules(rules []Rule) (Scanner, error) {
	r, err := Load(rules)
	if err != nil {
		return Scanner{}, err
	}

	s, err := yara.NewScanner(r)
	if err != nil {
		return Scanner{}, err
	}

	scanner := Scanner{scanner: s}
	return scanner, nil

}

// Load and compile yara rules.
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

// ScanFile performs a scan over a file path.
func (s Scanner) ScanFile(filepath string) ([]yara.MatchRule, error) {

	var m yara.MatchRules

	err := s.scanner.SetCallback(&m).ScanFile(filepath)

	return m, err
}

// ScanBytes performs a scan over a byte stream.
func (s Scanner) ScanBytes(buff []byte) ([]yara.MatchRule, error) {

	var m yara.MatchRules

	err := s.scanner.SetCallback(&m).ScanMem(buff)

	return m, err
}

// ScanProc performs a process scan.
func (s Scanner) ScanProc(pid int) ([]yara.MatchRule, error) {

	var m yara.MatchRules

	err := s.scanner.SetCallback(&m).ScanProc(pid)

	return m, err
}

// StringifyMatches return the list of matched yar rule names.
func (s Scanner) StringifyMatches(matches []yara.MatchRule) []string {

	var rules []string

	for _, match := range matches {
		rules = append(rules, match.Rule)
	}

	return rules
}
