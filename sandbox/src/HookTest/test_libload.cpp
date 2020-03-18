#include "header.h"

VOID
TestLibLoadHooks()
{
    HMODULE hModule;

    wprintf(L" ========= Testing libload opeations ========= \n\n");

    wprintf(L"[+] Calling LoadLibraryA\n");
    LoadLibraryA("advapi32.dll");

    wprintf(L"[+] Calling LoadLibraryW\n");
    LoadLibraryW(L"advapi32.dll");

    wprintf(L"[+] Calling LoadLibraryExA\n");
    LoadLibraryExA("kernelbase.dll", NULL, 0);

    wprintf(L"[+] Calling LoadLibraryExW\n");
    hModule = LoadLibraryExW(L"kernel32.dll", NULL, 0);

    wprintf(L"[+] Calling GetProcAddress\n");
    GetProcAddress(hModule, "WriteFile");
}