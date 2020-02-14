// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"bytes"
	"encoding/binary"
)

const (
	// ImageGuardCfInstrumented - Module performs control flow integrity checks
	// using system-supplied support.
	ImageGuardCfInstrumented = 0x00000100

	// ImageGuardCfWInstrumented - Module performs control flow and write
	// integrity checks.
	ImageGuardCfWInstrumented = 0x00000200

	// ImageGuardCfFunctionTablePresent - Module contains valid control flow
	// target metadata.
	ImageGuardCfFunctionTablePresent = 0x00000400

	// ImageGuardSecurityCookieUnused - Module does not make use of the /GS
	// security cookie.
	ImageGuardSecurityCookieUnused = 0x00000800

	// ImageGuardProtectDelayloadIAT - Module supports read only delay load IAT.
	ImageGuardProtectDelayloadIAT = 0x00001000

	// ImageGuardDelayloadIATInItsOwnSection - Delayload import table in its own
	// .didat section (with nothing else in it) that can be freely reprotected.
	ImageGuardDelayloadIATInItsOwnSection = 0x00002000

	// ImageGuardCfExportSuppressionInfoPresent - Module contains suppressed
	// export information. This also infers that the address taken IAT table is
	// also present in the load config.
	ImageGuardCfExportSuppressionInfoPresent = 0x00004000

	// ImageGuardCfEnableExportSuppression - Module enables suppression of exports.
	ImageGuardCfEnableExportSuppression = 0x00008000

	// ImageGuardCfLongjumpTablePresent - Module contains longjmp target information.
	ImageGuardCfLongjumpTablePresent = 0x00010000

	// ImageGuardCfFnctionTableSizeMask - Mask for the subfield that contains
	// the stride of Control Flow Guard function table entries (that is, the
	// additional count of bytes per table entry).
	ImageGuardCfFnctionTableSizeMask = 0xF0000000
)

// ImageLoadConfigCodeIntegrity Code Integrity in loadconfig (CI)
type ImageLoadConfigCodeIntegrity struct {
	Flags         uint16 // Flags to indicate if CI information is available, etc.
	Catalog       uint16 // 0xFFFF means not available
	CatalogOffset uint32
	Reserved      uint32 // Additional bitmask to be defined later
}

// ImageLoadConfigDirectory32 represents the Load Configuration Structure IMAGE_LOAD_CONFIG_DIRECTORY.
type ImageLoadConfigDirectory32 struct {
	Size uint32

	// Date and time stamp value. The value is represented in the number of
	// seconds that have elapsed since midnight (00:00:00), January 1, 1970,
	// Universal Coordinated Time, according to the system clock. The time stamp
	// can be printed by using the C runtime (CRT) time function.
	TimeDateStamp uint32

	// Major version number.
	MajorVersion uint32

	// Minor version number.
	MinorVersion uint32

	// The global loader flags to clear for this process as the loader starts the process.
	GlobalFlagsClear uint32

	// The global loader flags to set for this process as the loader starts the process.
	GlobalFlagsSet uint32

	// The default timeout value to use for this process's critical sections that are abandoned.
	CriticalSectionDefaultTimeout uint32

	// Memory that must be freed before it is returned to the system, in bytes.
	DeCommitFreeBlockThreshold uint32

	// Total amount of free memory, in bytes.
	DeCommitTotalFreeThreshold uint32

	// [x86 only] The VA of a list of addresses where the LOCK prefix is used so
	// that they can be replaced with NOP on single processor machines.
	LockPrefixTable uint32

	// Maximum allocation size, in bytes.
	MaximumAllocationSize uint32

	// Maximum virtual memory size, in bytes.
	VirtualMemoryThreshold uint32

	// Process heap flags that correspond to the first argument of the HeapCreate
	// function. These flags apply to the process heap that is created during
	// process startup.
	ProcessHeapFlags uint32

	// Setting this field to a non-zero value is equivalent to calling
	// SetProcessAffinityMask with this value during process startup (.exe only)
	ProcessAffinityMask uint32

	// The service pack version identifier.
	CSDVersion uint16

	// Must be zero.
	DependentLoadFlags uint16

	// Reserved for use by the system.
	EditList uint32

	// A pointer to a cookie that is used by Visual C++ or GS implementation.
	SecurityCookie uint32

	// [x86 only] The VA of the sorted table of RVAs of each valid, unique SE
	// handler in the image.
	SEHandlerTable uint32

	// [x86 only] The count of unique handlers in the table.
	SEHandlerCount uint32

	// The VA where Control Flow Guard check-function pointer is stored.
	GuardCFCheckFunctionPointer uint32

	// The VA where Control Flow Guard dispatch-function pointer is stored.
	GuardCFDispatchFunctionPointer uint32

	// The VA of the sorted table of RVAs of each Control Flow Guard function in
	// the image.
	GuardCFFunctionTable uint32

	// The count of unique RVAs in the above table.
	GuardCFFunctionCount uint32

	// Control Flow Guard related flags.
	GuardFlags uint32

	// Code integrity information.
	CodeIntegrity ImageLoadConfigCodeIntegrity

	// The VA where Control Flow Guard address taken IAT table is stored.
	GuardAddressTakenIatEntryTable uint32

	// The count of unique RVAs in the above table.
	GuardAddressTakenIatEntryCount uint32

	// The VA where Control Flow Guard long jump target table is stored.
	GuardLongJumpTargetTable uint32

	// The count of unique RVAs in the above table.
	GuardLongJumpTargetCount uint32

	DynamicValueRelocTable                   uint32
	CHPEMetadataPointer                      uint32
	GuardRFFailureRoutine                    uint32
	GuardRFFailureRoutineFunctionPointer     uint32
	DynamicValueRelocTableOffset             uint32
	DynamicValueRelocTableSection            uint16
	Reserved2                                uint16
	GuardRFVerifyStackPointerFunctionPointer uint32
	HotPatchTableOffset                      uint32
	Reserved3                                uint32
	EnclaveConfigurationPointer              uint32
}

// ImageLoadConfigDirectory64 represents the Load Configuration Structure IMAGE_LOAD_CONFIG_DIRECTORY.
type ImageLoadConfigDirectory64 struct {
	Size                                     uint32
	TimeDateStamp                            uint32                       // Date and time stamp value. The value is represented in the number of seconds that have elapsed since midnight (00:00:00), January 1, 1970, Universal Coordinated Time, according to the system clock. The time stamp can be printed by using the C runtime (CRT) time function.
	MajorVersion                             uint16                       // Major version number.
	MinorVersion                             uint16                       // Minor version number.
	GlobalFlagsClear                         uint32                       // The global loader flags to clear for this process as the loader starts the process.
	GlobalFlagsSet                           uint32                       // The global loader flags to set for this process as the loader starts the process.
	CriticalSectionDefaultTimeout            uint32                       // The default timeout value to use for this process's critical sections that are abandoned.
	DeCommitFreeBlockThreshold               uint64                       // Memory that must be freed before it is returned to the system, in bytes.
	DeCommitTotalFreeThreshold               uint64                       // Total amount of free memory, in bytes.
	LockPrefixTable                          uint64                       // [x86 only] The VA of a list of addresses where the LOCK prefix is used so that they can be replaced with NOP on single processor machines.
	MaximumAllocationSize                    uint64                       // Maximum allocation size, in bytes.
	VirtualMemoryThreshold                   uint64                       // Maximum virtual memory size, in bytes.
	ProcessAffinityMask                      uint64                       // Setting this field to a non-zero value is equivalent to calling SetProcessAffinityMask with this value during process startup (.exe only)
	ProcessHeapFlags                         uint32                       // Process heap flags that correspond to the first argument of the HeapCreate function. These flags apply to the process heap that is created during process startup.
	CSDVersion                               uint16                       // The service pack version identifier.
	DependentLoadFlags                       uint16                       // Must be zero.
	EditList                                 uint64                       // Reserved for use by the system.
	SecurityCookie                           uint64                       // A pointer to a cookie that is used by Visual C++ or GS implementation.
	SEHandlerTable                           uint64                       // [x86 only] The VA of the sorted table of RVAs of each valid, unique SE handler in the image.
	SEHandlerCount                           uint64                       // [x86 only] The count of unique handlers in the table.
	GuardCFCheckFunctionPointer              uint64                       // The VA where Control Flow Guard check-function pointer is stored.
	GuardCFDispatchFunctionPointer           uint64                       // The VA where Control Flow Guard dispatch-function pointer is stored.
	GuardCFFunctionTable                     uint64                       // The VA of the sorted table of RVAs of each Control Flow Guard function in the image.
	GuardCFFunctionCount                     uint64                       // The count of unique RVAs in the above table.
	GuardFlags                               uint32                       // Control Flow Guard related flags.
	CodeIntegrity                            ImageLoadConfigCodeIntegrity // Code integrity information.
	GuardAddressTakenIatEntryTable           uint64                       // The VA where Control Flow Guard address taken IAT table is stored.
	GuardAddressTakenIatEntryCount           uint64                       // The count of unique RVAs in the above table.
	GuardLongJumpTargetTable                 uint64                       // The VA where Control Flow Guard long jump target table is stored.
	GuardLongJumpTargetCount                 uint64                       // The count of unique RVAs in the above table.
	DynamicValueRelocTable                   uint64
	CHPEMetadataPointer                      uint64
	GuardRFFailureRoutine                    uint64
	GuardRFFailureRoutineFunctionPointer     uint64
	DynamicValueRelocTableOffset             uint32
	DynamicValueRelocTableSection            uint16
	Reserved2                                uint16
	GuardRFVerifyStackPointerFunctionPointer uint64
	HotPatchTableOffset                      uint32
	Reserved3                                uint32
	EnclaveConfigurationPointer              uint64
}

func (pe *File) parseLoadConfigDirectory(rva, size uint32) error {

	var loadConfig interface{}

	if pe.Is64 {
		loadCfg64 := ImageLoadConfigDirectory64{}
		loadCfgSize := uint32(binary.Size(loadCfg64))
		fileOffset := pe.getOffsetFromRva(rva)
		buf := bytes.NewReader(pe.data[fileOffset : fileOffset+loadCfgSize])
		err := binary.Read(buf, binary.LittleEndian, &loadCfg64)
		if err != nil {
			return err
		}
		loadConfig = loadCfg64
	} else {
		loadCfg32 := ImageLoadConfigDirectory32{}
		loadCfgSize := uint32(binary.Size(loadCfg32))
		fileOffset := pe.getOffsetFromRva(rva)
		buf := bytes.NewReader(pe.data[fileOffset : fileOffset+loadCfgSize])
		err := binary.Read(buf, binary.LittleEndian, &loadCfg32)
		if err != nil {
			return err
		}
		loadConfig = loadCfg32
	}

	pe.LoadConfig = loadConfig
	return nil
}
