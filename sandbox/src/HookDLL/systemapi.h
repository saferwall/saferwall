#pragma once

#include "stdafx.h"

//
// Prototypes
//

NTSTATUS
NTAPI
HookNtQuerySystemInformation(
    _In_ SYSTEM_INFORMATION_CLASS SystemInformationClass,
    _Out_writes_bytes_opt_(SystemInformationLength) PVOID SystemInformation,
    _In_ ULONG SystemInformationLength,
    _Out_opt_ PULONG ReturnLength);

NTSTATUS NTAPI
HookNtQueryVolumeInformationFile(
    _In_ HANDLE FileHandle,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _Out_writes_bytes_(Length) PVOID FsInformation,
    _In_ ULONG Length,
    _In_ FSINFOCLASS FsInformationClass);

NTSTATUS
NTAPI
HookNtLoadDriver(_In_ PUNICODE_STRING DriverServiceName);
