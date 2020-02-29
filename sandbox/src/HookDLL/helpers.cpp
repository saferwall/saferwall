#include "stdafx.h"

LPCWSTR FindFileName(LPCWSTR pPath)
{
	LPCWSTR pT = NULL;
	if (!pPath) {
		return NULL;
	}

	for (pT = pPath; *pPath; pPath++) {
		if ((pPath[0] == '\\' || pPath[0] == ':' || pPath[0] == '/')
			&& pPath[1] && pPath[1] != '\\' && pPath[1] != '/')
			pT = pPath + 1;
	}

	return pT;
}


WCHAR* MultiByteToWide(CHAR* lpMultiByteStr)
{
	//int Size = MultiByteToWideChar(CP_ACP, MB_ERR_INVALID_CHARS, szSource, strlen(szSource), NULL, 0);
	//WCHAR *wszDest = reinterpret_cast<WCHAR*>(RtlAllocateHeap(RtlProcessHeap(), 0, Size));
	//SecureZeroMemory(wszDest, Size);
	//MultiByteToWideChar(CP_ACP, MB_PRECOMPOSED, szSource, strlen(szSource), wszDest, Size);

		/* Get the required size */
	size_t iNumChars = strlen(lpMultiByteStr);

	/* Allocate new wide string */
	SIZE_T Size = (1 + iNumChars) * sizeof(WCHAR);

	WCHAR *lpWideCharStr = reinterpret_cast<WCHAR*>(RtlAllocateHeap(RtlProcessHeap(), 0, Size));
	WCHAR *It;
	It = lpWideCharStr;
	if (lpWideCharStr) {
		SecureZeroMemory(lpWideCharStr, Size);
		while (iNumChars) {

			*lpWideCharStr = *lpMultiByteStr;
			lpWideCharStr++;
			lpMultiByteStr++;
			iNumChars--;
		}

	}
	return It;

	//return wszDest;
}
