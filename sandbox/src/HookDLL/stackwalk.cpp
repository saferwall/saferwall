#include "stdafx.h"

CRITICAL_SECTION gDbgHelpLock;

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