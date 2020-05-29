# Portable Executable Parser

**peparser** is a go package for parsing the [portable executable](https://docs.microsoft.com/en-us/windows/win32/debug/pe-format) file format. This package was designed with malware analysis in mind, and being resistent to PE malformations.

## Features

- Works with PE32(x86) and PE32+(x64) binaries.
- MS DOS header
- Rich Header (verify checksum)
- NT Header (file header + optional header)
- Sections headers
- Data directories
    - Import Table
    - Export Table
    - Resource Table
    - Exceptions Table
    - Security Table
    - Relocations Table
    - Debug Table
    - TLS Table
    - Load Config Directory
    - Bound Import Table
    - Delay Import Table
    - COM Table
- Calculate Authentihash
- Calculate ImpHash
- Report several anomalies

# References

- [Peering Inside the PE: A Tour of the Win32 Portable Executable File Format by Matt Pietrek](http://bytepointer.com/resources/pietrek_peering_inside_pe.htm)
- [An In-Depth Look into the Win32 Portable Executable File Format - Part 1 by Matt Pietrek](http://www.delphibasics.info/home/delphibasicsarticles/anin-depthlookintothewin32portableexecutablefileformat-part1)
- [An In-Depth Look into the Win32 Portable Executable File Format - Part 2 by Matt Pietrek](http://www.delphibasics.info/home/delphibasicsarticles/anin-depthlookintothewin32portableexecutablefileformat-part2)
- [Portable Executable File Format](https://blog.kowalczyk.info/articles/pefileformat.html)
- [PE Format MSDN spec](https://docs.microsoft.com/en-us/windows/win32/debug/pe-format)