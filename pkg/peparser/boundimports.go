package pe

import (
	"bytes"
	"encoding/binary"
	"log"
)

const (
	// MaxStringLength represents the maximum length of a string to be retrieved
	// from the file. It's there to prevent loading massive amounts of data from
	// memory mapped files. Strings longer than 0x100B should be rather rare.
	MaxStringLength = uint32(0x100)

)
// ImageBoundImportDescriptor represents the IMAGE_BOUND_IMPORT_DESCRIPTOR.
type ImageBoundImportDescriptor struct {
	TimeDateStamp uint32 // is just the value from the Exports information of the DLL which is being imported from.
	OffsetModuleName uint16 //  offset of the DLL name counted from the beginning of the BOUND_IMPORT table
	NumberOfModuleForwarderRefs uint16  // number of forwards
	// Array of zero or more IMAGE_BOUND_FORWARDER_REF follows
}

// ImageBoundForwardedRef represents the IMAGE_BOUND_FORWARDER_REF.
type ImageBoundForwardedRef struct {
	TimeDateStamp uint32
	OffsetModuleName uint16
	Reserved uint16
}

// BoundImportDescriptorData represents the descripts in addition to forwarded refs.
type BoundImportDescriptorData struct {
	Struct ImageBoundImportDescriptor
	Name string
	ForwardedRefs []BoundForwardedRefData
}

// BoundForwardedRefData reprents the struct in addition to the dll name.
type BoundForwardedRefData struct {
	Struct ImageBoundForwardedRef
	Name string
}

func (pe *File) parseBoundImportDirectory(rva, size uint32) (err error) {
	var sectionsAfterOffset []uint32
	var safetyBoundary uint32
	var start = rva;

	for {
		bndDesc := ImageBoundImportDescriptor{}
		bndDescSize := uint32(binary.Size(bndDesc))
		buf := bytes.NewReader(pe.data[rva : rva+bndDescSize])
		err := binary.Read(buf, binary.LittleEndian, &bndDesc)
		// If the RVA is invalid all would blow up. Some EXEs seem to be
		// specially nasty and have an invalid RVA.
		if err != nil {
			return err
		}

		// If the structure is all zeros, we reached the end of the list.
		if bndDesc == (ImageBoundImportDescriptor{}) {
			break
		}

		rva += bndDescSize
		sectionsAfterOffset = nil

		fileOffset := pe.getOffsetFromRva(rva)
		section := pe.getSectionByRva(rva)
		if section == nil {
			safetyBoundary = pe.size - fileOffset
			for _, section := range pe.Sections {
				if section.PointerToRawData > fileOffset {
					sectionsAfterOffset = append(sectionsAfterOffset, section.PointerToRawData)
				}
			}
			if len(sectionsAfterOffset) > 0 {
				// Find the first section starting at a later offset than that
				// specified by 'rva'
				firstSectionAfterOffset := Min(sectionsAfterOffset)
				section = pe.getSectionByOffset(firstSectionAfterOffset)
				if section != nil {
					safetyBoundary = section.PointerToRawData - fileOffset
				}
			}
		} else {
			sectionLen := uint32(len(section.Data(0, 0, pe)))
			safetyBoundary = (section.PointerToRawData + sectionLen) - fileOffset
		}

		if section == nil {
			log.Printf("RVA of IMAGE_BOUND_IMPORT_DESCRIPTOR points to an invalid address: 0x%x", rva)
			return nil
		}

		bndFrwdRef := ImageBoundForwardedRef{}
		bndFrwdRefSize := uint32(binary.Size(bndFrwdRef))
		count := min(uint32(bndDesc.NumberOfModuleForwarderRefs), safetyBoundary/bndFrwdRefSize)

		var forwarderRefs[]BoundForwardedRefData
		for i := uint32(0) ; i < count ;i++ {
			buf := bytes.NewReader(pe.data[rva : rva+bndFrwdRefSize])
			err := binary.Read(buf, binary.LittleEndian, &bndFrwdRef)
			if err != nil {
				return err
			}

			rva += bndFrwdRefSize

			offset := start+uint32(bndFrwdRef.OffsetModuleName)
			DllNameBuff := string(pe.getStringFromData(0, pe.data[offset:offset+MaxStringLength]))
			DllName := string(DllNameBuff)

			// OffsetModuleName points to a DLL name. These shouldn't be too long.
			// Anything longer than a safety length of 128 will be taken to indicate
			// a corrupt entry and abort the processing of these entries.
			// Names shorter than 4 characters will be taken as invalid as well.
			if DllName != "" &&  (len(DllName) > 256 || !IsPrintable(DllName)) {
				break
			}

			forwarderRefs = append(forwarderRefs, BoundForwardedRefData {
				Struct: bndFrwdRef, Name:DllName})
		}

		offset := start+uint32(bndDesc.OffsetModuleName)
		DllNameBuff := pe.getStringFromData(0, pe.data[offset:offset+MaxStringLength])
		DllName := string(DllNameBuff)
		if DllName != "" &&  (len(DllName) > 256 || !IsPrintable(DllName)) {
			break
		}

		pe.BoundImports = append(pe.BoundImports, BoundImportDescriptorData {
				Struct: bndDesc,
				Name: DllName,
				ForwardedRefs: forwarderRefs})
	}

	return nil
}