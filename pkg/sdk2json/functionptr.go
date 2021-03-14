package main

import (
	"log"
	"regexp"
)

var (
	regFunctionPtr = `typedef[\w\s]+\(WINAPI \*(?P<Name>\w+)\)\(`
)

func parseFunctionPointers(data string) []string {

	var funcPtrs []string
	r := regexp.MustCompile(regFunctionPtr)
	matches := r.FindAllStringSubmatch(data, -1)
	for _, m := range matches {
		if len(m) > 0 {
			log.Println(m[1])
			funcPtrs = append(funcPtrs, m[1])
		}
	}

	return funcPtrs
}
