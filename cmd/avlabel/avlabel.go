// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/saferwall/saferwall/pkg/avlabel"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: avlabel engine <detection>")
		return
	}

	engine := os.Args[1]
	detection := os.Args[2]
	var parsed map[string]string

	switch engine {
	case "windefender":
		parsed = avlabel.ParseWindefender(detection)
	case "eset":
		parsed = avlabel.ParseEset(detection)
	case "avira":
		parsed = avlabel.ParseAvira(detection)
	}

	fmt.Print(parsed)

}
