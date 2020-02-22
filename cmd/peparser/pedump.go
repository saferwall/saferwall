
// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	peparser "github.com/saferwall/saferwall/pkg/peparser"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"encoding/json"
	"bytes"
)

var (
	all        bool
	verbose        bool
	dosHeader      bool
	richHeader      bool
	ntHeader     bool
	directories    bool
	sections       bool
)


func prettyPrint(buff []byte) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, buff, "", "\t")
    if error != nil {
        log.Println("JSON parse error: ", error)
        return string(buff)
    }

    return string(prettyJSON.Bytes())
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func parsePE(filename string, cmd *cobra.Command) {
	log.Printf("Processing filename %s", filename)

	pe, err := peparser.Open(filename)
	if err != nil {
		log.Printf("Error while opening file: %s, reason: %s", filename, err)
		return
	}
	defer pe.Close()

	err = pe.Parse()
	if err != nil {
		log.Printf("Error while parsing file: %s, reason: %s", filename, err)
		return
	}

	wantDosHeader, _ := cmd.Flags().GetBool("dosheader")
	if wantDosHeader {
		dosHeader, _ := json.Marshal(pe.DosHeader)
		fmt.Println(prettyPrint(dosHeader))
	}

	wantNtHeader, _ := cmd.Flags().GetBool("ntheader")
	if wantNtHeader {
		ntHeader, _ := json.Marshal(pe.NtHeader)
		fmt.Println(prettyPrint(ntHeader))
	}

	wantAll, _ := cmd.Flags().GetBool("all")
	if wantAll {
		dosHeader, _ := json.Marshal(pe.DosHeader)
		ntHeader, _ := json.Marshal(pe.NtHeader)
		sections, _ := json.Marshal(pe.Sections)
		fmt.Println(prettyPrint(dosHeader))
		fmt.Println(prettyPrint(ntHeader))
		fmt.Println(prettyPrint(sections))
	}

}

func parse(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// filePath points to a file.
	if !isDirectory(filePath) {
		parsePE(filePath, cmd)

	} else {
	// filePath points to a directory,
	// walk recursively through all files.
		fileList := []string{}
		filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
			if !isDirectory(path) {
				fileList = append(fileList, path)
			}
			return nil
		})
	
		for _, file := range fileList {
			parsePE(file, cmd)
		}
	}
}


func main() {

	var rootCmd = &cobra.Command{
		Use:   "pedumper",
		Short: "A Portable Executable file parser",
		Long:  "A PE-Parser built for speed and malware-analysis in mind by Saferwall",
		Run:   func(cmd *cobra.Command, args []string) {
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version number",
		Long:  "Print version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("You are using version 0.0.1")
		},
	}

	var parseCmd = &cobra.Command{
		Use:   "parse",
		Short: "Parses the file",
		Long:  "Parses the Portable Executable file",
		Args:  cobra.MinimumNArgs(1),
		Run: parse,
	}

	// Init root command.
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(parseCmd)

	// Init flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	parseCmd.Flags().BoolVarP(&dosHeader, "dosheader", "", false, "Dump DOS header")
	parseCmd.Flags().BoolVarP(&richHeader, "rich", "", false, "Dump Rich header")
	parseCmd.Flags().BoolVarP(&ntHeader, "ntheader", "", false, "Dump NT header")
	parseCmd.Flags().BoolVarP(&directories, "directories", "", false, "Dump data directories")
	parseCmd.Flags().BoolVarP(&sections, "sections", "", false, "Dump section headers")
	parseCmd.Flags().BoolVarP(&all, "all", "", false, "Dump everything")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
