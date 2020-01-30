package pe

// Image executable types
const (
	ImageDOSSignature   = 0x5A4D     // MZ
	ImageDOSZMSignature = 0x4D5A     // ZM
	ImageOS2Signature   = 0x454E     // The New Executable (abbreviated NE or NewEXE) is a 16-bit .exe file format, a successor to the DOS MZ executable format. It was used in Windows 1.0–3.x, multitasking MS-DOS 4.0, OS/2 1.x, and the OS/2 subset of Windows NT up to version 5.0 (Windows 2000). A NE is also called a segmented executable.[2]
	ImageOS2LESignature = 0x454C     // Linear Executable is an executable file format in the EXE family. It was used by 32-bit OS/2, by some DOS extenders, and by Microsoft Windows VxD files. It is an extension of MS-DOS EXE, and a successor to NE (New Executable).
	ImageVXDignature    = 0x584C     // There are two main varieties of it: LX (32-bit), and LE (mixed 16/32-bit).
	ImageTESignature    = 0x5A56     // Terse Executables have a 'VZ' signature
	ImageNTSignature    = 0x00004550 // PE00
)

// Optional Header magic
const (
	ImageNtOptionalHeader32Magic = 0x10b
	ImageNtOptionalHeader64Magic = 0x20b
	ImageROMOptionalHeaderMagic  = 0x10
)

// Image file machine types
const (
	ImageFileMachineUnknown   uint32 = 0x0    // The contents of this field are assumed to be applicable to any machine type
	ImageFileMachineAM33      uint32 = 0x1d3  // Matsushita AM33
	ImageFileMachineAMD64     uint32 = 0x8664 // x64
	ImageFileMachineARM       uint32 = 0x1c0  // ARM little endian
	ImageFileMachineARM64     uint32 = 0xaa64 // ARM64 little endian
	ImageFileMachineARMNT     uint32 = 0x1c4  // ARM Thumb-2 little endian
	ImageFileMachineEBC       uint32 = 0xebc  // EFI byte code
	ImageFileMachineI386      uint32 = 0x14c  // Intel 386 or later processors and compatible processors
	ImageFileMachineIA64      uint32 = 0x200  // Intel Itanium processor family
	ImageFileMachineM32R      uint32 = 0x9041 // Mitsubishi M32R little endian
	ImageFileMachineMIPS16    uint32 = 0x266  // MIPS16
	ImageFileMachineMIPSFPU   uint32 = 0x366  // MIPS with FPU
	ImageFileMachineMIPSFPU16 uint32 = 0x466  // MIPS16 with FPU
	ImageFileMachinePOWERPC   uint32 = 0x1f0  // Power PC little endian
	ImageFileMachinePOWERPCFP uint32 = 0x1f1  // Power PC with floating point support
	ImageFileMachineR4000     uint32 = 0x166  // MIPS little endian
	ImageFileMachineRISCV32   uint32 = 0x5032 // RISC-V 32-bit address space
	ImageFileMachineRISCV64   uint32 = 0x5064 // RISC-V 64-bit address space
	ImageFileMachineRISCV128  uint32 = 0x5128 // RISC-V 128-bit address space
	ImageFileMachineSH3       uint32 = 0x1a2  // Hitachi SH3
	ImageFileMachineSH3DSP    uint32 = 0x1a3  // Hitachi SH3 DSP
	ImageFileMachineSH4       uint32 = 0x1a6  // Hitachi SH4
	ImageFileMachineSH5       uint32 = 0x1a8  // Hitachi SH5
	ImageFileMachineTHUMB     uint32 = 0x1c2  // Thumb
	ImageFileMachineWCEMIPSV2 uint32 = 0x169  // MIPS little-endian WCE v2
)

// The Characteristics field contains flags that indicate attributes of the object or image file.
const (
	ImageFileRelocsStripped       = 0x0001 // Relocation info stripped from file.
	ImageFileExecutableImage      = 0x0002 // File is executable  (i.e. no unresolved external references).
	ImageFileLineNumsStripped     = 0x0004 // Line nunbers stripped from file.
	ImageFileLocalSymsStripped    = 0x0008 // Local symbols stripped from file.
	ImageFileAgressibeWsTrim      = 0x0010 // Aggressively trim working set
	ImageFileLargeAddressAware    = 0x0020 // App can handle >2gb addresses
	ImageFileBytesReservedLow     = 0x0080 // Bytes of machine word are reversed.
	ImageFile32BitMachine         = 0x0100 // 32 bit word machine.
	ImageFileDebugStripped        = 0x0200 // Debugging info stripped from file in .DBG file
	ImageFileRemovableRunFromSwap = 0x0400 // If Image is on removable media, copy and run from the swap file.
	ImageFileNetRunFromSwap       = 0x0800 // If Image is on Net, copy and run from the swap file.
	ImageFileSystem               = 0x1000 // System File.
	ImageFileDLL                  = 0x2000 // File is a DLL.
	ImageFileUpSystemOnly         = 0x4000 // File should only be run on a UP machine
	ImageFileBytesReservedHigh    = 0x8000 // Bytes of machine word are reversed.
)

// Subsystem values of an OptionalHeader
const (
	ImageSubsystemUnknown                = 0  // An unknown subsystem.
	ImageSubsystemNative                 = 1  // Device drivers and native Windows processes
	ImageSubsystemWindowsGui             = 2  // The Windows graphical user interface (GUI) subsystem.
	ImageSubsystemWindowsCui             = 3  // The Windows character subsystem
	ImageSubsystemOS2Cui                 = 5  // The OS/2 character subsystem.
	ImageSubsystemPosixCui               = 7  // The Posix character subsystem.
	ImageSubsystemNativeWindows          = 8  // Native Win9x driver
	ImageSubsystemWindowsCEGui           = 9  // Windows CE
	ImageSubsystemEFIApplication         = 10 // An Extensible Firmware Interface (EFI) application
	ImageSubsystemEFIBootServiceDriver   = 11 // An EFI driver with boot services
	ImageSubsystemEFIRuntimeDriver       = 12 // An EFI driver with run-time services
	ImageSubsystemEFIRom                 = 13 // An EFI ROM image .
	ImageSubsystemXBOX                   = 14 // XBOX.
	ImageSubsystemWindowsBootApplication = 16 // Windows boot application.
)

// DllCharacteristics values of an OptionalHeader
const (
	ImageDllCharacteristicsReserved1            = 0x0001 // Reserved, must be zero.
	ImageDllCharacteristicsReserved2            = 0x0002 // Reserved, must be zero.
	ImageDllCharacteristicsReserved4            = 0x0004 // Reserved, must be zero.
	ImageDllCharacteristicsReserved8            = 0x0008 // Reserved, must be zero.
	ImageDllCharacteristicsHighEntropyVA        = 0x0020 // Image can handle a high entropy 64-bit virtual address space
	ImageDllCharacteristicsDynamicBase          = 0x0040 // DLL can be relocated at load time.
	ImageDllCharacteristicsForceIntegrity       = 0x0080 // Code Integrity checks are enforced.
	ImageDllCharacteristicsNXCompact            = 0x0100 // Image is NX compatible.
	ImageDllCharacteristicsNoIsolation          = 0x0200 // Isolation aware, but do not isolate the image.
	ImageDllCharacteristicsNoSEH                = 0x0400 // Does not use structured exception (SE) handling. No SE handler may be called in this image.
	ImageDllCharacteristicsNoBind               = 0x0800 // Do not bind the image.
	ImageDllCharacteristicsAppContainer         = 0x1000 // Image must execute in an AppContainer
	ImageDllCharacteristicsWdmDriver            = 0x2000 // A WDM driver.
	ImageDllCharacteristicsGuardCF              = 0x4000 // Image supports Control Flow Guard.
	ImageDllCharacteristicsTerminalServiceAware = 0x8000 // Terminal Server aware.

)

// DataDirectory entries of an OptionalHeader
const (
	ImageDirectoryEntryExport       = 0  // Export Table
	ImageDirectoryEntryImport       = 1  // Import Table
	ImageDirectoryEntryResource     = 2  // Resource Table
	ImageDirectoryEntryException    = 3  // Exception Table
	ImageDirectoryEntryCertificate  = 4  // Certificate Directory
	ImageDirectoryEntryBaseReloc    = 5  // Base Relocation Table
	ImageDirectoryEntryDebug        = 6  // Debug
	ImageDirectoryEntryArchitecture = 7  // Architecture Specific Data
	ImageDirectoryEntryGlobalPtr    = 8  // The RVA of the value to be stored in the global pointer register.
	ImageDirectoryEntryTLS          = 9  // The thread local storage (TLS) table
	ImageDirectoryEntryLoadConfig   = 10 // The load configuration table
	ImageDirectoryEntryBoundImport  = 11 // The bound import table
	ImageDirectoryEntryIAT          = 12 // Import Address Table
	ImageDirectoryEntryDelayImport  = 13 // Delay Import Descriptor
	ImageDirectoryEntryCLR          = 14 // CLR Runtime Header
	ImageDirectoryEntryRESERVED     = 15 // Must be zero
	ImageNumberOfDirectoryEntries   = 16 // Tables count.
)

// ImageNtHeader represents the PE header and is the general term for a structure named IMAGE_NT_HEADERS.
type ImageNtHeader struct {
	// Signature is a DWORD containing the value 50h, 45h, 00h, 00h.
	Signature uint32
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

// OptionalHeader32 represents the PE32 format structure of the optional header.
// PE32 contains this additional field, which is absent in PE32+.
type OptionalHeader32 struct {
	Magic                       uint16            // The unsigned integer that identifies the state of the image file. The most common number is 0x10B, which identifies it as a normal executable file. 0x107 identifies it as a ROM image, and 0x20B identifies it as a PE32+ executable.
	MajorLinkerVersion          uint8             // The linker major version number.
	MinorLinkerVersion          uint8             // The linker minor version number.
	SizeOfCode                  uint32            // The size of the code (text) section, or the sum of all code sections if there are multiple sections.
	SizeOfInitializedData       uint32            // The size of the initialized data section, or the sum of all such sections if there are multiple data sections.
	SizeOfUninitializedData     uint32            // The size of the uninitialized data section (BSS), or the sum of all such sections if there are multiple BSS sections.
	AddressOfEntryPoint         uint32            // The address of the entry point relative to the image base when the executable file is loaded into memory. For program images, this is the starting address. For device drivers, this is the address of the initialization function. An entry point is optional for DLLs. When no entry point is present, this field must be zero.
	BaseOfCode                  uint32            // The address that is relative to the image base of the beginning-of-code section when it is loaded into memory.
	BaseOfData                  uint32            // The address that is relative to the image base of the beginning-of-data section when it is loaded into memory.
	ImageBase                   uint32            // The preferred address of the first byte of image when loaded into memory; must be a multiple of 64 K. The default for DLLs is 0x10000000. The default for Windows CE EXEs is 0x00010000. The default for Windows NT, Windows 2000, Windows XP, Windows 95, Windows 98, and Windows Me is 0x00400000.
	SectionAlignment            uint32            // The alignment (in bytes) of sections when they are loaded into memory. It must be greater than or equal to FileAlignment. The default is the page size for the architecture.
	FileAlignment               uint32            // The alignment factor (in bytes) that is used to align the raw data of sections in the image file. The value should be a power of 2 between 512 and 64 K, inclusive. The default is 512. If the SectionAlignment is less than the architecture's page size, then FileAlignment must match SectionAlignment.
	MajorOperatingSystemVersion uint16            // The major version number of the required operating system.
	MinorOperatingSystemVersion uint16            // The minor version number of the required operating system.
	MajorImageVersion           uint16            // The major version number of the image.
	MinorImageVersion           uint16            // The minor version number of the image.
	MajorSubsystemVersion       uint16            // The major version number of the subsystem.
	MinorSubsystemVersion       uint16            // The minor version number of the subsystem.
	Win32VersionValue           uint32            // Reserved, must be zero.
	SizeOfImage                 uint32            // The size (in bytes) of the image, including all headers, as the image is loaded in memory. It must be a multiple of SectionAlignment.
	SizeOfHeaders               uint32            // The combined size of an MS-DOS stub, PE header, and section headers rounded up to a multiple of FileAlignment.
	CheckSum                    uint32            // The image file checksum. The algorithm for computing the checksum is incorporated into IMAGHELP.DLL. The following are checked for validation at load time: all drivers, any DLL loaded at boot time, and any DLL that is loaded into a critical Windows process
	Subsystem                   uint16            // The subsystem that is required to run this image. For more information, see Windows Subsystem.
	DllCharacteristics          uint16            // For more information, see DLL Characteristics later in this specification.
	SizeOfStackReserve          uint32            // The size of the stack to reserve. Only SizeOfStackCommit is committed; the rest is made available one page at a time until the reserve size is reached.
	SizeOfStackCommit           uint32            // The size of the stack to commit.
	SizeOfHeapReserve           uint32            // The size of the local heap space to reserve. Only SizeOfHeapCommit is committed; the rest is made available one page at a time until the reserve size is reached.
	SizeOfHeapCommit            uint32            // The size of the local heap space to commit.
	LoaderFlags                 uint32            // Reserved, must be zero.
	NumberOfRvaAndSizes         uint32            // The number of data-directory entries in the remainder of the optional header. Each describes a location and size.
	DataDirectory               [16]DataDirectory // An array of 16 IMAGE_DATA_DIRECTORY structures.
}

// OptionalHeader64 represents the PE32+ format structure of the optional header.
type OptionalHeader64 struct {
	Magic                       uint16            // The unsigned integer that identifies the state of the image file. The most common number is 0x10B, which identifies it as a normal executable file. 0x107 identifies it as a ROM image, and 0x20B identifies it as a PE32+ executable.
	MajorLinkerVersion          uint8             // The linker major version number.
	MinorLinkerVersion          uint8             // The linker minor version number.
	SizeOfCode                  uint32            // The size of the code (text) section, or the sum of all code sections if there are multiple sections.
	SizeOfInitializedData       uint32            // The size of the initialized data section, or the sum of all such sections if there are multiple data sections.
	SizeOfUninitializedData     uint32            // The size of the uninitialized data section (BSS), or the sum of all such sections if there are multiple BSS sections.
	AddressOfEntryPoint         uint32            // The address of the entry point relative to the image base when the executable file is loaded into memory. For program images, this is the starting address. For device drivers, this is the address of the initialization function. An entry point is optional for DLLs. When no entry point is present, this field must be zero.
	BaseOfCode                  uint32            // The address that is relative to the image base of the beginning-of-code section when it is loaded into memory.
	ImageBase                   uint64            // In PE+, ImageBase is 8 bytes size.
	SectionAlignment            uint32            // The alignment (in bytes) of sections when they are loaded into memory. It must be greater than or equal to FileAlignment. The default is the page size for the architecture.
	FileAlignment               uint32            // The alignment factor (in bytes) that is used to align the raw data of sections in the image file. The value should be a power of 2 between 512 and 64 K, inclusive. The default is 512. If the SectionAlignment is less than the architecture's page size, then FileAlignment must match SectionAlignment.
	MajorOperatingSystemVersion uint16            // The major version number of the required operating system.
	MinorOperatingSystemVersion uint16            // The minor version number of the required operating system.
	MajorImageVersion           uint16            // The major version number of the image.
	MinorImageVersion           uint16            // The minor version number of the image.
	MajorSubsystemVersion       uint16            // The major version number of the subsystem.
	MinorSubsystemVersion       uint16            // The minor version number of the subsystem.
	Win32VersionValue           uint32            // Reserved, must be zero.
	SizeOfImage                 uint32            // The size (in bytes) of the image, including all headers, as the image is loaded in memory. It must be a multiple of SectionAlignment.
	SizeOfHeaders               uint32            // The combined size of an MS-DOS stub, PE header, and section headers rounded up to a multiple of FileAlignment.
	CheckSum                    uint32            // The image file checksum. The algorithm for computing the checksum is incorporated into IMAGHELP.DLL. The following are checked for validation at load time: all drivers, any DLL loaded at boot time, and any DLL that is loaded into a critical Windows process
	Subsystem                   uint16            // The subsystem that is required to run this image. For more information, see Windows Subsystem.
	DllCharacteristics          uint16            // For more information, see DLL Characteristics later in this specification.
	SizeOfStackReserve          uint64            // In PE+, this field is 8 bytes size.
	SizeOfStackCommit           uint64            // In PE+, this field is 8 bytes size.
	SizeOfHeapReserve           uint64            // In PE+, this field is 8 bytes size.
	SizeOfHeapCommit            uint64            // In PE+, this field is 8 bytes size.
	LoaderFlags                 uint32            // Reserved, must be zero.
	NumberOfRvaAndSizes         uint32            // The number of data-directory entries in the remainder of the optional header. Each describes a location and size.
	DataDirectory               [16]DataDirectory // An array of 16 IMAGE_DATA_DIRECTORY structures.
}

// DataDirectory represents an array of 16 IMAGE_DATA_DIRECTORY structures, 8 bytes apiece, each relating to an important data structure in the PE file.
type DataDirectory struct {
	VirtualAddress uint32 // The RVA of the data structure.
	Size           uint32 // The size in bytes of the data structure refered to.
}
