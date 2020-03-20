#include "stdafx.h"
#include "fileapi.h"

decltype(NtOpenFile) *TrueNtOpenFile = nullptr;
decltype(NtCreateFile) *TrueNtCreateFile = nullptr;
decltype(NtReadFile) *TrueNtReadFile = nullptr;
decltype(NtWriteFile) *TrueNtWriteFile = nullptr;
decltype(NtDeleteFile) *TrueNtDeleteFile = nullptr;
decltype(NtSetInformationFile) *TrueNtSetInformationFile = nullptr;
decltype(NtQueryDirectoryFile) *TrueNtQueryDirectoryFile = nullptr;
decltype(NtQueryInformationFile) *TrueNtQueryInformationFile = nullptr;

NTSTATUS NTAPI
HookNtOpenFile(
    _Out_ PHANDLE FileHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_ POBJECT_ATTRIBUTES ObjectAttributes,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _In_ ULONG ShareAccess,
    _In_ ULONG OpenOptions)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtOpenFile(FileHandle, DesiredAccess, ObjectAttributes, IoStatusBlock, ShareAccess, OpenOptions);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtOpenFile(%ws, DesiredAccess:0x%08x, FileAttributes: %lu, ShareAccess: %lu, OpenOptions:0x%08x), RETN: %p",
        ObjectAttributes->ObjectName->Buffer,
        DesiredAccess,
        ShareAccess,
        OpenOptions,
        _ReturnAddress());

    NTSTATUS Status =
        TrueNtOpenFile(FileHandle, DesiredAccess, ObjectAttributes, IoStatusBlock, ShareAccess, OpenOptions);

    ReleaseHookGuard();

    return Status;
}

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
    if (SfwIsCalledFromSystemMemory(4))
    {
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

    CaptureStackTrace();

    if (CreateOptions & FILE_DIRECTORY_FILE)
    {
        TraceAPI(
            L"CreateDirectory(%ws, DesiredAccess:0x%08x, FileAttributes: %lu, ShareAccess: %lu, CreateOptions:0x%08x), RETN: %p",
            ObjectAttributes->ObjectName->Buffer,
            DesiredAccess,
            FileAttributes,
            ShareAccess,
            CreateOptions,
            _ReturnAddress());
    }
    else
    {
        TraceAPI(
            L"NtCreateFile(%ws, DesiredAccess:0x%08x, FileAttributes: %lu, ShareAccess: %lu, CreateOptions:0x%08x), RETN: %p",
            ObjectAttributes->ObjectName->Buffer,
            DesiredAccess,
            FileAttributes,
            ShareAccess,
            CreateOptions,
            _ReturnAddress());
    }
    NTSTATUS Status = TrueNtCreateFile(
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

    ReleaseHookGuard();

    return Status;
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
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtWriteFile(
            FileHandle, Event, ApcRoutine, ApcContext, IoStatusBlock, Buffer, Length, ByteOffset, Key);
    }

    CaptureStackTrace();

    TraceAPI(L"NtWriteFile(FileHandle: %p), RETN: %p", FileHandle, _ReturnAddress());

    NTSTATUS Status =
        TrueNtWriteFile(FileHandle, Event, ApcRoutine, ApcContext, IoStatusBlock, Buffer, Length, ByteOffset, Key);

    ReleaseHookGuard();

    return Status;
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
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtReadFile(
            FileHandle, Event, ApcRoutine, IoStatusBlock, IoStatusBlock, Buffer, Length, ByteOffset, Key);
    }

    CaptureStackTrace();

    TraceAPI(L"NtReadFile(FileHandle: %p), RETN: %p", FileHandle, _ReturnAddress());

    NTSTATUS Status =
        TrueNtReadFile(FileHandle, Event, ApcRoutine, IoStatusBlock, IoStatusBlock, Buffer, Length, ByteOffset, Key);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS
NTAPI
HookNtDeleteFile(_In_ POBJECT_ATTRIBUTES ObjectAttributes)
/*
- NtDeleteFile
*/
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtDeleteFile(ObjectAttributes);
    }

    CaptureStackTrace();

    TraceAPI(L"NtDeleteFile(Filename:%ws), RETN: %p", ObjectAttributes->ObjectName->Buffer, _ReturnAddress());

    NTSTATUS Status = TrueNtDeleteFile(ObjectAttributes);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtSetInformationFile(
    _In_ HANDLE FileHandle,
    _Out_ PIO_STATUS_BLOCK IoStatusBlock,
    _In_reads_bytes_(Length) PVOID FileInformation,
    _In_ ULONG Length,
    _In_ FILE_INFORMATION_CLASS FileInformationClass)
/*
    - MoveFile => NtSetInformationFile(10)
    - DeleteFile => NtSetInformationFile(13)
*/
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtSetInformationFile(FileHandle, IoStatusBlock, FileInformation, Length, FileInformationClass);
    }
    CaptureStackTrace();

    TraceAPI(
        L"NtSetInformationFile(FileInformationClass: %d, FileInformation:0x%p, Length:0x%08x), RETN: %p",
        FileInformationClass,
        FileInformation,
        Length,
        _ReturnAddress());

    NTSTATUS Status =
        TrueNtSetInformationFile(FileHandle, IoStatusBlock, FileInformation, Length, FileInformationClass);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
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
    if (SfwIsCalledFromSystemMemory(5))
    {
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

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryDirectoryFile(FileHandle: %p, FileInformationClass: %d, Length:0x%08x), RETN: %p",
        FileHandle,
        FileInformationClass,
        Length,
        _ReturnAddress());

    NTSTATUS Status = TrueNtQueryDirectoryFile(
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

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
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
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtQueryInformationFile(FileHandle, IoStatusBlock, FileInformation, Length, FileInformationClass);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtQueryInformationFile(FileHandle: %p, FileInformationClass: %d, Length:0x%08x), RETN: %p",
        FileHandle,
        FileInformationClass,
        Length,
        _ReturnAddress());

    NTSTATUS Status =
        TrueNtQueryInformationFile(FileHandle, IoStatusBlock, FileInformation, Length, FileInformationClass);

    ReleaseHookGuard();

    return Status;
}
