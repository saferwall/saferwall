#include "stdafx.h"
#include <Windows.h>


//
// Globals
//

extern __vsnwprintf_fn_t _vsnwprintf;
extern REGHANDLE ProviderHandle;



VOID TraceAPI(PCWSTR Format, ...) {

	WCHAR Buffer[256];

	va_list arglist;
	va_start(arglist, Format);
	_vsnwprintf(Buffer, RTL_NUMBER_OF(Buffer), Format, arglist);
	EtwEventWriteString(ProviderHandle, 0, 0, Buffer);
	va_end(arglist);
}


VOID LogMessage(PCWSTR Format, ...) {
	WCHAR Buffer[256];

	va_list arglist;
	va_start(arglist, Format);
	_vsnwprintf(Buffer, RTL_NUMBER_OF(Buffer), Format, arglist);
	EtwEventWriteString(ProviderHandle, 0, 0, Buffer);
	va_end(arglist);
}