#include "stdafx.h"
#include "synchapi.h"

decltype(NtDelayExecution) *TrueNtDelayExecution = nullptr;

NTSTATUS WINAPI
HookNtDelayExecution(_In_ BOOLEAN Alertable, _In_opt_ PLARGE_INTEGER DelayInterval)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtDelayExecution(Alertable: 0x%d, DelayInterval:0x%, ReturnLength:0x%p), RETN: %p",
        Alertable,
        DelayInterval->QuadPart,
        _ReturnAddress());

    ReleaseHookGuard();

end:
    return TrueNtDelayExecution(Alertable, DelayInterval);
}
