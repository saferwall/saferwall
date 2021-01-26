// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package exiftool

import (
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
	strcase "github.com/stoewer/go-strcase"
)

const (
	// Command to invoke exiftool scanner
	Command = "exiftool"
)

// Scan a file using exiftool
// This will execute exigtool command line tool and read the stdout
func Scan(FilePath string) (map[string]string, error) {

	args := []string{FilePath}
	output, err := utils.ExecCommand(Command, args...)
	// exiftool returns exit status 1 for unknown files.
	if err != nil {
		return nil, err
	}

	return ParseOutput(output), nil
}

// ParseOutput convert exiftool output into map of string|string
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
			datas[strings.TrimSpace(strcase.UpperCamelCase(keyvalue[0]))] =
				strings.TrimSpace(keyvalue[1])
		}
	}

	return datas
}
