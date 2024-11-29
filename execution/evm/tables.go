package evm
// Instruction is a function type that takes an Interpreter and a generic type H.
type Instruction[H any] func(interpreter *Interpreter, h *H)

// InstructionTable is a list of 256 instructions mapped to EVM opcodes.
type InstructionTable[H any] [256]Instruction[H]

// DynInstruction is a function type signature for dynamic instructions.
type DynInstruction[H any] func(interpreter *Interpreter, h *H)

// BoxedInstruction wraps a DynInstruction in a pointer, enabling dynamic dispatch.
type BoxedInstruction[H any] *DynInstruction[H]

// BoxedInstructionTable is an array of 256 boxed instructions.
type BoxedInstructionTable[H any] [256]BoxedInstruction[H]

const (
	PlainTableMode = iota
	BoxedTableMode
)
type InstructionTables[H any] struct {
	PlainTable *InstructionTable[H]
	BoxedTable *BoxedInstructionTable[H]
	Mode       int // Indicates which table is in use (Plain or Boxed)
}
//Mode contains 0 and 1
// NewPlainInstructionTable creates an InstructionTables instance with a PlainTable.
func NewPlainInstructionTable[H any](table InstructionTable[H]) InstructionTables[H] {
    return InstructionTables[H]{PlainTable: &table}
}

// NewBoxedInstructionTable creates an InstructionTables instance with a BoxedTable.
func NewBoxedInstructionTable[H any](table BoxedInstructionTable[H]) InstructionTables[H] {
    return InstructionTables[H]{BoxedTable: &table}
}
