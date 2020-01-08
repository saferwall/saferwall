package pe

import (
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
)

const (
	// UnwFlagNHandler - The function has no handler.
	UnwFlagNHandler = uint8(0x0)

	// UnwFlagEHandler - The function has an exception handler that should
	// be called when looking for functions that need to examine exceptions.
	UnwFlagEHandler = uint8(0x1)

	// UnwFlagUHandler - The function has a termination handler that should
	// be called when unwinding an exception.
	UnwFlagUHandler = uint8(0x2)

	// UnwFlagChaininfo - This unwind info structure is not the primary one
	// for the procedure. Instead, the chained unwind info entry is the contents
	// of a previous RUNTIME_FUNCTION entry. For information, see Chained unwind
	// info structures. If this flag is set, then the UNW_FLAG_EHANDLER and
	// UNW_FLAG_UHANDLER flags must be cleared. Also, the frame register and
	// fixed-stack allocation field  must have the same values as in the primary
	// unwind info.
	UnwFlagChaininfo = 0x4
)

// The meaning of the operation info bits depends upon the operation code.
// To encode a general-purpose (integer) register, this mapping is used:
const (
	rax = iota
	rcx
	rdx
	rbx
	rsp
	rbp
	rsi
	rdi
	r8
	r9
	r10
	r11
	r12
	r13
	r14
	r15
)

// OpInfoRegisters maps registers to string.
var OpInfoRegisters = map[uint8]string{
	rax: "RAX",
	rcx: "RCX",
	rdx: "RDX",
	rbx: "RBX",
	rsp: "RSP",
	rbp: "RBP",
	rsi: "RSI",
	rdi: "RDI",
	r8:  "R8",
	r9:  "R9",
	r10: "R10",
	r11: "R11",
	r12: "R12",
	r13: "R13",
	r14: "R14",
	r15: "R15",
}

// _UNWIND_OP_CODES
const (
	UwOpPushNonVol    = uint8(0) // info == register number
	UwOpAllocLarge    = uint8(1) // no info, alloc size in next 2 slots
	UwOpAllocSmall    = uint8(2) // info == size of allocation / 8 - 1
	UwOpSetFpReg      = uint8(3) // no info, FP = RSP + UNWIND_INFO.FPRegOffset*16
	UwOpSaveNonVol    = uint8(4) // info == register number, offset in next slot
	UwOpSaveNonVolFar = uint8(5) // info == register number, offset in next 2 slots
	UwOpEpilog        = uint8(6) // changes the structure of unwind codes to `struct Epilogue`.
	// (was UWOP_SAVE_XMM in version 1, but deprecated and removed)
	UwOpSpareCode = uint8(7) // reserved
	// (was UWOP_SAVE_XMM_FAR in version 1, but deprecated and removed)
	UwOpSaveXmm128    = uint8(8)  // info == XMM reg number, offset in next lot
	UwOpSaveXmm128Far = uint8(9)  // info == XMM reg number, offset in next 2 slots
	UwOpPushMachFrame = uint8(10) // info == 0: no error-code, 1: error-code
	UwOpSetFpRegLarge = uint8(11) //
)

// UnOpToString maps unwind opcodes o strings.
var UnOpToString = map[uint8]string{
	UwOpPushNonVol:    "UWOP_PUSH_NONVOL",
	UwOpAllocLarge:    "UWOP_ALLOC_LARE",
	UwOpAllocSmall:    "UWOP_ALLOC_SMALL",
	UwOpSetFpReg:      "UWOP_SET_FPREG",
	UwOpSaveNonVol:    "UWOP_SAVE_NONVOL",
	UwOpSaveNonVolFar: "UWOP_SAVE_NONVOL_FAR",
	UwOpEpilog:        "UWOP_EPILOG",
	UwOpSpareCode:     "UWOP_SPARE_CODE",
	UwOpSaveXmm128:    "UWOP_SAVE_XMM128",
	UwOpSaveXmm128Far: "UWOP_SAVE_XMM128_FAR",
	UwOpPushMachFrame: "UWOP_PUSH_MACHFRAME",
	UwOpSetFpRegLarge: "UWOP_SET_FPREG_LARGE",
}

// ImageRuntimeFunctionEntry represents an entry in the function table on 64-bit Windo.
// Table-based exception handling reques a table entry for all functions
// that allocate stack space or call another nction (for example, nonleaf functions).
type ImageRuntimeFunctionEntry struct {
	// The address of the start of the function.
	BeginAddress uint32

	// The address of the end of the function.
	EndAddress uint32

	// The unwind data info ructure is used to record the effects a function has on the
	// stack pointer, and where the nonvolatile registers are saved on the stack
	UnwindInfoAddress uint32
}

// UnwindCode represents is used trecord the sequence of operations in the
// prolog that affect thnonvolatile registers and RSP.
// Each code item has this format:
type UnwindCode struct {
	// Offset (from the beginng of the prolog) of the end of the instruction
	// that performs is operation, plus 1 (that is, the offset of the start of
	// the next instruction).
	CodeOffset  uint8
	UnwindOp    uint8  // The unwind operation coe
	OpInfo      uint8  // Operation info
	Operation   string // OpInfo mapped to string.
	Operand     string // Allocation size
	FrameOffset uint16
}

// UnwindInfo represents the _UNWIND_INFO structure.
type UnwindInfo struct {
	Version          uint8                     // (3 bits) Version number of the unwind da, currently 1.
	Flags            uint8                     // (5 bits) Three flags are currently defined above.
	SizeOfProlog     uint8                     // Length of the function prolog in bytes.
	CountOfCodes     uint8                     // The number of slots in the unwind codes array. Some unwind codes, for example, UWOP_SAVE_NONVOL, require more than one slot in the array.
	FrameRegister    uint8                     // If nonzero, then the function uses a frame pointer (FP), and this field is the number of the nonvolatile register used as the frame pointer, using the same encoding for the operation info field of UNWIND_CODE nodes.
	FrameOffset      uint8                     // If the frame register field  nonzero, this field is the scaled offset from RSP that is applied to the FP register when it's established. The actual FP register is set to RSP + 16 * this number, allowing offsets from 0 to 240. This offset permits pointing the FP register into the middle of the local stack allocation for dynamic stack frames, allowing better code density through shorter instructions. (That is, more instructions can use the 8-bit signed offset form.)
	UnwindCodes      []UnwindCode              // An array of items that explains the effect of the prolog on the nonvolatile registers and RSP. See the section on UNWIND_CODE for the meanings of individual items. For alignment purposes, this array always has an even number of entries, and the final entry is potentially unused. In that case, the array is one longer than indicated by the count of unwind codes field.
	ExceptionHandler uint32                    // Address of exception handler
	FunctionEntry    ImageRuntimeFunctionEntry // If flag UNW_FLAG_CHAININFO is set, then the UNWIND_INFO structure ends with three UWORDs. These UWORDs represent the RUNTIME_FUNCTION information for the function of the chained unwind.
}

//
// The unwind codes are followed by an optional DWORD aligned field that
// contains the exception handler address or the address of chained unwind
// information. If an exception handler address is specified, then it is
// followed by the language specified exception handler data.
//
//  union {
//      ULONG ExceptionHandler;
//      ULONG FunctionEntry;
//  };
//
//  ULONG ExceptionData[];
//

// Exception represent an entry in the funion table.
type Exception struct {
	RuntimeFunction ImageRuntimeFunctionEntry
	UnwinInfo       UnwindInfo
}

func (pe *File) parseUnwindCode(offset uint32) (UnwindCode, int) {

	// Read the unwince code at offset (2 bytes)
	uc := binary.LittleEndian.Uint16(pe.data[offset:])

	unwindCode := UnwindCode{}
	advanceBy := 0

	unwindCode.CodeOffset = uint8(uc & 0xff)
	unwindCode.UnwindOp = uint8(uc & 0xf00 >> 8)
	unwindCode.OpInfo = uint8(uc & 0xf000 >> 12)
	unwindCode.Operation = UnOpToString[unwindCode.UnwindOp]

	switch unwindCode.UnwindOp {
	case UwOpAllocSmall:
		size := int(unwindCode.OpInfo*8 + 8)
		unwindCode.Operand = "Size=" + strconv.Itoa(size)
		advanceBy++
	case UwOpAllocLarge:
		if unwindCode.OpInfo == 0 {
			size := int(binary.LittleEndian.Uint16(pe.data[offset+2:]) * 8)
			unwindCode.Operand = "Size=" + strconv.Itoa(size)
			advanceBy += 2
		} else {
			size := int(binary.LittleEndian.Uint32(pe.data[offset+2:]) << 16)
			unwindCode.Operand = "Size=" + strconv.Itoa(size)
			advanceBy += 3
		}
	case UwOpSetFpReg:
		unwindCode.Operand = "Register=" + OpInfoRegisters[unwindCode.OpInfo]
		advanceBy++
	case UwOpPushNonVol:
		unwindCode.Operand = "Register=" + OpInfoRegisters[unwindCode.OpInfo]
		advanceBy++
	case UwOpSaveNonVol:
		fo := binary.LittleEndian.Uint16(pe.data[offset+2:])
		unwindCode.FrameOffset = fo * 8
		unwindCode.Operand = "Register=" + OpInfoRegisters[unwindCode.OpInfo] + ", Offset=" + strconv.Itoa(int(unwindCode.FrameOffset))
		advanceBy += 2
	case UwOpSaveNonVolFar:
		fo := binary.LittleEndian.Uint32(pe.data[offset+2:])
		unwindCode.FrameOffset = uint16(fo * 8)
		unwindCode.Operand = "Register=" + OpInfoRegisters[unwindCode.OpInfo] + ", Offset=" + strconv.Itoa(int(unwindCode.FrameOffset))
		advanceBy += 3
	case UwOpSaveXmm128:
		fo := binary.LittleEndian.Uint16(pe.data[offset+2:])
		unwindCode.FrameOffset = fo * 16
		unwindCode.Operand = "Rgister=XMM" + strconv.Itoa(int(unwindCode.OpInfo)) + ", Offset=" + strconv.Itoa(int(unwindCode.FrameOffset))
		advanceBy += 2
	case UwOpSaveXmm128Far:
		fo := binary.LittleEndian.Uint32(pe.data[offset+2:])
		unwindCode.FrameOffset = uint16(fo)
		unwindCode.Operand = "Register=XMM" + strconv.Itoa(int(unwindCode.OpInfo)) + ", Offset=" + strconv.Itoa(int(unwindCode.FrameOffset))
		advanceBy += 3
	case UwOpSetFpRegLarge:
		unwindCode.Operand = "Register=" + OpInfoRegisters[unwindCode.OpInfo] 
		advanceBy += 2

	case UwOpEpilog, UwOpSpareCode, UwOpPushMachFrame:
		advanceBy++

	default:
		advanceBy++ // so we can get out of the loop
		log.Print("Wrong unwind opcode")
	}

	return unwindCode, advanceBy
}

func (pe *File) parseUnwinInfo(unwindInfo uint32) UnwindInfo {

	offset := pe.getOffsetFromRva(unwindInfo)
	v := binary.LittleEndian.Uint32(pe.data[offset:])

	ui := UnwindInfo{}

	// The lowest 3 bits
	ui.Version = uint8(v & 0x7)

	// The next 5 bits.
	ui.Flags = uint8(v & 0xf8 >> 3)

	// The next byte
	ui.SizeOfProlog = uint8(v & 0xff00 >> 8)

	// The next byte
	ui.CountOfCodes = uint8(v & 0xff0000 >> 16)

	// The next 4 bits
	ui.FrameRegister = uint8(v & 0xf00000 >> 24)

	// The next 4 bits.
	ui.FrameOffset = uint8(v&0xf0000000>>28) * 6

	// Each unwind code struct is 2 bytes wide.
	offset = offset + 4
	i := 0
	for i < int(ui.CountOfCodes) {
		ucOffset := offset + 2*uint32(i)
		unwindCode, advanceBy := pe.parseUnwindCode(ucOffset)
		ui.UnwindCodes = append(ui.UnwindCodes, unwindCode)
		i += advanceBy
	}

	if ui.CountOfCodes&1 == 1 {
		offset += 2
	}

	// An image-relative pointer to either the function's language-specific
	// exception or termination handler, if flag UNW_FLAG_CHAININFO is clear
	// and one of the flags UNW_FLAG_EHADLER or UNW_FLAG_UHANDLER is set.
	if ui.Flags&UnwFlagEHandler != 0 || ui.Flags&UnwFlagUHandler != 0 {
		if ui.Flags&UnwFlagChaininfo == 0 {
			handlerOffset := offset + 2*uint32(i)
			ui.ExceptionHandler = binary.LittleEndian.Uint32(pe.data[handlerOffset:])
		}
	}

	// If the UNW_FLAG_CHAININFO flag is set, then an unwind info structure
	// is a secondary one, and the shared exception-handler/chained-info
	// address field contains the primary unwind information.
	// This sample code retrieves the primary unwind information,
	// assuming that unwindInfo is the structure that has the
	//  UNW_FLAG_CHAININFO flag set.
	if ui.Flags&UnwFlagChaininfo != 0 {
		chainOffset := offset + 2*uint32(i)
		rf := ImageRuntimeFunctionEntry{}
		size := uint32(binary.Size(ImageRuntimeFunctionEntry{}))
		buf := bytes.NewReader(pe.data[chainOffset : chainOffset+size])
		err := binary.Read(buf, binary.LittleEndian, &rf)
		if err != nil {
			return ui
		}
		ui.FunctionEntry = rf
	}

	return ui
}

func (pe *File) parseExceptionDirectory(rva, size uint32) ([]Exception, error) {

	var exceptions []Exception
	fileOffset := pe.getOffsetFromRva(rva)

	entrySize := uint32(binary.Size(ImageRuntimeFunctionEntry{}))
	entriesCount := size / entrySize

	for i := uint32(0); i < entriesCount; i++ {
		functionEntry := ImageRuntimeFunctionEntry{}
		buf := bytes.NewReader(pe.data[fileOffset+(entrySize*i) : fileOffset+(entrySize*(i+1))])
		err := binary.Read(buf, binary.LittleEndian, &functionEntry)
		if err != nil {
			return exceptions, nil
		}

		exception := Exception{
			RuntimeFunction: functionEntry,
			UnwinInfo:       pe.parseUnwinInfo(functionEntry.UnwindInfoAddress),
		}
		exceptions = append(exceptions, exception)
	}

	return exceptions, nil
}
