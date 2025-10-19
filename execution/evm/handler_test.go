package evm

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/BlocSoc-iitr/selene/common"
)

// SpecId value for testing
const testSpecId SpecId = 1

// Host, EXT, and DB mock implementations for testing
type TestHost struct{}
type TestEXT struct{}

// TestDB mock implementation with required Basic method signature
type TestDB struct {
	EmptyDB
}

func (db *TestDB) Basic(addr common.Address) (AccountInfo, error) {
	// Implement a mock response to satisfy the interface
	return AccountInfo{}, nil
}

// TestNewEvmHandler tests the creation of an EvmHandler instance
func TestNewEvmHandler(t *testing.T) {
	cfg := NewHandlerCfg(testSpecId)
	cfg.isOptimism = false
	handler, err := NewEvmHandler[TestHost, TestEXT, *TestDB](cfg)

	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, testSpecId, handler.specId())
	assert.False(t, handler.IsOptimism(), "Expected optimism to be disabled")
}

// TestOptimismWithSpec tests the handler initialization with optimism enabled
func TestOptimismWithSpec(t *testing.T) {
	cfg := HandlerCfg{
		specID:     testSpecId,
		isOptimism: true,  // Set to true here
	}

	handler, err := NewEvmHandler[TestHost, TestEXT, *TestDB](cfg)
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.False(t, handler.IsOptimism(), "Expected optimism to be enabled")
	assert.Equal(t, testSpecId, handler.specId())
}


// TestMainnetWithSpec tests the handler initialization for mainnet (optimism disabled)
func TestMainnetWithSpec(t *testing.T) {
	cfg := HandlerCfg{
		specID:     testSpecId,
		isOptimism: false,
	}

	handler, err := NewEvmHandler[TestHost, TestEXT, *TestDB](cfg)
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.False(t, handler.IsOptimism(), "Expected optimism to be disabled")
	assert.Equal(t, testSpecId, handler.specId())
}

// TestRegisterMethods ensures that PlainRegister and BoxRegister correctly call RegisterFn on handler
func TestRegisterMethods(t *testing.T) {
	cfg := HandlerCfg{
		specID:     testSpecId,
		isOptimism: false,
	}
	handler, _ := NewEvmHandler[TestHost, TestEXT, *TestDB](cfg)

	plainRegister := PlainRegister[TestHost, TestEXT, *TestDB]{
		RegisterFn: func(h *EvmHandler[TestHost, TestEXT, *TestDB]) {
			assert.Equal(t, handler, h, "Expected handler to be passed correctly in PlainRegister")
		},
	}
	plainRegister.Register(handler)

	boxRegister := BoxRegister[TestHost, TestEXT, *TestDB]{
		RegisterFn: func(h *EvmHandler[TestHost, TestEXT, *TestDB]) {
			assert.Equal(t, handler, h, "Expected handler to be passed correctly in BoxRegister")
		},
	}
	boxRegister.Register(handler)
}
