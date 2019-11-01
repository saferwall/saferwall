#pragma once

#include <ntddk.h>


//////////////////////////////////////////////////////////////////////////
// Definitions.
//////////////////////////////////////////////////////////////////////////

//
// _ASSERT
//
// This macro is identical to NT_ASSERT but works in fre builds as well.
//
// It is used for error checking in the driver in cases where
// we can't easily report the error to the user mode app, or the
// error is so severe that we should break in immediately to
// investigate.
//
// It's better than DbgBreakPoint because it provides additional info
// that can be dumped with .exr -1, and individual asserts can be disabled
// from kd using 'ahi' command.
//

#define _ASSERT(_exp) \
    ((!(_exp)) ? \
        (__annotation(L"Debug", L"AssertFail", L#_exp), \
         DbgRaiseAssertionFailure(), FALSE) : \
        TRUE)




#define _LogMsg(lvl, lvlname, format, ...)  \
		DbgPrintEx(DPFLTR_IHVDRIVER_ID, lvl , "Saferwall" __FUNCTION__ "[irql:%d,pid:%d][" lvlname "]: " format "\n", KeGetCurrentIrql(), PsGetCurrentProcessId(), __VA_ARGS__)

#define LOG_ERROR(format, ...)	_LogMsg(DPFLTR_ERROR_LEVEL,   "error",   format, __VA_ARGS__)
#define LOG_WARN(format, ...)	_LogMsg(DPFLTR_WARNING_LEVEL, "warning", format, __VA_ARGS__)
#define LOG_TRACE(format, ...)	_LogMsg(DPFLTR_TRACE_LEVEL,   "trace",   format, __VA_ARGS__)
#define LOG_INFO(format, ...)	_LogMsg(DPFLTR_INFO_LEVEL,    "info",    format, __VA_ARGS__)

