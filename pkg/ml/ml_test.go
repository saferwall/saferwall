// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package ml

import (
	"os"
	"testing"
)

const (
	server = "http://localhost:8000"
)

var peClassTests = []struct {
	in  string
	out ClassifierPrediction
}{
	{"pe-class-test1.json", ClassifierPrediction{
		Class:       "Label.MALICIOUS",
		Probability: 0.978692142794345,
		Score:       "Malicious (High Trust)",
		Sha256:      "4c728576bd65c8e8348410d1ab3bb5d6cae093985d9e82d3121295b16429b2db"},
	},
}


var stringRankerTests = []struct {
	in  string
	out StringsRanker
}{
	{"string-ranker-test1.json", StringsRanker{
		Strings: []string{"GetProcAddress", "LoadLibraryA", "GetProcessHeap",},
		Sha256:  "4c728576bd65c8e8348410d1ab3bb5d6cae093985d9e82d3121295b16429b2db"},
	},
}

// readAll reads the entire file into memory.
func readAll(filePath string) ([]byte, error) {
	// Start by getting a file descriptor over the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file size to know how much we need to allocate
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	// Read the whole binary
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
func TestPEClass(t *testing.T) {
	for _, tt := range peClassTests {
		t.Run(tt.in, func(t *testing.T) {
			buff, err := readAll(tt.in)
			if err != nil {
				t.Errorf("failed to read the file (%s): %v", tt.in, err)
			}
			got, err := PEClassPrediction(server, buff)
			if err != nil {
				t.Errorf("PEClassPrediction(%s) got %v, want %v",
					tt.in, err, tt.in)
			}
			if got != tt.out {
				t.Errorf("PEClassPrediction(%s) got %v, want %v",
					tt.in, got, tt.out)
			}
		})
	}
}

func TestStringRanker(t *testing.T) {
	for _, tt := range stringRankerTests {
		t.Run(tt.in, func(t *testing.T) {
			buff, err := readAll(tt.in)
			if err != nil {
				t.Errorf("failed to read the file (%s): %v", tt.in, err)
			}
			got, err := RankStrings(server, buff)
			if err != nil {
				t.Errorf("RankStrings(%s) got %v, want %v",
					tt.in, err, tt.in)
			}
			if got.Strings[1] != tt.out.Strings[0] {
				t.Errorf("RankStrings(%s) got %v, want %v",
					tt.in, got, tt.out)
			}
		})
	}
}