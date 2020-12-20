// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

// COM+ Header entry point flags.
const (
	// The image file contains IL code only, with no embedded native unmanaged
	// code except the start-up stub (which simply executes an indirect jump to
	// the CLR entry point).
	COMImageFlagsILOnly = 0x00000001

	// The image file can be loaded only into a 32-bit process.
	COMImageFlags32BitRequired = 0x00000002

	// This flag is obsolete and should not be set. Setting it—as the IL
	// assembler allows, using the .corflags directive—will render your module
	// unloadable.
	COMImageFlagILLibrary = 0x00000004

	// The image file is protected with a strong name signature.
	COMImageFlagsStrongNameSigned = 0x00000008

	// The executable’s entry point is an unmanaged method. The EntryPointToken/
	// EntryPointRVA field of the CLR header contains the RVA of this native
	// method. This flag was introduced in version 2.0 of the CLR.
	COMImageFlagsNativeEntrypoint = 0x00000010

	// The CLR loader and the JIT compiler are required to track debug
	// information about the methods. This flag is not used.
	COMImageFlagsTrackDebugData = 0x00010000

	// The image file can be loaded into any process, but preferably into a
	// 32-bit process. This flag can be only set together with flag
	// COMIMAGE_FLAGS_32BITREQUIRED. When set, these two flags mean the image
	// is platformneutral, but prefers to be loaded as 32-bit when possible.
	// This flag was introduced in CLR v4.0
	COMImageFlags32BitPreferred = 0x00020000
)

// V-table constants.
const (
	// V-table slots are 32-bits in size.
	CORVTable32Bit = 0x01

	// V-table slots are 64-bits in size.
	CORVTable64Bit = 0x02

	//  The thunk created by the common language runtime must provide data
	// marshaling between managed and unmanaged code.
	CORVTableFromUnmanaged = 0x04

	// The thunk created by the common language runtime must provide data
	// marshaling between managed and unmanaged code. Current appdomain should
	// be selected to dispatch the call.
	CORVTableFromUnmanagedRetainAppDomain = 0x08

	// Call most derived method described by
	CORVTableCallMostDerived = 0x10
)

// ImageDataDirectory represents the  directory format.
type ImageDataDirectory struct {

	// The relative virtual address of the table.
	VirtualAddress uint32

	// The size of the table, in bytes.
	Size uint32
}

// ImageCOR20Header represents the CLR 2.0 header structure.
type ImageCOR20Header struct {

	// Size of the header in bytes.
	Cb uint32

	// Major number of the minimum version of the runtime required to run the
	// program.
	MajorRuntimeVersion uint16

	// Minor number of the version of the runtime required to run the program.
	MinorRuntimeVersio uint16

	// RVA and size of the metadata.
	MetaData ImageDataDirectory

	// Bitwise flags indicating attributes of this executable.
	Flags uint32

	// Metadata identifier (token) of the entry point for the image file; can
	// be 0 for DLL images. This field identifies a method belonging to this
	// module or a module containing the entry point method.
	// In images of version 2.0 and newer, this field may contain RVA of the
	// embedded native entry point method.
	// union {
	//
	// If COMIMAGE_FLAGS_NATIVE_ENTRYPOINT is not set,
	// EntryPointToken represents a managed entrypoint.
	//	DWORD               EntryPointToken;
	//
	// If COMIMAGE_FLAGS_NATIVE_ENTRYPOINT is set,
	// EntryPointRVA represents an RVA to a native entrypoint
	//	DWORD               EntryPointRVA;
	//};
	EntryPointRVAorToken uint32

	// This is the blob of managed resources. Fetched using
	// code:AssemblyNative.GetResource and code:PEFile.GetResource and accessible
	// from managed code from System.Assembly.GetManifestResourceStream. The
	// metadata has a table that maps names to offsets into this blob, so
	// logically the blob is a set of resources.
	Resources ImageDataDirectory

	// RVA and size of the hash data for this PE file, used by the loader for
	// binding and versioning. IL assemblies can be signed with a public-private
	// key to validate who created it. The signature goes here if this feature
	// is used.
	StrongNameSignature ImageDataDirectory

	// RVA and size of the Code Manager table. In the existing releases of the
	// runtime, this field is reserved and must be set to 0.
	CodeManagerTable ImageDataDirectory

	// RVA and size in bytes of an array of virtual table (v-table) fixups.
	// Among current managed compilers, only the VC++ linker and the IL
	// assembler can produce this array.
	VTableFixups ImageDataDirectory

	// RVA and size of an array of addresses of jump thunks. Among managed
	// compilers, only the VC++ of versions pre-8.0 could produce this table,
	// which allows the export of unmanaged native methods embedded in the
	// managed PE file. In v2.0+ of CLR this entry is obsolete and must be set
	// to 0.
	ExportAddressTableJumps ImageDataDirectory

	// Reserved for precompiled images; set to 0
	// NGEN images it points at a code:CORCOMPILE_HEADER structure
	ManagedNativeHeader ImageDataDirectory
}

// ImageCORVTableFixup defines the v-table fixups that contains the
// initializing information necessary for the runtime to create the thunks.
// Non VOS v-table entries.  Define an array of these pointed to by
// IMAGE_COR20_HEADER.VTableFixups.  Each entry describes a contiguous array of
// v-table slots.  The slots start out initialized to the meta data token value
// for the method they need to call.  At image load time, the CLR Loader will
// turn each entry into a pointer to machine code for the CPU and can be
// called directly.
type ImageCORVTableFixup struct {
	RVA   uint32 // Offset of v-table array in image.
	Count uint16 // How many entries at location.
	Type  uint16 // COR_VTABLE_xxx type of entries.
}

// MetadataHeader consists of a storage signature and a storage header.
type MetadataHeader struct {

	// The storage signature, which must be 4-byte aligned:
	// ====================================================

	// ”Magic” signature for physical metadata, currently 0x424A5342, or, read
	// as characters, BSJB—the initials of four “founding fathers” Brian Harry,
	// Susa Radke-Sproull, Jason Zander, and Bill Evans, who started the
	// runtime development in 1998.
	Signature uint32

	// Major version.
	MajorVersion uint16

	// Minor version.
	MinorVersion uint16

	// Reserved; set to 0.
	ExtraData uint32

	// Length of the version string.
	VersionString uint32

	// Version string.
	Version string

	// The storage header follows the storage signature, aligned on a 4-byte
	// boundary.
	// ====================================================================

	// Reserved; set to 0.
	Flags uint8

	// Another byte used for [padding]

	// Number of streams.
	Streams uint16
}

// MetadataStreamHeader represents a Metadata Stream Header Structure.
type MetadataStreamHeader struct {
	// Offset Offset in the file for this stream.
	Offset uint32

	// Size of the stream in bytes.
	Size uint32

	// Name of the stream; a zero-terminated ASCII string no longer than 31
	// characters (plus zero terminator). The name might be shorter, in which
	// case the size of the stream header is correspondingly reduced, padded to
	// the 4-byte boundary.
	Name string
}

// CLRData embeds the Common Language Runtime Header structure as well as the
// Metadata header structure.
type CLRData struct {
	CLRHeader             *ImageCOR20Header
	MetadataHeader        *MetadataHeader
	MetadataStreamHeaders []*MetadataStreamHeader
}

func (pe *File) parseMetadataHeader(rva, size uint32, clr *CLRData) error {
	var err error
	mh := MetadataHeader{}

	offset := pe.getOffsetFromRva(rva)
	if mh.Signature, err = pe.ReadUint32(offset); err != nil {
		return err
	}
	if mh.MajorVersion, err = pe.ReadUint16(offset + 4); err != nil {
		return err
	}
	if mh.MinorVersion, err = pe.ReadUint16(offset + 6); err != nil {
		return err
	}
	if mh.ExtraData, err = pe.ReadUint32(offset + 8); err != nil {
		return err
	}
	if mh.VersionString, err = pe.ReadUint32(offset + 12); err != nil {
		return err
	}
	mh.Version, err = pe.getStringAtOffset(offset+16, mh.VersionString)
	if err != nil {
		return err
	}

	offset += 16 + mh.VersionString
	if mh.Flags, err = pe.ReadUint8(offset); err != nil {
		return err
	}

	if mh.Streams, err = pe.ReadUint16(offset + 2); err != nil {
		return err
	}

	clr.MetadataHeader = &mh

	// Immediately following the MetadataHeader is a series of Stream Headers.
	// A “stream” is to the metadata what a “section” is to the assembly. The
	// NumberOfStreams property indicates how many StreamHeaders to read.
	offset += 4
	for i := uint16(0); i < mh.Streams; i++ {
		sh := MetadataStreamHeader{}
		if sh.Offset, err = pe.ReadUint32(offset); err != nil {
			return err
		}
		if sh.Size, err = pe.ReadUint32(offset + 4); err != nil {
			return err
		}

		// Name requires a special treatement.
		offset += 8
		for j := uint32(0); j <= 32; j++ {
			var c uint8
			if c, err = pe.ReadUint8(offset); err != nil {
				return err
			}

			offset++
			if c == 0 && (j+1) % 4 == 0 {
				break
			}
			if c != 0 {
				sh.Name += string(c)
			}
		}

		clr.MetadataStreamHeaders = append(clr.MetadataStreamHeaders, &sh)

	}
	return nil
}

// The 15th directory entry of the PE header contains the RVA and size of the
// runtime header in the image file. The runtime header, which contains all of
// the runtime-specific data entries and other information, should reside in a
// read-only section of the image file. The IL assembler puts the common
// language runtime header in the .text section.
func (pe *File) parseCLRHeaderDirectory(rva, size uint32) error {

	clr := CLRData{}

	clrHeader := ImageCOR20Header{}
	offset := pe.getOffsetFromRva(rva)
	err := pe.structUnpack(&clrHeader, offset, size)
	if err != nil {
		return err
	}
	clr.CLRHeader = &clrHeader
	if clrHeader.MetaData.VirtualAddress != 0 && clrHeader.MetaData.Size != 0 {
		pe.parseMetadataHeader(clrHeader.MetaData.VirtualAddress,
			clrHeader.MetaData.Size, &clr)
	}

	pe.CLR = &clr
	return nil
}
