#pragma once

#include "stdafx.h"

// 
// Prototypes 
// 

NTSTATUS NTAPI HookNtQuerySystemInformation
(
	_In_ SYSTEM_INFORMATION_CLASS SystemInformationClass,
	_Out_writes_bytes_opt_(SystemInformationLength) PVOID SystemInformation,
	_In_ ULONG SystemInformationLength,
	_Out_opt_ PULONG ReturnLength
);

