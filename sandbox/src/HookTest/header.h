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
// WinCrypt
//
#include <Wincrypt.h>

//
// Defines.
//

#define TEST_FILE_HOOKS FALSE
#define TEST_LIB_LOAD_HOOKS FALSE
#define TEST_MEMORY_HOOKS FALSE
#define TEST_NETWORK_HOOKS FALSE
#define TEST_OLE_HOOKS FALSE
#define TEST_PROCESS_THREADS_HOOKS FALSE
#define TEST_REGISTRY_HOOKS FALSE
#define TEST_SYNC_HOOKS FALSE
#define TEST_WINSVC_HOOKS TRUE
#define TEST_WINCRYPT_HOOKS TRUE

//
// Prototypes
//

VOID
GetRandomString(PWCHAR Str, CONST INT Len);
VOID
GetRandomDir(PWSTR szPathOut);
VOID
GetRandomFilePath(PWSTR szPathOut);
DWORD
PrintError(const char *wszProcedureName);
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
VOID
TestWinSvcHooks();
VOID
TestWinCryptHooks();