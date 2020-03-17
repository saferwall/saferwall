#pragma once

#include "hooking.h"
#include "libloaderapi.h"
#include "fileapi.h"
#include "memoryapi.h"
#include "systemapi.h"
#include "synchapi.h"
#include "winregapi.h"
#include "processthreadsapi.h"
#include "ntifs.h"
#include "ole.h"
#include "net.h"

//
// Prototypes
//

BOOL
ProcessAttach();
BOOL
ProcessDetach();
BOOL
IsInsideHook();
VOID
EnterHookGuard();
VOID
ReleaseHookGuard();
BOOL
HookBegingTransation();
BOOL
HookCommitTransaction();
VOID
HookNtAPIs();
VOID
HookOleAPIs(BOOL Attach);
VOID
HookNetworkAPIs(BOOL Attach);
VOID
HookDll(PWCHAR DllName);
BOOL
SfwIsCalledFromSystemMemory(DWORD FramesToCapture);
NTSTATUS
SfwSymInit();

//
// Unfortunatelly sprintf-like functions are not exposed
// by ntdll.lib, which we're linking against.  We have to
// load them dynamically.
//

using __vsnwprintf_fn_t = int(__cdecl *)(wchar_t *buffer, size_t count, const wchar_t *format, ...);

using __snwprintf_fn_t = int(__cdecl *)(wchar_t *buffer, size_t count, const wchar_t *format, ...);

using strlen_fn_t = size_t(__cdecl *)(char const *buffer);

using pfn_wcsstr = wchar_t *(__cdecl *)(wchar_t *_String, wchar_t const *_SubStr);

//
// Structs
//

typedef struct HookInfo
{
    BOOL IsOleHooked;
    BOOL IsWinInetHooked;
    ULONG ExecutableModuleStart;
    ULONG ExecutableModuleEnd;

} HookInfo;

// MODULE_ENTRY contains basic information about a module
typedef struct _MODULE_ENTRY
{
    UNICODE_STRING BaseName; // BaseName of the module
    UNICODE_STRING FullName; // FullName of the module
    ULONG SizeOfImage;       // Size in bytes of the module
    PVOID BaseAddress;       // Base address of the module
    PVOID EntryPoint;        // Entrypoint of the module
} MODULE_ENTRY, *PMODULE_ENTRY;

// MODULE_INFORMATION_TABLE contains basic information about all the modules of a given process
typedef struct _MODULE_INFORMATION_TABLE
{
    ULONG Pid;             // PID of the process
    ULONG ModuleCount;     // Modules count for the above pointer
    PMODULE_ENTRY Modules; // Pointer to 0...* modules
} MODULE_INFORMATION_TABLE, *PMODULE_INFORMATION_TABLE;
