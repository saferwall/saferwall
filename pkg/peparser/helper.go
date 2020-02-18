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
	ErrInvalidElfanewValue = errors.New("Invalid e_lfanew value, probably not a PE file")

	// ErrImageOS2SignatureFound is returned when signature is for a NE file.
	ErrImageOS2SignatureFound = errors.New("Not a valid PE signature. Probably a NE file")

	// ErrImageOS2LESignatureFound is returned when signature is for a LE file.
	ErrImageOS2LESignatureFound = errors.New("Not a valid PE signature. Probably an LE file")

	// ErrImageVXDSignatureFound is returned when signature is for a NX file.
	ErrImageVXDSignatureFound = errors.New("Not a valid PE signature. Probably an LX file")

	// ErrImageTESignatureFound is returned when signature is for a TE file.
	ErrImageTESignatureFound = errors.New("Not a valid PE signature. Probably a TE file")

	// ErrImageNtSignatureNotFound is returned when PE magic signature is not found.
	ErrImageNtSignatureNotFound = errors.New("Not a valid PE signature. Magic not found")

	// ErrImageNtOptionalHeaderMagicNotFound is returned when optional header
	// magic is different from PE32/PE32+.
	ErrImageNtOptionalHeaderMagicNotFound = errors.New(`Not a valid PE signature.
	 Optional Header magic not found`)

	// ErrImageBaseNotAligned is reported when the image base is not aligned to 64 K.
	ErrImageBaseNotAligned = errors.New("Corrupt PE file. Image base not aligned to 64 K")

	// ErrImageBaseOverflow is reported when the image base is larger than
	// 80000000h/FFFF080000000000h in PE32/P32+.
	ErrImageBaseOverflow = errors.New("Corrupt PE file. Image base is overflow")

	// ErrInvalidSectionFileAlignment is reported when section alignment is less than a
	// PAGE_SIZE and section alignement != file alignment.
	ErrInvalidSectionFileAlignment = errors.New(`Corrupt PE file. Section alignment
	 is less than a PAGE_SIZE and section alignement != file alignment`)

	// ErrInvalidSizeOfImage is reported when SizeOfImage is nota multiple of
	// SectionAlignment.
	ErrInvalidSizeOfImage = errors.New(`Invalid SizeOfImage value, should be
	 multiple of SectionAlignment`)
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
		log.Println("data at RVA can't be fetched. Corrupt header?")
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

// PrettyMachineType returns the string representations
// of the `Machine` field of  the IMAGE_FILE_HEADER.
func (pe *File) PrettyMachineType() string {
	machineType := map[uint16]string{
		ImageFileMachineUnknown:   "Unknown",
		ImageFileMachineAM33:      "Matsushita AM33",
		ImageFileMachineAMD64:     "x64",
		ImageFileMachineARM:       "ARM little endian",
		ImageFileMachineARM64:     "ARM64 little endian",
		ImageFileMachineARMNT:     "ARM Thumb-2 little endian",
		ImageFileMachineEBC:       "EFI byte code",
		ImageFileMachineI386:      "Intel 386 or later / compatible processors",
		ImageFileMachineIA64:      "Intel Itanium processor family",
		ImageFileMachineM32R:      "Mitsubishi M32R little endian",
		ImageFileMachineMIPS16:    "MIPS16",
		ImageFileMachineMIPSFPU:   "MIPS with FPU",
		ImageFileMachineMIPSFPU16: "MIPS16 with FPU",
		ImageFileMachinePowerPC:   "Power PC little endian",
		ImageFileMachinePowerPCFP: "Power PC with floating point support",
		ImageFileMachineR4000:     "MIPS little endian",
		ImageFileMachineRISCV32:   "RISC-V 32-bit address space",
		ImageFileMachineRISCV64:   "RISC-V 64-bit address space",
		ImageFileMachineRISCV128:  "RISC-V 128-bit address space",
		ImageFileMachineSH3:       "Hitachi SH3",
		ImageFileMachineSH3DSP:    "Hitachi SH3 DSP",
		ImageFileMachineSH4:       "Hitachi SH4",
		ImageFileMachineSH5:       "Hitachi SH5",
		ImageFileMachineTHUMB:     "Thumb",
		ImageFileMachineWCEMIPSV2: "MIPS little-endian WCE v2",
	}

	return machineType[pe.NtHeader.FileHeader.Machine]
}

// PrettyImageFileCharacteristics returns the string representations
// of the `Characteristics` field of  the IMAGE_FILE_HEADER.
func (pe *File) PrettyImageFileCharacteristics() []string {
	var values []string
	fileHeaderCharacteristics := map[uint16]string{
		ImageFileRelocsStripped:       "RelocsStripped",
		ImageFileExecutableImage:      "ExecutableImage",
		ImageFileLineNumsStripped:     "LineNumsStripped",
		ImageFileLocalSymsStripped:    "LocalSymsStripped",
		ImageFileAgressibeWsTrim:      "AgressibeWsTrim",
		ImageFileLargeAddressAware:    "LargeAddressAware",
		ImageFileBytesReservedLow:     "BytesReservedLow",
		ImageFile32BitMachine:         "32BitMachine",
		ImageFileDebugStripped:        "DebugStripped",
		ImageFileRemovableRunFromSwap: "RemovableRunFromSwap",
		ImageFileSystem:               "FileSystem",
		ImageFileDLL:                  "DLL",
		ImageFileUpSystemOnly:         "UpSystemOnly",
		ImageFileBytesReservedHigh:    "BytesReservedHigh",
	}

	for k, s := range fileHeaderCharacteristics {
		if k&pe.NtHeader.FileHeader.Characteristics != 0 {
			values = append(values, s)
		}
	}
	return values
}

// PrettyDllCharacteristics returns the string representations
// of the `DllCharacteristics` field of ImageOptionalHeader.
func (pe *File) PrettyDllCharacteristics() []string {
	var values []string
	var characteristics uint16

	if pe.Is64 {
		characteristics =
			pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).DllCharacteristics
	} else {
		characteristics =
			pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).DllCharacteristics
	}

	imgDllCharacteristics := map[uint16]string{
		ImageDllCharacteristicsHighEntropyVA:        "HighEntropyVA",
		ImageDllCharacteristicsDynamicBase:          "DynamicBase",
		ImageDllCharacteristicsForceIntegrity:       "ForceIntegrity",
		ImageDllCharacteristicsNXCompact:            "NXCompact",
		ImageDllCharacteristicsNoIsolation:          "NoIsolation",
		ImageDllCharacteristicsNoSEH:                "NoSEH",
		ImageDllCharacteristicsNoBind:               "NoBind",
		ImageDllCharacteristicsAppContainer:         "AppContainer",
		ImageDllCharacteristicsWdmDriver:            "WdmDriver",
		ImageDllCharacteristicsGuardCF:              "GuardCF",
		ImageDllCharacteristicsTerminalServiceAware: "TerminalServiceAware",
	}

	for k, s := range imgDllCharacteristics {
		if k&characteristics != 0 {
			values = append(values, s)
		}
	}

	return values
}

// PrettySectionFlags returns the string representations of the
// `Flags` field of section header.
func (pe *File) PrettySectionFlags(curSectionFlag uint32) []string {
	var values []string

	sectionFlags := map[uint32]string{
		ImageScnReserved1:            "Reserved1",
		ImageScnReserved2:            "Reserved2",
		ImageScnReserved3:            "Reserved3",
		ImageScnReserved4:            "Reserved4",
		ImageScnTypeNoPad:            "No Padd",
		ImageScnReserved5:            "Reserved5",
		ImageScnCntCode:              "Contains Code",
		ImageScnCntInitializedData:   "Initialized Data",
		ImageScnCntUninitializedData: "Uninitialized Data",
		ImageScnLnkOther:             "Lnk Other",
		ImageScnLnkInfo:              "Lnk Info",
		ImageScnReserved6:            "Reserved6",
		ImageScnLnkRemove:            "LnkRemove",
		ImageScnLnkComdat:            "LnkComdat",
		ImageScnGpRel:                "GpReferenced",
		ImageScnMemPurgeable:         "Purgeable",
		ImageScnMemLocked:            "Locked",
		ImageScnMemPreload:           "Preload",
		ImageScnAlign1Bytes:          "Align1Bytes",
		ImageScnAlign2Bytes:          "Align2Bytes",
		ImageScnAlign4Bytes:          "Align4Bytes",
		ImageScnAlign8Bytes:          "Align8Bytes",
		ImageScnAlign16Bytes:         "Align16Bytes",
		ImageScnAlign32Bytes:         "Align32Bytes",
		ImageScnAlign64Bytes:         "Align64Bytes",
		ImageScnAlign128Bytes:        "Align128Bytes",
		ImageScnAlign256Bytes:        "Align265Bytes",
		ImageScnAlign512Bytes:        "Align512Bytes",
		ImageScnAlign1024Bytes:       "Align1024Bytes",
		ImageScnAlign2048Bytes:       "Align2048Bytes",
		ImageScnAlign4096Bytes:       "Align4096Bytes",
		ImageScnAlign8192Bytes:       "Align8192Bytes",
		ImageScnLnkMRelocOvfl:        "ExtendedReloc",
		ImageScnMemDiscardable:       "Discardable",
		ImageScnMemNotCached:         "NotCached",
		ImageScnMemNotPaged:          "NotPaged",
		ImageScnMemShared:            "Shared",
		ImageScnMemExecute:           "Executable",
		ImageScnMemRead:              "Readable",
		ImageScnMemWrite:             "Writable",
	}

	for k, s := range sectionFlags {
		if k&curSectionFlag != 0 {
			values = append(values, s)
		}
	}

	return values
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
	optionalHeaderOffset := pe.DosHeader.Elfanew + uint32(binary.Size(pe.NtHeader))

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
