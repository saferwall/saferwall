#include "header.h"

VOID
ErrorExit(const char *wszProcedureName)
{
    WCHAR wszMsgBuff[512]; // Buffer for text.
    DWORD dwChars;         // Number of chars returned.

    // Get the last error.
    DWORD dwErr = GetLastError();

    // Try to get the message from the system errors.
    dwChars = FormatMessage(
        FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS, NULL, dwErr, 0, wszMsgBuff, 512, NULL);

    if (0 == dwChars)
    {
        // The error code did not exist in the system errors.
        // Try Ntdsbmsg.dll for the error code.

        HINSTANCE hInst;

        // Load the library.
        hInst = LoadLibrary(L"Ntdsbmsg.dll");
        if (NULL == hInst)
        {
            printf("cannot load Ntdsbmsg.dll\n");
            exit(1); // Could 'return' instead of 'exit'.
        }

        // Try getting message text from ntdsbmsg.
        dwChars = FormatMessage(
            FORMAT_MESSAGE_FROM_HMODULE | FORMAT_MESSAGE_IGNORE_INSERTS, hInst, dwErr, 0, wszMsgBuff, 512, NULL);

        // Free the library.
        FreeLibrary(hInst);
    }

    // Display the error message, or generic text if not found.
    printf(
        "Function %s failed Error value: %d Message: %ws\n",
        wszProcedureName,
        dwErr,
        dwChars ? wszMsgBuff : L"Error message not found.");
    exit(1); // Could 'return' instead of 'exit'.
}

VOID
GetRandomString(PWCHAR Str, CONST INT Len)
{
    static CONST WCHAR AlphaNum[] = L"0123456789"
                                    "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
                                    "abcdefghijklmnopqrstuvwxyz";

    for (INT i = 0; i < Len; ++i)
    {
        Str[i] = AlphaNum[rand() % (wcslen(AlphaNum) - 1)];
    }

    Str[Len] = 0;
}

VOID
GetRandomDir(PWSTR szPathOut)
{
    DWORD Count;
    WCHAR TempPath[MAX_PATH] = L"";
    WCHAR RandomName[MAX_PATH] = L"";

    Count = GetTempPath(MAX_PATH, TempPath);
    if (!Count)
    {
        ErrorExit("GetTempPath");
    }
    GetRandomString(RandomName, 8);
    PathCombine(szPathOut, TempPath, RandomName);
}

VOID
GetRandomFilePath(PWSTR szPathOut)
{
    WCHAR szFilePath[MAX_PATH] = L"";
    GetRandomDir(szFilePath);

    CreateDirectory(szFilePath, NULL);

    WCHAR szFileName[MAX_PATH] = L"";
    GetRandomString(szFileName, 8);
    wcscat_s(szFileName, MAX_PATH, L".txt");

    PathCombine(szPathOut, szFilePath, szFileName);
}