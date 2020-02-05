// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
)

// ImageDelayImportDescriptor represents the _IMAGE_DELAYLOAD_DESCRIPTOR structure.
type ImageDelayImportDescriptor struct {
	// Must be zero.
	Attributes uint32

	// RVA to the name of the target library (NULL-terminate ASCII string).
	Name uint32

	// RVA to the HMODULE caching location (PHMODULE).
	ModuleHandleRVA uint32

	// RVA to the start of the IAT (PIMAGE_THUNK_DATA).
	ImportAddressTableRVA uint32

	// RVA to the start of the name table (PIMAGE_THUNK_DATA::AddressOfData).
	ImportNameTableRVA uint32

	// RVA to an optional bound IAT.
	BoundImportAddressTableRVA uint32

	// RVA to an optional unload info table.
	UnloadInformationTableRVA uint32

	// 0 if not bound, oherwise, date/time of the target DLL.
	TimeDateStamp uint32
}

// DelayImport represents an entry in the delay import table.
type DelayImport struct {
	Offset     uint32
	Name       string
	Functions  []*ImportFunction
	Descriptor ImageDelayImportDescriptor
}

// The delay-load directory table is the counterpart to the import directory table.
func (pe *File) parseDelayImportDirectory(rva, size uint32) error {
	for {
		importDelayDesc := ImageDelayImportDescriptor{}
		fileOffset := pe.getOffsetFromRva(rva)
		importDescSize := uint32(binary.Size(importDelayDesc))
		buf := bytes.NewReader(pe.data[fileOffset : fileOffset+importDescSize])
		err := binary.Read(buf, binary.LittleEndian, &importDelayDesc)
		// If the RVA is invalid all would blow up. Some EXEs seem to be
		// specially nasty and have an invalid RVA.
		if err != nil {
			return err
		}

		// If the structure is all zeros, we reached the end of the list.
		if importDelayDesc == (ImageDelayImportDescriptor{}) {
			break
		}

		rva += importDescSize

		// If the array of thunks is somewhere earlier than the import
		// descriptor we can set a maximum length for the array. Otherwise
		// just set a maximum length of the size of the file
		maxLen := uint32(len(pe.data)) - fileOffset
		if rva > importDelayDesc.ImportNameTableRVA ||
			rva > importDelayDesc.ImportAddressTableRVA {
			if rva < importDelayDesc.ImportNameTableRVA {
				maxLen = rva - importDelayDesc.ImportAddressTableRVA
			} else if rva < importDelayDesc.ImportAddressTableRVA {
				maxLen = rva - importDelayDesc.ImportNameTableRVA
			} else {
				maxLen = Max(rva-importDelayDesc.ImportNameTableRVA,
					rva-importDelayDesc.ImportAddressTableRVA)
			}
		}

		var importedFunctions []*ImportFunction
		if pe.Is64 {
			importedFunctions, err = pe.parseImports64(&importDelayDesc, maxLen)
		} else {
			importedFunctions, err = pe.parseImports32(&importDelayDesc, maxLen)
		}
		if err != nil {
			return err
		}

		dllName := pe.getStringAtRVA(importDelayDesc.Name, maxLen)
		if !IsValidDosFilename(dllName) {
			dllName = "*invalid*"
			continue
		}

		pe.DelayImports = append(pe.DelayImports, DelayImport{
			Offset:     fileOffset,
			Name:       string(dllName),
			Functions:  importedFunctions,
			Descriptor: importDelayDesc,
		})

	}

	return nil
}
