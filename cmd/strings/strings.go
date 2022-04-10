// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/saferwall/saferwall/internal/utils"
	s "github.com/saferwall/saferwall/pkg/strings"
)

// Result return the results
type Result struct {
	Encoding string
	Value    string
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

	// Read the file
	data, err := ioutil.ReadFile(os.Args[1])
	check(err)

	// Mininum string length
	n := 6

	// Get the strings
	asciiStrings := s.GetASCIIStrings(&data, n)
	wideStrings := s.GetUnicodeStrings(&data, n)
	asmStrings := s.GetAsmStrings(&data)
	fmt.Printf("Ascii: %d, Wide: %d\n", len(asciiStrings), len(wideStrings))

	// Remove duplicates
	uniqueASCII := utils.UniqueSlice(asciiStrings)
	uniqueWide := utils.UniqueSlice(wideStrings)
	uniqueAsm := utils.UniqueSlice(asmStrings)
	fmt.Printf("Unique Ascii: %d, Wide: %d", len(uniqueASCII), len(uniqueWide))

	var results []Result

	for _, str := range uniqueASCII {
		results = append(results, Result{"ascii", str})
	}

	for _, str := range uniqueWide {
		results = append(results, Result{"wide", str})
	}

	for _, str := range uniqueAsm {
		results = append(results, Result{"asm", str})
	}

	b, err := json.MarshalIndent(results, "", "  ")
	check(err)
	os.Stdout.Write(b)
}
