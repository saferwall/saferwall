// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"github.com/hillu/go-yara/v4"
	goyara "github.com/saferwall/saferwall/internal/yara"

	"errors"
	"flag"
	"log"
	"strconv"
	"strings"
)

type rules []goyara.Rule

func (r *rules) Set(arg string) error {
	if len(arg) == 0 {
		return errors.New("empty rule specification")
	}
	a := strings.SplitN(arg, ":", 2)
	switch len(a) {
	case 1:
		*r = append(*r, goyara.Rule{Filename: a[0]})
	case 2:
		*r = append(*r, goyara.Rule{Namespace: a[0], Filename: a[1]})
	}
	return nil
}

func (r *rules) String() string {
	var s string
	for _, rule := range *r {
		if len(s) > 0 {
			s += " "
		}
		if rule.Namespace != "" {
			s += rule.Namespace + ":"
		}
		s += rule.Filename
	}
	return s
}

func printMatches(m []yara.MatchRule, err error) {
	if err == nil {
		if len(m) > 0 {
			for _, match := range m {
				log.Printf("- [%s] %s ", match.Namespace, match.Rule)
			}
		} else {
			log.Print("no matches.")
		}
	} else {
		log.Printf("error: %s.", err)
	}
}

func main() {
	var (
		rules       rules
		processScan bool
		pids        []int
	)
	flag.BoolVar(&processScan, "processes", false, "scan processes instead of files")
	flag.Var(&rules, "rule", "add rule")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("no files or processes specified")
	}

	if processScan {
		for _, arg := range args {
			if pid, err := strconv.Atoi(arg); err != nil {
				log.Fatalf("Could not parse %s ad number", arg)
			} else {
				pids = append(pids, pid)
			}
		}
	}

	r, err := goyara.Load(rules)
	if err != nil {
		log.Fatalf("Could not parse %s ad number", err)
	}

	if processScan {
		s, _ := yara.NewScanner(r)
		for _, pid := range pids {
			log.Printf("Scanning process %d...", pid)
			var m yara.MatchRules
			err := s.SetCallback(&m).ScanProc(pid)
			printMatches(m, err)
		}
	} else {
		for _, filename := range args {
			log.Printf("Scanning file %s... ", filename)
			m, err := goyara.ScanFile(r, filename)
			printMatches(m, err)
		}
	}
}
