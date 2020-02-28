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
