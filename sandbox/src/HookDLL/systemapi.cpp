#include "stdafx.h"
#include "systemapi.h"

decltype(NtQuerySystemInformation) *TrueNtQuerySystemInformation = nullptr;
decltype(NtQueryVolumeInformationFile) *TrueNtQueryVolumeInformationFile = nullptr;
decltype(NtLoadDriver) *TrueNtLoadDriver = nullptr;

NTSTATUS NTAPI
HookNtQuerySystemInformation(
    _In_ SYSTEM_INFORMATION_CLASS SystemInformationClass,
    _Out_writes_bytes_opt_(SystemInformationLength) PVOID SystemInformation,
    _In_ ULONG SystemInformationLength,
    _Out_opt_ PULONG ReturnLength)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtQuerySystemInformation(
            SystemInformationClass, SystemInformation, SystemInformationLength, ReturnLength);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQuerySystemInformation(SystemInformationClass: %d, SystemInformation:0x%p, SystemInformationLength:0x%08x), RETN: %p",
        SystemInformationClass,
        SystemInformation,
        SystemInformationLength,
        _ReturnAddress());

    NTSTATUS Status =
        TrueNtQuerySystemInformation(SystemInformationClass, SystemInformation, SystemInformationLength, ReturnLength);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtQueryVolumeInformationFile(
    _In_ HANDLE FileHandle,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _Out_writes_bytes_(Length) PVOID FsInformation,
    _In_ ULONG Length,
    _In_ FSINFOCLASS FsInformationClass)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtQueryVolumeInformationFile(FileHandle, IoStatusBlock, FsInformation, Length, FsInformationClass);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryVolumeInformationFile(FileHandle: %p, FsInformationClass:%d), RETN: %p",
        FileHandle,
        FsInformationClass,
        _ReturnAddress());

    NTSTATUS Status =
        TrueNtQueryVolumeInformationFile(FileHandle, IoStatusBlock, FsInformation, Length, FsInformationClass);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtLoadDriver(_In_ PUNICODE_STRING DriverServiceName)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtLoadDriver(DriverServiceName);
    }

    CaptureStackTrace();

    if (DriverServiceName && DriverServiceName->Buffer)
        TraceAPI(
            L"NtLoadDriver(DriverServiceName: %ws), RETN: %p",

            _ReturnAddress());

    NTSTATUS Status = TrueNtLoadDriver(DriverServiceName);

    ReleaseHookGuard();

    return Status;
}
