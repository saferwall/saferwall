#include "stdafx.h"
#include "libloaderapi.h"

decltype(LdrLoadDll) *TrueLdrLoadDll = nullptr;
decltype(LdrGetProcedureAddressEx) *TrueLdrGetProcedureAddressEx = nullptr;
decltype(LdrGetDllHandleEx) *TrueLdrGetDllHandleEx = nullptr;

extern pfn_wcsstr _wcsstr;
extern HookInfo gHookInfo;

NTSTATUS
NTAPI
HookLdrLoadDll(PWSTR DllPath, PULONG DllCharacteristics, PUNICODE_STRING DllName, PVOID *DllHandle)
/*
- LdrLoadDll
    - LoadLibraryA -> LoadLibraryExA
    - LoadLibraryW -> LoadLibraryExW
*/
{
    if (IsInsideHook() || SfwIsCalledFromSystemMemory(4))
    {
        return TrueLdrLoadDll(DllPath, DllCharacteristics, DllName, DllHandle);
    }

    CaptureStackTrace();

    if (DllName && DllName->Buffer)
    {
        TraceAPI(L"LdrLoadDll(%ws), RETN: 0x%p", DllName->Buffer, _ReturnAddress());
    }

    NTSTATUS Status = TrueLdrLoadDll(DllPath, DllCharacteristics, DllName, DllHandle);

    if (NT_SUCCESS(Status))
    {
        if (DllName && DllName->Buffer)
        {
            HookDll(DllName->Buffer);
        }
    }

    ReleaseHookGuard();

    return Status;
}

NTSTATUS
NTAPI
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
        return TrueLdrGetProcedureAddressEx(DllHandle, ProcedureName, ProcedureNumber, ProcedureAddress, Flags);
    }

    CaptureStackTrace();

    if (ProcedureName && ProcedureName->Buffer)
        TraceAPI(L"LdrGetProcedureAddressEx(%ws) RETN: 0x%p", MultiByteToWide(ProcedureName->Buffer), _ReturnAddress());
    else
        TraceAPI(L"LdrGetProcedureAddressEx(Ordinal:0x%x), RETN: 0x%p", ProcedureNumber, _ReturnAddress());

    NTSTATUS Status = TrueLdrGetProcedureAddressEx(DllHandle, ProcedureName, ProcedureNumber, ProcedureAddress, Flags);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS
NTAPI
HookLdrGetDllHandleEx(
    _In_ ULONG Flags,
    _In_opt_ PWSTR DllPath,
    _In_opt_ PULONG DllCharacteristics,
    _In_ PUNICODE_STRING DllName,
    _Out_opt_ PVOID *DllHandle)
/*
- LdrGetDllHandle -> LdrGetDllHandleEx
*/
{
    if (IsInsideHook() || SfwIsCalledFromSystemMemory(3))
    {
        return TrueLdrGetDllHandleEx(Flags, DllPath, DllCharacteristics, DllName, DllHandle);
    }

    CaptureStackTrace();

    if (DllName && DllName->Buffer)
    {
        TraceAPI(L"LdrGetDllHandleEx(%ws), RETN: 0x%p", DllName->Buffer, _ReturnAddress());
    }

    NTSTATUS Status = TrueLdrGetDllHandleEx(Flags, DllPath, DllCharacteristics, DllName, DllHandle);

    ReleaseHookGuard();

    return Status;
}