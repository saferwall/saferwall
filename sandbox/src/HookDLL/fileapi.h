#pragma once

#include "stdafx.h"


NTSTATUS
WINAPI
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
	_In_ ULONG EaLength
);


NTSTATUS
NTAPI
HookNtReadFile(
	_In_ HANDLE FileHandle,
	_In_opt_ HANDLE Event,
	_In_opt_ PIO_APC_ROUTINE ApcRoutine,
	_In_opt_ PVOID ApcContext,
	_Out_ PIO_STATUS_BLOCK IoStatusBlock,
	_Out_writes_bytes_(Length) PVOID Buffer,
	_In_ ULONG Length,
	_In_opt_ PLARGE_INTEGER ByteOffset,
	_In_opt_ PULONG Key
);

NTSTATUS
NTAPI
HookNtWriteFile(
	_In_ HANDLE FileHandle,
	_In_opt_ HANDLE Event,
	_In_opt_ PIO_APC_ROUTINE ApcRoutine,
	_In_opt_ PVOID ApcContext,
	_Out_ PIO_STATUS_BLOCK IoStatusBlock,
	_In_reads_bytes_(Length) PVOID Buffer,
	_In_ ULONG Length,
	_In_opt_ PLARGE_INTEGER ByteOffset,
	_In_opt_ PULONG Key
);


NTSTATUS
NTAPI
HookNtDeleteFile(
	_In_ POBJECT_ATTRIBUTES ObjectAttributes
);


NTSTATUS
WINAPI
HookMoveFileWithProgressTransactedW(
	__in      LPWSTR lpExistingFileName,
	__in_opt  LPWSTR lpNewFileName,
	__in_opt  LPPROGRESS_ROUTINE lpProgressRoutine,
	__in_opt  LPVOID lpData,
	__in      DWORD dwFlags,
	__in	  HANDLE hTransaction
);

typedef NTSTATUS (WINAPI * pfnMoveFileWithProgressTransactedW)
(
	__in      LPWSTR lpExistingFileName,
	__in_opt  LPWSTR lpNewFileName,
	__in_opt  LPPROGRESS_ROUTINE lpProgressRoutine,
	__in_opt  LPVOID lpData,
	__in      DWORD dwFlags,
	__in	  HANDLE hTransaction
);

