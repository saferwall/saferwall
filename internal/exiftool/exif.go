// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package exiftool

import (
	"strings"
	"unicode"

	"github.com/saferwall/saferwall/internal/utils"
)

const (
	// Command to invoke exiftool scanner
	cmd = "exiftool"
)

// Scan a file using exiftool
// This will execute exigtool command line tool and read the stdout
func Scan(FilePath string) (map[string]string, error) {

	args := []string{FilePath}
	output, err := utils.ExecCommand(cmd, args...)
	// exiftool returns exit status 1 for unknown files.
	if err != nil {
		return nil, err
	}

	return ParseOutput(output), nil
}

// ParseOutput convert exiftool output into map of string|string.
func ParseOutput(exifout string) map[string]string {

	var ignoreTags = []string{
		"Directory",
		"File Name",
		"File Permissions",
	}

	lines := strings.Split(exifout, "\n")
	if utils.StringInSlice("File not found", lines) {
		return nil
	}

	datas := make(map[string]string, len(lines))
	for _, line := range lines {
		keyvalue := strings.Split(line, ":")
		if len(keyvalue) != 2 {
			continue
		}
		if !utils.StringInSlice(strings.TrimSpace(keyvalue[0]), ignoreTags) {
			datas[strings.TrimSpace(camelCase(keyvalue[0]))] =
				strings.TrimSpace(keyvalue[1])
		}
	}

	return datas
}

// camelCase convert a string to camelcase
func camelCase(s string) string {
	s = strings.TrimSpace(s)
	buffer := make([]rune, 0, len(s))
	stringIter(s, func(prev, curr, next rune) {
		if !isDelimiter(curr) {
			if isDelimiter(prev) || (prev == 0) {
				buffer = append(buffer, unicode.ToUpper(curr))
			} else if unicode.IsLower(prev) {
				buffer = append(buffer, curr)
			} else {
				buffer = append(buffer, unicode.ToLower(curr))
			}
		}
	})

	return string(buffer)
}

// isDelimiter checks if a character is some kind of whitespace or '_' or '-'.
func isDelimiter(ch rune) bool {
	return ch == '-' || ch == '_' || unicode.IsSpace(ch)
}

// stringIter iterates over a string, invoking the callback for every single rune in the string.
func stringIter(s string, callback func(prev, curr, next rune)) {
	var prev rune
	var curr rune
	for _, next := range s {
		if curr == 0 {
			prev = curr
			curr = next
			continue
		}

		callback(prev, curr, next)

		prev = curr
		curr = next
	}

	if len(s) > 0 {
		callback(prev, curr, 0)
	}
}
