package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"time"
)

// NodeError is the base error type for all node-related errors
type NodeError struct {
	errType string
	message string
	cause   error
}

func (e *NodeError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s", e.errType, e.cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.errType, e.message)
}
func (e *NodeError) Unwrap() error {
	return e.cause
}

// NewExecutionEvmError creates a new ExecutionEvmError
func NewExecutionEvmError(err error) *NodeError {
	return &NodeError{
		errType: "ExecutionEvmError",
		cause:   err,
	}
}

// NewExecutionError creates a new ExecutionError
func NewExecutionError(err error) *NodeError {
	return &NodeError{
		errType: "ExecutionError",
		message: "execution error",
		cause:   err,
	}
}

// NewOutOfSyncError creates a new OutOfSyncError
func NewOutOfSyncError(behind time.Duration) *NodeError {
	return &NodeError{
		errType: "OutOfSyncError",
		message: fmt.Sprintf("out of sync: %s behind", behind),
	}
}

// NewConsensusPayloadError creates a new ConsensusPayloadError
func NewConsensusPayloadError(err error) *NodeError {
	return &NodeError{
		errType: "ConsensusPayloadError",
		message: "consensus payload error",
		cause:   err,
	}
}

// NewExecutionPayloadError creates a new ExecutionPayloadError
func NewExecutionPayloadError(err error) *NodeError {
	return &NodeError{
		errType: "ExecutionPayloadError",
		message: "execution payload error",
		cause:   err,
	}
}

// NewConsensusClientCreationError creates a new ConsensusClientCreationError
func NewConsensusClientCreationError(err error) *NodeError {
	return &NodeError{
		errType: "ConsensusClientCreationError",
		message: "consensus client creation error",
		cause:   err,
	}
}

// NewExecutionClientCreationError creates a new ExecutionClientCreationError
func NewExecutionClientCreationError(err error) *NodeError {
	return &NodeError{
		errType: "ExecutionClientCreationError",
		message: "execution client creation error",
		cause:   err,
	}
}

// NewConsensusAdvanceError creates a new ConsensusAdvanceError
func NewConsensusAdvanceError(err error) *NodeError {
	return &NodeError{
		errType: "ConsensusAdvanceError",
		message: "consensus advance error",
		cause:   err,
	}
}

// NewConsensusSyncError creates a new ConsensusSyncError
func NewConsensusSyncError(err error) *NodeError {
	return &NodeError{
		errType: "ConsensusSyncError",
		message: "consensus sync error",
		cause:   err,
	}
}

// NewBlockNotFoundError creates a new BlockNotFoundError
func NewBlockNotFoundError(err error) *NodeError {
	return &NodeError{
		errType: "BlockNotFoundError",
		cause:   err,
	}
}

// EvmError represents errors from the EVM
type EvmError struct {
	revertData []byte
}

func (e *EvmError) Error() string {
	return "EVM error"
}

// DecodeRevertReason attempts to decode the revert reason from the given data
func DecodeRevertReason(data []byte) string {
	// This is a placeholder. In a real implementation, you'd decode the ABI-encoded revert reason.
	return string(data)
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToJSONRPCError converts a NodeError to a JSON-RPC error
func (e *NodeError) ToJSONRPCError() *JSONRPCError {
	switch e.errType {
	case "ExecutionEvmError":
		if evmErr, ok := e.cause.(*EvmError); ok {
			if evmErr.revertData != nil {
				msg := "execution reverted"
				if reason := DecodeRevertReason(evmErr.revertData); reason != "" {
					msg = fmt.Sprintf("%s: %s", msg, reason)
				}
				return &JSONRPCError{
					Code:    3, // Assuming 3 is the code for execution revert
					Message: msg,
					Data:    hexutil.Encode(evmErr.revertData),
				}
			}
		}
		return &JSONRPCError{
			Code:    -32000, // Generic server error
			Message: e.Error(),
		}
	default:
		return &JSONRPCError{
			Code:    -32000, // Generic server error
			Message: e.Error(),
		}
	}
}