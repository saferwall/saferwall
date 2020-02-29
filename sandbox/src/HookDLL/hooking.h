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


//
// Prototypes
//

BOOL ProcessAttach();
BOOL ProcessDetach();
BOOL IsInsideHook();
VOID ReleaseHookGuard();
VOID CaptureStackTrace();


//
// Unfortunatelly sprintf-like functions are not exposed
// by ntdll.lib, which we're linking against.  We have to
// load them dynamically.
//

using __vsnwprintf_fn_t = int(__cdecl*)(
	wchar_t *buffer,
	size_t count,
	const wchar_t *format,
	...
	);

using __snwprintf_fn_t = int(__cdecl*)(
	wchar_t *buffer,
	size_t count,
	const wchar_t *format,
	...
	);

using strlen_fn_t = size_t(__cdecl*)(
	char const *buffer
	);
