// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package packer

import (
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	// Command to invoke exiftool scanner
	Command = "/opt/die/diec.sh"
)

// Scan a file using exiftool
// This will execute exigtool command line tool and read the stdout
func Scan(FilePath string) ([]string, error) {

	args := []string{FilePath}
	output, err := utils.ExecCommand(Command, args...)
	if err != nil {
		return nil, err
	}

	return parseOutput(output), nil
}

// parseOutput parse TriD stdout, returns an array of strings
func parseOutput(tridout string) []string {
	keepLines := []string{}
	lines := strings.Split(tridout, "\n")
	for _, line := range lines {
		if line != "" {
			keepLines = append(keepLines, line)
		}

	}
	return keepLines
}
