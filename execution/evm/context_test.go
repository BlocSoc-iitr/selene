package evm

import (
    "testing"
    "math/big"

    "github.com/stretchr/testify/assert"
)

type MockDB struct {}

func (db *MockDB) Basic(address Address) (AccountInfo, error) {
    return AccountInfo{}, nil // Mock implementation
}

func (db *MockDB) BlockHash(number uint64) (B256, error) {
    return B256{}, nil // Mock implementation
}

func (db *MockDB) Storage(address Address, index U256) (U256, error) {
    return big.NewInt(0), nil // Mock implementation
}

func (db *MockDB) CodeByHash(codeHash B256) (Bytecode, error) {
    return Bytecode{}, nil // Mock implementation
}

func TestDefaultContext(t *testing.T) {
    ctx := DefaultContext[interface{}]()

    assert.NotNil(t, ctx)
    assert.IsType(t, Context[interface{}, *EmptyDB]{}, ctx)
    assert.NotNil(t, ctx.Evm)
    assert.Nil(t, ctx.External)
}

func TestNewContext(t *testing.T) {
    mockDB := &MockDB{}
    evmCtx := NewEvmContext(mockDB)
    externalData := struct{}{} // Example external data

    ctx := NewContext[interface{}, *MockDB](evmCtx, externalData)

    assert.NotNil(t, ctx)
    assert.Equal(t, evmCtx, ctx.Evm)
    assert.Equal(t, externalData, ctx.External)
}

func TestContextWithHandlerCfg(t *testing.T) {
    mockDB := &MockDB{}
    evmCtx := NewEvmContext(mockDB)
    externalData := struct{}{}
    cfg := HandlerCfg{} // Assuming you have a HandlerCfg struct

    cWithHandler := ContextWithHandlerCfg[interface{}, *MockDB]{
        Context: NewContext[interface{}, *MockDB](evmCtx, externalData),
        Cfg:     cfg,
    }

    assert.NotNil(t, cWithHandler)
    assert.Equal(t, cfg, cWithHandler.Cfg)
    assert.Equal(t, externalData, cWithHandler.Context.External)
}

// func TestSetSpecId(t *testing.T) {
//     js := NewJournalState(SpecId(0), nil) // Initialize with some SpecId
//     specId := SpecId(5) // Set a specific value for SpecId

//     js.SetSpecId(specId)

//     assert.Equal(t, specId, js.Spec) // This should now pass
// }
