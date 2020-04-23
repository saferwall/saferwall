// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

const (
	// DansSignature ('DanS' as dword) is where the rich header struct starts.
	DansSignature = 0x536E6144

	// RichSignature ('0x68636952' as dword) is where the rich header struct ends.
	RichSignature = "Rich"
)

// CompID represents the `@comp.id` structure.
type CompID struct {
	// The minor version information for the compiler used when building the product.
	MinorCV uint16

	// Provides information about the identity or type of the objects used to
	// build the PE32.
	ProdID uint16

	// Indicates how often the object identified by the former two fields is
	// referenced by this PE32 file.
	Count uint32
}

// RichHeader is a structure that is written right after the MZ DOS header.
// It consists of pairs of 4-byte integers. And it is also
// encrypted using a simple XOR operation using the checksum as the key.
// The data between the magic values encodes the ‘bill of materials’ that were
// collected by the linker to produce the binary.
type RichHeader struct {
	XorKey  uint32
	CompIDs []CompID
	Raw     []byte
}

// ParseRichHeader parses the rich header struct.
func (pe *File) ParseRichHeader() error {

	rh := RichHeader{}
	fileHeaderOffset := pe.DosHeader.AddressOfNewEXEHeader + uint32(binary.Size(pe.NtHeader))
	richSigOffset := bytes.Index(pe.data[:fileHeaderOffset], []byte(RichSignature))

	// For example, .NET executable files do not use the MSVC linker and these
	// executables do not contain a detectable Rich Header.
	if richSigOffset < 0 {
		log.Print("Rich header not found")
	}

	// the DWORD following the "Rich" sequence is the XOR key stored by and
	// calculated by the linker. It is actually a checksum of the DOS header with
	// the e_lfanew zeroed out, and additionally includes the values of the
	// unencrypted "Rich" array. Using a checksum with encryption will not only
	// obfuscate the values, but it also serves as a rudimentary digital
	// signature. If the checksum is calculated from scratch once the values
	// have been decrypted, but doesn't match the stored key, it can be assumed
	// the structure had been tampered with. For those that go the extra step to
	// recalculate the checksum/key, this simple protection mechanism can be bypassed.
	rh.XorKey = binary.LittleEndian.Uint32(pe.data[richSigOffset+4:])

	// To decrypt the array, start with the DWORD just prior to the `Rich` sequence
	// and XOR it with the key. Continue the loop backwards, 4 bytes at a time,
	// until the sequence `DanS` is decrypted.
	var decRichHeader []uint32
	dansSigOffset := -1
	for it := 0; it < 0x100; it += 4 {
		buff := binary.LittleEndian.Uint32(pe.data[richSigOffset-4-it:])
		res := buff ^ rh.XorKey
		if res == DansSignature {
			dansSigOffset = richSigOffset - it - 4
			break
		}

		decRichHeader = append(decRichHeader, res)
	}

	// Probe we successfuly found the `DanS` magic.
	if dansSigOffset == -1 {
		return errors.New("Rich Header Found, but could not locate DanS Signature")
	}

	// Anomaly check: dansSigOffset is usually found in offset 0x80.
	if dansSigOffset != 0x80 {
		pe.Anomalies = append(pe.Anomalies, AnoDanSMagicOffset)
	}

	rh.Raw = pe.data[dansSigOffset : richSigOffset+8]

	// reverse the decrypted rich header
	for i, j := 0, len(decRichHeader)-1; i < j; i, j = i+1, j-1 {
		decRichHeader[i], decRichHeader[j] = decRichHeader[j], decRichHeader[i]
	}

	// After the `DanS` signature, there are some zero-padded In practice,
	// Microsoft seems to have wanted the entries to begin on a 16-byte
	// (paragraph) boundary, so the 3 leading padding DWORDs can be safely
	// skipped as not belonging to the data.
	if decRichHeader[0] != 0 || decRichHeader[1] != 0 || decRichHeader[2] != 0 {
		return errors.New("Rich header: 3 leading padding DWORDs not not found")
	}

	// The array stores entries that are 8-bytes each, broken into 3 members.
	// Each entry represents either a tool that was employed as part of building
	// the executable or a statistic.
	lenCompIDs := len(decRichHeader)
	for i := 3; i < lenCompIDs; i += 2 {
		cid := CompID{}
		compid := make([]byte, binary.Size(cid))
		binary.LittleEndian.PutUint32(compid, decRichHeader[i])
		binary.LittleEndian.PutUint32(compid[4:], decRichHeader[i+1])
		buf := bytes.NewReader(compid)
		err := binary.Read(buf, binary.LittleEndian, &cid)
		if err != nil {
			return err
		}

		rh.CompIDs = append(rh.CompIDs, cid)
	}

	pe.RichHeader = rh
	return nil
}
