package pe

// Image executable types
const (
	ImageDOSSignature   = 0x5A4D // MZ
	ImageDOSZMSignature = 0x4D5A // ZM

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
