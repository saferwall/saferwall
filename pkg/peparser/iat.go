// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"encoding/binary"
)

// IATEntry represents an entry inside the IAT.
type IATEntry struct {
	Index   uint32
	Rva     uint32
	Value   interface{}
	Meaning string
}

// The structure and content of the import address table are identical to those
// of the import lookup table, until the file is bound. During binding, the entries
//  in the import address table are overwritten with the 32-bit (for PE32) or
// 64-bit (for PE32+) addresses of the symbols that are being imported.
// These addresses are the actual memory addresses of the symbols, although
// technically they are still called “virtual addresses.” The loader typically
// processes the binding.
//
// The Import Address Table is there to to only trigger Copy On Write for as
// few pages as possible (those being the actual Import Address Table pages themselves).
// This is, partially the reason there's that extra level of indirection in the PE to begin with.
func (pe *File) parseIATDirectory(rva, size uint32) error {

	var entries []IATEntry
	var index uint32

	startRva := rva

	for startRva+size > rva {
		ie := IATEntry{}
		offset := pe.getOffsetFromRva(rva)
		if pe.Is64 {
			ie.Value = binary.LittleEndian.Uint64(pe.data[offset:])
			ie.Rva = rva
			rva += 8
		} else {
			ie.Value = binary.LittleEndian.Uint32(pe.data[offset:])
			ie.Rva = rva
			rva += 4
		}
		ie.Index = index
		entries = append(entries, ie)
		index++
	}

	return nil
}
