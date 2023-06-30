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

// Registry APIs
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

// File APIs
var (
	fileCreateAPIs = []string{"CreateFileA", "CreateFileW", "CreateDirectory", "CreateDirectoryExA", "CreateDirectoryExW"}
	fileOpenAPIs   = []string{"OpenFile"}
	fileDeleteAPIS = []string{"DeleteFileA", "DeleteFileW"}
	fileWriteAPIs  = []string{"WriteFile", "WriteFileEx"}
	fileReadAPIs   = []string{"ReadFile", "ReadFileEx"}
	fileCopyAPIS   = []string{"CopyFileA", "CopyFileW", "CopyFileExA", "CopyFileExW"}
	fileMoveAPIs   = []string{"MoveFileA", "MoveFileW", "MoveFileWithProgressA", "MoveFileWithProgressW"}
	fileAPIs       = utils.ConcatMultipleSlices([][]string{fileCreateAPIs, fileOpenAPIs, fileDeleteAPIS, fileWriteAPIs, fileReadAPIs})
	fileHandlesMap = make(map[string]string)
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

func fileHandleToStr(hFile string) string {

	filePath, ok := fileHandlesMap[hFile]
	if ok {
		return filePath
	}
	return ""
}

func summarizeFileAPI(w32api Win32API) Event {
	event := Event{}

	// lpFileName points to the name of the file or device to be created or opened. Yo
	lpFileName := w32api.getParamValueByName("lpFileName").(string)

	// hFile represents a handle to the file or I/O device.
	hFile := w32api.getParamValueByName("hFile").(string)

	// lpPathName points to the path of the directory to be created.
	lpPathName := w32api.getParamValueByName("lpPathName").(string)

	// lpNewDirectory points to the path of the directory to be created.
	lpNewDirectory := w32api.getParamValueByName("lpNewDirectory").(string)

	// lpExistingFileName points to the name of an existing file.
	lpExistingFileName := w32api.getParamValueByName("lpExistingFileName").(string)

	// lpNewFileName points to the name of the new file.
	lpNewFileName := w32api.getParamValueByName("lpNewFileName").(string)

	// Th return value of the API, which is a handle in the case of file APIs.
	returnedHandle := w32api.ReturnValue

	if utils.StringInSlice(w32api.Name, fileCreateAPIs) {

		event.Operation = "create"

		// Either a file or a directory creation.
		if lpFileName != "" {
			event.Path = lpFileName
		} else {
			// The Ex version of create directory have a different param name.
			if lpPathName != "" {
				event.Path = lpPathName
			} else {
				event.Path = lpNewDirectory
			}
		}

		// Save the mapping between the handle and its equivalent path.
		fileHandlesMap[returnedHandle] = event.Path

	} else if utils.StringInSlice(w32api.Name, fileOpenAPIs) {

		event.Operation = "open"
		event.Path = lpFileName

		// Save the mapping between the handle and its equivalent path.
		fileHandlesMap[returnedHandle] = event.Path

	} else if utils.StringInSlice(w32api.Name, fileReadAPIs) {

		event.Operation = "read"
		event.Path = fileHandleToStr(hFile)

	} else if utils.StringInSlice(w32api.Name, fileWriteAPIs) {

		event.Operation = "write"
		event.Path = fileHandleToStr(hFile)

	} else if utils.StringInSlice(w32api.Name, fileDeleteAPIS) {

		event.Operation = "delete"
		event.Path = lpFileName

	} else if utils.StringInSlice(w32api.Name, fileCopyAPIS) {
		event.Operation = "copy"
		event.Path = lpExistingFileName + "->" + lpNewFileName

	} else if utils.StringInSlice(w32api.Name, fileMoveAPIs) {
		event.Operation = "move"
		event.Path = lpExistingFileName + "->" + lpNewFileName

	}
	return event
}
