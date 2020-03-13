#include "stdafx.h"
#include "net.h"

extern pfnInternetOpenA TrueInternetOpenA;
extern pfnInternetConnectA TrueInternetConnectA;
extern pfnInternetConnectW TrueInternetConnectW;
extern pfnHttpOpenRequestA TrueHttpOpenRequestA;
extern pfnHttpOpenRequestW TrueHttpOpenRequestW;
extern pfnHttpSendRequestA TrueHttpSendRequestA;
extern pfnHttpSendRequestW TrueHttpSendRequestW;

HINTERNET WINAPI
HookInternetOpenA(
    _In_opt_ LPCSTR lpszAgent,
    _In_ DWORD dwAccessType,
    _In_opt_ LPCSTR lpszProxy,
    _In_opt_ LPCSTR lpszProxyBypass,
    _In_ DWORD dwFlags)
{
    /*
        InternetOpenW -> InternetOpenA.
    */
    if (IsInsideHook())
    {
        return TrueInternetOpenA(lpszAgent, dwAccessType, lpszProxy, lpszProxyBypass, dwFlags);
    }

    CaptureStackTrace();

    TraceAPI(L"InternetOpenA(Agent: %s), RETN: 0x%p", MultiByteToWide((CHAR *)lpszAgent), _ReturnAddress());

    HINTERNET hSession = TrueInternetOpenA(lpszAgent, dwAccessType, lpszProxy, lpszProxyBypass, dwFlags);

    ReleaseHookGuard();

    return hSession;
}

HINTERNET WINAPI
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
        return TrueInternetConnectA(
            hInternet, lpszServerName, nServerPort, lpszUserName, lpszPassword, dwService, dwFlags, dwContext);
    }

    CaptureStackTrace();

    TraceAPI(
        L"InternetConnectA(ServerName: %s), RETN: 0x%p", MultiByteToWide((CHAR *)lpszServerName), _ReturnAddress());

    HINTERNET hConnect = TrueInternetConnectA(
        hInternet, lpszServerName, nServerPort, lpszUserName, lpszPassword, dwService, dwFlags, dwContext);

    ReleaseHookGuard();

    return hConnect;
}

HINTERNET WINAPI
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
        return TrueInternetConnectW(
            hInternet, lpszServerName, nServerPort, lpszUserName, lpszPassword, dwService, dwFlags, dwContext);
    }

    CaptureStackTrace();

    TraceAPI(L"InternetConnectW(ServerName: %ws), RETN: 0x%p", lpszServerName, _ReturnAddress());

    HINTERNET hConnect = TrueInternetConnectW(
        hInternet, lpszServerName, nServerPort, lpszUserName, lpszPassword, dwService, dwFlags, dwContext);

    ReleaseHookGuard();

    return hConnect;
}

HINTERNET WINAPI
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
        return TrueHttpOpenRequestA(
            hConnect, lpszVerb, lpszObjectName, lpszVersion, lpszReferrer, lplpszAcceptTypes, dwFlags, dwContext);
    }

    CaptureStackTrace();

    TraceAPI(L"HttpOpenRequestA(Method: %s), RETN: 0x%p", MultiByteToWide((CHAR *)lpszVerb), _ReturnAddress());

    HINTERNET hRequest = TrueHttpOpenRequestA(
        hConnect, lpszVerb, lpszObjectName, lpszVersion, lpszReferrer, lplpszAcceptTypes, dwFlags, dwContext);

    ReleaseHookGuard();

    return hRequest;
}

HINTERNET WINAPI
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
        return TrueHttpOpenRequestW(
            hConnect, lpszVerb, lpszObjectName, lpszVersion, lpszReferrer, lplpszAcceptTypes, dwFlags, dwContext);
    }

    CaptureStackTrace();

    TraceAPI(L"HttpOpenRequestW(Method: %ws), RETN: 0x%p", lpszVerb, _ReturnAddress());

    HINTERNET hRequest = TrueHttpOpenRequestW(
        hConnect, lpszVerb, lpszObjectName, lpszVersion, lpszReferrer, lplpszAcceptTypes, dwFlags, dwContext);

    ReleaseHookGuard();

    return hRequest;
}

BOOL WINAPI
HookHttpSendRequestA(
    _In_ HINTERNET hRequest,
    _In_reads_opt_(dwHeadersLength) LPCSTR lpszHeaders,
    _In_ DWORD dwHeadersLength,
    _In_reads_bytes_opt_(dwOptionalLength) LPVOID lpOptional,
    _In_ DWORD dwOptionalLength)
{
    if (IsInsideHook())
    {
        return TrueHttpSendRequestA(hRequest, lpszHeaders, dwHeadersLength, lpOptional, dwOptionalLength);
    }

    CaptureStackTrace();

    TraceAPI(L"HttpSendRequestA(hRequest: %p), RETN: 0x%p", hRequest, _ReturnAddress());

    BOOL bSend = TrueHttpSendRequestA(hRequest, lpszHeaders, dwHeadersLength, lpOptional, dwOptionalLength);

    ReleaseHookGuard();

    return bSend;
}

BOOL WINAPI
HookHttpSendRequestW(
    _In_ HINTERNET hRequest,
    _In_reads_opt_(dwHeadersLength) LPCWSTR lpszHeaders,
    _In_ DWORD dwHeadersLength,
    _In_reads_bytes_opt_(dwOptionalLength) LPVOID lpOptional,
    _In_ DWORD dwOptionalLength)
{
    if (IsInsideHook())
    {
        return TrueHttpSendRequestW(hRequest, lpszHeaders, dwHeadersLength, lpOptional, dwOptionalLength);
    }

    CaptureStackTrace();

    TraceAPI(L"HttpSendRequestW(hRequest: %p), RETN: 0x%p", hRequest, _ReturnAddress());

    BOOL bSend = TrueHttpSendRequestW(hRequest, lpszHeaders, dwHeadersLength, lpOptional, dwOptionalLength);

    ReleaseHookGuard();

    return bSend;
}