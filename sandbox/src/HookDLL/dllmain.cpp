// dllmain.cpp : Defines the entry point for the DLL application.
#include "stdafx.h"




BOOL APIENTRY DllMain( HMODULE hModule,
                       DWORD  dwReason,
                       LPVOID lpReserved
                     )
{
	//
	// When creating a 32-bit target process from a 64-bit parent process or
	// creating a 64-bit target process from a 32-bit parent process, we must
	// check if the current process is helper process or not.

	//if (DetourIsHelperProcess()) {
	//	return TRUE;
	//}

	switch (dwReason)
	{
	case DLL_PROCESS_ATTACH:

		//
		// Restore the contents in memory import table after a process was started
		// with DetourCreateProcessWithDllEx or DetourCreateProcessWithDlls.
		// As we are mapping the dll from kernel, we don't need to call DetourRestoreAfterWith();
		// 
		OutputDebugStringA("HookDLL" DETOURS_STRINGIFY(DETOURS_BITS) ".dll:"
                                                                    " DllMain DLL_PROCESS_ATTACH\n");
		ProcessAttach();
		break;

	case DLL_PROCESS_DETACH:

		//
		// Removes all the hooks.
		//

		ProcessDetach();
		break;


	case DLL_THREAD_ATTACH:
		OutputDebugStringA("HookDLL" DETOURS_STRINGIFY(DETOURS_BITS) ".dll:"
			" DllMain DLL_THREAD_ATTACH\n");
		break;

	case DLL_THREAD_DETACH:
		OutputDebugStringA("HookDLL" DETOURS_STRINGIFY(DETOURS_BITS) ".dll:"
			" DllMain DLL_THREAD_DETACH\n");
		break;
	}

	return TRUE;
}

