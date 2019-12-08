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
	BOOLEAN SymLinkCreated = FALSE;
	BOOLEAN CreateProcessCallbackCreate = FALSE;

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
		LOG_ERROR("Couldn't create the device object");
		goto Exit;
	}


	//
	// Initialize the driver object with this driver's entry points.
	//

	DriverObject->MajorFunction[IRP_MJ_CREATE] = DispatchCreateClose;
	DriverObject->MajorFunction[IRP_MJ_CLOSE] = DispatchCreateClose;
	DriverObject->MajorFunction[IRP_MJ_DEVICE_CONTROL] = DispatchDeviceControl;
	DriverObject->DriverUnload = UnloadDriver;

	//
	// Initialize a Unicode String containing the Win32 name
	// for our device.
	//

	RtlInitUnicodeString(&ntWin32NameString, DOS_DEVICE_NAME);

	//
	// Create a symbolic link between our device name  and the Win32 name
	//

	ntStatus = IoCreateSymbolicLink(&ntWin32NameString, &ntUnicodeString);

	if (!NT_SUCCESS(ntStatus))
	{
		LOG_ERROR("Couldn't create symbolic link\n");
		goto Exit;
		
	}
	SymLinkCreated = TRUE;


	//
	// Initialize injection.
	//
	__debugbreak();
	InjInitialize(RegistryPath);

	//
	// Registers a process notification callback that notifies the us
	// when a process is created or exits.
	//

	ntStatus = PsSetCreateProcessNotifyRoutineEx(CreateProcessNotifyRoutine, FALSE);
	if (!NT_SUCCESS(ntStatus))
	{
		if (ntStatus == STATUS_ACCESS_DENIED) {
			// Project Properties -> Linker -> All Options then add /INTEGRITYCHECK
			// See: https://www.osronline.com/showthread.cfm?link=169632
			// See: https://msdn.microsoft.com/en-us/library/dn195769.aspx?f=255&MSPPError=-2147217396
			LOG_ERROR("PsSetCreateProcessNotifyRoutineEx() failed; ensure /INTEGRITYCHECK linker flag was used during linking");
		}
		else {
			LOG_ERROR("Unable to add process creation notification routine");
		}


		//DbgPrintEx(DPFLTR_IHVDRIVER_ID, DPFLTR_ERROR_LEVEL, "ObCallbackTest: DriverEntry: PsSetCreateProcessNotifyRoutineEx(2) returned 0x%x\n", ntStatus);
		goto Exit;
	}
	CreateProcessCallbackCreate = TRUE;

	//
	// Registers an image load notification callback.
	//

	ntStatus = PsSetLoadImageNotifyRoutine(&LoadImageNotifyRoutine);

	if (!NT_SUCCESS(ntStatus))
	{
		LOG_ERROR("Unable to add image load notification routine");
		goto Exit;
	}


Exit:

	if (!NT_SUCCESS(ntStatus))
	{
		//
		// Delete everything that this routine has allocated.
		//

		ntStatus = PsRemoveLoadImageNotifyRoutine(&LoadImageNotifyRoutine);
		_ASSERT(ntStatus == STATUS_SUCCESS);


		if (CreateProcessCallbackCreate == TRUE)
		{
			ntStatus = PsSetCreateProcessNotifyRoutineEx(CreateProcessNotifyRoutine, TRUE);
			_ASSERT(ntStatus == STATUS_SUCCESS);
			CreateProcessCallbackCreate = FALSE;
		}

		if (SymLinkCreated == TRUE)
		{
			IoDeleteSymbolicLink(&ntWin32NameString);
		}

		if (deviceObject != NULL)
		{
			IoDeleteDevice(deviceObject);
		}
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

	//
	// Removes the LoadImageNotifyRoutine.
	//
	PsRemoveLoadImageNotifyRoutine(&LoadImageNotifyRoutine);

	//
	// Removes the CreateProcessNotifyRoutine.
	//
	PsSetCreateProcessNotifyRoutineEx(&CreateProcessNotifyRoutine, TRUE);

	//
	// Release memory of all injection-info entries.
	//
	InjDestroy();

}



NTSTATUS
DispatchCreateClose(
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
DispatchDeviceControl(
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

