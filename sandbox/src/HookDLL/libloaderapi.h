#pragma once

#include "stdafx.h"

//
// Prototypes
//

NTSTATUS
NTAPI
HookLdrLoadDll(PWSTR DllPath, PULONG DllCharacteristics, PUNICODE_STRING DllName, PVOID *DllHandle);

NTSTATUS
NTAPI
HookLdrGetProcedureAddressEx(
    PVOID DllHandle,
    PANSI_STRING ProcedureName,
    ULONG ProcedureNumber,
    PVOID *ProcedureAddress,
    ULONG Flags);

NTSTATUS
NTAPI
HookLdrGetDllHandleEx(
    _In_ ULONG Flags,
    _In_opt_ PWSTR DllPath,
    _In_opt_ PULONG DllCharacteristics,
    _In_ PUNICODE_STRING DllName,
    _Out_opt_ PVOID *DllHandle);