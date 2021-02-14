package main

import (
	"regexp"
	"strings"
)

var (
	//regStructs = `typedef[\w() ]*?struct`
	regStructs = `typedef[\w() ]*?struct[\w\s]*{`

	regParseStruct = `typedef ([\w() ]*?)struct( [\w]+)*((.|\n)+?})([\w\n ,*_]+?;)(\n} [\w,* ]+;)*`

	// DWORD cbOutQue;
	// DWORD fCtsHold : 1;
	// WCHAR wcProvChar[1];
	// ULONG iHash; // index of hash object
	// _Field_size_(cbBuffer) PUCHAR pbBuffer;
	// SIZE_T dwAvailVirtual;
	regStructMember = `(?P<Type>[A-Za-z_]+)[\s]+(?P<Name>[\w]+)(?P<ArraySize>\[\w+\])*(?P<BitPack>[ :\d]+)*;`
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
	PointerAlias string
}

func findClosingBracket(text []byte, openPos int) int {
	closePos := openPos
	counter := 1
	for counter > 0 {
		closePos++
		c := text[closePos]
		if c == '{' {
			counter++
		} else if c == '}' {
			counter--
		}
	}
	return closePos
}

func findClosingSemicolon(text []byte, pos int) int {
	for text[pos] != ';' {
		pos++
	}
	return pos
}

func parseStruct(structBeg, structBody, structEnd string) Struct {

	winStruct := Struct{}

	// Get struct members
	r := regexp.MustCompile(regStructMember)
	matches := r.FindAllStringSubmatch(structBody, -1)
	paramsMap := make(map[string]string)
	for _, match := range matches {
		for i, name := range r.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}

		}
		sm := StructMember{
			Type: paramsMap["Type"],
			Name: paramsMap["Name"],
		}
		winStruct.Members = append(winStruct.Members, sm)
	}

	// Get struct name
	structEnd = spaceFieldsJoin(structEnd)
	n := strings.Split(structEnd, ",")
	if len(n) > 0 {
		winStruct.Name = n[0]
	}

	if len(n) == 2 && strings.HasPrefix(n[1], "*") {
		winStruct.PointerAlias = n[1][1:]
	}

	// Case 1:
	// typedef struct _INTERNET_BUFFERSA {
	// 	DWORD dwStructSize;
	// 	DWORD dwOffsetLow;
	// } INTERNET_BUFFERSA, * LPINTERNET_BUFFERSA;

	// Case2 :
	// typedef struct {
	// 	BOOL    fAccepted;
	// 	BOOL    fLeashed;
	// }
	// InternetCookieHistory;

	return winStruct
}

func getAllStructs(data []byte) ([]string, []Struct) {

	var winstructs []Struct
	var strStructs []string

	re := regexp.MustCompile(regStructs)
	matches := re.FindAllStringIndex(string(data), -1)
	for _, m := range matches {

		endPos := findClosingBracket(data, m[1])
		endStruct := findClosingSemicolon(data, endPos+1)

		structBeg := string(data[m[0]:m[1]])
		structBody := string(data[m[1]:endPos])
		structEnd := string(data[endPos+1 : endStruct])
		strStruct := string(data[m[0] : endStruct+1])
		// log.Println(structBeg)
		// log.Println(structBody)
		// log.Println(structEnd)
		// log.Println(strStruct)
		structObj := parseStruct(structBeg, structBody, structEnd)
		winstructs = append(winstructs, structObj)
		strStructs = append(strStructs, strStruct)
	}
	return strStructs, winstructs
}
