package pe


import (
	"encoding/binary"
)

// This data directory is set to all zeros if the target architecture 
// (for example, I386 or AMD64) does not use the concept of a global pointer.
func (pe *File)parseGlobalPtrDirectory(rva, size uint32) error {

	// RVA of the value to be stored in the global pointer register. 
	offset := pe.getOffsetFromRva(rva)
	pe.GlobalPtr = binary.LittleEndian.Uint32(pe.data[offset:])
	
	// The size must be set to 0.
	return nil
}