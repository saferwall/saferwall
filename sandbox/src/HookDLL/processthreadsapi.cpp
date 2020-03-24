#include "stdafx.h"
#include "processthreadsapi.h"

decltype(NtCreateUserProcess) *TrueNtCreateUserProcess = nullptr;
decltype(NtCreateThread) *TrueNtCreateThread = nullptr;
decltype(NtCreateThreadEx) *TrueNtCreateThreadEx = nullptr;
decltype(NtResumeThread) *TrueNtResumeThread = nullptr;
decltype(NtSuspendThread) *TrueNtSuspendThread = nullptr;
decltype(NtOpenThread) *TrueNtOpenThread = nullptr;
decltype(NtOpenProcess) *TrueNtOpenProcess = nullptr;
decltype(NtTerminateProcess) *TrueNtTerminateProcess = nullptr;
decltype(NtContinue) *TrueNtContinue = nullptr;
decltype(NtGetContextThread) *TrueNtGetContextThread = nullptr;
decltype(NtSetContextThread) *TrueNtSetContextThread = nullptr;


BOOL bFirstTime = TRUE;
extern HookInfo gHookInfo;

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
    if (SfwIsCalledFromSystemMemory(5))
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
    if (SfwIsCalledFromSystemMemory(5))
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
    if (SfwIsCalledFromSystemMemory(5))
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
    if (SfwIsCalledFromSystemMemory(5))
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
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtResumeThread(ThreadHandle, PreviousSuspendCount);
    }

    CaptureStackTrace();

    TraceAPI(L"NtResumeThread(ThreadHandle: %p), RETN: %p", ThreadHandle, _ReturnAddress());

    NTSTATUS Status = TrueNtResumeThread(ThreadHandle, PreviousSuspendCount);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS
NTAPI
HookNtOpenThread(
    _Out_ PHANDLE ThreadHandle,
    _In_ ACCESS_MASK DesiredAccess,
    _In_ POBJECT_ATTRIBUTES ObjectAttributes,
    _In_opt_ PCLIENT_ID ClientId)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtOpenThread(ThreadHandle, DesiredAccess, ObjectAttributes, ClientId);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtOpenThread(DesiredAccess: 0x%x, UniqueProcess:  0x%x), RETN: %p",
        DesiredAccess,
        ClientId->UniqueProcess,
        _ReturnAddress());

    NTSTATUS Status = TrueNtOpenThread(ThreadHandle, DesiredAccess, ObjectAttributes, ClientId);

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
    if (SfwIsCalledFromSystemMemory(5))
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
    if (SfwIsCalledFromSystemMemory(5))
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

    if (bFirstTime)
    {
        HANDLE ModuleHandle = NULL;
        UNICODE_STRING ModulePath;

        RtlInitUnicodeString(&ModulePath, (PWSTR)L"ole32.dll");
        Status = LdrGetDllHandle(NULL, 0, &ModulePath, &ModuleHandle);
        if (Status == STATUS_SUCCESS && !gHookInfo.IsOleHooked)
        {
            HookOleAPIs(TRUE);
        }

        RtlInitUnicodeString(&ModulePath, (PWSTR)L"wininet.dll");
        Status = LdrGetDllHandle(NULL, 0, &ModulePath, &ModuleHandle);
        if (Status == STATUS_SUCCESS && !gHookInfo.IsWinInetHooked)
        {
            HookNetworkAPIs(TRUE);
        }

		SfwSymInit();
        bFirstTime = FALSE;

        Status = TrueNtContinue(ContextRecord, TestAlert);
        return Status;
    }

	if (SfwIsCalledFromSystemMemory(1))
    {
        Status = TrueNtContinue(ContextRecord, TestAlert);
        return Status;
    }

    CaptureStackTrace();

    TraceAPI(L"NtContinue(ContextRecord: 0x%p, TestAlert: %d), RETN: %p", ContextRecord, TestAlert, _ReturnAddress());

    Status =  TrueNtContinue(ContextRecord, TestAlert);
    return Status;
}

NTSTATUS NTAPI
HookNtGetContextThread(_In_ HANDLE ThreadHandle, _Inout_ PCONTEXT ThreadContext) {
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtGetContextThread(ThreadHandle, ThreadContext);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtGetContextThread(ThreadHandle: 0x%x, ThreadContext:  0x%p), RETN: %p",
        ThreadHandle,
        ThreadContext,
        _ReturnAddress());

    NTSTATUS Status = TrueNtGetContextThread(ThreadHandle, ThreadContext);

    ReleaseHookGuard();

    return Status;
}

NTSTATUS
NTAPI
HookNtSetContextThread(_In_ HANDLE ThreadHandle, _In_ PCONTEXT ThreadContext)
{
    if (SfwIsCalledFromSystemMemory(5))
    {
        return TrueNtSetContextThread(ThreadHandle, ThreadContext);
    }

    CaptureStackTrace();

    TraceAPI(
        L"NtSetContextThread(ThreadHandle: 0x%x, ThreadContext:  0x%p), RETN: %p",
        ThreadHandle,
        ThreadContext,
        _ReturnAddress());

    NTSTATUS Status = TrueNtSetContextThread(ThreadHandle, ThreadContext);

    ReleaseHookGuard();

    return Status;
}