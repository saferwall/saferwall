// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// The following values are defined for the Type field of the debug directory entry:
const (
	ImageDebugTypeUnknown              = 0  // An unknown value that is ignored by all tools.
	ImageDebugTypeCOFF                 = 1  // The COFF debug information (line numbers, symbol table, and string table). This type of debug information is also pointed to by fields in the file headers.
	ImageDebugTypeCodeview             = 2  // The Visual C++ debug information.
	ImageDebugTypeFPO                  = 3  // The frame pointer omission (FPO) information. This information tells the debugger how to interpret nonstandard stack frames, which use the EBP register for a purpose other than as a frame pointer.
	ImageDebugTypeMisc                 = 4  // The location of DBG file.
	ImageDebugTypeException            = 5  // A copy of .pdata section.
	ImageDebugTypeFixup                = 6  // Reserved.
	ImageDebugTypeOmapToSrc            = 7  // The mapping from an RVA in image to an RVA in source image.
	ImageDebugTypeOmapFromSrc          = 8  // The mapping from an RVA in source image to an RVA in image.
	ImageDebugTypeBorland              = 9  // Reserved for Borland.
	ImageDebugTypeReserved10           = 10 // Reserved.
	ImageDebugTypeClsid                = 11 // Reserved.
	ImageDebugTypeVCFeature            = 12 // Visual C++ features (/GS counts /sdl counts and guardN counts)
	ImageDebugTypePOGO                 = 13 // Profile Guided Optimization
	ImageDebugTypeILTCG                = 14 // Incremental Link Time Code Generation (iLTCG)
	ImageDebugTypeMPX                  = 15 // Intel MPX
	ImageDebugTypeRepro                = 16 // PE determinism or reproducibility.
	ImageDebugTypeExDllCharacteristics = 20 // Extended DLL characteristics bits.
)

const (
	// CVSignatureRSDS represents the CodeView signature 'SDSR'
	CVSignatureRSDS = 0x53445352
	// CVSignatureNB10 represents the CodeView signature 'NB10'
	CVSignatureNB10 = 0x3031424e
)

// ImageDebugDirectory represents the IMAGE_DEBUG_DIRECTORY structure.
// This directory indicates what form of debug information is present
// and where it is. This directory consists of an array of debug directory
// entries whose location and size are indicated in the image optional header.
type ImageDebugDirectory struct {
	Characteristics  uint32 // Reserved, must be 0.
	TimeDateStamp    uint32 // The time and date that the debug data was created.
	MajorVersion     uint16 // The major version number of the debug data format.
	MinorVersion     uint16 // The minor version number of the debug data format.
	Type             uint32 // The format of debugging information. This field enables support of multiple debuggers.
	SizeOfData       uint32 // The size of the debug data (not including the debug directory itself).
	AddressOfRawData uint32 //The address of the debug data when loaded, relative to the image base.
	PointerToRawData uint32 // The file pointer to the debug data.
}

// DebugEntry wraps ImageDebugDirectory to include debug directory type.
type DebugEntry struct {
	Data ImageDebugDirectory // Points to the image debug entry structure.
	Type uint32              // Debug type.
	Info interface{}         // Holds specific information about the debug type entry.
}

// GUID is a 128-bit value consisting of one group of 8 hexadecimal digits,
// followed by three groups of 4 hexadecimal digits each, followed by one
//group of 12 hexadecimal digits.
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// CvInfoPDB70 represents the the CodeView data block of a PDB 7.0 file.
type CvInfoPDB70 struct {
	CvSignature uint32 // CodeView signature, equal to `RSDS`
	Signature   GUID   // A unique identifier, which changes with every rebuild of the executable and PDB file.
	Age         uint32 // Ever-incrementing value, which is initially set to 1 and incremented every time when a part of the PDB file is updated without rewriting the whole file.
	PDBFileName string // Null-terminated name of the PDB file. It can also contain full or partial path to the file.
}

// CVHeader represents the the CodeView header struct to the PDB 2.0 file.
type CVHeader struct {
	Signature uint32 // CodeView signature, equal to `NB10`
	Offset    uint32 // CodeView offset. Set to 0, because debug information is stored in a separate file.
}

// CvInfoPDB20 represents the the CodeView data block of a PDB 2.0 file.
type CvInfoPDB20 struct {
	CvHeader    CVHeader
	Signature   uint32 //The time when debug information was created (in seconds since 01.01.1970)
	Age         uint32 // Ever-incrementing value, which is initially set to 1 and incremented every time when a part of the PDB file is updated without rewriting the whole file.
	PDBFileName string // Null-terminated name of the PDB file. It can also contain full or partial path to the file.
}

// ImageDebugMisc represents the IMAGE_DEBUG_MISC structure.
type ImageDebugMisc struct {
	DataType uint32  // The type of data carried in the `Data` field.
	Length   uint32  // The length of this structure in bytes, including the entire Data field and its NUL terminator (rounded to four byte multiple.)
	Unicode  bool    // The encoding of the Data field. True if data is unicode string
	Reserved [3]byte // Reserved
	Data     string  // Actual data
}

func (pe *File) parseDebugDirectory(rva, size uint32) error {

	// Define some vars.
	debugDir := ImageDebugDirectory{}
	debugEntry := DebugEntry{}
	errorMsg := fmt.Sprintf("Invalid debug information. Can't read data at RVA: 0x%x", rva)
	debugDirSize := uint32(binary.Size(debugDir))
	debugDirsCount := size / debugDirSize

	for i := uint32(0); i < debugDirsCount; i++ {
		offset := pe.getOffsetFromRva(rva + debugDirSize*i)
		buf := bytes.NewReader(pe.data[offset : offset+debugDirSize])
		err := binary.Read(buf, binary.LittleEndian, &debugDir)
		if err != nil {
			return errors.New(errorMsg)
		}

		if debugDir.Type == ImageDebugTypeCodeview {
			debugSignature := binary.LittleEndian.Uint32(pe.data[debugDir.PointerToRawData:])
			if debugSignature == CVSignatureRSDS {
				// PDB 7.0
				pdb := CvInfoPDB70{CvSignature: debugSignature}

				// Guid
				offset := debugDir.PointerToRawData
				buff := bytes.NewReader(pe.data[offset+4 : offset+4+16])
				err = binary.Read(buff, binary.LittleEndian, &pdb.Signature)
				if err != nil {
					continue
				}
				// Age
				pdb.Age = binary.LittleEndian.Uint32(pe.data[offset+20:])

				// PDB file name
				pdbFilenameSize := debugDir.SizeOfData - 24
				// pdbFileName_size can be negative here, as seen in the malware sample with hash
				// MD5: 7c297600870d026c014d42596bb9b5fd
				// SHA256: 83f4e63681fcba8a9d7bbb1688c71981b1837446514a1773597e0192bba9fac3
				// Checking for positive size here to ensure proper parsing.
				if pdbFilenameSize > 0 {
					buff = bytes.NewReader(pe.data[offset+24 : offset+24+pdbFilenameSize])
					pdbFilename := make([]byte, pdbFilenameSize)
					err = binary.Read(buff, binary.LittleEndian, &pdbFilename)
					if err != nil {
						continue
					}
					pdb.PDBFileName = string(pdbFilename)
				}

				// Include these extra informations
				debugEntry.Info = pdb

			} else if debugSignature == CVSignatureNB10 {
				// PDB 2.0
				cvHeader := CVHeader{}
				offset := debugDir.PointerToRawData
				buf := bytes.NewReader(pe.data[offset : offset+8])
				err := binary.Read(buf, binary.LittleEndian, &cvHeader)
				if err != nil {
					continue
				}

				pdb := CvInfoPDB20{CvHeader: cvHeader}
				pdb.Signature = binary.LittleEndian.Uint32(pe.data[offset+8:])
				pdb.Age = binary.LittleEndian.Uint32(pe.data[offset+12:])
				pdbFilenameSize := debugDir.SizeOfData - 16
				if pdbFilenameSize > 0 {
					buff := bytes.NewReader(pe.data[offset+16 : offset+16+pdbFilenameSize])
					pdbFilename := make([]byte, pdbFilenameSize)
					err = binary.Read(buff, binary.LittleEndian, &pdbFilename)
					if err != nil {
						continue
					}
					pdb.PDBFileName = string(pdbFilename)
				}

				// Include these extra informations
				debugEntry.Info = pdb
			}
		} else if debugDir.Type == ImageDebugTypeMisc {
			break
		}

		debugEntry.Type = debugDir.Type
		debugEntry.Data = debugDir
		pe.Debugs = append(pe.Debugs, debugEntry)
	}

	return nil
}
