#include "monitor.h"

VOID
PushMonitoredProcessEntry(PSINGLE_LIST_ENTRY ListHead, PMONITORED_PROCESS_ENTRY Entry)
{
    PushEntryList(ListHead, &(Entry->SingleListEntry));
}

PMONITORED_PROCESS_ENTRY
PopMonitoredProcessEntry(PSINGLE_LIST_ENTRY ListHead)
{
    PSINGLE_LIST_ENTRY SingleListEntry;
    SingleListEntry = PopEntryList(ListHead);
    return CONTAINING_RECORD(SingleListEntry, MONITORED_PROCESS_ENTRY, SingleListEntry);
}

NTSTATUS
MonAddProcessToMonitoredList(ULONG TargetPid)
/*++

Routine Description:
    Add the target pid to the process monitoring list.

Arguments:
    TargetPid - Process Identifier of the app to monitor.

Return Value:
    STATUS_SUCCESS if no error was encountered; otherwise, relevant NTSTATUS code.

--*/
{
    if (TargetPid == 0)
    {
        return STATUS_INVALID_PARAMETER;
    }

    if (MonIsProcessMonitored(TargetPid))
    {
        return STATUS_SUCCESS;
    }

	PMONITORED_PROCESS_ENTRY NewEntry;
    NewEntry =
        (PMONITORED_PROCESS_ENTRY)ExAllocatePoolWithTag(NonPagedPool, sizeof(MONITORED_PROCESS_ENTRY), MONIT_POOL_TAG);
    if (NewEntry == NULL)
        return STATUS_NO_MEMORY;

	NewEntry->ProcessId = TargetPid;
    PushMonitoredProcessEntry(&gMonitoredProcessList->SingleListEntry, NewEntry);
    return STATUS_SUCCESS;
}

BOOLEAN
MonIsProcessMonitored(ULONG ProcessId)
/*++

Routine Description:
    Returns TRUE if the process is being monitored.

Arguments:
    ProcessId - Process Identifier.

Return Value:
    TRUE if process is found, FALSE otherwise.

--*/
{
    PMONITORED_PROCESS_ENTRY pTmpSListEntry;

	if (ProcessId == 0)
	{
		return FALSE;
	}

    pTmpSListEntry = gMonitoredProcessList;
    while (pTmpSListEntry != NULL)
    {
        if (pTmpSListEntry->ProcessId == ProcessId)
        {
            return TRUE;
		}

        pTmpSListEntry = (PMONITORED_PROCESS_ENTRY)(pTmpSListEntry->SingleListEntry.Next);
    }

    return FALSE;
}
