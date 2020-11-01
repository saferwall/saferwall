package gib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// NGramScores is a match between an ngram and it's score in from a text corpora
type NGramScores map[string][3]float64

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
