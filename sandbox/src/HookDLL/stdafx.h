// stdafx.h : include file for standard system include files,
// or project specific include files that are used frequently, but
// are changed infrequently
//

#pragma once

#include "targetver.h"

#define WIN32_LEAN_AND_MEAN             // Exclude rarely-used stuff from Windows headers



// reference additional headers your program requires here
#include "ntdll.h"

//
// For program instrumentation.
//
#include <detours.h>


//
// Include support for ETW logging.
// Note that following functions are mocked, because they're
// located in advapi32.dll.  Fortunatelly, advapi32.dll simply
// redirects calls to these functions to the ntdll.dll.
//

#define EventActivityIdControl  EtwEventActivityIdControl
#define EventEnabled            EtwEventEnabled
#define EventProviderEnabled    EtwEventProviderEnabled
#define EventRegister           EtwEventRegister
#define EventSetInformation     EtwEventSetInformation
#define EventUnregister         EtwEventUnregister
#define EventWrite              EtwEventWrite
#define EventWriteEndScenario   EtwEventWriteEndScenario
#define EventWriteEx            EtwEventWriteEx
#define EventWriteStartScenario EtwEventWriteStartScenario
#define EventWriteString        EtwEventWriteString
#define EventWriteTransfer      EtwEventWriteTransfer

#include <evntprov.h>


//
// For Stack walking
//

#ifdef _DEBUG
#include <dbghelp.h>
#endif // DEBUG








#include "hooking.h"
#include "logging.h"

//#include <shlwapi.h>
//#include <psapi.h>
#include <intrin.h>
//#include <stdio.h>
#include <stdlib.h>

//
//
//
//#pragma comment (lib, "Shlwapi.lib")


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

//extern "C" size_t __cdecl wcslen(wchar_t const* const string)
//{
//    wchar_t const* string_it = string;
//    while (*string_it++) { }
//    return static_cast<size_t>(string_it - string - 1);
//}

