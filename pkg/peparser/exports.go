// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

const (
	maxExportedSymbols = 0x2000
)

// ImageExportDirectory represents the IMAGE_EXPORT_DIRECTORY structure.
// The export directory table contains address information that is used
// to resolve imports to the entry points within this image.
type ImageExportDirectory struct {
	// Reserved, must be 0.
	Characteristics uint32

	// The time and date that the export data was created.
	TimeDateStamp uint32

	// The major version number.
	//The major and minor version numbers can be set by the user.
	MajorVersion uint16

	// The minor version number.
	MinorVersion uint16

	// The address of the ASCII string that contains the name of the DLL.
	// This address is relative to the image base.
	Name uint32

	// The starting ordinal number for exports in this image. This field
	// specifies the starting ordinal number for the export address table.
	// It is usually set to 1.
	Base uint32

	// The number of entries in the export address table.
	NumberOfFunctions uint32

	// The number of entries in the name pointer table.
	// This is also the number of entries in the ordinal table.
	NumberOfNames uint32

	//The address of the export address table, relative to the image base.
	AddressOfFunctions uint32

	// The address of the export name pointer table, relative to the image base.
	// The table size is given by the Number of Name Pointers field.
	AddressOfNames uint32

	// The address of the ordinal table, relative to the image base.
	AddressOfNameOrdinals uint32
}

// ExportFunction represents an imported function in the export table.
type ExportFunction struct {
	Ordinal      uint32
	FunctionRVA  uint32
	NameOrdinal  uint32
	NameRVA      uint32
	Name         string
	Forwarder    string
	ForwarderRVA uint32
}

func (pe *File) parseExportDirectory(rva, size uint32) error {

	// Define some vars.
	exportDir := ImageExportDirectory{}
	errorMsg := fmt.Sprintf("Error parsing export directory at RVA: 0x%x", rva)

	fileOffset := pe.getOffsetFromRva(rva)
	buf := bytes.NewReader(pe.data[fileOffset : fileOffset+size])
	err := binary.Read(buf, binary.LittleEndian, &exportDir)
	if err != nil {
		return errors.New(errorMsg)
	}

	// We keep track of the bytes left in the file and use it to set a upper
	// bound in the number of items that can be read from the different arrays.
	lengthUntilEOF := func(rva uint32) uint32 {
		return pe.size - pe.getOffsetFromRva(rva)
	}
	var length uint32
	var addressOfNames []byte

	if exportDir.NumberOfFunctions == 0 {
		log.Println("Export Directory counts zero number of functions")
		return nil
	}

	length = min(lengthUntilEOF(exportDir.AddressOfNames),
		exportDir.NumberOfNames*4)
	addressOfNames, err = pe.getData(exportDir.AddressOfNames, length)
	if err != nil {
		return errors.New(errorMsg)
	}

	length = min(lengthUntilEOF(exportDir.AddressOfNameOrdinals),
		exportDir.NumberOfNames*4)
	addressOfNameOrdinals, err := pe.getData(exportDir.AddressOfNameOrdinals, length)

	if err != nil {
		return errors.New(errorMsg)
	}

	length = min(lengthUntilEOF(exportDir.AddressOfFunctions),
		exportDir.NumberOfFunctions*4)
	addressOfFunctions, err := pe.getData(exportDir.AddressOfFunctions, length)
	if err != nil {
		return errors.New(errorMsg)
	}

	maxFailedEntries := 10
	var forwarderStr string
	var forwarderOffset uint32
	safetyBoundary := pe.size // overly generous upper bound
	symbolCounts := make(map[uint32]int)
	parsingFailed := false

	// read the image export directory
	section := pe.getSectionByRva(exportDir.AddressOfNames)
	if section != nil {
		safetyBoundary = (section.VirtualAddress +
			uint32(len(section.Data(0, 0, pe)))) - exportDir.AddressOfNames
	}

	numNames := min(exportDir.NumberOfNames, safetyBoundary/4)
	var symbolAddress uint32
	for i := uint32(0); i < numNames; i++ {

		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			if recover() != nil {
				err = errors.New("array index out of bounds")
			}
		}()

		symbolOrdinal := binary.LittleEndian.Uint16(addressOfNameOrdinals[i*2:])
		if symbolOrdinal*4 < uint16(len(addressOfFunctions)) {
			symbolAddress = binary.LittleEndian.Uint32(addressOfFunctions[symbolOrdinal*4:])
		} else {
			return errors.New(`Corrupt? a bad pointer... we assume it's
			 all useless, no exports`)
		}

		if symbolAddress == 0 {
			continue
		}

		// If the function's RVA points within the export directory
		// it will point to a string with the forwarded symbol's string
		// instead of pointing the the function start address.
		if symbolAddress >= rva && symbolAddress < rva+size {
			forwarderStr = pe.getStringAtRVA(symbolAddress, 0x100000)
			forwarderOffset = pe.getOffsetFromRva(symbolAddress)
		} else {
			forwarderStr = ""
			fileOffset = 0
		}

		symbolNameAddress := binary.LittleEndian.Uint32(addressOfNames[i*4:])
		if symbolNameAddress == 0 {
			maxFailedEntries--
			if maxFailedEntries <= 0 {
				parsingFailed = true
				break
			}
		}
		symbolName := pe.getStringAtRVA(symbolNameAddress, 0x100000)
		if !IsValidFunctionName(symbolName) {
			parsingFailed = true
			break
		}

		symbolNameOffset := pe.getOffsetFromRva(symbolNameAddress)
		if symbolNameOffset == 0 {
			maxFailedEntries--
			if maxFailedEntries <= 0 {
				parsingFailed = true
				break
			}
		}

		// File 0b1d3d3664915577ab9a32188d29bbf3542b86c7b9ce333e245496c3018819f1
		// was being parsed as potentially containing millions of exports.
		// Checking for duplicates addresses the issue.
		symbolCounts[symbolAddress]++
		if symbolCounts[symbolAddress] > 0x200 {
			log.Printf(`Export directory contains more than 0x200 repeated
			 entries: %d, Name:%s, Address:0x%x`, symbolCounts[symbolAddress],
				string(symbolName), symbolAddress)
			break
		}
		if len(symbolCounts) > maxExportedSymbols {
			log.Printf(`Export directory contains more than %d symbol entries.
			 Assuming corrupt.`, maxExportedSymbols)
			break
		}
		newExport := ExportFunction{
			Name:         symbolName,
			NameRVA:      symbolNameAddress,
			NameOrdinal:  uint32(symbolOrdinal),
			Ordinal:      exportDir.Base + uint32(symbolOrdinal),
			FunctionRVA:  symbolAddress,
			Forwarder:    forwarderStr,
			ForwarderRVA: forwarderOffset,
		}

		pe.Exports = append(pe.Exports, newExport)
	}

	if parsingFailed {
		log.Printf(`RVA AddressOfNames in the export directory points to an
		 invalid address: 0x%x`, exportDir.AddressOfNames)
	}

	maxFailedEntries = 10
	section = pe.getSectionByRva(exportDir.AddressOfFunctions)

	// Overly generous upper bound
	safetyBoundary = pe.size
	if section != nil {
		safetyBoundary = section.VirtualAddress +
			uint32(len(section.Data(0, 0, pe))) - exportDir.AddressOfNames
	}
	parsingFailed = false
	ordinals := make(map[uint32]bool)
	for _, exp := range pe.Exports {
		ordinals[exp.Ordinal] = true
	}
	numNames = min(exportDir.NumberOfFunctions, safetyBoundary/4)
	for i := uint32(0); i < numNames; i++ {
		value := i + exportDir.Base
		if ordinals[value] {
			continue
		}

		symbolAddress = binary.LittleEndian.Uint32(addressOfFunctions[i*4:])
		if symbolAddress == 0 {
			continue
		}

		// Checking for forwarder again.
		if symbolAddress >= rva && symbolAddress < rva+size {
			forwarderStr = pe.getStringAtRVA(symbolAddress, 0x100000)
			forwarderOffset = pe.getOffsetFromRva(symbolAddress)
		} else {
			forwarderStr = ""
			fileOffset = 0
		}

		// File 0b1d3d3664915577ab9a32188d29bbf3542b86c7b9ce333e245496c3018819f1
		// was being parsed as potentially containing millions of exports.
		// Checking for duplicates addresses the issue.
		symbolCounts[symbolAddress]++
		if symbolCounts[symbolAddress] > 0x200 {
			log.Printf(`Export directory contains more than 0x200 repeated
			 entries: %d, Address:0x%x`, symbolCounts[symbolAddress],
				symbolAddress)
			break
		}
		if len(symbolCounts) > maxExportedSymbols {
			log.Printf(`Export directory contains more than %d ordinal entries.
			 Assuming corrupt.`, maxExportedSymbols)
			break
		}
		newExport := ExportFunction{
			Ordinal:      exportDir.Base + i,
			FunctionRVA:  symbolAddress,
			Forwarder:    forwarderStr,
			ForwarderRVA: forwarderOffset,
		}

		pe.Exports = append(pe.Exports, newExport)

	}

	return nil
}
