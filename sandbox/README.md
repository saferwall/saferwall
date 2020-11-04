# Introduction

Saferwall sandbox is a tool written in C focused primarly in analyzing malware.
The current implementation is targetting Windows OS for the momen, specifically Windows 7 x64.

## Architecture

- Driver:
    - Intercept newly created process / modules, and inject a DLL via APC..
- DLL:
	- Copy the original page of the API.
    - Performs inling hooking using Windows Detours library.
    - Send ioctl with code_page and data_page to be hooked.
- Hypervisor (modded KVM):
    - Use EPT techniques to hide user mode hooks.
    - Use EPTP switching VMFUNC to avoid traping all VM-Exit, handle them inside our driver running in the guest instead.
		1. Modify KVM to enable VE ( EPTP switching )
		2. Shadow the IDT, and provide original copy in `sidt`.
		3. Handle EPT violations inside the guest (in IDT VE exception vector)
	- Hook page:
		- allocate code_page and data_page.
		- Update EPT mapping for the 3 EPT pointers and invept.
	- IDT #VE:
		- Handle EPT violation.
		- if violation happens because of bad access and in a page we are hooking:
			- vmfunc to WR EPT if read|write violation
			- vmfunc to EXEC EPT if exec violation

## Features

- Invisible inline user hooks.
- Kernel mode malware support.
- Track child processes and follows code injection.
- Extract all files writen to disk.
- Memory dumps/unpacking.
- Fix IAT for PE dumps.
- Resistent to anti-VM detection techniques.
- User simulater running inside the guest.

## Hiding

- Driver running in the guest (registry keys, ..)
- Processes running inside the guest (User simulator)
- DLL injected.


## Hooked APIs

### Libload

- LdrLoadDll
- LdrGetProcedureAddressEx
- LdrGetDllHandleEx

### Files

- NtOpenFile
- NtCreateFile
- NtReadFile
- NtWriteFile
- NtDeleteFile
- NtSetInformationFile
- NtQueryDirectoryFile
- NtQueryInformationFile
- NtQueryFullAttributesFile (Add)
- MoveFileWithProgressTransactedW

### Memory

- NtProtectVirtualMemory
- NtQueryVirtualMemory
- NtReadVirtualMemory
- NtWriteVirtualMemory
- NtMapViewOfSection
- NtUnmapViewOfSection
- NtAllocateVirtualMemory
- NtFreeVirtualMemory

### Registry

- NtOpenKey
- NtOpenKeyEx
- NtCreateKey
- NtQueryValueKey
- NtSetValueKey
- NtDeleteKey
- NtDeleteValueKey


### Process/Threads

- NtOpenProcess
- NtTerminateProcess
- NtCreateUserProcess
- NtCreateThread
- NtCreateThreadEx
- NtSuspendThread
- NtResumeThread

### Network

- InternetOpenA
- InternetConnectA
- InternetConnectW
- HttpOpenRequestA
- HttpOpenRequestW
- HttpSendRequestA
- HttpSendRequestW
- InternetReadFile

### Service (to add)

- OpenSCManagerW
- CreateServiceW
- OpenServiceW
- StartServiceW
- ControlService
- DeleteService
- EnumServicesStatusW

### Crypto

- RtlDecompressBuffer

### Synchronization

- NtDelayExecution

### System

- NtQuerySystemInformation
- NtLoadDriver
- NtQueryVolumeInformationFile
- NtDeviceIoControlFile

## Apps running inside the VM

- Visual C++ Redistributable Package from 2005 -> 2019.
- Microsoft Office

## VM Hardening

- Disabled Windows Defender
- Disabled Windows Update
- Disabled Windows Firewall.
- Disable windows security center notifications.
- Disable Turning off the display
- Prevent Hard Drives from Going to Sleep.
- Change Power Settings to High Performance
- Turn Off Search Indexing
- Disable unnecessary services