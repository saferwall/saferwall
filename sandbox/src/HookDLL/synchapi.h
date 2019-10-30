#pragma once

#include "stdafx.h"

// 
// Prototypes 
// 


NTSTATUS NTAPI HookNtDelayExecution
(
	_In_ BOOLEAN Alertable,
	_In_opt_ PLARGE_INTEGER DelayInterval
);
