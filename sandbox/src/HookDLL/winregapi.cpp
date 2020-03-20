#include "stdafx.h"
#include "synchapi.h"

decltype(NtOpenKey) *TrueNtOpenKey = nullptr;
decltype(NtOpenKeyEx) *TrueNtOpenKeyEx = nullptr;
decltype(NtCreateKey) *TrueNtCreateKey = nullptr;
decltype(NtQueryValueKey) *TrueNtQueryValueKey = nullptr;
decltype(NtDeleteKey) *TrueNtDeleteKey = nullptr;
decltype(NtDeleteValueKey) *TrueNtDeleteValueKey = nullptr;
decltype(NtSetValueKey) *TrueNtSetValueKey = nullptr;

NTSTATUS NTAPI
HookNtOpenKey(_Out_ PHANDLE KeyHandle, _In_ ACCESS_MASK DesiredAccess, _In_ POBJECT_ATTRIBUTES ObjectAttributes)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtOpenKey(KeyHandle, DesiredAccess, ObjectAttributes);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtOpenKey(DesiredAccess: 0x%x, ObjectName:%ws, ReturnLength:0x%p), _ReturnAddress: 0x%p",
        DesiredAccess,
        ObjectAttributes->ObjectName->Buffer,
        _ReturnAddress());

    NTSTATUS Status = TrueNtOpenKey(KeyHandle, DesiredAccess, ObjectAttributes);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtOpenKeyEx(
    _Out_ PHANDLE KeyHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_ POBJECT_ATTRIBUTES ObjectAttributes,
    _In_ ULONG OpenOptions)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtOpenKeyEx(KeyHandle, DesiredAccess, ObjectAttributes, OpenOptions);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtOpenKeyEx(DesiredAccess: 0x%d, ObjectName:%ws, ReturnLength:0x%p), RETN: 0x%p",
        DesiredAccess,
        ObjectAttributes->ObjectName->Buffer,
        _ReturnAddress());

    NTSTATUS Status = TrueNtOpenKeyEx(KeyHandle, DesiredAccess, ObjectAttributes, OpenOptions);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtCreateKey(
    _Out_ PHANDLE KeyHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_ POBJECT_ATTRIBUTES ObjectAttributes,
    _Reserved_ ULONG TitleIndex,
    _In_opt_ PUNICODE_STRING Class,
    _In_ ULONG CreateOptions,
    _Out_opt_ PULONG Disposition)
{
    if (SfwIsCalledFromSystemMemory(4))
    {
        return TrueNtCreateKey(
            KeyHandle, DesiredAccess, ObjectAttributes, TitleIndex, Class, CreateOptions, Disposition);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtCreateKey(DesiredAccess: 0x%d, ObjectName:%ws, CreateOptions: %ul, ReturnLength:0x%p), RETN: 0x%p",
        DesiredAccess,
        ObjectAttributes->ObjectName->Buffer,
        CreateOptions,
        _ReturnAddress());

    NTSTATUS Status =
        TrueNtCreateKey(KeyHandle, DesiredAccess, ObjectAttributes, TitleIndex, Class, CreateOptions, Disposition);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtQueryValueKey(
    _In_ HANDLE KeyHandle,
    _In_ PUNICODE_STRING ValueName,
    _In_ KEY_VALUE_INFORMATION_CLASS KeyValueInformationClass,
    _Out_writes_bytes_opt_(Length) PVOID KeyValueInformation,
    _In_ ULONG Length,
    _Out_ PULONG ResultLength)
/*
    RegQueryValueA -> RegQueryValueExA -> RegQueryValueExW -> NtQueryValueKey
*/
{
    if (SfwIsCalledFromSystemMemory(4))
    {
        return TrueNtQueryValueKey(
            KeyHandle, ValueName, KeyValueInformationClass, KeyValueInformation, Length, ResultLength);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryValueKey(KeyHandle: 0x%d, ValueName:%ws), RETN: %p", KeyHandle, ValueName->Buffer, _ReturnAddress());

    NTSTATUS Status =
        TrueNtQueryValueKey(KeyHandle, ValueName, KeyValueInformationClass, KeyValueInformation, Length, ResultLength);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtDeleteKey(_In_ HANDLE KeyHandle)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtDeleteKey(KeyHandle);
    }

    CaptureStackTrace();

    TraceAPI(L"NtDeleteKey(KeyHandle: 0x%d), RETN: %p", KeyHandle, _ReturnAddress());

    NTSTATUS Status = TrueNtDeleteKey(KeyHandle);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtDeleteValueKey(_In_ HANDLE KeyHandle, _In_ PUNICODE_STRING ValueName)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtDeleteValueKey(KeyHandle, ValueName);
    }

    CaptureStackTrace();

    TraceAPI(L"NtDeleteValueKey(KeyHandle: 0x%d, ValueName: %ws), RETN: %p", KeyHandle, ValueName, _ReturnAddress());

    NTSTATUS Status = TrueNtDeleteValueKey(KeyHandle, ValueName);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtSetValueKey(
    _In_ HANDLE KeyHandle,
    _In_ PUNICODE_STRING ValueName,
    _In_opt_ ULONG TitleIndex,
    _In_ ULONG Type,
    _In_reads_bytes_opt_(DataSize) PVOID Data,
    _In_ ULONG DataSize)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtSetValueKey(KeyHandle, ValueName, TitleIndex, Type, Data, DataSize);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtSetValueKey(KeyHandle: 0x%d, ValueName:%ws), RETN: %p", KeyHandle, ValueName->Buffer, _ReturnAddress());

    NTSTATUS Status = TrueNtSetValueKey(KeyHandle, ValueName, TitleIndex, Type, Data, DataSize);

    ReleaseHookGuard();

    return Status;
}