// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/saferwall/saferwall/internal/trid"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: trid <filepath>")
		return
	}

	// get exiftool output
	r, err := trid.Scan(os.Args[1])
	check(err)

	// pretty print the results
	b, err := json.MarshalIndent(r, "", "  ")
	check(err)
	os.Stdout.Write(b)
}
