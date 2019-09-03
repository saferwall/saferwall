package pe

import (
	"bytes"
	"encoding/binary"
)

const (
	ImageRelBasedAbsolute      = 0  // The base relocation is skipped. This type can be used to pad a block.
	ImageRelBasedHigh          = 1  // The base relocation adds the high 16 bits of the difference to the 16-bit field at offset. The 16-bit field represents the high value of a 32-bit word.
	ImageRelBasedLow           = 2  // The base relocation adds the low 16 bits of the difference to the 16-bit field at offset. The 16-bit field represents the low half of a 32-bit word.
	ImageRelBasedHighLow       = 3  // The base relocation applies all 32 bits of the difference to the 32-bit field at offset.
	ImageRelBasedHighAdj       = 4  // The base relocation adds the high 16 bits of the difference to the 16-bit field at offset. The 16-bit field represents the high value of a 32-bit word. The low 16 bits of the 32-bit value are stored in the 16-bit word that follows this base relocation. This means that this base relocation occupies two slots.
	ImageRelBasedMipsJmpAddr   = 5  // The relocation interpretation is dependent on the machine type. When the machine type is MIPS, the base relocation applies to a MIPS jump instruction.
	ImageRelBasedArmMov32      = 5  // This relocation is meaningful only when the machine type is ARM or Thumb. The base relocation applies the 32-bit address of a symbol across a consecutive MOVW/MOVT instruction pair.
	ImageRelBasedRiscvHigh20   = 5  // This relocation is only meaningful when the machine type is RISC-V. The base relocation applies to the high 20 bits of a 32-bit absolute address.
	ImageRelReserved           = 6  // Reserved, must be zero.
	ImageRelBasedThumbMov32    = 7  // This relocation is meaningful only when the machine type is Thumb. The base relocation applies the 32-bit address of a symbol to a consecutive MOVW/MOVT instruction pair.
	ImageRelBasedRiscvLow12i   = 7  // This relocation is only meaningful when the machine type is RISC-V. The base relocation applies to the low 12 bits of a 32-bit absolute address formed in RISC-V I-type instruction format.
	ImageRelBasedRiscvLow12s   = 8  // This relocation is only meaningful when the machine type is RISC-V. The base relocation applies to the low 12 bits of a 32-bit absolute address formed in RISC-V S-type instruction format.
	ImageRelBasedMipsJmpAddr16 = 9  // The relocation is only meaningful when the machine type is MIPS. The base relocation applies to a MIPS16 jump instruction.
	ImageRelBasedDir64         = 10 // The base relocation applies the difference to the 64-bit field at offset.
)

// ImageBaseRelocation
type ImageBaseRelocation struct {
	VirtualAddress uint32 // The image base plus the page RVA is added to each offset to create the VA where the base relocation must be applied.
	SizeOfBlock    uint32 // The total number of bytes in the base relocation block, including the Page RVA and Block Size fields and the Type/Offset fields that follow.

}

type Relocation struct {
	Data ImageBaseRelocation // Pointd to the ImageBaseRelocation structure.
	//Entries
}

func (pe *File) parseRelocDirectory(rva, size uint32) error {

	relocSize := uint32(binary.Size(ImageBaseRelocation{}))
	end := rva + size
	for rva < end {
		baseReloc := ImageBaseRelocation{}
		offset := pe.getOffsetFromRva(rva)
		buff := bytes.NewReader(pe.data[offset : offset+relocSize])
		err := binary.Read(buff, binary.LittleEndian, &baseReloc)
		if err != nil {
			return err
		}

		pe.Relocations = append(pe.Relocations, Relocation{
			Data: baseReloc,
		})
	}

	return nil

}
	}

	return relocEntries, nil
}

func (pe *File) parseRelocDirectory(rva, size uint32) error {

	relocSize := uint32(binary.Size(ImageBaseRelocation{}))
	end := rva + size
	for rva < end {
		baseReloc := ImageBaseRelocation{}
		offset := pe.getOffsetFromRva(rva)
		buff := bytes.NewReader(pe.data[offset : offset+relocSize])
		err := binary.Read(buff, binary.LittleEndian, &baseReloc)
		if err != nil {
			return err
		}

		// VirtualAddress must lie within the Image
		if baseReloc.VirtualAddress > pe.OptionalHeader.SizeOfImage {
			return errors.New("Invalid relocation information. VirtualAddress outside of Image")
		}

		// SizeOfBlock must be less or equal than the size of the image
		// (It's a rather loose sanity test)
		if baseReloc.SizeOfBlock > pe.OptionalHeader.SizeOfImage {
			return errors.New("Invalid relocation information. SizeOfBlock too large")
		}

		relocEntries, err := pe.parseRelocations(rva+relocSize, baseReloc.VirtualAddress, baseReloc.SizeOfBlock-relocSize)
		if err != nil {
			return err
		}

		pe.Relocations = append(pe.Relocations, Relocation{
			Data: baseReloc,
			Entries: relocEntries,
		})

		if baseReloc.SizeOfBlock == 0 {
			break
		}
		rva += baseReloc.SizeOfBlock
	}

	return nil

}
