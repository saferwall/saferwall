// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

// Image executable types
const (

	// The DOS MZ executable format is the executable file format used
	// for .EXE files in DOS.
	ImageDOSSignature   = 0x5A4D     // MZ
	ImageDOSZMSignature = 0x4D5A     // ZM

	// The New Executable (abbreviated NE or NewEXE) is a 16-bit .exe file 
	// format, a successor to the DOS MZ executable format. It was used in
	// Windows 1.0–3.x, multitasking MS-DOS 4.0, OS/2 1.x, and the OS/2 subset
	// of Windows NT up to version 5.0 (Windows 2000). A NE is also called a
	// segmented executable.
	ImageOS2Signature   = 0x454E
	
	// Linear Executable is an executable file format in the EXE family.
	// It was used by 32-bit OS/2, by some DOS extenders, and by Microsoft
	// Windows VxD files. It is an extension of MS-DOS EXE, and a successor
	// to NE (New Executable).
	ImageOS2LESignature = 0x454C
	
	// There are two main varieties of LE executables:
	// LX (32-bit), and LE (mixed 16/32-bit).
	ImageVXDSignature    = 0x584C
	
	// Terse Executables have a 'VZ' signature.
	ImageTESignature    = 0x5A56
	
	// The Portable Executable (PE) format is a file format for executables,
	// object code, DLLs and others used in 32-bit and 64-bit versions of
	// Windows operating systems.
	ImageNTSignature    = 0x00004550 // PE00
)

// Optional Header magic
const (
	ImageNtOptionalHeader32Magic = 0x10b
	ImageNtOptionalHeader64Magic = 0x20b
	ImageROMOptionalHeaderMagic  = 0x10
)

// Image file machine types
const (
	ImageFileMachineUnknown   = uint16(0x0)    // The contents of this field are assumed to be applicable to any machine type
	ImageFileMachineAM33      = uint16(0x1d3)  // Matsushita AM33
	ImageFileMachineAMD64     = uint16(0x8664) // x64
	ImageFileMachineARM       = uint16(0x1c0)  // ARM little endian
	ImageFileMachineARM64     = uint16(0xaa64) // ARM64 little endian
	ImageFileMachineARMNT     = uint16(0x1c4)  // ARM Thumb-2 little endian
	ImageFileMachineEBC       = uint16(0xebc)  // EFI byte code
	ImageFileMachineI386      = uint16(0x14c)  // Intel 386 or later processors and compatible processors
	ImageFileMachineIA64      = uint16(0x200)  // Intel Itanium processor family
	ImageFileMachineM32R      = uint16(0x9041) // Mitsubishi M32R little endian
	ImageFileMachineMIPS16    = uint16(0x266)  // MIPS16
	ImageFileMachineMIPSFPU   = uint16(0x366)  // MIPS with FPU
	ImageFileMachineMIPSFPU16 = uint16(0x466)  // MIPS16 with FPU
	ImageFileMachinePowerPC   = uint16(0x1f0)  // Power PC little endian
	ImageFileMachinePowerPCFP = uint16(0x1f1)  // Power PC with floating point support
	ImageFileMachineR4000     = uint16(0x166)  // MIPS little endian
	ImageFileMachineRISCV32   = uint16(0x5032) // RISC-V 32-bit address space
	ImageFileMachineRISCV64   = uint16(0x5064) // RISC-V 64-bit address space
	ImageFileMachineRISCV128  = uint16(0x5128) // RISC-V 128-bit address space
	ImageFileMachineSH3       = uint16(0x1a2)  // Hitachi SH3
	ImageFileMachineSH3DSP    = uint16(0x1a3)  // Hitachi SH3 DSP
	ImageFileMachineSH4       = uint16(0x1a6)  // Hitachi SH4
	ImageFileMachineSH5       = uint16(0x1a8)  // Hitachi SH5
	ImageFileMachineTHUMB     = uint16(0x1c2)  // Thumb
	ImageFileMachineWCEMIPSv2 = uint16(0x169)  // MIPS little-endian WCE v2
)

// The Characteristics field contains flags that indicate attributes of the object or image file.
const (
	ImageFileRelocsStripped       = 0x0001 // Relocation info stripped from file.
	ImageFileExecutableImage      = 0x0002 // File is executable  (i.e. no unresolved external references).
	ImageFileLineNumsStripped     = 0x0004 // Line numbers stripped from file.
	ImageFileLocalSymsStripped    = 0x0008 // Local symbols stripped from file.
	ImageFileAgressibeWsTrim      = 0x0010 // Aggressively trim working set
	ImageFileLargeAddressAware    = 0x0020 // App can handle >2gb addresses
	ImageFileBytesReservedLow     = 0x0080 // Bytes of machine word are reversed.
	ImageFile32BitMachine         = 0x0100 // 32 bit word machine.
	ImageFileDebugStripped        = 0x0200 // Debugging info stripped from file in .DBG file
	ImageFileRemovableRunFromSwap = 0x0400 // If Image is on removable media, copy and run from the swap file.
	ImageFileNetRunFromSwap       = 0x0800 // If Image is on Net, copy and run from the swap file.
	ImageFileSystem               = 0x1000 // System File.
	ImageFileDLL                  = 0x2000 // File is a DLL.
	ImageFileUpSystemOnly         = 0x4000 // File should only be run on a UP machine
	ImageFileBytesReservedHigh    = 0x8000 // Bytes of machine word are reversed.
)

// Subsystem values of an OptionalHeader
const (
	ImageSubsystemUnknown                = 0  // An unknown subsystem.
	ImageSubsystemNative                 = 1  // Device drivers and native Windows processes
	ImageSubsystemWindowsGUI             = 2  // The Windows graphical user interface (GUI) subsystem.
	ImageSubsystemWindowsCUI             = 3  // The Windows character subsystem
	ImageSubsystemOS2CUI                 = 5  // The OS/2 character subsystem.
	ImageSubsystemPosixCUI               = 7  // The Posix character subsystem.
	ImageSubsystemNativeWindows          = 8  // Native Win9x driver
	ImageSubsystemWindowsCEGUI           = 9  // Windows CE
	ImageSubsystemEFIApplication         = 10 // An Extensible Firmware Interface (EFI) application
	ImageSubsystemEFIBootServiceDriver   = 11 // An EFI driver with boot services
	ImageSubsystemEFIRuntimeDriver       = 12 // An EFI driver with run-time services
	ImageSubsystemEFIRom                 = 13 // An EFI ROM image .
	ImageSubsystemXBOX                   = 14 // XBOX.
	ImageSubsystemWindowsBootApplication = 16 // Windows boot application.
)

// DllCharacteristics values of an OptionalHeader
const (
	ImageDllCharacteristicsReserved1            = 0x0001 // Reserved, must be zero.
	ImageDllCharacteristicsReserved2            = 0x0002 // Reserved, must be zero.
	ImageDllCharacteristicsReserved4            = 0x0004 // Reserved, must be zero.
	ImageDllCharacteristicsReserved8            = 0x0008 // Reserved, must be zero.
	ImageDllCharacteristicsHighEntropyVA        = 0x0020 // Image can handle a high entropy 64-bit virtual address space
	ImageDllCharacteristicsDynamicBase          = 0x0040 // DLL can be relocated at load time.
	ImageDllCharacteristicsForceIntegrity       = 0x0080 // Code Integrity checks are enforced.
	ImageDllCharacteristicsNXCompact            = 0x0100 // Image is NX compatible.
	ImageDllCharacteristicsNoIsolation          = 0x0200 // Isolation aware, but do not isolate the image.
	ImageDllCharacteristicsNoSEH                = 0x0400 // Does not use structured exception (SE) handling. No SE handler may be called in this image.
	ImageDllCharacteristicsNoBind               = 0x0800 // Do not bind the image.
	ImageDllCharacteristicsAppContainer         = 0x1000 // Image must execute in an AppContainer
	ImageDllCharacteristicsWdmDriver            = 0x2000 // A WDM driver.
	ImageDllCharacteristicsGuardCF              = 0x4000 // Image supports Control Flow Guard.
	ImageDllCharacteristicsTerminalServiceAware = 0x8000 // Terminal Server aware.

)

// DataDirectory entries of an OptionalHeader
const (
	ImageDirectoryEntryExport       = 0  // Export Table
	ImageDirectoryEntryImport       = 1  // Import Table
	ImageDirectoryEntryResource     = 2  // Resource Table
	ImageDirectoryEntryException    = 3  // Exception Table
	ImageDirectoryEntryCertificate  = 4  // Certificate Directory
	ImageDirectoryEntryBaseReloc    = 5  // Base Relocation Table
	ImageDirectoryEntryDebug        = 6  // Debug
	ImageDirectoryEntryArchitecture = 7  // Architecture Specific Data
	ImageDirectoryEntryGlobalPtr    = 8  // The RVA of the value to be stored in the global pointer register.
	ImageDirectoryEntryTLS          = 9  // The thread local storage (TLS) table
	ImageDirectoryEntryLoadConfig   = 10 // The load configuration table
	ImageDirectoryEntryBoundImport  = 11 // The bound import table
	ImageDirectoryEntryIAT          = 12 // Import Address Table
	ImageDirectoryEntryDelayImport  = 13 // Delay Import Descriptor
	ImageDirectoryEntryCLR          = 14 // CLR Runtime Header
	ImageDirectoryEntryReserved     = 15 // Must be zero
	ImageNumberOfDirectoryEntries   = 16 // Tables count.
)
