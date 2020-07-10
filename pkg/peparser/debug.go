// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// The following values are defined for the Type field of the debug directory entry:
const (
	// An unknown value that is ignored by all tools.
	ImageDebugTypeUnknown = 0

	// The COFF debug information (line numbers, symbol table, and string table).
	// This type of debug information is also pointed to by fields in the file headers.
	ImageDebugTypeCOFF = 1

	// The Visual C++ debug information.
	ImageDebugTypeCodeview = 2

	// The frame pointer omission (FPO) information. This information tells the
	// debugger how to interpret nonstandard stack frames, which use the EBP
	// register for a purpose other than as a frame pointer.
	ImageDebugTypeFPO = 3

	// The location of DBG file.
	ImageDebugTypeMisc = 4

	// A copy of .pdata section.
	ImageDebugTypeException = 5

	// Reserved.
	ImageDebugTypeFixup = 6

	// The mapping from an RVA in image to an RVA in source image.
	ImageDebugTypeOmapToSrc = 7

	// The mapping from an RVA in source image to an RVA in image.
	ImageDebugTypeOmapFromSrc = 8

	// Reserved for Borland.
	ImageDebugTypeBorland = 9

	// Reserved.
	ImageDebugTypeReserved10 = 10

	// Reserved.
	ImageDebugTypeClsid = 11

	// Visual C++ features (/GS counts /sdl counts and guardN counts)
	ImageDebugTypeVCFeature = 12

	// Profile Guided Optimization
	ImageDebugTypePOGO = 13

	// Incremental Link Time Code Generation (iLTCG)
	ImageDebugTypeILTCG = 14

	// Intel MPX
	ImageDebugTypeMPX = 15

	// PE determinism or reproducibility.
	ImageDebugTypeRepro = 16

	// Extended DLL characteristics bits.
	ImageDebugTypeExDllCharacteristics = 20
)

const (
	// CVSignatureRSDS represents the CodeView signature 'SDSR'.
	CVSignatureRSDS = 0x53445352

	// CVSignatureNB10 represents the CodeView signature 'NB10'.
	CVSignatureNB10 = 0x3031424e
)

const (
	// FrameFPO indicates a frame of type FPO.
	FrameFPO = 0x0

	// FrameTrap indicates a frame of type Trap.
	FrameTrap = 0x1

	// FrameTSS indicates a frame of type TSS.
	FrameTSS = 0x2

	// FrameNonFPO indicates a frame of type Non-FPO.
	FrameNonFPO = 0x3
)

const (
	POGOTypePGU  = 0x50475500
	POGzOTypePGI = 0x50474900
	POGOTypePGO  = 0x50474F00
	POGOTypeLTCG = 0x4c544347
)

// ImageDebugDirectory represents the IMAGE_DEBUG_DIRECTORY structure.
// This directory indicates what form of debug information is present
// and where it is. This directory consists of an array of debug directory
// entries whose location and size are indicated in the image optional header.
type ImageDebugDirectory struct {
	// Reserved, must be 0.
	Characteristics uint32

	// The time and date that the debug data was created.
	TimeDateStamp uint32

	// The major version number of the debug data format.
	MajorVersion uint16

	// The minor version number of the debug data format.
	MinorVersion uint16

	// The format of debugging information. This field enables support of
	// multiple debuggers.
	Type uint32

	// The size of the debug data (not including the debug directory itself).
	SizeOfData uint32

	//The address of the debug data when loaded, relative to the image base.
	AddressOfRawData uint32

	// The file pointer to the debug data.
	PointerToRawData uint32
}

// DebugEntry wraps ImageDebugDirectory to include debug directory type.
type DebugEntry struct {
	// Points to the image debug entry structure.
	Struct ImageDebugDirectory

	// Holds specific information about the debug type entry.
	Info interface{}
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
	// CodeView signature, equal to `RSDS`
	CvSignature uint32

	// A unique identifier, which changes with every rebuild of the executable and PDB file.
	Signature GUID

	// Ever-incrementing value, which is initially set to 1 and incremented every
	// time when a part of the PDB file is updated without rewriting the whole file.
	Age uint32

	// Null-terminated name of the PDB file. It can also contain full or partial
	// path to the file.
	PDBFileName string
}

// CVHeader represents the the CodeView header struct to the PDB 2.0 file.
type CVHeader struct {
	// CodeView signature, equal to `NB10`
	Signature uint32

	// CodeView offset. Set to 0, because debug information is stored in a separate file.
	Offset uint32
}

// CvInfoPDB20 represents the the CodeView data block of a PDB 2.0 file.
type CvInfoPDB20 struct {
	// Points to the CodeView header structure.
	CvHeader CVHeader

	// The time when debug information was created (in seconds since 01.01.1970)
	Signature uint32

	// Ever-incrementing value, which is initially set to 1 and incremented every
	// time when a part of the PDB file is updated without rewriting the whole file.
	Age uint32

	// Null-terminated name of the PDB file. It can also contain full or partial
	// path to the file.
	PDBFileName string
}

// FPOData Represents the stack frame layout for a function on an x86 computer when frame pointer omission (FPO) optimization is used. The structure is used to locate the base of the call frame.
type FPOData struct {
	// The offset of the first byte of the function code.
	OffStart uint32

	// The number of bytes in the function.
	ProcSize uint32

	// The number of local variables.
	NumLocals uint32

	// The size of the parameters, in DWORDs.
	NunParams uint32

	// The number of bytes in the function prolog code.
	NumProlog uint32

	// The number of registers saved.
	NumRegs uint8

	// A variable that indicates whether the function uses structured exception handling.
	HasSEH uint8

	// A variable that indicates whether the EBP register has been allocated.
	UseBP uint8

	// Reserved for future use.
	Reserved uint8

	// A variable that indicates the frame type.
	FrameType uint8
}

type ImagePGOItem struct {
	Rva  uint32
	Size uint32
	Name string
}

type POGO struct {
	Signature uint32 // _IMAGE_POGO_INFO
	Entries   []ImagePGOItem
}

type VCFeature struct {
	PreVC11 uint32 `json:"Pre VC 11"`
	CCpp uint32 `json:"C/C++"`
	Gs uint32 `json:"/GS"`
	Sdl uint32 `json:"/sdl"`
	GuardN uint32
}

// ImageDebugMisc represents the IMAGE_DEBUG_MISC structure.
type ImageDebugMisc struct {
	DataType uint32 // The type of data carried in the `Data` field.

	// The length of this structure in bytes, including the entire Data field
	// and its NUL terminator (rounded to four byte multiple.)
	Length uint32

	// The encoding of the Data field. True if data is unicode string
	Unicode bool

	// Reserved
	Reserved [3]byte

	// Actual data
	Data string
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
		err := pe.structUnpack(&debugDir, offset, debugDirSize)
		if err != nil {
			return errors.New(errorMsg)
		}

		switch debugDir.Type {
		case ImageDebugTypeCodeview:
			debugSignature, err := pe.ReadUint32(debugDir.PointerToRawData)
			if err != nil {
				continue
			}

			if debugSignature == CVSignatureRSDS {
				// PDB 7.0
				pdb := CvInfoPDB70{CvSignature: CVSignatureRSDS}

				// Guid
				offset := debugDir.PointerToRawData + 4
				guidSize := uint32(binary.Size(pdb.Signature))
				err = pe.structUnpack(&pdb.Signature, offset, guidSize)
				if err != nil {
					continue
				}
				// Age
				offset += guidSize
				pdb.Age, err = pe.ReadUint32(offset)
				if err != nil {
					continue
				}
				offset += 4

				// PDB file name
				pdbFilenameSize := debugDir.SizeOfData - 24 - 1

				// pdbFileName_size can be negative here, as seen in the malware
				// sample with MD5 hash: 7c297600870d026c014d42596bb9b5fd
				// Checking for positive size here to ensure proper parsing.
				if pdbFilenameSize > 0 {
					pdbFilename := make([]byte, pdbFilenameSize)
					err = pe.structUnpack(&pdbFilename, offset, pdbFilenameSize)
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
				err = pe.structUnpack(&cvHeader, offset, size)
				if err != nil {
					continue
				}

				pdb := CvInfoPDB20{CvHeader: cvHeader}

				// Signature
				pdb.Signature, err = pe.ReadUint32(offset + 8)
				if err != nil {
					continue
				}

				// Age
				pdb.Age, err = pe.ReadUint32(offset + 12)
				if err != nil {
					continue
				}
				offset += 16

				pdbFilenameSize := debugDir.SizeOfData - 16 - 1
				if pdbFilenameSize > 0 {
					pdbFilename := make([]byte, pdbFilenameSize)
					err = pe.structUnpack(&pdbFilename, offset, pdbFilenameSize)
					if err != nil {
						continue
					}
					pdb.PDBFileName = string(pdbFilename)
				}

				// Include these extra informations
				debugEntry.Info = pdb
			}
		case ImageDebugTypePOGO:
			pogoSignature, err := pe.ReadUint32(debugDir.PointerToRawData)
			if err != nil {
				continue
			}

			pogo := POGO{}

			switch pogoSignature {
			case POGOTypePGU:
			case POGzOTypePGI:
			case POGOTypePGO:
			case POGOTypeLTCG:
				pogo.Signature = pogoSignature
				offset = debugDir.PointerToRawData + 4
				c := uint32(0)
				for c < debugDir.SizeOfData-4 {

					pogoEntry := ImagePGOItem{}
					pogoEntry.Rva, err = pe.ReadUint32(offset)
					if err != nil {
						break
					}
					pogoEntry.Size, err = pe.ReadUint32(offset + 4)
					if err != nil {
						break
					}

					pogoEntry.Name = string(pe.getStringFromData(0, pe.data[offset+8:offset+8+32]))

					pogo.Entries = append(pogo.Entries, pogoEntry)
					c += 8 + uint32(len(pogoEntry.Name)) + 4
					offset += 8 + uint32(len(pogoEntry.Name)) + 4
				}

				debugEntry.Info = pogo
			}

		case ImageDebugTypeVCFeature:
			vcf := VCFeature{}
			size := uint32(binary.Size(vcf))
			err := pe.structUnpack(&vcf, debugDir.PointerToRawData, size)
			if err != nil {
				continue
			}
			debugEntry.Info = vcf
		}

		debugEntry.Struct = debugDir
		pe.Debugs = append(pe.Debugs, debugEntry)
	}

	return nil
}

// SectionAttributeDescription maps a section attribute to a friendly name.
func SectionAttributeDescription(section string) string {
	sectionNameMap := map[string]string{
		".CRT$XCA":      "First C++ Initializer",
		".CRT$XCAA":     "Startup C++ Initializer",
		".CRT$XCZ":      "Last C++ Initializer",
		".CRT$XDA":      "First Dynamic TLS Initializer",
		".CRT$XDZ":      "Last Dynamic TLS Initializer",
		".CRT$XIA":      "First C Initializer",
		".CRT$XIAA":     "Startup C Initializer",
		".CRT$XIAB":     "PGO C Initializer",
		".CRT$XIAC":     "Post-PGO C Initializer",
		".CRT$XIC":      "CRT C Initializers",
		".CRT$XIYA":     "VCCorLib Threading Model Initializer",
		".CRT$XIYAA":    "XAML Designer Threading Model Override Initializer",
		".CRT$XIYB":     "VCCorLib Main Initializer",
		".CRT$XIZ":      "Last C Initializer",
		".CRT$XLA":      "First Loader TLS Callback",
		".CRT$XLC":      "CRT TLS Constructor",
		".CRT$XLD":      "CRT TLS Terminator",
		".CRT$XLZ":      "Last Loader TLS Callback",
		".CRT$XPA":      "First Pre-Terminator",
		".CRT$XPB":      "CRT ConcRT Pre-Terminator",
		".CRT$XPX":      "CRT Pre-Terminators",
		".CRT$XPXA":     "CRT stdio Pre-Terminator",
		".CRT$XPZ":      "Last Pre-Terminator",
		".CRT$XTA":      "First Terminator",
		".CRT$XTZ":      "Last Terminator",
		".CRTMA$XCA":    "First Managed C++ Initializer",
		".CRTMA$XCZ":    "Last Managed C++ Initializer",
		".CRTVT$XCA":    "First Managed VTable Initializer",
		".CRTVT$XCZ":    "Last Managed VTable Initializer",
		".rtc$IAA":      "First RTC Initializer",
		".rtc$IZZ":      "Last RTC Initializer",
		".rtc$TAA":      "First RTC Terminator",
		".rtc$TZZ":      "Last RTC Terminator",
		".text$x":       "EH Filters",
		".text$di":      "MSVC Dynamic Initializers",
		".text$yd":      "MSVC Destructors",
		".text$mn":      "Contains EP",
		".00cfg":        "CFG Check Functions Pointers",
		".rdata$T":      "TLS Header",
		".rdata$r":      "RTTI Data",
		".data$r":       "RTTI Type Descriptors",
		".rdata$sxdata": "Safe SEH",
		".rdata$zzzdbg": "Debug Data",
		".idata$2":      "Import Descriptors",
		".idata$3":      "Final Null Entry",
		".idata$4":      "INT Array",
		".idata$5":      "IAT Array",
		".idata$6":      "Symbol and DLL names",
		".rsrc$01":      "Resources Header",
		".rsrc$02":      "Resources Data",
	}

	if val, ok := sectionNameMap[section]; ok {
		return val
	}

	return "?"
}
