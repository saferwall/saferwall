// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/saferwall/saferwall/internal/utils"
)

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

// Networking APIs
var (
	netWinHTTPAPIs = []string{"WinHttpConnect"}
	netWinINetAPIs = []string{"InternetConnectA", "InternetConnectW"}
	netWinsockAPIs = []string{"getaddrinfo", "GetAddrInfoW", "GetAddrInfoExA", "GetAddrInfoExW"}
	netAPIs        = utils.ConcatMultipleSlices([][]string{netWinINetAPIs, netWinHTTPAPIs, netWinsockAPIs})
	netHandlesMap  = make(map[string]string)
)

// Default port numbers for WinINet.
const (
	InvalidPortNumber = 0
	DefaultFtpPort    = 21
	DefaultGopherPort = 70
	DefaultHTTPPort   = 80
	DefaultHTTPSPort  = 443
	DefaultSocksPort  = 1080
)

// Track unique system events.
var (
	uniqueEvents = make(map[string]bool)
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
			if param.Annotation == APIParamAnnotationIn ||
				param.Annotation == APIParamAnnotationOut {
				return param.Value
			} else {
				// TODO:  Do we want to always return the OUT value,
				// or we can be flexible and add an argument in the function
				// that influence the decision.
				return param.Value
			}
		}
	}

	return ""
}

func summarizeRegAPI(w32api Win32API) Event {
	event := Event{Type: registryEventType, ProcessID: w32api.ProcessID}

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
	event := Event{Type: fileEventType, ProcessID: w32api.ProcessID}

	// lpFileName points to the name of the file or device to be created or opened.
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

func summarizeNetworkAPI(w32api Win32API) Event {
	event := Event{Type: networkEventType, ProcessID: w32api.ProcessID}

	if utils.StringInSlice(w32api.Name, netWinsockAPIs) {
		// pNodeName contains a host (node) name or a numeric host address string.
		pNodeName := w32api.getParamValueByName("pNodeName").(string)

		// pName contains a host (node) name or a numeric host address string.
		pName := w32api.getParamValueByName("pName").(string)

		// pServiceName contains either a service name or port number represented as a string.
		pServiceName := w32api.getParamValueByName("pServiceName").(string)

		if pNodeName != "" {
			event.Path = pNodeName
		} else {
			event.Path = pName
		}

		event.Path += ":" + pServiceName
		event.Operation = "socket"
	} else {

		// lpszServerName specifies the host name of an Internet server.
		lpszServerName := w32api.getParamValueByName("lpszServerName").(string)

		// nServerPort represents the TCP/IP port on the server.
		nServerPort := w32api.getParamValueByName("nServerPort").(string)

		// pswzServerName contains the host name of an HTTP server.
		pswzServerName := w32api.getParamValueByName("pswzServerName").(string)

		// The return value of the API, which is a handle in the case of network APIs.
		returnedHandle := w32api.ReturnValue

		serverPort, _ := strconv.ParseInt(nServerPort, 0, 64)

		if utils.StringInSlice(w32api.Name, netWinHTTPAPIs) {
			event.Path = pswzServerName
		} else if utils.StringInSlice(w32api.Name, netWinINetAPIs) {
			event.Path = lpszServerName
		}

		switch serverPort {
		case InvalidPortNumber:
			// Uses the default port for the service specified by dwService.
			dwService := w32api.getParamValueByName("dwService").(string)
			svcPort, _ := strconv.ParseInt(dwService, 0, 64)
			switch svcPort {
			case 0x1:
				event.Operation = "ftp"
				serverPort = 21
			case 0x2:
				event.Operation = "gopher"
				serverPort = 70
			case 0x3:
				event.Operation = "http"
				serverPort = 80
			}
		case DefaultFtpPort:
			event.Operation = "ftp"
		case DefaultHTTPPort:
			event.Operation = "http"
		case DefaultHTTPSPort:
			event.Operation = "https"
		case DefaultSocksPort:
			event.Operation = "socks"
		}

		event.Path += ":" + strconv.Itoa(int(serverPort))

		// Save the mapping between the handle and its equivalent path.
		netHandlesMap[returnedHandle] = event.Path
	}

	return event
}

func (s *Service) isNewEvent(event Event) bool {
	eventKey := fmt.Sprintf("%s-%s-%s", event.ProcessID,
		strings.ToLower(event.Path), event.Operation)
	_, ok := uniqueEvents[eventKey]
	if ok {
		return false
	} else {
		uniqueEvents[eventKey] = true
		return true
	}
}

func (s *Service) curateAPIEvents(w32apis []Win32API) []byte {
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
