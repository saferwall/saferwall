#ifndef _NTDLL_H
#define _NTDLL_H

//
// Remap definitions.
//

#ifdef NTDLL_NO_INLINE_INIT_STRING
#define PHNT_NO_INLINE_INIT_STRING
#endif

//
// Hack, because prototype in PH's headers and evntprov.h
// don't match.
//

#define EtwEventRegister __EtwEventRegisterIgnored

//
// Provides access to the Win32 API as well as the `NTSTATUS` values. 
//
#include <phnt_windows.h>

//
// Provides access to the entire Native API.
//

#include <phnt.h>

#undef  EtwEventRegister

#endif
