// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"encoding/binary"
)

// ImageDelayImportDescriptor represents the _IMAGE_DELAYLOAD_DESCRIPTOR structure.
type ImageDelayImportDescriptor struct {
	// As yet, no attribute flags are defined. The linker sets this field to zero
	// in the image. This field can be used to extend the record by indicating
	// the presence of new fields, or it can be used to indicate behaviors to
	// the delay or unload helper functions.
	Attributes uint32

	// The name of the DLL to be delay-loaded resides in the read-only data
	// section of the image. It is referenced through the szName field.
	Name uint32

	// The handle of the DLL to be delay-loaded is in the data section of the
	// image. The phmod field points to the handle. The supplied delay-load
	// helper uses this location to store the handle to the loaded DLL.
	ModuleHandleRVA uint32

	// The delay import address table (IAT) is referenced by the delay import
	// descriptor through the pIAT field. The delay-load helper updates these
	// pointers with the real entry points so that the thunks are no longer in
	// the calling loop
	ImportAddressTableRVA uint32

	// The delay import name table (INT) contains the names of the imports that
	// might require loading. They are ordered in the same fashion as the
	// function pointers in the IAT.
	ImportNameTableRVA uint32

	// The delay bound import address table (BIAT) is an optional table of
	// IMAGE_THUNK_DATA items that is used along with the timestamp field
	// of the delay-load directory table by a post-process binding phase.
	BoundImportAddressTableRVA uint32

	// The delay unload import address table (UIAT) is an optional table of
	// IMAGE_THUNK_DATA items that the unload code uses to handle an explicit
	// unload request. It consists of initialized data in the read-only section
	// that is an exact copy of the original IAT that referred the code to the
	// delay-load thunks. On the unload request, the library can be freed,
	// the *phmod cleared, and the UIAT written over the IAT to restore
	// everything to its preload state.
	UnloadInformationTableRVA uint32

	// 0 if not bound, otherwise, date/time of the target DLL.
	TimeDateStamp uint32
}

// DelayImport represents an entry in the delay import table.
type DelayImport struct {
	Offset     uint32
	Name       string
	Functions  []*ImportFunction
	Descriptor ImageDelayImportDescriptor
}

// Delay-Load Import Tables tables were added to the image to support a uniform
// mechanism for applications to delay the loading of a DLL until the first call
// into that DLL. The delay-load directory table is the counterpart to the
// import directory table.
func (pe *File) parseDelayImportDirectory(rva, size uint32) error {
	for {
		importDelayDesc := ImageDelayImportDescriptor{}
		fileOffset := pe.getOffsetFromRva(rva)
		importDescSize := uint32(binary.Size(importDelayDesc))
		err := pe.structUnpack(&importDelayDesc, fileOffset, importDescSize)

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

		nameRVA := uint32(0)
		if importDelayDesc.Attributes == 0 {
			nameRVA = importDelayDesc.Name -
				pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).ImageBase
		} else {
			nameRVA = importDelayDesc.Name
		}
		dllName := pe.getStringAtRVA(nameRVA, maxLen)
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

// GetDelayImportEntryInfoByRVA return an import function + index of the entry given
// an RVA.
func (pe *File) GetDelayImportEntryInfoByRVA(rva uint32) (DelayImport, int) {
	for _, imp := range pe.DelayImports {
		for i, entry := range imp.Functions {
			if entry.ThunkRVA == rva {
				return imp, i
			}
		}
	}

	return DelayImport{}, 0
}
