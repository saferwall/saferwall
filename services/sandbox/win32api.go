// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

// Win32APIParam describes parameter information for a given Win32 API.
type Win32APIParam struct {
	// SAL annotation.
	Annotation string `json:"anno"`
	// Name of the parameter.
	Name string `json:"name"`
	// Value of the parameter. This can be either a string or a slice of bytes.
	// This field is mutually exclusive with the In and Out values.
	Value interface{} `json:"val,omitempty"`
	// Win32 API sometimes uses IN and OUT annotations, so instead of having
	// one `value`, we separate the `in` and `out`. Occasionally, a function can
	// both reads from and writes to buffer, so ValueIn and ValueOut are filled.
	// The function reads from the buffer.
	InValue interface{} `json:"in_val,omitempty"`
	// The function writes to the buffer.
	OutValue interface{} `json:"out_val,omitempty"`

	// An ID is attributed to track BYTE* parameters that spans over 4KB of data.
	// If the buffer is either IN or OUT, the ID will be on `BuffID`, otherwise:
	// BuffIDIn and BufferIdOut
	BufID    string `json:"buf_id,omitempty"`
	InBufID  string `json:"in_buf_id_in,omitempty"`
	OutBufID string `json:"out_buf_id,omitempty"`
}

// Win32API represents a Win32 API event.
type Win32API struct {
	// Timestamp of the trace.
	Timestamp int64 `json:"ts"`
	// Name of the API.
	Name string `json:"api"`
	// List of its parameters.
	Parameters []Win32APIParam `json:"params"`
	// Process Identifier responsible for generating the API.
	ProcessID string `json:"pid"`
	// Thread Identifier responsible for generating the API.
	ThreadID string `json:"tid"`
	// The name of the process that corresponds to the process ID.
	ProcessName string `json:"proc"`
	// Return value of the API.
	ReturnValue string `json:"ret"`
}

// Win32APIBuffer represents a Win32 API large buffer of parameter of type BYTE*.
type Win32APIBuffer struct {
	// Name of the buffer.
	Name string
	// Content of the buffer.
	Content []byte
}

// SAL Win32 API annotation for function parameters.
const (
	APIParamAnnotationIn       = "in"
	APIParamAnnotationOut      = "out"
	APIParamAnnotationInOut    = "in_out"
	APIParamAnnotationReserved = "reserved"
)

func curateAPIEvents(w32apis []Win32API) []byte {
	var curatedAPIs []interface{}
	for _, w32api := range w32apis {
		curatedAPI := make(map[string]interface{})
		curatedAPI["name"] = w32api.Name
		curatedAPI["ts"] = w32api.Timestamp
		curatedAPI["pid"] = w32api.ProcessID
		curatedAPI["tid"] = w32api.ThreadID
		curatedAPIArgs := make([]map[string]interface{}, len(w32api.Parameters))
		for i, w32Param := range w32api.Parameters {
			curatedAPIArgs[i] = make(map[string]interface{})
			if w32Param.Annotation == APIParamAnnotationIn ||
				w32Param.Annotation == APIParamAnnotationOut ||
				w32Param.Annotation == APIParamAnnotationReserved {
				curatedAPIArgs[i]["val"] = w32Param.Value
				if w32Param.BufID != "" {
					curatedAPIArgs[i]["buf_id"] = w32Param.BufID
				}
			} else {
				curatedAPIArgs[i]["in"] = w32Param.InValue
				curatedAPIArgs[i]["out"] = w32Param.OutValue
				if w32Param.InBufID != "" {
					curatedAPIArgs[i]["in_buf_id"] = w32Param.InBufID
				}
				if w32Param.OutBufID != "" {
					curatedAPIArgs[i]["out_buf_id"] = w32Param.OutBufID
				}
			}
		}
		curatedAPI["ret"] = w32api.ReturnValue
		curatedAPI["args"] = curatedAPIArgs
		curatedAPIs = append(curatedAPIs, curatedAPI)
	}

	return toJSON(curatedAPIs)
}
