// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/saferwall/saferwall/internal/utils"
	"github.com/saferwall/saferwall/pkg/avlabel"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: avlabel engine </path/text/file>")
		return
	}

	engine := os.Args[1]
	filePath := os.Args[2]
	var parsed avlabel.Detection

	data, err := utils.ReadAll(filePath)
	if err != nil {
		log.Fatal(err)
	}
	content := string(data)
	content = strings.Replace(content, "\r\n", "\n", -1)
	detections := strings.Split(content, "\n")
	var count int
	for _, detection := range detections {
		switch engine {
		case "windefender":
			parsed = avlabel.ParseWindefender(detection)
		case "eset":
			parsed = avlabel.ParseEset(detection)
		case "avira":
			parsed = avlabel.ParseAvira(detection)
		}
		if parsed == (avlabel.Detection{}) {
			count += 1
		}
		fmt.Printf("%s => %v\n", detection, parsed)
	}

	fmt.Printf("\n Count of misses: %d", count)

}
