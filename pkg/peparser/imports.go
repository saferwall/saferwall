// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	imageOrdinalFlag32   = uint32(0x80000000)
	imageOrdinalFlag64   = uint64(0x8000000000000000)
	maxRepeatedAddresses = uint32(0xF)
	maxAddressSpread     = uint32(0x8000000)
	addressMask32        = uint32(0x7fffffff)
	addressMask64        = uint64(0x7fffffffffffffff)
	maxDllLength         = 0x200
	maxImportNameLength  = 0x200
)

var (
	// AnoNullImportTable is reported when the ILT or the IAT are empty.
	AnoNullImportTable = "Damaged Import Table information. ILT and/or IAT" +
		" appear to be broken"

	// AnoInvalidThunkAddressOfData is reported when thunk address is too spread out.
	AnoInvalidThunkAddressOfData = "Thunk Address Of Data too spread out"

	// AnoManyRepeatedEntries is reported when import directory contains many
	// entries have the same RVA.
	AnoManyRepeatedEntries = "Import directory contains many repeated entries"

	// AnoAddressOfDataBeyondLimits is reported when Thunk AddressOfData goes
	// beyond limites.
	AnoAddressOfDataBeyondLimits = "Thunk AddressOfData beyond limits"
)

// ImageImportDescriptor describes the remainder of the import information.
// The import directory table contains address information that is used to
// resolve fixup references to the entry points within a DLL image.
// It consists of an array of import directory entries, one entry for each DLL
// to which the image refers. The last directory entry is empty (filled with
// null values), which indicates the end of the directory table.
type ImageImportDescriptor struct {
	// The RVA of the import lookup/name table (INT). This table contains a name
	// or ordinal for each import. The INT is an array of IMAGE_THUNK_DATA structs.
	OriginalFirstThunk uint32

	// The stamp that is set to zero until the image is bound. After the image
	// is bound, this field is set to the time/data stamp of the DLL.
	TimeDateStamp uint32

	// The index of the first forwarder reference (-1 if no forwarders).
	ForwarderChain uint32

	// The address of an ASCII string that contains the name of the DLL.
	// This address is relative to the image base.
	Name uint32

	// The RVA of the import address table (IAT). The contents of this table are
	// identical to the contents of the import lookup table until the image is bound.
	FirstThunk uint32
}

// ImageThunkData32 corresponds to one imported function from the executable.
// The entries are an array of 32-bit numbers for PE32 or an array of 64-bit
// numbers for PE32+. The ends of both arrays are indicated by an
// IMAGE_THUNK_DATA element with a value of zero.
// The IMAGE_THUNK_DATA union is a DWORD with these interpretations:
// DWORD Function;       // Memory address of the imported function
// DWORD Ordinal;        // Ordinal value of imported API
// DWORD AddressOfData;  // RVA to an IMAGE_IMPORT_BY_NAME with the imported API name
// DWORD ForwarderString;// RVA to a forwarder string
type ImageThunkData32 struct {
	AddressOfData uint32
}

// ImageThunkData64 is the PE32+ version of IMAGE_THUNK_DATA.
type ImageThunkData64 struct {
	AddressOfData uint64
}

// ImportFunction represents an imported function in the import table.
type ImportFunction struct {
	// An ASCII string that contains the name to import. This is the string that
	// must be matched to the public name in the DLL. This string is case
	// sensitive and terminated by a null byte.
	Name string

	// An index into the export name pointer table. A match is attempted first
	// with this value. If it fails, a binary search is performed on the DLL's
	// export name pointer table.
	Hint uint16

	// If this is true, import by ordinal. Otherwise, import by name.
	ByOrdinal bool

	// A 16-bit ordinal number. This field is used only if the Ordinal/Name Flag
	// bit field is 1 (import by ordinal). Bits 30-15 or 62-15 must be 0.
	Ordinal uint32

	// Name Thunk Value (OFT)
	OriginalThunkValue   uint64

	// Address Thunk Value (FT)
	ThunkValue uint64
}

// Import represents an empty entry in the emport table.
type Import struct {
	Offset     uint32
	Name       string
	Functions  []*ImportFunction
	Descriptor ImageImportDescriptor
}

func (pe *File) parseImportDirectory(rva, size uint32) (err error) {

	for {
		importDesc := ImageImportDescriptor{}
		fileOffset := pe.getOffsetFromRva(rva)
		importDescSize := uint32(binary.Size(importDesc))
		buf := bytes.NewReader(pe.data[fileOffset : fileOffset+importDescSize])
		err := binary.Read(buf, binary.LittleEndian, &importDesc)

		// If the RVA is invalid all would blow up. Some EXEs seem to be
		// specially nasty and have an invalid RVA.
		if err != nil {
			return err
		}

		// If the structure is all zeros, we reached the end of the list.
		if importDesc == (ImageImportDescriptor{}) {
			break
		}

		rva += importDescSize

		// If the array of thunks is somewhere earlier than the import
		// descriptor we can set a maximum length for the array. Otherwise
		// just set a maximum length of the size of the file
		maxLen := uint32(len(pe.data)) - fileOffset
		if rva > importDesc.OriginalFirstThunk || rva > importDesc.FirstThunk {
			if rva < importDesc.OriginalFirstThunk {
				maxLen = rva - importDesc.FirstThunk
			} else if rva < importDesc.FirstThunk {
				maxLen = rva - importDesc.OriginalFirstThunk
			} else {
				maxLen = Max(rva-importDesc.OriginalFirstThunk,
					rva-importDesc.FirstThunk)
			}
		}

		var importedFunctions []*ImportFunction
		if pe.Is64 {
			importedFunctions, err = pe.parseImports64(&importDesc, maxLen)
		} else {
			importedFunctions, err = pe.parseImports32(&importDesc, maxLen)
		}
		if err != nil {
			return err
		}

		dllName := pe.getStringAtRVA(importDesc.Name, maxDllLength)
		if !IsValidDosFilename(dllName) {
			dllName = "*invalid*"
			continue
		}

		pe.Imports = append(pe.Imports, Import{
			Offset:     fileOffset,
			Name:       string(dllName),
			Functions:  importedFunctions,
			Descriptor: importDesc,
		})
	}

	return nil
}

func (pe *File) getImportTable32(rva uint32, maxLen uint32,
	isOldDelayImport bool) ([]*ImageThunkData32, error) {

	// Setup variables
	thunkTable := make(map[uint32]*ImageThunkData32)
	retVal := make([]*ImageThunkData32, 0)
	minAddressOfData := ^uint32(0)
	maxAddressOfData := uint32(0)
	repeatedAddress := uint32(0)
	var size uint32 = 4
	addressesOfData := make(map[uint32]bool)

	startRVA := rva

	if rva == 0 {
		return []*ImageThunkData32{}, nil
	}
	for {
		if rva >= startRVA+maxLen {
			log.Println("Error parsing the import table. Entries go beyond bounds.")
			break
		}

		// if we see too many times the same entry we assume it could be
		// a table containing bogus data (with malicious intent or otherwise)
		if repeatedAddress >= maxRepeatedAddresses {
			if !stringInSlice(AnoManyRepeatedEntries, pe.Anomalies) {
				pe.Anomalies = append(pe.Anomalies, AnoManyRepeatedEntries)
			}
		}

		// if the addresses point somewhere but the difference between the
		// highest and lowest address is larger than maxAddressSpread we assume
		// a bogus table as the addresses should be contained within a module
		if maxAddressOfData-minAddressOfData > maxAddressSpread {
			if !stringInSlice(AnoInvalidThunkAddressOfData, pe.Anomalies) {
				pe.Anomalies = append(pe.Anomalies, AnoInvalidThunkAddressOfData)
			}
		}

		// In its original incarnation in Visual C++ 6.0, all ImgDelayDescr
		// fields containing addresses used virtual addresses, rather than RVAs.
		// That is, they contained actual addresses where the delayload data
		// could be found. These fields are DWORDs, the size of a pointer on the x86.
		// Now fast-forward to IA-64 support. All of a sudden, 4 bytes isn't
		// enough to hold a complete address. At this point, Microsoft did the
		// correct thing and changed the fields containing addresses to RVAs.
		offset := uint32(0)
		if isOldDelayImport {
			oh32 := pe.NtHeader.OptionalHeader.(ImageOptionalHeader32)
			newRVA := rva - oh32.ImageBase
			offset = pe.getOffsetFromRva(newRVA)
			if offset == ^uint32(0) {
				return []*ImageThunkData32{}, nil
			}
		} else {
			offset = pe.getOffsetFromRva(rva)
			if offset == ^uint32(0) {
				return []*ImageThunkData32{}, nil
			}
		}

		// Read the image thunk data.
		thunk := ImageThunkData32{}
		buf := bytes.NewReader(pe.data[offset : offset+size])
		err := binary.Read(buf, binary.LittleEndian, &thunk)
		if err != nil {
			msg := fmt.Sprintf(`Error parsing the import table. 
				Invalid data at RVA: 0x%x`, rva)
			return []*ImageThunkData32{}, errors.New(msg)
		}

		if thunk == (ImageThunkData32{}) {
			break
		}

		// Check if the AddressOfData lies within the range of RVAs that it's
		// being scanned, abort if that is the case, as it is very unlikely
		// to be legitimate data.
		// Seen in PE with SHA256:
		// 5945bb6f0ac879ddf61b1c284f3b8d20c06b228e75ae4f571fa87f5b9512902c
		if thunk.AddressOfData >= startRVA && thunk.AddressOfData <= rva {
			log.Printf("Error parsing the import table. "+
				"AddressOfData overlaps with THUNK_DATA for THUNK at: "+
				"RVA 0x%x", rva)
			break
		}

		// If the entry looks like could be an ordinal.
		if thunk.AddressOfData&imageOrdinalFlag32 > 0 {
			// but its value is beyond 2^16, we will assume it's a
			// corrupted and ignore it altogether
			if thunk.AddressOfData&0x7fffffff > 0xffff {
				if !stringInSlice(AnoAddressOfDataBeyondLimits, pe.Anomalies) {
					pe.Anomalies = append(pe.Anomalies, AnoAddressOfDataBeyondLimits)
				}
			}
		} else {
			// and if it looks like it should be an RVA keep track of the RVAs seen
			// and store them to study their  properties. When certain non-standard
			// features are detected the parsing will be aborted
			_, ok := addressesOfData[thunk.AddressOfData]
			if ok {
				repeatedAddress++
			} else {
				addressesOfData[thunk.AddressOfData] = true
			}

			if thunk.AddressOfData > maxAddressOfData {
				maxAddressOfData = thunk.AddressOfData
			}

			if thunk.AddressOfData < minAddressOfData {
				minAddressOfData = thunk.AddressOfData
			}
		}

		thunkTable[rva] = &thunk
		retVal = append(retVal, &thunk)
		rva += size
	}
	return retVal, nil
}

func (pe *File) getImportTable64(rva uint32, maxLen uint32,
	isOldDelayImport bool) ([]*ImageThunkData64, error) {

	// Setup variables
	thunkTable := make(map[uint32]*ImageThunkData64)
	retVal := make([]*ImageThunkData64, 0)
	minAddressOfData := ^uint64(0)
	maxAddressOfData := uint64(0)
	repeatedAddress := uint64(0)
	var size uint32 = 8
	addressesOfData := make(map[uint64]bool)

	startRVA := rva
	for {
		if rva >= startRVA+maxLen {
			log.Println("Error parsing the import table. Entries go beyond bounds.")
			break
		}

		// if we see too many times the same entry we assume it could be
		// a table containing bogus data (with malicious intent or otherwise)
		if repeatedAddress >= uint64(maxRepeatedAddresses) {
			if !stringInSlice(AnoManyRepeatedEntries, pe.Anomalies) {
				pe.Anomalies = append(pe.Anomalies, AnoManyRepeatedEntries)
			}
		}

		// if the addresses point somewhere but the difference between the highest
		// and lowest address is larger than maxAddressSpread we assume a bogus
		// table as the addresses should be contained within a module
		if maxAddressOfData-minAddressOfData > uint64(maxAddressSpread) {
			if !stringInSlice(AnoInvalidThunkAddressOfData, pe.Anomalies) {
				pe.Anomalies = append(pe.Anomalies, AnoInvalidThunkAddressOfData)
			}
		}

		// In its original incarnation in Visual C++ 6.0, all ImgDelayDescr
		// fields containing addresses used virtual addresses, rather than RVAs.
		// That is, they contained actual addresses where the delayload data
		// could be found. These fields are DWORDs, the size of a pointer on the x86.
		// Now fast-forward to IA-64 support. All of a sudden, 4 bytes isn't
		// enough to hold a complete address. At this point, Microsoft did the
		// correct thing and changed the fields containing addresses to RVAs.
		offset := uint32(0)
		if isOldDelayImport {
			oh64 := pe.NtHeader.OptionalHeader.(ImageOptionalHeader64)
			newRVA := rva - uint32(oh64.ImageBase)
			offset = pe.getOffsetFromRva(newRVA)
			if offset == ^uint32(0) {
				return []*ImageThunkData64{}, nil
			}
		} else {
			offset = pe.getOffsetFromRva(rva)
			if offset == ^uint32(0) {
				return []*ImageThunkData64{}, nil
			}
		}

		// Read the image thunk data.
		thunk := ImageThunkData64{}
		buf := bytes.NewReader(pe.data[offset : offset+size])
		err := binary.Read(buf, binary.LittleEndian, &thunk)
		if err != nil {
			msg := fmt.Sprintf(`Error parsing the import table. 
				Invalid data at RVA: 0x%x`, rva)
			return []*ImageThunkData64{}, errors.New(msg)
		}

		if thunk == (ImageThunkData64{}) {
			break
		}

		// Check if the AddressOfData lies within the range of RVAs that it's
		// being scanned, abort if that is the case, as it is very unlikely
		// to be legitimate data.
		// Seen in PE with SHA256:
		// 5945bb6f0ac879ddf61b1c284f3b8d20c06b228e75ae4f571fa87f5b9512902c
		if thunk.AddressOfData >= uint64(startRVA) &&
			thunk.AddressOfData <= uint64(rva) {
			log.Printf("Error parsing the import table. "+
				"AddressOfData overlaps with THUNK_DATA for THUNK at: "+
				"RVA 0x%x", rva)
			break
		}

		// If the entry looks like could be an ordinal
		if thunk.AddressOfData&imageOrdinalFlag64 > 0 {
			// but its value is beyond 2^16, we will assume it's a
			// corrupted and ignore it altogether
			if thunk.AddressOfData&0x7fffffff > 0xffff {
				if !stringInSlice(AnoAddressOfDataBeyondLimits, pe.Anomalies) {
					pe.Anomalies = append(pe.Anomalies, AnoAddressOfDataBeyondLimits)
				}
			}
			// and if it looks like it should be an RVA
		} else {
			// keep track of the RVAs seen and store them to study their
			// properties. When certain non-standard features are detected
			// the parsing will be aborted
			_, ok := addressesOfData[thunk.AddressOfData]
			if ok {
				repeatedAddress++
			} else {
				addressesOfData[thunk.AddressOfData] = true
			}

			if thunk.AddressOfData > maxAddressOfData {
				maxAddressOfData = thunk.AddressOfData
			}

			if thunk.AddressOfData < minAddressOfData {
				minAddressOfData = thunk.AddressOfData
			}
		}

		thunkTable[rva] = &thunk
		retVal = append(retVal, &thunk)
		rva += size
	}
	return retVal, nil
}

func (pe *File) parseImports32(importDesc interface{}, maxLen uint32) (
	[]*ImportFunction, error) {

	var OriginalFirstThunk uint32
	var FirstThunk uint32
	var isOldDelayImport bool

	switch desc := importDesc.(type) {
	case *ImageImportDescriptor:
		OriginalFirstThunk = desc.OriginalFirstThunk
		FirstThunk = desc.FirstThunk
	case *ImageDelayImportDescriptor:
		OriginalFirstThunk = desc.ImportNameTableRVA
		FirstThunk = desc.ImportAddressTableRVA
		if desc.Attributes == 0 {
			isOldDelayImport = true
		}
	}

	// Import Lookup Table. Contains ordinals or pointers to strings.
	ilt, err := pe.getImportTable32(OriginalFirstThunk, maxLen, isOldDelayImport)
	if err != nil {
		return []*ImportFunction{}, err
	}

	// Import Address Table. May have identical content to ILT if PE file is
	// not bound. It will contain the address of the imported symbols once
	// the binary is loaded or if it is already bound.
	iat, err := pe.getImportTable32(FirstThunk, maxLen, isOldDelayImport)
	if err != nil {
		return []*ImportFunction{}, err
	}

	// Some DLLs has IAT or ILT with nil type.
	if len(iat) == 0 && len(ilt) == 0 {
		pe.Anomalies = append(pe.Anomalies, AnoNullImportTable)
		return []*ImportFunction{}, nil
	}

	var table []*ImageThunkData32
	if len(ilt) > 0 {
		table = ilt
	} else if len(iat) > 0 {
		table = iat
	} else {
		return []*ImportFunction{}, err
	}

	importedFunctions := make([]*ImportFunction, 0)
	numInvalid := uint32(0)
	for idx := uint32(0); idx < uint32(len(table)); idx++ {
		imp := ImportFunction{}
		if table[idx].AddressOfData > 0 {
			// If imported by ordinal, we will append the ordinal number
			if table[idx].AddressOfData&imageOrdinalFlag32 > 0 {
				imp.ByOrdinal = true
				imp.Ordinal = table[idx].AddressOfData & uint32(0xffff)
			} else {
				imp.ByOrdinal = false
				if isOldDelayImport {
					table[idx].AddressOfData -=
						pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).ImageBase
				}

				// Original Thunk
				if uint32(len(ilt)) != 0 {
					imp.OriginalThunkValue = uint64(ilt[idx].AddressOfData & addressMask32)
				}

				// Thunk
				if uint32(len(iat)) != 0 {
					imp.ThunkValue = uint64(iat[idx].AddressOfData & addressMask32)
				}

				// Thunk
				hintNameTableRva := table[idx].AddressOfData & addressMask32
				off := pe.getOffsetFromRva(hintNameTableRva)
				imp.Hint = binary.LittleEndian.Uint16(pe.data[off:])
				imp.Name = pe.getStringAtRVA(table[idx].AddressOfData+2,
					maxImportNameLength)
				if !IsValidFunctionName(imp.Name) {
					imp.Name = "*invalid*"
				}
			}
		}

		// This file bfe97192e8107d52dd7b4010d12b2924 has an invalid table built
		// in a way that it's parseable but contains invalid entries that lead
		// pefile to take extremely long amounts of time to parse. It also leads
		// to extreme memory consumption. To prevent similar cases, if invalid
		// entries are found in the middle of a table the parsing will be aborted.
		hasName := len(imp.Name) > 0
		if imp.Ordinal == 0 && !hasName {
			return []*ImportFunction{}, errors.New("Must have either an ordinal or a name in an import")
		}
		// Some PEs appear to interleave valid and invalid imports. Instead of
		// aborting the parsing altogether we will simply skip the invalid entries.
		// Although if we see 1000 invalid entries and no legit ones, we abort.
		if imp.Name == "*invalid*" {
			if numInvalid > 1000 && numInvalid == idx {
				return []*ImportFunction{}, errors.New("Too many invalid names, aborting parsing")
			}
			numInvalid++
			continue
		}

		if imp.Ordinal > 0 || hasName {
			importedFunctions = append(importedFunctions, &imp)
		}
	}

	return importedFunctions, nil
}

func (pe *File) parseImports64(importDesc interface{}, maxLen uint32) ([]*ImportFunction, error) {

	var OriginalFirstThunk uint32
	var FirstThunk uint32
	var isOldDelayImport bool

	switch desc := importDesc.(type) {
	case *ImageImportDescriptor:
		OriginalFirstThunk = desc.OriginalFirstThunk
		FirstThunk = desc.FirstThunk
	case *ImageDelayImportDescriptor:
		OriginalFirstThunk = desc.ImportNameTableRVA
		FirstThunk = desc.ImportAddressTableRVA
		if desc.Attributes == 0 {
			isOldDelayImport = true
		}
	}

	// Import Lookup Table. Contains ordinals or pointers to strings.
	ilt, err := pe.getImportTable64(OriginalFirstThunk, maxLen, isOldDelayImport)
	if err != nil {
		return []*ImportFunction{}, err
	}

	// Import Address Table. May have identical content to ILT if PE file is
	// not bound. It will contain the address of the imported symbols once
	// the binary is loaded or if it is already bound.
	iat, err := pe.getImportTable64(FirstThunk, maxLen, isOldDelayImport)
	if err != nil {
		return []*ImportFunction{}, err
	}

	// Would crash if IAT or ILT had nil type
	if len(iat) == 0 && len(ilt) == 0 {
		return []*ImportFunction{}, errors.New(
			"Damaged Import Table information. ILT and/or IAT appear to be broken")
	}

	var table []*ImageThunkData64
	if len(ilt) > 0 {
		table = ilt
	} else if len(iat) > 0 {
		table = iat
	} else {
		return []*ImportFunction{}, err
	}

	importedFunctions := make([]*ImportFunction, 0)
	numInvalid := uint32(0)
	for idx := uint32(0); idx < uint32(len(table)); idx++ {
		imp := ImportFunction{}
		if table[idx].AddressOfData > 0 {
			// If imported by ordinal, we will append the ordinal number
			if table[idx].AddressOfData&imageOrdinalFlag64 > 0 {
				imp.ByOrdinal = true
				imp.Ordinal = uint32(table[idx].AddressOfData) & uint32(0xffff)
			} else {
				imp.ByOrdinal = false
				if isOldDelayImport {
					table[idx].AddressOfData -=
						pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).ImageBase
				}

				// Original Thunk
				if uint32(len(ilt)) != 0 {
					imp.OriginalThunkValue = ilt[idx].AddressOfData & addressMask64
				}

				// Thunk
				if uint32(len(iat)) != 0 {
					imp.ThunkValue = iat[idx].AddressOfData & addressMask64
				}

				hintNameTableRva := table[idx].AddressOfData & addressMask64
				off := pe.getOffsetFromRva(uint32(hintNameTableRva))
				imp.Hint = binary.LittleEndian.Uint16(pe.data[off:])
				imp.Name = pe.getStringAtRVA(uint32(table[idx].AddressOfData+2),
					maxImportNameLength)
				if !IsValidFunctionName(imp.Name) {
					imp.Name = "*invalid*"
				}
			}
		}


		// This file bfe97192e8107d52dd7b4010d12b2924 has an invalid table built
		// in a way that it's parseable but contains invalid entries that lead
		// pefile to take extremely long amounts of time to parse. It also leads
		// to extreme memory consumption. To prevent similar cases, if invalid
		// entries are found in the middle of a table the parsing will be aborted.
		hasName := len(imp.Name) > 0
		if imp.Ordinal == 0 && !hasName {
			return []*ImportFunction{}, errors.New("Must have either an ordinal or a name in an import")
		}
		// Some PEs appear to interleave valid and invalid imports. Instead of
		// aborting the parsing altogether we will simply skip the invalid entries.
		// Although if we see 1000 invalid entries and no legit ones, we abort.
		if imp.Name == "*invalid*" {
			if numInvalid > 1000 && numInvalid == idx {
				return []*ImportFunction{}, errors.New("Too many invalid names, aborting parsing")
			}
			numInvalid++
			continue
		}

		if imp.Ordinal > 0 || hasName {
			importedFunctions = append(importedFunctions, &imp)
		}
	}

	return importedFunctions, nil
}

// md5hash hashes using md5 algorithm.
func md5hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// ImpHash calculates the import hash.
// Algorithm:
// ==========
// Resolving ordinals to function names when they appear
// Converting both DLL names and function names to all lowercase
// Removing the file extensions from imported module names
// Building and storing the lowercased string . in an ordered list
// Generating the MD5 hash of the ordered list
func (pe *File) ImpHash() (string, error) {
	if len(pe.Imports) == 0 {
		return "", errors.New("No imports found")
	}

	extensions := []string{"ocx", "sys", "dll"}
	var impStrs []string

	for _, imp := range pe.Imports {
		var libname string
		parts := strings.Split(imp.Name, ".")
		if len(parts) > 0 && stringInSlice(parts[1], extensions) {
			libname = parts[0]
		}

		libname = strings.ToLower(libname)

		for _, function := range imp.Functions {
			var funcname string
			if function.ByOrdinal {
				funcname = OrdLookup(libname, uint64(function.Ordinal), true)
			} else {
				funcname = function.Name
			}

			if funcname == "" {
				continue
			}

			impStr := fmt.Sprintf("%s.%s", libname, strings.ToLower(funcname))
			impStrs = append(impStrs, impStr)
		}
	}

	hash := md5hash(strings.Join(impStrs, ","))
	return hash, nil
}
