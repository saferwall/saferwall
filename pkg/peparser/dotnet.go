// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

// COM+ Header entry point flags.
const (
	COMImageFlagsILOnly           = 0x00000001
	COMImageFlags32BitRequired    = 0x00000002
	COMImageFlagILLibrary         = 0x00000004
	COMImageFlagsStrongNameSigned = 0x00000008
	COMImageFlagsNativeEntrypoint = 0x00000010
	COMImageFlagsTrackDebugData   = 0x00010000
	COMImageFlags32BitPreferred   = 0x00020000
)

// ImageDataDirectory represents the  directory format.
type ImageDataDirectory struct {
	VirtualAddress uint32 // The relative virtual address of the table.
	Size           uint32 // The size of the table, in bytes.
}

// ImageCOR20Header represents the CLR 2.0 header structure.
type ImageCOR20Header struct {
	// Header versioning
	Cb                  uint32 // Size of this structure (0x48)
	MajorRuntimeVersion uint16 //Major version of the CLR runtime
	MinorRuntimeVersio  uint16 // Minor version of the CLR runtime

	// Symbol table and startup information

	// RVA to, and size of, the executables meta-data.
	MetaData ImageDataDirectory

	// Bitwise flags indicating attributes of this executable.
	Flags uint32

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
	// meta data has a table that maps names to offsets into this blob, so
	// logically the blob is a set of resources.
	Resources ImageDataDirectory

	// IL assemblies can be signed with a public-private key to validate who
	// created it. The signature goes here if this feature is used.
	StrongNameSignature ImageDataDirectory

	CodeManagerTable ImageDataDirectory // Depricated, not used

	// Used for manged codee that has unmaanaged code inside it
	// (or exports methods as unmanaged entry points)
	VTableFixups            ImageDataDirectory
	ExportAddressTableJumps ImageDataDirectory

	// null for ordinary IL images.
	// NGEN images it points at a code:CORCOMPILE_HEADER structure
	ManagedNativeHeader ImageDataDirectory
}

func (pe *File) parseCLRHeaderDirectory(rva, size uint32) error {

	clrHeader := ImageCOR20Header{}
	offset := pe.getOffsetFromRva(rva)
	err := pe.structUnpack(&clrHeader, offset, size)
	if err != nil {
		return err
	}

	pe.CLRHeader = &clrHeader
	return nil
}
