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
ReleaseHookGuard();
VOID
CaptureStackTrace();
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
} HookInfo;