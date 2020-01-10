package pe

import (
	"bytes"
	"encoding/binary"
)

// ImageDelayImportDescriptor represents the _IMAGE_DELAYLOAD_DESCRIPTOR structure.
type ImageDelayImportDescriptor struct {
	Attributes uint32						// Must be zero.
    Name uint32                       		// RVA to the name of the target library (NULL-terminate ASCII string)
    ModuleHandleRVA uint32                  // RVA to the HMODULE caching location (PHMODULE)
    ImportAddressTableRVA uint32           	// RVA to the start of the IAT (PIMAGE_THUNK_DATA)
    ImportNameTableRVA uint32               // RVA to the start of the name table (PIMAGE_THUNK_DATA::AddressOfData)
    BoundImportAddressTableRVA uint32       // RVA to an optional bound IAT
    UnloadInformationTableRVA uint32        // RVA to an optional unload info table
    TimeDateStamp uint32                    // 0 if not bound, oherwise, date/time of the target DLL
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
				maxLen = Max(rva-importDelayDesc.ImportNameTableRVA, rva-importDelayDesc.ImportAddressTableRVA)
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