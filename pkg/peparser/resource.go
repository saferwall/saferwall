package pe

import (
	"bytes"
	"encoding/binary"
	"log"
	"unsafe"
)

const (
	maxAllowedEntries = 0x1000
)

// Predefined Resource Types
var (
	RTCursor       = makeIntResource(1)
	RTBitmap       = makeIntResource(2)
	RTIcon         = makeIntResource(3)
	RTMenu         = makeIntResource(4)
	RTDialog       = makeIntResource(5)
	RTString       = makeIntResource(6)
	RTFontdir      = makeIntResource(7)
	RTFont         = makeIntResource(8)
	RTAccelerator  = makeIntResource(9)
	RTRCdata       = makeIntResource(10)
	RTMessagetable = makeIntResource(11)
)

// ImageResourceDirectory represents the IMAGE_RESOURCE_DIRECTORY.
// This data structure should be considered the heading of a table because the table actually consists of directory entries.
type ImageResourceDirectory struct {
	Characteristics      uint32 // Resource flags. This field is reserved for future use. It is currently set to zero.
	TimeDateStamp        uint32 // The time that the resource data was created by the resource compiler.
	MajorVersion         uint16 // The major version number, set by the user.
	MinorVersion         uint16 // The minor version number, set by the user.
	NumberOfNamedEntries uint16 // The number of directory entries immediately following the table that use strings to identify Type, Name, or Language entries (depending on the level of the table).
	NumberOfIDEntries    uint16 // The number of directory entries immediately following the Name entries that use numeric IDs for Type, Name, or Language entries.
}

// ImageResourceDirectoryEntry represents an entry in the resource directory entries.
type ImageResourceDirectoryEntry struct {
	Name         uint32 // is used to identify either a type of resource, a resource name, or a resource's language ID.
	OffsetToData uint32 // is always used to point to a sibling in the tree, either a directory node or a leaf node.
}

// ImageResourceDataEntry: Each Resource Data entry describes an actual unit of raw data in the Resource Data area.
type ImageResourceDataEntry struct {
	OffsetToData uint32 // The address of a unit of resource data in the Resource Data area.
	Size         uint32 // The size, in bytes, of the resource data that is pointed to by the Data RVA field.
	CodePage     uint32 // The code page that is used to decode code point values within the resource data. Typically, the code page would be the Unicode code page.
	Reserved     uint32 // Reserved, must be 0.
}

// makeIntResource mimics the MAKEINTRESOURCE macro.
func makeIntResource(id uintptr) *uint16 {
	return (*uint16)(unsafe.Pointer(id))
}

func (pe *File) parseResourceDirectory(rva, size uint32) error {

	// Get the resource directory structure, that is, the header
	// If the table preceding the actual entries
	resourceDir := ImageResourceDirectory{}
	resourceDirSize := uint32(binary.Size(resourceDir))
	offset := pe.getOffsetFromRva(rva)
	buff := bytes.NewReader(pe.data[offset : offset+resourceDirSize])
	err := binary.Read(buff, binary.LittleEndian, &resourceDir)
	if err != nil {
		// If we can't parse resources directory then silently return.
		// This directory does not necessarily have to be valid to
		// still have a valid PE file
		log.Println("Invalid resources directory. Can't parse directory data")
		return nil
	}

	// Advance the RVA to the position immediately following the directory
	// table header and pointing to the first entry in the table
	rva += resourceDirSize

	numberOfEntries := resourceDir.NumberOfNamedEntries + resourceDir.NumberOfIDEntries

	// Set a hard limit on the maximum reasonable number of entries
	if numberOfEntries > maxAllowedEntries {
		log.Printf("Error parsing the resources directory. The directory contains %d entries", numberOfEntries)
		return nil
	}

	return nil

}
