#include "stdafx.h"
#include "fileapi.h"


decltype(NtCreateFile)* TrueNtCreateFile = nullptr;
decltype(NtReadFile)* TrueNtReadFile = nullptr;
decltype(NtWriteFile)* TrueNtWriteFile = nullptr;
decltype(NtDeleteFile)* TrueNtDeleteFile = nullptr;
pfnMoveFileWithProgressTransactedW TrueMoveFileWithProgressTransactedW = nullptr;



NTSTATUS NTAPI HookNtCreateFile(_Out_ PHANDLE FileHandle,
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
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	if (CreateOptions & FILE_DIRECTORY_FILE) {
		TraceAPI(L"CreateDirectory(%ws, DesiredAccess:0x%08x, CreateOptions:0x%08x), RETN: %p",
			ObjectAttributes->ObjectName->Buffer,
			DesiredAccess, CreateOptions, _ReturnAddress());
	}

	TraceAPI(L"NtCreateFile(%ws, DesiredAccess:0x%08x, CreateOptions:0x%08x), RETN: %p",
		ObjectAttributes->ObjectName->Buffer, DesiredAccess, CreateOptions, _ReturnAddress());
	
	ReleaseHookGuard();

end:
	return TrueNtCreateFile(FileHandle, DesiredAccess, ObjectAttributes, IoStatusBlock, AllocationSize,
		FileAttributes, ShareAccess, CreateDisposition, CreateOptions, EaBuffer, EaLength);
}



NTSTATUS NTAPI HookNtWriteFile(
	_In_ HANDLE FileHandle,
	_In_opt_ HANDLE Event,
	_In_opt_ PIO_APC_ROUTINE ApcRoutine,
	_In_opt_ PVOID ApcContext,
	_Out_ PIO_STATUS_BLOCK IoStatusBlock,
	_In_reads_bytes_(Length) PVOID Buffer,
	_In_ ULONG Length,
	_In_opt_ PLARGE_INTEGER ByteOffset,
	_In_opt_ PULONG Key
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtWriteFile(FileHandle: %p), RETN: %p", FileHandle, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtWriteFile(FileHandle, Event, ApcRoutine, ApcContext,
		IoStatusBlock, Buffer, Length, ByteOffset, Key);
}


NTSTATUS NTAPI HookNtReadFile(
	_In_ HANDLE FileHandle,
	_In_opt_ HANDLE Event,
	_In_opt_ PIO_APC_ROUTINE ApcRoutine,
	_In_opt_ PVOID ApcContext,
	_Out_ PIO_STATUS_BLOCK IoStatusBlock,
	_Out_writes_bytes_(Length) PVOID Buffer,
	_In_ ULONG Length,
	_In_opt_ PLARGE_INTEGER ByteOffset,
	_In_opt_ PULONG Key
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtReadFile(FileHandle: %p), RETN: %p",
		FileHandle, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtReadFile(FileHandle, Event, ApcRoutine, IoStatusBlock, IoStatusBlock,
		Buffer, Length, ByteOffset, Key);
}



NTSTATUS
WINAPI
HookNtDeleteFile(
	_In_ POBJECT_ATTRIBUTES ObjectAttributes
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtDeleteFile(Filename:%ws), RETN: %p", ObjectAttributes->ObjectName->Buffer, _ReturnAddress());

	ReleaseHookGuard();
end:
	return TrueNtDeleteFile(ObjectAttributes);
}

NTSTATUS WINAPI HookMoveFileWithProgressTransactedW(
	__in      LPWSTR lpExistingFileName,
	__in_opt  LPWSTR lpNewFileName,
	__in_opt  LPPROGRESS_ROUTINE lpProgressRoutine,
	__in_opt  LPVOID lpData,
	__in      DWORD dwFlags,
	__in	  HANDLE hTransaction)
{

	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();
	
	TraceAPI(L"MoveFileWithProgressTransactedW(lpExistingFileName:%ws, lpNewFileName:%ws), RETN: %p",
		lpExistingFileName, lpNewFileName, _ReturnAddress());

	ReleaseHookGuard();
end:
	return TrueMoveFileWithProgressTransactedW(lpExistingFileName, lpNewFileName,
		lpProgressRoutine, lpData, dwFlags, hTransaction);
}