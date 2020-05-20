// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
)

// ImageDosHeader represents the DOS stub of a PE.
type ImageDosHeader struct {
	// Magic number.
	Magic uint16

	// Bytes on last page of file.
	BytesOnLastPageOfFile uint16

	// Pages in file.
	PagesInFile uint16

	// Relocations.
	Relocations uint16

	// Size of header in paragraphs.
	SizeOfHeader uint16

	// Minimum extra paragraphs needed.
	MinExtraParagraphsNeeded uint16

	// Maximum extra paragraphs needed.
	MaxExtraParagraphsNeeded uint16

	// Initial (relative) SS value.
	InitialSS uint16

	// Initial SP value.
	InitialSP uint16

	// Checksum.
	Checksum uint16

	// Initial IP value.
	InitialIP uint16

	// Initial (relative) CS value.
	InitialCS uint16

	// File address of relocation table.
	AddressOfRelocationTable uint16

	// Overlay number.
	OverlayNumber uint16

	// Reserved words.
	ReservedWords1 [4]uint16

	// OEM identifier.
	OEMIdentifier uint16

	// OEM information.
	OEMInformation uint16

	// Reserved words.
	ReservedWords2 [10]uint16

	// File address of new exe header (Elfanew).
	AddressOfNewEXEHeader uint32
}

// ParseDOSHeader parses the DOS header stub. Every PE file begins with a small
// MS-DOS stub. The need for this arose in the early days of Windows, before a
// significant number of consumers were running it. When executed on a machine
// without Windows, the program could at least print out amessage saying that
// Windows was required to run the executable.
func (pe *File) ParseDOSHeader() (err error) {
	offset := 0
	size := binary.Size(pe.DosHeader)
	buf := bytes.NewReader(pe.data[offset : offset+size])
	err = binary.Read(buf, binary.LittleEndian, &pe.DosHeader)
	if err != nil {
		return err
	}

	// It can be ZM on an (non-PE) EXE.
	// These executables still work under XP via ntvdm.
	if pe.DosHeader.Magic != ImageDOSSignature &&
		pe.DosHeader.Magic != ImageDOSZMSignature {
		return ErrDOSMagicNotFound
	}

	// `e_lfanew` is the only required element (besides the signature) of the
	// DOS header to turn the EXE into a PE. It is is a relative offset to the
	// NT Headers. It can't be null (signatures would overlap).
	// Can be 4 at minimum.
	if pe.DosHeader.AddressOfNewEXEHeader < 4 ||
		pe.DosHeader.AddressOfNewEXEHeader > pe.size {
		return ErrInvalidElfanewValue
	}

	// tiny pe has a e_lfanew of 4, which means the NT Headers is overlapping
	// the DOS Header.
	if pe.DosHeader.AddressOfNewEXEHeader <= 0x3c {
		pe.Anomalies = append(pe.Anomalies, AnoPEHeaderOverlapDOSHeader)
	}

	return nil
}
