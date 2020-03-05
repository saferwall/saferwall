#pragma once

//
// Include files.
//

#include <ntddk.h> // various NT definitions
#include "injection.h"
#include "shared.h"

//
// Defines.
//

#define NT_DEVICE_NAME L"\\Device\\SAFERWALL_SANDBOX"
#define DOS_DEVICE_NAME L"\\DosDevices\\SaferwallSandbox"

//
// Prototypes.
//

NTSTATUS
DispatchCreateClose(PDEVICE_OBJECT DeviceObject, PIRP Irp);

NTSTATUS
DispatchDeviceControl(PDEVICE_OBJECT DeviceObject, PIRP Irp);

VOID
UnloadDriver(_In_ PDRIVER_OBJECT DriverObject);

VOID NTAPI
CreateProcessNotifyRoutine(
    _Inout_ PEPROCESS Process,
    _In_ HANDLE ProcessId,
    _In_opt_ PPS_CREATE_NOTIFY_INFO CreateInfo);

VOID NTAPI
LoadImageNotifyRoutine(_In_opt_ PUNICODE_STRING FullImageName, _In_ HANDLE ProcessId, _In_ PIMAGE_INFO ImageInfo);