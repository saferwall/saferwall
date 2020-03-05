#include "stdafx.h"
//#include "ole.h"
//
//decltype(CoCreateInstanceEx) *TrueCoCreateInstanceEx = nullptr;
//
//HRESULT
//HookCoCreateInstanceEx(
//    _In_ REFCLSID Clsid,
//    _In_opt_ IUnknown *punkOuter,
//    _In_ DWORD dwClsCtx,
//    _In_opt_ COSERVERINFO *pServerInfo,
//    _In_ DWORD dwCount,
//    _Inout_updates_(dwCount) MULTI_QI *pResults)
//{
//    WCHAR szGuidW[40] = {0};
//
//    if (IsInsideHook())
//    {
//        goto end;
//    }
//
//    CaptureStackTrace();
//
//	StringFromGUID2(Clsid, szGuidW, 40);
//    TraceAPI(L"CoCreateInstanceEx(szGuidW: %ws), RETN: 0x%p", _ReturnAddress());
//
//    ReleaseHookGuard();
//end:
//    return TrueCoCreateInstanceEx(
//        Clsid, punkOuter, dwClsCtx, pServerInfo, dwCount, pResults);
//}
//
