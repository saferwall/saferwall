#pragma once

#include "stdafx.h"

// Prototypes
NTSTATUS WINAPI HookLdrLoadDll(
	PWSTR  DllPath,
	PULONG  DllCharacteristics,
	PUNICODE_STRING DllName,
	PVOID * DllHandle
);


NTSTATUS WINAPI HookLdrGetProcedureAddress(
	PVOID DllHandle,
	PANSI_STRING ProcedureName,
	ULONG ProcedureNumber,
	PVOID *ProcedureAddress
);

