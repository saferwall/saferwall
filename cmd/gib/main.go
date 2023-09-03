// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

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
	opts := gib.Options{Dataset: "./pkg/gib/data/ngram.json"}
	isGibberish, err := gib.NewScorer(&opts)
	if err != nil {
		log.Fatalf("NewScorer() failed with: %v", err)
	}

	// Will return `True`.
	fmt.Println(isGibberish(randomString))

	// Will return `False`.
	fmt.Println(isGibberish(nonRandomString))
}
