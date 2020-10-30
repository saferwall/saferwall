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
	`^M{0,4}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$`,
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
	// Make a Regex to say we only want letters and numbers
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	processedString := reg.ReplaceAllString(s, "")
	return strings.ToLower(processedString)
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
