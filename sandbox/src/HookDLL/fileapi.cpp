#include "stdafx.h"
#include "fileapi.h"

decltype(NtCreateFile) *TrueNtCreateFile = nullptr;
decltype(NtReadFile) *TrueNtReadFile = nullptr;
decltype(NtWriteFile) *TrueNtWriteFile = nullptr;
decltype(NtDeleteFile) *TrueNtDeleteFile = nullptr;
decltype(NtSetInformationFile) *TrueNtSetInformationFile = nullptr;
decltype(NtQueryDirectoryFile) *TrueNtQueryDirectoryFile = nullptr;
decltype(NtQueryInformationFile) *TrueNtQueryInformationFile = nullptr;

NTSTATUS NTAPI
HookNtCreateFile(
    _Out_ PHANDLE FileHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_ POBJECT_ATTRIBUTES ObjectAttributes,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _In_opt_ PLARGE_INTEGER AllocationSize,
    _In_ ULONG FileAttributes,
    _In_ ULONG ShareAccess,
    _In_ ULONG CreateDisposition,
    _In_ ULONG CreateOptions,
    _In_reads_bytes_opt_(EaLength) PVOID EaBuffer,
    _In_ ULONG EaLength)
/*
- NtCreateFile
    - CreateFileA -> CreateFileW
*/
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    if (CreateOptions & FILE_DIRECTORY_FILE)
    {
        TraceAPI(
            L"CreateDirectory(%ws, DesiredAccess:0x%08x, CreateOptions:0x%08x), RETN: %p",
            ObjectAttributes->ObjectName->Buffer,
            DesiredAccess,
            CreateOptions,
            _ReturnAddress());
    }
    else
    {
        TraceAPI(
            L"NtCreateFile(%ws, DesiredAccess:0x%08x, CreateOptions:0x%08x), RETN: %p",
            ObjectAttributes->ObjectName->Buffer,
            DesiredAccess,
            CreateOptions,
            _ReturnAddress());
    }

    ReleaseHookGuard();

end:
    return TrueNtCreateFile(
        FileHandle,
        DesiredAccess,
        ObjectAttributes,
        IoStatusBlock,
        AllocationSize,
        FileAttributes,
        ShareAccess,
        CreateDisposition,
        CreateOptions,
        EaBuffer,
        EaLength);
}

NTSTATUS NTAPI
HookNtWriteFile(
    _In_ HANDLE FileHandle,
    _In_opt_ HANDLE Event,
    _In_opt_ PIO_APC_ROUTINE ApcRoutine,
    _In_opt_ PVOID ApcContext,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _In_reads_bytes_(Length) PVOID Buffer,
    _In_ ULONG Length,
    _In_opt_ PLARGE_INTEGER ByteOffset,
    _In_opt_ PULONG Key)
/*
- NtWriteFile
    - WriteFile
    - WriteFileEx
*/
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"NtWriteFile(FileHandle: %p), RETN: %p", FileHandle, _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtWriteFile(FileHandle, Event, ApcRoutine, ApcContext, IoStatusBlock, Buffer, Length, ByteOffset, Key);
}

NTSTATUS NTAPI
HookNtReadFile(
    _In_ HANDLE FileHandle,
    _In_opt_ HANDLE Event,
    _In_opt_ PIO_APC_ROUTINE ApcRoutine,
    _In_opt_ PVOID ApcContext,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _Out_writes_bytes_(Length) PVOID Buffer,
    _In_ ULONG Length,
    _In_opt_ PLARGE_INTEGER ByteOffset,
    _In_opt_ PULONG Key)
/*
- NtReadFile
    - ReadFile
    - ReadFileEx
*/
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"NtReadFile(FileHandle: %p), RETN: %p", FileHandle, _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtReadFile(FileHandle, Event, ApcRoutine, IoStatusBlock, IoStatusBlock, Buffer, Length, ByteOffset, Key);
}

NTSTATUS
WINAPI
HookNtDeleteFile(_In_ POBJECT_ATTRIBUTES ObjectAttributes)
/*
- NtDeleteFile
*/
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"NtDeleteFile(Filename:%ws), RETN: %p", ObjectAttributes->ObjectName->Buffer, _ReturnAddress());

    ReleaseHookGuard();
end:
    return TrueNtDeleteFile(ObjectAttributes);
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
        L"NtSetInformationFile(FileInformationClass: %d, FileInformation:0x%p, Length:0x%08x), RETN: %p",
        FileInformationClass,
        FileInformation,
        Length,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtSetInformationFile(FileHandle, IoStatusBlock, FileInformation, Length, FileInformationClass);
}

NTSTATUS WINAPI
HookNtQueryDirectoryFile(
    _In_ HANDLE FileHandle,
    _In_opt_ HANDLE Event,
    _In_opt_ PIO_APC_ROUTINE ApcRoutine,
    _In_opt_ PVOID ApcContext,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _Out_writes_bytes_(Length) PVOID FileInformation,
    _In_ ULONG Length,
    _In_ FILE_INFORMATION_CLASS FileInformationClass,
    _In_ BOOLEAN ReturnSingleEntry,
    _In_opt_ PUNICODE_STRING FileName,
    _In_ BOOLEAN RestartScan)
/*
- FindFirstFileA->FindFirstFileExW -> NtQueryDirectoryFile
- FindFirstFileW->FindFirstFileExW -> NtQueryDirectoryFile
- FindFirstFileExA -> FindFirstFileExW -> NtQueryDirectoryFile
- FindNextFileA -> NtQueryDirectoryFile
- FindNextFileW -> NtQueryDirectoryFile
*/
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryDirectoryFile(FileHandle: %p, FileInformationClass: %d, Length:0x%08x), RETN: %p",
        FileHandle,
        FileInformationClass,
        Length,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtQueryDirectoryFile(
        FileHandle,
        Event,
        ApcRoutine,
        ApcContext,
        IoStatusBlock,
        FileInformation,
        Length,
        FileInformationClass,
        ReturnSingleEntry,
        FileName,
        RestartScan);
}

NTSTATUS WINAPI
HookNtQueryInformationFile(
    _In_ HANDLE FileHandle,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _Out_writes_bytes_(Length) PVOID FileInformation,
    _In_ ULONG Length,
    _In_ FILE_INFORMATION_CLASS FileInformationClass)
/*
- GetFileSize -> GetFileSizeEx -> NtQueryInformationFile
- GetFileSizeEx -> NtQueryInformationFile.
*/
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryInformationFile(FileHandle: %p, FileInformationClass: %d, Length:0x%08x), RETN: %p",
        FileHandle,
        FileInformationClass,
        Length,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtQueryInformationFile(FileHandle, IoStatusBlock, FileInformation, Length, FileInformationClass);
}
