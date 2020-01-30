package pe

import (
	"bytes"
	"encoding/binary"
)

// ImageDosHeader represents the DOS stub of a PE.
type ImageDosHeader struct {
	// Magic number
	Emagic uint16 `json:"magic"`

	// Bytes on last page of file
	Ecblp uint16 `json:"bytesOnLastPageOfFile"`

	// Pages in file
	Ecp uint16 `json:"pagesInFile"`

	// Relocations
	Ecrlc uint16 `json:"relocations"`

	// Size of header in paragraphs
	Ecparhdr uint16 `json:"headerSizeInParagraphs"`

	// Minimum extra paragraphs needed
	Eminalloc uint16 `json:"minExtraParagraphsNeeded"`

	// Maximum extra paragraphs needed
	Emaxalloc uint16 `json:"maxExtraParagraphsNeeded"`

	// Initial (relative) SS value
	Ess uint16 `json:"initialSS"`

	// Initial SP value
	Esp uint16 `json:"initialSP"`

	// Checksum
	Ecsum uint16 `json:"checksum"`

	// Initial IP value
	Eip uint16 `json:"initialIP"`

	// Initial (relative) CS value
	Ecs uint16 `json:"initialCS"`

	// File address of relocation table
	Elfarlc uint16 `json:"fileAddressOfRelocationTable"`

	// Overlay number
	Eovno uint16 `json:"overlayNumber"`

	// Reserved words
	Eres [4]uint16 `json:"reservedWords4"`

	// OEM identifier (for e_oeminfo)
	Eoemid uint16 `json:"OEMIdentifier"`

	// OEM information; e_oemid specific
	Eoeminfo uint16 `json:"OEMInformation"`

	// Reserved words
	Eres2 [10]uint16 `json:"reservedWords10"`

	// File address of new exe header
	Elfanew uint32 `json:"fileAddressOfNewEXEHeader"`
}

// Every PE file begins with a small MS-DOS stud. The need for this arose in the
// early days of Windows, before a significant number of consumers were running
// it. When executed on a machine without Windows, the program could at least
// print out amessage saying that Windows was required to run the executable.
func (pe *File) parseDosHeader() (err error) {
	offset := 0
	size := binary.Size(pe.DosHeader)
	buf := bytes.NewReader(pe.data[offset : offset+size])
	err = binary.Read(buf, binary.LittleEndian, &pe.DosHeader)
	if err != nil {
		return err
	}

	// it can be ZM on an (non-PE) EXE.
	// These executables still work under XP via ntvdm.
	if pe.DosHeader.Emagic != ImageDOSSignature &&
		pe.DosHeader.Emagic != ImageDOSZMSignature {
		return ErrDOSMagicNotFound
	}

	// `e_lfanew` is the only required element (besides the signature) of the
	// DOS header to turn the EXE into a PE. It is is a relative offset to the
	// NT Headers. It can't be null (signatures would overlap).
	// Can be 4 at minimum.
	if pe.DosHeader.Elfanew < 4 ||
		pe.DosHeader.Elfanew > uint32(len(pe.data)) {
		return ErrInvalidElfanewValue
	}

	return nil
}
