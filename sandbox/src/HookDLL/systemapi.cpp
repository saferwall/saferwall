#include "stdafx.h"
#include "systemapi.h"

decltype(NtQuerySystemInformation) *TrueNtQuerySystemInformation = nullptr;
decltype(NtSetInformationFile) *TrueNtSetInformationFile = nullptr;

NTSTATUS WINAPI
HookNtQuerySystemInformation(
    _In_ SYSTEM_INFORMATION_CLASS SystemInformationClass,
    _Out_writes_bytes_opt_(SystemInformationLength) PVOID SystemInformation,
    _In_ ULONG SystemInformationLength,
    _Out_opt_ PULONG ReturnLength)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQuerySystemInformation(SystemInformationClass: %d, SystemInformation:0x%p, SystemInformationLength:0x%08x), RETN: %p",
        SystemInformationClass,
        SystemInformation,
        SystemInformationLength,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtQuerySystemInformation(
        SystemInformationClass, SystemInformation, SystemInformationLength, ReturnLength);
}

NTSTATUS WINAPI
HookNtSetInformationFile(
    _In_ HANDLE FileHandle,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _In_reads_bytes_(Length) PVOID FileInformation,
    _In_ ULONG Length,
    _In_ FILE_INFORMATION_CLASS FileInformationClass)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtSetInformationFile(FileInformationClass: %d, FileInformation:0x%p, Length:0x%08x, ReturnLength:0x%p), RETN: %p",
        FileInformationClass,
        FileInformation,
        Length,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtSetInformationFile(FileHandle, IoStatusBlock, FileInformation, Length, FileInformationClass);
}
