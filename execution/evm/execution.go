package evm

type LastFrameReturnHandle[EXT any, DB Database] func(*Context[EXT, DB], *FrameResult) EvmError
type ExecuteFrameHandle[EXT any, DB Database] func(*Frame, *SharedMemory, *InstructionTables[Context[EXT, DB]], *Context[EXT, DB]) (InterpreterAction, EvmError)
type FrameCallHandle[EXT any, DB Database] func(*Context[EXT, DB], *CallInputs) (FrameOrResult, EvmError)
type FrameCallReturnHandle[EXT any, DB Database] func(*Context[EXT, DB], *CallFrame, InterpreterResult) (CallOutcome, EvmError)
type InsertCallOutcomeHandle[EXT any, DB Database] func(*Context[EXT, DB], *Frame, *SharedMemory, CallOutcome) EvmError
type FrameCreateHandle[EXT any, DB Database] func(*Context[EXT, DB], *CreateInputs) (FrameOrResult, EvmError)
type FrameCreateReturnHandle[EXT any, DB Database] func(*Context[EXT, DB], *CreateFrame, InterpreterResult) (CreateOutcome, EvmError)
type InsertCreateOutcomeHandle[EXT any, DB Database] func(*Context[EXT, DB], *Frame, CreateOutcome) EvmError
type FrameEOFCreateHandle[EXT any, DB Database] func(*Context[EXT, DB], EOFCreateInputs) (FrameOrResult, EvmError)
type FrameEOFCreateReturnHandle[EXT any, DB Database] func(*Context[EXT, DB], *EOFCreateFrame, InterpreterResult) (CreateOutcome, EvmError)
type InsertEOFCreateOutcomeHandle[EXT any, DB Database] func(*Context[EXT, DB], *CallFrame, InterpreterResult) (CallOutcome, EvmError)

type ExecutionHandler[EXT interface{}, DB Database] struct {
	LastFrameReturn        LastFrameReturnHandle[EXT, DB]
	ExecuteFrame           ExecuteFrameHandle[EXT, DB]
	Call                   FrameCallHandle[EXT, DB]
	CallReturn             FrameCallReturnHandle[EXT, DB]
	InsertCallOutcome      InsertCallOutcomeHandle[EXT, DB]
	Create                 FrameCreateHandle[EXT, DB]
	CreateReturn           FrameCreateReturnHandle[EXT, DB]
	InsertCreateOutcome    InsertCreateOutcomeHandle[EXT, DB]
	EOFCreate              FrameEOFCreateHandle[EXT, DB]
	EOFCreateReturn        FrameEOFCreateReturnHandle[EXT, DB]
	InsertEOFCreateOutcome InsertEOFCreateOutcomeHandle[EXT, DB]
}

