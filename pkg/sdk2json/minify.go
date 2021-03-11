// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"strings"
)

const (
	paramIn uint8 = iota
	paramOut
	paramInOut
	paramReserved
)

// StructUnionMemberMini represents a struct or a union member.
type StructUnionMemberMini struct {
	Name       string          `json:"name"`
	X86Offset  uint8           `json:"x86off"`
	X86Size    uint8           `json:"x86size"`
	X64Offset  uint8           `json:"x64off"`
	X64Size    uint8           `json:"x64size"`
	Type       uint8           `json:"type"` // Help interpret the value.
	Definition StructUnionMini `json:"def,omitempty"`
}

// StructUnionMini represents a complex type like a Union or a Struct.
type StructUnionMini struct {
	Name    string                  `json:"name"`
	Members []StructUnionMemberMini `json:"members"`
	X86Size uint8                   `json:"x86size"`
	X64Size uint8                   `json:"x64size"`
}

// APIParamMini represents a paramter of a Win32 API.
type APIParamMini struct {
	Annotation        uint8  `json:"anno"`
	Type              uint8  `json:"type"`
	Name              string `json:"name"`
	BufferSizeOrIndex uint8   `json:"buffsize_or_idx"`
}

// APIMini represents information about a Win32 API.
type APIMini struct {
	ReturnValueType bool           `json:"retType"` // Return value type.
	PropertyCount   int            `json:"p_count"`
	Params          []APIParamMini `json:"params"` // API Arguments.
}

var (
	reAnnotationIn = regexp.MustCompile(`(?i)_In_|IN|_In_opt[\w]+|In_reads[\w()]+`)
	// __out_data_source
	reAnnotationOut      = regexp.MustCompile(`(?i)_Out_|OUT|_Out_opt[\w]+|_Out_writes[\w()]+|_Outptr_`)
	reAnnotationIntOut   = regexp.MustCompile(`(?i)_Inout[\w]+`)
	reAnnotationReserved = regexp.MustCompile(`(?i)Reserved`)

	reOutWritesBytesTo    = `\w+\((?P<s>[*\w]+), (?P<c>[*\w]+)\)`
	reInOutReadWriteBytes = `\w+\((?P<s>\w+)\)`
)

func findParamIndexByName(api API, target string) int {
	for i, param := range api.Params {
		if param.Name == target {
			return i
		}
	}
	return 0
}

func getNameFromAnnotation(param APIParam) string {
	m := regSubMatchToMapString(reOutWritesBytesTo, param.Annotation)
	if len(m) > 0 {
		return m["c"]
	}

	m = regSubMatchToMapString(reInOutReadWriteBytes, param.Annotation)
	if len(m) > 0 {
		return m["s"]
	}

	return ""
}

func getBytePtrIndex(api API, param APIParam, dt dataType,
	parammini *APIParamMini) uint8 {
	if dt.Kind == typeBytePtr {
		// log.Printf("API: %s, Name: %s, Type: %s, Anno: %s\n", api.Name,
		// 	param.Name, param.Type, param.Annotation)

		name := getNameFromAnnotation(param)
		if strings.HasPrefix(name, "*") {
			name = name[1:]
		}
		// log.Println(name)
		idx := findParamIndexByName(api, name)
		return uint8(idx)

	} else if dt.Name == "LPVOID" {
		// Unfortunately MS is not really consistent about data types.
		// LPVOID sometimes points to an DWORD_PTR, and sometimes it points to
		// []bytes, we try to make a guess based on the annotation.
		name := getNameFromAnnotation(param)
		if len(name) > 0 {
			// log.Printf("API: %s, Name: %s, Type: %s, Anno: %s\n", api.Name,
			// 	param.Name, param.Type, param.Annotation)
			if strings.HasPrefix(name, "*") {
				name = name[1:]
			}
			idx := findParamIndexByName(api, name)
			parammini.Type = typeBytePtr
			return uint8(idx)
		}

	}

	return dt.Size

}

func minifyAPIs(apis map[string]map[string]API) map[string]map[string]APIMini {
	mapis := make(map[string]map[string]APIMini)
	for dllname, v := range apis {
		if _, ok := mapis[dllname]; !ok {
			mapis[dllname] = make(map[string]APIMini)
		}
		for apiname, vv := range v {

			// Return type.
			returnType := false
			if vv.ReturnValueType == "VOID" {
				returnType = true
			}

			propertiesCount := 0
			var paramsMini []APIParamMini
			for _, param := range vv.Params {
				parammini := APIParamMini{}
				if reAnnotationIn.MatchString(param.Annotation) {
					parammini.Annotation = paramIn
					propertiesCount++
				} else if reAnnotationOut.MatchString(param.Annotation) {
					parammini.Annotation = paramOut
					propertiesCount++
				} else if reAnnotationIntOut.MatchString(param.Annotation) {
					parammini.Annotation = paramInOut
					propertiesCount += 2
				} else if reAnnotationReserved.MatchString(param.Annotation) {
					parammini.Annotation = paramReserved
					propertiesCount++
				} else {
					// If we don't know, take it as in:
					parammini.Annotation = paramIn
				}

				// Get the param type.
				dataType := typefromString(param.Type)
				parammini.Type = dataType.Kind
				parammini.BufferSizeOrIndex = getBytePtrIndex(vv, param, dataType, &parammini)
				parammini.Name = param.Name
				paramsMini = append(paramsMini, parammini)
			}
			apiMini := APIMini{}
			apiMini.Params = paramsMini
			apiMini.ReturnValueType = returnType
			apiMini.PropertyCount = propertiesCount
			mapis[dllname][apiname] = apiMini
		}
	}

	return mapis
}

func minifyStructAndUnions(winStructs []Struct) []StructUnionMemberMini {
	structsAndUnionsMini := []StructUnionMemberMini{}
	for _, winStruct := range winStructs {
		structUnionMini := StructUnionMini{}
		structUnionMini.Name = winStruct.Name

		//
		// Because of structure padding, the code below calculates the index of
		// where each individual member is located as well as the size of the
		// structure in both x86 and x64.
		// Here are a few rules the Microsoft C compiler arranges structures in
		// memory:
		//
		// 1. each individual member, there will be padding so that to make it
		// start at an address that is divisible by its size.
		// e.g on 64 bit system,int should start at address divisible by 4, and
		// long by 8, short by 2.
		// 2. char and char[] are special, could be any memory address, so they
		// don't need padding before them.
		// 3. For struct, other than the alignment need for each individual
		// member, the size of whole struct itself will be aligned to a size
		// divisible by size of largest individual member, by padding at end.
		// e.g if struct's largest member is long then divisible by 8, int then
		// by 4, short then by 2.
		//

		x86Size := uint8(0)
		x64Size := uint8(0)
		x86Offset := uint8(0)
		x64Offset := uint8(0)
		largestMemSizex86, largestMemSizex64 := winStruct.Max()
		for i, winStructMember := range winStruct.Members {
			miniMember := StructUnionMemberMini{}
			miniMember.Name = winStructMember.Name
			miniMember.X86Offset = x86Offset
			miniMember.X64Offset = x64Offset
			miniMember.X86Size, miniMember.X64Size = winStructMember.Size()
			countBytesPaddingx86, countBytesPaddingx64 := winStruct.structMemberPadding(i, largestMemSizex86, largestMemSizex64)
			x86Offset += countBytesPaddingx86 + miniMember.X86Size
			x64Offset += countBytesPaddingx64 + miniMember.X64Size
			x86Size += miniMember.X86Size + countBytesPaddingx86
			x64Size += miniMember.X64Size + countBytesPaddingx64

		}
		structUnionMini.X86Size = x86Size
		structUnionMini.X64Size = x64Size
	}

	return structsAndUnionsMini
}
