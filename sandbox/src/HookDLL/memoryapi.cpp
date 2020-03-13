#include "stdafx.h"
#include "memoryapi.h"

decltype(NtAllocateVirtualMemory) *TrueNtAllocateVirtualMemory = nullptr;
decltype(NtProtectVirtualMemory) *TrueNtProtectVirtualMemory = nullptr;
decltype(NtQueryVirtualMemory) *TrueNtQueryVirtualMemory = nullptr;
decltype(NtReadVirtualMemory) *TrueNtReadVirtualMemory = nullptr;
decltype(NtWriteVirtualMemory) *TrueNtWriteVirtualMemory = nullptr;
decltype(NtFreeVirtualMemory) *TrueNtFreeVirtualMemory = nullptr;
decltype(NtMapViewOfSection) *TrueNtMapViewOfSection = nullptr;
decltype(NtUnmapViewOfSection) *TrueNtUnmapViewOfSection = nullptr;

NTSTATUS NTAPI
HookNtAllocateVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _Inout_ PVOID *BaseAddress,
    _In_ ULONG_PTR ZeroBits,
    _Inout_ PSIZE_T RegionSize,
    _In_ ULONG AllocationType,
    _In_ ULONG Protect)
{
    if (IsInsideHook())
    {
        return TrueNtAllocateVirtualMemory(ProcessHandle, BaseAddress, ZeroBits, RegionSize, AllocationType, Protect);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtAllocateVirtualMemory(ProcessHandle:0x%p, AllocationType:%lu, Protect:%lu), RETN: 0x%p",
        ProcessHandle,
        AllocationType,
        Protect,
        _ReturnAddress());

    NTSTATUS Status =
        TrueNtAllocateVirtualMemory(ProcessHandle, BaseAddress, ZeroBits, RegionSize, AllocationType, Protect);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtProtectVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _Inout_ PVOID *BaseAddress,
    _Inout_ PSIZE_T RegionSize,
    _In_ ULONG NewProtect,
    _Out_ PULONG OldProtect)
{
    if (IsInsideHook())
    {
        return TrueNtProtectVirtualMemory(ProcessHandle, BaseAddress, RegionSize, NewProtect, OldProtect);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtProtectVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, RegionSize:0x%lu NewProtect: %lu, OldProtect:%lu), RETN: 0x%p",
        ProcessHandle,
        *BaseAddress,
        *RegionSize,
        NewProtect,
        *OldProtect,
        _ReturnAddress());

    NTSTATUS Status = TrueNtProtectVirtualMemory(ProcessHandle, BaseAddress, RegionSize, NewProtect, OldProtect);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtQueryVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _In_opt_ PVOID BaseAddress,
    _In_ MEMORY_INFORMATION_CLASS MemoryInformationClass,
    _Out_writes_bytes_(MemoryInformationLength) PVOID MemoryInformation,
    _In_ SIZE_T MemoryInformationLength,
    _Out_opt_ PSIZE_T ReturnLength)
/*
- VirtualQuery -> VirtualQueryEx -> NtQueryVirtualMemory
*/
{
    if (IsInsideHook())
    {
        return TrueNtQueryVirtualMemory(
            ProcessHandle,
            BaseAddress,
            MemoryInformationClass,
            MemoryInformation,
            MemoryInformationLength,
            ReturnLength);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, MemoryInformationClass:%d, MemoryInformationLength: %lu), RETN: %p",
        ProcessHandle,
        BaseAddress,
        MemoryInformationClass,
        MemoryInformationLength,
        _ReturnAddress());

    NTSTATUS Status = TrueNtQueryVirtualMemory(
        ProcessHandle, BaseAddress, MemoryInformationClass, MemoryInformation, MemoryInformationLength, ReturnLength);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtReadVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _In_opt_ PVOID BaseAddress,
    _Out_writes_bytes_(BufferSize) PVOID Buffer,
    _In_ SIZE_T BufferSize,
    _Out_opt_ PSIZE_T NumberOfBytesRead)
{
    if (IsInsideHook())
    {
        return TrueNtReadVirtualMemory(ProcessHandle, BaseAddress, Buffer, BufferSize, NumberOfBytesRead);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtReadVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, BufferSize:0x%08lu), RETN: %p",
        ProcessHandle,
        BaseAddress,
        BufferSize,
        _ReturnAddress());

    NTSTATUS Status = TrueNtReadVirtualMemory(ProcessHandle, BaseAddress, Buffer, BufferSize, NumberOfBytesRead);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtWriteVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _In_opt_ PVOID BaseAddress,
    _In_reads_bytes_(BufferSize) PVOID Buffer,
    _In_ SIZE_T BufferSize,
    _Out_opt_ PSIZE_T NumberOfBytesWritten)
{
    if (IsInsideHook())
    {
        return TrueNtWriteVirtualMemory(ProcessHandle, BaseAddress, Buffer, BufferSize, NumberOfBytesWritten);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtWriteVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, BufferSize:0x%08lu), RETN: %p",
        ProcessHandle,
        BaseAddress,
        BufferSize,
        _ReturnAddress());

    ReleaseHookGuard();

    NTSTATUS Status = TrueNtWriteVirtualMemory(ProcessHandle, BaseAddress, Buffer, BufferSize, NumberOfBytesWritten);

    return Status;
}

NTSTATUS NTAPI
HookNtFreeVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _Inout_ PVOID *BaseAddress,
    _Inout_ PSIZE_T RegionSize,
    _In_ ULONG FreeType)
{
    if (IsInsideHook())
    {
        return TrueNtFreeVirtualMemory(ProcessHandle, BaseAddress, RegionSize, FreeType);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtFreeVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, RegionSize:0x%08lu), RETN: %p",
        ProcessHandle,
        *BaseAddress,
        *RegionSize,
        _ReturnAddress());

    NTSTATUS Status = TrueNtFreeVirtualMemory(ProcessHandle, BaseAddress, RegionSize, FreeType);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtMapViewOfSection(
    _In_ HANDLE SectionHandle,
    _In_ HANDLE ProcessHandle,
    _Inout_ _At_(*BaseAddress, _Readable_bytes_(*ViewSize) _Writable_bytes_(*ViewSize) _Post_readable_byte_size_(*ViewSize))
    PVOID *BaseAddress,
    _In_ ULONG_PTR ZeroBits,
    _In_ SIZE_T CommitSize,
    _Inout_opt_ PLARGE_INTEGER SectionOffset,
    _Inout_ PSIZE_T ViewSize,
    _In_ SECTION_INHERIT InheritDisposition,
    _In_ ULONG AllocationType,
    _In_ ULONG Win32Protect)
{
    if (IsInsideHook())
    {
        return TrueNtMapViewOfSection(
            SectionHandle,
            ProcessHandle,
            BaseAddress,
            ZeroBits,
            CommitSize,
            SectionOffset,
            ViewSize,
            InheritDisposition,
            AllocationType,
            Win32Protect);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtMapViewOfSection(SectionHandle: 0x%p, ProcessHandle:0x%p, BaseAddress:0x%p, AllocationType:0x%08lu), RETN: %p",
        SectionHandle,
        ProcessHandle,
        *BaseAddress,
        AllocationType,
        _ReturnAddress());

    NTSTATUS Status = TrueNtMapViewOfSection(
        SectionHandle,
        ProcessHandle,
        BaseAddress,
        ZeroBits,
        CommitSize,
        SectionOffset,
        ViewSize,
        InheritDisposition,
        AllocationType,
        Win32Protect);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtUnmapViewOfSection(_In_ HANDLE ProcessHandle, _In_opt_ PVOID BaseAddress)
{
    if (IsInsideHook())
    {
        return TrueNtUnmapViewOfSection(ProcessHandle, BaseAddress);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtUnmapViewOfSection(ProcessHandle:0x%p, BaseAddress:0x%p), RETN: %p",
        ProcessHandle,
        BaseAddress,
        _ReturnAddress());

    NTSTATUS Status = TrueNtUnmapViewOfSection(ProcessHandle, BaseAddress);

    ReleaseHookGuard();

    return Status;
}
