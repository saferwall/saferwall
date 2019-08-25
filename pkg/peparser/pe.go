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

// ImageDosHeader represents the DOS stub of a PE.
type ImageDosHeader struct {
	Emagic    uint16     // Magic number
	Ecblp     uint16     // Bytes on last page of file
	Ecp       uint16     // Pages in file
	Ecrlc     uint16     // Relocations
	Ecparhdr  uint16     // Size of header in paragraphs
	Eminalloc uint16     // Minimum extra paragraphs needed
	Emaxalloc uint16     // Maximum extra paragraphs needed
	Ess       uint16     // Initial (relative) SS value
	Esp       uint16     // Initial SP value
	Ecsum     uint16     // Checksum
	Eip       uint16     // Initial IP value
	Ecs       uint16     // Initial (relative) CS value
	Elfarlc   uint16     // File address of relocation table
	Eovno     uint16     // Overlay number
	Eres      [4]uint16  // Reserved words
	Eoemid    uint16     // OEM identifier (for e_oeminfo)
	Eoeminfo  uint16     // OEM information; e_oemid specific
	Eres2     [10]uint16 // Reserved words
	Elfanew   uint32     // File address of new exe header
}

// ImageNtHeader represents the PE header and is the general term for a structure named IMAGE_NT_HEADERS.
type ImageNtHeader struct {
	Signature uint32 // Signature is a DWORD containing the value 50h, 45h, 00h, 00h.
}



// ImageFileHeader contains info about the physical layout and properties of the file.
type ImageFileHeader struct {
	Machine              uint16 // The number that identifies the type of target machine
	NumberOfSections     uint16 // The number of sections. This indicates the size of the section table, which immediately follows the headers
	TimeDateStamp        uint32 // The low 32 bits of the number of seconds since 00:00 January 1, 1970 (a C run-time time_t value), that indicates when the file was created.
	PointerToSymbolTable uint32 // The file offset of the COFF symbol table, or zero if no COFF symbol table is present. This value should be zero for an image because COFF debugging information is deprecated.
	NumberOfSymbols      uint32 // The number of entries in the symbol table. This data can be used to locate the string table, which immediately follows the symbol table. This value should be zero for an image because COFF debugging information is deprecated.
	SizeOfOptionalHeader uint16 // The size of the optional header, which is required for executable files but not for object files. This value should be zero for an object file.
	Characteristics      uint16 // The flags that indicate the attributes of the file.
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
	Magic                       uint16 // The unsigned integer that identifies the state of the image file. The most common number is 0x10B, which identifies it as a normal executable file. 0x107 identifies it as a ROM image, and 0x20B identifies it as a PE32+ executable.
	MajorLinkerVersion          uint8  // The linker major version number.
	MinorLinkerVersion          uint8  // The linker minor version number.
	SizeOfCode                  uint32 // The size of the code (text) section, or the sum of all code sections if there are multiple sections.
	SizeOfInitializedData       uint32 // The size of the initialized data section, or the sum of all such sections if there are multiple data sections.
	SizeOfUninitializedData     uint32 // The size of the uninitialized data section (BSS), or the sum of all such sections if there are multiple BSS sections.
	AddressOfEntryPoint         uint32 // The address of the entry point relative to the image base when the executable file is loaded into memory. For program images, this is the starting address. For device drivers, this is the address of the initialization function. An entry point is optional for DLLs. When no entry point is present, this field must be zero.
	BaseOfCode                  uint32 // The address that is relative to the image base of the beginning-of-code section when it is loaded into memory.
	ImageBase                   uint64 // In PE+, ImageBase is 8 bytes size.
	SectionAlignment            uint32 // The alignment (in bytes) of sections when they are loaded into memory. It must be greater than or equal to FileAlignment. The default is the page size for the architecture.
	FileAlignment               uint32 // The alignment factor (in bytes) that is used to align the raw data of sections in the image file. The value should be a power of 2 between 512 and 64 K, inclusive. The default is 512. If the SectionAlignment is less than the architecture's page size, then FileAlignment must match SectionAlignment.
	MajorOperatingSystemVersion uint16 // The major version number of the required operating system.
	MinorOperatingSystemVersion uint16 // The minor version number of the required operating system.
	MajorImageVersion           uint16 // The major version number of the image.
	MinorImageVersion           uint16 // The minor version number of the image.
	MajorSubsystemVersion       uint16 // The major version number of the subsystem.
	MinorSubsystemVersion       uint16 // The minor version number of the subsystem.
	Win32VersionValue           uint32 // Reserved, must be zero.
	SizeOfImage                 uint32 // The size (in bytes) of the image, including all headers, as the image is loaded in memory. It must be a multiple of SectionAlignment.
	SizeOfHeaders               uint32 // The combined size of an MS-DOS stub, PE header, and section headers rounded up to a multiple of FileAlignment.
	CheckSum                    uint32 // The image file checksum. The algorithm for computing the checksum is incorporated into IMAGHELP.DLL. The following are checked for validation at load time: all drivers, any DLL loaded at boot time, and any DLL that is loaded into a critical Windows process
	Subsystem                   uint16 // The subsystem that is required to run this image. For more information, see Windows Subsystem.
	DllCharacteristics          uint16 // For more information, see DLL Characteristics later in this specification.
	SizeOfStackReserve          uint64 // In PE+, this field is 8 bytes size.
	SizeOfStackCommit           uint64 // In PE+, this field is 8 bytes size.
	SizeOfHeapReserve           uint64 // In PE+, this field is 8 bytes size.
	SizeOfHeapCommit            uint64 // In PE+, this field is 8 bytes size.
	LoaderFlags                 uint32 // Reserved, must be zero.
	NumberOfRvaAndSizes         uint32 // The number of data-directory entries in the remainder of the optional header. Each describes a location and size.
	DataDirectory               [16]DataDirectory
}

// DataDirectory represents an array of 16 IMAGE_DATA_DIRECTORY structures, 8 bytes apiece, each relating to an important data structure in the PE file.
type DataDirectory struct {
	VirtualAddress uint32 // The RVA of the data structure.
	Size           uint32 // The size in bytes of the data structure refered to.
}
