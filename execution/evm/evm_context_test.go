package evm

import (
	"testing"
	"math/big"
)

// MockDatabase is a mock implementation of the Database interface for testing.
type MockDatabase struct {
    BasicFunc     func(Address) (AccountInfo, error)
    BlockHashFunc func(uint64) (B256, error)
    StorageFunc   func(Address, U256) (U256, error)
    CodeByHashFunc func(B256) (Bytecode, error)
}

func (mdb *MockDatabase) Basic(address Address) (AccountInfo, error) {
    return mdb.BasicFunc(address)
}

func (mdb *MockDatabase) BlockHash(number uint64) (B256, error) {
    return mdb.BlockHashFunc(number)
}

func (mdb *MockDatabase) Storage(address Address, index U256) (U256, error) {
    return mdb.StorageFunc(address, index)
}

func (mdb *MockDatabase) CodeByHash(codeHash B256) (Bytecode, error) {
    return mdb.CodeByHashFunc(codeHash)
}

// TestNewEvmContext tests the NewEvmContext function.
func TestNewEvmContext(t *testing.T) {
	db := &MockDatabase{
		BasicFunc: func(address Address) (AccountInfo, error) {
			return AccountInfo{}, nil
		},
		BlockHashFunc: func(number uint64) (B256, error) {
			return B256{}, nil
		},
		StorageFunc: func(address Address, index U256) (U256, error) {
			return big.NewInt(0), nil
		},
		CodeByHashFunc: func(codeHash B256) (Bytecode, error) {
			return Bytecode{}, nil
		},
	}

	// Create a new EvmContext
	evmCtx := NewEvmContext(db)

	// Check if the Inner and Precompiles are initialized correctly
	if evmCtx.Inner.DB != db {
		t.Errorf("Expected DB to be the provided mock database, got: %v", evmCtx.Inner.DB)
	}
	if evmCtx.Precompiles.Inner.Owned == nil {
		t.Error("Expected Precompiles Owned to be initialized, got nil")
	}
}

// TestWithNewEvmDB tests the WithNewEvmDB function.
func TestWithNewEvmDB(t *testing.T) {
	db1 := &MockDatabase{
		BasicFunc: func(address Address) (AccountInfo, error) {
			return AccountInfo{}, nil
		},
	}

	db2 := &MockDatabase{
		BasicFunc: func(address Address) (AccountInfo, error) {
			return AccountInfo{}, nil
		},
	}

	// Create a new EvmContext with the first DB
	evmCtx1 := NewEvmContext(db1)

	// Create a new EvmContext with the second DB
	evmCtx2 := WithNewEvmDB(evmCtx1, db2)

	// Check if the new EvmContext's DB is the second DB
	if evmCtx2.Inner.DB != db2 {
		t.Errorf("Expected DB to be the second mock database, got: %v", evmCtx2.Inner.DB)
	}
}


