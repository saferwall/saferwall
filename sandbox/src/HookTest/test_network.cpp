#include "header.h"


VOID
TestNetworkHooks()
{
    wprintf(L" ========= Testing network opeations ========= \n\n");

	wprintf(L"[+] Calling InternetOpenW\n");
    HINTERNET hSession = InternetOpen(
        L"Mozilla/5.0", // User-Agent
        INTERNET_OPEN_TYPE_PRECONFIG,
        NULL,
        NULL,
        0);

	wprintf(L"[+] Calling InternetConnectW\n");
    HINTERNET hConnect = InternetConnect(
        hSession,
        L"www.google.com", // HOST
        0,
        L"",
        L"",
        INTERNET_SERVICE_HTTP,
        0,
        0);

	wprintf(L"[+] Calling HttpOpenRequestW\n");
    HINTERNET hHttpFile = HttpOpenRequest(
        hConnect,
        L"GET", // METHOD
        L"/",   // URI
        NULL,
        NULL,
        NULL,
        0,
        0);

	wprintf(L"[+] Calling HttpSendRequestW\n");
    BOOL Success = HttpSendRequest(hHttpFile, NULL, 0, 0, 0);
    if (Success == FALSE)
    {
        wprintf(L"Failed Request\n");
    }

    HttpSendRequestA(hHttpFile, NULL, 0, 0, 0);

    CHAR szBuffer[4096 * 10] = "";
    DWORD dwRead = 0;
    while (InternetReadFile(hHttpFile, szBuffer, sizeof(szBuffer) - 1, &dwRead) && dwRead)
    {
        szBuffer[dwRead] = 0;
        dwRead = 0;
    }
    InternetCloseHandle(hHttpFile);
    InternetCloseHandle(hConnect);
    InternetCloseHandle(hSession);
}