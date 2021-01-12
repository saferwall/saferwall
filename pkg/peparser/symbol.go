package pe

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

const (
	// Type Representation

	// ImageSymTypeNull indicates no type information or unknown base type.
	// Microsoft tools use this setting.
	ImageSymTypeNull = 0

	// ImageSymTypeVoid indicates no type no valid type; used with void pointers and functions.
	ImageSymTypeVoid = 1

	// ImageSymTypeChar indicates a character (signed byte).
	ImageSymTypeChar = 2

	// ImageSymTypeShort indicates a 2-byte signed integer.
	ImageSymTypeShort = 3

	// ImageSymTypeInt indicates a natural integer type (normally 4 bytes in
	// Windows).
	ImageSymTypeInt = 4

	// ImageSymTypeLong indicates a 4-byte signed integer.
	ImageSymTypeLong = 5

	// ImageSymTypeFloat indicates a 4-byte floating-point number.
	ImageSymTypeFloat = 6

	// ImageSymTypeDouble indicates an 8-byte floating-point number.
	ImageSymTypeDouble = 7

	// ImageSymTypeStruct indicates a structure.
	ImageSymTypeStruct = 8

	// ImageSymTypeUnion indicates a union.
	ImageSymTypeUnion = 9

	// ImageSymTypeEnum indicates an enumerated type.
	ImageSymTypeEnum = 10

	// ImageSymTypeMoe A member of enumeration (a specific value).
	ImageSymTypeMoe = 11

	// ImageSymTypeByte indicates a byte; unsigned 1-byte integer.
	ImageSymTypeByte = 12

	// ImageSymTypeWord indicates a word; unsigned 2-byte integer.
	ImageSymTypeWord = 13

	// ImageSymTypeUint indicates an unsigned integer of natural size
	// (normally, 4 bytes).
	ImageSymTypeUint = 14

	// ImageSymTypeDword indicates an unsigned 4-byte integer.
	ImageSymTypeDword = 15
)

var (
	errCOFFTableNotPresent   = errors.New("PE image does not countains a COFF symbol table")
	errNoCOFFStringInTable   = errors.New("PE image got a PointerToSymbolTable but no string in the COFF string table")
	errCOFFSymbolOutOfBounds = errors.New("COFF symbol offset out of bounds")
)

// COFFSymbol represents an entry in the COFF symbol table, which it is an
// array of records, each 18 bytes long. Each record is either a standard or
// auxiliary symbol-table record. A standard record defines a symbol or name
// and has the following format.
type COFFSymbol struct {
	// The name of the symbol, represented by a union of three structures. An
	// array of 8 bytes is used if the name is not more than 8 bytes long.
	// union {
	//    BYTE     ShortName[8];
	//    struct {
	//        DWORD   Short;     // if 0, use LongName
	//        DWORD   Long;      // offset into string table
	//    } Name;
	//    DWORD   LongName[2];    // PBYTE  [2]
	// } N;
	Name [8]byte

	// The value that is associated with the symbol. The interpretation of this
	// field depends on SectionNumber and StorageClass. A typical meaning is
	// the relocatable address.
	Value uint32

	// The signed integer that identifies the section, using a one-based index
	// into the section table. Some values have special meaning, as defined in section 5.4.2, "Section Number Values."
	SectionNumber int16

	// A number that represents type. Microsoft tools set this field to 0x20 (function) or 0x0 (not a function). For more information, see Type Representation.
	Type uint16

	// An enumerated value that represents storage class. For more information, see Storage Class.
	StorageClass uint8

	// The number of auxiliary symbol table entries that follow this record.
	NumberOfAuxSymbols uint8
}

// COFF holds properties related to the COFF format.
type COFF struct {
	SymbolTable       []COFFSymbol
	StringTable       []string
	StringTableOffset uint32
	StringTableM      map[uint32]string // Map the symbol offset => symbol name.
}

// ParseCOFFSymbolTable parses the COFF symbol table. The symbol table is
// inherited from the traditional COFF format. It is distinct from Microsoft
// Visual C++ debug information. A file can contain both a COFF symbol table
// and Visual C++ debug information, and the two are kept separate. Some
// Microsoft tools use the symbol table for limited but important purposes,
// such as communicating COMDAT information to the linker. Section names and
// file names, as well as code and data symbols, are listed in the symbol table.
func (pe *File) ParseCOFFSymbolTable() error {
	pointerToSymbolTable := pe.NtHeader.FileHeader.PointerToSymbolTable
	if pointerToSymbolTable == 0 {
		return errCOFFTableNotPresent
	}

	size := uint32(binary.Size(COFFSymbol{}))
	symCount := pe.NtHeader.FileHeader.NumberOfSymbols
	offset := pe.NtHeader.FileHeader.PointerToSymbolTable
	symbols := make([]COFFSymbol, symCount)

	for i := uint32(0); i < symCount; i++ {
		err := pe.structUnpack(&symbols[i], offset, size)
		if err != nil {
			return err
		}
		offset += size
	}
	pe.COFF.SymbolTable = symbols

	// Get the COFF string table.
	pe.COFFStringTable()

	return nil
}

// COFFStringTable retrieves the list of strings in the COFF string table if
// any.
func (pe *File) COFFStringTable() error {
	m := make(map[uint32]string)
	pointerToSymbolTable := pe.NtHeader.FileHeader.PointerToSymbolTable
	if pointerToSymbolTable == 0 {
		return errCOFFTableNotPresent
	}

	// COFF String Table immediately following the COFF symbol table. The
	// position of this table is found by taking the symbol table address in
	// the COFF header and adding the number of symbols multiplied by the size
	// of a symbol.
	size := uint32(binary.Size(COFFSymbol{}))
	symCount := pe.NtHeader.FileHeader.NumberOfSymbols
	offset := pointerToSymbolTable + (size * symCount)
	pe.COFF.StringTableOffset = offset

	// At the beginning of the COFF string table are 4 bytes that contain the
	// total size (in bytes) of the rest of the string table. This size
	// includes the size field itself, so that the value in this location would
	// be 4 if no strings were present.
	strTableSize, err := pe.ReadUint32(offset)
	if err != nil {
		return err
	}
	if strTableSize <= 4 {
		return errNoCOFFStringInTable
	}
	offset += 4

	// Following the size are null-terminated strings that are pointed to by
	// symbols in the COFF symbol table. We create a map to map offset to
	// string.
	end := offset + strTableSize
	for offset <= end {
		len, str := pe.readASCIIStringAtOffset(offset, 0x30)
		if len == 0 {
			break
		}
		m[offset] = str
		offset += len + 1
		pe.COFF.StringTable = append(pe.COFF.StringTable, str)
	}

	pe.COFF.StringTableM = m
	return nil
}

func (symbol *COFFSymbol) String(pe *File) (string, error) {
	// contain the name itself, if it is not more than 8 bytes long, or the
	// ShortName field gives an offset into the string table. To determine
	// whether the name itself or an offset is given, test the first 4
	// bytes for equality to zero.
	var short, long uint32
	highDw := bytes.NewBuffer(symbol.Name[4:])
	lowDw := bytes.NewBuffer(symbol.Name[:4])
	errl := binary.Read(lowDw, binary.LittleEndian, &short)
	errh := binary.Read(highDw, binary.LittleEndian, &long)
	if errl != nil || errh != nil {
		return "", errCOFFSymbolOutOfBounds
	}

	// if 0, use LongName.
	if short != 0 {
		name := strings.Replace(string(symbol.Name[:]), "\x00", "", -1)
		return name, nil
	}

	// Long name offset to the string table.
	strOff := pe.COFF.StringTableOffset + long
	name := pe.COFF.StringTableM[strOff]
	return name, nil
}

// PrettyCOFFTypeRepresentation returns the string representation of the `Type`
// field of a COFF table entry.
func (pe *File) PrettyCOFFTypeRepresentation(k uint8) string {
	coffSymTypeMap := map[uint8]string{
		ImageSymTypeNull:   "Null",
		ImageSymTypeVoid:   "Void",
		ImageSymTypeChar:   "Char",
		ImageSymTypeShort:  "Short",
		ImageSymTypeInt:    "Int",
		ImageSymTypeLong:   "Long",
		ImageSymTypeFloat:  "Float",
		ImageSymTypeDouble: "Double",
		ImageSymTypeStruct: "Struct",
		ImageSymTypeUnion:  "Union",
		ImageSymTypeEnum:   "Enum",
		ImageSymTypeMoe:    "Moe",
		ImageSymTypeByte:   "Byte",
		ImageSymTypeWord:   "Word",
		ImageSymTypeUint:   "Uint",
		ImageSymTypeDword:  "Dword",
	}

	if value, ok := coffSymTypeMap[k]; ok {
		return value
	}
	return ""
}
