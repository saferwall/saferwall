package gib

import (
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
