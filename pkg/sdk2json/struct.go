package main

import (
	"strings"

	"github.com/dlclark/regexp2"
)

var (
	regStructs = `typedef [\w() ]*struct [\w]+[\n\s]+{(.|\n)+?} (?!DUMMYSTRUCTNAME|DUMMYUNIONNAME)[\w, *]+;`

	regParseStruct = `typedef [\w() ]*struct ([\w]+)[\n\s]+{((.|\n)+?)} (?!DUMMYSTRUCTNAME|DUMMYUNIONNAME)([\w, *]+);`

	regStructMember = `(?P<Type>[A-Z]+)[\s]+(?P<Name>[\w]+);`
)

// StructMember represents a member of a structure.
type StructMember struct {
	Name string
	Type string
}

// Struct represents a C data type structure.
type Struct struct {
	Name         string
	Members      []StructMember
	Alias        string
	PointerAlias string
}

func parseStruct(def string) Struct {

	winStruct := Struct{}
	r := regexp2.MustCompile(regParseStruct, 0)
	if m, _ := r.FindStringMatch(def); m != nil {

		//log.Printf("Struct definition: %v\n", m.String())
		gps := m.Groups()
		winStruct.Name = gps[1].Capture.String()

		if winStruct.Name == "_SERVICE_STATUS" {
			winStruct.Name = "_SERVICE_STATUS"
		}

		// Parse struct members
		members := strings.Split(gps[2].Capture.String(), "\n")
		for _, member := range members {
			member = standardizeSpaces(member)
			if member != "" && !strings.HasPrefix(member, "//") {
				//log.Println(member)
				m := regSubMatchToMapString(regStructMember, member)
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

	r := regexp2.MustCompile(regStructs, 0)
	matches := regexp2FindAllString(r, string(data))
	for _, m := range matches {
		structObj := parseStruct(m)
		winstructs = append(winstructs, structObj)
	}

	return matches, winstructs
}
