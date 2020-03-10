#pragma once

#include "stdafx.h"
#include <unknwn.h>

//
// Prototypes
//

HRESULT
HookCoCreateInstanceEx(
    _In_ REFCLSID Clsid,
    _In_opt_ IUnknown *punkOuter,
    _In_ DWORD dwClsCtx,
    _In_opt_ COSERVERINFO *pServerInfo,
    _In_ DWORD dwCount,
    _Inout_updates_(dwCount) MULTI_QI *pResults);


using pfnStringFromGUID2 =
    int(__stdcall *)(_In_ REFGUID rguid, _Out_writes_to_(cchMax, return ) LPOLESTR lpsz, _In_ int cchMax);

using pfnCoCreateInstanceEx = HRESULT(__stdcall *)(
    REFCLSID Clsid,
    IUnknown *punkOuter,
    DWORD dwClsCtx,
    COSERVERINFO *pServerInfo,
    DWORD dwCount,
    MULTI_QI *pResults);

using pfnStringFromCLSID = HRESULT(__stdcall *)(_In_ REFGUID rguid, _Outptr_ LPOLESTR FAR *lplpsz);

using pfnCoTaskMemFree = VOID(__stdcall *)(_Frees_ptr_opt_ LPVOID pv);
