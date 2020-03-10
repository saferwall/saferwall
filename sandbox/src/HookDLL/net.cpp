#include "stdafx.h"
#include "net.h"

extern pfnInternetOpenA TrueInternetOpenA;

HINTERNET
HookInternetOpenA(
    _In_opt_ LPCSTR lpszAgent,
    _In_ DWORD dwAccessType,
    _In_opt_ LPCSTR lpszProxy,
    _In_opt_ LPCSTR lpszProxyBypass,
    _In_ DWORD dwFlags)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"InternetOpenA(zAgent: %s), RETN: 0x%p", lpszAgent, _ReturnAddress());

    ReleaseHookGuard();
end:
    return TrueInternetOpenA(lpszAgent, dwAccessType, lpszProxy, lpszProxyBypass, dwFlags);
}
