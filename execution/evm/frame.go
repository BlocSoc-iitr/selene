package evm
type CallFrame struct {
	ReturnMemoryRange Range
	FrameData         FrameData
}
type CreateFrame struct {
	CreatedAddress Address
	FrameData      FrameData
}

type EOFCreateFrame struct {
	CreatedAddress Address
	FrameData      FrameData
}

type FrameData struct {
	Checkpoint  JournalCheckpoint
	Interpreter Interpreter
}

type CallOutcome struct {
	Result       InterpreterResult
	MemoryOffset Range
}

type CreateOutcome struct {
	Result  InterpreterResult
	Address *Address
}

type FrameType uint8

const (
	Call_Frame FrameType = iota
	Create_Frame
	EOFCreate_Frame
)

type Frame struct {
	Type      FrameType
	Call      *CallFrame
	Create    *CreateFrame
	EOFCreate *EOFCreateFrame
}

// FrameResultType represents the type of FrameResult.
type FrameResultType int

const (
	Call FrameResultType = iota
	Create
	EOFCreate
)

// FrameResult represents the outcome of a frame, which can be Call, Create, or EOFCreate.
type FrameResult struct {
	ResultType FrameResultType
	Call       *CallOutcome
	Create     *CreateOutcome
}

type FrameOrResultType uint8

const (
	Frame_FrameOrResult FrameOrResultType = iota
	Result_FrameOrResult
)

type FrameOrResult struct {
	Type   FrameOrResultType
	Frame  Frame
	Result FrameResult
}

func (f *FrameResult) Gas() *Gas {
	if f.ResultType == Call {
		return &f.Call.Result.Gas
	} else if f.ResultType == Create {
		return &f.Create.Result.Gas
	} else {
		return &f.Create.Result.Gas
	}
}
type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
}