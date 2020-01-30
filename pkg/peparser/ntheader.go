// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
)

// ImageNtHeader represents the PE header and is the general term for a structure named IMAGE_NT_HEADERS.
type ImageNtHeader struct {
	// Signature is a DWORD containing the value 50h, 45h, 00h, 00h.
	Signature uint32

	// IMAGE_NT_HEADERS includes the IMAGE_FILE_HEADER structure.
	FileHeader ImageFileHeader
}

// ImageFileHeader contains info about the physical layout and properties of the file.
type ImageFileHeader struct {
	// The number that identifies the type of target machine.
	Machine uint16

	// The number of sections. This indicates the size of the section table,
	// which immediately follows the headers.
	NumberOfSections uint16

	// // The low 32 bits of the number of seconds since 00:00 January 1, 1970
	// (a C run-time time_t value), that indicates when the file was created.
	TimeDateStamp uint32

	// // The file offset of the COFF symbol table, or zero if no COFF symbol
	// table is present. This value should be zero for an image because COFF
	// debugging information is deprecated.
	PointerToSymbolTable uint32

	// The number of entries in the symbol table. This data can be used to
	// locate the string table, which immediately follows the symbol table.
	// This value should be zero for an image because COFF debugging information
	// is deprecated.
	NumberOfSymbols uint32

	// The size of the optional header, which is required for executable files
	// but not for object files. This value should be zero for an object file.
	SizeOfOptionalHeader uint16

	// The flags that indicate the attributes of the file.
	Characteristics uint16
}

// The IMAGE_NT_HEADERS structure is the primary location where specifics of
// the PE file are stored. Its offset is given by the e_lfanew field in the
// IMAGE_DOS_HEADER at the beginning of the file.
func (pe *File) parseNtHeader() (err error) {
	ntHeaderOffset := pe.DosHeader.Elfanew
	signature := binary.LittleEndian.Uint32(pe.data[ntHeaderOffset:])

	// Probe for PE signature.
	if signature == ImageOS2Signature {
		return ErrImageOS2SignatureFound
	}
	if signature == ImageOS2LESignature {
		return ErrImageOS2LESignatureFound
	}
	if signature == ImageVXDignature {
		return ErrImageVXDSignatureFound
	}
	if signature == ImageTESignature {
		return ErrImageTESignatureFound
	}

	// This is the smallest requirement for a valid PE.
	if signature != ImageNTSignature {
		return ErrImageNtSignatureNotFound
	}

	// The file header structure contains some basic information about the file;
	// most importantly, a field describing the size of the optional data that
	// follows it. In PE files, this optional data is very much required, but is
	// still called the IMAGE_OPTIONAL_HEADER.
	size := uint32(binary.Size(pe.NtHeader))
	buf := bytes.NewReader(pe.data[ntHeaderOffset : ntHeaderOffset+size])
	err = binary.Read(buf, binary.LittleEndian, &pe.NtHeader)
	if err != nil {
		return err
	}

	return nil
}
