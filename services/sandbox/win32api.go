// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"strconv"

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
	regKeyHandlesMap = make(map[uint32]string)
)

// ReservedRegKeyHandleType represents the type of a reserved registry key handle.
type ReservedRegKeyHandleType uint32

// Reserved registry key handles.
const (
	HKeyClassesRoot              = 0x80000000
	HKeyCurrentUser              = 0x80000001
	HKeyLocalMachine             = 0x80000002
	HKeyUsers                    = 0x80000003
	HKeyPerformanceData          = 0x80000004
	HKeyCurrentConfig            = 0x80000005
	HKeyDynData                  = 0x80000006
	HKeyCurrentUserLocalSettings = 0x80000007
	HKeyPerformanceText          = 0x80000050
	HKeyPerformanceNlsText       = 0x80000060
)

// String returns the string representation of a reserved registry key handle.
func (regKeyHandle ReservedRegKeyHandleType) String() string {
	regHandlesMap := map[ReservedRegKeyHandleType]string{
		HKeyClassesRoot:              "HKEY_CLASSES_ROOT",
		HKeyCurrentUser:              "HKEY_CURRENT_USER",
		HKeyLocalMachine:             "HKEY_LOCAL_MACHINE",
		HKeyUsers:                    "HKEY_USERS",
		HKeyPerformanceData:          "HKEY_PERFORMANCE_DATA",
		HKeyCurrentConfig:            "HKEY_CURRENT_CONFIG",
		HKeyDynData:                  "HKEY_DYN_DATA",
		HKeyCurrentUserLocalSettings: "HKEY_CURRENT_USER_LOCAL_SETTINGS",
		HKeyPerformanceText:          "HKEY_PERFORMANCE_TEXT",
		HKeyPerformanceNlsText:       "HKEY_PERFORMANCE_NLSTEXT",
	}

	if val, ok := regHandlesMap[regKeyHandle]; ok {
		return val
	}

	return ""
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
	phkResult := w32api.getParamValueByName("phkResult").(uint32)

	if utils.StringInSlice(w32api.Name, regCreateAPIs) {

		event.Operation = "create"

		event.Path = regKeyHandleToStr(hKeyStr) + "/" + lpSubKey

		// Save the mapping between the handle and its equivalent path.
		regKeyHandlesMap[phkResult] = event.Path

	} else if utils.StringInSlice(w32api.Name, regOpenAPIs) {

		event.Operation = "open"

		event.Path = regKeyHandleToStr(hKeyStr) + "/" + lpSubKey

		// Save the mapping between the handle and its equivalent path.
		regKeyHandlesMap[phkResult] = event.Path

	} else if utils.StringInSlice(w32api.Name, regSetAPIs) {

		event.Operation = "write"

		if utils.StringInSlice(w32api.Name, regSetKeyValueAPIs) {
			event.Path = regKeyHandleToStr(hKeyStr) + "/" + lpSubKey + "//" + lpValueName
		} else if utils.StringInSlice(w32api.Name, regSetValueAPIs) {
			event.Path = regKeyHandleToStr(hKeyStr) + "//" + lpValueName
		}
	} else if utils.StringInSlice(w32api.Name, regDeleteAPIs) {

		event.Operation = "delete"

		event.Path = regKeyHandleToStr(hKeyStr) + "/" + lpSubKey
		if lpValueName != "" {
			event.Path += "//" + lpValueName
		}

		// Save the mapping between the handle and its equivalent path.
		regKeyHandlesMap[phkResult] = event.Path
	}
	return event
}

func regKeyHandleToStr(hKey string) string {
	// Convert it to an integer.
	hKeyInt, _ := strconv.Atoi(hKey)

	// Try to see if we are dealing with a pre-defined registry key handle.
	hKeyStr := ReservedRegKeyHandleType(hKeyInt).String()
	if hKeyStr != "" {
		return hKeyStr
	}

	// hKey should be resolved.
	regKey, ok := regKeyHandlesMap[uint32(hKeyInt)]
	if !ok {
		return ""
	}

	return regKey
}
