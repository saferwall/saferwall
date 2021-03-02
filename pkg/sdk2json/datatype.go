// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// https://docs.microsoft.com/en-us/cpp/cpp/fundamental-types-cpp
// https://docs.microsoft.com/en-us/cpp/cpp/data-type-range

package main

import (
	"regexp"
	"strings"
)

const (
	oneByte = iota
	twoBytes
	fourBytes
	eightBytes
)

const (
	typeScalar      uint8 = iota // Scalar types: int, long, DWORD, ...
	typePtrScalar                // Pointers to scalar types: int*, ...
	typeASCIIStr                 // NULL terminated ascii string.
	typeWideStr                  // NULL terminated wide string.
	typeArrASCIIStr              // Array of null terminated ascii strings.
	typeArrWideStr               // Array of null terminated wide strings.
	typeStruct                   // Struct.
	typeBytePtr                  // Byte pointer, size is known.
	typeVoidPtr                  // Void*, not known to what it points to.
)

type dataType struct {

	// Name of the type: PVOID, DWORD_PTR, ...
	Name string
	// Kind of data type: typeScalar, ...
	Kind uint8
	// Size in bytes, for pointers it holds the size of the type it points to, not the size of the pointer itself.
	Size int8
}

var (
	regAllTypedef = `(?m)^typedef(\s)+(?P<Source>[\w\s]+)+(\s)+(?P<Target>[*\w\s,]+);`

	// 1 byte types,
	oneByteTypes = []string{
		"bool", "char", "char8_t", "unsigned char", "signed char", "__int8",
	}

	// 2 bytes types.
	twoByteTypes = []string{
		"char16_t", "__int16", "short", "unsigned short", "wchar_t", "__wchar_t", "short int", "signed short", "signed short int",
	}

	// 4 bytes types.
	fourByteTypes = []string{
		"char32_t", "float", "__int32", "int", "signed", "unsigned",
		"signed int", "unsigned int", "long", "long int", "signed long",
		"signed long int", "unsigned long", "unsigned long int",
	}

	// 8 bytes types.
	eightByteTypes = []string{
		"double", "__int64", "long double", "long long", "unsigned long long",
		"unsigned long long int", "long long int", "signed long long",
		"signed long long int",
	}

	// Void* types + HANDLE-types + XYZ_PTR like types.
	// The `Handle` types was declated with `DECLARE_HANDLE` macro
	// instead of a direct typedef. We hardcode them here for now.
	voidPtrTypes = []string{
		"void*", "VOID*", "HANDLE", "HKEY", "HMETAFILE", "HINSTANCE", "HRGN", "HRSRC", "HSPRITE", "HLSURF", "HSTR", "HTASK", "SC_HANDLE", "ULONG_PTR", "LONG_PTR",
	}

	// ASCIIStrTypes represents the null terminated string types.
	ASCIIStrTypes = []string{
		"_Null_terminated_ CHAR*", "_Null_terminated_ char*",
	}

	// wide null terminated string types.
	wideStrTypes = []string{
		"_Null_terminated_ WCHAR*", "_Null_terminated_ wchar_t*",
	}

	// Array of null terminated ascii string types.
	arrASCIIStrTypes = []string{
		"_Null_terminated_ PSTR*",
	}

	// Array of null terminated wide string types.
	arrWideStrTypes = []string{
		"_Null_terminated_ PWSTR*",
	}

	// Byte pointer types.
	bytePtrType = []string{
		"BYTE*",
	}

	// Maps a type to its typedef alias.
	typedefs = map[string]string{}

	// Maps a type to its dataType object.
	dataTypes = map[string]dataType{}
)

// Fill in built in types.
func initBuiltInTypes() {
	for _, t := range oneByteTypes {
		dataTypes[t] = dataType{Name: t, Size: 1, Kind: typeScalar}
		dataTypes[t+"*"] = dataType{Name: t + "*", Size: 1, Kind: typePtrScalar}
	}

	for _, t := range twoByteTypes {
		dataTypes[t] = dataType{Name: t, Size: 2, Kind: typeScalar}
		dataTypes[t+"*"] = dataType{Name: t + "*", Size: 2, Kind: typePtrScalar}
	}

	for _, t := range fourByteTypes {
		dataTypes[t] = dataType{Name: t, Size: 4, Kind: typeScalar}
		dataTypes[t+"*"] = dataType{Name: t + "*", Size: 4, Kind: typePtrScalar}
	}

	for _, t := range eightByteTypes {
		dataTypes[t] = dataType{Name: t, Size: 8, Kind: typeScalar}
		dataTypes[t+"*"] = dataType{Name: t + "*", Size: 8, Kind: typePtrScalar}
	}

	for _, t := range voidPtrTypes {
		dataTypes[t] = dataType{Name: t, Size: 0, Kind: typeVoidPtr}
		dataTypes[t+"*"] = dataType{Name: t + "*", Size: 0, Kind: typePtrScalar}
	}

	for _, t := range ASCIIStrTypes {
		dataTypes[t] = dataType{Name: t, Size: 0, Kind: typeASCIIStr}
	}

	for _, t := range wideStrTypes {
		dataTypes[t] = dataType{Name: t, Size: 0, Kind: typeWideStr}
	}

	for _, t := range arrASCIIStrTypes {
		dataTypes[t] = dataType{Name: t, Size: 0, Kind: typeArrASCIIStr}
	}

	for _, t := range arrWideStrTypes {
		dataTypes[t] = dataType{Name: t, Size: 0, Kind: typeArrWideStr}
	}

	for _, t := range bytePtrType {
		dataTypes[t] = dataType{Name: t, Size: 0, Kind: typeBytePtr}
	}
}

// Create custom data types: CHAR, DWORD ..
// usually typedefs to built in types.
func initCustomTypes() {
	// We repeat this process 2 times as some types won't be know only after
	// first iteration.
	for i := 0; i < 3; i++ {

		for k, v := range typedefs {
			// No need to go further if the type is already known.
			if _, ok := dataTypes[k]; ok {
				continue
			}
			
			// Search in our typedef map
			val, ok := dataTypes[v]
			if !ok {
				// Take out the `*` and look up again.
				val, ok := dataTypes[v[:len(v)-1]]
				if ok {
					if val.Kind == typeScalar {
						dt := dataType{Name: k, Size: val.Size, Kind: typePtrScalar}
						dataTypes[k] = dt
					}
				}
			} else {
				dt := dataType{Name: k, Size: val.Size, Kind: val.Kind}
				dataTypes[k] = dt

				// We found a new type, if it is a scalar, let's add its
				// pointer version as well.
				if val.Kind == typeScalar {
					if _, ok := dataTypes[k+"*"]; !ok {
						dt := dataType{
							Name: k + "*", Size: val.Size, Kind: typePtrScalar}
						dataTypes[k+"*"] = dt
					}
				} else if val.Kind == typeASCIIStr {
					if _, ok := dataTypes[k+"*"]; !ok {
						dt := dataType{
							Name: k + "*", Size: val.Size, Kind: typeArrASCIIStr}
						dataTypes[k+"*"] = dt
					}
				} else if val.Kind == typeWideStr {
					if _, ok := dataTypes[k+"*"]; !ok {
						dt := dataType{
							Name: k + "*", Size: val.Size, Kind: typeArrWideStr}
						dataTypes[k+"*"] = dt
					}
				}
			}
		}
	}
}

func typefromString(t string) dataType {

	// Remove non-important C language modifiers like CONST ...
	t = strings.ReplaceAll(t, "CONST ", "")
	t = strings.ReplaceAll(t, " FAR", "")
	t = strings.ReplaceAll(t, " NEAR", "")
	t = spaceFieldsJoin(t)

	if dt, ok := dataTypes[t]; ok {
		return dt
	}

	//log.Println(t)

	return dataType{Name: t, Kind: typeStruct, }
}

func parseTypedefs(data []byte) {

	// Retrieve all typedeffed names.
	re := regexp.MustCompile(regAllTypedef)
	matches := re.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		// Strip extra white spaces from typedef statement.
		srcName := standardizeSpaces(match[2])
		newName := standardizeSpaces(match[4])

		// the newName in typedef could include multiple names:
		// i.e:typedef _Null_terminated_ CHAR *NPSTR, *LPSTR, *PSTR;
		elements := strings.Split(newName, ",")
		for _, val := range elements {
			src := srcName
			dest := spaceFieldsJoin(val)

			// Take out some modifiers like `CONST`, `near`, `far`,
			// as they don't affect the type, but more a hint for the compiler.
			src = strings.ReplaceAll(src, "CONST ", "")
			src = strings.ReplaceAll(src, " far", "")
			src = strings.ReplaceAll(src, " near", "")

			// When the data type is a pointer.
			if strings.HasPrefix(dest, "*") {
				src += "*"
				dest = dest[1:]
			}

			typedefs[dest] = src

		}
	}
}
