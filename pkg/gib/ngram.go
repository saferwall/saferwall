// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package gib implements gibberish strings evaluation.
// file ngram.go implements N-Grams data structures and utility functions.
package gib

// NGramSet defines a set of ngrams and their respective strings.
// Each Ngram is a key and the values are represented by a slice
// of strings pertaining to that ngram.
type NGramSet struct {
	Set map[string][]string
}

// NewNGramSet creates a new instance of ngram set.
func NewNGramSet() NGramSet {
	set := make(map[string][]string)
	return NGramSet{
		Set: set,
	}
}

// Add a new string to an ngram set.
func (n *NGramSet) Add(ngram string, s string) {
	n.Set[ngram] = append(n.Set[ngram], s)
}

// NewNGramDict creates a new instance of ngram dict.
func NewNGramDict(keys []string, values []Score) NGramScores {

	// if len(keys) != len(values) throw an error
	var dict = make(NGramScores)
	for i, k := range keys {
		dict[k] = values[i]
	}

	return dict
}
