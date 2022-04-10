// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package ml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var testFixturesPath = "../../testdata/"

var peClassTests = []struct {
	in  string
	out ClassifierPrediction
}{
	{
		path.Join(testFixturesPath, "pe-class-test1.json"), ClassifierPrediction{
			Class:       "Label.MALICIOUS",
			Probability: 0.978692142794345,
			Score:       "Malicious (High Trust)",
			SHA256:      "4c728576bd65c8e8348410d1ab3bb5d6cae093985d9e82d3121295b16429b2db"},
	},
}

var stringRankerTests = []struct {
	in  string
	out StringsRanker
}{
	{path.Join(testFixturesPath, "string-ranker-test1.json"), StringsRanker{
		Strings: []string{"GetProcAddress", "LoadLibraryA", "GetProcessHeap"},
		SHA256:  "4c728576bd65c8e8348410d1ab3bb5d6cae093985d9e82d3121295b16429b2db"},
	},
}

// readAll reads the entire file into memory.
func readAll(filePath string) ([]byte, error) {
	// Start by getting a file descriptor over the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
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
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(ClassifierPrediction{
			Class:       "Label.MALICIOUS",
			Probability: 0.978692142794345,
			Score:       "Malicious (High Trust)",
			SHA256:      "4c728576bd65c8e8348410d1ab3bb5d6cae093985d9e82d3121295b16429b2db",
		})
		fmt.Fprintln(w, string(response))
	}))
	defer s.Close()
	for _, tt := range peClassTests {
		buff, err := readAll(tt.in)
		if err != nil {
			t.Fatalf("failed to read the file (%s): %v", tt.in, err)
		}
		got, err := PEClassPrediction(s.URL, buff)
		if err != nil {
			t.Fatalf("PEClassPrediction(%s) failed with error %s", tt.in, err)
		}
		if got != tt.out {
			t.Fatalf("PEClassPrediction(%s) failed expected %v got %v", tt.in, tt.out, got)
		}

	}
}

func TestStringRanker(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(StringsRanker{
			Strings: []string{"GetProcAddress", "LoadLibraryA", "GetProcessHeap"},
			SHA256:  "4c728576bd65c8e8348410d1ab3bb5d6cae093985d9e82d3121295b16429b2db",
		})
		fmt.Fprintln(w, string(response))
	}))
	defer s.Close()
	for _, tt := range stringRankerTests {
		buff, err := readAll(tt.in)
		if err != nil {
			t.Fatalf("failed to read the file (%s): %v", tt.in, err)
		}
		got, err := RankStrings(s.URL, buff)
		if err != nil {
			t.Fatalf("RankStrings(%s) failed with error %s", tt.in, err.Error())
		}
		if got.SHA256 != tt.out.SHA256 {
			t.Fatalf("RankStrings(%s) failed with expected SHA256 : %s got SHA256 %s", tt.in, tt.out.SHA256, got.SHA256)
		}
		for i, v := range tt.out.Strings {
			if v != got.Strings[i] {
				t.Fatalf("RankStrings(%s) failed at index %d expected %s got %s", tt.in, i, tt.out.Strings[i], got.Strings[i])
			}
		}

	}
}
