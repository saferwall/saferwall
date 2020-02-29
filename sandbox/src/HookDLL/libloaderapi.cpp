#include "stdafx.h"
#include "libloaderapi.h"

decltype(LdrLoadDll) *TrueLdrLoadDll = nullptr;
decltype(LdrGetProcedureAddressEx) *TrueLdrGetProcedureAddressEx = nullptr;

NTSTATUS
WINAPI
HookLdrLoadDll(PWSTR DllPath, PULONG DllCharacteristics, PUNICODE_STRING DllName, PVOID *DllHandle)
/*
- LdrLoadDll
    - LoadLibraryA -> LoadLibraryExA
    - LoadLibraryW -> LoadLibraryExW
*/
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    if (DllName && DllName->Buffer)
    {
        TraceAPI(L"LdrLoadDll(%ws), RETN: 0x%p", DllName->Buffer, _ReturnAddress());
    }

    ReleaseHookGuard();
end:
    return TrueLdrLoadDll(DllPath, DllCharacteristics, DllName, DllHandle);
}

NTSTATUS
WINAPI
HookLdrGetProcedureAddressEx(
    PVOID DllHandle,
    PANSI_STRING ProcedureName,
    ULONG ProcedureNumber,
    PVOID *ProcedureAddress,
    ULONG Flags)
/*
- LdrGetProcedureAddressEx
    - GetProcAddress
    - LdrGetProcedureAddress
 */
{
    if (IsInsideHook())
    {
        goto end;
    }
    CaptureStackTrace();

    if (ProcedureName && ProcedureName->Buffer)
        TraceAPI(L"LdrGetProcedureAddressEx(%ws) RETN: 0x%p", MultiByteToWide(ProcedureName->Buffer), _ReturnAddress());
    else
        TraceAPI(L"LdrGetProcedureAddressEx(Ordinal:0x%x), RETN: 0x%p", ProcedureNumber, _ReturnAddress());

    ReleaseHookGuard();
end:
    return TrueLdrGetProcedureAddressEx(DllHandle, ProcedureName, ProcedureNumber, ProcedureAddress, Flags);
}
