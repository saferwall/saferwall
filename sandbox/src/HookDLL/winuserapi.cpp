#include "stdafx.h"
#include "winuserapi.h"

extern pfnSetWindowsHookExW TrueSetWindowsHookExW;

HHOOK
WINAPI
HookSetWindowsHookExW(_In_ int idHook, _In_ HOOKPROC lpfn, _In_opt_ HINSTANCE hmod, _In_ DWORD dwThreadId)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueSetWindowsHookExW(idHook, lpfn, hmod, dwThreadId);
    }

    CaptureStackTrace();

    TraceAPI(L"SetWindowsHookExW(idHook: %d, lpfn: %p,  ThreadID: %d), RETN: 0x%p", idHook, lpfn, dwThreadId, _ReturnAddress());

    HHOOK hHook = TrueSetWindowsHookExW(idHook, lpfn, hmod, dwThreadId);

    ReleaseHookGuard();

    return hHook;
}
