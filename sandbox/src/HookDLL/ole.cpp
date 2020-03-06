#include "stdafx.h"
#include "ole.h"

extern pfnCoCreateInstanceEx TrueCoCreateInstanceEx = nullptr;
extern pfnStringFromGUID2 _StringFromGUID2;
extern pfnStringFromCLSID _StringFromCLSID;
extern pfnCoTaskMemFree _CoTaskMemFree;


HRESULT
HookCoCreateInstanceEx(
    _In_ REFCLSID Clsid,
    _In_opt_ IUnknown *punkOuter,
    _In_ DWORD dwClsCtx,
    _In_opt_ COSERVERINFO *pServerInfo,
    _In_ DWORD dwCount,
    _Inout_updates_(dwCount) MULTI_QI *pResults)
{
   /* if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();*/

	//_StringFromGUID2(Clsid, szGuidW, 100);
//    OLECHAR *guidString;
//    _StringFromCLSID(Clsid, &guidString);
//    TraceAPI(L"CoCreateInstanceEx(szGuidW: %ws), RETN: 0x%p", guidString, _ReturnAddress());
//    _CoTaskMemFree(guidString);
//
//    ReleaseHookGuard();
//end:
    return TrueCoCreateInstanceEx(
        Clsid, punkOuter, dwClsCtx, pServerInfo, dwCount, pResults);
}


