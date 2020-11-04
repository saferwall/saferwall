#include "header.h"

DWORD WINAPI
ThreadFunc(void *data)
{
    // Do stuff.  This will be the first function called on the new thread.
    // When this function returns, the thread goes away.  See MSDN for more details.
    return 0;
}


VOID
TestProcessThreadHooks()
{

    STARTUPINFO info = {sizeof(info)};
    DWORD dwPid;
    PROCESS_INFORMATION processInfo;

    wprintf(L"[+] Calling CreateProcess\n");
    if (CreateProcess(L"C:\\Windows\\notepad.exe", NULL, NULL, NULL, TRUE, 0, NULL, NULL, &info, &processInfo))
    {
        dwPid = GetProcessId(processInfo.hProcess);
        WaitForSingleObject(processInfo.hProcess, INFINITE);
    }

    wprintf(L"[+] Calling CreateThread\n");
    HANDLE hThread = CreateThread(NULL, 0, ThreadFunc, NULL, 0, NULL);
    if (hThread)
    {
        // Optionally do stuff, such as wait on the thread.
    }
    wprintf(L"[+] Calling OpenProcess\n");
    HANDLE hProcess = OpenProcess(PROCESS_ALL_ACCESS, FALSE, dwPid);

    CloseHandle(processInfo.hProcess);
    CloseHandle(processInfo.hThread);
    CloseHandle(hThread);
}