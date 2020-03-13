#include "stdafx.h"
#include "processthreadsapi.h"

decltype(NtCreateUserProcess) *TrueNtCreateUserProcess = nullptr;
decltype(NtCreateThread) *TrueNtCreateThread = nullptr;
decltype(NtCreateThreadEx) *TrueNtCreateThreadEx = nullptr;
decltype(NtResumeThread) *TrueNtResumeThread = nullptr;
decltype(NtSuspendThread) *TrueNtSuspendThread = nullptr;
decltype(NtOpenProcess) *TrueNtOpenProcess = nullptr;
decltype(NtTerminateProcess) *TrueNtTerminateProcess = nullptr;
decltype(NtContinue) *TrueNtContinue = nullptr;

BOOL bFirstTime = TRUE;

NTSTATUS NTAPI
HookNtCreateUserProcess(
    _Out_ PHANDLE ProcessHandle,
    _Out_ PHANDLE ThreadHandle,
    _In_ ACCESS_MASK ProcessDesiredAccess,
    _In_ ACCESS_MASK ThreadDesiredAccess,
    _In_opt_ POBJECT_ATTRIBUTES ProcessObjectAttributes,
    _In_opt_ POBJECT_ATTRIBUTES ThreadObjectAttributes,
    _In_ ULONG ProcessFlags,          // PROCESS_CREATE_FLAGS_*
    _In_ ULONG ThreadFlags,           // THREAD_CREATE_FLAGS_*
    _In_opt_ PVOID ProcessParameters, // PRTL_USER_PROCESS_PARAMETERS
    _Inout_ PPS_CREATE_INFO CreateInfo,
    _In_opt_ PPS_ATTRIBUTE_LIST AttributeList)
{
    if (IsInsideHook())
    {
        return TrueNtCreateUserProcess(
            ProcessHandle,
            ThreadHandle,
            ProcessDesiredAccess,
            ThreadDesiredAccess,
            ProcessObjectAttributes,
            ThreadObjectAttributes,
            ProcessFlags,
            ThreadFlags,
            ProcessParameters,
            CreateInfo,
            AttributeList);
    }

    CaptureStackTrace();

    TraceAPI(L"NtCreateUserProcess(%ws), RETN: %p", AttributeList->Attributes[0].Value, _ReturnAddress());

    NTSTATUS Status = TrueNtCreateUserProcess(
        ProcessHandle,
        ThreadHandle,
        ProcessDesiredAccess,
        ThreadDesiredAccess,
        ProcessObjectAttributes,
        ThreadObjectAttributes,
        ProcessFlags,
        ThreadFlags,
        ProcessParameters,
        CreateInfo,
        AttributeList);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtCreateThread(
    _Out_ PHANDLE ThreadHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_opt_ POBJECT_ATTRIBUTES ObjectAttributes,
    _In_ HANDLE ProcessHandle,
    _Out_ PCLIENT_ID ClientId,
    _In_ PCONTEXT ThreadContext,
    _In_ PINITIAL_TEB InitialTeb,
    _In_ BOOLEAN CreateSuspended)
{
    if (IsInsideHook())
    {
        return TrueNtCreateThread(
            ThreadHandle,
            DesiredAccess,
            ObjectAttributes,
            ProcessHandle,
            ClientId,
            ThreadContext,
            InitialTeb,
            CreateSuspended);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtCreateThread(DesiredAccess: %d, ProcessHandle: %p, CreateSuspended: %d), RETN: %p",
        DesiredAccess,
        ProcessHandle,
        CreateSuspended,
        _ReturnAddress());

    NTSTATUS Status = TrueNtCreateThread(
        ThreadHandle,
        DesiredAccess,
        ObjectAttributes,
        ProcessHandle,
        ClientId,
        ThreadContext,
        InitialTeb,
        CreateSuspended);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtCreateThreadEx(
    _Out_ PHANDLE ThreadHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_opt_ POBJECT_ATTRIBUTES ObjectAttributes,
    _In_ HANDLE ProcessHandle,
    _In_ PVOID StartRoutine, // PUSER_THREAD_START_ROUTINE
    _In_opt_ PVOID Argument,
    _In_ ULONG CreateFlags, // THREAD_CREATE_FLAGS_*
    _In_ SIZE_T ZeroBits,
    _In_ SIZE_T StackSize,
    _In_ SIZE_T MaximumStackSize,
    _In_opt_ PPS_ATTRIBUTE_LIST AttributeList)
{
    if (IsInsideHook())
    {
        return TrueNtCreateThreadEx(
            ThreadHandle,
            DesiredAccess,
            ObjectAttributes,
            ProcessHandle,
            StartRoutine,
            Argument,
            CreateFlags,
            ZeroBits,
            StackSize,
            MaximumStackSize,
            AttributeList);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtCreateThreadEx(DesiredAccess: %d, ProcessHandle: %p, StartRoutine: %p, CreateFlags: %lu), RETN: %p",
        DesiredAccess,
        ProcessHandle,
        StartRoutine,
        CreateFlags,
        _ReturnAddress());

    NTSTATUS Status = TrueNtCreateThreadEx(
        ThreadHandle,
        DesiredAccess,
        ObjectAttributes,
        ProcessHandle,
        StartRoutine,
        Argument,
        CreateFlags,
        ZeroBits,
        StackSize,
        MaximumStackSize,
        AttributeList);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtSuspendThread(_In_ HANDLE ThreadHandle, _Out_opt_ PULONG PreviousSuspendCount)
{
    if (IsInsideHook())
    {
        return TrueNtSuspendThread(ThreadHandle, PreviousSuspendCount);
    }

    CaptureStackTrace();

    TraceAPI(L"NtSuspendThread(ThreadHandle: %p), RETN: %p", ThreadHandle, _ReturnAddress());

    NTSTATUS Status = TrueNtSuspendThread(ThreadHandle, PreviousSuspendCount);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtResumeThread(_In_ HANDLE ThreadHandle, _Out_opt_ PULONG PreviousSuspendCount)
{
    if (IsInsideHook())
    {
        return TrueNtResumeThread(ThreadHandle, PreviousSuspendCount);
    }

    CaptureStackTrace();

    TraceAPI(L"NtResumeThread(ThreadHandle: %p), RETN: %p", ThreadHandle, _ReturnAddress());

    NTSTATUS Status = TrueNtResumeThread(ThreadHandle, PreviousSuspendCount);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtOpenProcess(
    _Out_ PHANDLE ProcessHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_ POBJECT_ATTRIBUTES ObjectAttributes,
    _In_opt_ PCLIENT_ID ClientId)
{
    if (IsInsideHook())
    {
        return TrueNtOpenProcess(ProcessHandle, DesiredAccess, ObjectAttributes, ClientId);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtOpenProcess(DesiredAccess: 0x%x, UniqueProcess:  0x%x), RETN: %p",
        DesiredAccess,
        ClientId->UniqueProcess,
        _ReturnAddress());

    NTSTATUS Status = TrueNtOpenProcess(ProcessHandle, DesiredAccess, ObjectAttributes, ClientId);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtTerminateProcess(_In_opt_ HANDLE ProcessHandle, _In_ NTSTATUS ExitStatus)
{
    if (IsInsideHook())
    {
        return TrueNtTerminateProcess(ProcessHandle, ExitStatus);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtTerminateProcess(ProcessHandle: 0x%p, ExitStatus: %d), RETN: %p",
        ProcessHandle,
        ExitStatus,
        _ReturnAddress());

    NTSTATUS Status = TrueNtTerminateProcess(ProcessHandle, ExitStatus);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS NTAPI
HookNtContinue(_In_ PCONTEXT ContextRecord, _In_ BOOLEAN TestAlert)
{
    NTSTATUS Status;

    CaptureStackTrace();

    TraceAPI(L"NtContinue(ContextRecord: 0x%p, TestAlert: %d), RETN: %p", ContextRecord, TestAlert, _ReturnAddress());

    if (bFirstTime)
    {
        HANDLE ModuleHandle = NULL;
        UNICODE_STRING ModulePath;

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
        bFirstTime = FALSE;
    }

	return TrueNtContinue(ContextRecord, TestAlert);
}
