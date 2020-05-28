// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"encoding/binary"
	"errors"
	"log"
	"strings"
)

const (
	// TinyPESize On Windows XP (x32) the smallest PE executable is 97 bytes.
	TinyPESize = 97

	// FileAlignmentHardcodedValue represents the value which PointerToRawData
	// should be at least equal or bigger to, or it will be rounded to zero.
	// According to http://corkami.blogspot.com/2010/01/parce-que-la-planche-aura-brule.html
	// if PointerToRawData is less that 0x200 it's rounded to zero.
	FileAlignmentHardcodedValue = 0x200
)

// Errors
var (

	// ErrInvalidPESize is returned when the file size is less that the smallest
	// PE file size possible.ErrImageOS2SignatureFound
	ErrInvalidPESize = errors.New("Not a PE file, smaller than tiny PE")

	// ErrDOSMagicNotFound is returned when file is potentially a ZM executable.
	ErrDOSMagicNotFound = errors.New("DOS Header magic not found")

	// ErrInvalidElfanewValue is returned when e_lfanew is larger than file size.
	ErrInvalidElfanewValue = errors.New(
		"Invalid e_lfanew value. Probably not a PE file")

	// ErrImageOS2SignatureFound is returned when signature is for a NE file.
	ErrImageOS2SignatureFound = errors.New(
		"Not a valid PE signature. Probably a NE file")

	// ErrImageOS2LESignatureFound is returned when signature is for a LE file.
	ErrImageOS2LESignatureFound = errors.New(
		"Not a valid PE signature. Probably an LE file")

	// ErrImageVXDSignatureFound is returned when signature is for a LX file.
	ErrImageVXDSignatureFound = errors.New(
		"Not a valid PE signature. Probably an LX file")

	// ErrImageTESignatureFound is returned when signature is for a TE file.
	ErrImageTESignatureFound = errors.New(
		"Not a valid PE signature. Probably a TE file")

	// ErrImageNtSignatureNotFound is returned when PE magic signature is not found.
	ErrImageNtSignatureNotFound = errors.New(
		"Not a valid PE signature. Magic not found")

	// ErrImageNtOptionalHeaderMagicNotFound is returned when optional header
	// magic is different from PE32/PE32+.
	ErrImageNtOptionalHeaderMagicNotFound = errors.New(
		"Not a valid PE signature. Optional Header magic not found")

	// ErrImageBaseNotAligned is reported when the image base is not aligned to 64K.
	ErrImageBaseNotAligned = errors.New(
		"Corrupt PE file. Image base not aligned to 64 K")

	// AnoImageBaseOverflow is reported when the image base + SizeOfImage is
	// larger than 80000000h/FFFF080000000000h in PE32/P32+.
	AnoImageBaseOverflow = "Image base beyong allowed address"

	// ErrInvalidSectionFileAlignment is reported when section alignment is less than a
	// PAGE_SIZE and section alignement != file alignment.
	ErrInvalidSectionFileAlignment = errors.New("Corrupt PE file. Section " + 
	"alignment is less than a PAGE_SIZE and section alignement != file alignment")

	// AnoInvalidSizeOfImage is reported when SizeOfImage is not multiple of
	// SectionAlignment.
	AnoInvalidSizeOfImage = "Invalid SizeOfImage value, should be multiple " +
	 "of SectionAlignment"
)

// Max returns the larger of x or y.
func Max(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// Min returns the min number in a slice.
func Min(values []uint32) uint32 {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// IsValidDosFilename returns true if the DLL name is likely to be valid.
// Valid FAT32 8.3 short filename characters according to:
// http://en.wikipedia.org/wiki/8.3_filename
// The filename length is not checked because the DLLs filename
// can be longer that the 8.3
func IsValidDosFilename(filename string) bool {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numerals := "0123456789"
	special := "!#$%&'()-@^_`{}~+,.;=[]\\/"
	charset := alphabet + numerals + special
	for _, c := range filename {
		if !strings.Contains(charset, string(c)) {
			return false
		}
	}
	return true
}

// IsValidFunctionName checks if an imported name uses the valid accepted
// characters expected in mangled function names. If the symbol's characters
// don't fall within this charset we will assume the name is invalid.
func IsValidFunctionName(functionName string) bool {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numerals := "0123456789"
	special := "_?@$()<>"
	charset := alphabet + numerals + special
	for _, c := range charset {
		if !strings.Contains(charset, string(c)) {
			return false
		}
	}
	return true
}

// IsPrintable checks weather a string is printable.
func IsPrintable(s string) bool {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numerals := "0123456789"
	whitespace := " \t\n\r\v\f"
	special := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	charset := alphabet + numerals + special + whitespace
	for _, c := range charset {
		if !strings.Contains(charset, string(c)) {
			return false
		}
	}
	return true
}

// getSectionByRva returns the section containing the given address.
func (pe *File) getSectionByRva(rva uint32) *ImageSectionHeader {
	for _, section := range pe.Sections {
		if section.Contains(rva, pe) {
			return &section
		}
	}
	return nil
}

func (pe *File) getSectionByOffset(offset uint32) *ImageSectionHeader {
	for _, section := range pe.Sections {
		if section.PointerToRawData == 0 {
			continue
		}

		adjustedPointer := pe.adjustFileAlignment(section.PointerToRawData)
		if adjustedPointer <= offset &&
			offset < (adjustedPointer+section.SizeOfRawData) {
			return &section
		}
	}
	return nil
}

// getOffsetFromRva returns the file offset corresponding to this RVA.
func (pe *File) getOffsetFromRva(rva uint32) uint32 {

	// Given a RVA, this method will find the section where the
	// data lies and return the offset within the file.
	section := pe.getSectionByRva(rva)
	if section == nil {
		if rva < uint32(len(pe.data)) {
			return rva
		}
		log.Println("Data at RVA can't be fetched. Corrupt header?")
		return ^uint32(0)
	}
	sectionAlignment := pe.adjustSectionAlignment(section.VirtualAddress)
	fileAlignment := pe.adjustFileAlignment(section.PointerToRawData)
	return rva - sectionAlignment + fileAlignment
}

func (pe *File) getRvaFromOffset(offset uint32) uint32 {
	section := pe.getSectionByOffset(offset)
	minAddr := ^uint32(0)
	if section == nil {

		if len(pe.Sections) == 0 {
			return offset
		}

		for _, section := range pe.Sections {
			vaddr := pe.adjustSectionAlignment(section.VirtualAddress)
			if vaddr < minAddr {
				minAddr = vaddr
			}
		}
		// Assume that offset lies within the headers
		// The case illustrating this behavior can be found at:
		// http://corkami.blogspot.com/2010/01/hey-hey-hey-whats-in-your-head.html
		// where the import table is not contained by any section
		// hence the RVA needs to be resolved to a raw offset
		if offset < minAddr {
			return offset
		}

		log.Println("data at Offset can't be fetched. Corrupt header?")
		return ^uint32(0)
	}
	sectionAlignment := pe.adjustSectionAlignment(section.VirtualAddress)
	fileAlignment := pe.adjustFileAlignment(section.PointerToRawData)
	return offset - fileAlignment + sectionAlignment
}

// getStringAtRVA returns an ASCII string located at the given address.
func (pe *File) getStringAtRVA(rva, maxLen uint32) string {
	if rva == 0 {
		return ""
	}

	section := pe.getSectionByRva(rva)
	if section == nil {
		s := pe.getStringFromData(0, []byte(pe.data[rva:rva+maxLen]))
		return string(s)
	}
	s := pe.getStringFromData(0, section.Data(rva, maxLen, pe))
	return string(s)
}

// getStringFromData returns ASCII string from within the data.
func (pe *File) getStringFromData(offset uint32, data []byte) []byte {
	if offset > uint32(len(data)) {
		return nil
	}

	end := offset
	for end < uint32(len(data)) {
		if data[end] == 0 {
			break
		}
		end++
	}
	return data[offset:end]
}

// getData returns the data given an RVA regardless of the section where it lies on.
func (pe *File) getData(rva, length uint32) ([]byte, error) {

	// Given a RVA and the size of the chunk to retrieve, this method
	// will find the section where the data lies and return the data.
	section := pe.getSectionByRva(rva)

	var end uint32
	if length > 0 {
		end = rva + length
	} else {
		end = 0
	}

	if section == nil {
		if rva < uint32(len(pe.Header)) {
			return pe.Header[rva:end], nil
		}

		// Before we give up we check whether the file might
		// contain the data anyway. There are cases of PE files
		// without sections that rely on windows loading the first
		// 8291 bytes into memory and assume the data will be there
		// A functional file with these characteristics is:
		// MD5: 0008892cdfbc3bda5ce047c565e52295
		// SHA-1: c7116b9ff950f86af256defb95b5d4859d4752a9

		if rva < uint32(len(pe.data)) {
			return pe.data[rva:end], nil
		}

		return nil, errors.New("data at RVA can't be fetched. Corrupt header?")
	}
	return section.Data(rva, length, pe), nil
}

// The alignment factor (in bytes) that is used to align the raw data of sections
// in the image file. The value should be a power of 2 between 512 and 64 K,
// inclusive. The default is 512. If the SectionAlignment is less than the
// architecture's page size, then FileAlignment must match SectionAlignment.
func (pe *File) adjustFileAlignment(va uint32) uint32 {

	var fileAlignment uint32
	switch pe.Is64 {
	case true:
		fileAlignment = pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).FileAlignment
	case false:
		fileAlignment = pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).FileAlignment
	}

	if fileAlignment > FileAlignmentHardcodedValue && fileAlignment%2 != 0 {
		pe.Anomalies = append(pe.Anomalies, ErrInvalidFileAlignment)
	}

	if fileAlignment < FileAlignmentHardcodedValue {
		return va
	}

	// round it to 0x200 if not power of 2.
	return (va / 0x200) * 0x200

}

// The alignment (in bytes) of sections when they are loaded into memory
// It must be greater than or equal to FileAlignment. The default is the
// page size for the architecture.
func (pe *File) adjustSectionAlignment(va uint32) uint32 {
	var fileAlignment, sectionAlignment uint32

	switch pe.Is64 {
	case true:
		fileAlignment = pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).FileAlignment
		sectionAlignment = pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).FileAlignment
	case false:
		fileAlignment = pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).SectionAlignment
		sectionAlignment = pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).SectionAlignment
	}

	if fileAlignment < FileAlignmentHardcodedValue &&
		fileAlignment != sectionAlignment {
		pe.Anomalies = append(pe.Anomalies, ErrInvalidSectionAlignment)
	}

	// 0x200 is the minimum valid FileAlignment according to the documentation
	// although ntoskrnl.exe has an alignment of 0x80 in some Windows versions

	if sectionAlignment != 0 && va%sectionAlignment != 0 {
		return sectionAlignment * (va / sectionAlignment)
	}
	return va
}

// stringInSlice checks weather a string exists in a slice of strings.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// intInSlice checks weather a uint32 exists in a slice of uint32.
func intInSlice(a uint32, list []uint32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IsDriver returns true if the PE file is a Windows driver.
func (pe *File)IsDriver() bool {

	// Checking that the ImageBase field of the OptionalHeader is above or
	// equal to 0x80000000 (that is, whether it lies in the upper 2GB of
	//the address space, normally belonging to the kernel) is not a
	// reliable enough indicator.  For instance, PEs that play the invalid
	// ImageBase trick to get relocated could be incorrectly assumed to be
	// drivers.

	// Checking if any section characteristics have the IMAGE_SCN_MEM_NOT_PAGED 
	// flag set is not reliable either.

	// If the import directory was not parsed (fast_load = True); do it now.
	if len(pe.Imports) == 0 {
		pe.ParseDataDirectories()
	}

	// If there's still no import directory (the PE doesn't have one or it's
	// malformed), give up.
	if len(pe.Imports) == 0 {
		return false
	}

	// DIRECTORY_ENTRY_IMPORT will now exist, although it may be empty.
	// If it imports from "ntoskrnl.exe" or other kernel components it should
	// be a driver.
	systemDLLs := []string {"ntoskrnl.exe", "hal.dll", "ndis.sys",
	 "bootvid.dll", "kdcom.dll"}
	for _, dll := range pe.Imports {
		if stringInSlice(strings.ToLower(dll.Name), systemDLLs) {
			return true
		}
	}

	// If still we couldn't tell, check common driver section with combinaison 
	// of IMAGE_SUBSYSTEM_NATIVE or IMAGE_SUBSYSTEM_NATIVE_WINDOWS.
	subsystem := uint16(0)
	oh32 := ImageOptionalHeader32{}
	oh64 := ImageOptionalHeader64{}
	switch pe.Is64 {
	case true:
		oh64 = pe.NtHeader.OptionalHeader.(ImageOptionalHeader64)
		subsystem = oh64.Subsystem
	case false:
		oh32 = pe.NtHeader.OptionalHeader.(ImageOptionalHeader32)
		subsystem = oh32.Subsystem
	}
	commonDriverSectionNames := []string {"page", "paged", "nonpage", "init"}
	for _, section := range pe.Sections {
		s := strings.ToLower(section.NameString())
		s = strings.Replace(s, "\x00", "", -1)
		if stringInSlice(s, commonDriverSectionNames) && 
			(subsystem & ImageSubsystemNativeWindows != 0 || 
				subsystem & ImageSubsystemNative != 0 ) { 
			return true
		}

	}

	return false
}

// IsDLL returns true if the PE file is a standard DLL.
func (pe *File)IsDLL() bool {
	if pe.NtHeader.FileHeader.Characteristics & ImageFileDLL != 0 {
		return true
	}
	return false
}


// IsEXE returns true if the PE file is a standard executable.
func (pe *File)IsEXE() bool {

	// Returns true only if the file has the IMAGE_FILE_EXECUTABLE_IMAGE flag set
	// and the IMAGE_FILE_DLL not set and the file does not appear to be a driver either.
	if pe.IsDLL() || pe.IsDriver() {
		return false
	}

	if pe.NtHeader.FileHeader.Characteristics & ImageFileExecutableImage == 0 {
		return false
	}

	return true
}

// padOrTrim returns (size) bytes from input (bb)
// Short bb gets zeros prefixed, Long bb gets left/MSB bits trimmed
func padOrTrim(bb []byte, size int) []byte {
	l := len(bb)
	if l == size {
		return bb
	}
	if l > size {
		return bb[l-size:]
	}
	tmp := make([]byte, size)
	copy(tmp[size-l:], bb)
	return tmp
}

// Checksum calculates the PE checksum as generated by CheckSumMappedFile().
func (pe *File) Checksum() uint32 {
	checksum := 0
	max := 0x100000000
	currentDword := uint32(0)

	// Get the Checksum offset.
	optionalHeaderOffset := pe.DosHeader.AddressOfNewEXEHeader + uint32(binary.Size(pe.NtHeader))

	// `CheckSum` field position in optional PE headers is always 64 for PE and PE+.
	checksumOffset := optionalHeaderOffset + 64

	// Verify the data is DWORD-aligned and add padding if needed
	remainder := pe.size % 4
	dataLen := pe.size
	if remainder > 0 {
		dataLen = pe.size + (4 - remainder)
	}

	for i := uint32(0); i < dataLen/4; i++ {
		// Skip the checksum field
		if i*4 == checksumOffset {
			continue
		}

		// Did we reach the last dword ?
		if i+1 == dataLen/4 && remainder > 0 {
			bb := pe.data[i*4 : i*4+(4-remainder)]
			lastDword := padOrTrim(bb, 4)
			currentDword = binary.LittleEndian.Uint32(lastDword)
		} else {
			currentDword = binary.LittleEndian.Uint32(pe.data[i*4:])
		}

		checksum += int(currentDword)
		if checksum > max {
			checksum = (checksum & 0xffffffff) + (checksum >> 32)
		}
	}

	checksum = (checksum & 0xffff) + (checksum >> 16)
	checksum = checksum + (checksum >> 16)
	checksum = checksum & 0xffff

	// The length is the one of the original data, not the padded one
	checksum += int(pe.size)

	return uint32(checksum)
}

// ReadUint32 read a uint32 from a buffer.
func ReadUint32 (buff []byte, offset uint32) (uint32, error) {
	if offset > uint32(len(buff)) {
		return binary.LittleEndian.Uint32(buff[offset:]), nil
	}
	return 0, errors.New("")
}