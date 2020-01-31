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

	// OptionalHeader is of type *OptionalHeader32 or *OptionalHeader64.
	OptionalHeader interface{}
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

// ImageOptionalHeader32 represents the PE32 format structure of the optional header.
// PE32 contains this additional field, which is absent in PE32+.
type ImageOptionalHeader32 struct {

	// The unsigned integer that identifies the state of the image file.
	// The most common number is 0x10B, which identifies it as a normal
	// executable file. 0x107 identifies it as a ROM image, and 0x20B identifies
	// it as a PE32+ executable.
	Magic uint16

	// The linker major version number.
	MajorLinkerVersion uint8

	// The linker minor version number.
	MinorLinkerVersion uint8

	// The size of the code (text) section, or the sum of all code sections
	// if there are multiple sections.
	SizeOfCode uint32

	// The size of the initialized data section, or the sum of all such
	// sections if there are multiple data sections.
	SizeOfInitializedData uint32

	// The size of the uninitialized data section (BSS), or the sum of all
	// such sections if there are multiple BSS sections.
	SizeOfUninitializedData uint32

	// The address of the entry point relative to the image base when the
	// executable file is loaded into memory. For program images, this is the
	// starting address. For device drivers, this is the address of the
	// initialization function. An entry point is optional for DLLs. When no
	// entry point is present, this field must be zero.
	AddressOfEntryPoint uint32

	// The address that is relative to the image base of the beginning-of-code
	// section when it is loaded into memory.
	BaseOfCode uint32

	// The address that is relative to the image base of the beginning-of-data
	// section when it is loaded into memory.
	BaseOfData uint32

	// The preferred address of the first byte of image when loaded into memory;
	// must be a multiple of 64 K. The default for DLLs is 0x10000000. The
	//default for Windows CE EXEs is 0x00010000. The default for Windows NT,
	// Windows 2000, Windows XP, Windows 95, Windows 98, and Windows Me is 0x00400000.
	ImageBase uint32

	// The alignment (in bytes) of sections when they are loaded into memory.
	// It must be greater than or equal to FileAlignment. The default is the
	// page size for the architecture.
	SectionAlignment uint32

	// The alignment factor (in bytes) that is used to align the raw data of
	// sections in the image file. The value should be a power of 2 between 512
	// and 64 K, inclusive. The default is 512. If the SectionAlignment is less
	// than the architecture's page size, then FileAlignment must match SectionAlignment.
	FileAlignment uint32

	// The major version number of the required operating system.
	MajorOperatingSystemVersion uint16

	// The minor version number of the required operating system.
	MinorOperatingSystemVersion uint16

	// The major version number of the image.
	MajorImageVersion uint16

	// The minor version number of the image.
	MinorImageVersion uint16

	// The major version number of the subsystem.
	MajorSubsystemVersion uint16

	// The minor version number of the subsystem.
	MinorSubsystemVersion uint16

	// Reserved, must be zero.
	Win32VersionValue uint32

	// The size (in bytes) of the image, including all headers, as the image
	// is loaded in memory. It must be a multiple of SectionAlignment.
	SizeOfImage uint32

	// The combined size of an MS-DOS stub, PE header, and section headers
	// rounded up to a multiple of FileAlignment.
	SizeOfHeaders uint32

	// The image file checksum. The algorithm for computing the checksum is
	// incorporated into IMAGHELP.DLL. The following are checked for validation
	// at load time: all drivers, any DLL loaded at boot time, and any DLL
	// that is loaded into a critical Windows process
	CheckSum uint32

	// The subsystem that is required to run this image.
	Subsystem uint16

	// For more information, see DLL Characteristics later in this specification.
	DllCharacteristics uint16

	// The size of the stack to reserve. Only SizeOfStackCommit is committed;
	// the rest is made available one page at a time until the reserve size is reached.
	SizeOfStackReserve uint32

	// The size of the stack to commit.
	SizeOfStackCommit uint32

	// The size of the local heap space to reserve. Only SizeOfHeapCommit is
	// committed; the rest is made available one page at a time until the
	// reserve size is reached.
	SizeOfHeapReserve uint32

	// The size of the local heap space to commit.
	SizeOfHeapCommit uint32

	// Reserved, must be zero.
	LoaderFlags uint32

	// The number of data-directory entries in the remainder of the optional
	// header. Each describes a location and size.
	NumberOfRvaAndSizes uint32

	// An array of 16 IMAGE_DATA_DIRECTORY structures.
	DataDirectory [16]DataDirectory
}

// ImageOptionalHeader64 represents the PE32+ format structure of the optional header.
type ImageOptionalHeader64 struct {
	// The unsigned integer that identifies the state of the image file.
	// The most common number is 0x10B, which identifies it as a normal
	// executable file. 0x107 identifies it as a ROM image, and 0x20B identifies
	// it as a PE32+ executable.
	Magic uint16

	// The linker major version number.
	MajorLinkerVersion uint8

	// The linker minor version number.
	MinorLinkerVersion uint8

	// The size of the code (text) section, or the sum of all code sections
	// if there are multiple sections.
	SizeOfCode uint32

	// The size of the initialized data section, or the sum of all such
	// sections if there are multiple data sections.
	SizeOfInitializedData uint32

	// The size of the uninitialized data section (BSS), or the sum of all
	// such sections if there are multiple BSS sections.
	SizeOfUninitializedData uint32

	// The address of the entry point relative to the image base when the
	// executable file is loaded into memory. For program images, this is the
	// starting address. For device drivers, this is the address of the
	// initialization function. An entry point is optional for DLLs. When no
	// entry point is present, this field must be zero.
	AddressOfEntryPoint uint32

	// The address that is relative to the image base of the beginning-of-code
	// section when it is loaded into memory.
	BaseOfCode uint32

	// In PE+, ImageBase is 8 bytes size.
	ImageBase uint64

	// The alignment (in bytes) of sections when they are loaded into memory.
	// It must be greater than or equal to FileAlignment. The default is the
	// page size for the architecture.
	SectionAlignment uint32

	// The alignment factor (in bytes) that is used to align the raw data of
	// sections in the image file. The value should be a power of 2 between 512
	// and 64 K, inclusive. The default is 512. If the SectionAlignment is less
	// than the architecture's page size, then FileAlignment must match SectionAlignment.
	FileAlignment uint32

	// The major version number of the required operating system.
	MajorOperatingSystemVersion uint16

	// The minor version number of the required operating system.
	MinorOperatingSystemVersion uint16

	// The major version number of the image.
	MajorImageVersion uint16

	// The minor version number of the image.
	MinorImageVersion uint16

	// The major version number of the subsystem.
	MajorSubsystemVersion uint16

	// The minor version number of the subsystem.
	MinorSubsystemVersion uint16

	// Reserved, must be zero.
	Win32VersionValue uint32

	// The size (in bytes) of the image, including all headers, as the image
	// is loaded in memory. It must be a multiple of SectionAlignment.
	SizeOfImage uint32

	// The combined size of an MS-DOS stub, PE header, and section headers
	// rounded up to a multiple of FileAlignment.
	SizeOfHeaders uint32

	// The image file checksum. The algorithm for computing the checksum is
	// incorporated into IMAGHELP.DLL. The following are checked for validation
	// at load time: all drivers, any DLL loaded at boot time, and any DLL
	// that is loaded into a critical Windows process
	CheckSum uint32

	// The subsystem that is required to run this image.
	Subsystem uint16

	// For more information, see DLL Characteristics later in this specification.
	DllCharacteristics uint16

	// The size of the stack to reserve. Only SizeOfStackCommit is committed;
	// the rest is made available one page at a time until the reserve size is reached.
	SizeOfStackReserve uint64

	// The size of the stack to commit.
	SizeOfStackCommit uint64

	// The size of the local heap space to reserve. Only SizeOfHeapCommit is
	// committed; the rest is made available one page at a time until the
	// reserve size is reached.
	SizeOfHeapReserve uint64

	// The size of the local heap space to commit.
	SizeOfHeapCommit uint64

	// Reserved, must be zero.
	LoaderFlags uint32

	// The number of data-directory entries in the remainder of the optional
	// header. Each describes a location and size.
	NumberOfRvaAndSizes uint32

	// An array of 16 IMAGE_DATA_DIRECTORY structures.
	DataDirectory [16]DataDirectory
}

// DataDirectory represents an array of 16 IMAGE_DATA_DIRECTORY structures,
// 8 bytes apiece, each relating to an important data structure in the PE file.
type DataDirectory struct {
	VirtualAddress uint32 // The RVA of the data structure.
	Size           uint32 // The size in bytes of the data structure refered to.
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
	pe.NtHeader.Signature = signature

	// The file header structure contains some basic information about the file;
	// most importantly, a field describing the size of the optional data that
	// follows it. In PE files, this optional data is very much required, but is
	// still called the IMAGE_OPTIONAL_HEADER.
	size := uint32(binary.Size(pe.NtHeader.FileHeader))
	fileHeaderOffset := ntHeaderOffset + 4
	buf := bytes.NewReader(pe.data[fileHeaderOffset : fileHeaderOffset+size])
	err = binary.Read(buf, binary.LittleEndian, &pe.NtHeader.FileHeader)
	if err != nil {
		return err
	}

	// The optional header could be either for a PE or PE+ file.
	// Do not include the interface of optionheader in the Size(),
	// Otherwise, it won't work.
	oh32 := ImageOptionalHeader32{}
	oh64 := ImageOptionalHeader64{}

	optHeaderOffset := ntHeaderOffset + uint32(binary.Size(pe.NtHeader.FileHeader) + 4)
	magic := binary.LittleEndian.Uint16(pe.data[optHeaderOffset:])

	// Probes for PE32/PE32+ optional header magic.
	if magic != ImageNtOptionalHeader32Magic && magic != ImageNtOptionalHeader64Magic {
		return ErrImageNtOptionalHeaderMagicNotFound
	}

	// Are we dealing with a PE64 optional header.
	switch magic {
	case ImageNtOptionalHeader64Magic:
		size := uint32(binary.Size(oh64))
		buf := bytes.NewReader(pe.data[optHeaderOffset : optHeaderOffset+size])
		err := binary.Read(buf, binary.LittleEndian, &oh64)
		if err != nil {
			return err
		}
		pe.Is64 = true
		pe.NtHeader.OptionalHeader = oh64
	case ImageNtOptionalHeader32Magic:
		size := uint32(binary.Size(oh32))
		buf := bytes.NewReader(pe.data[optHeaderOffset : optHeaderOffset+size])
		err := binary.Read(buf, binary.LittleEndian, &oh32)
		if err != nil {
			return err
		}
		pe.Is32 = true
		pe.NtHeader.OptionalHeader = oh32
	}


	// ImageBase should be multiple of 10000h.
	if pe.Is64 && oh64.ImageBase%0x10000 != 0 {
		return ErrImageBaseNotAligned
	}
	if pe.Is32 && oh32.ImageBase%0x10000 != 0 {
		return ErrImageBaseNotAligned
	}

	// ImageBase can be any value as long as:
	// ImageBase + SizeOfImage < 80000000h for PE32.
	if pe.Is32 &&
		oh32.ImageBase+oh32.SizeOfImage >= 0x80000000 {
		return ErrImageBaseOverflow
	}

	// SizeOfImage must be a multiple of the section alignment.
	if (pe.Is32 && oh32.SizeOfImage%oh32.SectionAlignment != 0) || 
		(pe.Is64 && oh64.SizeOfImage%oh64.SectionAlignment != 0){
		return ErrInvalidSizeOfImage
	}

	return nil
}