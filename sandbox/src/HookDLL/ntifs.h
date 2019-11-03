#pragma once

#include "stdafx.h"


NTSTATUS WINAPI HookRtlDecompressBuffer(
	_In_ USHORT CompressionFormat,
	_Out_writes_bytes_to_(UncompressedBufferSize, *FinalUncompressedSize) PUCHAR UncompressedBuffer,
	_In_ ULONG UncompressedBufferSize,
	_In_reads_bytes_(CompressedBufferSize) PUCHAR CompressedBuffer,
	_In_ ULONG CompressedBufferSize,
	_Out_ PULONG FinalUncompressedSize
);
