// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

// WriteStrSliceToFile writes a slice of string line by line to a file.
func WriteStrSliceToFile(filename string, data []string) (int, error) {
	// Open a new file for writing only
	file, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Create a new writer.
	w := bufio.NewWriter(file)
	nn := 0
	for _, s := range data {
		n, _ := w.WriteString(s + "\n")
		nn += n
	}

	w.Flush()
	return nn, nil
}

// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {

	var (
		part   []byte
		prefix bool
	)

	// Start by getting a file descriptor over the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// StringInSlice returns whether or not a string exists in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func regSubMatchToMapString(regEx, s string) (paramsMap map[string]string) {

	r := regexp.MustCompile(regEx)
	match := r.FindStringSubmatch(s)

	paramsMap = make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// Group multi-whitespaces to one whitespace.
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// Strip all whitespaces.
func spaceFieldsJoin(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// Remove C language comments.
// Removes both single line and multi-line comments.
func stripComments(s string) string {

	// Remove first the single line ones.
	regSingleLine := regexp.MustCompile(`//.*`)
	s = regSingleLine.ReplaceAllString(s, "")

	// Then the multi-lines ones.
	regMultiLine := regexp.MustCompile(`/\*([^*]|[\r\n]|(\*+([^*/]|[\r\n])))*\*+/`)
	s = regMultiLine.ReplaceAllString(s, "")
	return s
}

func findClosingBracket(text []byte, openPos int) int {
	closePos := openPos
	counter := 1
	for counter > 0 {
		closePos++
		c := text[closePos]
		if c == '{' {
			counter++
		} else if c == '}' {
			counter--
		}
	}
	return closePos
}

func findClosingSemicolon(text []byte, pos int) int {
	for text[pos] != ';' {
		pos++
	}
	return pos
}
