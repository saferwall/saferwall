package main

import (
	"log"
	"regexp"
	"strings"
)

var (
	regStructs = `typedef[\w() ]*?struct[\w\s]*{`

	//
	// Case 1:
	// typedef struct _INTERNET_BUFFERSA {
	// 	DWORD dwStructSize;
	// 	DWORD dwOffsetLow;
	// } INTERNET_BUFFERSA, * LPINTERNET_BUFFERSA;
	//

	//
	// Case2 :
	// typedef struct {
	// 	BOOL    fAccepted;
	// 	BOOL    fLeashed;
	// }
	// InternetCookieHistory;
	//

	regParseStruct = `typedef ([\w() ]*?)struct( [\w]+)*((.|\n)+?})([\w\n ,*_]+?;)(\n} [\w,* ]+;)*`

	// DWORD cbOutQue;
	// DWORD fCtsHold : 1;
	// WCHAR wcProvChar[1];
	// ULONG iHash; // index of hash object
	// _Field_size_(cbBuffer) PUCHAR pbBuffer;
	// SIZE_T dwAvailVirtual;
	// PWSTR *rgpszFunctions;
	regStructMember = `(?P<Type>[\w]+[\s*]+)(?P<Name>[\w]+)(?P<ArraySize>\[\w+\])*(?P<BitPack>[ :\d]+)*`
)

// StructMember represents a member of a structure.
type StructMember struct {
	// The name of the structure member.
	// When the member itself represents a structure/union, we use the name of the structure/union otherwise `anonymous`.
	Name string `json:"name"`
	// The type of the member: DWORD, int, char*, ...
	// Or `_structure` / `_union` for complexe types.
	Type string `json:"type"`
	// For complex types, `Body`describes the struct/union members.
	Body []StructMember `json:"body,omitempty"`
}

// Struct represents a C data type structure.
type Struct struct {
	Name         string `json:"name"`
	Members      []StructMember
	PointerAlias string `json:"pointer_alias"`
}

// Delete all white spaces from a C structure.
func stripStruct(s string) string {
	s = stripComments(s)
	s = standardizeSpaces(s)
	s = strings.ReplaceAll(s, "; ", ";")
	s = strings.ReplaceAll(s, " { ", "{")
	s = strings.ReplaceAll(s, " } ", "}")
	s = strings.ReplaceAll(s, " : ", ":")
	s = strings.ReplaceAll(s, ": ", ":")
	return s
}

func parseStructBody(body string) []StructMember {

	var structMembers []StructMember

	pos := 0
	log.Println(body)
	endPos := len(body) - 1
	for pos < endPos {
		sm := StructMember{}
		semiColPos := strings.Index(body[pos:], ";")
		if semiColPos < 0 {
			break
		}
		memberStr := body[pos : pos+semiColPos]
		mu := strings.Index(memberStr, "union{")
		ms := strings.Index(memberStr, "struct{")
		if mu < 0 && ms < 0 {
			mMap := regSubMatchToMapString(regStructMember, memberStr)
			sm.Type = spaceFieldsJoin(mMap["Type"])
			sm.Name = mMap["Name"]
			pos += semiColPos + 1 // for the ;
		} else {
			l := 0
			// Union inside the struct OR Union comes first then struct.
			if (mu >= 0 && ms < 0) || (mu >= 0 && ms >= 0 && mu < ms) {
				sm.Type = "_union"
				l = len("union{") + mu

			} else if (ms >= 0 && mu < 0) || (mu >= 0 && ms >= 0 && mu < ms) {
				// Struct inside the struc OR Struct comes first then union.
				sm.Type = "_struct"
				l = len("struct{") + ms
			}

			endStructPos := findClosingBracket([]byte(body), pos+l+1)+1
			semiColPos = findClosingSemicolon([]byte(body), endStructPos)
			structBody := body[pos+l : endStructPos-1]
			sm.Name = spaceFieldsJoin(body[endStructPos:semiColPos])
			sm.Body = parseStructBody(structBody)
			pos = semiColPos + 1 // for the ;
		}

		structMembers = append(structMembers, sm)
	}

	return structMembers
}

func parseStruct(structBeg, structBody, structEnd string) Struct {

	winStruct := Struct{}

	// Start by deleteing unecessery characters like comments and whitespaces.
	structBody = stripStruct(structBody)

	if strings.Contains(structBody, "CachingFlags") {
		log.Println("w9raf 3and 7addak")
	}

	// Get struct members
	winStruct.Members = parseStructBody(structBody)

	// Get struct name
	structEnd = spaceFieldsJoin(structEnd)
	n := strings.Split(structEnd, ",")
	if len(n) > 0 {
		winStruct.Name = n[0]
	}

	if len(n) == 2 && strings.HasPrefix(n[1], "*") {
		winStruct.PointerAlias = n[1][1:]
	}

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
		structObj := parseStruct(structBeg, structBody, structEnd)
		winstructs = append(winstructs, structObj)
		strStructs = append(strStructs, strStruct)
	}
	return strStructs, winstructs
}
