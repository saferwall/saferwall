// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

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

var (
	depth = 0
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
	RTGroupCursor  = makeIntResource(12)
	RTGroupIcon    = makeIntResource(14)
	RTVersion      = makeIntResource(16)
	RTDlgInclude   = makeIntResource(17)
	RTPlugPlay     = makeIntResource(19)
	RTVxd          = makeIntResource(20)
	RTAniCursor    = makeIntResource(21)
	RTAniIcon      = makeIntResource(22)
	RTHtml         = makeIntResource(23)
	RTManifest     = makeIntResource(24)
)

// ImageResourceDirectory represents the IMAGE_RESOURCE_DIRECTORY.
// This data structure should be considered the heading of a table because the
// table actually consists of directory entries.
type ImageResourceDirectory struct {
	// Resource flags. This field is reserved for future use. It is currently
	// set to zero.
	Characteristics uint32

	// The time that the resource data was created by the resource compiler.
	TimeDateStamp uint32

	// The major version number, set by the user.
	MajorVersion uint16

	// The minor version number, set by the user.
	MinorVersion uint16

	// The number of directory entries immediately following the table that use
	// strings to identify Type, Name, or Language entries (depending on the
	// level of the table).
	NumberOfNamedEntries uint16

	// The number of directory entries immediately following the Name entries
	// that use numeric IDs for Type, Name, or Language entries.
	NumberOfIDEntries uint16
}

// ImageResourceDirectoryEntry represents an entry in the resource directory entries.
type ImageResourceDirectoryEntry struct {
	// is used to identify either a type of resource, a resource name, or a
	// resource's language ID.
	Name uint32

	//is always used to point to a sibling in the tree, either a directory node
	// or a leaf node.
	OffsetToData uint32
}

// ImageResourceDataEntry Each Resource Data entry describes an actual unit of
// raw data in the Resource Data area.
type ImageResourceDataEntry struct {
	// The address of a unit of resource data in the Resource Data area.
	OffsetToData uint32

	// The size, in bytes, of the resource data that is pointed to by the Data
	// RVA field.
	Size uint32

	// The code page that is used to decode code point values within the
	// resource data. Typically, the code page would be the Unicode code page.
	CodePage uint32

	// Reserved, must be 0.
	Reserved uint32
}

// ResourceDirectory represents resource directory information.
type ResourceDirectory struct {
	// IMAGE_RESOURCE_DIRECTORY structure
	Struct ImageResourceDirectory

	// list of entries
	Entries []ResourceDirectoryEntry
}

// ResourceDirectoryEntry represents a resource directory entry.
type ResourceDirectoryEntry struct {
	// IMAGE_RESOURCE_DIRECTORY_ENTRY structure.
	Struct ImageResourceDirectoryEntry

	// If the resource is identified by name this attribute will contain the
	// name string. Empty string otherwise. If identified by id, the id is
	//available at .Id field.
	Name string

	// The resource identifier.
	ID uint32

	// If this entry has a lower level directory this attribute will point to
	// the ResourceDirData instance representing it.
	Directory ResourceDirectory

	// If this entry has no further lower directories and points to the actual
	// resource data, this attribute will reference the corresponding
	// ResourceDataEntry instance.
	Data ResourceDataEntry
}

// ResourceDataEntry represents a resource data entry.
type ResourceDataEntry struct {

	// IMAGE_RESOURCE_DATA_ENTRY structure.
	Struct ImageResourceDataEntry

	// Primary language ID
	Lang    uint32
	Sublang uint32 // Sublanguage ID
}

// makeIntResource mimics the MAKEINTRESOURCE macro.
func makeIntResource(id uintptr) *uint16 {
	return (*uint16)(unsafe.Pointer(id))
}

func (pe *File) readUnicodeStringAtRVA(rva uint32, maxLength uint32) string {
	unicodeString := ""
	offset := pe.getOffsetFromRva(rva)
	buff := pe.data[offset : offset+(maxLength*2)]
	for i := uint32(0); i < maxLength*2; i += 2 {
		unicodeString += string(buff[i])
	}
	return unicodeString

}

func (pe *File) parseResourceDataEntry(rva uint32) *ImageResourceDataEntry {
	dataEntry := ImageResourceDataEntry{}
	dataEntrySize := uint32(binary.Size(dataEntry))
	offset := pe.getOffsetFromRva(rva)
	buff := bytes.NewReader(pe.data[offset : offset+dataEntrySize])
	err := binary.Read(buff, binary.LittleEndian, &dataEntry)
	if err != nil {
		log.Println("Error parsing a resource directory data entry, the RVA is invalid")
		return nil
	}
	return &dataEntry
}

func (pe *File) parseResourceDirectoryEntry(rva uint32) *ImageResourceDirectoryEntry {
	resource := ImageResourceDirectoryEntry{}
	resourceSize := uint32(binary.Size(resource))
	offset := pe.getOffsetFromRva(rva)
	buff := bytes.NewReader(pe.data[offset : offset+resourceSize])
	err := binary.Read(buff, binary.LittleEndian, &resource)
	if err != nil {
		return nil
	}

	if resource == (ImageResourceDirectoryEntry{}) {
		return nil
	}

	// resource.NameOffset = resource.Name & 0x7FFFFFFF

	// resource.__pad = resource.Name & 0xFFFF0000
	// resource.Id = resource.Name & 0x0000FFFF

	// resource.DataIsDirectory = (resource.OffsetToData & 0x80000000) >> 31
	// resource.OffsetToDirectory = resource.OffsetToData & 0x7FFFFFFF

	return &resource
}

func (pe *File) parseResourceDirectory(rva, size, baseRVA, level uint32) (ResourceDirectory, error) {
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
		return ResourceDirectory{}, nil
	}

	if baseRVA == 0 {
		baseRVA = rva
	}

	// Advance the RVA to the position immediately following the directory
	// table header and pointing to the first entry in the table
	rva += resourceDirSize

	numberOfEntries := int(resourceDir.NumberOfNamedEntries + resourceDir.NumberOfIDEntries)
	var dirEntries []ResourceDirectoryEntry

	// Set a hard limit on the maximum reasonable number of entries
	if numberOfEntries > maxAllowedEntries {
		log.Printf("Error parsing the resources directory. The directory contains %d entries", numberOfEntries)
		return ResourceDirectory{}, nil
	}

	for i := 0; i < numberOfEntries; i++ {
		res := pe.parseResourceDirectoryEntry(rva)
		if res == nil {
			log.Println("Error parsing a resource directory entry, the RVA is invalid")
			break
		}

		nameIsString := (res.Name & 0x80000000) >> 31
		entryName := ""
		entryID := uint32(0)
		if nameIsString == 0 {
			entryID = res.Name
		} else {
			nameOffset := res.Name & 0x7FFFFFFF
			uStringOffset := pe.getOffsetFromRva(baseRVA + nameOffset)
			maxLen := binary.LittleEndian.Uint16(pe.data[uStringOffset:])
			entryName = pe.readUnicodeStringAtRVA(baseRVA+nameOffset+2, uint32(maxLen))
		}

		dataIsDirectory := (res.OffsetToData & 0x80000000) >> 31
		OffsetToDirectory := res.OffsetToData & 0x7FFFFFFF
		if dataIsDirectory > 0 {
			// One trick malware can do is to recursively reference
			// the next directory. This causes hilarity to ensue when
			// trying to parse everything correctly.
			// If the original RVA given to this function is equal to
			// the next one to parse, we assume that it's a trick.
			// Instead of raising a PEFormatError this would skip some
			// reasonable data so we just break.
			// 9ee4d0a0caf095314fd7041a3e4404dc is the offending sample.
			level++
			directoryEntry, _ := pe.parseResourceDirectory(
				baseRVA+OffsetToDirectory,
				size-(rva+baseRVA),
				baseRVA,
				level)

			dirEntries = append(dirEntries, ResourceDirectoryEntry{
				Struct:    *res,
				Name:      entryName,
				ID:        entryID,
				Directory: directoryEntry})
		} else {
			// data is entry
			dataEntryStruct := pe.parseResourceDataEntry(baseRVA + OffsetToDirectory)
			entryData := ResourceDataEntry{
				Struct:  *dataEntryStruct,
				Lang:    res.Name & 0x3ff,
				Sublang: res.Name >> 10,
			}

			dirEntries = append(dirEntries, ResourceDirectoryEntry{
				Struct: *res,
				Name:   entryName,
				ID:     entryID,
				Data:   entryData})
		}

		rva += uint32(binary.Size(res))
	}

	return ResourceDirectory{
		Struct:  resourceDir,
		Entries: dirEntries,
	}, nil
}
