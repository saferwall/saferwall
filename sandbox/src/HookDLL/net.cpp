#include "stdafx.h"
#include "net.h"

extern pfnInternetOpenA TrueInternetOpenA;
extern pfnInternetConnectA TrueInternetConnectA;
extern pfnInternetConnectW TrueInternetConnectW;
extern pfnHttpOpenRequestA TrueHttpOpenRequestA;
extern pfnHttpOpenRequestW TrueHttpOpenRequestW;

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

    TraceAPI(L"InternetConnectA(ServerName: %s), RETN: 0x%p", lpszServerName, _ReturnAddress());

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

    TraceAPI(L"InternetConnectW(ServerName: %ws), RETN: 0x%p", lpszServerName, _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueInternetConnectW(
        hInternet, lpszServerName, nServerPort, lpszUserName, lpszPassword, dwService, dwFlags, dwContext);
}


HINTERNET
HookHttpOpenRequestA(
    _In_ HINTERNET hConnect,
    _In_opt_ LPCSTR lpszVerb,
    _In_opt_ LPCSTR lpszObjectName,
    _In_opt_ LPCSTR lpszVersion,
    _In_opt_ LPCSTR lpszReferrer,
    _In_opt_z_ LPCSTR FAR *lplpszAcceptTypes,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"HttpOpenRequestA(Method: %s), RETN: 0x%p", lpszVerb, _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueHttpOpenRequestA(hConnect, lpszVerb, lpszObjectName, lpszVersion, lpszReferrer, lplpszAcceptTypes, dwFlags, dwContext);
}


HINTERNET
HookHttpOpenRequestW(
    _In_ HINTERNET hConnect,
    _In_opt_ LPCWSTR lpszVerb,
    _In_opt_ LPCWSTR lpszObjectName,
    _In_opt_ LPCWSTR lpszVersion,
    _In_opt_ LPCWSTR lpszReferrer,
    _In_opt_z_ LPCWSTR FAR *lplpszAcceptTypes,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"HttpOpenRequestW(Method: %ws), RETN: 0x%p", lpszVerb, _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueHttpOpenRequestW(
        hConnect, lpszVerb, lpszObjectName, lpszVersion, lpszReferrer, lplpszAcceptTypes, dwFlags, dwContext);
}