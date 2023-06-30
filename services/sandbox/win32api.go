// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"github.com/saferwall/saferwall/internal/utils"
)

// Win32APIParam describes parameter information for a given Win32 API.
type Win32APIParam struct {
	// SAL annotation.
	Annotation string `json:"anno"`
	// Name of the parameter.
	Name string `json:"name"`
	// Value of the parameter. This can be either a string or a slice of bytes.
	Value interface{} `json:"value"`
}

// Win32API represents a Win32 API event.
type Win32API struct {
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

var (
	regCreateAPIs      = []string{"RegCreateKeyA", "RegCreateKeyW", "RegCreateKeyExA", "RegCreateKeyExW"}
	regOpenAPIs        = []string{"RegOpenKeyA", "RegOpenKeyW", "RegOpenKeyExA", "RegOpenKeyExW"}
	regSetValueAPIs    = []string{"RegSetValueA", "RegSetValueW", "RegSetValueExA", "RegSetValueExW"}
	regSetKeyValueAPIs = []string{"RegSetKeyValueA", "RegSetKeyValueW"}
	regSetAPIs         = utils.ConcatMultipleSlices([][]string{regSetValueAPIs, regSetKeyValueAPIs})
	regDeleteAPIs      = []string{"RegDeleteKeyA", "RegDeleteKeyW", "RegDeleteKeyExA", "RegDeleteKeyExW",
		"RegDeleteValueA", "RegDeleteValueW"}
	regAPIs          = utils.ConcatMultipleSlices([][]string{regCreateAPIs, regOpenAPIs, regSetAPIs, regDeleteAPIs})
	regKeyHandlesMap = make(map[string]string)
)

// Reserved registry key handles.
func init() {

	// We fill in both x86 and x64 reserved registry key handles.
	regKeyHandlesMap["0x80000000"] = "HKEY_CLASSES_ROOT"
	regKeyHandlesMap["0xffffffff80000000"] = "HKEY_CLASSES_ROOT"

	regKeyHandlesMap["0x80000001"] = "HKEY_CURRENT_USER"
	regKeyHandlesMap["0xffffffff80000001"] = "HKEY_CURRENT_USER"

	regKeyHandlesMap["0x80000002"] = "HKEY_LOCAL_MACHINE"
	regKeyHandlesMap["0xffffffff80000002"] = "HKEY_LOCAL_MACHINE"

	regKeyHandlesMap["0x800000003"] = "HKEY_USERS"
	regKeyHandlesMap["0xffffffff80000003"] = "HKEY_USERS"

	regKeyHandlesMap["0x800000004"] = "HKEY_PERFORMANCE_DATA"
	regKeyHandlesMap["0xffffffff80000004"] = "HKEY_PERFORMANCE_DATA"

	regKeyHandlesMap["0x800000005"] = "HKEY_CURRENT_CONFIG"
	regKeyHandlesMap["0xffffffff80000005"] = "HKEY_CURRENT_CONFIG"

	regKeyHandlesMap["0x800000006"] = "HKEY_DYN_DATA"
	regKeyHandlesMap["0xffffffff80000006"] = "HKEY_DYN_DATA"

	regKeyHandlesMap["0x800000007"] = "HKEY_CURRENT_USER_LOCAL_SETTINGS"
	regKeyHandlesMap["0xffffffff80000007"] = "HKEY_CURRENT_USER_LOCAL_SETTINGS"

	regKeyHandlesMap["0x800000050"] = "HKEY_PERFORMANCE_TEXT"
	regKeyHandlesMap["0xffffffff80000050"] = "HKEY_PERFORMANCE_TEXT"

	regKeyHandlesMap["0x800000060"] = "HKEY_PERFORMANCE_NLSTEXT"
	regKeyHandlesMap["0xffffffff80000060"] = "HKEY_PERFORMANCE_NLSTEXT"
}

func (w32api Win32API) getParamValueByName(paramName string) interface{} {

	for _, param := range w32api.Parameters {
		if param.Name == paramName {
			return param.Value
		}
	}

	return ""
}

func summarizeRegAPI(w32api Win32API) Event {
	event := Event{}

	// hKey is either a handle returned by on of registry creation APIs;
	// or it can be one of the predefined keys.
	hKeyStr := w32api.getParamValueByName("hKey").(string)

	// lpSubKey is subkey of the key identified by the hKey parameter.
	lpSubKey := w32api.getParamValueByName("lpSubKey").(string)

	// lpValueName is the name of the registry value whose data is to be updated.
	lpValueName := w32api.getParamValueByName("lpValueName").(string)

	// phkResult is a pointer to a variable that receives a handle to the
	// opened or created key.
	phkResult := w32api.getParamValueByName("phkResult").(string)

	if utils.StringInSlice(w32api.Name, regCreateAPIs) {

		event.Operation = "create"

		event.Path = regKeyHandleToStr(hKeyStr) + "\\" + lpSubKey

		// Save the mapping between the handle and its equivalent path.
		regKeyHandlesMap[phkResult] = event.Path

	} else if utils.StringInSlice(w32api.Name, regOpenAPIs) {

		event.Operation = "open"

		event.Path = regKeyHandleToStr(hKeyStr) + "\\" + lpSubKey

		// Save the mapping between the handle and its equivalent path.
		regKeyHandlesMap[phkResult] = event.Path

	} else if utils.StringInSlice(w32api.Name, regSetAPIs) {

		event.Operation = "write"

		if utils.StringInSlice(w32api.Name, regSetKeyValueAPIs) {
			event.Path = regKeyHandleToStr(hKeyStr) + "\\" + lpSubKey + "\\\\" + lpValueName
		} else if utils.StringInSlice(w32api.Name, regSetValueAPIs) {
			event.Path = regKeyHandleToStr(hKeyStr) + "\\\\" + lpValueName
		}
	} else if utils.StringInSlice(w32api.Name, regDeleteAPIs) {

		event.Operation = "delete"

		event.Path = regKeyHandleToStr(hKeyStr) + "\\" + lpSubKey
		if lpValueName != "" {
			event.Path += "\\\\" + lpValueName
		}

	}
	return event
}

func regKeyHandleToStr(hKey string) string {

	regKey, ok := regKeyHandlesMap[hKey]
	if ok {
		return regKey
	}
	return ""
}
