package pe

import (
	"errors"
	"log"
)

const (
	// TinyPESize On Windows XP (x32) the smallest PE executable is 97 bytes.
	TinyPESize = 97

	// FileAlignmentHardcodedValue represents the value which PointerToRawData
	// should be at least equal or bigger to, or it will be rounded to zero.
	// According to http://corkami.blogspot.com/2010/01/parce-que-la-planche-aura-brule.html
	// if PointerToRawData is less that 0x200 it's rounded to zero.
	FileAlignmentHardcodedValue = 0x200
)

// Errors
var (

	// ErrInvalidPESize is returned when the file size is less that the smallest PE file size possible.ErrImageOS2SignatureFound
	ErrInvalidPESize = errors.New("Not a PE file, smaller than tiny PE")

	// ErrDOSMagicNotFound is returned when file is potentially a ZM executable.
	ErrDOSMagicNotFound = errors.New("DOS Header magic not found")

	// ErrInvalidElfanewValue is returned when e_lfanew is larger than file size.
	ErrInvalidElfanewValue = errors.New("Invalid e_lfanew value, probably not a PE file")

	// ErrImageOS2SignatureFound is returned when signature is for a NE file.
	ErrImageOS2SignatureFound = errors.New("Not a valid PE signature. Probably a NE file")

	// ErrImageOS2LESignatureFound is returned when signature is for a LE file.
	ErrImageOS2LESignatureFound = errors.New("Not a valid PE signature. Probably an LE file")

	// ErrImageVXDSignatureFound is returned when signature is for a NX file.
	ErrImageVXDSignatureFound = errors.New("Not a valid PE signature. Probably an LX file")

	// ErrImageTESignatureFound is returned when signature is for a TE file.
	ErrImageTESignatureFound = errors.New("Not a valid PE signature. Probably a TE file")

	// ErrImageNtSignatureNotFound is returned when PE magic signature is not found.
	ErrImageNtSignatureNotFound = errors.New("Not a valid PE signature. Magic not found")

	// ErrImageNtOptionalHeaderMagicNotFound is returned when optional header magic is different from PE32/PE32+.
	ErrImageNtOptionalHeaderMagicNotFound = errors.New("Not a valid PE signature. Optional Header magic not found")

	// ErrImageBaseNotAligned is reported when the image base is not aligned to 64 K.
	ErrImageBaseNotAligned = errors.New("Corrupt PE file. Image base not aligned to 64 K")

	// ErrImageBaseOverflow is reported when the image base is larger than 80000000h/FFFF080000000000h in PE32/P32+.
	ErrImageBaseOverflow = errors.New("Corrupt PE file. Image base is overflow")

	// ErrInvalidMajorSubsystemVersion is reported when major subsystem version is less than 3.
	ErrInvalidMajorSubsystemVersion = errors.New("Corrupt PE file. Optional header major subsystem version is less than 3")

	// ErrInvalidSectionFileAlignment is reported when section alignment is less than a PAGE_SIZE and section alignement != file alignment.
	ErrInvalidSectionFileAlignment = errors.New("Corrupt PE file. Ssection alignment is less than a PAGE_SIZE and section alignement != file alignment")
)
