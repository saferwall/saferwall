#pragma once
#include "stdafx.h"

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
// Defines
//
#define NtCurrentThread()         ((HANDLE)(LONG_PTR)-2)


//
// Prototypes
//

VOID SetupHook();
VOID Unhook();
VOID GetStackWalk();
BOOL IsInsideHook();
VOID ReleaseHookGuard();
PWCHAR MultiByteToWide(PCHAR lpMultiByteStr);
LPCWSTR FindFileName(LPCWSTR pPath);
VOID PrintStackTrace();
VOID CaptureStackTrace();