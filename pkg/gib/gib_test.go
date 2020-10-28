package gib

import (
	"bufio"
	"os"
	"testing"
)

func TestNGrams(t *testing.T) {

	testCases := []struct {
		input  string
		output []string
		n      int
	}{
		{
			"hello world",
			[]string{"hel", "ell", "llo", "lo ", "o w", " wo", "wor", "orl", "rld"},
			3,
		}, {
			"hello world",
			[]string{"hell", "ello", "llo ", "lo w", "o wo", " wor", "worl", "orld"},
			4,
		}, {
			"saferwall",
			[]string{"safer", "aferw", "ferwa", "erwal", "rwall"},
			5,
		}, {
			"",
			[]string{},
			5,
		},
	}

	for _, tt := range testCases {
		res := ngrams(tt.input, tt.n)

		for i, s := range tt.output {
			if s != res[i] {
				t.Errorf("bad ngram exepceted %s got %s", s, res[i])
			}
		}
	}
}

func TestAllNGrams(t *testing.T) {
	expectedThreeNgrams, err := readLines("./testdata/all-three-ngrams.txt")
	if err != nil {
		t.Fatal("failed to read file with ", err)
	}
	allPossibleThreeNgrams := allNgrams(3)
	for i, c := range expectedThreeNgrams {
		if c != allPossibleThreeNgrams[i] {
			t.Fatalf("wrong possible ngram expected %s got %s", c, allPossibleThreeNgrams[i])
		}
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
