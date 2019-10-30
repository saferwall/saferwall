#include "stdafx.h"
#include "systemapi.h"

decltype(NtQuerySystemInformation)* TrueNtQuerySystemInformation = nullptr;


NTSTATUS WINAPI HookTrueNtQuerySystemInformation(
	_In_ SYSTEM_INFORMATION_CLASS SystemInformationClass,
	_Out_writes_bytes_opt_(SystemInformationLength) PVOID SystemInformation,
	_In_ ULONG SystemInformationLength,
	_Out_opt_ PULONG ReturnLength
)
{

	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	TraceAPI(L"NtQuerySystemInformation(SystemInformationClass: 0x%d, SystemInformation:0x%p, SystemInformationLength:0x%08x, ReturnLength:0x%p), RETN: %p",
		SystemInformationClass, SystemInformation, SystemInformationLength, ReturnLength, _ReturnAddress());

	ReleaseHookGuard();

end:
 return TrueNtQuerySystemInformation(SystemInformationClass, SystemInformation, SystemInformationLength, ReturnLength);

}
