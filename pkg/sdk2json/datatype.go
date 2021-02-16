// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// https://docs.microsoft.com/en-us/cpp/cpp/fundamental-types-cpp
// https://docs.microsoft.com/en-us/cpp/cpp/data-type-range

package main

import (
	"log"
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
	typeScalar      uint8 = iota // Scalar types: int, char, ...
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
	// Size in bytes, for pointers it holds the size of the type it points to, not the size of the pointer itself. For void*, we set it to -1.
	Size int8
}

var (
	// regAllTypedef = `(?m)^typedef(\s)+(?P<Source>[\w\s]+)+(\s)+(?P<Target>[*\w]+);`
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

	// Void* types. The `Handle` types was declated with `DECLARE_HANDLE` macro
	// instead of a direct typedef. We hardcode here for now.
	voidPtrTypes = []string{
		"void*", "VOID*", "HKEY", "HMETAFILE", "HINSTANCE", "HRGN", "HRSRC",
		"HSPRITE", "HLSURF", "HSTR", "HTASK",
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
	}

	for _, t := range twoByteTypes {
		dataTypes[t] = dataType{Name: t, Size: 2, Kind: typeScalar}
	}

	for _, t := range fourByteTypes {
		dataTypes[t] = dataType{Name: t, Size: 4, Kind: typeScalar}
	}

	for _, t := range eightByteTypes {
		dataTypes[t] = dataType{Name: t, Size: 8, Kind: typeScalar}
	}

	for _, t := range voidPtrTypes {
		dataTypes[t] = dataType{Name: t, Size: -1, Kind: typeVoidPtr}
	}
}

// Create custom data types: CHAR, DWORD ..
// usually typedefs to built in types.
func initCustomTypes() {
	// We repeat this process 2 times as some types won't be know only after
	// first iteration.
	for i := 0; i < 3; i++ {
		log.Println(len(dataTypes))

		for k, v := range typedefs {
			// No need to go further if the type is already known.
			if _, ok := dataTypes[k]; ok {
				continue
			}

			// Search in our typedef map
			val, ok := dataTypes[v]
			if !ok {
				log.Printf("We dont have %s for now\n", v)
				// Take out the `*` and look up again.
				val, ok := dataTypes[v[:len(v)-1]]
				if ok {
					dt := dataType{Name: k, Size: val.Size, Kind: typePtrScalar}
					dataTypes[k] = dt
					continue
				}
				//
			} else {
				log.Printf("%s is found on our map\n", v)
				dt := dataType{Name: k, Size: val.Size, Kind: val.Kind}
				dataTypes[k] = dt
			}
		}
	}
	log.Println("Done")

}

func typefromString(t string) dataType {

	if dt, ok := dataTypes[t]; ok {
		return dt
	}

	log.Println(t)

	return dataType{}
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
			// When the data type is a pointer.
			if strings.HasPrefix(dest, "*") {
				src += "*"
				dest = dest[1:]
			}
			typedefs[dest] = src

		}
	}
}
