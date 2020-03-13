#pragma once

#include "stdafx.h"
#include <WinInet.h>

//
// Prototypes
//

using pfnInternetOpenA = HINTERNET(__stdcall *)(
    _In_opt_ LPCSTR lpszAgent,
    _In_ DWORD dwAccessType,
    _In_opt_ LPCSTR lpszProxy,
    _In_opt_ LPCSTR lpszProxyBypass,
    _In_ DWORD dwFlags);

using pfnInternetConnectA = HINTERNET(__stdcall *)(
    _In_ HINTERNET hInternet,
    _In_ LPCSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCSTR lpszUserName,
    _In_opt_ LPCSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);

using pfnInternetConnectW = HINTERNET(__stdcall *)(
    _In_ HINTERNET hInternet,
    _In_ LPCWSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCWSTR lpszUserName,
    _In_opt_ LPCWSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);

using pfnHttpOpenRequestA = HINTERNET(__stdcall *)(
    _In_ HINTERNET hConnect,
    _In_opt_ LPCSTR lpszVerb,
    _In_opt_ LPCSTR lpszObjectName,
    _In_opt_ LPCSTR lpszVersion,
    _In_opt_ LPCSTR lpszReferrer,
    _In_opt_z_ LPCSTR FAR *lplpszAcceptTypes,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);

using pfnHttpOpenRequestW = HINTERNET(__stdcall *)(
    _In_ HINTERNET hConnect,
    _In_opt_ LPCWSTR lpszVerb,
    _In_opt_ LPCWSTR lpszObjectName,
    _In_opt_ LPCWSTR lpszVersion,
    _In_opt_ LPCWSTR lpszReferrer,
    _In_opt_z_ LPCWSTR FAR *lplpszAcceptTypes,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);

//////////////////////////////////////////////////////////////

HINTERNET WINAPI HookInternetOpenA(
    _In_opt_ LPCSTR lpszAgent,
    _In_ DWORD dwAccessType,
    _In_opt_ LPCSTR lpszProxy,
    _In_opt_ LPCSTR lpszProxyBypass,
    _In_ DWORD dwFlags);

HINTERNET WINAPI
HookInternetConnectA(
    _In_ HINTERNET hInternet,
    _In_ LPCSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCSTR lpszUserName,
    _In_opt_ LPCSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);

HINTERNET WINAPI
HookInternetConnectW(
    _In_ HINTERNET hInternet,
    _In_ LPCWSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCWSTR lpszUserName,
    _In_opt_ LPCWSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);


HINTERNET __stdcall HookHttpOpenRequestA(
    _In_ HINTERNET hConnect,
    _In_opt_ LPCSTR lpszVerb,
    _In_opt_ LPCSTR lpszObjectName,
    _In_opt_ LPCSTR lpszVersion,
    _In_opt_ LPCSTR lpszReferrer,
    _In_opt_z_ LPCSTR FAR *lplpszAcceptTypes,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);

HINTERNET __stdcall HookHttpOpenRequestW(
    _In_ HINTERNET hConnect,
    _In_opt_ LPCWSTR lpszVerb,
    _In_opt_ LPCWSTR lpszObjectName,
    _In_opt_ LPCWSTR lpszVersion,
    _In_opt_ LPCWSTR lpszReferrer,
    _In_opt_z_ LPCWSTR FAR *lplpszAcceptTypes,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);