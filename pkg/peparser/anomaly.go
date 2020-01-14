package pe

import (
	"encoding/binary"
)

// Anomalies found in a PE
var (
	// NumberOfSections is reported when number of sections is larger or equal than 10.
	AnoNumberOfSections10Plus = "Number of sections is 10+"

	// AnoNumberOfSectionsNull is reported when sections count's is 0.
	AnoNumberOfSectionsNull = "Number of sections is 0"

	// AnoSizeOfOptionalHeaderNull is reported when size of optional header is 0.
	AnoSizeOfOptionalHeaderNull = "Size of optional header is 0"

	// AnoUncommonSizeOfOptionalHeader32 is reported when size of optional header for PE32 is larger than 0xE0.
	AnoUncommonSizeOfOptionalHeader32 = "Size of optional header is larger than 0xE0 (PE32)"

	// AnoUncommonSizeOfOptionalHeader64 is reported when size of optional header for PE32+ is larger than 0xF0.
	AnoUncommonSizeOfOptionalHeader64 = "Size of optional header is larger than 0xF0 (PE32+)"

	// AnoAddressOfEntryPointNull is reported when address of entry point is 0.
	AnoAddressOfEntryPointNull = "Address of entry point is 0."

	// AnoAddressOfEPLessSizeOfHeaders is reported when address of entry point is smaller than size of headers, the file cannot run under Windows.
	AnoAddressOfEPLessSizeOfHeaders = "Address of entry point is smaller than size of headers, the file cannot run under Windows 8"

	// AnoImageBaseNull is reported when the image base is null
	AnoImageBaseNull = "Image base is 0"

	// AnoDanSMagicOffset is reported when the `DanS` magic offset is different than 0x80.
	AnoDanSMagicOffset = "`DanS` magic offset is different than 0x80."

	// ErrInvalidFileAlignment is reported when file alignment is larger than 0x200 and not a power of 2.
	ErrInvalidFileAlignment = "FileAlignment larger than 0x200 and not a power of 2"

	// ErrInvalidSectionAlignment is reported when file alignment is lesser than 0x200 and different from section alignment.
	ErrInvalidSectionAlignment = "FileAlignment lesser than 0x200 and different from section alignment"
)

// GetAnomalies reportes anomalies found in a PE binary.
// Code was explicitely put in a separate function to not process it if user is not interested.
func (pe *File) GetAnomalies() error {

	// ******************************* Anomalies in File header *******************************
	// An application for Windows NT typically has the nine predefined sections named:
	// .text, .bss, .rdata, .data, .rsrc, .edata, .idata, .pdata, and .debug.
	// Some applications do not need all of these sections,
	//  while others may define still more sections to suit their specific needs.
	// NumberOfSections can be up to 96 under XP.
	// NumberOfSections can be up to 65535 under Vista and later.
	if pe.FileHeader.NumberOfSections >= 10 {
		pe.Anomalies = append(pe.Anomalies, AnoNumberOfSections10Plus)
	}

	// NumberOfSections can be null with low alignment PEs
	// and in this case, the values are just checked but not really used (under XP)
	if pe.FileHeader.NumberOfSections == 0 {
		pe.Anomalies = append(pe.Anomalies, AnoNumberOfSectionsNull)
	}

	// SizeOfOptionalHeader is not the size of the optional header, but the delta between
	// the top of the Optional header and the start of the section table.
	// Thus, it can be null (the section table will overlap the Optional Header,
	// or can be null when no sections are present)
	if pe.FileHeader.SizeOfOptionalHeader == 0 {
		pe.Anomalies = append(pe.Anomalies, AnoSizeOfOptionalHeaderNull)
	}

	// SizeOfOptionalHeader can be bigger than the file
	// (the section table will be in virtual space, full of zeroes), but can't be negative.
	// Do some check here.

	// SizeOfOptionalHeader standard value is 0xE0 for PE32.
	if !pe.Is64 && pe.FileHeader.SizeOfOptionalHeader > uint16(binary.Size(pe.OptionalHeader)) {
		pe.Anomalies = append(pe.Anomalies, AnoUncommonSizeOfOptionalHeader32)
	}

	// SizeOfOptionalHeader standard value is 0xF0 for PE32+.
	if pe.Is64 && pe.FileHeader.SizeOfOptionalHeader > uint16(binary.Size(pe.OptionalHeader64)) {
		pe.Anomalies = append(pe.Anomalies, AnoUncommonSizeOfOptionalHeader64)
	}

	// **************************** Anomalies in Optional header ****************************
	// Under Windows 8, AddressOfEntryPoint is not allowed to be smaller than SizeOfHeaders,
	// except if it's null.
	if pe.OptionalHeader.AddressOfEntryPoint != 0 &&
		pe.OptionalHeader.AddressOfEntryPoint < pe.OptionalHeader.SizeOfHeaders {
		pe.Anomalies = append(pe.Anomalies, AnoAddressOfEPLessSizeOfHeaders)
	}

	// AddressOfEntryPoint can be null in DLLs: in this case, DllMain is just not called.
	// can be null
	if pe.OptionalHeader.AddressOfEntryPoint == 0 {
		pe.Anomalies = append(pe.Anomalies, AnoAddressOfEntryPointNull)
	}

	// ImageBase can be null, under XP. In this case, the binary will be relocated to 10000h
	if (pe.Is64 && pe.OptionalHeader64.ImageBase == 0) || (!pe.Is64 && pe.OptionalHeader.ImageBase == 0) {
		pe.Anomalies = append(pe.Anomalies, AnoImageBaseNull)
	}

	// MajorSubsystemVersion/MinorSubsystemVersion
	// Todo :MajorSubsystemVersion.MinorSubsystemVersion has to be at least 3.10

	// (Win32VersionValue)
	// officially defined as ''reserved'' and should be null
	// if non null, it overrides MajorVersion/MinorVersion/BuildNumber/PlatformId OperatingSystem Versions values located in the PEB, after loading.

	// SizeOfImage
	// normally equal the total virtual size of all sections + headers
	// 65535sects.exe has a SizeOfImage of 0x7027A000 (W7 only)
	// tinyXP has a SizeOfImage of 0x2e. it only covers up to the EntryPoint value.

	// todo
	// SizeOfImage
	// normally equal the total virtual size of all sections + headers
	// 65535sects.exe has a SizeOfImage of 0x7027A000 (W7 only)
	// tinyXP has a SizeOfImage of 0x2e. it only covers up to the EntryPoint value.

	// todo:
	// CheckSum
	// simple algorithm
	// required for drivers only


	return nil

}
