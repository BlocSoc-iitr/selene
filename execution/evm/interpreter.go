package evm
import (
	"encoding/json"
	"math"
)
type Interpreter struct {
	InstructionPointer *byte // Using *byte as equivalent to *const u8 in Rust
	Gas                Gas
	Contract           Contract
	InstructionResult  InstructionResult
	Bytecode           Bytes
	IsEof              bool
	// Is init flag for eof create
	IsEofInit bool
	// Shared memory.
	// Note: This field is only set while running the interpreter loop.
	// Otherwise, it is taken and replaced with empty shared memory.
	SharedMemory SharedMemory
	// Stack.
	Stack Stack
	// EOF function stack.
	FunctionStack FunctionStack
	// The return data buffer for internal calls.
	// It has multi usage:
	// - It contains the output bytes of call sub call.
	// - When this interpreter finishes execution it contains the output bytes of this contract.
	ReturnDataBuffer Bytes
	// Whether the interpreter is in "staticcall" mode, meaning no state changes can happen.
	IsStatic bool
	// Actions that the EVM should do.
	// Set inside CALL or CREATE instructions and RETURN or REVERT instructions.
	// Additionally, those instructions will set InstructionResult to CallOrCreate/Return/Revert so we know the reason.
	NextAction InterpreterAction
}
type Gas struct {
	Limit     uint64
	Remaining uint64
	Refunded  int64
}
type Contract struct {
	Input           Bytes
	Bytecode        Bytecode
	Hash            *B256
	TargetAddress   Address
	BytecodeAddress *Address
	Caller          Address
	CallValue       U256
}
type InstructionResult uint8
const (
	// Success codes
	Continue InstructionResult = iota // 0x00
	Stop
	Return
	SelfDestruct
	ReturnContract

	// Revert codes
	Revert = 0x10 // revert opcode
	CallTooDeep
	OutOfFunds
	CreateInitCodeStartingEF00
	InvalidEOFInitCode
	InvalidExtDelegateCallTarget

	// Actions
	CallOrCreate = 0x20

	// Error codes
	OutOfGas = 0x50
	MemoryOOG
	MemoryLimitOOG
	PrecompileOOG
	InvalidOperandOOG
	OpcodeNotFound
	CallNotAllowedInsideStatic
	StateChangeDuringStaticCall
	InvalidFEOpcode
	InvalidJump
	NotActivated
	StackUnderflow
	StackOverflow
	OutOfOffset
	CreateCollision
	OverflowPayment
	PrecompileError
	NonceOverflow
	CreateContractSizeLimit
	CreateContractStartingWithEF
	CreateInitCodeSizeLimit
	FatalExternalError
	ReturnContractInNotInitEOF
	EOFOpcodeDisabledInLegacy
	EOFFunctionStackOverflow
	EofAuxDataOverflow
	EofAuxDataTooSmall
	InvalidEXTCALLTarget
)
// MarshalJSON customizes the JSON serialization of InstructionResult.
func (ir InstructionResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(ir))
}

// UnmarshalJSON customizes the JSON deserialization of InstructionResult.
func (ir *InstructionResult) UnmarshalJSON(data []byte) error {
	var value uint8
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*ir = InstructionResult(value)
	return nil
}
type SharedMemory struct {
	Buffer         []byte
	Checkpoints    []int
	LastCheckpoint int
	MemoryLimit    uint64 //although implemented as feature memory_limit but kept p=optional for now
}
func (s *SharedMemory) NewContext() {
    newCheckpoint := len(s.Buffer)
    s.Checkpoints = append(s.Checkpoints, newCheckpoint)
    s.LastCheckpoint = newCheckpoint
}
func (s SharedMemory) New() SharedMemory {
    return s.WithCapacity(4 * 1024)
}
func (s SharedMemory) WithCapacity(capacity uint) SharedMemory {
    return SharedMemory{
        Buffer: make([]byte, capacity),
        Checkpoints: make([]int, 32),
        LastCheckpoint: 0,
        MemoryLimit: math.MaxUint64,
    }
}

func (s SharedMemory) NewWithMemoryLimit(memoryLimit uint64) SharedMemory {
    return SharedMemory{
        Buffer: make([]byte, 4 * 1024),
        Checkpoints: make([]int, 32),
        LastCheckpoint: 0,
        MemoryLimit: memoryLimit,
    }
}
func (s *SharedMemory) FreeContext() {
    if len(s.Checkpoints) > 0 {
        // Pop the last element from Checkpoints
        oldCheckpoint := s.Checkpoints[len(s.Checkpoints)-1]
        s.Checkpoints = s.Checkpoints[:len(s.Checkpoints)-1]

        // Update LastCheckpoint to the last value in Checkpoints, or 0 if empty
        if len(s.Checkpoints) > 0 {
            s.LastCheckpoint = s.Checkpoints[len(s.Checkpoints)-1]
        } else {
            s.LastCheckpoint = 0 // default value if no checkpoints are left
        }

        // Safely resize Buffer to the old checkpoint length
        if oldCheckpoint < len(s.Buffer) {
            s.Buffer = s.Buffer[:oldCheckpoint]
        }
    }
}
type Stack struct {
	Data []U256
}
type FunctionStack struct {
	ReturnStack    []FunctionReturnFrame `json:"return_stack"`
	CurrentCodeIdx uint64                `json:"current_code_idx"` // Changed to uint64 to match Rust's usize
}

// MarshalJSON customizes the JSON serialization of FunctionStack.
func (fs FunctionStack) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ReturnStack    []FunctionReturnFrame `json:"return_stack"`
		CurrentCodeIdx uint64                `json:"current_code_idx"`
	}{
		ReturnStack:    fs.ReturnStack,
		CurrentCodeIdx: fs.CurrentCodeIdx,
	})
}

// UnmarshalJSON customizes the JSON deserialization of FunctionStack.
func (fs *FunctionStack) UnmarshalJSON(data []byte) error {
	var temp struct {
		ReturnStack    []FunctionReturnFrame `json:"return_stack"`
		CurrentCodeIdx uint64                `json:"current_code_idx"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	fs.ReturnStack = temp.ReturnStack
	fs.CurrentCodeIdx = temp.CurrentCodeIdx
	return nil
}

// FunctionReturnFrame represents a frame in the function return stack for the EVM interpreter.
type FunctionReturnFrame struct {
	// The index of the code container that this frame is executing.
	Idx uint64 `json:"idx"` // Changed to uint64 to match Rust's usize
	// The program counter where frame execution should continue.
	PC uint64 `json:"pc"` // Changed to uint64 to match Rust's usize
}

// MarshalJSON customizes the JSON serialization of FunctionReturnFrame.
func (frf FunctionReturnFrame) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Idx uint64 `json:"idx"`
		PC  uint64 `json:"pc"`
	}{
		Idx: frf.Idx,
		PC:  frf.PC,
	})
}

// UnmarshalJSON customizes the JSON deserialization of FunctionReturnFrame.
func (frf *FunctionReturnFrame) UnmarshalJSON(data []byte) error {
	var temp struct {
		Idx uint64 `json:"idx"`
		PC  uint64 `json:"pc"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	frf.Idx = temp.Idx
	frf.PC = temp.PC
	return nil
}

