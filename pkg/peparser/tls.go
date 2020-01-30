// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
)

// TLSDirectory represents tls directory information with callback entries.
type TLSDirectory struct {
	Struct    interface{} // of type *IMAGE_TLS_DIRECTORY32 or *IMAGE_TLS_DIRECTORY64 structure.
	Callbacks interface{} // of type uint32 or uint64.
}

// ImageTLSDirectory32 represents the IMAGE_TLS_DIRECTORY32 structure.
// It Points to the Thread Local Storage initialization section.
type ImageTLSDirectory32 struct {
	StartAddressOfRawData uint32 // The starting address of the TLS template. The template is a block of data that is used to initialize TLS data.
	EndAddressOfRawData   uint32 // The address of the last byte of the TLS, except for the zero fill. As with the Raw Data Start VA field, this is a VA, not an RVA.
	AddressOfIndex        uint32 // The location to receive the TLS index, which the loader assigns. This location is in an ordinary data section, so it can be given a symbolic name that is accessible to the program.
	AddressOfCallBacks    uint32 // The pointer to an array of TLS callback functions. The array is null-terminated, so if no callback function is supported, this field points to 4 bytes set to zero.
	SizeOfZeroFill        uint32 // The size in bytes of the template, beyond the initialized data delimited by the Raw Data Start VA and Raw Data End VA fields. The total template size should be the same as the total size of TLS data in the image file. The zero fill is the amount of data that comes after the initialized nonzero data.
	Characteristics       uint32 // The four bits [23:20] describe alignment info. Possible values are those defined as IMAGE_SCN_ALIGN_*, which are also used to describe alignment of section in object files. The other 28 bits are reserved for future use.
}

// ImageTLSDirectory64 represents the IMAGE_TLS_DIRECTORY64 structure.
// It Points to the Thread Local Storage initialization section.
type ImageTLSDirectory64 struct {
	StartAddressOfRawData uint64 // The starting address of the TLS template. The template is a block of data that is used to initialize TLS data.
	EndAddressOfRawData   uint64 // The address of the last byte of the TLS, except for the zero fill. As with the Raw Data Start VA field, this is a VA, not an RVA.
	AddressOfIndex        uint64 // The location to receive the TLS index, which the loader assigns. This location is in an ordinary data section, so it can be given a symbolic name that is accessible to the program.
	AddressOfCallBacks    uint64 // The pointer to an array of TLS callback functions. The array is null-terminated, so if no callback function is supported, this field points to 4 bytes set to zero.
	SizeOfZeroFill        uint32 // The size in bytes of the template, beyond the initialized data delimited by the Raw Data Start VA and Raw Data End VA fields. The total template size should be the same as the total size of TLS data in the image file. The zero fill is the amount of data that comes after the initialized nonzero data.
	Characteristics       uint32 // The four bits [23:20] describe alignment info. Possible values are those defined as IMAGE_SCN_ALIGN_*, which are also used to describe alignment of section in object files. The other 28 bits are reserved for future use.
}

func (pe *File) parseTLSDirectory(rva, size uint32) (TLSDirectory, error) {

	tls := TLSDirectory{}

	if pe.Is64 {
		tlsDir := ImageTLSDirectory64{}
		tlsSize := uint32(binary.Size(ImageTLSDirectory64{}))
		fileOffset := pe.getOffsetFromRva(rva)

		buf := bytes.NewReader(pe.data[fileOffset : fileOffset+tlsSize])
		err := binary.Read(buf, binary.LittleEndian, &tlsDir)
		if err != nil {
			return tls, err
		}

		rvaAddressOfCallBacks := uint32(tlsDir.AddressOfCallBacks - pe.OptionalHeader64.ImageBase)
		offset := pe.getOffsetFromRva(rvaAddressOfCallBacks)
		tls.Struct = tlsDir

		var callbacks []uint64
		for i := 0; ; i++ {
			c := binary.LittleEndian.Uint64(pe.data[offset+(uint32(i)*4):])
			if c == 0 {
				break
			}
			callbacks = append(callbacks, c)
		}

		tls.Callbacks = callbacks
	} else {
		tlsDir := ImageTLSDirectory32{}
		tlsSize := uint32(binary.Size(ImageTLSDirectory32{}))
		fileOffset := pe.getOffsetFromRva(rva)

		buf := bytes.NewReader(pe.data[fileOffset : fileOffset+tlsSize])
		err := binary.Read(buf, binary.LittleEndian, &tlsDir)
		if err != nil {
			return tls, err
		}

		// 94a9dc17d47b03f6fb01cb639e25503b37761b452e7c07ec6b6c2280635f1df9
		// Callbacks may be empty
		var callbacks []uint32
		if tlsDir.AddressOfCallBacks != 0 {
			rvaAddressOfCallBacks := tlsDir.AddressOfCallBacks - pe.OptionalHeader.ImageBase
			offset := pe.getOffsetFromRva(rvaAddressOfCallBacks)
			tls.Struct = tlsDir

			for i := 0; ; i++ {
				c := binary.LittleEndian.Uint32(pe.data[offset+(uint32(i)*4):])
				if c == 0 {
					break
				}
				callbacks = append(callbacks, c)
			}
		}

		tls.Callbacks = callbacks
	}

	return tls, nil
}
