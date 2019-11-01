#pragma once

#include "stdafx.h"

// Prototypes
NTSTATUS WINAPI HookNtCreateUserProcess(
	_Out_ PHANDLE ProcessHandle,
	_Out_ PHANDLE ThreadHandle,
	_In_ ACCESS_MASK ProcessDesiredAccess,
	_In_ ACCESS_MASK ThreadDesiredAccess,
	_In_opt_ POBJECT_ATTRIBUTES ProcessObjectAttributes,
	_In_opt_ POBJECT_ATTRIBUTES ThreadObjectAttributes,
	_In_ ULONG ProcessFlags, // PROCESS_CREATE_FLAGS_*
	_In_ ULONG ThreadFlags, // THREAD_CREATE_FLAGS_*
	_In_opt_ PVOID ProcessParameters, // PRTL_USER_PROCESS_PARAMETERS
	_Inout_ PPS_CREATE_INFO CreateInfo,
	_In_opt_ PPS_ATTRIBUTE_LIST AttributeList
);

