package pe

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"sort"

	mmap "github.com/edsrzf/mmap-go"
)

// A File represents an open PE file.
type File struct {
	DosHeader        ImageDosHeader

	data mmap.MMap
}

// Open opens the named file using os.Open and prepares it for use as a PE binary.
func Open(name string) (File, error) {

	// Init an File instance
	file := File{}

	f, err := os.Open(name)
	if err != nil {
		return file, err
	}

	// Memory map the file insead of using read/write.
	data, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		f.Close()
		return file, err
	}

	file.data = data
	return file, nil
}



// Parse performs the file parsing for a PE binary.
func (pe *File) Parse() error {

	// Probes for the smallest PE size
	if len(pe.data) < TinyPESize {
		return ErrInvalidPESize
	}

	// Parse the DOS header
	err := pe.parseDosHeader()
	if err != nil {
		return err
	}

	// Parse the NT header
	err = pe.parseNtHeader()
	if err != nil {
		return err
	}

	return nil
}

func (pe *File) parseDosHeader() (err error) {
	offset := 0
	size := binary.Size(pe.DosHeader)
	buf := bytes.NewReader(pe.data[offset : offset+size])
	err = binary.Read(buf, binary.LittleEndian, &pe.DosHeader)
	if err != nil {
		return err
	}

	// it can be ZM on an (non-PE) EXE.
	// These executables still work under XP via ntvdm.
	if pe.DosHeader.Emagic != ImageDOSSignature && pe.DosHeader.Emagic != ImageDOSZMSignature {
		return ErrDOSMagicNotFound
	}

	// e_lfanew  is the only required element (besides the signature) of the DOS HEADER to turn the EXE into a PE.
	// is a relative offset to the NT Headers.
	// can't be null (signatures would overlap)
	// can be 4 at minimum
	if pe.DosHeader.Elfanew < 4 || pe.DosHeader.Elfanew > uint32(len(pe.data)) {
		return ErrInvalidElfanewValue
	}

	return nil
}


func (pe *File) parseNtHeader() (err error) {
	ntHeaderOffset := pe.DosHeader.Elfanew
	size := uint32(binary.Size(pe.NtHeader))
	buf := bytes.NewReader(pe.data[ntHeaderOffset : ntHeaderOffset+size])
	err = binary.Read(buf, binary.LittleEndian, &pe.NtHeader)
	if err != nil {
		return err
	}

	// Probe for PE signature.
	if pe.NtHeader.Signature == ImageOS2Signature {
		return ErrImageOS2SignatureFound
	}
	if pe.NtHeader.Signature == ImageOS2LESignature {
		return ErrImageOS2LESignatureFound
	}
	if pe.NtHeader.Signature == ImageVXDignature {
		return ErrImageVXDSignatureFound
	}
	if pe.NtHeader.Signature == ImageTESignature {
		return ErrImageTESignatureFound
	}

	// this is the smallest requirement for a valid PE
	if pe.NtHeader.Signature != ImageNTSignature {
		return ErrImageNtSignatureNotFound
	}

	return nil
}