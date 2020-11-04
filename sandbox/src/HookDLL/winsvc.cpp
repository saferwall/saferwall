#include "stdafx.h"
#include "winsvc.h"

extern pfnCreateServiceW TrueCreateServiceW;

SC_HANDLE
WINAPI
HookCreateServiceW(
	_In_        SC_HANDLE    hSCManager,
	_In_        LPCWSTR     lpServiceName,
	_In_opt_    LPCWSTR     lpDisplayName,
	_In_        DWORD        dwDesiredAccess,
	_In_        DWORD        dwServiceType,
	_In_        DWORD        dwStartType,
	_In_        DWORD        dwErrorControl,
	_In_opt_    LPCWSTR     lpBinaryPathName,
	_In_opt_    LPCWSTR     lpLoadOrderGroup,
	_Out_opt_   LPDWORD      lpdwTagId,
	_In_opt_    LPCWSTR     lpDependencies,
	_In_opt_    LPCWSTR     lpServiceStartName,
	_In_opt_    LPCWSTR     lpPassword)
{
	if (SfwIsCalledFromSystemMemory(5))
	{
		return TrueCreateServiceW(hSCManager, lpServiceName, lpDisplayName, dwDesiredAccess,
			dwServiceType, dwStartType, dwErrorControl, lpBinaryPathName, lpLoadOrderGroup,
			lpdwTagId, lpDependencies, lpServiceStartName, lpPassword);
	}

	CaptureStackTrace();

	TraceAPI(L"CreateServiceW(ServiceName: %ws, lpDisplayName: %ws, dwServiceType: %d, dwStartType: %d), RETN: 0x%p", 
		lpServiceName, lpDisplayName, dwServiceType, dwStartType, _ReturnAddress());

	SC_HANDLE hService = TrueCreateServiceW(hSCManager, lpServiceName, lpDisplayName, dwDesiredAccess,
		dwServiceType, dwStartType, dwErrorControl, lpBinaryPathName, lpLoadOrderGroup,
		lpdwTagId, lpDependencies, lpServiceStartName, lpPassword);

	ReleaseHookGuard();

	return hService;
}
