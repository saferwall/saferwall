#include "stdafx.h"
#include "systemapi.h"

decltype(NtQuerySystemInformation) *TrueNtQuerySystemInformation = nullptr;
decltype(NtLoadDriver) *TrueNtLoadDriver = nullptr;

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
HookNtLoadDriver(
	_In_ PUNICODE_STRING DriverServiceName)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

	if (DriverServiceName && DriverServiceName->Buffer)
    TraceAPI(
        L"NtLoadDriver(DriverServiceName: %ws), RETN: %p",

        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtLoadDriver(DriverServiceName);
}
