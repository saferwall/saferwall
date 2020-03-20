// dllmain.cpp : Defines the entry point for the DLL application.
#include "stdafx.h"

BOOL APIENTRY
DllMain(HMODULE hModule, DWORD dwReason, LPVOID lpReserved)
{
    //
    // When creating a 32-bit target process from a 64-bit parent process or
    // creating a 64-bit target process from a 32-bit parent process, we must
    // check if the current process is helper process or not.

    // if (DetourIsHelperProcess()) {
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

        ProcessAttach();
        LogMessage(L"DllMain DLL_PROCESS_ATTACH Done");
        break;

    case DLL_PROCESS_DETACH:

        //
        // Removes all the hooks.
        //

        ProcessDetach();
        LogMessage(L"DllMain DLL_PROCESS_DETACH Done");
        break;

    case DLL_THREAD_ATTACH:
        // LogMessage(L"DllMain DLL_THREAD_ATTACH");
        break;

    case DLL_THREAD_DETACH:
        // LogMessage(L"DllMain DLL_THREAD_DETACH");
        break;
    }

    return TRUE;
}
