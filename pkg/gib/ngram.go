// Package gib implements gibberish strings evaluation.
// file ngram.go implements N-Grams data structures and utility functions.
package gib

// NGramData defines an entry in an n-gram table that contains frequency statistics
// and IDF (inverse document frequency) derived from example data.
type NGramData struct {
	StringFrequency float64
	TotalFrequency  float64
	IDF             float64
}
