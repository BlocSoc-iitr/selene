package evm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing
func mockLoadPrecompiles() ContextPrecompiles[Database] {
	return ContextPrecompiles[Database]{}
}

func mockLoadAccounts(ctx *Context[any, Database]) error {
	return nil
}

func mockDeductCaller(ctx *Context[any, Database]) EVMResultGeneric[ResultAndState, DatabaseError] {
	return EVMResultGeneric[ResultAndState, DatabaseError]{Value: ResultAndState{}, Err: nil}
}

func mockValidateInitialTxGas(env *Env) (uint64, error) {
	return 100000, nil
}

func mockTxAgainstState(ctx *Context[any, Database]) error {
	return nil
}

func mockValidateEnv(env *Env) error {
	return nil
}

// Tests for PreExecutionHandler
func TestPreExecutionHandler(t *testing.T) {
	handler := PreExecutionHandler[any, Database]{
		LoadPrecompiles: mockLoadPrecompiles,
		LoadAccounts:    mockLoadAccounts,
		DeductCaller:    mockDeductCaller,
	}

	// Test LoadPrecompiles
	precompiles := handler.LoadPrecompiles()
	assert.NotNil(t, precompiles, "LoadPrecompiles should return a valid ContextPrecompiles")

	// Test LoadAccounts
	ctx := &Context[any, Database]{}
	err := handler.LoadAccounts(ctx)
	assert.NoError(t, err, "LoadAccounts should not return an error")

	// Test DeductCaller
	result := handler.DeductCaller(ctx)
	assert.NoError(t, result.Err, "DeductCaller should not return an error")
}

// Tests for ValidationHandler
func TestValidationHandler(t *testing.T) {
	handler := ValidationHandler[any, Database]{
		InitialTxGas:   mockValidateInitialTxGas,
		TxAgainstState: mockTxAgainstState,
		Env:            mockValidateEnv,
	}

	env := &Env{}

	// Test ValidateInitialTxGas
	gas, err := handler.InitialTxGas(env)
	assert.NoError(t, err, "InitialTxGas should not return an error")
	assert.Equal(t, uint64(100000), gas, "Expected gas value should be 100000")

	// Test ValidateTxEnvAgainstState
	ctx := &Context[any, Database]{}
	err = handler.TxAgainstState(ctx)
	assert.NoError(t, err, "TxAgainstState should not return an error")

	// Test ValidateEnv
	err = handler.Env(env)
	assert.NoError(t, err, "ValidateEnv should not return an error")
}

// Tests with errors
func TestPreExecutionHandlerWithError(t *testing.T) {
	handler := PreExecutionHandler[any, Database]{
		LoadPrecompiles: mockLoadPrecompiles,
		LoadAccounts: func(ctx *Context[any, Database]) error {
			return errors.New("mock error loading accounts")
		},
		DeductCaller: mockDeductCaller,
	}

	ctx := &Context[any, Database]{}
	// Test LoadAccounts with an error
	err := handler.LoadAccounts(ctx)
	assert.Error(t, err, "Expected error from LoadAccounts")

	// Test DeductCaller with a mock error
	handler.DeductCaller = func(ctx *Context[any, Database]) EVMResultGeneric[ResultAndState, DatabaseError] {
		return EVMResultGeneric[ResultAndState, DatabaseError]{Err: errors.New("mock error deducting caller")}
	}

	result := handler.DeductCaller(ctx)
	assert.Error(t, result.Err, "Expected error from DeductCaller")
}

func TestValidationHandlerWithError(t *testing.T) {
	handler := ValidationHandler[any, Database]{
		InitialTxGas: func(env *Env) (uint64, error) {
			return 0, errors.New("mock initial gas error")
		},
		TxAgainstState: mockTxAgainstState,
		Env:            mockValidateEnv,
	}

	env := &Env{}
	// Test ValidateInitialTxGas with an error
	_, err := handler.InitialTxGas(env)
	assert.Error(t, err, "Expected error from InitialTxGas")
}