#include "header.h"

#define SVCNAME L"saferwall"

VOID
TestWinSvcHooks()
{
	wprintf(L" ========= Testing winsvc opeations ========= \n\n");

	wprintf(L"[+] Calling CreateServiceW\n");

	SC_HANDLE schSCManager, schService;
	WCHAR szPath[MAX_PATH];

	if (!GetModuleFileName(NULL, szPath, MAX_PATH))
	{
		printf("Cannot install service (%d)\n", GetLastError());
		return;
	}

	// Get a handle to the SCM database. 
	schSCManager = OpenSCManager(
		NULL,                    // local computer
		NULL,                    // ServicesActive database 
		SC_MANAGER_ALL_ACCESS);  // full access rights 

	if (NULL == schSCManager)
	{
		PrintError("OpenSCManager");
		return;
	}

	// Check if the service it is already exists.
	schService = OpenService(schSCManager, SVCNAME, SERVICE_ALL_ACCESS);
	if (NULL != schService) {
		DeleteService(schService);
	}

	// Create the service
	schService = CreateService(
		schSCManager,              // SCM database 
		SVCNAME,                   // name of service 
		SVCNAME,                   // service name to display 
		SERVICE_ALL_ACCESS,        // desired access 
		SERVICE_WIN32_OWN_PROCESS, // service type 
		SERVICE_DEMAND_START,      // start type 
		SERVICE_ERROR_NORMAL,      // error control type 
		szPath,                    // path to service's binary 
		NULL,                      // no load ordering group 
		NULL,                      // no tag identifier 
		NULL,                      // no dependencies 
		NULL,                      // LocalSystem account 
		NULL);                     // no password 

	if (schService == NULL)
	{
		CloseServiceHandle(schSCManager);
		PrintError("CreateService");
	}
	else printf("Service installed successfully\n");

	CloseServiceHandle(schService);
	CloseServiceHandle(schSCManager);
	
}