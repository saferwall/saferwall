package main

import (
	"fmt"
	"log"

	"github.com/saferwall/saferwall/pkg/gib"
)

func main() {

	// A couple test cases.
	randomString := "asdqwfbeqbfuilac"
	nonRandomString := "CreateNewUser"

	// Create a new gibberish detector.
	// You can always provide your own
	// dataset by passing a &gib.Options{Dataset: /path/to/dataset.json}
	isGibberish, err := gib.NewScorer(nil)
	if err != nil {
		log.Fatalf("NewScorer() failed with: %v", err)
	}

	// Will return `True`.
	fmt.Println(isGibberish(randomString))

	// Will return `False`.
	fmt.Println(isGibberish(nonRandomString))
}
