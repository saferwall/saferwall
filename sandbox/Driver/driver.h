#pragma once

//
// Include files.
//

#include <ntddk.h>          // various NT definitions



//
// Defines.
//

#define NT_DEVICE_NAME      L"\\Device\\SAFERWALL_SANDBOX"
#define DOS_DEVICE_NAME     L"\\DosDevices\\SaferwallSandbox"

#if DBG
#define LOG(_x_) \
                DbgPrint("SIOCTL.SYS: ");\
                DbgPrint _x_;

#else
#define SIOCTL_KDPRINT(_x_)
#endif



//
// Prototypes.
//

NTSTATUS
DeviceCreateClose(
	PDEVICE_OBJECT DeviceObject,
	PIRP Irp
);

NTSTATUS
IoctlDeviceControl(
	PDEVICE_OBJECT DeviceObject,
	PIRP Irp
);


VOID
UnloadDriver(
	_In_ PDRIVER_OBJECT DriverObject
);