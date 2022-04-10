// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package magic

import (
	"strings"

	"github.com/saferwall/saferwall/internal/utils"
)

const (
	// Command to invoke the file tool
	Command = "file"
)

// Scan a file using linux `file` tool.
// This will execute file command line tool and read the stdout.
func Scan(FilePath string) (string, error) {

	args := []string{FilePath}
	output, err := utils.ExecCommand(Command, args...)
	if err != nil {
		return "", err
	}

	return ParseOutput(output), nil
}

// ParseOutput convert exiftool output into map of string|string.
func ParseOutput(fileout string) string {
	lines := strings.Split(fileout, ": ")
	if len(lines) > 1 {
		return strings.TrimSuffix(lines[1], "\n")
	}
	return ""
}
