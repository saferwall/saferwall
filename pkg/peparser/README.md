# Portable Executable Parser

**peparser** is a go package for parsing the [portable executable](https://docs.microsoft.com/en-us/windows/win32/debug/pe-format) file format. This package was designed with malware analysis in mind, and being resistent to PE malformations.

## Features

- :heavy_check_mark: Works with PE32/PE32+ file fomat.
- :heavy_check_mark: Supports Intel x86/AMD64/ARM7ARM7 Thumb/ARM8-64/IA64/CHPE architectures.
- :heavy_check_mark: MS DOS header
- :heavy_check_mark: Rich Header (calculate checksum)
- :heavy_check_mark: NT Header (file header + optional header)
- :heavy_check_mark: Sections headers
- Data directories
    - :heavy_check_mark: Import Table + ImpHash calculation.
    - :heavy_check_mark: Export Table
    - :heavy_check_mark: Resource Table
    - :heavy_check_mark: Exceptions Table
    - :heavy_check_mark: Security Table + Authentihash calculation.
    - :heavy_check_mark: Relocations Table
    - :heavy_check_mark: Debug Table (CODEVIEW, POGO, VC FEATURE, REPRO, FPO, EXDLL CHARACTERISTICS debug types).
    - :heavy_check_mark: TLS Table
    - :heavy_check_mark: Load Config Directory (SEH, GFID, GIAT, Guard LongJumps, CHPE, Dynamic Value Reloc Table, Enclave Configuration, Volatile Metadata tables).
    - :heavy_check_mark: Bound Import Table
    - :heavy_check_mark: Delay Import Table
    - :heavy_check_mark: COM Table
- :heavy_check_mark: Report several anomalies

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
    pe, err := peparser.Open("C:\\Binaries\\notepad.exe")
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
- section entropy

# References

- [Peering Inside the PE: A Tour of the Win32 Portable Executable File Format by Matt Pietrek](http://bytepointer.com/resources/pietrek_peering_inside_pe.htm)
- [An In-Depth Look into the Win32 Portable Executable File Format - Part 1 by Matt Pietrek](http://www.delphibasics.info/home/delphibasicsarticles/anin-depthlookintothewin32portableexecutablefileformat-part1)
- [An In-Depth Look into the Win32 Portable Executable File Format - Part 2 by Matt Pietrek](http://www.delphibasics.info/home/delphibasicsarticles/anin-depthlookintothewin32portableexecutablefileformat-part2)
- [Portable Executable File Format](https://blog.kowalczyk.info/articles/pefileformat.html)
- [PE Format MSDN spec](https://docs.microsoft.com/en-us/windows/win32/debug/pe-format)