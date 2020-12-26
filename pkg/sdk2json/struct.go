package main

import (
	"fmt"
	"strings"


	"github.com/dlclark/regexp2"
)

var (
	RegStructs = `typedef [\w() ]*struct [\w]+[\n\s]+{(.|\n)+?} (?!DUMMYSTRUCTNAME|DUMMYUNIONNAME)[\w, *]+;`

	RegParseStruct = `typedef [\w() ]*struct ([\w]+[\n\s]+){((.|\n)+?)} (?!DUMMYSTRUCTNAME|DUMMYUNIONNAME)([\w, *]+);`

	RegStructMember = `(?P<Type>[A-Z]+)[\s]+(?P<Name>[a-zA-Z]+); `
)

type StructMember struct {
	Name string
	Type string
}

type Struct struct {
	Name         string
	Members      []StructMember
	Alias        string
	PointerAlias string
}

func parseStruct(def string) Struct {

	winStruct := Struct{}
	r := regexp2.MustCompile(RegParseStruct, 0)
	if m, _ := r.FindStringMatch(def); m != nil {

		fmt.Printf("Struct definition: %v\n", m.String())
		gps := m.Groups()
		winStruct.Name = gps[1].Capture.String()

		// Parse struct members
		members := strings.Split(gps[2].Capture.String(), "\n")
		for _, member := range members {
			member = standardizeSpaces(member)
			if member != "" && !strings.HasPrefix(member, "//") {
				fmt.Println(member)
				m := regSubMatchToMapString(RegStructMember, member)
				sm := StructMember{
					Type: m["Type"],
					Name: m["Name"],
				}
				winStruct.Members = append(winStruct.Members, sm)
			}
		}
		winStruct.Alias = gps[4].Capture.String()
	}

	return winStruct
}

func getAllStructs(data string) ([]string, []Struct) {

	var winstructs []Struct

	r := regexp2.MustCompile(RegStructs, 0)
	matches := regexp2FindAllString(r, string(data))
	for _, m := range matches {
		structObj := parseStruct(m)
		winstructs = append(winstructs, structObj)
	}

	return matches, winstructs
}
