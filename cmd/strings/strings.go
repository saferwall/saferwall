// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/saferwall/saferwall/internal/utils"
	"github.com/saferwall/saferwall/pkg/gib"
	s "github.com/saferwall/saferwall/pkg/strings"
)

// Result return the results
type Result struct {
	Encoding    string
	Value       string
	IsGibberish bool
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: crypto <filepath>")
		return
	}

	data, err := utils.ReadAll(os.Args[1])
	check(err)

	// Create a new gibberish detector.
	opts := gib.Options{Dataset: "./pkg/gib/data/ngram.json"}
	isGibberish, err := gib.NewScorer(&opts)
	if err != nil {
		log.Fatalf("NewScorer() failed with: %v", err)
	}

	// Minimum string length
	n := 8

	// Get the strings
	asciiStrings := s.GetASCIIStrings(&data, n)
	wideStrings := s.GetUnicodeStrings(&data, n)
	asmStrings := s.GetAsmStrings(&data)
	fmt.Printf("Ascii: %d, Wide: %d\n", len(asciiStrings), len(wideStrings))

	// Remove duplicates
	uniqueASCII := utils.UniqueSlice(asciiStrings)
	uniqueWide := utils.UniqueSlice(wideStrings)
	uniqueAsm := utils.UniqueSlice(asmStrings)
	fmt.Printf("Unique Ascii: %d, Wide: %d\n", len(uniqueASCII), len(uniqueWide))

	var results []Result

	for _, str := range uniqueASCII {
		isGib, err := isGibberish(str)
		if err != nil {
			fmt.Println(str)
			isGib = true
		}
		results = append(results, Result{"ascii", str, isGib})
	}

	for _, str := range uniqueWide {
		isGib, _ := isGibberish(str)
		if err != nil {
			fmt.Println(str)
			isGib = true
		}
		results = append(results, Result{"wide", str, isGib})
	}

	for _, str := range uniqueAsm {
		isGib, _ := isGibberish(str)
		results = append(results, Result{"asm", str, isGib})
	}

	b, err := json.MarshalIndent(results, "", "  ")
	check(err)
	os.Stdout.Write(b)
}
