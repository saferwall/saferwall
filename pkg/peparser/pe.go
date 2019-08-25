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
