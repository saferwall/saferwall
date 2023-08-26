// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package gib implements a gibberish string evaluator.
package gib

import (
	"errors"
	"math"
	"strings"
)

const (
	// MinLength represents minimal length of a string to process.
	MinLength = 6
	// DefaultNgramLength is the default ngram length of the prepared dataset.
	DefaultNgramLength = 4
	// Dataset is the file path to the ngram dataset collected from a corpora.
	Dataset = "./data/ngram.json"
)

// Options provides different option to create a new scorer.
type Options struct {
	Dataset string
}

var (
	lowerCaseLetters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i",
		"j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w",
		"x", "y", "z"}
	scoreLenThreshold = 25.
	scoreLenPenalty   = 0.9233
	scoreRepPenalty   = 0.9674
	// MinScore represents the absolute minimal score for any given string
	MinScore = 8.6
)

// ngramsFromString returns a list of all n-grams of length n in a given string s.
func ngramsFromString(s string, n int) []string {

	var ngrams = make([]string, 0, len(s))

	for i := 0; i < len(s)-n+1; i++ {
		ngrams = append(ngrams, s[i:i+n])
	}

	return ngrams
}

// allNgrams returns a list of all possible n-grams.
func allNgrams(n int) []string {
	if n == 0 {
		return []string{}
	} else if n == 1 {
		return lowerCaseLetters
	}

	newNgrams := make([]string, 0)
	for _, letter := range lowerCaseLetters {
		for _, ngram := range allNgrams(n - 1) {
			newNgrams = append(newNgrams, letter+ngram)
		}
	}
	return newNgrams
}

// ngramIDFValue computes scores using modified TF-IDF.
func ngramIDFValue(totalStrings, stringFreq, totalFreq, maxFreq float64) float64 {
	return math.Log2(totalStrings / (1. + stringFreq))
}

// highestIDF computes highest idf value in map of ngram frequencies.
func highestIDF(ngramFreq NGramScores) float64 {

	max := 0.
	for _, ngram := range ngramFreq {
		max = math.Max(max, ngram.IDF())
	}
	return max
}

// highestFreq computes highest total frequency of any n-gram in a map of
// n-gram score values for a given corpus.
func highestFreq(ngramFreq NGramScores) float64 {
	max := 0.
	for _, ngram := range ngramFreq {
		max = math.Max(max, ngram.TotalFrequency())
	}
	return max
}

// nGramValues computes n-gram statistics across a given corpus of strings.
func nGramValues(corpus []string, n int, reAdjust bool) NGramScores {
	var counts = make(map[string]int, n)
	var occurrences = NewNGramSet()
	var numStrings int

	for _, s := range corpus {
		s = strings.ToLower(s)
		numStrings++
		for _, ngram := range ngramsFromString(s, n) {
			occurrences.Add(ngram, s)
			counts[ngram]++
		}
	}

	keys := allNgrams(n)
	values := make([]Score, len(keys))

	generatedNGrams := NewNGramDict(keys, values)
	maxFreq := 0
	// computes max count and assign it as max frequency of ngram
	for _, k := range counts {
		maxFreq = int(math.Max(float64(k), float64(maxFreq)))
	}

	for ngram, strings := range occurrences.Set {
		stringFreq := len(strings)
		totalFreq := counts[ngram]
		score := ngramIDFValue(float64(numStrings), float64(stringFreq),
			float64(totalFreq), float64(maxFreq))
		generatedNGrams[ngram] = [3]float64{
			float64(stringFreq),
			float64(totalFreq),
			score,
		}
	}

	if reAdjust {
		maxIDF := math.Ceil(highestIDF(generatedNGrams))
		for ngram, value := range generatedNGrams {
			if value.IDF() == 0 {
				generatedNGrams[ngram] = [3]float64{
					0.,
					0.,
					maxIDF,
				}
			}
		}
	}

	return generatedNGrams
}

// TFIDFScoreFunction generates a function that computes a score given a string.
func TFIDFScoreFunction(ngramFreq NGramScores, n int, lenThres float64,
	lenPenalty float64, repPenalty float64) func(string) float64 {

	// Formula
	// S : a string to score
	// NGramFreq : map of NGramData
	// NGramLen : the n-gram length
	// MaxFreq : max frequency of any n-gram
	// LenPenalty : pow(max,0,numNGrams 0 lenThres), lenPenalty)
	// NGramScoreSum : 0
	// for every n-gram in S:
	// 	c = count of times the n-gram appears in S
	// idf = IDF score of the n-gram from the ngramFreq map
	// tf = 0.5 + 0.5*(c/maxFreq)
	// repPenalty = pow(c,repPenalty)
	// ngramScoreSum += (tf * idf * repPenalty)
	// finalScore = (ngramScoreSum + lenPenalty) / (1 + numNGrams)

	maxFreq := highestFreq(ngramFreq)
	ngramLen := n

	score := func(s string) float64 {
		s = sanitize(s)
		ngramsInStr := ngramsFromString(s, ngramLen)
		ngramCounts := make(map[string]int)

		for _, ngram := range ngramsInStr {
			ngramCounts[ngram]++
		}
		numNGrams := len(ngramsInStr)
		lengthPenalty := math.Pow(math.Max(0., float64(numNGrams)-lenThres),
			lenPenalty)
		// compute the scores
		//	scores := make([]float64, 0)
		score := lengthPenalty
		for n, c := range ngramCounts {
			sc := ngramFreq.IDF(n) * math.Pow(float64(c), repPenalty) *
				(0.5 + 0.5*(float64(c)/maxFreq))
				//	scores = append(scores, sc)
			score += sc
		}

		return score / (1. + float64(numNGrams))
	}

	return score
}

// NewScorer creates a new scoring function
func NewScorer(opts *Options) (func(string) (bool, error), error) {

	var ngramFreq NGramScores
	var err error

	if opts != nil {
		ngramFreq, err = loadDataset(opts.Dataset)
		if err != nil {
			return nil, errors.New("failed to load dataset with error " +
				err.Error())
		}
	} else {
		ngramFreq, err = loadDataset(Dataset)
		if err != nil {
			return nil, errors.New("failed to load dataset with error " +
				err.Error())
		}
	}

	tfidfScorer := TFIDFScoreFunction(ngramFreq, DefaultNgramLength,
		scoreLenThreshold, scoreLenPenalty, scoreRepPenalty)

	scorer := func(s string) (bool, error) {
		s = sanitize(s)
		if len(s) < MinLength {
			return false,
				errors.New("string to score is too short min length is 6")
		}
		score := tfidfScorer(s)
		result := score > MinScore
		return result, nil
	}

	return scorer, nil
}
