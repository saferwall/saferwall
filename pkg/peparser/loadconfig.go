// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

const (

	// The GuardFlags field contains a combination of one or more of the
	// following flags and subfields:

	// ImageGuardCfInstrumented indicates that the module performs control flow
	// integrity checks using system-supplied support.
	ImageGuardCfInstrumented = 0x00000100

	// ImageGuardCfWInstrumented indicates that the module performs control
	// flow and write integrity checks.
	ImageGuardCfWInstrumented = 0x00000200

	// ImageGuardCfFunctionTablePresent indicates that the module contains
	// valid control flow target metadata.
	ImageGuardCfFunctionTablePresent = 0x00000400

	// ImageGuardSecurityCookieUnused indicates that the module does not make
	// use of the /GS security cookie.
	ImageGuardSecurityCookieUnused = 0x00000800

	// ImageGuardProtectDelayloadIAT indicates that the module supports read
	// only delay load IAT.
	ImageGuardProtectDelayloadIAT = 0x00001000

	// ImageGuardDelayloadIATInItsOwnSection indicates that the Delayload
	// import table in its own .didat section (with nothing else in it) that
	// can be freely reprotected.
	ImageGuardDelayloadIATInItsOwnSection = 0x00002000

	// ImageGuardCfExportSuppressionInfoPresent indicates that the module
	// contains suppressed export information. This also infers that the
	// address taken IAT table is also present in the load config.
	ImageGuardCfExportSuppressionInfoPresent = 0x00004000

	// ImageGuardCfEnableExportSuppression indicates that the module enables
	// suppression of exports.
	ImageGuardCfEnableExportSuppression = 0x00008000

	// ImageGuardCfLongjumpTablePresent indicates that the module contains
	// longjmp target information.
	ImageGuardCfLongjumpTablePresent = 0x00010000

	// ImageGuardCfFnctionTableSizeMask indicates that the mask for the
	// subfield that contains the stride of Control Flow Guard function table
	// entries (that is, the additional count of bytes per table entry).
	ImageGuardCfFnctionTableSizeMask = 0xF0000000

	// ImageGuardCfFnctionTableSizeShift indicates the shift to right-justify
	// Guard CF function table stride.
	ImageGuardCfFnctionTableSizeShift = 28

	// ImageGuardFlagFIDSupressed indicates that the call target is explicitly
	// suppressed (do not treat it as valid for purposes of CFG)
	ImageGuardFlagFIDSupressed = 0x1

	// ImageGuardFlagExportSupressed indicates that the call target is export
	// suppressed. See Export suppression for more details
	ImageGuardFlagExportSupressed = 0x2
)

type CFGFunction struct {
	Target      uint32
	Flags       *uint8
	Description string
}

type CFGIATEntry struct {
	RVA         uint32
	IATValue    uint32
	INTValue    uint32
	Description string
}

type LoadConfig struct {
	LoadCfgStruct interface{}
	SEH           []uint32
	GFIDS         []CFGFunction
	CFGIAT        []CFGIATEntry
	CFGLongJump   []uint32
	CHPE          HybridPE
}

// https://www.virtualbox.org/svn/vbox/trunk/include/iprt/formats/pecoff.h

// ImageLoadConfigCodeIntegrity Code Integrity in loadconfig (CI).
type ImageLoadConfigCodeIntegrity struct {
	Flags         uint16 // Flags to indicate if CI information is available, etc.
	Catalog       uint16 // 0xFFFF means not available
	CatalogOffset uint32
	Reserved      uint32 // Additional bitmask to be defined later
}

// ImageLoadConfigDirectory32v1 size is 0x40.
type ImageLoadConfigDirectory32v1 struct {
	Size                          uint32
	TimeDateStamp                 uint32
	MajorVersion                  uint16
	MinorVersion                  uint16
	GlobalFlagsClear              uint32
	GlobalFlagsSet                uint32
	CriticalSectionDefaultTimeout uint32
	DeCommitFreeBlockThreshold    uint32
	DeCommitTotalFreeThreshold    uint32
	LockPrefixTable               uint32
	MaximumAllocationSize         uint32
	VirtualMemoryThreshold        uint32
	ProcessHeapFlags              uint32
	ProcessAffinityMask           uint32
	CSDVersion                    uint16
	DependentLoadFlags            uint16
	EditList                      uint32
	SecurityCookie                uint32
}

// ImageLoadConfigDirectory32v2 size is 0x48.
type ImageLoadConfigDirectory32v2 struct {
	Size                          uint32
	TimeDateStamp                 uint32
	MajorVersion                  uint16
	MinorVersion                  uint16
	GlobalFlagsClear              uint32
	GlobalFlagsSet                uint32
	CriticalSectionDefaultTimeout uint32
	DeCommitFreeBlockThreshold    uint32
	DeCommitTotalFreeThreshold    uint32
	LockPrefixTable               uint32
	MaximumAllocationSize         uint32
	VirtualMemoryThreshold        uint32
	ProcessHeapFlags              uint32
	ProcessAffinityMask           uint32
	CSDVersion                    uint16
	DependentLoadFlags            uint16
	EditList                      uint32
	SecurityCookie                uint32
	SEHandlerTable                uint32
	SEHandlerCount                uint32
}

// ImageLoadConfigDirectory32v3 size is 0x5C.
type ImageLoadConfigDirectory32v3 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint32
	DeCommitTotalFreeThreshold     uint32
	LockPrefixTable                uint32
	MaximumAllocationSize          uint32
	VirtualMemoryThreshold         uint32
	ProcessHeapFlags               uint32
	ProcessAffinityMask            uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint32
	SecurityCookie                 uint32
	SEHandlerTable                 uint32
	SEHandlerCount                 uint32
	GuardCFCheckFunctionPointer    uint32
	GuardCFDispatchFunctionPointer uint32
	GuardCFFunctionTable           uint32
	GuardCFFunctionCount           uint32
	GuardFlags                     uint32
}

// ImageLoadConfigDirectory32v4 size is 0x6c
// since Windows 10 (preview 9879)
type ImageLoadConfigDirectory32v4 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint32
	DeCommitTotalFreeThreshold     uint32
	LockPrefixTable                uint32
	MaximumAllocationSize          uint32
	VirtualMemoryThreshold         uint32
	ProcessHeapFlags               uint32
	ProcessAffinityMask            uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint32
	SecurityCookie                 uint32
	SEHandlerTable                 uint32
	SEHandlerCount                 uint32
	GuardCFCheckFunctionPointer    uint32
	GuardCFDispatchFunctionPointer uint32
	GuardCFFunctionTable           uint32
	GuardCFFunctionCount           uint32
	GuardFlags                     uint32
	CodeIntegrity                  ImageLoadConfigCodeIntegrity
}

// ImageLoadConfigDirectory32v5 size is 0x78
// since Windows 10 build 14286 (or maybe earlier).
type ImageLoadConfigDirectory32v5 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint32
	DeCommitTotalFreeThreshold     uint32
	LockPrefixTable                uint32
	MaximumAllocationSize          uint32
	VirtualMemoryThreshold         uint32
	ProcessHeapFlags               uint32
	ProcessAffinityMask            uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint32
	SecurityCookie                 uint32
	SEHandlerTable                 uint32
	SEHandlerCount                 uint32
	GuardCFCheckFunctionPointer    uint32
	GuardCFDispatchFunctionPointer uint32
	GuardCFFunctionTable           uint32
	GuardCFFunctionCount           uint32
	GuardFlags                     uint32
	CodeIntegrity                  ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable uint32
	GuardAddressTakenIatEntryCount uint32
	GuardLongJumpTargetTable       uint32
	GuardLongJumpTargetCount       uint32
}

// ImageLoadConfigDirectory32v6 size is 0x80
// since Windows 10 build 14383 (or maybe earlier).
type ImageLoadConfigDirectory32v6 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint32
	DeCommitTotalFreeThreshold     uint32
	LockPrefixTable                uint32
	MaximumAllocationSize          uint32
	VirtualMemoryThreshold         uint32
	ProcessHeapFlags               uint32
	ProcessAffinityMask            uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint32
	SecurityCookie                 uint32
	SEHandlerTable                 uint32
	SEHandlerCount                 uint32
	GuardCFCheckFunctionPointer    uint32
	GuardCFDispatchFunctionPointer uint32
	GuardCFFunctionTable           uint32
	GuardCFFunctionCount           uint32
	GuardFlags                     uint32
	CodeIntegrity                  ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable uint32
	GuardAddressTakenIatEntryCount uint32
	GuardLongJumpTargetTable       uint32
	GuardLongJumpTargetCount       uint32
	DynamicValueRelocTable         uint32
	HybridMetadataPointer          uint32
}

// ImageLoadConfigDirectory32v7 size is 0x90
// since Windows 10 build 14901 (or maybe earlier)
type ImageLoadConfigDirectory32v7 struct {
	Size                                 uint32
	TimeDateStamp                        uint32
	MajorVersion                         uint16
	MinorVersion                         uint16
	GlobalFlagsClear                     uint32
	GlobalFlagsSet                       uint32
	CriticalSectionDefaultTimeout        uint32
	DeCommitFreeBlockThreshold           uint32
	DeCommitTotalFreeThreshold           uint32
	LockPrefixTable                      uint32
	MaximumAllocationSize                uint32
	VirtualMemoryThreshold               uint32
	ProcessHeapFlags                     uint32
	ProcessAffinityMask                  uint32
	CSDVersion                           uint16
	DependentLoadFlags                   uint16
	EditList                             uint32
	SecurityCookie                       uint32
	SEHandlerTable                       uint32
	SEHandlerCount                       uint32
	GuardCFCheckFunctionPointer          uint32
	GuardCFDispatchFunctionPointer       uint32
	GuardCFFunctionTable                 uint32
	GuardCFFunctionCount                 uint32
	GuardFlags                           uint32
	CodeIntegrity                        ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable       uint32
	GuardAddressTakenIatEntryCount       uint32
	GuardLongJumpTargetTable             uint32
	GuardLongJumpTargetCount             uint32
	DynamicValueRelocTable               uint32
	CHPEMetadataPointer                  uint32
	GuardRFFailureRoutine                uint32
	GuardRFFailureRoutineFunctionPointer uint32
	DynamicValueRelocTableOffset         uint32
	DynamicValueRelocTableSection        uint16
	Reserved2                            uint16
}

// ImageLoadConfigDirectory32v8 size is 0x98
// since Windows 10 build 15002 (or maybe earlier).
type ImageLoadConfigDirectory32v8 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint32
	DeCommitTotalFreeThreshold               uint32
	LockPrefixTable                          uint32
	MaximumAllocationSize                    uint32
	VirtualMemoryThreshold                   uint32
	ProcessHeapFlags                         uint32
	ProcessAffinityMask                      uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint32
	SecurityCookie                           uint32
	SEHandlerTable                           uint32
	SEHandlerCount                           uint32
	GuardCFCheckFunctionPointer              uint32
	GuardCFDispatchFunctionPointer           uint32
	GuardCFFunctionTable                     uint32
	GuardCFFunctionCount                     uint32
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint32
	GuardAddressTakenIatEntryCount           uint32
	GuardLongJumpTargetTable                 uint32
	GuardLongJumpTargetCount                 uint32
	DynamicValueRelocTable                   uint32
	CHPEMetadataPointer                      uint32
	GuardRFFailureRoutine                    uint32
	GuardRFFailureRoutineFunctionPointer     uint32
	DynamicValueRelocTableOffset             uint32
	DynamicValueRelocTableSection            uint16
	Reserved2                                uint16
	GuardRFVerifyStackPointerFunctionPointer uint32
	HotPatchTableOffset                      uint32
}

// ImageLoadConfigDirectory32v9 size is 0xA0
// since Windows 10 build 16237 (or maybe earlier).
type ImageLoadConfigDirectory32v9 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint32
	DeCommitTotalFreeThreshold               uint32
	LockPrefixTable                          uint32
	MaximumAllocationSize                    uint32
	VirtualMemoryThreshold                   uint32
	ProcessHeapFlags                         uint32
	ProcessAffinityMask                      uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint32
	SecurityCookie                           uint32
	SEHandlerTable                           uint32
	SEHandlerCount                           uint32
	GuardCFCheckFunctionPointer              uint32
	GuardCFDispatchFunctionPointer           uint32
	GuardCFFunctionTable                     uint32
	GuardCFFunctionCount                     uint32
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint32
	GuardAddressTakenIatEntryCount           uint32
	GuardLongJumpTargetTable                 uint32
	GuardLongJumpTargetCount                 uint32
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

// ImageLoadConfigDirectory32v10 size is 0xA4
// since Windows 10 build 18362  (or maybe earlier).
type ImageLoadConfigDirectory32v10 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint32
	DeCommitTotalFreeThreshold               uint32
	LockPrefixTable                          uint32
	MaximumAllocationSize                    uint32
	VirtualMemoryThreshold                   uint32
	ProcessHeapFlags                         uint32
	ProcessAffinityMask                      uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint32
	SecurityCookie                           uint32
	SEHandlerTable                           uint32
	SEHandlerCount                           uint32
	GuardCFCheckFunctionPointer              uint32
	GuardCFDispatchFunctionPointer           uint32
	GuardCFFunctionTable                     uint32
	GuardCFFunctionCount                     uint32
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint32
	GuardAddressTakenIatEntryCount           uint32
	GuardLongJumpTargetTable                 uint32
	GuardLongJumpTargetCount                 uint32
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
	VolatileMetadataPointer                  uint32
}

// ImageLoadConfigDirectory32v11 size is 0xAC
// since Windows 10 build 19564   (or maybe earlier).
type ImageLoadConfigDirectory32v11 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint32
	DeCommitTotalFreeThreshold               uint32
	LockPrefixTable                          uint32
	MaximumAllocationSize                    uint32
	VirtualMemoryThreshold                   uint32
	ProcessHeapFlags                         uint32
	ProcessAffinityMask                      uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint32
	SecurityCookie                           uint32
	SEHandlerTable                           uint32
	SEHandlerCount                           uint32
	GuardCFCheckFunctionPointer              uint32
	GuardCFDispatchFunctionPointer           uint32
	GuardCFFunctionTable                     uint32
	GuardCFFunctionCount                     uint32
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint32
	GuardAddressTakenIatEntryCount           uint32
	GuardLongJumpTargetTable                 uint32
	GuardLongJumpTargetCount                 uint32
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
	VolatileMetadataPointer                  uint32
	GuardEHContinuationTable                 uint32
	GuardEHContinuationCount                 uint32
}

// ImageLoadConfigDirectory32v12 size is 0xB8
// since Visual C++ 2019 / RS5_IMAGE_LOAD_CONFIG_DIRECTORY32.
type ImageLoadConfigDirectory32v12 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint32
	DeCommitTotalFreeThreshold               uint32
	LockPrefixTable                          uint32
	MaximumAllocationSize                    uint32
	VirtualMemoryThreshold                   uint32
	ProcessHeapFlags                         uint32
	ProcessAffinityMask                      uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint32
	SecurityCookie                           uint32
	SEHandlerTable                           uint32
	SEHandlerCount                           uint32
	GuardCFCheckFunctionPointer              uint32
	GuardCFDispatchFunctionPointer           uint32
	GuardCFFunctionTable                     uint32
	GuardCFFunctionCount                     uint32
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint32
	GuardAddressTakenIatEntryCount           uint32
	GuardLongJumpTargetTable                 uint32
	GuardLongJumpTargetCount                 uint32
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
	VolatileMetadataPointer                  uint32
	GuardEHContinuationTable                 uint32
	GuardEHContinuationCount                 uint32
	GuardXFGCheckFunctionPointer             uint32
	GuardXFGDispatchFunctionPointer          uint32
	GuardXFGTableDispatchFunctionPointer     uint32
}

// ImageLoadConfigDirectory32 Contains the load configuration data of an image
// for x86 binaries.
type ImageLoadConfigDirectory32 struct {
	// The actual size of the structure inclusive. May differ from the size
	// given in the data directory for Windows XP and earlier compatibility.
	Size uint32

	// Date and time stamp value.
	TimeDateStamp uint32

	// Major version number.
	MajorVersion uint16

	// Minor version number.
	MinorVersion uint16

	// The global loader flags to clear for this process as the loader starts
	// the process.
	GlobalFlagsClear uint32

	// The global loader flags to set for this process as the loader starts the
	// process.
	GlobalFlagsSet uint32

	// The default timeout value to use for this process's critical sections
	// that are abandoned.
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

	DynamicValueRelocTable uint32

	// Not sure when this was renamed from HybridMetadataPointer.
	CHPEMetadataPointer uint32

	GuardRFFailureRoutine                    uint32
	GuardRFFailureRoutineFunctionPointer     uint32
	DynamicValueRelocTableOffset             uint32
	DynamicValueRelocTableSection            uint16
	Reserved2                                uint16
	GuardRFVerifyStackPointerFunctionPointer uint32
	HotPatchTableOffset                      uint32
	Reserved3                                uint32
	EnclaveConfigurationPointer              uint32
	VolatileMetadataPointer                  uint32
	GuardEHContinuationTable                 uint32
	GuardEHContinuationCount                 uint32
	GuardXFGCheckFunctionPointer             uint32
	GuardXFGDispatchFunctionPointer          uint32
	GuardXFGTableDispatchFunctionPointer     uint32
}

// ImageLoadConfigDirectory64v2 is the first structure for x64.
// No _IMAGE_LOAD_CONFIG_DIRECTORY64_V1 exists
type ImageLoadConfigDirectory64v2 struct {
	Size                          uint32
	TimeDateStamp                 uint32
	MajorVersion                  uint16
	MinorVersion                  uint16
	GlobalFlagsClear              uint32
	GlobalFlagsSet                uint32
	CriticalSectionDefaultTimeout uint32
	DeCommitFreeBlockThreshold    uint64
	DeCommitTotalFreeThreshold    uint64
	LockPrefixTable               uint64
	MaximumAllocationSize         uint64
	VirtualMemoryThreshold        uint64
	ProcessAffinityMask           uint64
	ProcessHeapFlags              uint32
	CSDVersion                    uint16
	DependentLoadFlags            uint16
	EditList                      uint64
	SecurityCookie                uint64
	SEHandlerTable                uint64
	SEHandlerCount                uint64
}

// ImageLoadConfigDirectory64v3 added #pragma pack(4).
type ImageLoadConfigDirectory64v3 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint64
	DeCommitTotalFreeThreshold     uint64
	LockPrefixTable                uint64
	MaximumAllocationSize          uint64
	VirtualMemoryThreshold         uint64
	ProcessAffinityMask            uint64
	ProcessHeapFlags               uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint64
	SecurityCookie                 uint64
	SEHandlerTable                 uint64
	SEHandlerCount                 uint64
	GuardCFCheckFunctionPointer    uint64
	GuardCFDispatchFunctionPointer uint64
	GuardCFFunctionTable           uint64
	GuardCFFunctionCount           uint64
	GuardFlags                     uint32
}

// ImageLoadConfigDirectory64v4 for binaries compiled
// since Windows 10 build 9879 (or maybe earlier).
type ImageLoadConfigDirectory64v4 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint64
	DeCommitTotalFreeThreshold     uint64
	LockPrefixTable                uint64
	MaximumAllocationSize          uint64
	VirtualMemoryThreshold         uint64
	ProcessAffinityMask            uint64
	ProcessHeapFlags               uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint64
	SecurityCookie                 uint64
	SEHandlerTable                 uint64
	SEHandlerCount                 uint64
	GuardCFCheckFunctionPointer    uint64
	GuardCFDispatchFunctionPointer uint64
	GuardCFFunctionTable           uint64
	GuardCFFunctionCount           uint64
	GuardFlags                     uint32
	CodeIntegrity                  ImageLoadConfigCodeIntegrity
}

// ImageLoadConfigDirectory64v5 for binaries compiled
// since Windows 10 build 14286 (or maybe earlier).
type ImageLoadConfigDirectory64v5 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint64
	DeCommitTotalFreeThreshold     uint64
	LockPrefixTable                uint64
	MaximumAllocationSize          uint64
	VirtualMemoryThreshold         uint64
	ProcessAffinityMask            uint64
	ProcessHeapFlags               uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint64
	SecurityCookie                 uint64
	SEHandlerTable                 uint64
	SEHandlerCount                 uint64
	GuardCFCheckFunctionPointer    uint64
	GuardCFDispatchFunctionPointer uint64
	GuardCFFunctionTable           uint64
	GuardCFFunctionCount           uint64
	GuardFlags                     uint32
	CodeIntegrity                  ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable uint64
	GuardAddressTakenIatEntryCount uint64
	GuardLongJumpTargetTable       uint64
	GuardLongJumpTargetCount       uint64
}

// ImageLoadConfigDirectory64v6 for binaries compiled
// since Windows 10 build 14393 (or maybe earlier).
type ImageLoadConfigDirectory64v6 struct {
	Size                           uint32
	TimeDateStamp                  uint32
	MajorVersion                   uint16
	MinorVersion                   uint16
	GlobalFlagsClear               uint32
	GlobalFlagsSet                 uint32
	CriticalSectionDefaultTimeout  uint32
	DeCommitFreeBlockThreshold     uint64
	DeCommitTotalFreeThreshold     uint64
	LockPrefixTable                uint64
	MaximumAllocationSize          uint64
	VirtualMemoryThreshold         uint64
	ProcessAffinityMask            uint64
	ProcessHeapFlags               uint32
	CSDVersion                     uint16
	DependentLoadFlags             uint16
	EditList                       uint64
	SecurityCookie                 uint64
	SEHandlerTable                 uint64
	SEHandlerCount                 uint64
	GuardCFCheckFunctionPointer    uint64
	GuardCFDispatchFunctionPointer uint64
	GuardCFFunctionTable           uint64
	GuardCFFunctionCount           uint64
	GuardFlags                     uint32
	CodeIntegrity                  ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable uint64
	GuardAddressTakenIatEntryCount uint64
	GuardLongJumpTargetTable       uint64
	GuardLongJumpTargetCount       uint64
	DynamicValueRelocTable         uint64
	HybridMetadataPointer          uint64
}

// ImageLoadConfigDirectory64v7 for binaries compiled
// since Windows 10 build 14901 (or maybe earlier).
type ImageLoadConfigDirectory64v7 struct {
	Size                                 uint32
	TimeDateStamp                        uint32
	MajorVersion                         uint16
	MinorVersion                         uint16
	GlobalFlagsClear                     uint32
	GlobalFlagsSet                       uint32
	CriticalSectionDefaultTimeout        uint32
	DeCommitFreeBlockThreshold           uint64
	DeCommitTotalFreeThreshold           uint64
	LockPrefixTable                      uint64
	MaximumAllocationSize                uint64
	VirtualMemoryThreshold               uint64
	ProcessAffinityMask                  uint64
	ProcessHeapFlags                     uint32
	CSDVersion                           uint16
	DependentLoadFlags                   uint16
	EditList                             uint64
	SecurityCookie                       uint64
	SEHandlerTable                       uint64
	SEHandlerCount                       uint64
	GuardCFCheckFunctionPointer          uint64
	GuardCFDispatchFunctionPointer       uint64
	GuardCFFunctionTable                 uint64
	GuardCFFunctionCount                 uint64
	GuardFlags                           uint32
	CodeIntegrity                        ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable       uint64
	GuardAddressTakenIatEntryCount       uint64
	GuardLongJumpTargetTable             uint64
	GuardLongJumpTargetCount             uint64
	DynamicValueRelocTable               uint64
	CHPEMetadataPointer                  uint64
	GuardRFFailureRoutine                uint64
	GuardRFFailureRoutineFunctionPointer uint64
	DynamicValueRelocTableOffset         uint32
	DynamicValueRelocTableSection        uint16
	Reserved2                            uint16
}

// ImageLoadConfigDirectory64v8 for binaries compiled
// since Windows 10 build 15002 (or maybe earlier).
// #pragma pack(4) available here.
type ImageLoadConfigDirectory64v8 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint64
	DeCommitTotalFreeThreshold               uint64
	LockPrefixTable                          uint64
	MaximumAllocationSize                    uint64
	VirtualMemoryThreshold                   uint64
	ProcessAffinityMask                      uint64
	ProcessHeapFlags                         uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint64
	SecurityCookie                           uint64
	SEHandlerTable                           uint64
	SEHandlerCount                           uint64
	GuardCFCheckFunctionPointer              uint64
	GuardCFDispatchFunctionPointer           uint64
	GuardCFFunctionTable                     uint64
	GuardCFFunctionCount                     uint64
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint64
	GuardAddressTakenIatEntryCount           uint64
	GuardLongJumpTargetTable                 uint64
	GuardLongJumpTargetCount                 uint64
	DynamicValueRelocTable                   uint64
	CHPEMetadataPointer                      uint64
	GuardRFFailureRoutine                    uint64
	GuardRFFailureRoutineFunctionPointer     uint64
	DynamicValueRelocTableOffset             uint32
	DynamicValueRelocTableSection            uint16
	Reserved2                                uint16
	GuardRFVerifyStackPointerFunctionPointer uint64
	HotPatchTableOffset                      uint32
}

// ImageLoadConfigDirectory64v9 for binaries compiled
// since Windows 10 build 15002 (or maybe earlier).
// #pragma pack(4) was taken.
type ImageLoadConfigDirectory64v9 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint64
	DeCommitTotalFreeThreshold               uint64
	LockPrefixTable                          uint64
	MaximumAllocationSize                    uint64
	VirtualMemoryThreshold                   uint64
	ProcessAffinityMask                      uint64
	ProcessHeapFlags                         uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint64
	SecurityCookie                           uint64
	SEHandlerTable                           uint64
	SEHandlerCount                           uint64
	GuardCFCheckFunctionPointer              uint64
	GuardCFDispatchFunctionPointer           uint64
	GuardCFFunctionTable                     uint64
	GuardCFFunctionCount                     uint64
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint64
	GuardAddressTakenIatEntryCount           uint64
	GuardLongJumpTargetTable                 uint64
	GuardLongJumpTargetCount                 uint64
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

// ImageLoadConfigDirectory64v10 for binaries compiled
// since Windows 10 build 18362 (or maybe earlier).
type ImageLoadConfigDirectory64v10 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint64
	DeCommitTotalFreeThreshold               uint64
	LockPrefixTable                          uint64
	MaximumAllocationSize                    uint64
	VirtualMemoryThreshold                   uint64
	ProcessAffinityMask                      uint64
	ProcessHeapFlags                         uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint64
	SecurityCookie                           uint64
	SEHandlerTable                           uint64
	SEHandlerCount                           uint64
	GuardCFCheckFunctionPointer              uint64
	GuardCFDispatchFunctionPointer           uint64
	GuardCFFunctionTable                     uint64
	GuardCFFunctionCount                     uint64
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint64
	GuardAddressTakenIatEntryCount           uint64
	GuardLongJumpTargetTable                 uint64
	GuardLongJumpTargetCount                 uint64
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
	VolatileMetadataPointer                  uint64
}

// ImageLoadConfigDirectory64v11 for binaries compiled
// since Windows 10 build 19534 (or maybe earlier).
type ImageLoadConfigDirectory64v11 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint64
	DeCommitTotalFreeThreshold               uint64
	LockPrefixTable                          uint64
	MaximumAllocationSize                    uint64
	VirtualMemoryThreshold                   uint64
	ProcessAffinityMask                      uint64
	ProcessHeapFlags                         uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint64
	SecurityCookie                           uint64
	SEHandlerTable                           uint64
	SEHandlerCount                           uint64
	GuardCFCheckFunctionPointer              uint64
	GuardCFDispatchFunctionPointer           uint64
	GuardCFFunctionTable                     uint64
	GuardCFFunctionCount                     uint64
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint64
	GuardAddressTakenIatEntryCount           uint64
	GuardLongJumpTargetTable                 uint64
	GuardLongJumpTargetCount                 uint64
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
	VolatileMetadataPointer                  uint64
	GuardEHContinuationTable                 uint64
	GuardEHContinuationCount                 uint64
}

// ImageLoadConfigDirectory64v12 for binaries compiled
// since Visual C++ 2019 / RS5_IMAGE_LOAD_CONFIG_DIRECTORY64.
type ImageLoadConfigDirectory64v12 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint64
	DeCommitTotalFreeThreshold               uint64
	LockPrefixTable                          uint64
	MaximumAllocationSize                    uint64
	VirtualMemoryThreshold                   uint64
	ProcessAffinityMask                      uint64
	ProcessHeapFlags                         uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint64
	SecurityCookie                           uint64
	SEHandlerTable                           uint64
	SEHandlerCount                           uint64
	GuardCFCheckFunctionPointer              uint64
	GuardCFDispatchFunctionPointer           uint64
	GuardCFFunctionTable                     uint64
	GuardCFFunctionCount                     uint64
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint64
	GuardAddressTakenIatEntryCount           uint64
	GuardLongJumpTargetTable                 uint64
	GuardLongJumpTargetCount                 uint64
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
	VolatileMetadataPointer                  uint64
	GuardEHContinuationTable                 uint64
	GuardEHContinuationCount                 uint64
	GuardXFGCheckFunctionPointer             uint64
	GuardXFGDispatchFunctionPointer          uint64
	GuardXFGTableDispatchFunctionPointer     uint64
}

// ImageLoadConfigDirectory64 Contains the load configuration data of an image
// for x64 binaries.
type ImageLoadConfigDirectory64 struct {
	Size                                     uint32
	TimeDateStamp                            uint32
	MajorVersion                             uint16
	MinorVersion                             uint16
	GlobalFlagsClear                         uint32
	GlobalFlagsSet                           uint32
	CriticalSectionDefaultTimeout            uint32
	DeCommitFreeBlockThreshold               uint64
	DeCommitTotalFreeThreshold               uint64
	LockPrefixTable                          uint64
	MaximumAllocationSize                    uint64
	VirtualMemoryThreshold                   uint64
	ProcessAffinityMask                      uint64
	ProcessHeapFlags                         uint32
	CSDVersion                               uint16
	DependentLoadFlags                       uint16
	EditList                                 uint64
	SecurityCookie                           uint64
	SEHandlerTable                           uint64
	SEHandlerCount                           uint64
	GuardCFCheckFunctionPointer              uint64
	GuardCFDispatchFunctionPointer           uint64
	GuardCFFunctionTable                     uint64
	GuardCFFunctionCount                     uint64
	GuardFlags                               uint32
	CodeIntegrity                            ImageLoadConfigCodeIntegrity
	GuardAddressTakenIatEntryTable           uint64
	GuardAddressTakenIatEntryCount           uint64
	GuardLongJumpTargetTable                 uint64
	GuardLongJumpTargetCount                 uint64
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
	VolatileMetadataPointer                  uint64
	GuardEHContinuationTable                 uint64
	GuardEHContinuationCount                 uint64
	GuardXFGCheckFunctionPointer             uint64
	GuardXFGDispatchFunctionPointer          uint64
	GuardXFGTableDispatchFunctionPointer     uint64
}

type HybridPE struct {
	Version                                  uint32
	CHPECodeAddressRangeOffset               uint32
	CHPECodeAddressRangeCount                uint32
	WowA64ExceptionHandlerFunctionPtr        uint32
	WowA64DispatchCallFunctionPtr            uint32
	WowA64DispatchIndirectCallFunctionPtr    uint32
	WowA64DispatchIndirectCallCfgFunctionPtr uint32
	WowA64DispatchRetFunctionPtr             uint32
	WowA64DispatchRetLeafFunctionPtr         uint32
	WowA64DispatchJumpFunctionPtr            uint32
	CompilerIATPointer                       uint32
	WowA64RDTSCFunctionPtr                   uint32
}

// The load configuration structure (IMAGE_LOAD_CONFIG_DIRECTORY) was formerly
// used in very limited cases in the Windows NT operating system itself to
// describe various features too difficult or too large to describe in the file
// header or optional header of the image. Current versions of the Microsoft
// linker and Windows XP and later versions of Windows use a new version of this
// structure for 32-bit x86-based systems that include reserved SEH technology.
// The data directory entry for a pre-reserved SEH load configuration structure
// must specify a particular size of the load configuration structure because
// the operating system loader always expects it to be a certain value. In that
// regard, the size is really only a version check. For compatibility with
// Windows XP and earlier versions of Windows, the size must be 64 for x86 images.
func (pe *File) parseLoadConfigDirectory(rva, size uint32) error {

	// As the laod config structure changes over time,
	// we first read it size to figure out which one we have to cast against.
	fileOffset := pe.getOffsetFromRva(rva)
	structSize, err := pe.ReadUint32(fileOffset)
	if err != nil {
		return err
	}

	// Use this helper function to print struct size.
	// PrintLoadConfigStruct()
	var loadCfg interface{}

	if pe.Is32 {
		switch int(structSize) {
		case 0x40:
			loadCfgv1 := ImageLoadConfigDirectory32v1{}
			err = pe.structUnpack(&loadCfgv1, fileOffset, structSize)
			loadCfg = loadCfgv1
		case 0x48:
			loadCfgv2 := ImageLoadConfigDirectory32v2{}
			err = pe.structUnpack(&loadCfgv2, fileOffset, structSize)
			loadCfg = loadCfgv2
		case 0x5c:
			loadCfgv3 := ImageLoadConfigDirectory32v3{}
			err = pe.structUnpack(&loadCfgv3, fileOffset, structSize)
			loadCfg = loadCfgv3
		case 0x68:
			loadCfgv4 := ImageLoadConfigDirectory32v4{}
			err = pe.structUnpack(&loadCfgv4, fileOffset, structSize)
			loadCfg = loadCfgv4
		case 0x78:
			loadCfgv5 := ImageLoadConfigDirectory32v5{}
			err = pe.structUnpack(&loadCfgv5, fileOffset, structSize)
			loadCfg = loadCfgv5
		case 0x80:
			loadCfgv6 := ImageLoadConfigDirectory32v6{}
			err = pe.structUnpack(&loadCfgv6, fileOffset, structSize)
			loadCfg = loadCfgv6
		case 0x90:
			loadCfgv7 := ImageLoadConfigDirectory32v7{}
			err = pe.structUnpack(&loadCfgv7, fileOffset, structSize)
			loadCfg = loadCfgv7
		case 0x98:
			loadCfgv8 := ImageLoadConfigDirectory32v8{}
			err = pe.structUnpack(&loadCfgv8, fileOffset, structSize)
			loadCfg = loadCfgv8
		case 0xa0:
			loadCfgv9 := ImageLoadConfigDirectory32v9{}
			err = pe.structUnpack(&loadCfgv9, fileOffset, structSize)
			loadCfg = loadCfgv9
		case 0xa4:
			loadCfgv10 := ImageLoadConfigDirectory32v10{}
			err = pe.structUnpack(&loadCfgv10, fileOffset, structSize)
			loadCfg = loadCfgv10
		case 0xac:
			loadCfgv11 := ImageLoadConfigDirectory32v11{}
			err = pe.structUnpack(&loadCfgv11, fileOffset, structSize)
			loadCfg = loadCfgv11
		case 0xb8:
			loadCfgv12 := ImageLoadConfigDirectory32v12{}
			err = pe.structUnpack(&loadCfgv12, fileOffset, structSize)
			loadCfg = loadCfgv12
		default:
			// We use the oldest load config.
			loadCfg32 := ImageLoadConfigDirectory32v1{}
			err = pe.structUnpack(&loadCfg32, fileOffset, structSize)
			loadCfg = loadCfg32
		}
	} else {
		switch int(structSize) {
		case 0x70:
			loadCfgv2 := ImageLoadConfigDirectory64v2{}
			err = pe.structUnpack(&loadCfgv2, fileOffset, structSize)
			loadCfg = loadCfgv2
		case 0x94:
			loadCfgv3 := ImageLoadConfigDirectory64v3{}
			err = pe.structUnpack(&loadCfgv3, fileOffset, structSize)
			loadCfg = loadCfgv3
		case 0xa0:
			loadCfgv4 := ImageLoadConfigDirectory64v4{}
			err = pe.structUnpack(&loadCfgv4, fileOffset, structSize)
			loadCfg = loadCfgv4
		case 0xc0:
			loadCfgv5 := ImageLoadConfigDirectory64v5{}
			err = pe.structUnpack(&loadCfgv5, fileOffset, structSize)
			loadCfg = loadCfgv5
		case 0xd0:
			loadCfgv6 := ImageLoadConfigDirectory64v6{}
			err = pe.structUnpack(&loadCfgv6, fileOffset, structSize)
			loadCfg = loadCfgv6
		case 0xe8:
			loadCfgv7 := ImageLoadConfigDirectory64v7{}
			err = pe.structUnpack(&loadCfgv7, fileOffset, structSize)
			loadCfg = loadCfgv7
		case 0xf4:
			loadCfgv8 := ImageLoadConfigDirectory64v8{}
			err = pe.structUnpack(&loadCfgv8, fileOffset, structSize)
			loadCfg = loadCfgv8
		case 0x100:
			loadCfgv9 := ImageLoadConfigDirectory64v9{}
			err = pe.structUnpack(&loadCfgv9, fileOffset, structSize)
			loadCfg = loadCfgv9
		case 0x108:
			loadCfgv10 := ImageLoadConfigDirectory64v10{}
			err = pe.structUnpack(&loadCfgv10, fileOffset, structSize)
			loadCfg = loadCfgv10
		case 0x118:
			loadCfgv11 := ImageLoadConfigDirectory64v11{}
			err = pe.structUnpack(&loadCfgv11, fileOffset, structSize)
			loadCfg = loadCfgv11
		case 0x130:
			loadCfgv12 := ImageLoadConfigDirectory64v12{}
			err = pe.structUnpack(&loadCfgv12, fileOffset, structSize)
			loadCfg = loadCfgv12
		default:
			// We use the oldest load config.
			loadCfg64 := ImageLoadConfigDirectory64v2{}
			err = pe.structUnpack(&loadCfg64, fileOffset, structSize)
			loadCfg = loadCfg64
		}
	}

	if err != nil {
		return err
	}

	// Save the load config struct.
	pe.LoadConfig.LoadCfgStruct = loadCfg

	// Retrieve SEH handlers if there are any..
	if pe.Is32 {
		handlers := pe.getSEHHandlers()
		pe.LoadConfig.SEH = handlers
	}

	// Retrieve Control Flow Guard Function Targets if there are any.
	pe.LoadConfig.GFIDS = pe.getControlFlowGuardFunctions()

	// Retrieve Control Flow Guard IAT entries if there are any.
	pe.LoadConfig.CFGIAT = pe.getControlFlowGuardIAT()

	// Retrive Long jump target functions if there are any.
	pe.LoadConfig.CFGLongJump = pe.getLongJumpTargetTable()

	// Retrieve compiled hybrid PE metadata if there are any.
	pe.LoadConfig.CHPE = pe.getHybridPE()

	return nil
}

// PrintLoadConfigStruct will print size of each load config struct.
func PrintLoadConfigStruct() {
	fmt.Printf("ImageLoadConfigDirectory32v1 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v1{}))
	fmt.Printf("ImageLoadConfigDirectory32v2 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v2{}))
	fmt.Printf("ImageLoadConfigDirectory32v3 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v3{}))
	fmt.Printf("ImageLoadConfigDirectory32v4 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v4{}))
	fmt.Printf("ImageLoadConfigDirectory32v5 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v5{}))
	fmt.Printf("ImageLoadConfigDirectory32v6 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v6{}))
	fmt.Printf("ImageLoadConfigDirectory32v7 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v7{}))
	fmt.Printf("ImageLoadConfigDirectory32v8 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v8{}))
	fmt.Printf("ImageLoadConfigDirectory32v9 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v9{}))
	fmt.Printf("ImageLoadConfigDirectory32v10 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v10{}))
	fmt.Printf("ImageLoadConfigDirectory32v11 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v11{}))
	fmt.Printf("ImageLoadConfigDirectory32v12 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory32v12{}))

	fmt.Printf("ImageLoadConfigDirectory64v2 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v2{}))
	fmt.Printf("ImageLoadConfigDirectory64v3 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v3{}))
	fmt.Printf("ImageLoadConfigDirectory64v4 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v4{}))
	fmt.Printf("ImageLoadConfigDirectory64v5 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v5{}))
	fmt.Printf("ImageLoadConfigDirectory64v6 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v6{}))
	fmt.Printf("ImageLoadConfigDirectory64v7 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v7{}))
	fmt.Printf("ImageLoadConfigDirectory64v8 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v8{}))
	fmt.Printf("ImageLoadConfigDirectory64v9 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v9{}))
	fmt.Printf("ImageLoadConfigDirectory64v10 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v10{}))
	fmt.Printf("ImageLoadConfigDirectory64v11 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v11{}))
	fmt.Printf("ImageLoadConfigDirectory64v12 size : 0x%x\n", binary.Size(ImageLoadConfigDirectory64v12{}))

}

// StringifyGuardFlags returns list of strings which describes the GuardFlags.
func StringifyGuardFlags(flags uint32) []string {
	var values []string
	guardFlagMap := map[uint32]string{
		ImageGuardCfInstrumented:                 "Instrumented",
		ImageGuardCfWInstrumented:                "WriteInstrumented",
		ImageGuardCfFunctionTablePresent:         "TargetMetadata",
		ImageGuardSecurityCookieUnused:           "SecurityCookieUnused",
		ImageGuardProtectDelayloadIAT:            "DelayloadIAT",
		ImageGuardDelayloadIATInItsOwnSection:    "DelayloadIATInItsOwnSection",
		ImageGuardCfExportSuppressionInfoPresent: "ExportSuppressionInfoPresent",
		ImageGuardCfEnableExportSuppression:      "EnableExportSuppression",
		ImageGuardCfLongjumpTablePresent:         "LongjumpTablePresent",
	}

	for k, s := range guardFlagMap {
		if k&flags != 0 {
			values = append(values, s)
		}
	}
	return values
}

func (pe *File) getSEHHandlers() []uint32 {

	var handlers []uint32
	v := reflect.ValueOf(pe.LoadConfig.LoadCfgStruct)

	// SEHandlerCount is found in index 19 of the struct.
	if v.NumField() > 19 {
		SEHandlerCount := uint32(v.Field(19).Uint())
		if SEHandlerCount > 0 {
			SEHandlerTable := uint32(v.Field(18).Uint())
			imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).ImageBase
			rva := SEHandlerTable - imageBase
			for i := uint32(0); i < SEHandlerCount; i++ {
				offset := pe.getOffsetFromRva(rva + i*4)
				handler, err := pe.ReadUint32(offset)
				if err != nil {
					return handlers
				}

				handlers = append(handlers, handler)
			}
		}
	}

	return handlers
}

func (pe *File) getControlFlowGuardFunctions() []CFGFunction {

	v := reflect.ValueOf(pe.LoadConfig.LoadCfgStruct)
	var GFIDS []CFGFunction
	var err error

	// GuardCFFunctionCount is found in index 23 of the struct.
	if v.NumField() > 23 {
		// The GFIDS table is an array of 4 + n bytes, where n is given by :
		// ((GuardFlags & IMAGE_GUARD_CF_FUNCTION_TABLE_SIZE_MASK) >>
		// IMAGE_GUARD_CF_FUNCTION_TABLE_SIZE_SHIFT).

		// This allows for extra metadata to be attached to CFG call targets in
		// the future. The only currently defined metadata is an optional 1-byte
		// extra flags field (“GFIDS flags”) that is attached to each GFIDS
		// entry if any call targets have metadata.
		GuardFlags := v.Field(24).Uint()
		n := (GuardFlags & ImageGuardCfFnctionTableSizeMask) >>
			ImageGuardCfFnctionTableSizeShift
		GuardCFFunctionCount := v.Field(23).Uint()
		if GuardCFFunctionCount > 0 {
			if pe.Is32 {
				GuardCFFunctionTable := uint32(v.Field(22).Uint())
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).ImageBase
				rva := GuardCFFunctionTable - imageBase
				offset := pe.getOffsetFromRva(rva)
				for i := uint32(1); i <= uint32(GuardCFFunctionCount); i++ {
					cfgFunction := CFGFunction{}
					var cfgFlags uint8
					cfgFunction.Target, err = pe.ReadUint32(offset)
					if err != nil {
						return GFIDS
					}
					if n > 0 {
						pe.structUnpack(&cfgFlags, offset+4, uint32(n))
						cfgFunction.Flags = &cfgFlags
						if cfgFlags == ImageGuardFlagFIDSupressed ||
							cfgFlags == ImageGuardFlagExportSupressed {
							exportName := pe.GetExportFunctionByRVA(cfgFunction.Target)
							cfgFunction.Description = exportName.Name
						}
					}

					GFIDS = append(GFIDS, cfgFunction)
					offset += 4 + uint32(n)
				}
			} else {
				GuardCFFunctionTable := v.Field(22).Uint()
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).ImageBase
				rva := uint32(GuardCFFunctionTable - imageBase)
				offset := pe.getOffsetFromRva(rva)
				for i := uint64(1); i <= GuardCFFunctionCount; i++ {
					var cfgFlags uint8
					cfgFunction := CFGFunction{}
					cfgFunction.Target, err = pe.ReadUint32(offset)
					if err != nil {
						return GFIDS
					}
					if n > 0 {
						pe.structUnpack(&cfgFlags, offset+4, uint32(n))
						cfgFunction.Flags = &cfgFlags
						if cfgFlags == ImageGuardFlagFIDSupressed ||
							cfgFlags == ImageGuardFlagExportSupressed {
							exportName := pe.GetExportFunctionByRVA(cfgFunction.Target)
							cfgFunction.Description = exportName.Name
						}
					}

					GFIDS = append(GFIDS, cfgFunction)
					offset += 4 + uint32(n)
				}
			}

		}
	}
	return GFIDS
}

func (pe *File) getControlFlowGuardIAT() []CFGIATEntry {

	v := reflect.ValueOf(pe.LoadConfig.LoadCfgStruct)
	var GFGIAT []CFGIATEntry
	var err error

	// GuardAddressTakenIatEntryCount is found in index 27 of the struct.
	if v.NumField() > 27 {
		// An image that supports CFG ES includes a GuardAddressTakenIatEntryTable
		// whose count is provided by the GuardAddressTakenIatEntryCount as part
		// of its load configuration directory. This table is structurally
		// formatted the same as the GFIDS table. It uses the same GuardFlags
		// IMAGE_GUARD_CF_FUNCTION_TABLE_SIZE_MASK mechanism to encode extra
		// optional metadata bytes in the address taken IAT table, though all
		// metadata bytes must be zero for the address taken IAT table and are
		// reserved.
		GuardFlags := v.Field(24).Uint()
		n := (GuardFlags & ImageGuardCfFnctionTableSizeMask) >>
			ImageGuardCfFnctionTableSizeShift
		GuardAddressTakenIatEntryCount := v.Field(27).Uint()
		if GuardAddressTakenIatEntryCount > 0 {
			if pe.Is32 {
				GuardAddressTakenIatEntryTable := uint32(v.Field(26).Uint())
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).ImageBase
				rva := GuardAddressTakenIatEntryTable - imageBase
				offset := pe.getOffsetFromRva(rva)
				for i := uint32(1); i <= uint32(GuardAddressTakenIatEntryCount); i++ {
					cfgIATEntry := CFGIATEntry{}
					cfgIATEntry.RVA, err = pe.ReadUint32(offset)
					if err != nil {
						return GFGIAT
					}
					imp, index := pe.GetImportEntryInfoByRVA(cfgIATEntry.RVA)
					cfgIATEntry.INTValue = uint32(imp.Functions[index].OriginalThunkValue)
					cfgIATEntry.IATValue = uint32(imp.Functions[index].ThunkValue)
					cfgIATEntry.Description = imp.Name + "!" + imp.Functions[index].Name
					GFGIAT = append(GFGIAT, cfgIATEntry)
					offset += 4 + uint32(n)
				}
			} else {
				GuardAddressTakenIatEntryTable := v.Field(26).Uint()
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).ImageBase
				rva := uint32(GuardAddressTakenIatEntryTable - imageBase)
				offset := pe.getOffsetFromRva(rva)
				for i := uint64(1); i <= GuardAddressTakenIatEntryCount; i++ {
					cfgIATEntry := CFGIATEntry{}
					cfgIATEntry.RVA, err = pe.ReadUint32(offset)
					if err != nil {
						return GFGIAT
					}
					imp, index := pe.GetImportEntryInfoByRVA(cfgIATEntry.RVA)
					cfgIATEntry.INTValue = uint32(imp.Functions[index].OriginalThunkValue)
					cfgIATEntry.IATValue = uint32(imp.Functions[index].ThunkValue)
					cfgIATEntry.Description = imp.Name + "!" + imp.Functions[index].Name
					GFGIAT = append(GFGIAT, cfgIATEntry)
					GFGIAT = append(GFGIAT, cfgIATEntry)
					offset += 4 + uint32(n)
				}
			}

		}
	}
	return GFGIAT
}

func (pe *File) getLongJumpTargetTable() []uint32 {

	v := reflect.ValueOf(pe.LoadConfig.LoadCfgStruct)
	var longJumpTargets []uint32

	// GuardLongJumpTargetCount is found in index 29 of the struct.
	if v.NumField() > 29 {
		// The long jump table represents a sorted array of RVAs that are valid
		// long jump targets. If a long jump target module sets
		// IMAGE_GUARD_CF_LONGJUMP_TABLE_PRESENT in its GuardFlags field, then
		// all long jump targets must be enumerated in the LongJumpTargetTable.
		GuardFlags := v.Field(24).Uint()
		n := (GuardFlags & ImageGuardCfFnctionTableSizeMask) >>
			ImageGuardCfFnctionTableSizeShift
		GuardLongJumpTargetCount := v.Field(29).Uint()
		if GuardLongJumpTargetCount > 0 {
			if pe.Is32 {
				GuardLongJumpTargetTable := uint32(v.Field(28).Uint())
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).ImageBase
				rva := GuardLongJumpTargetTable - imageBase
				offset := pe.getOffsetFromRva(rva)
				for i := uint32(1); i <= uint32(GuardLongJumpTargetCount); i++ {
					target, err := pe.ReadUint32(offset)
					if err != nil {
						return longJumpTargets
					}
					longJumpTargets = append(longJumpTargets, target)
					offset += 4 + uint32(n)
				}
			} else {
				GuardLongJumpTargetTable := v.Field(26).Uint()
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).ImageBase
				rva := uint32(GuardLongJumpTargetTable - imageBase)
				offset := pe.getOffsetFromRva(rva)
				for i := uint64(1); i <= GuardLongJumpTargetCount; i++ {
					target, err := pe.ReadUint32(offset)
					if err != nil {
						return longJumpTargets
					}
					longJumpTargets = append(longJumpTargets, target)
					offset += 4 + uint32(n)
				}
			}

		}
	}
	return longJumpTargets
}

func (pe *File) getHybridPE() HybridPE {
	v := reflect.ValueOf(pe.LoadConfig.LoadCfgStruct)
	hybridPE := HybridPE{}

	// CHPEMetadataPointer is found in index 31 of the struct.
	if v.NumField() > 31 {
		CHPEMetadataPointer := v.Field(31).Uint()
		if CHPEMetadataPointer != 0 {
			structSize := uint32(binary.Size(hybridPE))
			if pe.Is32 {
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader32).ImageBase
				rva := uint32(CHPEMetadataPointer) - imageBase
				fileOffset := pe.getOffsetFromRva(rva)
				err := pe.structUnpack(&hybridPE, fileOffset, structSize)
				if err != nil {
					return hybridPE
				}
			} else {
				imageBase := pe.NtHeader.OptionalHeader.(ImageOptionalHeader64).ImageBase
				rva := CHPEMetadataPointer - imageBase
				err := pe.structUnpack(&hybridPE, uint32(rva), structSize)
				if err != nil {
					return hybridPE
				}
			}

		}
	}

	return hybridPE
}
