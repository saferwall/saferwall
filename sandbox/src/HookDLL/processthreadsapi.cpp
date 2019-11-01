#include "stdafx.h"
#include "processthreadsapi.h"


decltype(NtCreateUserProcess)* TrueNtCreateUserProcess = nullptr;



NTSTATUS NTAPI HookNtCreateUserProcess(
	_Out_ PHANDLE ProcessHandle,
	_Out_ PHANDLE ThreadHandle,
	_In_ ACCESS_MASK ProcessDesiredAccess,
	_In_ ACCESS_MASK ThreadDesiredAccess,
	_In_opt_ POBJECT_ATTRIBUTES ProcessObjectAttributes,
	_In_opt_ POBJECT_ATTRIBUTES ThreadObjectAttributes,
	_In_ ULONG ProcessFlags, // PROCESS_CREATE_FLAGS_*
	_In_ ULONG ThreadFlags, // THREAD_CREATE_FLAGS_*
	_In_opt_ PVOID ProcessParameters, // PRTL_USER_PROCESS_PARAMETERS
	_Inout_ PPS_CREATE_INFO CreateInfo,
	_In_opt_ PPS_ATTRIBUTE_LIST AttributeList
)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtCreateUserProcess(%ws), RETN: %p", AttributeList->Attributes[0].Value, _ReturnAddress());

	ReleaseHookGuard();

end:
	return TrueNtCreateUserProcess(ProcessHandle, ThreadHandle, ProcessDesiredAccess,
		ThreadDesiredAccess, ProcessObjectAttributes, ThreadObjectAttributes, ProcessFlags,
		ThreadFlags, ProcessParameters, CreateInfo, AttributeList);
}