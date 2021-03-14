// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"regexp"
	"strings"
)

const (
	// RegAPIs is a regex that extract API prototypes.
	RegAPIs = `(_Success_|HANDLE|INTERNETAPI|WINHTTPAPI|BOOLAPI|BOOL|STDAPI|WINUSERAPI|WINBASEAPI|WINADVAPI|NTSTATUS|_Must_inspect_result_|BOOLEAN)[\d\w\s\)\(,\[\]\!*+=&<>/|]+;`

	// RegProto extracts API information.
	RegProto = `(?P<Attr>WINBASEAPI|WINADVAPI)?( )?(?P<RetValType>[A-Z]+) (?P<CallConv>WINAPI|APIENTRY) (?P<ApiName>[a-zA-Z0-9]+)( )?\((?P<Params>.*)\);`

	// RegAPIParams parses params.
	RegAPIParams = `(?P<Anno>_In_|IN|OUT|_In_opt_|_Inout_opt_|_Out_|_Inout_|_Out_opt_|_Outptr_opt_|_Reserved_|_(O|o)ut[\w(),+ *]+|_In[\w()]+|_When[\w() =,!*]+) (?P<Type>[\w *]+) (?P<Name>[*a-zA-Z0-9]+)`

	// RegParam extacts API parameters.
	RegParam = `, `
)

// APIParam represents a paramter of a Win32 API.
type APIParam struct {
	Annotation string `json:"anno"`
	Type       string `json:"type"`
	Name       string `json:"name"`
}

// API represents information about a Win32 API.
type API struct {
	Attribute         string     `json:"-"`        // Microsoft-specific attribute.
	CallingConvention string     `json:"callconv"` // Calling Convention.
	Name              string     `json:"name"`     // Name of the API.
	ReturnValueType   string     `json:"retVal"`   // Return value type.
	Params            []APIParam `json:"params"`   // API Arguments.
	CountParams       uint8      `json:"-"`        // Count of Params.
}

func parseAPIParameter(params string) APIParam {
	m := regSubMatchToMapString(RegAPIParams, params)
	apiParam := APIParam{
		Annotation: m["Anno"],
		Name:       m["Name"],
		Type:       m["Type"],
	}

	// move the `*` to the type.
	if strings.HasPrefix(apiParam.Name, "*") {
		apiParam.Name = apiParam.Name[1:]
		apiParam.Type += "*"
	}

	return apiParam
}

func parseAPI(apiPrototype string) API {
	m := regSubMatchToMapString(RegProto, apiPrototype)
	api := API{
		Attribute:         m["Attr"],
		CallingConvention: m["CallConv"],
		Name:              m["ApiName"],
		ReturnValueType:   m["RetValType"],
	}

	// Treat the VOID case.
	if m["Params"] == " VOID " {
		api.CountParams = 0
		return api
	}

	if api.Name == "" || api.CallingConvention == "" {
		log.Printf("Failed to parse: %s", apiPrototype)
		return api
	}

	re := regexp.MustCompile(RegParam)
	split := re.Split(m["Params"], -1)
	for i, v := range split {
		// Quick hack:
		ss := strings.Split(standardizeSpaces(v), " ")
		if len(ss) == 2 {
			// Force In for API without annotations.
			v = "_In_ " + v
		} else {
			if i+1 < len(split) {
				vv := standardizeSpaces(split[i+1])
				if !strings.HasPrefix(vv, "In") &&
					!strings.HasPrefix(vv, "Out") &&
					!strings.HasPrefix(vv, "_In") &&
					!strings.HasPrefix(vv, "IN") &&
					!strings.HasPrefix(vv, "OUT") &&
					!strings.HasPrefix(vv, "_Reserved") &&
					!strings.HasPrefix(vv, "_When") &&
					!strings.HasPrefix(vv, "__out") &&
					!strings.HasPrefix(vv, "_Out") {
					v += ", " + split[i+1]
					split[i+1] = v
					continue
				}
			}
		}
		api.Params = append(api.Params, parseAPIParameter("_"+v))
		api.CountParams++
	}
	return api
}
