#include "stdafx.h"
#include "ntifs.h"

decltype(RtlDecompressBuffer) *TrueRtlDecompressBuffer = nullptr;

NTSTATUS WINAPI
HookRtlDecompressBuffer(
    _In_ USHORT CompressionFormat,
    _Out_writes_bytes_to_(UncompressedBufferSize, *FinalUncompressedSize) PUCHAR UncompressedBuffer,
    _In_ ULONG UncompressedBufferSize,
    _In_reads_bytes_(CompressedBufferSize) PUCHAR CompressedBuffer,
    _In_ ULONG CompressedBufferSize,
    _Out_ PULONG FinalUncompressedSize)
{
    if (IsInsideHook())
    {
        goto end;
    }

    CaptureStackTrace();

    TraceAPI(L"RtlDecompressBuffer(CompressionFormat: %hu), RETN: 0x%p", _ReturnAddress());

    ReleaseHookGuard();
end:
    return TrueRtlDecompressBuffer(
        CompressionFormat,
        UncompressedBuffer,
        UncompressedBufferSize,
        CompressedBuffer,
        CompressedBufferSize,
        FinalUncompressedSize);
}
