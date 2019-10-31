#include "stdafx.h"
#include "synchapi.h"

decltype(NtOpenKey)* TrueNtOpenKey = nullptr;
decltype(NtOpenKeyEx)* TrueNtOpenKeyEx = nullptr;
decltype(NtCreateKey)* TrueNtCreateKey = nullptr;
decltype(NtQueryValueKey)* TrueNtQueryValueKey = nullptr;
decltype(NtDeleteKey)* TrueNtDeleteKey = nullptr;
decltype(NtDeleteValueKey)* TrueNtDeleteValueKey = nullptr;


NTSTATUS WINAPI HookNtOpenKey(
	_Out_ PHANDLE KeyHandle,
	_In_ ACCESS_MASK DesiredAccess,
	_In_ POBJECT_ATTRIBUTES ObjectAttributes
)
{

	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtOpenKey(DesiredAccess: 0x%d, ObjectName:0x%ws, ReturnLength:0x%p), RETN: %p",
		DesiredAccess, ObjectAttributes->ObjectName->Buffer, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtOpenKey(KeyHandle, DesiredAccess, ObjectAttributes);

}



NTSTATUS WINAPI HookNtOpenKeyEx(
	_Out_ PHANDLE KeyHandle,
	_In_ ACCESS_MASK DesiredAccess,
	_In_ POBJECT_ATTRIBUTES ObjectAttributes,
	_In_ ULONG OpenOptions
)
{

	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtOpenKeyEx(DesiredAccess: 0x%d, ObjectName:0x%ws, ReturnLength:0x%p), RETN: %p",
		DesiredAccess, ObjectAttributes->ObjectName->Buffer, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtOpenKeyEx(KeyHandle, DesiredAccess, ObjectAttributes, OpenOptions);

}



NTSTATUS WINAPI HookNtCreateKey(
	_Out_ PHANDLE KeyHandle,
	_In_ ACCESS_MASK DesiredAccess,
	_In_ POBJECT_ATTRIBUTES ObjectAttributes,
	_Reserved_ ULONG TitleIndex,
	_In_opt_ PUNICODE_STRING Class,
	_In_ ULONG CreateOptions,
	_Out_opt_ PULONG Disposition
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtCreateKey(DesiredAccess: 0x%d, ObjectName:0x%ws, CreateOptions: %ul, ReturnLength:0x%p), RETN: %p",
		DesiredAccess, ObjectAttributes->ObjectName->Buffer, CreateOptions, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtCreateKey(KeyHandle, DesiredAccess, ObjectAttributes, TitleIndex, Class, CreateOptions, Disposition);

}



NTSTATUS WINAPI HookNtQueryValueKey(
	_In_ HANDLE KeyHandle,
	_In_ PUNICODE_STRING ValueName,
	_In_ KEY_VALUE_INFORMATION_CLASS KeyValueInformationClass,
	_Out_writes_bytes_opt_(Length) PVOID KeyValueInformation,
	_In_ ULONG Length,
	_Out_ PULONG ResultLength
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtQueryValueKey(KeyHandle: 0x%d, ValueName:0x%ws), RETN: %p",
		KeyHandle, ValueName->Buffer, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtQueryValueKey(KeyHandle, ValueName, KeyValueInformationClass, KeyValueInformation, Length, ResultLength);

}



NTSTATUS WINAPI HookNtDeleteKey(
	_In_ HANDLE KeyHandle
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtDeleteKey(KeyHandle: 0x%d), RETN: %p",
		KeyHandle, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtDeleteKey(KeyHandle);

}


NTSTATUS WINAPI HookNtDeleteValueKey(
	_In_ HANDLE KeyHandle,
	_In_ PUNICODE_STRING ValueName
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtDeleteValueKey(KeyHandle: 0x%d, ValueName: %ws), RETN: %p",
		KeyHandle, ValueName, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtDeleteValueKey(KeyHandle, ValueName);

}