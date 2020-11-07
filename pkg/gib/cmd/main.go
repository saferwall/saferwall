package main

import (
	"fmt"

	"github.com/saferwall/saferwall/pkg/gib"
)

func main() {

	// A couple test cases.
	randomString := "asdqwfbeqbfuilac"
	nonRandomString := "CreateNewUser"

	// Create a new gibberish detector.
	// You can always provide your own
	// dataset by passing a &gib.Options{Dataset: /path/to/dataset.json}
	isGibberish := gib.NewScorer(nil)

	// Will return `True`.
	fmt.Println(isGibberish(randomString))

	// Will return `False`.
	fmt.Println(isGibberish(nonRandomString))
}
