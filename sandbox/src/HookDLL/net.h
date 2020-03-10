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

//////////////////////////////////////////////////////////////

HINTERNET
HookInternetOpenA(
    _In_opt_ LPCSTR lpszAgent,
    _In_ DWORD dwAccessType,
    _In_opt_ LPCSTR lpszProxy,
    _In_opt_ LPCSTR lpszProxyBypass,
    _In_ DWORD dwFlags);

HINTERNET
HookInternetConnectA(
    _In_ HINTERNET hInternet,
    _In_ LPCSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCSTR lpszUserName,
    _In_opt_ LPCSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);

HINTERNET
HookInternetConnectW(
    _In_ HINTERNET hInternet,
    _In_ LPCWSTR lpszServerName,
    _In_ INTERNET_PORT nServerPort,
    _In_opt_ LPCWSTR lpszUserName,
    _In_opt_ LPCWSTR lpszPassword,
    _In_ DWORD dwService,
    _In_ DWORD dwFlags,
    _In_opt_ DWORD_PTR dwContext);