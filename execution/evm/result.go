package evm
import "github.com/ethereum/go-ethereum/core/types"
type ResultAndState struct{
	Result ExecutionResult
	State  EvmState
}
type ExecutionResult struct {
    Type         string         // "Success", "Revert", or "Halt"
    Reason       interface{}    // SuccessReason or HaltReason
    GasUsed      uint64
    GasRefunded  uint64        // Only for Success
    Logs         []types.Log   // Only for Success
    Output       Output        // Only for Success and Revert
}
type SuccessReason string
const (
    SStop              SuccessReason = "Stop"
    SReturn            SuccessReason = "Return"
    SSelfDestruct      SuccessReason = "SelfDestruct"
    SEofReturnContract SuccessReason = "EofReturnContract"
)
type Output struct {
    Type    string
    Data    []byte
    Address *Address // Only for Create type
}
type HaltReason string
const (
    HOutOfGas                    HaltReason = "OutOfGas"
    HOpcodeNotFound             HaltReason = "OpcodeNotFound"
    HInvalidFEOpcode            HaltReason = "InvalidFEOpcode"
    HInvalidJump                HaltReason = "InvalidJump"
    HNotActivated               HaltReason = "NotActivated"
    HStackUnderflow             HaltReason = "StackUnderflow"
    HStackOverflow              HaltReason = "StackOverflow"
    HOutOfOffset                HaltReason = "OutOfOffset"
    HCreateCollision            HaltReason = "CreateCollision"
    HPrecompileError            HaltReason = "PrecompileError"
    HNonceOverflow              HaltReason = "NonceOverflow"
    HCreateContractSizeLimit    HaltReason = "CreateContractSizeLimit"
    HCreateContractStartingWithEF HaltReason = "CreateContractStartingWithEF"
    HCreateInitCodeSizeLimit    HaltReason = "CreateInitCodeSizeLimit"
    HOverflowPayment           HaltReason = "OverflowPayment"
    HStateChangeDuringStaticCall HaltReason = "StateChangeDuringStaticCall"
    HCallNotAllowedInsideStatic HaltReason = "CallNotAllowedInsideStatic"
    HOutOfFunds                HaltReason = "OutOfFunds"
    HCallTooDeep               HaltReason = "CallTooDeep"
    HEofAuxDataOverflow        HaltReason = "EofAuxDataOverflow"
    HEofAuxDataTooSmall        HaltReason = "EofAuxDataTooSmall"
    HEOFFunctionStackOverflow  HaltReason = "EOFFunctionStackOverflow"
)
type EvmError struct {
    Message string
    Data    []byte
}

// Error implements the error interface
func (e EvmError) Error() string {
    return e.Message
}
