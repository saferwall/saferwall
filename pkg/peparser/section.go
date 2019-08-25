package pe


import (
	"fmt"
	"reflect"
)

// ImageSectionHeader is part of the section table , in fact section table is an array of Image Section Header
// each contains information about one section of the whole file such as attribute,virtual offset.
// the array size is the number of sections in the file .
// Binary Spec : each struct is 40 byte and there is no padding .
type ImageSectionHeader struct {
	Name                 [8]uint8 // Name of the section.
	VirtualSize          uint32  // The total size of the section when loaded into memory. If this value is greater than SizeOfRawData, the section is zero-padded. This field is valid only for executable images and should be set to zero for object files.
	VirtualAddress       uint32  // For executable images, the address of the first byte of the section relative to the image base when the section is loaded into memory. For object files, this field is the address of the first byte before relocation is applied; for simplicity, compilers should set this to zero. Otherwise, it is an arbitrary value that is subtracted from offsets during relocation.
	SizeOfRawData        uint32  // The size of the section (for object files) or the size of the initialized data on disk (for image files). For executable images, this must be a multiple of FileAlignment from the optional header. If this is less than VirtualSize, the remainder of the section is zero-filled. Because the SizeOfRawData field is rounded but the VirtualSize field is not, it is possible for SizeOfRawData to be greater than VirtualSize as well. When a section contains only uninitialized data, this field should be zero.
	PointerToRawData     uint32  // The file pointer to the beginning of relocation entries for the section. This is set to zero for executable images or if there are no relocations.
	PointerToRelocations uint32  // The size of the section (for object files) or the size of the initialized data on disk (for image files). For executable images, this must be a multiple of FileAlignment from the optional header. If this is less than VirtualSize, the remainder of the section is zero-filled. Because the SizeOfRawData field is rounded but the VirtualSize field is not, it is possible for SizeOfRawData to be greater than VirtualSize as well. When a section contains only uninitialized data, this field should be zero.
	PointerToLineNumbers uint32  // The file pointer to the beginning of line-number entries for the section. This is set to zero if there are no COFF line numbers. This value should be zero for an image because COFF debugging information is deprecated.
	NumberOfRelocations  uint16  // The number of relocation entries for the section. This is set to zero for executable images.
	NumberOfLineNumbers  uint16  // The number of line-number entries for the section. This value should be zero for an image because COFF debugging information is deprecated.
	Characteristics      uint32  // The flags that describe the characteristics of the section. For more information, see Section Flags.
}

// NameString returns string represntation of a ImageSectionHeader.Name field.
func (section *ImageSectionHeader) NameString() string {
    return fmt.Sprintf("%s", section.Name)
}

// NextHeaderAddr returns the VirtualAddress of the next section.
func (section *ImageSectionHeader) NextHeaderAddr(pe* File) uint32 {
	for i, currentSection := range pe.Sections {
		if i == len(pe.Sections) {
			return 0
		}

		if reflect.DeepEqual(section, &currentSection){
			return pe.Sections[i+1].VirtualAddress
		}
	}

	return 0
}

// Contains checks whether the section contains a given RVA.
func (section *ImageSectionHeader) Contains(rva uint32, pe* File) bool {

	// Check if the SizeOfRawData is realistic. If it's bigger than the size of
	// the whole PE file minus the start address of the section it could be
	// either truncated or the SizeOfRawData contains a misleading value.
	// In either of those cases we take the VirtualSize.


	var size uint32
	adjustedPointer := pe.adjustFileAlignment(section.PointerToRawData)
	if uint32(len(pe.data)) - adjustedPointer < section.SizeOfRawData {
		size = section.VirtualSize
	} else {
		size = Max(section.SizeOfRawData, section.VirtualSize)
	}
	vaAdj := pe.adjustSectionAlignment(section.VirtualAddress)

	// Check whether there's any section after the current one that starts before the
	// calculated end for the current one. If so, cut the current section's size
	// to fit in the range up to where the next section starts.
	if section.NextHeaderAddr(pe) != 0 &&
		section.NextHeaderAddr(pe) > section.VirtualAddress &&
		vaAdj + size > section.NextHeaderAddr(pe) {
		size = section.NextHeaderAddr(pe) - vaAdj
	}

	return vaAdj <= rva && rva < vaAdj + size
}

// Data returns a data chunk from a section
func (section *ImageSectionHeader) Data(start, length uint32, pe* File) []byte {
	
	pointerToRawDataAdj := pe.adjustFileAlignment(section.PointerToRawData)
	virtualAddressAdj := pe.adjustSectionAlignment(section.VirtualAddress)
	
	var offset uint32
	if start == 0 {
		offset = pointerToRawDataAdj
	} else {
		offset = (start - virtualAddressAdj) + pointerToRawDataAdj
	}

	var end uint32
	if start != 0 {
		end = offset + length
	} else {
		end = offset + section.SizeOfRawData
	}

	// PointerToRawData is not adjusted here as we might want to read any possible extra bytes
	// that might get cut off by aligning the start (and hence cutting something off the end)
	if end > section.PointerToRawData + section.SizeOfRawData {
		end = section.PointerToRawData + section.SizeOfRawData
	}

	return pe.data[offset:end]
}


// byVirtualAddress sorts all sections by Virtual Address.
type byVirtualAddress []ImageSectionHeader

func (s byVirtualAddress) Len() int { return len(s) }
func (s byVirtualAddress) Less(i, j int) bool { return s[i].VirtualAddress < s[j].VirtualAddress }
func (s byVirtualAddress) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
