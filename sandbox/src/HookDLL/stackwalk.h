#pragma once

#define MAX_FRAME 10


typedef struct _STACKTRACE
{
    //
    // Number of frames in Frames array.
    //
    UINT FrameCount;

    //
    // PC-Addresses of frames. Index 0 contains the topmost frame.
    //
    ULONGLONG Frames[MAX_FRAME];
} STACKTRACE, *PSTACKTRACE;



VOID
CaptureStackTrace();
VOID
AllocateSpaceSymbol();
VOID
SfwCaptureStackFrames(PSTACKTRACE StackTrace, UINT MaxFrames);
BOOL
SfwIsCalledFromSystemMemory(UINT FramesToCapture);
NTSTATUS
SfwSymInit();