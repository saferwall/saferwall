#include "stdafx.h"
#include "libloaderapi.h"


decltype(LdrLoadDll)* TrueLdrLoadDll = nullptr;
decltype(LdrGetProcedureAddress)* TrueLdrGetProcedureAddress = nullptr;



NTSTATUS WINAPI HookLdrLoadDll(PWSTR  DllPath, PULONG  DllCharacteristics, PUNICODE_STRING DllName, PVOID* DllHandle)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}

	GetStackWalk();

	if (DllName && DllName->Buffer) {
		TraceAPI(L"LdrLoadDll(%ws), RETN: 0x%p", DllName->Buffer, _ReturnAddress());
	}

	ReleaseHookGuard();
end:
	return TrueLdrLoadDll(DllPath, DllCharacteristics, DllName, DllHandle);
}


NTSTATUS WINAPI HookLdrGetProcedureAddress(PVOID DllHandle, PANSI_STRING ProcedureName, ULONG ProcedureNumber, PVOID *ProcedureAddress)
{
	if (IsInsideHook() == FALSE) {
		goto end;
	}
	GetStackWalk();

	TraceAPI(L"LdrGetProcedureAddress(%ws), RETN: 0x%p", MultiByteToWide(ProcedureName->Buffer), _ReturnAddress());

	ReleaseHookGuard();
end:
	return TrueLdrGetProcedureAddress(DllHandle, ProcedureName, ProcedureNumber, ProcedureAddress);
}
