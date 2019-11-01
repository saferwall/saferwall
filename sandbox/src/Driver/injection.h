#pragma once

#include <ntddk.h>
#include <ntimage.h>
#include "shared.h"

//////////////////////////////////////////////////////////////////////////
// Definitions.
//////////////////////////////////////////////////////////////////////////

#define INJ_MEMORY_TAG ' wfS'

#if defined(_M_AMD64) || defined(_M_ARM64)
# define INJ_CONFIG_SUPPORTS_WOW64
#endif


//////////////////////////////////////////////////////////////////////////
// Enumerations.
//////////////////////////////////////////////////////////////////////////

typedef enum _INJ_SYSTEM_DLL
{
	INJ_NOTHING_LOADED = 0x0000,
	INJ_SYSARM32_NTDLL_LOADED = 0x0001,
	INJ_SYCHPE32_NTDLL_LOADED = 0x0002,
	INJ_SYSWOW64_NTDLL_LOADED = 0x0004,
	INJ_SYSTEM32_NTDLL_LOADED = 0x0008,
	INJ_SYSTEM32_WOW64_LOADED = 0x0010,
	INJ_SYSTEM32_WOW64WIN_LOADED = 0x0020,
	INJ_SYSTEM32_WOW64CPU_LOADED = 0x0040,
	INJ_SYSTEM32_WOWARMHW_LOADED = 0x0080,
	INJ_SYSTEM32_XTAJIT_LOADED = 0x0100,
} INJ_SYSTEM_DLL;


typedef enum _INJ_ARCHITECTURE
{
	InjArchitectureX86,
	InjArchitectureX64,
	InjArchitectureARM32,
	InjArchitectureARM64,
	InjArchitectureMax,

#if defined(_M_IX86)
	InjArchitectureNative = InjArchitectureX86
#elif defined (_M_AMD64)
	InjArchitectureNative = InjArchitectureX64
#elif defined (_M_ARM64)
	InjArchitectureNative = InjArchitectureARM64
#endif
} INJ_ARCHITECTURE;


//////////////////////////////////////////////////////////////////////////
// Structures.
//////////////////////////////////////////////////////////////////////////

typedef struct _INJ_SYSTEM_DLL_DESCRIPTOR
{
	UNICODE_STRING  DllPath;
	INJ_SYSTEM_DLL  Flag;
}	INJ_SYSTEM_DLL_DESCRIPTOR, *PINJ_SYSTEM_DLL_DESCRIPTOR;


typedef struct _INJ_THUNK
{
	PVOID           Buffer;
	USHORT          Length;
} INJ_THUNK, *PINJ_THUNK;


typedef struct _INJECTION_INFO
{
	LIST_ENTRY  ListEntry;

	//
	// Process ID.
	//

	HANDLE      ProcessId;

	//
	// Combination of INJ_SYSTEM_DLL flags indicating
	// which DLLs has been already loaded into this
	// process.
	//

	ULONG       LoadedDlls;

	//
	// If true, the process has been already injected.
	//

	BOOLEAN     IsInjected;

	//
	// If true, trigger of the queued user APC will be
	// immediately forced upon next kernel->user transition.
	//

	BOOLEAN     ForceUserApc;

	//
	// Address of LdrLoadDll routine within ntdll.dll
	// (which ntdll.dll is selected is based on the INJ_METHOD).
	//

	PVOID       LdrLoadDllRoutineAddress;

} INJECTION_INFO, *PINJECTION_INFO;


typedef struct _INJ_SETTINGS
{
	//
	// Paths to the inject DLLs for each architecture.
	// Unsupported architectures (either by OS or the
	// method) can have empty string.
	//

	UNICODE_STRING  DllPath[InjArchitectureMax];

} INJ_SETTINGS, *PINJ_SETTINGS;


//
// Taken from ReactOS, used by InjpInitializeDllPaths.
//

typedef union
{
	WCHAR Name[sizeof(ULARGE_INTEGER) / sizeof(WCHAR)];
	ULARGE_INTEGER Alignment;
} ALIGNEDNAME;





//////////////////////////////////////////////////////////////////////////
// ke.h
//////////////////////////////////////////////////////////////////////////

typedef enum _KAPC_ENVIRONMENT
{
	OriginalApcEnvironment,
	AttachedApcEnvironment,
	CurrentApcEnvironment,
	InsertApcEnvironment
} KAPC_ENVIRONMENT;


typedef
VOID
(NTAPI *PKNORMAL_ROUTINE)(
	_In_ PVOID NormalContext,
	_In_ PVOID SystemArgument1,
	_In_ PVOID SystemArgument2
	);

typedef
VOID
(NTAPI *PKKERNEL_ROUTINE)(
	_In_ PKAPC Apc,
	_Inout_ PKNORMAL_ROUTINE* NormalRoutine,
	_Inout_ PVOID* NormalContext,
	_Inout_ PVOID* SystemArgument1,
	_Inout_ PVOID* SystemArgument2
	);

typedef
VOID
(NTAPI *PKRUNDOWN_ROUTINE) (
	_In_ PKAPC Apc
	);

NTKERNELAPI
VOID
NTAPI
KeInitializeApc(
	_Out_ PRKAPC Apc,
	_In_ PETHREAD Thread,
	_In_ KAPC_ENVIRONMENT Environment,
	_In_ PKKERNEL_ROUTINE KernelRoutine,
	_In_opt_ PKRUNDOWN_ROUTINE RundownRoutine,
	_In_opt_ PKNORMAL_ROUTINE NormalRoutine,
	_In_opt_ KPROCESSOR_MODE ApcMode,
	_In_opt_ PVOID NormalContext
);

NTKERNELAPI
BOOLEAN
NTAPI
KeInsertQueueApc(
	_Inout_ PRKAPC Apc,
	_In_opt_ PVOID SystemArgument1,
	_In_opt_ PVOID SystemArgument2,
	_In_ KPRIORITY Increment
);

NTKERNELAPI
BOOLEAN
NTAPI
KeAlertThread(
	_Inout_ PKTHREAD Thread,
	_In_ KPROCESSOR_MODE AlertMode
);

NTKERNELAPI
BOOLEAN
NTAPI
KeTestAlertThread(
	_In_ KPROCESSOR_MODE AlertMode
);


//////////////////////////////////////////////////////////////////////////
// ps.h
//////////////////////////////////////////////////////////////////////////

NTKERNELAPI
PVOID
NTAPI
PsGetProcessWow64Process(
	_In_ PEPROCESS Process
);

NTKERNELAPI
PCHAR
NTAPI
PsGetProcessImageFileName(
	_In_ PEPROCESS Process
);

NTKERNELAPI
BOOLEAN
NTAPI
PsIsProtectedProcess(
	_In_ PEPROCESS Process
);

NTKERNELAPI
USHORT
NTAPI
PsWow64GetProcessMachine(
	_In_ PEPROCESS Process
);


//////////////////////////////////////////////////////////////////////////
// ntrtl.h
//////////////////////////////////////////////////////////////////////////
#define RTL_DUPLICATE_UNICODE_STRING_NULL_TERMINATE (0x00000001)
#define RTL_DUPLICATE_UNICODE_STRING_ALLOCATE_NULL_STRING (0x00000002)

NTSYSAPI
NTSTATUS
NTAPI
RtlDuplicateUnicodeString(
	_In_ ULONG Flags,
	_In_ PUNICODE_STRING StringIn,
	_Out_ PUNICODE_STRING StringOut
);



NTSYSAPI
PVOID
NTAPI
RtlImageDirectoryEntryToData(
	_In_ PVOID BaseOfImage,
	_In_ BOOLEAN MappedAsImage,
	_In_ USHORT DirectoryEntry,
	_Out_ PULONG Size
);



//////////////////////////////////////////////////////////////////////////
// Prototypes.
//////////////////////////////////////////////////////////////////////////

PINJECTION_INFO
NTAPI
InjFindInjectionInfo(
	_In_ HANDLE ProcessId
);


VOID
NTAPI
InjRemoveInjectionInfoByProcessId(
	_In_ HANDLE ProcessId,
	_In_ BOOLEAN FreeMemory
);

BOOLEAN
NTAPI
InjCanInject(
	_In_ PINJECTION_INFO InjectionInfo
);

VOID
NTAPI
InjIsWantedSystemDllBeingLoaded(
	_In_ PINJECTION_INFO InjectionInfo,
	_In_ PUNICODE_STRING FullImageName,
	_In_ PIMAGE_INFO ImageInfo
);

NTSTATUS
NTAPI
InjInitialize(
	_In_ PUNICODE_STRING RegistryPath
);


BOOLEAN
NTAPI
RtlxSuffixUnicodeString(
	_In_ PUNICODE_STRING String1,
	_In_ PUNICODE_STRING String2,
	_In_ BOOLEAN CaseInSensitive
);


VOID
NTAPI
InjpInjectApcNormalRoutine(
	_In_ PVOID NormalContext,
	_In_ PVOID SystemArgument1,
	_In_ PVOID SystemArgument2
);


NTSTATUS
NTAPI
InjQueueApc(
	_In_ KPROCESSOR_MODE ApcMode,
	_In_ PKNORMAL_ROUTINE NormalRoutine,
	_In_ PVOID NormalContext,
	_In_ PVOID SystemArgument1,
	_In_ PVOID SystemArgument2
);


NTSTATUS
NTAPI
InjpInject(
	_In_ PINJECTION_INFO InjectionInfo,
	_In_ INJ_ARCHITECTURE Architecture,
	_In_ HANDLE SectionHandle,
	_In_ SIZE_T SectionSize
);

NTSTATUS
NTAPI
InjInject(
	_In_ PINJECTION_INFO InjectionInfo
);