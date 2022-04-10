// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package gib heuristic.go implements heuristic pattern matching on strings.
package gib

import (
	"regexp"
	"strings"
)

var simplePatterns = []string{
	`\A[^eariotnslcu]+`, // Lack of any of the first 10 most-used letters in English.
	`(.){5,}`,           // Repeated single characters: 5 or more in row.
	`(.){2,}(.){2,}`,    // repeated sequences
	"abcdef",
	"bcdefg",
	"cdefgh",
	"defghi",
	"efghij",
	"fghijk",
	"ghijkl",
	"hijklm",
	"ijklmn",
	"jklmno",
	"klmnop",
	"lmnopq",
	"mnopqr",
	"nopqrs",
	"opqrst",
	"pqrstu",
	"qrstuv",
	"rstuvw",
	"stuvwx",
	"tuvwxy",
	"uvwxyz",
	"|[asdfjkl]{8}",
}

func sanitize(s string) string {
	// Make a Regex to say we only want letters and numbers.
	s = strings.ToLower(s)
	reg := regexp.MustCompile("[^a-zA-Z]+")
	processedString := reg.ReplaceAllString(s, "")
	return processedString
}

func simpleNonSense(text string) bool {
	matchers := make([]*regexp.Regexp, 0, len(simplePatterns))

	for _, pattern := range simplePatterns {
		p := regexp.MustCompile(pattern)
		matchers = append(matchers, p)
	}

	for _, matcher := range matchers {
		if matcher.MatchString(text) {
			return true
		}
	}
	return false
}
