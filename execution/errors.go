package execution

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ExecutionError represents various execution-related errors
type ExecutionError struct {
	Kind    string
	Details interface{}
}

func (e *ExecutionError) Error() string {
	switch e.Kind {
	case "InvalidAccountProof":
		return fmt.Sprintf("invalid account proof for string: %v", e.Details)
	case "InvalidStorageProof":
		details := e.Details.([]interface{})
		return fmt.Sprintf("invalid storage proof for string: %v, slot: %v", details[0], details[1])
	case "CodeHashMismatch":
		details := e.Details.([]interface{})
		return fmt.Sprintf("code hash mismatch for string: %v, found: %v, expected: %v", details[0], details[1], details[2])
	case "ReceiptRootMismatch":
		return fmt.Sprintf("receipt root mismatch for tx: %v", e.Details)
	case "MissingTransaction":
		return fmt.Sprintf("missing transaction for tx: %v", e.Details)
	case "NoReceiptForTransaction":
		return fmt.Sprintf("could not prove receipt for tx: %v", e.Details)
	case "MissingLog":
		details := e.Details.([]interface{})
		return fmt.Sprintf("missing log for transaction: %v, index: %v", details[0], details[1])
	case "TooManyLogsToProve":
		details := e.Details.([]interface{})
		return fmt.Sprintf("too many logs to prove: %v, current limit is: %v", details[0], details[1])
	case "IncorrectRpcNetwork":
		return "execution RPC is for the incorrect network"
	case "InvalidBaseGasFee":
		details := e.Details.([]interface{})
		return fmt.Sprintf("Invalid base gas fee selene %v vs rpc endpoint %v at block %v", details[0], details[1], details[2])
	case "InvalidGasUsedRatio":
		details := e.Details.([]interface{})
		return fmt.Sprintf("Invalid gas used ratio of selene %v vs rpc endpoint %v at block %v", details[0], details[1], details[2])
	case "BlockNotFoundError":
		return fmt.Sprintf("Block %v not found", e.Details)
	case "EmptyExecutionPayload":
		return "Selene Execution Payload is empty"
	case "InvalidBlockRange":
		details := e.Details.([]interface{})
		return fmt.Sprintf("User query for block %v but selene oldest block is %v", details[0], details[1])
	default:
		return "unknown execution error"
	}
}

// Helper functions to create specific ExecutionError instances
func NewInvalidAccountProofError(string string) error {
	return &ExecutionError{"InvalidAccountProof", string}
}

func NewInvalidStorageProofError(string string, slot [32]byte) error {
	return &ExecutionError{"InvalidStorageProof", []interface{}{string, slot}}
}

func NewCodeHashMismatchError(string string, found [32]byte, expected [32]byte) error {
	return &ExecutionError{"CodeHashMismatch", []interface{}{string, found, expected}}
}

func NewReceiptRootMismatchError(tx [32]byte) error {
	return &ExecutionError{"ReceiptRootMismatch", tx}
}

func NewMissingTransactionError(tx [32]byte) error {
	return &ExecutionError{"MissingTransaction", tx}
}

func NewNoReceiptForTransactionError(tx [32]byte) error {
	return &ExecutionError{"NoReceiptForTransaction", tx}
}

func NewMissingLogError(tx [32]byte, index uint64) error {
	return &ExecutionError{"MissingLog", []interface{}{tx, index}}
}

func NewTooManyLogsToProveError(count int, limit int) error {
	return &ExecutionError{"TooManyLogsToProve", []interface{}{count, limit}}
}

func NewIncorrectRpcNetworkError() error {
	return &ExecutionError{"IncorrectRpcNetwork", nil}
}

func NewInvalidBaseGasFeeError(selene uint64, rpc uint64, block uint64) error {
	return &ExecutionError{"InvalidBaseGasFee", []interface{}{selene, rpc, block}}
}

func NewInvalidGasUsedRatioError(seleneRatio float64, rpcRatio float64, block uint64) error {
	return &ExecutionError{"InvalidGasUsedRatio", []interface{}{seleneRatio, rpcRatio, block}}
}

func NewBlockNotFoundError(block uint64) error {
	return &ExecutionError{"BlockNotFoundError", block}
}

func NewEmptyExecutionPayloadError() error {
	return &ExecutionError{"EmptyExecutionPayload", nil}
}

func NewInvalidBlockRangeError(queryBlock uint64, oldestBlock uint64) error {
	return &ExecutionError{"InvalidBlockRange", []interface{}{queryBlock, oldestBlock}}
}

// EvmError represents EVM-related errors
type EvmError struct {
	Kind    string
	Details interface{}
}

func (e *EvmError) Error() string {
	switch e.Kind {
	case "Revert":
		return fmt.Sprintf("execution reverted: %v", e.Details)
	case "Generic":
		return fmt.Sprintf("evm error: %v", e.Details)
	case "RpcError":
		return fmt.Sprintf("rpc error: %v", e.Details)
	default:
		return "unknown evm error"
	}
}

// Helper functions for creating specific EVM errors
func NewRevertError(data []byte) error {
	return &EvmError{"Revert", data}
}

func NewGenericError(message string) error {
	return &EvmError{"Generic", message}
}

func NewRpcError(report error) error {
	return &EvmError{"RpcError", report}
}

func DecodeRevertReason(data []byte) string {
	reason, err := abi.UnpackRevert(data)
	if err != nil {
		reason = string(err.Error())
	}
	return reason
}
