// HookDLL.cpp : Defines the exported functions for the DLL application.
//

#include "stdafx.h"
#include "hooking.h"

//
// Defines
//

#define DETOUR_END

//
// Lib Load APIs.
//

//
// Globals
//
extern decltype(LdrLoadDll) *TrueLdrLoadDll;
extern decltype(LdrGetProcedureAddressEx) *TrueLdrGetProcedureAddressEx;
extern decltype(LdrGetDllHandleEx) *TrueLdrGetDllHandleEx;

extern decltype(NtCreateFile) *TrueNtCreateFile;
extern decltype(NtReadFile) *TrueNtReadFile;
extern decltype(NtWriteFile) *TrueNtWriteFile;
extern decltype(NtDeleteFile) *TrueNtDeleteFile;
extern decltype(NtSetInformationFile) *TrueNtSetInformationFile;
extern decltype(NtQueryDirectoryFile) *TrueNtQueryDirectoryFile;
extern decltype(NtQueryInformationFile) *TrueNtQueryInformationFile;

extern decltype(NtProtectVirtualMemory) *TrueNtProtectVirtualMemory;
extern decltype(NtQueryVirtualMemory) *TrueNtQueryVirtualMemory;
extern decltype(NtReadVirtualMemory) *TrueNtReadVirtualMemory;
extern decltype(NtWriteVirtualMemory) *TrueNtWriteVirtualMemory;
extern decltype(NtMapViewOfSection) *TrueNtMapViewOfSection;
extern decltype(NtUnmapViewOfSection) *TrueNtUnmapViewOfSection;
extern decltype(NtAllocateVirtualMemory) *TrueNtAllocateVirtualMemory;
extern decltype(NtFreeVirtualMemory) *TrueNtFreeVirtualMemory;

extern decltype(NtOpenKey) *TrueNtOpenKey;
extern decltype(NtOpenKeyEx) *TrueNtOpenKeyEx;
extern decltype(NtCreateKey) *TrueNtCreateKey;
extern decltype(NtQueryValueKey) *TrueNtQueryValueKey;
extern decltype(NtSetValueKey) *TrueNtSetValueKey;
extern decltype(NtDeleteKey) *TrueNtDeleteKey;
extern decltype(NtDeleteValueKey) *TrueNtDeleteValueKey;

extern decltype(NtOpenProcess) *TrueNtOpenProcess;
extern decltype(NtTerminateProcess) *TrueNtTerminateProcess;
extern decltype(NtCreateUserProcess) *TrueNtCreateUserProcess;
extern decltype(NtCreateThread) *TrueNtCreateThread;
extern decltype(NtCreateThreadEx) *TrueNtCreateThreadEx;
extern decltype(NtSuspendThread) *TrueNtSuspendThread;
extern decltype(NtResumeThread) *TrueNtResumeThread;
extern decltype(NtContinue) *TrueNtContinue;

extern decltype(NtQuerySystemInformation) *TrueNtQuerySystemInformation;
extern decltype(RtlDecompressBuffer) *TrueRtlDecompressBuffer;
extern decltype(NtDelayExecution) *TrueNtDelayExecution;
extern decltype(NtLoadDriver) *TrueNtLoadDriver;


__vsnwprintf_fn_t _vsnwprintf = nullptr;
__snwprintf_fn_t _snwprintf = nullptr;
strlen_fn_t _strlen = nullptr;
pfn_wcsstr _wcsstr = nullptr;
pfnStringFromGUID2 _StringFromGUID2 = nullptr;
pfnStringFromCLSID _StringFromCLSID = nullptr;
pfnCoTaskMemFree _CoTaskMemFree = nullptr;
pfnCoCreateInstanceEx TrueCoCreateInstanceEx = nullptr;
pfnInternetOpenA TrueInternetOpenA = nullptr;
pfnInternetConnectA TrueInternetConnectA = nullptr;
pfnInternetConnectW TrueInternetConnectW = nullptr;
pfnHttpOpenRequestA TrueHttpOpenRequestA = nullptr;
pfnHttpOpenRequestW TrueHttpOpenRequestW = nullptr;

CRITICAL_SECTION gDbgHelpLock, gInsideHookLock;
BOOL gInsideHook = FALSE;
DWORD dwTlsIndex;

//
// ETW provider GUID and global provider handle.
// GUID:
//   {a4b4ba50-a667-43f5-919b-1e52a6d69bd5}
//

GUID ProviderGuid = {0xa4b4ba50, 0xa667, 0x43f5, {0x91, 0x9b, 0x1e, 0x52, 0xa6, 0xd6, 0x9b, 0xd5}};

REGHANDLE ProviderHandle;
#define ATTACH(x) DetAttach(&(PVOID &)True##x, Hook##x, #x)
#define DETACH(x) DetDetach(&(PVOID &)True##x, Hook##x, #x)


typedef struct _STACKTRACE
{
    //
    // Number of frames in Frames array.
    //
    UINT FrameCount;

    //
    // PC-Addresses of frames. Index 0 contains the topmost frame.
    //
    ULONGLONG Frames[ANYSIZE_ARRAY];
} STACKTRACE, *PSTACKTRACE;


VOID
CaptureStackTrace()
{
    return;
    PCONTEXT InitialContext = NULL;
    // STACKTRACE StackTrace;
    UINT MaxFrames = 50;
    STACKFRAME64 StackFrame;
    DWORD MachineType = 0;
    CONTEXT Context = {};

    if (InitialContext == NULL)
    {
        //
        // Use current context.
        //
        // N.B. GetThreadContext cannot be used on the current thread.
        // Capture own context - on i386, there is no API for that.
        //
#ifdef _M_IX86
        ZeroMemory(&Context, sizeof(CONTEXT));

        Context.ContextFlags = CONTEXT_CONTROL;

        //
        // Those three registers are enough.
        //
        __asm {
		Label:
			mov[Context.Ebp], ebp;
			mov[Context.Esp], esp;
			mov eax, [Label];
			mov[Context.Eip], eax;
        }
#else
        RtlCaptureContext(&Context);
#endif
    }
    else
    {
        CopyMemory(&Context, InitialContext, sizeof(CONTEXT));
    }
    //
    // Set up stack frame.
    //
    ZeroMemory(&StackFrame, sizeof(STACKFRAME64));

#ifdef _M_IX86
    MachineType = IMAGE_FILE_MACHINE_I386;
    StackFrame.AddrPC.Offset = Context.Eip;
    StackFrame.AddrPC.Mode = AddrModeFlat;
    StackFrame.AddrFrame.Offset = Context.Ebp;
    StackFrame.AddrFrame.Mode = AddrModeFlat;
    StackFrame.AddrStack.Offset = Context.Esp;
    StackFrame.AddrStack.Mode = AddrModeFlat;
#elif _M_X64
    MachineType = IMAGE_FILE_MACHINE_AMD64;
    StackFrame.AddrPC.Offset = Context.Rip;
    StackFrame.AddrPC.Mode = AddrModeFlat;
    StackFrame.AddrFrame.Offset = Context.Rsp;
    StackFrame.AddrFrame.Mode = AddrModeFlat;
    StackFrame.AddrStack.Offset = Context.Rsp;
    StackFrame.AddrStack.Mode = AddrModeFlat;
#elif _M_IA64
    MachineType = IMAGE_FILE_MACHINE_IA64;
    StackFrame.AddrPC.Offset = Context.StIIP;
    StackFrame.AddrPC.Mode = AddrModeFlat;
    StackFrame.AddrFrame.Offset = Context.IntSp;
    StackFrame.AddrFrame.Mode = AddrModeFlat;
    StackFrame.AddrBStore.Offset = Context.RsBSP;
    StackFrame.AddrBStore.Mode = AddrModeFlat;
    StackFrame.AddrStack.Offset = Context.IntSp;
    StackFrame.AddrStack.Mode = AddrModeFlat;
#else
#    error "Unsupported platform"
#endif

    //
    // Allocate a buffer large enough to hold the symbol information on the stack and get
    // a pointer to the buffer. We also have to set the size of the symbol structure itself
    // and the number of bytes reserved for the name.
    //

    char buffer[sizeof(SYMBOL_INFO) + MAX_SYM_NAME * sizeof(WCHAR)];
    PSYMBOL_INFO pSymbol = (PSYMBOL_INFO)buffer;
    pSymbol->SizeOfStruct = sizeof(SYMBOL_INFO);
    pSymbol->MaxNameLen = MAX_SYM_NAME;

    DWORD64 dwDisplacement = 0;

    UINT FrameCount = 0;
    WCHAR pszFilename[MAX_PATH + 1];
    WCHAR *ModuleName = NULL;
    DWORD dwResult = 0;

    //
    // Dbghelp is is singlethreaded, so acquire a lock.
    //
    // Note that the code assumes that
    // SymInitialize( GetCurrentProcess(), NULL, TRUE ) has
    // already been called.
    //
    EnterCriticalSection(&gDbgHelpLock);

    while (FrameCount < MaxFrames)
    {
        if (!StackWalk64(
                MachineType,
                NtCurrentProcess(),
                NtCurrentThread(),
                &StackFrame,
                MachineType == IMAGE_FILE_MACHINE_I386 ? NULL : &Context,
                NULL,
                SymFunctionTableAccess64,
                SymGetModuleBase64,
                NULL))
        {
            //
            // Maybe it failed, maybe we have finished walking the stack.
            //
            break;
        }

        if (StackFrame.AddrPC.Offset != 0)
        {
            //
            // Valid frame.
            //
            // StackTrace->Frames[StackTrace->FrameCount++] = StackFrame.AddrPC.Offset;
            FrameCount++;

            dwResult = GetMappedFileNameW(NtCurrentProcess(), (LPVOID)StackFrame.AddrPC.Offset, pszFilename, MAX_PATH);
            if (dwResult)
            {
                ModuleName = (WCHAR *)FindFileName(pszFilename);
            }
            else
            {
                ModuleName = (WCHAR *)L"N/A";
            }

            //
            // Retrieves symbol information for the specified address.
            //
            if (SymFromAddr(NtCurrentProcess(), StackFrame.AddrPC.Offset, &dwDisplacement, pSymbol))
            {
                LogMessage(
                    L"Module: %s, SymbolName:%ws, SymbolAddress:0x%08llx, Offset:0x%p",
                    ModuleName,
                    MultiByteToWide(pSymbol->Name),
                    pSymbol->Address,
                    StackFrame.AddrPC.Offset);
            }
            else
            {
                LogMessage(
                    L"Module: %s, SymbolName:N/A, SymbolAddress: N/A, Offset:0x%p",
                    ModuleName,
                    StackFrame.AddrPC.Offset);
            }
        }
        else
        {
            //
            // Base reached.
            //
            break;
        }
    }

    LeaveCriticalSection(&gDbgHelpLock);
}

BOOL
IsInsideHook()
/*++

Routine Description:

    This function checks if are already inside a hook handler.
    This helps avoid infinite recursions which happens in hooking
    as some APIs inside the hook handler end up calling functions
    which are detoured as well.

    There are some few issues you have to be concerned about
    if you are injecting a 64bits DLL inside a WoW64 process.
        1.  Implicit TLS (__declspec(thread)) relies heavily on the
            CRT, which is not available to us.
        2.  Explicit TLS APIs (TlsAlloc() / TlsFree(), etc.) are
            implemented entirely in kernel32.dll, whose 64-bit
            version is not loaded into WoW64 processes.

    In our case, we alweays injects DLL of the same architecture
    as the process. So it should be safe to use TLS. The TLS
    allocation should happen before attacking the hooks as TlsAlloc
    end up calling RtlAllocateHeap() which might be hooked as well.

Return Value:
    TRUE: if we are inside a hook handler.
    FALSE: otherwise.
--*/
{
    if (!TlsGetValue(dwTlsIndex))
    {
        TlsSetValue(dwTlsIndex, (LPVOID)TRUE);
        return FALSE;
    }
    return TRUE;
}

VOID
ReleaseHookGuard()
{
    TlsSetValue(dwTlsIndex, (LPVOID)FALSE);
}

LONG
CheckDetourAttach(LONG err)
{
    switch (err)
    {
    case ERROR_INVALID_BLOCK: /*printf("ERROR_INVALID_BLOCK: The function referenced is too small to be detoured.");*/
        break;
    case ERROR_INVALID_HANDLE: /*printf("ERROR_INVALID_HANDLE: The ppPointer parameter is null or points to a null
                                  pointer.");*/
        break;
    case ERROR_INVALID_OPERATION: /*	printf("ERROR_INVALID_OPERATION: No pending transaction exists."); */
        break;
    case ERROR_NOT_ENOUGH_MEMORY: /*printf("ERROR_NOT_ENOUGH_MEMORY: Not enough memory exists to complete the
                                     operation.");*/
        break;
    case NO_ERROR:
        break;
    default: /*printf("CheckDetourAttach failed with unknown error code.");*/
        break;
    }
    return err;
}

static const char *
DetRealName(const char *psz)
{
    const char *pszBeg = psz;
    // Move to end of name.
    while (*psz)
    {
        psz++;
    }
    // Move back through A-Za-z0-9 names.
    while (psz > pszBeg && ((psz[-1] >= 'A' && psz[-1] <= 'Z') || (psz[-1] >= 'a' && psz[-1] <= 'z') ||
                            (psz[-1] >= '0' && psz[-1] <= '9')))
    {
        psz--;
    }
    return psz;
}

VOID
DetAttach(PVOID *ppvReal, PVOID pvMine, PCCH psz)
{
    PVOID pvReal = NULL;
    if (ppvReal == NULL)
    {
        ppvReal = &pvReal;
    }

    LONG l = DetourAttach(ppvReal, pvMine);
    if (l != NO_ERROR)
    {
        WCHAR Buffer[128];
        _snwprintf(
            Buffer,
            RTL_NUMBER_OF(Buffer),
            L"Detour Attach failed: %ws: error %d",
            MultiByteToWide((CHAR *)DetRealName(psz)),
            l);
        EtwEventWriteString(ProviderHandle, 0, 0, Buffer);
        // Decode((PBYTE)*ppvReal, 3);
    }
}

VOID
DetDetach(PVOID *ppvReal, PVOID pvMine, PCCH psz)
{
    LONG l = DetourDetach(ppvReal, pvMine);
    if (l != NO_ERROR)
    {
        WCHAR Buffer[128];
        _snwprintf(
            Buffer,
            RTL_NUMBER_OF(Buffer),
            L"Detour Detach failed: %s: error %d",
            MultiByteToWide((CHAR *)DetRealName(psz)),
            l);
        EtwEventWriteString(ProviderHandle, 0, 0, Buffer);
    }
}

PVOID
GetAPIAddress(PSTR FunctionName, PWSTR ModuleName)
{
    NTSTATUS Status;

    ANSI_STRING RoutineName;
    RtlInitAnsiString(&RoutineName, FunctionName);

    UNICODE_STRING ModulePath;
    RtlInitUnicodeString(&ModulePath, ModuleName);

    HANDLE ModuleHandle = NULL;
    Status = LdrGetDllHandle(NULL, 0, &ModulePath, &ModuleHandle);
    if (Status != STATUS_SUCCESS)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"LdrGetDllHandle failed");
        return NULL;
    }

    PVOID Address;
    Status = LdrGetProcedureAddress(ModuleHandle, &RoutineName, 0, &Address);
    if (Status != STATUS_SUCCESS)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"LdrGetProcedureAddress failed");
        return NULL;
    }

    return Address;
}

BOOL
ProcessAttach()
{
    //
    // Register ETW provider.
    //

    EtwEventRegister(&ProviderGuid, NULL, NULL, &ProviderHandle);

    //
    // Set up the symbol options so that we can gather information from the current
    // executable's PDB files, as well as the Microsoft symbol servers.  We also want
    // to undecorate the symbol names we're returned.  If you want, you can add other
    // symbol servers or paths via a semi-colon separated list in SymInitialized.
    //

    if (!SymInitialize(NtCurrentProcess(), NULL, TRUE))
    {
        LogMessage(L"SymInitialize returned error : %d", GetLastError());
        return STATUS_UNSUCCESSFUL;
    }
    SymSetOptions(SYMOPT_DEFERRED_LOADS | SYMOPT_INCLUDE_32BIT_MODULES | SYMOPT_UNDNAME | SYMOPT_LOAD_LINES);

    //
    // Allocate a TLS index.
    //

    if ((dwTlsIndex = TlsAlloc()) == TLS_OUT_OF_INDEXES)
        LogMessage(L"TlsAlloc failed");

    //
    // Save real API addresses.
    //

    TrueLdrLoadDll = LdrLoadDll;
    TrueLdrGetProcedureAddressEx = LdrGetProcedureAddressEx;
    TrueLdrGetDllHandleEx = LdrGetDllHandleEx;
    TrueNtCreateFile = NtCreateFile;
    TrueNtWriteFile = NtWriteFile;
    TrueNtReadFile = NtReadFile;
    TrueNtDeleteFile = NtDeleteFile;
    TrueNtSetInformationFile = NtSetInformationFile;
    TrueNtQueryDirectoryFile = NtQueryDirectoryFile;
    TrueNtQueryInformationFile = NtQueryInformationFile;
    TrueNtDelayExecution = NtDelayExecution;
    TrueNtProtectVirtualMemory = NtProtectVirtualMemory;
    TrueNtQueryVirtualMemory = NtQueryVirtualMemory;
    TrueNtReadVirtualMemory = NtReadVirtualMemory;
    TrueNtWriteVirtualMemory = NtWriteVirtualMemory;
    TrueNtFreeVirtualMemory = NtFreeVirtualMemory;
    TrueNtMapViewOfSection = NtMapViewOfSection;
    TrueNtUnmapViewOfSection = NtUnmapViewOfSection;
    TrueNtAllocateVirtualMemory = NtAllocateVirtualMemory;
    TrueNtProtectVirtualMemory = NtProtectVirtualMemory;
    TrueNtOpenKey = NtOpenKey;
    TrueNtOpenKeyEx = NtOpenKeyEx;
    TrueNtCreateKey = NtCreateKey;
    TrueNtQueryValueKey = NtQueryValueKey;
    TrueNtSetValueKey = NtSetValueKey;
    TrueNtDeleteKey = NtDeleteKey;
    TrueNtDeleteValueKey = NtDeleteValueKey;
    TrueNtCreateUserProcess = NtCreateUserProcess;
    TrueNtCreateThread = NtCreateThread;
    TrueNtCreateThreadEx = NtCreateThreadEx;
    TrueNtResumeThread = NtResumeThread;
    TrueNtSuspendThread = NtSuspendThread;
    TrueNtOpenProcess = NtOpenProcess;
    TrueNtTerminateProcess = NtTerminateProcess;
    TrueNtContinue = NtContinue;
    TrueRtlDecompressBuffer = RtlDecompressBuffer;
    TrueNtQuerySystemInformation = NtQuerySystemInformation;
    TrueNtLoadDriver = NtLoadDriver;

    //
    // Resolve the ones not exposed by ntdll.
    //

    _vsnwprintf = (__vsnwprintf_fn_t)GetAPIAddress((PSTR) "_vsnwprintf", (PWSTR)L"ntdll.dll");
    if (_vsnwprintf == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"_vsnwprintf() is NULL");
    }
    _snwprintf = (__snwprintf_fn_t)GetAPIAddress((PSTR) "_snwprintf", (PWSTR)L"ntdll.dll");
    if (_vsnwprintf == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"_snwprintf() is NULL");
    }

    _wcsstr = (pfn_wcsstr)GetAPIAddress((PSTR) "wcsstr", (PWSTR)L"ntdll.dll");
    if (_wcsstr == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"wcsstr() is NULL");
    }

    //
    // Initializes a critical section objects.
    // Uesd for capturing stack trace and IsInsideHook.
    //

    InitializeCriticalSection(&gDbgHelpLock);
    InitializeCriticalSection(&gInsideHookLock);

    //
    // Attach Native APIs.
    //

    HookNtAPIs();

    return TRUE;
}

BOOL
ProcessDetach()
{
    HookBegingTransation();

    //DETACH(LdrLoadDll);
    //DETACH(LdrGetProcedureAddressEx);
    //DETACH(LdrGetDllHandleEx);

    //
    // File APIs.
    //

   /* DETACH(NtCreateFile);
    DETACH(NtReadFile);
    DETACH(NtWriteFile);
    DETACH(NtDeleteFile);
    DETACH(NtSetInformationFile);
    DETACH(NtQueryDirectoryFile);
    DETACH(NtQueryInformationFile);*/
    // DETACH(MoveFileWithProgressTransactedW);

    //
    // Registry APIs.
    //

    //DETACH(NtOpenKey);
    //DETACH(NtOpenKeyEx);
    //DETACH(NtCreateKey);
    //DETACH(NtQueryValueKey);
    //DETACH(NtSetValueKey);
    //DETACH(NtDeleteKey);
    //DETACH(NtDeleteValueKey);

    //
    // Process/Thread APIs.
    //

    //DETACH(NtOpenProcess);
    //DETACH(NtTerminateProcess);
    //DETACH(NtCreateUserProcess);
    //DETACH(NtCreateThread);
    //DETACH(NtCreateThreadEx);
    //DETACH(NtSuspendThread);
    //DETACH(NtResumeThread);
    // DETACH(NtContinue);


    //
    // System APIs.
    //

    //DETACH(NtQuerySystemInformation);
    //DETACH(RtlDecompressBuffer);
    DETACH(NtDelayExecution);
    //DETACH(NtLoadDriver);

    //
    // Memory APIs.
    //

    //DETACH(NtQueryVirtualMemory);
    //DETACH(NtReadVirtualMemory);
    //DETACH(NtWriteVirtualMemory);
    //DETACH(NtFreeVirtualMemory);
    //DETACH(NtMapViewOfSection);
    //// DETACH(NtAllocateVirtualMemory);
    //DETACH(NtUnmapViewOfSection);
    //DETACH(NtProtectVirtualMemory);

    HookCommitTransaction();

    TlsFree(dwTlsIndex);
    SymCleanup(NtCurrentProcess());
    EtwEventUnregister(ProviderHandle);

    EtwEventWriteString(ProviderHandle, 0, 0, L"Detached success");

    return STATUS_SUCCESS;
}

BOOL
HookBegingTransation()
{
    LONG Status;

    //
    // Begin a new transaction for attaching detours.
    //

	Status = DetourTransactionBegin();
    if (Status != NO_ERROR)
    {
        LogMessage(L"DetourTransactionBegin() failed with %d", Status);
        return FALSE;
    }

    //
    // Enlist a thread for update in the current transaction.
    //

	Status = DetourUpdateThread(NtCurrentThread());
    if (Status != NO_ERROR)
    {
        LogMessage(L"DetourUpdateThread() failed with %d", Status);
        return FALSE;
    }

    return TRUE;
}

BOOL
HookCommitTransaction()
{
    /*
    Commit the current transaction.
	*/

    PVOID *ppbFailedPointer = NULL;

    LONG error = DetourTransactionCommitEx(&ppbFailedPointer);
    if (error != NO_ERROR)
    {
        LogMessage(
            L"Attach transaction failed to commit. Error %d (%p/%p)", error, ppbFailedPointer, *ppbFailedPointer);
        return FALSE;
    }

    EtwEventWriteString(ProviderHandle, 0, 0, L"Detours Attached");
    return TRUE;
}

VOID
HookOleAPIs(BOOL Attach)
{
	_StringFromCLSID = (pfnStringFromCLSID)GetAPIAddress((PSTR) "StringFromCLSID", (PWSTR)L"ole32.dll");
    if (_StringFromCLSID == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"StringFromCLSID() is NULL");
    }

	_CoTaskMemFree = (pfnCoTaskMemFree)GetAPIAddress((PSTR) "CoTaskMemFree", (PWSTR)L"ole32.dll");
    if (_CoTaskMemFree == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"CoTaskMemFree() is NULL");
    }

    TrueCoCreateInstanceEx = (pfnCoCreateInstanceEx)GetAPIAddress((PSTR) "CoCreateInstanceEx", (PWSTR)L"ole32.dll");
    if (TrueCoCreateInstanceEx == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"CoCreateInstanceEx() is NULL");
    }

    HookBegingTransation();

    if (Attach)
    {
        ATTACH(CoCreateInstanceEx);
        
    }
    else
    {
        DETACH(CoCreateInstanceEx);
    }

    HookCommitTransaction();
}

VOID
HookNetworkAPIs(BOOL Attach)
{
    TrueInternetOpenA = (pfnInternetOpenA)GetAPIAddress((PSTR) "InternetOpenA", (PWSTR)L"wininet.dll");
    if (TrueInternetOpenA == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"InternetOpenA() is NULL");
    }

	TrueInternetConnectA = (pfnInternetConnectA)GetAPIAddress((PSTR) "InternetConnectA", (PWSTR)L"wininet.dll");
    if (TrueInternetOpenA == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"InternetConnectA() is NULL");
    }

	TrueInternetConnectW = (pfnInternetConnectW)GetAPIAddress((PSTR) "InternetConnectW", (PWSTR)L"wininet.dll");
    if (TrueInternetOpenA == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"InternetConnectW() is NULL");
    }

	TrueHttpOpenRequestA = (pfnHttpOpenRequestA)GetAPIAddress((PSTR) "HttpOpenRequestA", (PWSTR)L"wininet.dll");
    if (TrueHttpOpenRequestA == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"HttpOpenRequestA() is NULL");
    }

	TrueHttpOpenRequestW = (pfnHttpOpenRequestW)GetAPIAddress((PSTR) "HttpOpenRequestW", (PWSTR)L"wininet.dll");
    if (TrueHttpOpenRequestW == NULL)
    {
        EtwEventWriteString(ProviderHandle, 0, 0, L"HttpOpenRequestW() is NULL");
    }

    HookBegingTransation();

    if (Attach)
    {
        ATTACH(InternetOpenA);
        ATTACH(InternetConnectA);
        ATTACH(InternetConnectW);
        ATTACH(HttpOpenRequestA);
        ATTACH(HttpOpenRequestW);
    }
    else
    {
        DETACH(InternetOpenA);
        DETACH(InternetConnectA);
        DETACH(InternetConnectW);
        DETACH(HttpOpenRequestA);
        DETACH(HttpOpenRequestW);
    }

    HookCommitTransaction();
}

VOID
HookNtAPIs()
{
    LogMessage(L"HookNtAPIs Begin");

    HookBegingTransation();

    //
    // Lib Load APIs.
    //

    //ATTACH(LdrLoadDll);
    //ATTACH(LdrGetProcedureAddressEx);
    //ATTACH(LdrGetDllHandleEx);

    //
    // File APIs.
    //

   /* ATTACH(NtCreateFile);
    ATTACH(NtReadFile);
    ATTACH(NtWriteFile);
    ATTACH(NtDeleteFile);
    ATTACH(NtSetInformationFile);
    ATTACH(NtQueryDirectoryFile);
    ATTACH(NtQueryInformationFile);
*/
    //
    // Registry APIs.
    //

    //ATTACH(NtOpenKey);
    //ATTACH(NtOpenKeyEx);
    //ATTACH(NtCreateKey);
    //ATTACH(NtQueryValueKey);
    //ATTACH(NtSetValueKey);
    //ATTACH(NtDeleteKey);
    //ATTACH(NtDeleteValueKey);

    //
    // Process/Thread APIs.
    //

    //ATTACH(NtOpenProcess);
    //ATTACH(NtTerminateProcess);
    //ATTACH(NtCreateUserProcess);
    //ATTACH(NtCreateThread);
    //ATTACH(NtCreateThreadEx);
    //ATTACH(NtSuspendThread);
    //ATTACH(NtResumeThread);
    //ATTACH(NtContinue);

    //
    // System APIs.
    //

    //ATTACH(NtQuerySystemInformation);
    //ATTACH(RtlDecompressBuffer);
    ATTACH(NtDelayExecution);
    //ATTACH(NtLoadDriver);

    //
    // Memory APIs.
    //

  /*  ATTACH(NtQueryVirtualMemory);
    ATTACH(NtReadVirtualMemory);
    ATTACH(NtWriteVirtualMemory);
    ATTACH(NtFreeVirtualMemory);
    ATTACH(NtMapViewOfSection);*/
    // ATTACH(NtAllocateVirtualMemory);
    //ATTACH(NtUnmapViewOfSection);
    //ATTACH(NtProtectVirtualMemory);

    HookCommitTransaction();

LogMessage(L"HookNtAPIs End");
}