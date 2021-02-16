// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
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
	X86Offset  int             `json:"x86off"`
	X86Size    int             `json:"x86size"`
	X64Offset  int             `json:"x64off"`
	X64Size    int             `json:"x64size"`
	Type       uint8           `json:"type"` // Help interpret the value.
	Definition StructUnionMini `json:"def,omitempty"`
}

// StructUnionMini represents a complex type like a Union or a Struct.
type StructUnionMini struct {
	Name    string                  `json:"name"`
	Members []StructUnionMemberMini `json:"members"`
	X86Size int                     `json:"x86size"`
	X64Size int                     `json:"x64size"`
}

// APIParamMini represents a paramter of a Win32 API.
type APIParamMini struct {
	Annotation        uint8  `json:"anno"`
	Type              uint8  `json:"type"`
	Name              string `json:"name"`
	BufferSizeOrIndex int8   `json:"buffsize_or_idx,omitempty"`
}

// APIMini represents information about a Win32 API.
type APIMini struct {
	ReturnValueType bool           `json:"retType"` // Return value type.
	Params          []APIParamMini `json:"params"`  // API Arguments.
}

var (
	reAnnotationIn       = regexp.MustCompile(`_In_|IN|_In_opt[\w]+|In_reads[\w()]+`)
	reAnnotationOut      = regexp.MustCompile(`_Out_|OUT|_Out_opt[\w]+|_Out_writes[\w()]+|_Outptr_`)
	reAnnotationIntOut   = regexp.MustCompile(`_Inout[\w]+`)
	reAnnotationReserved = regexp.MustCompile(`Reserved`)
)


func minifyAPIs(apis map[string]map[string]API) map[string]map[string]APIMini {
	mapis := make(map[string]map[string]APIMini)
	for dllname, v := range apis {
		if _, ok := mapis[dllname]; !ok {
			mapis[dllname] = make(map[string]APIMini)
		}
		for apiname, vv := range v {

			returnType := false
			if vv.ReturnValueType == "VOID" {
				returnType = true
			}

			copy := APIMini{
				ReturnValueType: returnType}

			var paramsMini []APIParamMini
			for _, param := range vv.Params {
				parammini := APIParamMini{}
				if reAnnotationIn.MatchString(param.Annotation) {
					parammini.Annotation = paramIn
				} else if reAnnotationOut.MatchString(param.Annotation) {
					parammini.Annotation = paramOut
				} else if reAnnotationIntOut.MatchString(param.Annotation) {
					parammini.Annotation = paramInOut
				} else if reAnnotationReserved.MatchString(param.Annotation) {
					parammini.Annotation = paramReserved
				} else {
					// If we don't know, take it as in:
					parammini.Annotation = paramIn
				}

				// Get the param type.
				dataType := typefromString(param.Type)
				parammini.Type = dataType.Kind
				parammini.BufferSizeOrIndex = dataType.Size
				parammini.Name = param.Name
				paramsMini = append(paramsMini, parammini)
			}
			copy.Params = paramsMini
			mapis[dllname][apiname] = copy
		}
	}

	return mapis
}
