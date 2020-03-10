#include "stdafx.h"
#include "net.h"

extern pfnInternetOpenA TrueInternetOpenA;
extern pfnInternetConnectA TrueInternetConnectA;
extern pfnInternetConnectW TrueInternetConnectW;

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

HINTERNET
HookInternetConnectA(
    _In_ HINTERNET hInternet,
    _In_ LPCSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCSTR lpszUserName,
    _In_opt_ LPCSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"InternetConnectA(szServerName: %s), RETN: 0x%p", lpszServerName, _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueInternetConnectA(
        hInternet, lpszServerName, nServerPort, lpszUserName, lpszPassword, dwService, dwFlags, dwContext);
}


HINTERNET
HookInternetConnectW(
    _In_ HINTERNET hInternet,
    _In_ LPCWSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCWSTR lpszUserName,
    _In_opt_ LPCWSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"InternetConnectW(szServerName: %ws), RETN: 0x%p", lpszServerName, _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueInternetConnectW(
        hInternet, lpszServerName, nServerPort, lpszUserName, lpszPassword, dwService, dwFlags, dwContext);
}