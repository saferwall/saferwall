package main

import (
	"regexp"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	paramIn uint8 = iota
	paramOut
	paramInOut
	paramReserved
)

const (
	typeImm uint8 = iota
	typePtrImm
	typePtrASCIIStr
	typePtrWideStr
	typeArrASCIIStr
	typeArrWideStr
	typePtrStruct
	typeBytePtr
)

// APIParamMini represents a paramter of a Win32 API.
type APIParamMini struct {
	Annotation uint8  `json:"anno"`
	Type       uint8  `json:"type"`
	Name       string `json:"name"`
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

var (
	// Because widening conversions are always safe, we will treat paramteres
	// either 4 or 8 bytes long depending on what the process is running on.
	immTypes = []string{"int", "DWORD", "WORD", "UINT", "ULONG_PTR", "DWORD_PTR", "HOOKPROC", "HRESULT", "HINSTANCE", "DWORD_PTR",
		"INTERNET_PORT", "BOOL", "ULONG", "SIZE_T", "HKEY", "HINTERNET", "HANDLE", "HMODULE", "SC_HANDLE", "REGSAM", "HHOOK",
		"BCRYPT_KEY_HANDLE", "BCRYPT_HASH_HANDLE",
		"PVOID", "VOID*", "CONST BYTE*", "LPVOID", "PVOID", "LPBYTE", "LPCVOID"}

	//  Pointers to immediate values.
	immPtrTypes = []string{"PBYTE", "LPDWORD", "SIZE_T*", "PHKEY", "PUCHAR", "PHANDLE", "ULONG*"}

	// Void pointers, required further parsing from SAL to read the size from the parameter.
	bytePtrTypes = []string{}

	// Pointer to ascii strings.
	asciiStrTypes = []string{"LPCSTR", "LPSTR"}

	// Pointer to wide stings.
	wideStrTypes = []string{"LPCWSTR", "LPWSTR"}

	// Array of ascii strings.
	arrOfASCIIStrTypes = []string{"LPCSTR FAR *", "LPCSTR*"}

	// Aarray of wide strings.
	arrOfWideStrTypes = []string{"LPCWSTR FAR *", "LPCWSTR*"}

	// Pointer to a struct.
	ptrStructTypes = []string{"LPURL_COMPONENTSA", "LPINTERNET_BUFFERSW", "LPSERVICE_STATUS", "PFILETIME", "LPSTARTUPINFOW", "LPPROCESS_INFORMATION",
		"LPSECURITY_ATTRIBUTES", "LPPROCESSENTRY32W", "PLUID", "LPINTERNET_BUFFERSW", "LPENUM_SERVICE_STATUSW", "LPSTARTUPINFOA", "SERVICE_TABLE_ENTRYW", "LPURL_COMPONENTSW",
		"CONST LPSECURITY_ATTRIBUTES", "CONST SERVICE_TABLE_ENTRYW*",
		"LPENUM_SERVICE_STATUSA", "LPINTERNET_BUFFERSA", "PMEMORY_BASIC_INFORMATION"}
)

func convertStrType(t string) uint8 {

	if utils.StringInSlice(t, immTypes) {
		return typeImm
	} else if utils.StringInSlice(t, immPtrTypes) {
		return typePtrImm
	} else if utils.StringInSlice(t, bytePtrTypes) {
		return typeBytePtr
	} else if utils.StringInSlice(t, asciiStrTypes) {
		return typePtrASCIIStr
	} else if utils.StringInSlice(t, wideStrTypes) {
		return typePtrWideStr
	} else if utils.StringInSlice(t, arrOfASCIIStrTypes) {
		return typeArrASCIIStr
	} else if utils.StringInSlice(t, arrOfWideStrTypes) {
		return typeArrWideStr
	} else if utils.StringInSlice(t, ptrStructTypes) {
		return typePtrStruct
	} else {
		return 0
	}
}

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
				parammini.Name = param.Name
				parammini.Type = convertStrType(param.Type)
				paramsMini = append(paramsMini, parammini)
			}
			copy.Params = paramsMini
			mapis[dllname][apiname] = copy
		}
	}

	return mapis
}
