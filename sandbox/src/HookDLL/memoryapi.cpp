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
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtAllocateVirtualMemory(ProcessHandle:0x%p, AllocationType:%lu, Protect:%lu), RETN: 0x%p",
        ProcessHandle,
        AllocationType,
        Protect,
        _ReturnAddress()); // , Protect

    ReleaseHookGuard();
end:
    return TrueNtAllocateVirtualMemory(ProcessHandle, BaseAddress, ZeroBits, RegionSize, AllocationType, Protect);
}

NTSTATUS WINAPI
HookNtProtectVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _Inout_ PVOID *BaseAddress,
    _Inout_ PSIZE_T RegionSize,
    _In_ ULONG NewProtect,
    _Out_ PULONG OldProtect)
{
    if (IsInsideHook())
    {
        goto end;
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

    ReleaseHookGuard();
end:
    return TrueNtProtectVirtualMemory(ProcessHandle, BaseAddress, RegionSize, NewProtect, OldProtect);
    ;
}

NTSTATUS WINAPI
HookNtQueryVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _In_opt_ PVOID BaseAddress,
    _In_ MEMORY_INFORMATION_CLASS MemoryInformationClass,
    _Out_writes_bytes_(MemoryInformationLength) PVOID MemoryInformation,
    _In_ SIZE_T MemoryInformationLength,
    _Out_opt_ PSIZE_T ReturnLength)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, MemoryInformationClass:%d, MemoryInformationLength: %lu), RETN: %p",
        ProcessHandle,
        BaseAddress,
        MemoryInformationClass,
        MemoryInformationLength,
        _ReturnAddress());
    ReleaseHookGuard();
end:
    return TrueNtQueryVirtualMemory(
        ProcessHandle, BaseAddress, MemoryInformationClass, MemoryInformation, MemoryInformationLength, ReturnLength);
}

NTSTATUS WINAPI
HookNtReadVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _In_opt_ PVOID BaseAddress,
    _Out_writes_bytes_(BufferSize) PVOID Buffer,
    _In_ SIZE_T BufferSize,
    _Out_opt_ PSIZE_T NumberOfBytesRead)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtReadVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, BufferSize:0x%08lu), RETN: %p",
        ProcessHandle,
        BaseAddress,
        BufferSize,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtReadVirtualMemory(ProcessHandle, BaseAddress, Buffer, BufferSize, NumberOfBytesRead);
    ;
}

NTSTATUS WINAPI
HookNtWriteVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _In_opt_ PVOID BaseAddress,
    _In_reads_bytes_(BufferSize) PVOID Buffer,
    _In_ SIZE_T BufferSize,
    _Out_opt_ PSIZE_T NumberOfBytesWritten)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtWriteVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, BufferSize:0x%08lu), RETN: %p",
        ProcessHandle,
        BaseAddress,
        BufferSize,
        _ReturnAddress());

    ReleaseHookGuard();
end:
    return TrueNtWriteVirtualMemory(ProcessHandle, BaseAddress, Buffer, BufferSize, NumberOfBytesWritten);
    ;
}

NTSTATUS WINAPI
HookNtFreeVirtualMemory(
    _In_ HANDLE ProcessHandle,
    _Inout_ PVOID *BaseAddress,
    _Inout_ PSIZE_T RegionSize,
    _In_ ULONG FreeType)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtFreeVirtualMemory(ProcessHandle:0x%p, BaseAddress:0x%p, RegionSize:0x%08lu), RETN: %p",
        ProcessHandle,
        *BaseAddress,
        *RegionSize,
        _ReturnAddress());

    ReleaseHookGuard();
end:
    return TrueNtFreeVirtualMemory(ProcessHandle, BaseAddress, RegionSize, FreeType);
    ;
}

NTSTATUS WINAPI
HookNtMapViewOfSection(
    _In_ HANDLE SectionHandle,
    _In_ HANDLE ProcessHandle,
    _Inout_
        _At_(*BaseAddress, _Readable_bytes_(*ViewSize) _Writable_bytes_(*ViewSize) _Post_readable_byte_size_(*ViewSize))
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
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtMapViewOfSection(SectionHandle: 0x%p, ProcessHandle:0x%p, BaseAddress:0x%p, AllocationType:0x%08lu), RETN: %p",
        SectionHandle,
        ProcessHandle,
        *BaseAddress,
        AllocationType,
        _ReturnAddress());

    ReleaseHookGuard();

end:
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

NTSTATUS WINAPI
HookNtUnmapViewOfSection(_In_ HANDLE ProcessHandle, _In_opt_ PVOID BaseAddress)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtUnmapViewOfSection(ProcessHandle:0x%p, BaseAddress:0x%p), RETN: %p",
        ProcessHandle,
        BaseAddress,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtUnmapViewOfSection(ProcessHandle, BaseAddress);
}
