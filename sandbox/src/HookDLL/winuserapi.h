#pragma once

#include "stdafx.h"

//
// Prototypes.
//

using pfnSetWindowsHookExW =
    HHOOK(WINAPI *)(_In_ int idHook, _In_ HOOKPROC lpfn, _In_opt_ HINSTANCE hmod, _In_ DWORD dwThreadId);


HHOOK
WINAPI
HookSetWindowsHookExW(_In_ int idHook, _In_ HOOKPROC lpfn, _In_opt_ HINSTANCE hmod, _In_ DWORD dwThreadId);