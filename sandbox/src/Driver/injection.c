#include "injection.h"

//
// DOS Device Prefix \??\
//

ALIGNEDNAME ObpDosDevicesShortNamePrefix = {{L'\\', L'?', L'?', L'\\'}};
UNICODE_STRING ObpDosDevicesShortName = {
    sizeof(ObpDosDevicesShortNamePrefix), // Length
    sizeof(ObpDosDevicesShortNamePrefix), // MaximumLength
    (PWSTR)&ObpDosDevicesShortNamePrefix  // Buffer
};

//////////////////////////////////////////////////////////////////////////
// Private constant variables.
//////////////////////////////////////////////////////////////////////////

ANSI_STRING LdrLoadDllRoutineName = RTL_CONSTANT_STRING("LdrLoadDll");

UCHAR InjpThunkX86[] = {
    //
    0x83, 0xec, 0x08,             // sub    esp,0x8
    0x0f, 0xb7, 0x44, 0x24, 0x14, // movzx  eax,[esp + 0x14]
    0x66, 0x89, 0x04, 0x24,       // mov    [esp],ax
    0x66, 0x89, 0x44, 0x24, 0x02, // mov    [esp + 0x2],ax
    0x8b, 0x44, 0x24, 0x10,       // mov    eax,[esp + 0x10]
    0x89, 0x44, 0x24, 0x04,       // mov    [esp + 0x4],eax
    0x8d, 0x44, 0x24, 0x14,       // lea    eax,[esp + 0x14]
    0x50,                         // push   eax
    0x8d, 0x44, 0x24, 0x04,       // lea    eax,[esp + 0x4]
    0x50,                         // push   eax
    0x6a, 0x00,                   // push   0x0
    0x6a, 0x00,                   // push   0x0
    0xff, 0x54, 0x24, 0x1c,       // call   [esp + 0x1c]
    0x83, 0xc4, 0x08,             // add    esp,0x8
    0xc2, 0x0c, 0x00,             // ret    0xc
};                                //

UCHAR InjpThunkX64[] = {
    //
    0x48, 0x83, 0xec, 0x38,             // sub    rsp,0x38
    0x48, 0x89, 0xc8,                   // mov    rax,rcx
    0x66, 0x44, 0x89, 0x44, 0x24, 0x20, // mov    [rsp+0x20],r8w
    0x66, 0x44, 0x89, 0x44, 0x24, 0x22, // mov    [rsp+0x22],r8w
    0x4c, 0x8d, 0x4c, 0x24, 0x40,       // lea    r9,[rsp+0x40]
    0x48, 0x89, 0x54, 0x24, 0x28,       // mov    [rsp+0x28],rdx
    0x4c, 0x8d, 0x44, 0x24, 0x20,       // lea    r8,[rsp+0x20]
    0x31, 0xd2,                         // xor    edx,edx
    0x31, 0xc9,                         // xor    ecx,ecx
    0xff, 0xd0,                         // call   rax
    0x48, 0x83, 0xc4, 0x38,             // add    rsp,0x38
    0xc2, 0x00, 0x00,                   // ret    0x0
};                                      //

UCHAR InjpThunkARM32[] = {
    //
    0x1f, 0xb5,             // push   {r0-r4,lr}
    0xad, 0xf8, 0x08, 0x20, // strh   r2,[sp,#8]
    0xad, 0xf8, 0x0a, 0x20, // strh   r2,[sp,#0xA]
    0x03, 0x91,             // str    r1,[sp,#0xC]
    0x02, 0xaa,             // add    r2,sp,#8
    0x00, 0x21,             // movs   r1,#0
    0x04, 0x46,             // mov    r4,r0
    0x6b, 0x46,             // mov    r3,sp
    0x00, 0x20,             // movs   r0,#0
    0xa0, 0x47,             // blx    r4
    0x1f, 0xbd,             // pop    {r0-r4,pc}
};                          //

UCHAR InjpThunkARM64[] = {
    //
    0xfe, 0x0f, 0x1f, 0xf8, // str    lr,[sp,#-0x10]!
    0xff, 0x83, 0x00, 0xd1, // sub    sp,sp,#0x20
    0xe9, 0x03, 0x00, 0xaa, // mov    x9,x0
    0xe2, 0x13, 0x00, 0x79, // strh   w2,[sp,#8]
    0x00, 0x00, 0x80, 0xd2, // mov    x0,#0
    0xe2, 0x17, 0x00, 0x79, // strh   w2,[sp,#0xA]
    0xe2, 0x23, 0x00, 0x91, // add    x2,sp,#8
    0xe1, 0x0b, 0x00, 0xf9, // str    x1,[sp,#0x10]
    0x01, 0x00, 0x80, 0xd2, // mov    x1,#0
    0xe3, 0x03, 0x00, 0x91, // mov    x3,sp
    0x20, 0x01, 0x3f, 0xd6, // blr    x9
    0xff, 0x83, 0x00, 0x91, // add    sp,sp,#0x20
    0xfe, 0x07, 0x41, 0xf8, // ldr    lr,[sp],#0x10
    0xc0, 0x03, 0x5f, 0xd6, // ret
};                          //
//
// Paths can have format "\Device\HarddiskVolume3\Windows\System32\ntdll.dll",
// so only the end of the string is compared.
//

INJ_SYSTEM_DLL_DESCRIPTOR InjpSystemDlls[] = {
    {RTL_CONSTANT_STRING(L"\\SysArm32\\ntdll.dll"), INJ_SYSARM32_NTDLL_LOADED},
    {RTL_CONSTANT_STRING(L"\\SyChpe32\\ntdll.dll"), INJ_SYCHPE32_NTDLL_LOADED},
    {RTL_CONSTANT_STRING(L"\\SysWow64\\ntdll.dll"), INJ_SYSWOW64_NTDLL_LOADED},
    {RTL_CONSTANT_STRING(L"\\System32\\ntdll.dll"), INJ_SYSTEM32_NTDLL_LOADED},
    {RTL_CONSTANT_STRING(L"\\System32\\wow64.dll"), INJ_SYSTEM32_WOW64_LOADED},
    {RTL_CONSTANT_STRING(L"\\System32\\wow64win.dll"), INJ_SYSTEM32_WOW64WIN_LOADED},
    {RTL_CONSTANT_STRING(L"\\System32\\wow64cpu.dll"), INJ_SYSTEM32_WOW64CPU_LOADED},
    {RTL_CONSTANT_STRING(L"\\System32\\wowarmhw.dll"), INJ_SYSTEM32_WOWARMHW_LOADED},
    {RTL_CONSTANT_STRING(L"\\System32\\xtajit.dll"), INJ_SYSTEM32_XTAJIT_LOADED},
};

//////////////////////////////////////////////////////////////////////////
// Variables.
//////////////////////////////////////////////////////////////////////////

LIST_ENTRY InjInfoListHead;
BOOLEAN g_InjIsWindows7;
UNICODE_STRING InjDllPath[InjArchitectureMax];

INJ_THUNK InjThunk[InjArchitectureMax] = {
    {InjpThunkX86, sizeof(InjpThunkX86)},
    {InjpThunkX64, sizeof(InjpThunkX64)},
    {InjpThunkARM32, sizeof(InjpThunkARM32)},
    {InjpThunkARM64, sizeof(InjpThunkARM64)},
};

//////////////////////////////////////////////////////////////////////////
// Helper functions.
//////////////////////////////////////////////////////////////////////////

PVOID
NTAPI
RtlxFindExportedRoutineByName(_In_ PVOID DllBase, _In_ PANSI_STRING ExportName)
{
    //
    // RtlFindExportedRoutineByName is not exported by ntoskrnl until Win10.
    // Following code is borrowed from ReactOS.
    //

    PULONG NameTable;
    PUSHORT OrdinalTable;
    PIMAGE_EXPORT_DIRECTORY ExportDirectory;
    LONG Low = 0, Mid = 0, High, Ret;
    USHORT Ordinal;
    PVOID Function;
    ULONG ExportSize;
    PULONG ExportTable;

    //
    // Get the export directory.
    //

    ExportDirectory = RtlImageDirectoryEntryToData(DllBase, TRUE, IMAGE_DIRECTORY_ENTRY_EXPORT, &ExportSize);

    if (!ExportDirectory)
    {
        return NULL;
    }

    //
    // Setup name tables.
    //

    NameTable = (PULONG)((ULONG_PTR)DllBase + ExportDirectory->AddressOfNames);
    OrdinalTable = (PUSHORT)((ULONG_PTR)DllBase + ExportDirectory->AddressOfNameOrdinals);

    //
    // Do a binary search.
    //

    High = ExportDirectory->NumberOfNames - 1;
    while (High >= Low)
    {
        //
        // Get new middle value.
        //

        Mid = (Low + High) >> 1;

        //
        // Compare name.
        //

        Ret = strcmp(ExportName->Buffer, (PCHAR)DllBase + NameTable[Mid]);
        if (Ret < 0)
        {
            //
            // Update high.
            //
            High = Mid - 1;
        }
        else if (Ret > 0)
        {
            //
            // Update low.
            //
            Low = Mid + 1;
        }
        else
        {
            //
            // We got it.
            //
            break;
        }
    }

    //
    // Check if we couldn't find it.
    //

    if (High < Low)
    {
        return NULL;
    }

    //
    // Otherwise, this is the ordinal.
    //

    Ordinal = OrdinalTable[Mid];

    //
    // Validate the ordinal.
    //

    if (Ordinal >= ExportDirectory->NumberOfFunctions)
    {
        return NULL;
    }

    //
    // Resolve the address and write it.
    //

    ExportTable = (PULONG)((ULONG_PTR)DllBase + ExportDirectory->AddressOfFunctions);
    Function = (PVOID)((ULONG_PTR)DllBase + ExportTable[Ordinal]);

    //
    // We found it!
    //

    NT_ASSERT((Function < (PVOID)ExportDirectory) || (Function > (PVOID)((ULONG_PTR)ExportDirectory + ExportSize)));

    return Function;
}

BOOLEAN
NTAPI
RtlxSuffixUnicodeString(_In_ PUNICODE_STRING String1, _In_ PUNICODE_STRING String2, _In_ BOOLEAN CaseInSensitive)
{
    //
    // RtlSuffixUnicodeString is not exported by ntoskrnl until Win10.
    //

    return String2->Length >= String1->Length &&
           RtlCompareUnicodeStrings(
               String2->Buffer + (String2->Length - String1->Length) / sizeof(WCHAR),
               String1->Length / sizeof(WCHAR),
               String1->Buffer,
               String1->Length / sizeof(WCHAR),
               CaseInSensitive) == 0;
}

PINJECTION_INFO
NTAPI
InjFindInjectionInfo(_In_ HANDLE ProcessId)
{
    PLIST_ENTRY NextEntry = InjInfoListHead.Flink;

    while (NextEntry != &InjInfoListHead)
    {
        PINJECTION_INFO InjectionInfo = CONTAINING_RECORD(NextEntry, INJECTION_INFO, ListEntry);
        if (InjectionInfo->ProcessId == ProcessId)
        {
            return InjectionInfo;
        }

        NextEntry = NextEntry->Flink;
    }

    return NULL;
}

VOID NTAPI
InjRemoveInjectionInfo(_In_ PINJECTION_INFO InjectionInfo, _In_ BOOLEAN FreeMemory)
{
    RemoveEntryList(&InjectionInfo->ListEntry);

    if (FreeMemory)
    {
        ExFreePoolWithTag(InjectionInfo, INJ_MEMORY_TAG);
    }
}

VOID NTAPI
InjRemoveInjectionInfoByProcessId(_In_ HANDLE ProcessId, _In_ BOOLEAN FreeMemory)
{
    PINJECTION_INFO InjectionInfo = InjFindInjectionInfo(ProcessId);

    if (InjectionInfo)
    {
        InjRemoveInjectionInfo(InjectionInfo, FreeMemory);
    }
}

BOOLEAN
NTAPI
InjCanInject(_In_ PINJECTION_INFO InjectionInfo)
{
    //
    // DLLs that need to be loaded in the native process
    // (i.e.: x64 process on x64 Windows, x86 process on
    // x86 Windows) before we can safely load our DLL.
    //

    ULONG RequiredDlls = INJ_SYSTEM32_NTDLL_LOADED;

#if defined(INJ_CONFIG_SUPPORTS_WOW64)

    if (PsGetProcessWow64Process(PsGetCurrentProcess()))
    {
        //
        // DLLs that need to be loaded in the Wow64 process
        // before we can safely load our DLL.
        //

        RequiredDlls |= INJ_SYSTEM32_NTDLL_LOADED;
        RequiredDlls |= INJ_SYSTEM32_WOW64_LOADED;
        RequiredDlls |= INJ_SYSTEM32_WOW64WIN_LOADED;

#    if defined(_M_AMD64)

        RequiredDlls |= INJ_SYSTEM32_WOW64CPU_LOADED;
        RequiredDlls |= INJ_SYSWOW64_NTDLL_LOADED;

#    elif defined(_M_ARM64)

        switch (PsWow64GetProcessMachine(PsGetCurrentProcess()))
        {
        case IMAGE_FILE_MACHINE_I386:

            //
            // Emulated x86 processes can load either SyCHPE32\ntdll.dll or
            // SysWOW64\ntdll.dll - depending on whether "hybrid execution
            // mode" is enabled or disabled.
            //
            // PsWow64GetProcessNtdllType(Process) can provide this information,
            // by returning EPROCESS->Wow64Process.NtdllType.  Unfortunatelly,
            // that function is not exported and EPROCESS is not documented.
            //
            // The solution here is to pick the Wow64 NTDLL which is already
            // loaded and set it as "required".
            //

            RequiredDlls |= InjectionInfo->LoadedDlls & (INJ_SYSWOW64_NTDLL_LOADED | INJ_SYCHPE32_NTDLL_LOADED);
            RequiredDlls |= INJ_SYSTEM32_XTAJIT_LOADED;
            break;

        case IMAGE_FILE_MACHINE_ARMNT:
            RequiredDlls |= INJ_SYSARM32_NTDLL_LOADED;
            RequiredDlls |= INJ_SYSTEM32_WOWARMHW_LOADED;
            break;

        case IMAGE_FILE_MACHINE_ARM64:
            break;
        }

#    endif
    }

#endif

    return (InjectionInfo->LoadedDlls & RequiredDlls) == RequiredDlls;
}

VOID NTAPI
InjIsWantedSystemDllBeingLoaded(
    _In_ PINJECTION_INFO InjectionInfo,
    _In_ PUNICODE_STRING FullImageName,
    _In_ PIMAGE_INFO ImageInfo)
{
    for (ULONG Index = 0; Index < RTL_NUMBER_OF(InjpSystemDlls); Index += 1)
    {
        PUNICODE_STRING SystemDllPath = &InjpSystemDlls[Index].DllPath;

        if (RtlxSuffixUnicodeString(SystemDllPath, FullImageName, TRUE))
        {
            PVOID LdrLoadDllRoutineAddress =
                RtlxFindExportedRoutineByName(ImageInfo->ImageBase, &LdrLoadDllRoutineName);

            ULONG DllFlag = InjpSystemDlls[Index].Flag;
            InjectionInfo->LoadedDlls |= DllFlag;

            switch (DllFlag)
            {
            case INJ_SYSARM32_NTDLL_LOADED:
            case INJ_SYCHPE32_NTDLL_LOADED:
            case INJ_SYSWOW64_NTDLL_LOADED:
            case INJ_SYSTEM32_NTDLL_LOADED:
                InjectionInfo->LdrLoadDllRoutineAddress = LdrLoadDllRoutineAddress;
                break;
            default:
                break;
            }

            //
            // Break the for-loop.
            //

            break;
        }
    }
}

NTSTATUS
NTAPI
InjpJoinPath(_In_ PUNICODE_STRING Directory, _In_ PUNICODE_STRING Filename, _Inout_ PUNICODE_STRING FullPath)
{
    UNICODE_STRING UnicodeBackslash = RTL_CONSTANT_STRING(L"\\");

    BOOLEAN DirectoryEndsWithBackslash = Directory->Length > 0 && Directory->Buffer[Directory->Length - 1] == L'\\';

    if (FullPath->MaximumLength < Directory->Length ||
        FullPath->MaximumLength - Directory->Length - (!DirectoryEndsWithBackslash ? 1 : 0) < Filename->Length)
    {
        return STATUS_DATA_ERROR;
    }

    RtlCopyUnicodeString(FullPath, Directory);

    if (!DirectoryEndsWithBackslash)
    {
        RtlAppendUnicodeStringToString(FullPath, &UnicodeBackslash);
    }

    RtlAppendUnicodeStringToString(FullPath, Filename);

    return STATUS_SUCCESS;
}

NTSTATUS
NTAPI
InjCreateSettings(_In_ PUNICODE_STRING RegistryPath, _Inout_ PINJ_SETTINGS Settings)
{
    //
    // In the "ImagePath" key of the RegistryPath, there
    // is a full path of this driver file.  Fetch it.
    //

    NTSTATUS Status;

    UNICODE_STRING ValueName = RTL_CONSTANT_STRING(L"ImagePath");

    OBJECT_ATTRIBUTES ObjectAttributes;
    InitializeObjectAttributes(&ObjectAttributes, RegistryPath, OBJ_CASE_INSENSITIVE | OBJ_KERNEL_HANDLE, NULL, NULL);

    HANDLE KeyHandle;
    Status = ZwOpenKey(&KeyHandle, KEY_READ, &ObjectAttributes);

    if (!NT_SUCCESS(Status))
    {
        return Status;
    }

    //
    // Save all information on stack - simply fail if path
    // is too long.
    //

    UCHAR KeyValueInformationBuffer[sizeof(KEY_VALUE_FULL_INFORMATION) + sizeof(WCHAR) * 128];
    PKEY_VALUE_FULL_INFORMATION KeyValueInformation = (PKEY_VALUE_FULL_INFORMATION)KeyValueInformationBuffer;

    ULONG ResultLength;
    Status = ZwQueryValueKey(
        KeyHandle,
        &ValueName,
        KeyValueFullInformation,
        KeyValueInformation,
        sizeof(KeyValueInformationBuffer),
        &ResultLength);

    ZwClose(KeyHandle);

    //
    // Check for succes.  Also check if the value is of expected
    // type and whether the path has a meaninful length.
    //

    if (!NT_SUCCESS(Status) || KeyValueInformation->Type != REG_EXPAND_SZ ||
        KeyValueInformation->DataLength < sizeof(ObpDosDevicesShortNamePrefix))
    {
        return Status;
    }

    //
    // Save pointer to the fetched ImagePath value and test if
    // the path starts with "\??\" prefix - if so, skip it.
    //

    PWCHAR ImagePathValue = (PWCHAR)((PUCHAR)KeyValueInformation + KeyValueInformation->DataOffset);
    ULONG ImagePathValueLength = KeyValueInformation->DataLength;

    if (*(PULONGLONG)(ImagePathValue) == ObpDosDevicesShortNamePrefix.Alignment.QuadPart)
    {
        ImagePathValue += ObpDosDevicesShortName.Length / sizeof(WCHAR);
        ImagePathValueLength -= ObpDosDevicesShortName.Length;
    }

    //
    // Cut the string by the last '\' character, leaving there
    // only the directory path.
    //

    PWCHAR LastBackslash = wcsrchr(ImagePathValue, L'\\');

    if (!LastBackslash)
    {
        return STATUS_DATA_ERROR;
    }

    *LastBackslash = UNICODE_NULL;

    UNICODE_STRING Directory;
    RtlInitUnicodeString(&Directory, ImagePathValue);

    //
    // Finally, fill all the buffers...
    //

#define INJ_DLL_X86_NAME L"HookDLL-x86.dll"
    UNICODE_STRING InjDllNameX86 = RTL_CONSTANT_STRING(INJ_DLL_X86_NAME);
    InjpJoinPath(&Directory, &InjDllNameX86, &Settings->DllPath[InjArchitectureX86]);
    LOG_INFO("DLL path (x86):   '%wZ'", &Settings->DllPath[InjArchitectureX86]);

#define INJ_DLL_X64_NAME L"HookDLL-x64.dll"
    UNICODE_STRING InjDllNameX64 = RTL_CONSTANT_STRING(INJ_DLL_X64_NAME);
    InjpJoinPath(&Directory, &InjDllNameX64, &Settings->DllPath[InjArchitectureX64]);
    LOG_INFO("DLL path (x64):   '%wZ'", &Settings->DllPath[InjArchitectureX64]);

#define INJ_DLL_ARM32_NAME L"HookDLL-ARM.dll"
    UNICODE_STRING InjDllNameARM32 = RTL_CONSTANT_STRING(INJ_DLL_ARM32_NAME);
    InjpJoinPath(&Directory, &InjDllNameARM32, &Settings->DllPath[InjArchitectureARM32]);
    LOG_INFO("DLL path (ARM32): '%wZ'", &Settings->DllPath[InjArchitectureARM32]);

#define INJ_DLL_ARM64_NAME L"HookDLL-ARM64.dll"
    UNICODE_STRING InjDllNameARM64 = RTL_CONSTANT_STRING(INJ_DLL_ARM64_NAME);
    InjpJoinPath(&Directory, &InjDllNameARM64, &Settings->DllPath[InjArchitectureARM64]);
    LOG_INFO("DLL path (ARM64): '%wZ'", &Settings->DllPath[InjArchitectureARM64]);

    return STATUS_SUCCESS;
}

NTSTATUS
NTAPI
InjInitialize(_In_ PUNICODE_STRING RegistryPath)
{
    //
    // Create injection settings.
    //

    INJ_SETTINGS Settings;

    WCHAR BufferDllPathX86[128];
    Settings.DllPath[InjArchitectureX86].Length = 0;
    Settings.DllPath[InjArchitectureX86].MaximumLength = sizeof(BufferDllPathX86);
    Settings.DllPath[InjArchitectureX86].Buffer = BufferDllPathX86;

    WCHAR BufferDllPathX64[128];
    Settings.DllPath[InjArchitectureX64].Length = 0;
    Settings.DllPath[InjArchitectureX64].MaximumLength = sizeof(BufferDllPathX64);
    Settings.DllPath[InjArchitectureX64].Buffer = BufferDllPathX64;

    WCHAR BufferDllPathARM32[128];
    Settings.DllPath[InjArchitectureARM32].Length = 0;
    Settings.DllPath[InjArchitectureARM32].MaximumLength = sizeof(BufferDllPathARM32);
    Settings.DllPath[InjArchitectureARM32].Buffer = BufferDllPathARM32;

    WCHAR BufferDllPathARM64[128];
    Settings.DllPath[InjArchitectureARM64].Length = 0;
    Settings.DllPath[InjArchitectureARM64].MaximumLength = sizeof(BufferDllPathARM64);
    Settings.DllPath[InjArchitectureARM64].Buffer = BufferDllPathARM64;

    NTSTATUS ntStatus;
    ntStatus = InjCreateSettings(RegistryPath, &Settings);
    if (!NT_SUCCESS(ntStatus))
    {
        LOG_ERROR("Failed to create injection settings");
        return ntStatus;
    }

    //
    // Initialize injection info linked list.
    //

    InitializeListHead(&InjInfoListHead);

    ULONG Flags = RTL_DUPLICATE_UNICODE_STRING_NULL_TERMINATE | RTL_DUPLICATE_UNICODE_STRING_ALLOCATE_NULL_STRING;

    PINJ_SETTINGS pSettings = &Settings;
    for (ULONG Architecture = 0; Architecture < InjArchitectureMax; Architecture += 1)
    {
        ntStatus = RtlDuplicateUnicodeString(Flags, &pSettings->DllPath[Architecture], &InjDllPath[Architecture]);
        if (!NT_SUCCESS(ntStatus))
        {
            return ntStatus;
        }
    }

    //
    // Check if we're running on Windows 7.
    //

    RTL_OSVERSIONINFOW VersionInformation = {0};
    VersionInformation.dwOSVersionInfoSize = sizeof(VersionInformation);
    RtlGetVersion(&VersionInformation);

    if (VersionInformation.dwMajorVersion == 6 && VersionInformation.dwMinorVersion == 1)
    {
        LOG_INFO("Current system is Windows 7");
        g_InjIsWindows7 = TRUE;
    }

    return STATUS_SUCCESS;
}

VOID NTAPI InjDestroy(VOID)
{
    //
    // Release memory of all injection-info entries.
    //

    PLIST_ENTRY NextEntry = InjInfoListHead.Flink;

    while (NextEntry != &InjInfoListHead)
    {
        PINJECTION_INFO InjectionInfo = CONTAINING_RECORD(NextEntry, INJECTION_INFO, ListEntry);
        NextEntry = NextEntry->Flink;

        ExFreePoolWithTag(InjectionInfo, INJ_MEMORY_TAG);
    }

    //
    // Release memory of all buffers.
    //

    for (ULONG Architecture = 0; Architecture < InjArchitectureMax; Architecture += 1)
    {
        RtlFreeUnicodeString(&InjDllPath[Architecture]);
    }
}

VOID NTAPI
InjpInjectApcNormalRoutine(_In_ PVOID NormalContext, _In_ PVOID SystemArgument1, _In_ PVOID SystemArgument2)
{
    UNREFERENCED_PARAMETER(SystemArgument1);
    UNREFERENCED_PARAMETER(SystemArgument2);

    PINJECTION_INFO InjectionInfo = NormalContext;
    InjInject(InjectionInfo);
}

VOID NTAPI
InjInjectApcKernelRoutine(
    _In_ PKAPC Apc,
    _Inout_ PKNORMAL_ROUTINE *NormalRoutine,
    _Inout_ PVOID *NormalContext,
    _Inout_ PVOID *SystemArgument1,
    _Inout_ PVOID *SystemArgument2)
{
    UNREFERENCED_PARAMETER(NormalRoutine);
    UNREFERENCED_PARAMETER(NormalContext);
    UNREFERENCED_PARAMETER(SystemArgument1);
    UNREFERENCED_PARAMETER(SystemArgument2);

    //
    // Common kernel routine for both user-mode and
    // kernel-mode APCs queued by the InjpQueueApc
    // function.  Just release the memory of the APC
    // structure and return back.
    //

    ExFreePoolWithTag(Apc, INJ_MEMORY_TAG);
}

NTSTATUS
NTAPI
InjQueueApc(
    _In_ KPROCESSOR_MODE ApcMode,
    _In_ PKNORMAL_ROUTINE NormalRoutine,
    _In_ PVOID NormalContext,
    _In_ PVOID SystemArgument1,
    _In_ PVOID SystemArgument2)
{
    //
    // Allocate memory for the KAPC structure.
    //

    PKAPC Apc = ExAllocatePoolWithTag(NonPagedPoolNx, sizeof(KAPC), INJ_MEMORY_TAG);

    if (!Apc)
    {
        return STATUS_INSUFFICIENT_RESOURCES;
    }

    //
    // Initialize and queue the APC.
    //

    KeInitializeApc(
        Apc,                        // Apc
        PsGetCurrentThread(),       // Thread
        OriginalApcEnvironment,     // Environment
        &InjInjectApcKernelRoutine, // KernelRoutine
        NULL,                       // RundownRoutine
        NormalRoutine,              // NormalRoutine
        ApcMode,                    // ApcMode
        NormalContext);             // NormalContext

    BOOLEAN Inserted = KeInsertQueueApc(
        Apc,             // Apc
        SystemArgument1, // SystemArgument1
        SystemArgument2, // SystemArgument2
        0);              // Increment

    if (!Inserted)
    {
        ExFreePoolWithTag(Apc, INJ_MEMORY_TAG);
        return STATUS_UNSUCCESSFUL;
    }

    return STATUS_SUCCESS;
}

NTSTATUS
NTAPI
InjInject(_In_ PINJECTION_INFO InjectionInfo)
{
    NTSTATUS Status;

    //
    // Create memory space for injection-specific data,
    // such as path to the to-be-injected DLL.  Memory
    // of this section will be eventually mapped to the
    // injected process.
    //
    // Note that this memory is created using sections
    // instead of ZwAllocateVirtualMemory, mainly because
    // function ZwProtectVirtualMemory is not exported
    // by ntoskrnl.exe until Windows 8.1.  In case of
    // sections, the effect of memory protection change
    // is achieved by remaping the section with different
    // protection type.
    //

    OBJECT_ATTRIBUTES ObjectAttributes;
    InitializeObjectAttributes(&ObjectAttributes, NULL, OBJ_KERNEL_HANDLE, NULL, NULL);

    HANDLE SectionHandle;
    SIZE_T SectionSize = PAGE_SIZE;
    LARGE_INTEGER MaximumSize;
    MaximumSize.QuadPart = SectionSize;
    Status = ZwCreateSection(
        &SectionHandle,
        GENERIC_READ | GENERIC_WRITE,
        &ObjectAttributes,
        &MaximumSize,
        PAGE_EXECUTE_READWRITE,
        SEC_COMMIT,
        NULL);

    if (!NT_SUCCESS(Status))
    {
        return Status;
    }

    INJ_ARCHITECTURE Architecture = InjArchitectureMax;

#if defined(_M_IX86)
    Architecture = InjArchitectureX86;

#elif defined(_M_AMD64)
    Architecture = PsGetProcessWow64Process(PsGetCurrentProcess()) ? InjArchitectureX86 : InjArchitectureX64;

#elif defined(_M_ARM64)
    switch (PsWow64GetProcessMachine(PsGetCurrentProcess()))
    {
    case IMAGE_FILE_MACHINE_I386:
        Architecture = InjArchitectureX86;
        break;

    case IMAGE_FILE_MACHINE_ARMNT:
        Architecture = InjArchitectureARM32;
        break;

    case IMAGE_FILE_MACHINE_ARM64:
        Architecture = InjArchitectureARM64;
        break;
    }
#endif

    NT_ASSERT(Architecture != InjArchitectureMax);

    InjpInject(InjectionInfo, Architecture, SectionHandle, SectionSize);

    ZwClose(SectionHandle);

    if (NT_SUCCESS(Status) && InjectionInfo->ForceUserApc)
    {
        //
        // Sets CurrentThread->ApcState.UserApcPending to TRUE.
        // This causes the queued user APC to be triggered immediately
        // on next transition of this thread to the user-mode.
        //

        KeTestAlertThread(UserMode);
    }

    return Status;
}

NTSTATUS
NTAPI
InjpInject(
    _In_ PINJECTION_INFO InjectionInfo,
    _In_ INJ_ARCHITECTURE Architecture,
    _In_ HANDLE SectionHandle,
    _In_ SIZE_T SectionSize)
{
    NTSTATUS Status;

    //
    // First, map this section with read-write access.
    //

    PVOID SectionMemoryAddress = NULL;
    Status = ZwMapViewOfSection(
        SectionHandle,
        ZwCurrentProcess(),
        &SectionMemoryAddress,
        0,
        SectionSize,
        NULL,
        &SectionSize,
        ViewUnmap,
        0,
        PAGE_READWRITE);

    if (!NT_SUCCESS(Status))
    {
        goto Exit;
    }

    //
    // Code of the APC routine (ApcNormalRoutine defined in the
    // "shellcode" above) starts at the SectionMemoryAddress.
    // Copy the shellcode to the allocated memory.
    //

    PVOID ApcRoutineAddress = SectionMemoryAddress;
    RtlCopyMemory(ApcRoutineAddress, InjThunk[Architecture].Buffer, InjThunk[Architecture].Length);

    //
    // Fill the data of the ApcContext.
    //

    PWCHAR DllPath = (PWCHAR)((PUCHAR)SectionMemoryAddress + InjThunk[Architecture].Length);
    RtlCopyMemory(DllPath, InjDllPath[Architecture].Buffer, InjDllPath[Architecture].Length);

    //
    // Unmap the section and map it again, but now
    // with read-execute (no write) access.
    //

    ZwUnmapViewOfSection(ZwCurrentProcess(), SectionMemoryAddress);

    SectionMemoryAddress = NULL;
    Status = ZwMapViewOfSection(
        SectionHandle,
        ZwCurrentProcess(),
        &SectionMemoryAddress,
        0,
        PAGE_SIZE,
        NULL,
        &SectionSize,
        ViewUnmap,
        0,
        PAGE_EXECUTE_READ);

    if (!NT_SUCCESS(Status))
    {
        goto Exit;
    }

    //
    // Reassign remapped address.
    //

    ApcRoutineAddress = SectionMemoryAddress;
    DllPath = (PWCHAR)((PUCHAR)SectionMemoryAddress + InjThunk[Architecture].Length);

    PVOID ApcContext = (PVOID)InjectionInfo->LdrLoadDllRoutineAddress;
    PVOID ApcArgument1 = (PVOID)DllPath;
    PVOID ApcArgument2 = (PVOID)InjDllPath[Architecture].Length;

#if defined(INJ_CONFIG_SUPPORTS_WOW64)

    if (PsGetProcessWow64Process(PsGetCurrentProcess()))
    {
        //
        // The ARM32 ntdll.dll uses "BLX" instruction for calling the ApcRoutine.
        // This instruction can change the processor state (between Thumb & ARM),
        // based on the LSB (least significant bit).  If this bit is 0, the code
        // will run in the ARM instruction set.  If this bit is 1, the code will
        // run in Thumb instruction set.  Because Windows can run only in the Thumb
        // instruction set, we have to ensure this bit is set.  Otherwise, Windows
        // would raise STATUS_ILLEGAL_INSTRUCTION upon execution of the ApcRoutine.
        //

        if (Architecture == InjArchitectureARM32)
        {
            ApcRoutineAddress = (PVOID)((ULONG_PTR)ApcRoutineAddress | 1);
        }

        //
        // PsWrapApcWow64Thread essentially assigns wow64.dll!Wow64ApcRoutine
        // to the NormalRoutine.  This Wow64ApcRoutine (which is 64-bit code)
        // in turn calls KiUserApcDispatcher (in 32-bit ntdll.dll) which finally
        // calls our provided ApcRoutine.
        //

        PsWrapApcWow64Thread(&ApcContext, &ApcRoutineAddress);
    }

#endif

    PKNORMAL_ROUTINE ApcRoutine = (PKNORMAL_ROUTINE)(ULONG_PTR)ApcRoutineAddress;

    Status = InjQueueApc(UserMode, ApcRoutine, ApcContext, ApcArgument1, ApcArgument2);

    if (!NT_SUCCESS(Status))
    {
        //
        // If injection failed for some reason, unmap the section.
        //

        ZwUnmapViewOfSection(ZwCurrentProcess(), SectionMemoryAddress);
    }

Exit:
    return Status;
}
