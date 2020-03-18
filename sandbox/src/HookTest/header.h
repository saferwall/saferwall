#pragma once

#include <Windows.h>
#include <stdlib.h>
#include <time.h>

//
// Strings, Path & IO
//

#include <stdio.h>
#include <shlwapi.h>
#pragma comment(lib, "shlwapi.lib")

//
// Network
//

#include <wininet.h>
#pragma comment(lib, "Wininet.lib")

// COM
#include <comdef.h>
#include <Wbemidl.h>
#pragma comment(lib, "wbemuuid.lib")

//
// Defines.
//

#define TEST_FILE_HOOKS TRUE
#define TEST_LIB_LOAD_HOOKS TRUE
#define TEST_MEMORY_HOOKS TRUE
#define TEST_NETWORK_HOOKS TRUE
#define TEST_OLE_HOOKS TRUE
#define TEST_PROCESS_THREADS_HOOKS TRUE
#define TEST_REGISTRY_HOOKS TRUE
#define TEST_SYNC_HOOKS TRUE

//
// Prototypes
//

VOID
GetRandomString(PWCHAR Str, CONST INT Len);
VOID
GetRandomDir(PWSTR szPathOut);
VOID
GetRandomFilePath(PWSTR szPathOut);
VOID
ErrorExit(const char *wszProcedureName);
VOID
TestFileHooks();
VOID
TestLibLoadHooks();
VOID
TestMemoryHooks();
VOID
TestNetworkHooks();
VOID
TestOleHooks();
VOID
TestRegistryHooks();