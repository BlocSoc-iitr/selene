package evm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Rename the mock database to avoid conflicts
type TestMockDatabase struct{}
type MockExternal struct{}

func (db *TestMockDatabase) Basic(address Address) (AccountInfo, error) {
	return AccountInfo{}, nil
}

func (db *TestMockDatabase) BlockHash(number uint64) (B256, error) {
	return B256{}, nil
}

func (db *TestMockDatabase) Storage(address Address, index U256) (U256, error) {
	return big.NewInt(0), nil
}

func (db *TestMockDatabase) CodeByHash(codeHash B256) (Bytecode, error) {
	return Bytecode{}, nil
}

// Test for creating EvmBuilder with default configurations
func TestNewDefaultEvmBuilder(t *testing.T) {
	builder := NewDefaultEvmBuilder[MockExternal]()
	assert.NotNil(t, builder, "EvmBuilder should not be nil")
	assert.IsType(t, &EvmBuilder[MockExternal, *EmptyDB]{}, builder, "Builder type should match EvmBuilder")
}

// Test for changing the database within EvmBuildertype MockExternal struct{}

// Test for changing the database within EvmBuilder
func TestWithNewDB(t *testing.T) {
	builder := NewDefaultEvmBuilder[MockExternal]() // Use the mock external struct
	newDB := &TestMockDatabase{} // Use the renamed mock database

	// Ensure that the context is properly set up
	builder.context.External = MockExternal{} // Set the external field to an instance of MockExternal

	builderWithNewDB := WithNewDB(builder, newDB)

	assert.IsType(t, &EvmBuilder[MockExternal, *TestMockDatabase]{}, builderWithNewDB, "Builder should be updated with new DB type")
	assert.NotEqual(t, builder.context, builderWithNewDB.context, "Context should be updated with new DB")
}
// Test for setting environment in EvmBuilder
func TestWithEnv(t *testing.T) {
	builder := NewDefaultEvmBuilder[MockExternal]()
	env := &Env{} // Use *Env to match the type expected by WithEnv

	builderWithEnv := builder.WithEnv(env)
	assert.Equal(t, env, builderWithEnv.context.Evm.Inner.Env, "Environment should be set correctly in EvmBuilder")
}

// Test for building an Evm instance
func TestBuild(t *testing.T) {
	builder := NewDefaultEvmBuilder[MockExternal]()
	evm := builder.Build()

	assert.NotNil(t, evm, "Evm instance should be successfully built")
}

// Test for updating context and handler configuration
func TestWithContextWithHandlerCfg(t *testing.T) {
	builder := NewDefaultEvmBuilder[MockExternal]()
	mockContextWithHandlerCfg := ContextWithHandlerCfg[MockExternal, *TestMockDatabase]{
		Context: NewContext[MockExternal, *TestMockDatabase](NewEvmContext(&TestMockDatabase{}), MockExternal{}),
		Cfg:     NewHandlerCfg(LATEST),
	}

	builderWithNewCfg := WithContextWithHandlerCfg(builder, mockContextWithHandlerCfg)
	assert.Equal(t, mockContextWithHandlerCfg.Context, builderWithNewCfg.context, "Context should be updated correctly")
	assert.Equal(t, mockContextWithHandlerCfg.Cfg, builderWithNewCfg.handler.Cfg, "Handler configuration should be updated correctly")
}
