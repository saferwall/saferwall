// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package trid

import (
	"strings"

	"github.com/saferwall/saferwall/internal/utils"
)

const (
	// Command to invoke TriD scanner.
	tridCmd = "trid"
)

// Scan a file using TRiD Scanner
// This will execute trid command line tool and read the stdout
func Scan(FilePath string) ([]string, error) {

	args := []string{FilePath}
	output, err := utils.ExecCommand(tridCmd, args...)
	if err != nil {
		return []string{}, err
	}
	return parseOutput(output), nil

}

// parseOutput parse TriD stdout, returns an array of strings
func parseOutput(tridout string) []string {

	keepLines := []string{}
	lines := strings.Split(tridout, "\n")
	if utils.StringInSlice("Error: found no file(s) to analyze!", lines) {
		return nil
	}
	lines = lines[6:]

	for _, line := range lines {
		if len(strings.TrimSpace(line)) != 0 {
			keepLines = append(keepLines, strings.TrimSpace(line))
		}
	}

	return keepLines
}
