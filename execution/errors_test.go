package execution

import (
	"errors"
	"testing"

	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/BlocSoc-iitr/selene/consensus/types"
	"github.com/stretchr/testify/assert"
)

func TestExecutionErrors(t *testing.T) {
	// Test InvalidAccountProofError
	address := types.Address{0x01, 0x02}
	err := NewInvalidAccountProofError(address)
	assert.EqualError(t, err, "invalid account proof for string: [1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]")

	// Test InvalidStorageProofError
	slot := consensus_core.Bytes32{0x0a}
	err = NewInvalidStorageProofError(address, slot)
	assert.EqualError(t, err, "invalid storage proof for string: [1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0], slot: [10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]")

	// Test CodeHashMismatchError
	found := consensus_core.Bytes32{0x03}
	expected := consensus_core.Bytes32{0x04}
	err = NewCodeHashMismatchError(address, found, expected)
	assert.EqualError(t, err, "code hash mismatch for string: [1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0], found: [3 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0], expected: [4 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]")

	// Test ReceiptRootMismatchError
	tx := consensus_core.Bytes32{0x05}
	err = NewReceiptRootMismatchError(tx)
	assert.EqualError(t, err, "receipt root mismatch for tx: [5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]")

	// Test MissingTransactionError
	err = NewMissingTransactionError(tx)
	assert.EqualError(t, err, "missing transaction for tx: [5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]")

	// Test MissingLogError
	err = NewMissingLogError(tx, 3)
	assert.EqualError(t, err, "missing log for transaction: [5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0], index: 3")

	// Test TooManyLogsToProveError
	err = NewTooManyLogsToProveError(5000, 1000)
	assert.EqualError(t, err, "too many logs to prove: 5000, current limit is: 1000")

	// Test InvalidBaseGasFeeError
	err = NewInvalidBaseGasFeeError(1000, 2000, 123456)
	assert.EqualError(t, err, "Invalid base gas fee selene 1000 vs rpc endpoint 2000 at block 123456")

	// Test BlockNotFoundError
	err = NewBlockNotFoundError(123456)
	assert.EqualError(t, err, "Block 123456 not found")

	// Test EmptyExecutionPayloadError
	err = NewEmptyExecutionPayloadError()
	assert.EqualError(t, err, "Selene Execution Payload is empty")
}

func TestEvmErrors(t *testing.T) {
	// Test RevertError
	data := []byte{0x08, 0xc3, 0x79, 0xa0}
	err := NewRevertError(data)
	assert.EqualError(t, err, "execution reverted: [8 195 121 160]")

	// Test GenericError
	err = NewGenericError("generic error")
	assert.EqualError(t, err, "evm error: generic error")

	// Test RpcError
	rpcErr := errors.New("rpc connection failed")
	err = NewRpcError(rpcErr)
	assert.EqualError(t, err, "rpc error: rpc connection failed")
}

func TestDecodeRevertReason(t *testing.T) {
	// Test successful revert reason decoding
	reasonData := []byte{0x08, 0xc3, 0x79, 0xa0}
	reason := DecodeRevertReason(reasonData)
	assert.NotEmpty(t, reason, "Revert reason should be decoded")

	// Test invalid revert data
	invalidData := []byte{0x00}
	reason = DecodeRevertReason(invalidData)
	assert.Contains(t, reason, "invalid data for unpacking")
}