#include "stdafx.h"
#include "synchapi.h"

decltype(NtDelayExecution) *TrueNtDelayExecution = nullptr;

NTSTATUS NTAPI
HookNtDelayExecution(_In_ BOOLEAN Alertable, _In_opt_ PLARGE_INTEGER DelayInterval)
{
    if (IsInsideHook())
    {
        return TrueNtDelayExecution(Alertable, DelayInterval);
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
