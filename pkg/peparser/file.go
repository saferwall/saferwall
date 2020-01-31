// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"io"
	"os"

	mmap "github.com/edsrzf/mmap-go"
)

// A File represents an open PE file.
type File struct {
	DosHeader        ImageDosHeader
	NtHeader         ImageNtHeader
	Sections         []ImageSectionHeader
	Imports          []Import
	Exports          []ExportFunction
	Debugs           []DebugEntry
	Relocations      []Relocation
	Resources        ResourceDirectory
	TLS              TLSDirectory
	LoadConfig       interface{}
	Exceptions       []Exception
	Certificates     Certificate
	DelayImports     []DelayImport
	BoundImports     []BoundImportDescriptorData
	GlobalPtr        uint32
	RichHeader       RichHeader
	CLRHeader        ImageCOR20Header

	Header    []byte
	data      mmap.MMap
	closer    io.Closer
	Is64      bool
	Is32      bool
	Anomalies []string
	size      uint32
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
	file.size = uint32(len(file.data))
	return file, nil
}

// Parse performs the file parsing for a PE binary.
func (pe *File) Parse() error {

	// check for the smallest PE size.
	if len(pe.data) < TinyPESize {
		return ErrInvalidPESize
	}

	// Parse the DOS header.
	err := pe.parseDosHeader()
	if err != nil {
		return err
	}

	// Parse the NT header.
	err = pe.parseNtHeader()
	if err != nil {
		return err
	}


	// Parse the Section Header.
	err = pe.parseSectionHeader()
	if err != nil {
		return err
	}

	// Parse the Data Directory entries.
	err = pe.parseDataDirectories()
	if err != nil {
		return err
	}

	return nil
}

func (pe *File) parseDataDirectories() (err error) {
	oh32 := ImageOptionalHeader32{}
	oh64 := ImageOptionalHeader64{}
	switch pe.Is64 {
	case true:
		oh64 = pe.NtHeader.OptionalHeader.(ImageOptionalHeader64)
	case false:
		oh32 = pe.NtHeader.OptionalHeader.(ImageOptionalHeader32)		
	}

	if pe.Is64 {
		importDirEntry := oh64.DataDirectory[ImageDirectoryEntryImport]
		if importDirEntry.VirtualAddress != 0 {
			err = pe.parseImportDirectory(importDirEntry.VirtualAddress, importDirEntry.Size)
		}

		exportDirEntry := oh64.DataDirectory[ImageDirectoryEntryExport]
		if exportDirEntry.VirtualAddress != 0 {
			err = pe.parseExportDirectory(exportDirEntry.VirtualAddress, exportDirEntry.Size)
		}

		debugDirEntry := oh64.DataDirectory[ImageDirectoryEntryDebug]
		if debugDirEntry.VirtualAddress != 0 {
			err = pe.parseDebugDirectory(debugDirEntry.VirtualAddress, debugDirEntry.Size)
		}

		relocDirEntry := oh64.DataDirectory[ImageDirectoryEntryBaseReloc]
		if relocDirEntry.VirtualAddress != 0 {
			err = pe.parseRelocDirectory(relocDirEntry.VirtualAddress, relocDirEntry.Size)
		}

		rsrcDirEntry := oh64.DataDirectory[ImageDirectoryEntryResource]
		if relocDirEntry.VirtualAddress != 0 {
			pe.Resources, err = pe.parseResourceDirectory(rsrcDirEntry.VirtualAddress, rsrcDirEntry.Size, 0, 0)
		}

		tlsDirEntry := oh64.DataDirectory[ImageDirectoryEntryTLS]
		if tlsDirEntry.VirtualAddress != 0 {
			pe.TLS, err = pe.parseTLSDirectory(tlsDirEntry.VirtualAddress, tlsDirEntry.Size)
		}

		loadConfigDirEntry := oh64.DataDirectory[ImageDirectoryEntryLoadConfig]
		if tlsDirEntry.VirtualAddress != 0 {
			pe.LoadConfig, err = pe.parseLoadConfigDirectory(loadConfigDirEntry.VirtualAddress, loadConfigDirEntry.Size)
		}

		exceptionDirEntry := oh64.DataDirectory[ImageDirectoryEntryException]
		if exceptionDirEntry.VirtualAddress != 0 {
			pe.Exceptions, err = pe.parseExceptionDirectory(exceptionDirEntry.VirtualAddress, exceptionDirEntry.Size)
		}

		certificateDirEntry := oh64.DataDirectory[ImageDirectoryEntryCertificate]
		if certificateDirEntry.VirtualAddress != 0 {
			pe.Certificates, err = pe.parseSecurityDirectory(certificateDirEntry.VirtualAddress, certificateDirEntry.Size)
		}

		delayImportDirEntry := oh64.DataDirectory[ImageDirectoryEntryDelayImport]
		if delayImportDirEntry.VirtualAddress != 0 {
			err = pe.parseDelayImportDirectory(delayImportDirEntry.VirtualAddress, delayImportDirEntry.Size)
		}

		boundImportDirEntry := oh64.DataDirectory[ImageDirectoryEntryBoundImport]
		if boundImportDirEntry.VirtualAddress != 0 {
			err = pe.parseBoundImportDirectory(boundImportDirEntry.VirtualAddress, boundImportDirEntry.Size)
		}

		globalPtrDirEntry := oh64.DataDirectory[ImageDirectoryEntryGlobalPtr]
		if globalPtrDirEntry.VirtualAddress != 0 {
			err = pe.parseGlobalPtrDirectory(globalPtrDirEntry.VirtualAddress, globalPtrDirEntry.Size)
		}

		iatDirEntry := oh64.DataDirectory[ImageDirectoryEntryIAT]
		if iatDirEntry.VirtualAddress != 0 {
			err = pe.parseIATDirectory(iatDirEntry.VirtualAddress, iatDirEntry.Size)
		}

		clrHeaderDirEntry := oh64.DataDirectory[ImageDirectoryEntryCLR]
		if clrHeaderDirEntry.VirtualAddress != 0 {
			err = pe.parseCLRHeaderDirectory(clrHeaderDirEntry.VirtualAddress, clrHeaderDirEntry.Size)
		}
	}

	if pe.Is32 {
		importDirEntry := oh32.DataDirectory[ImageDirectoryEntryImport]
		if importDirEntry.VirtualAddress != 0 {
			err = pe.parseImportDirectory(importDirEntry.VirtualAddress, importDirEntry.Size)
		}

		exportDirEntry := oh32.DataDirectory[ImageDirectoryEntryExport]
		if exportDirEntry.VirtualAddress != 0 {
			err = pe.parseExportDirectory(exportDirEntry.VirtualAddress, exportDirEntry.Size)
		}

		debugDirEntry := oh32.DataDirectory[ImageDirectoryEntryDebug]
		if debugDirEntry.VirtualAddress != 0 {
			err = pe.parseDebugDirectory(debugDirEntry.VirtualAddress, debugDirEntry.Size)
		}

		relocDirEntry := oh32.DataDirectory[ImageDirectoryEntryBaseReloc]
		if relocDirEntry.VirtualAddress != 0 {
			err = pe.parseRelocDirectory(relocDirEntry.VirtualAddress, relocDirEntry.Size)
		}

		rsrcDirEntry := oh32.DataDirectory[ImageDirectoryEntryResource]
		if rsrcDirEntry.VirtualAddress != 0 {
			pe.Resources, err = pe.parseResourceDirectory(rsrcDirEntry.VirtualAddress, rsrcDirEntry.Size, 0, 0)
		}

		tlsDirEntry := oh32.DataDirectory[ImageDirectoryEntryTLS]
		if tlsDirEntry.VirtualAddress != 0 {
			pe.TLS, err = pe.parseTLSDirectory(tlsDirEntry.VirtualAddress, tlsDirEntry.Size)
		}

		loadConfigDirEntry := oh32.DataDirectory[ImageDirectoryEntryLoadConfig]
		if tlsDirEntry.VirtualAddress != 0 {
			pe.LoadConfig, err = pe.parseLoadConfigDirectory(loadConfigDirEntry.VirtualAddress, loadConfigDirEntry.Size)
		}

		certificateDirEntry := oh32.DataDirectory[ImageDirectoryEntryCertificate]
		if certificateDirEntry.VirtualAddress != 0 {
			pe.Certificates, err = pe.parseSecurityDirectory(certificateDirEntry.VirtualAddress, certificateDirEntry.Size)
		}

		delayImportDirEntry := oh32.DataDirectory[ImageDirectoryEntryDelayImport]
		if delayImportDirEntry.VirtualAddress != 0 {
			err = pe.parseDelayImportDirectory(delayImportDirEntry.VirtualAddress, delayImportDirEntry.Size)
		}

		boundImportDirEntry := oh32.DataDirectory[ImageDirectoryEntryBoundImport]
		if boundImportDirEntry.VirtualAddress != 0 {
			err = pe.parseBoundImportDirectory(boundImportDirEntry.VirtualAddress, boundImportDirEntry.Size)
		}

		globalPtrDirEntry := oh32.DataDirectory[ImageDirectoryEntryGlobalPtr]
		if globalPtrDirEntry.VirtualAddress != 0 {
			err = pe.parseGlobalPtrDirectory(globalPtrDirEntry.VirtualAddress, globalPtrDirEntry.Size)
		}

		iatDirEntry := oh32.DataDirectory[ImageDirectoryEntryIAT]
		if iatDirEntry.VirtualAddress != 0 {
			err = pe.parseIATDirectory(iatDirEntry.VirtualAddress, iatDirEntry.Size)
		}

		clrHeaderDirEntry := oh32.DataDirectory[ImageDirectoryEntryCLR]
		if clrHeaderDirEntry.VirtualAddress != 0 {
			err = pe.parseCLRHeaderDirectory(clrHeaderDirEntry.VirtualAddress, clrHeaderDirEntry.Size)
		}
	}

	return err
}
