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

	"github.com/saferwall/saferwall/pkg/crypto"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
		return
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: crypto <filepath>")
		return
	}

	// read the file
	data, err := ioutil.ReadFile(os.Args[1])
	check(err)

	// compute the hashes
	r := crypto.HashBytes(data)

	// pretty print the results
	b, err := json.MarshalIndent(r, "", "  ")
	check(err)
	os.Stdout.Write(b)
}
