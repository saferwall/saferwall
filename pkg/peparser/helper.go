package pe

import (
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

	// ErrInvalidPESize is returned when the file size is less that the smallest PE file size possible.ErrImageOS2SignatureFound
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

	// ErrImageNtOptionalHeaderMagicNotFound is returned when optional header magic is different from PE32/PE32+.
	ErrImageNtOptionalHeaderMagicNotFound = errors.New("Not a valid PE signature. Optional Header magic not found")

	// ErrImageBaseNotAligned is reported when the image base is not aligned to 64 K.
	ErrImageBaseNotAligned = errors.New("Corrupt PE file. Image base not aligned to 64 K")

	// ErrImageBaseOverflow is reported when the image base is larger than 80000000h/FFFF080000000000h in PE32/P32+.
	ErrImageBaseOverflow = errors.New("Corrupt PE file. Image base is overflow")

	// ErrInvalidMajorSubsystemVersion is reported when major subsystem version is less than 3.
	ErrInvalidMajorSubsystemVersion = errors.New("Corrupt PE file. Optional header major subsystem version is less than 3")

	// ErrInvalidSectionFileAlignment is reported when section alignment is less than a PAGE_SIZE and section alignement != file alignment.
	ErrInvalidSectionFileAlignment = errors.New("Corrupt PE file. Ssection alignment is less than a PAGE_SIZE and section alignement != file alignment")
)


// Max returns the larger of x or y.
func Max(x, y uint32) uint32 {
    if x < y {
        return y
    }
    return x
}

// Min returns the min number in a slice.
func Min(values []uint32) uint32 {
	min := values[0]
	for _, v := range values {
			if (v < min) {
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
func IsValidDosFilename (filename string) bool {
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
func IsValidFunctionName (functionName string) bool {
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
		if adjustedPointer <= offset && offset < (adjustedPointer + section.SizeOfRawData) {
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
		if offset < minAddr { return offset }

		log.Println("data at Offset can't be fetched. Corrupt header?")
		return ^uint32(0)
	}
	sectionAlignment := pe.adjustSectionAlignment(section.VirtualAddress)
	fileAlignment := pe.adjustFileAlignment(section.PointerToRawData)
	return offset - fileAlignment + sectionAlignment
}


// getStringAtRVA returns an ASCII string located at the given address.
func (pe *File) getStringAtRVA(rva uint32) string {
	if rva == 0 {
		return ""
	}

	section := pe.getSectionByRva(rva)
	if section == nil {
		s :=  pe.getStringFromData(0, []byte(pe.data[rva:rva+0x100000]))
		return string(s)
	}
	s :=  pe.getStringFromData(0, section.Data(rva, uint32(0x100000), pe))
	return string(s)
}


// getStringFromData; Get an ASCII string from within the data.
func (pe *File) getStringFromData(offset uint32, data []byte) ([]byte) {
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

	fileAlignment := pe.OptionalHeader.FileAlignment
	if fileAlignment > FileAlignmentHardcodedValue && fileAlignment % 2 != 0 {
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
	sectionAlignment := pe.OptionalHeader.SectionAlignment
	fileAlignment := pe.OptionalHeader.FileAlignment
	if fileAlignment < FileAlignmentHardcodedValue &&
		fileAlignment != sectionAlignment {
			pe.Anomalies = append(pe.Anomalies, ErrInvalidSectionAlignment) 
	}
	
	// 0x200 is the minimum valid FileAlignment according to the documentation
	// although ntoskrnl.exe has an alignment of 0x80 in some Windows versions

	if sectionAlignment != 0 && va % sectionAlignment != 0 {
		return sectionAlignment * ( va / sectionAlignment )
	}
	return va
}
