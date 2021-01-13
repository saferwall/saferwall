# Portable Executable Parser

**peparser** is a go package for parsing the [portable executable](https://docs.microsoft.com/en-us/windows/win32/debug/pe-format) file format. This package was designed with malware analysis in mind, and being resistent to PE malformations.

## Features

- Works with PE32/PE32+ file fomat.
- Supports Intel x86/AMD64/ARM7ARM7 Thumb/ARM8-64/IA64/CHPE architectures.
- MS DOS header.
- Rich Header (calculate checksum).
- NT Header (file header + optional header).
- COFF symbol table and string table.
- Sections headers + entropy calculation. 
- Data directories
    - Import Table + ImpHash calculation.
    - Export Table
    - Resource Table
    - Exceptions Table
    - Security Table + Authentihash calculation.
    - Relocations Table
    - Debug Table (CODEVIEW, POGO, VC FEATURE, REPRO, FPO, EXDLL CHARACTERISTICS debug types).
    - TLS Table
    - Load Config Directory (SEH, GFID, GIAT, Guard LongJumps, CHPE, Dynamic Value Reloc Table, Enclave Configuration, Volatile Metadata tables).
    - Bound Import Table
    - Delay Import Table
    - COM Table (CLR Metadata Header, Metadata Table Streams)
- Report several anomalies

## Installing

Using peparser is easy. First, use `go get` to install the latest version
of the library. This command will install the `peparser` generator executable
along with the library and its dependencies:

    go get -u github.com/saferwall/saferwall/pkg/peparser

Next, include `peparser` in your application:

```go
import "github.com/saferwall/saferwall/pkg/peparser"
```

## Using the library

```go
package main

import (
	"github.com/saferwall/saferwall/pkg/peparser"
)

func main() {
    pe, err := peparser.New("C:\\Binaries\\notepad.exe", nil)
	if err != nil {
		log.Fatalf("Error while opening file: %s, reason: %s", filename, err)
    }
    
    err = pe.Parse()
    if err != nil {
        log.Fatalf("Error while opening file: %s, reason: %s", filename, err)
    }
```

## Todo:

- imports MS-styled names demangling
- PE: VB5 and VB6 typical structures: project info, DLLCall-imports, referenced modules, object table

# References

- [Peering Inside the PE: A Tour of the Win32 Portable Executable File Format by Matt Pietrek](http://bytepointer.com/resources/pietrek_peering_inside_pe.htm)
- [An In-Depth Look into the Win32 Portable Executable File Format - Part 1 by Matt Pietrek](http://www.delphibasics.info/home/delphibasicsarticles/anin-depthlookintothewin32portableexecutablefileformat-part1)
- [An In-Depth Look into the Win32 Portable Executable File Format - Part 2 by Matt Pietrek](http://www.delphibasics.info/home/delphibasicsarticles/anin-depthlookintothewin32portableexecutablefileformat-part2)
- [Portable Executable File Format](https://blog.kowalczyk.info/articles/pefileformat.html)
- [PE Format MSDN spec](https://docs.microsoft.com/en-us/windows/win32/debug/pe-format)