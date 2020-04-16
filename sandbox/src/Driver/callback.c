#include "injection.h"
#include "shared.h"
#include "monitor.h"
#include <ntddk.h>

extern BOOLEAN g_InjIsWindows7;
extern LIST_ENTRY InjInfoListHead;

VOID
CreateProcessNotifyRoutine(_Inout_ PEPROCESS Process, _In_ HANDLE ProcessId, _In_opt_ PPS_CREATE_NOTIFY_INFO CreateInfo)
{
    if (CreateInfo != NULL)
    {
        LOG_INFO(
            "CreateProcessNotifyRoutine: process %p (ID 0x%p) created\n"
            "    creator %Ix:%Ix, command line %wZ, file name %wZ (FileOpenNameAvailable: %d)",
            Process,
            (PVOID)ProcessId,
            (ULONG_PTR)CreateInfo->CreatingThreadId.UniqueProcess,
            (ULONG_PTR)CreateInfo->CreatingThreadId.UniqueThread,
            CreateInfo->CommandLine,
            CreateInfo->ImageFileName,
            CreateInfo->FileOpenNameAvailable);

        UNICODE_STRING ProcessNameToWatch = RTL_CONSTANT_STRING(L"wrar59b3_1831105618.exe");
        if (!RtlxSuffixUnicodeString(&ProcessNameToWatch, (PUNICODE_STRING)CreateInfo->ImageFileName, TRUE))
        {
            return;
        }

        PINJECTION_INFO CapturedInjectionInfo;

        CapturedInjectionInfo = ExAllocatePoolWithTag(NonPagedPoolNx, sizeof(INJECTION_INFO), INJ_MEMORY_TAG);
        if (!CapturedInjectionInfo)
        {
            return;
            // return STATUS_INSUFFICIENT_RESOURCES;
        }

        RtlZeroMemory(CapturedInjectionInfo, sizeof(INJECTION_INFO));

        CapturedInjectionInfo->ProcessId = ProcessId;
        CapturedInjectionInfo->ForceUserApc = TRUE;

        InsertTailList(&InjInfoListHead, &CapturedInjectionInfo->ListEntry);
    }
    else
    {
        LOG_INFO("CreateProcessNotifyRoutine: process %p (ID 0x%p) terminated", Process, (PVOID)ProcessId);
    }
}

VOID NTAPI
LoadImageNotifyRoutine(_In_opt_ PUNICODE_STRING FullImageName, _In_ HANDLE ProcessId, _In_ PIMAGE_INFO ImageInfo)
{
    UNREFERENCED_PARAMETER(ProcessId);
    UNREFERENCED_PARAMETER(ImageInfo);

    //
    // Check if current process is injected.
    //

    PINJECTION_INFO InjectionInfo = InjFindInjectionInfo(ProcessId);

    if (!InjectionInfo || InjectionInfo->IsInjected)
    {
        return;
    }

    LOG_INFO(
        "LoadImageNotifyRoutine: Process: %s, FullImageName %wZ",
        PsGetProcessImageFileName(PsGetCurrentProcess()),
        FullImageName);

    if (PsIsProtectedProcess(PsGetCurrentProcess()))
    {
        //
        // Protected processes throw code-integrity error when
        // they are injected.  Signing policy can be changed, but
        // it requires hacking with lots of internal and Windows-
        // version-specific structures.  Simly don't inject such
        // processes.
        //
        // See Blackbone project (https://github.com/DarthTon/Blackbone)
        // if you're interested how protection can be temporarily
        // disabled on such processes.  (Look for BBSetProtection).
        //

        LOG_WARN(
            "Ignoring protected process (PID: %u, Name: '%s')",
            (ULONG)(ULONG_PTR)ProcessId,
            PsGetProcessImageFileName(PsGetCurrentProcess()));

        InjRemoveInjectionInfoByProcessId(ProcessId, TRUE);

        return;
    }

    if (!InjCanInject(InjectionInfo))
    {
        //
        // This process is in early stage - important DLLs (such as
        // ntdll.dll - or wow64.dll in case of Wow64 process) aren't
        // properly initialized yet.  We can't inject the DLL until
        // they are.
        //
        // Check if any of the system DLLs we're interested in is being
        // currently loaded - if so, mark that information down into the
        // LoadedDlls field.
        //

        InjIsWantedSystemDllBeingLoaded(InjectionInfo, FullImageName, ImageInfo);
    }
    else
    {
#if defined(INJ_CONFIG_SUPPORTS_WOW64)
        if (g_InjIsWindows7 && PsGetProcessWow64Process(PsGetCurrentProcess()))
        {
            //
            // On Windows 7, if we're injecting DLL into Wow64 process using
            // the "thunk method", we have additionaly postpone the load after
            // these system DLLs.
            //
            // This is because on Windows 7, these DLLs are loaded as part of
            // the wow64!ProcessInit routine, therefore the Wow64 subsystem
            // is not fully initialized to execute our injected Wow64ApcRoutine.
            //

            UNICODE_STRING System32Kernel32Path = RTL_CONSTANT_STRING(L"\\System32\\kernel32.dll");
            UNICODE_STRING SysWOW64Kernel32Path = RTL_CONSTANT_STRING(L"\\SysWOW64\\kernel32.dll");
            UNICODE_STRING System32User32Path = RTL_CONSTANT_STRING(L"\\System32\\user32.dll");
            UNICODE_STRING SysWOW64User32Path = RTL_CONSTANT_STRING(L"\\SysWOW64\\user32.dll");

            if (RtlxSuffixUnicodeString(&System32Kernel32Path, FullImageName, TRUE) ||
                RtlxSuffixUnicodeString(&SysWOW64Kernel32Path, FullImageName, TRUE) ||
                RtlxSuffixUnicodeString(&System32User32Path, FullImageName, TRUE) ||
                RtlxSuffixUnicodeString(&SysWOW64User32Path, FullImageName, TRUE))
            {
                LOG_INFO("Postponing injection (%wZ)", FullImageName);
                return;
            }
        }
#endif

        //
        // All necessary DLLs are loaded - perform the injection.
        //
        // Note that injection is done via kernel-mode APC, because
        // InjInject calls ZwMapViewOfSection and MapViewOfSection
        // might be already on the callstack.  Because MapViewOfSection
        // locks the EPROCESS->AddressCreationLock, we would be risking
        // deadlock by calling InjInject directly.
        //

#if defined(INJ_CONFIG_SUPPORTS_WOW64)
        LOG_INFO(
            "Injecting (PID: %u, Wow64: %s, Name: '%s')\n",
            (ULONG)(ULONG_PTR)ProcessId,
            PsGetProcessWow64Process(PsGetCurrentProcess()) ? "TRUE" : "FALSE",
            PsGetProcessImageFileName(PsGetCurrentProcess()));
#else
        LOG_INFO(
            "injecting (PID: %u, Name: '%s')\n",
            (ULONG)(ULONG_PTR)ProcessId,
            PsGetProcessImageFileName(PsGetCurrentProcess()));

#endif

        InjQueueApc(KernelMode, &InjpInjectApcNormalRoutine, InjectionInfo, NULL, NULL);

        //
        // Mark that this process is injected.
        //

        InjectionInfo->IsInjected = TRUE;

		//
		// Add to the list of monitored processes.
		//

		MonAddProcessToMonitoredList((ULONG)(ULONG_PTR)ProcessId);
    }
}

VOID NTAPI
CreateThreadNotifyRoutine(HANDLE ProcessId, HANDLE ThreadId, BOOLEAN Create)
{
    if (Create)
    {

		HANDLE CurrentProcessId = PsGetCurrentProcessId();

		if (CurrentProcessId != ProcessId)
        {
             LOG_INFO(
                "CreateThreadNotifyRoutine: (ProcessId 0x%p, ThreadId 0x%p) created", (PVOID)ProcessId,
                (PVOID)ThreadId);
            LOG_INFO("Thread injection");
		}

	}
}