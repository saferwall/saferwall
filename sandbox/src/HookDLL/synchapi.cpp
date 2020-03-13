#include "stdafx.h"
#include "synchapi.h"

decltype(NtDelayExecution) *TrueNtDelayExecution = nullptr;

NTSTATUS WINAPI
HookNtDelayExecution(_In_ BOOLEAN Alertable, _In_opt_ PLARGE_INTEGER DelayInterval)
{
    if (IsInsideHook())
    {
        return TrueNtDelayExecution(Alertable, DelayInterval);
    }

    if (DelayInterval->QuadPart == 3)
    {
        NTSTATUS Status;

        UNICODE_STRING ModulePath;
        HANDLE ModuleHandle = NULL;

        RtlInitUnicodeString(&ModulePath, (PWSTR)L"ole32.dll");
        Status = LdrGetDllHandle(NULL, 0, &ModulePath, &ModuleHandle);
        if (Status == STATUS_SUCCESS)
        {
            LogMessage(L"Attaching to ole32");
            HookOleAPIs(TRUE);
            LogMessage(L"Hooked OLE");
        }

         RtlInitUnicodeString(&ModulePath, (PWSTR)L"wininet.dll");
         Status = LdrGetDllHandle(NULL, 0, &ModulePath, &ModuleHandle);
         if (Status == STATUS_SUCCESS)
        {
            LogMessage(L"Attaching to wininet");
            HookNetworkAPIs(TRUE);
            LogMessage(L"Hooked wininet");
        }
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtDelayExecution(Alertable: %d, DelayInterval: %I64u), RETN: %p",
        Alertable,
        DelayInterval->QuadPart,
        _ReturnAddress());
    NTSTATUS Status = TrueNtDelayExecution(Alertable, DelayInterval);
    ReleaseHookGuard();
    return Status;
}
