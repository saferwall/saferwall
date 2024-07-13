// Copyright 2018 Saferwall. All rights reserved.
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

// Shorten returns a short version of the full magic output.
func Shorten(magicResult string) string {
	var typeMap = map[string]string{
		// binary executables
		"PE32":                "pe",
		"ELF":                 "elf",
		"Mach-O":              "mach-o",
		"MS Windows shortcut": "lnk",
		"MS-DOS":              "msdos",

		// documents
		"PDF document":                        "pdf",
		"Rich Text Format":                    "rtf",
		"Microsoft Word 2007+":                "ooxml",
		"Microsoft Excel 2007+":               "ooxml",
		"Microsoft PowerPoint 2007+":          "ooxml",
		"Composite Document File V2 Document": "ole2",

		// images and media
		"PC bitmap":           "bmp",
		"JPEG image data":     "jpeg",
		"PNG image data":      "png",
		"GIF image data":      "gif",
		"SVG Scalable Vector": "svg",
		"Macromedia Flash":    "swf",

		// archives
		"Zip archive data":          "zip",
		"RAR archive data":          "rar",
		"7-zip archive data":        "7-zip",
		"gzip compressed data":      "gzip",
		"bzip2 compressed data":     "bzip2",
		"tar archive":               "tar",
		"XZ compressed data":        "xz",
		"Java archive data (JAR)":   "jar",
		"Microsoft Cabinet archive": "cab",

		// misc
		"ISO 9660 CD-ROM": "iso",

		// text-based: xml, html, js, hta, swf, ....
		"ASCII text":              "txt",
		"Unicode text":            "txt",
		"ISO-8859 text":           "txt",
		"Unicode (with BOM) text": "txt",
	}

	for k, v := range typeMap {
		if strings.Contains(magicResult, k) {
			return v
		}
	}

	// `file` returns = "data" for any binary data.
	if magicResult == "data" {
		return magicResult
	}

	return "unknown"
}

// Scan a file using linux `file` tool.
// This will execute file command line tool and read the stdout.
func Scan(FilePath string) (string, error) {

	args := []string{FilePath}
	output, err := utils.ExecCmd(Command, args...)
	if err != nil {
		return "", err
	}

	return ParseOutput(output), nil
}

// ScanBytes a memory buffer using linux `file` tool.
func ScanBytes(data []byte) (string, error) {

	filePath, err := utils.CreateTempFile(data)
	if err != nil {
		return "", err
	}

	args := []string{filePath}
	output, err := utils.ExecCmd(Command, args...)
	if err != nil {
		return "", err
	}

	res := ParseOutput(output)
	defer utils.DeleteFile(filePath)
	return res, nil

}

// ParseOutput convert exiftool output into map of string|string.
func ParseOutput(fileout string) string {
	lines := strings.Split(fileout, ": ")
	if len(lines) > 1 {
		return strings.TrimSuffix(lines[1], "\n")
	}
	return ""
}
