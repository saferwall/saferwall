# Introduction

Saferwall's sandbox is a tool written in C focused primarly in analyzing malware.

## Motivation

The current implementation is targetting Windows OS for the moment.

## Architecture

- Driver:
    - Intercept newly created process / modules, and inject a DLL via APC.
- DLL:
    - Performs inling hooking using Windows Detours library.
- Hypervisor:
    - KVM hiding user mode hooks using EPT shadowing.
    - Usage of EPTP switching VMFUNC to avoid traping all VM-Exit to the host but handle them inside our driving running in the guest.

1. Modify KVM to enable VE ( EPTP swicthing )
2. Create EPTs indentity mapping (All access, RW, Execute only)
3. IDT needs to be shadowed
4. Driver inside the guest hook IDT VE exception, handle EPT inside the guest.

Before the DLL do the inline hook, you need to make a copy of the page.

- On process start:
	- Wait untill requires DLLs are loaded.
	- Inject DLL via APC
	- Before you inline hook, you need 
	- From Kernel, get VA for
	- VMCALL (code_page, data_page)


## Features

- Invisible hooks and resistent to anti-sandbox tricks.
- Track child processes and follows code injection.
- User simulater running inside the guest.
- Extract all files writen to disk.
- Unpack the file and fix its IAT.


## Hooked APIs

### Libload

- LdrLoadDll
- LdrGetProcedureAddressEx

### Files

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

### Crypto

- RtlDecompressBuffer

### Synchronization

- NtDelayExecution

### System

- NtQuerySystemInformation
- NtLoadDriver


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