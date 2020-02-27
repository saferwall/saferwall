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
		//printf("\nprocess attach\n");
		SetupHook();
	case DLL_THREAD_ATTACH:
		//printf("\nthread attach\n");

	case DLL_THREAD_DETACH:
		//printf("\nthread detach\n");

	case DLL_PROCESS_DETACH:
		//printf("\nProcess detach\n");

		//SymCleanup(GetCurrentProcess());
		//Unhook();
		break;
	}

	return TRUE;
}

