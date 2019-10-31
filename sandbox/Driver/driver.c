#include "driver.h"

NTSTATUS
DriverEntry(
	_In_ PDRIVER_OBJECT   DriverObject,
	_In_ PUNICODE_STRING  RegistryPath
)
/*++

Routine Description:
	This routine is called by the Operating System to initialize the driver.

	It creates the device object, fills in the dispatch entry points and
	completes the initialization.

Arguments:
	DriverObject - a pointer to the object that represents this device
	driver.

	RegistryPath - a pointer to our Services key in the registry.

Return Value:
	STATUS_SUCCESS if initialized; an error otherwise.

--*/

{
	NTSTATUS        ntStatus;
	PDEVICE_OBJECT  deviceObject = NULL;	// ptr to device object
	UNICODE_STRING  ntUnicodeString;		// NT Device Name "\Device\SIOCTL"
	UNICODE_STRING  ntWin32NameString;		// Win32 Name "\DosDevices\IoctlTest"

	RtlInitUnicodeString(&ntUnicodeString, NT_DEVICE_NAME);

	UNREFERENCED_PARAMETER(RegistryPath);

	//
	// Creates a device object for use by our driver.
	// 

	ntStatus = IoCreateDevice(
		DriverObject,                   // Our Driver Object
		0,                              // We don't use a device extension
		&ntUnicodeString,               // Device name "\Device\SIOCTL"
		FILE_DEVICE_UNKNOWN,            // Device type
		FILE_DEVICE_SECURE_OPEN,        // Device characteristics
		FALSE,                          // Not an exclusive device
		&deviceObject);                 // Returned ptr to Device Object

	if (!NT_SUCCESS(ntStatus))
	{
		LOG(("Couldn't create the device object\n"));
		return ntStatus;
	}


	//
	// Initialize the driver object with this driver's entry points.
	//

	DriverObject->MajorFunction[IRP_MJ_CREATE] = DeviceCreateClose;
	DriverObject->MajorFunction[IRP_MJ_CLOSE] = DeviceCreateClose;
	DriverObject->MajorFunction[IRP_MJ_DEVICE_CONTROL] = IoctlDeviceControl;
	DriverObject->DriverUnload = UnloadDriver;

	//
	// Initialize a Unicode String containing the Win32 name
	// for our device.
	//

	RtlInitUnicodeString(&ntWin32NameString, DOS_DEVICE_NAME);

	//
	// Create a symbolic link between our device name  and the Win32 name
	//

	ntStatus = IoCreateSymbolicLink(
		&ntWin32NameString, &ntUnicodeString);

	if (!NT_SUCCESS(ntStatus))
	{
		//
		// Delete everything that this routine has allocated.
		//
		LOG(("Couldn't create symbolic link\n"));
		IoDeleteDevice(deviceObject);
	}


	return ntStatus;



}


VOID
UnloadDriver(
	_In_ PDRIVER_OBJECT DriverObject
)
/*++

Routine Description:

	This routine is called by the I/O system to unload the driver.

	Any resources previously allocated must be freed.

Arguments:

	DriverObject - a pointer to the object that represents our driver.

Return Value:

	None
--*/

{
	PDEVICE_OBJECT deviceObject = DriverObject->DeviceObject;
	UNICODE_STRING uniWin32NameString;

	PAGED_CODE();

	//
	// Create counted string version of our Win32 device name.
	//

	RtlInitUnicodeString(&uniWin32NameString, DOS_DEVICE_NAME);


	//
	// Delete the link from our device name to a name in the Win32 namespace.
	//

	IoDeleteSymbolicLink(&uniWin32NameString);

	if (deviceObject != NULL)
	{
		IoDeleteDevice(deviceObject);
	}

}



NTSTATUS
DeviceCreateClose(
	PDEVICE_OBJECT DeviceObject,
	PIRP Irp
)
/*++

Routine Description:

	This routine is called by the I/O system when the our device is opened or
	closed.

	No action is performed other than completing the request successfully.

Arguments:

	DeviceObject - a pointer to the object that represents the device
	that I/O is to be done on.

	Irp - a pointer to the I/O Request Packet for this request.

Return Value:

	NT status code

--*/

{
	UNREFERENCED_PARAMETER(DeviceObject);

	PAGED_CODE();

	Irp->IoStatus.Status = STATUS_SUCCESS;
	Irp->IoStatus.Information = 0;

	IoCompleteRequest(Irp, IO_NO_INCREMENT);

	return STATUS_SUCCESS;
}


NTSTATUS
IoctlDeviceControl(
	PDEVICE_OBJECT DeviceObject,
	PIRP Irp
)

/*++

Routine Description:

	This routine is called by the I/O system to perform a device I/O
	control function.

Arguments:

	DeviceObject - a pointer to the object that represents the device
		that I/O is to be done on.

	Irp - a pointer to the I/O Request Packet for this request.

Return Value:

	NT status code

--*/

{
	PIO_STACK_LOCATION  irpSp;// Pointer to current stack location
	NTSTATUS            ntStatus = STATUS_SUCCESS;// Assume success
	ULONG               inBufLength; // Input buffer length
	ULONG               outBufLength; // Output buffer length

	UNREFERENCED_PARAMETER(DeviceObject);

	PAGED_CODE();

	irpSp = IoGetCurrentIrpStackLocation(Irp);
	inBufLength = irpSp->Parameters.DeviceIoControl.InputBufferLength;
	outBufLength = irpSp->Parameters.DeviceIoControl.OutputBufferLength;

	if (!inBufLength || !outBufLength)
	{
		ntStatus = STATUS_INVALID_PARAMETER;
		goto End;
	}

	//
	// Determine which I/O control code was specified.
	//

	switch (irpSp->Parameters.DeviceIoControl.IoControlCode)
	{
	}

End:
	//
	// Finish the I/O operation by simply completing the packet and returning
	// the same status as in the packet itself.
	//

	Irp->IoStatus.Status = ntStatus;

	IoCompleteRequest(Irp, IO_NO_INCREMENT);

	return ntStatus;
}