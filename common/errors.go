package common

import "fmt"

type BlockNotFoundError struct {
	Block BlockTag
}

func NewBlockNotFoundError(block BlockTag) BlockNotFoundError {
	return BlockNotFoundError{Block: block}
}

func (e BlockNotFoundError) Error() string {
	return fmt.Sprintf("block not available: %s", e.Block)
}

// need to confirm how such primitive types will be imported
type hash [32]byte

type SlotNotFoundError struct {
	slot hash
}

func NewSlotNotFoundError(slot hash) SlotNotFoundError {
	return SlotNotFoundError{slot: slot}
}

func (e SlotNotFoundError) Error() string {
	return fmt.Sprintf("slot not available: %s", e.slot)
}

type RpcError struct {
	method string
	error  error
}

func NewRpcError(method string, err error) RpcError {
	return RpcError{
		method: method,
		error:  err,
	}
}

func (e RpcError) Error() string {
	return fmt.Sprintf("rpc error on method: %s, message: %s", e.method, e.error.Error())
}
