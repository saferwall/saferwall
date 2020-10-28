// Package gib implements a gibberish string evaluator.
package gib

import (
	"math"
	"strings"
)

var (
	lowerCaseLetters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
)

// ngrams returns a list of all n-grams of length n in a given string s.
func ngrams(s string, n int) []string {

	var ngrams []string = make([]string, 0, len(s))

	for i := 0; i < len(s)-n+1; i++ {
		ngrams = append(ngrams, s[i:i+n])
	}

	return ngrams
}

// allNgrams returns a list of all possible n-grams
func allNgrams(n int) []string {
	if n == 0 {
		return []string{}
	} else if n == 1 {
		return lowerCaseLetters
	}

	newNgrams := make([]string, 0, 26)
	for _, letter := range lowerCaseLetters {
		for _, ngram := range allNgrams(n - 1) {
			newNgrams = append(newNgrams, letter+ngram)
		}
	}
	return newNgrams
}

// ngramIDFValue computes scores using modified TF-IDF
func ngramIDFValue(totalStrings, stringFreq, totalFreq, maxFreq float64) float64 {
	return math.Log2(totalStrings / (1. + stringFreq))
}

// highestIDF computes highest idf value in map of ngram frequencies
func highestIDF(ngramFreq map[string]NGramData) float64 {

	max := 0.
	for _, ngram := range ngramFreq {
		max = math.Max(max, ngram.IDF)
	}
	return max
}

// highestFreq computes highest total frequency of any n-gram in a map of n-gram score
// values for a given corpus.
func highestFreq(ngramFreq map[string]NGramData) float64 {
	max := 0.
	for _, ngram := range ngramFreq {
		max = math.Max(max, ngram.TotalFrequency)
	}
	return max
}

// nGramValues computes n-gram statistics across a given corpus of strings
func nGramValues(corpus []string, n int, reAdjust bool) map[string]NGramData {
	var counts map[string]int
	var occurrences = NewNGramSet()
	var numStrings int

	for _, s := range corpus {
		s = strings.ToLower(s)
		numStrings++
		for _, ngram := range ngrams(s, n) {
			occurrences.Add(ngram, s)
			counts[ngram]++
		}
	}

	keys := allNgrams(n)
	values := make([]NGramData, len(keys))

	generatedNGrams := NewNGramDict(keys, values)
	maxFreq := 0
	// computes max count and assign it as max frequency of ngram
	for _, k := range counts {
		if k > maxFreq {
			maxFreq = k
		}
	}

	for ngram, strings := range occurrences.Set {
		stringFreq := len(strings)
		totalFreq := counts[ngram]
		score := ngramIDFValue(float64(numStrings), float64(stringFreq), float64(totalFreq), float64(maxFreq))
		generatedNGrams[ngram] = NGramData{
			StringFrequency: float64(stringFreq),
			TotalFrequency:  float64(totalFreq),
			IDF:             score,
		}
	}

	if reAdjust {
		maxIDF := math.Ceil(highestIDF(generatedNGrams))
		for ngram, value := range generatedNGrams {
			if value.IDF == 0 {
				generatedNGrams[ngram] = NGramData{
					StringFrequency: 0.,
					TotalFrequency:  0.,
					IDF:             maxIDF,
				}
			}
		}
	}

	return generatedNGrams
}
