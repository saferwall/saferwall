package gib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// NGramScores is a match between an ngram and it's score in from a text corpora
type NGramScores map[string][3]float64

// IsNGram checks if a given string is a valid ngram in our dataset
func (ns NGramScores) IsNGram(s string) bool {

	_, ok := ns[s]
	return ok
}

// Frequency returns the string computed frequency in the dataset
func (ns NGramScores) Frequency(s string) float64 {
	score, ok := ns[s]
	if !ok {
		return 0.
	}
	return score[0]
}

// TotalFrequency returns the string count in the dataset
func (ns NGramScores) TotalFrequency(s string) float64 {
	score, ok := ns[s]
	if !ok {
		return 0.
	}
	return score[1]
}

// IDF returns the IDF score in the corpus of a given ngram
func (ns NGramScores) IDF(s string) float64 {
	score, ok := ns[s]
	if !ok {
		return 0.
	}
	return score[2]
}

func loadDataset(filename string) (NGramScores, error) {

	// Open our jsonFile
	jsonFile, err := os.Open("./data/ngram.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we initialize our Users array
	// var ng NGramScore
	var ngrams NGramScores
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &ngrams)
	if err != nil {
		return nil, err
	}
	return ngrams, nil

}
