#include "header.h"

VOID
TestSyncHooks()
{
    UINT delayInMillis = 3;
    LARGE_INTEGER DelayInterval;
    LONGLONG llDelay = delayInMillis * 1LL;
    DelayInterval.QuadPart = llDelay;
    static NTSTATUS(__stdcall * NtDelayExecution)(IN BOOLEAN Alertable, IN PLARGE_INTEGER DelayInterval) = (NTSTATUS(
        __stdcall *)(BOOLEAN, PLARGE_INTEGER))GetProcAddress(GetModuleHandle(L"ntdll.dll"), "NtDelayExecution");
    NtDelayExecution(FALSE, &DelayInterval);
}