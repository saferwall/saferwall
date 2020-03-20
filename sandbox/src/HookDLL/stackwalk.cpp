#include "stdafx.h"

#define SIZE_SYMBOL sizeof(SYMBOL_INFO) + MAX_SYM_NAME * sizeof(WCHAR)

extern DWORD dwTlsIndex;

CRITICAL_SECTION gDbgHelpLock;
PVOID gSymbolBuffer;
PVOID gFrames[MAX_FRAME];

VOID
AllocateSpaceSymbol()
{
    //
    // Allocate a buffer large enough to hold the symbol information on the stack and get
    // a pointer to the buffer. We also have to set the size of the symbol structure itself
    // and the number of bytes reserved for the name.
    //

    gSymbolBuffer = RtlAllocateHeap(RtlProcessHeap(), 0, SIZE_SYMBOL);
}

VOID
CaptureStackTrace()
{
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

        Context.ContextFlags = CONTEXT_CONTROL | CONTEXT_INTEGER;

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

    RtlZeroMemory(gSymbolBuffer, SIZE_SYMBOL);
    PSYMBOL_INFO pSymbol = (PSYMBOL_INFO)gSymbolBuffer;
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
                &Context,
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

VOID
SfwCaptureStackFrames(PSTACKTRACE StackTrace, UINT MaxFrames)
{
    PCONTEXT InitialContext = NULL;
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

        Context.ContextFlags = CONTEXT_CONTROL | CONTEXT_INTEGER;

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

    RtlZeroMemory(gSymbolBuffer, SIZE_SYMBOL);
    RtlZeroMemory(gFrames, MAX_FRAME * sizeof(PVOID));

    PSYMBOL_INFO pSymbol = (PSYMBOL_INFO)gSymbolBuffer;
    pSymbol->SizeOfStruct = sizeof(SYMBOL_INFO);
    pSymbol->MaxNameLen = MAX_SYM_NAME;

    //
    // Dbghelp is is singlethreaded, so acquire a lock.
    //
    // Note that the code assumes that
    // SymInitialize( GetCurrentProcess(), NULL, TRUE ) has
    // already been called.
    //
    EnterCriticalSection(&gDbgHelpLock);

    while (StackTrace->FrameCount < MaxFrames)
    {
        if (!StackWalk64(
                MachineType,
                NtCurrentProcess(),
                NtCurrentThread(),
                &StackFrame,
                &Context,
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
            StackTrace->Frames[StackTrace->FrameCount++] = StackFrame.AddrPC.Offset;
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
SfwIsCalledFromSystemMemory(UINT FramesToCapture)
{
    //
    // Are we called from inside a our own hook handler.
    //
    if (IsInsideHook())
    {
        return TRUE;
    }

    //
    // Capture up to 5 stack frames from the current call stack.
    // We're going to skip the first two stack frame returned
    // because that's the SfwIsCalledFromSystemMemory() and the
    // Hook Handler function itself, which we don't care about.
    //

	STACKTRACE StackTrace = {0};
    FramesToCapture += 4;
    SfwCaptureStackFrames(&StackTrace, FramesToCapture);

    //
    // Get the PEB.
    //
#if defined(_WIN64)
    PPEB pPeb = (PPEB)__readgsqword(0x60);

#elif defined(_WIN32)
    PPEB pPeb = (PPEB)__readfsdword(0x30);
#endif

	BOOL bFound = FALSE;
	PPEB_LDR_DATA pLdrData = NULL;
    PLIST_ENTRY pEntry, pHeadEntry = NULL;
    PLDR_DATA_TABLE_ENTRY pLdrEntry = NULL;
    pLdrData = pPeb->Ldr;

    //
    // Iterate over our stack frame and check 
    // if our address belongs to a subsystem DLL.
	// We skip the first 3 entries:
	// Hook Handler -> SfwIsCalledFromSystemMemory -> SfwCaptureStackFrame.
    //
    for (ULONG i = 3; i < StackTrace.FrameCount; i++)
    {
        bFound = FALSE;
        pHeadEntry = &pLdrData->InMemoryOrderModuleList;
        pEntry = pHeadEntry->Flink;

        while (pEntry != pHeadEntry)
        {
            // Retrieve the current LDR_DATA_TABLE_ENTRY
            pLdrEntry = CONTAINING_RECORD(pEntry, LDR_DATA_TABLE_ENTRY, InMemoryOrderLinks);

            // Exluce the main module code in the search.
            if (wcscmp(pLdrEntry->FullDllName.Buffer, pPeb->ProcessParameters->ImagePathName.Buffer) == 0)
            {
                pEntry = pEntry->Flink;
                continue;
            }

            // Fill the MODULE_ENTRY with the LDR_DATA_TABLE_ENTRY information
            if (StackTrace.Frames[i] >= (ULONGLONG)pLdrEntry->DllBase &&
                StackTrace.Frames[i] <= (ULONGLONG)pLdrEntry->DllBase + pLdrEntry->SizeOfImage)
            {
                bFound = TRUE;
                break;
            }

            // Iterate to the next entry.
            pEntry = pEntry->Flink;
        }

        if (bFound)
        {
            continue;
        }
        else
        {
            return FALSE;
        }
    }
    ReleaseHookGuard();
    return TRUE;
}

NTSTATUS
SfwSymInit()
{
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

    DWORD SymOptions = SymGetOptions();
    SymOptions |= SYMOPT_LOAD_LINES;
    SymOptions |= SYMOPT_FAIL_CRITICAL_ERRORS;
    SymOptions = SymSetOptions(SymOptions);

    return STATUS_SUCCESS;
}
