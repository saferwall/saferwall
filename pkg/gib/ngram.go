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

// NGramSet defines a set of ngrams and their respective strings.
// Each Ngram is a key and the values are represented by a slice
// of strings pertaining to that ngram.
type NGramSet struct {
	Set map[string][]string
}

// NewNGramData creates a new instance of an ngram data tuple
func NewNGramData() NGramData {
	return NGramData{}
}

// NewNGramSet creates a new instance of ngram set
func NewNGramSet() NGramSet {
	set := make(map[string][]string, 0)

	return NGramSet{
		Set: set,
	}
}

// Add a new string to an ngram set
func (n *NGramSet) Add(ngram string, s string) {
	n.Set[ngram] = append(n.Set[ngram], s)
}

// NGramDict is a dictionary (map) of ngrams and their statistics
type NGramDict map[string]NGramData

// NewNGramDict creates a new instance of ngram dict
func NewNGramDict(keys []string, values []NGramData) NGramDict {

	// if len(keys) != len(values) throw an error
	dict := make(map[string]NGramData, len(keys))
	for i, k := range keys {
		dict[k] = values[i]
	}

	return dict
}
