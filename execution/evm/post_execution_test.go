package evm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementations for testing
func mockReimburseCaller(ctx *Context[any, Database], gas *Gas) EVMResultGeneric[struct{}, any] {
	return EVMResultGeneric[struct{}, any]{Value: struct{}{}, Err: nil}
}

func mockRewardBeneficiary(ctx *Context[any, Database], gas *Gas) EVMResultGeneric[struct{}, any] {
	return EVMResultGeneric[struct{}, any]{Value: struct{}{}, Err: nil}
}

func mockOutput(ctx *Context[any, Database], frameResult FrameResult) (ResultAndState, EvmError) {
	return ResultAndState{}, EvmError{}
}

func mockEnd(ctx *Context[any, Database], result EVMResultGeneric[ResultAndState, EvmError]) (ResultAndState, EvmError) {
	return ResultAndState{}, EvmError{}
}

func mockClear(ctx *Context[any, Database]) {}

// Test the PostExecutionHandler with mock functions
func TestPostExecutionHandler(t *testing.T) {
	handler := PostExecutionHandler[any, Database]{
		ReimburseCaller:   mockReimburseCaller,
		RewardBeneficiary: mockRewardBeneficiary,
		Output:            mockOutput,
		End:               mockEnd,
		Clear:             mockClear,
	}

	ctx := &Context[any, Database]{}
	gas := &Gas{}

	// Test ReimburseCaller
	t.Run("ReimburseCaller Success", func(t *testing.T) {
		result := handler.ReimburseCaller(ctx, gas)
		assert.NoError(t, result.Err, "ReimburseCaller should not return an error")
	})

	// Test RewardBeneficiary
	t.Run("RewardBeneficiary Success", func(t *testing.T) {
		result := handler.RewardBeneficiary(ctx, gas)
		assert.NoError(t, result.Err, "RewardBeneficiary should not return an error")
	})

	// Test Output
	t.Run("Output Success", func(t *testing.T) {
		_, err := handler.Output(ctx, FrameResult{})
		assert.Empty(t, err.Message, "Output should not return an error")
	})

	// Create a new result for End, converting to the expected type
	endResult := EVMResultGeneric[ResultAndState, EvmError]{Value: ResultAndState{}}

	// Test End
	t.Run("End Success", func(t *testing.T) {
		_, err := handler.End(ctx, endResult)
		assert.Empty(t, err.Message, "End should not return an error")
	})

	// Test Clear
	t.Run("Clear Success", func(t *testing.T) {
		assert.NotPanics(t, func() { handler.Clear(ctx) }, "Clear should not panic")
	})
}

// Test with errors
func TestPostExecutionHandlerWithError(t *testing.T) {
	handler := PostExecutionHandler[any, Database]{
		ReimburseCaller: func(ctx *Context[any, Database], gas *Gas) EVMResultGeneric[struct{}, any] {
			return EVMResultGeneric[struct{}, any]{Err: errors.New("mock error")}
		},
		RewardBeneficiary: func(ctx *Context[any, Database], gas *Gas) EVMResultGeneric[struct{}, any] {
			return EVMResultGeneric[struct{}, any]{Err: errors.New("mock reward error")}
		},
		Output: func(ctx *Context[any, Database], frameResult FrameResult) (ResultAndState, EvmError) {
			return ResultAndState{}, EvmError{Message: "mock output error"}
		},
		End: func(ctx *Context[any, Database], result EVMResultGeneric[ResultAndState, EvmError]) (ResultAndState, EvmError) {
			return ResultAndState{}, EvmError{Message: "mock end error"}
		},
		Clear: mockClear,
	}

	ctx := &Context[any, Database]{}
	gas := &Gas{}

	// Test ReimburseCaller with an error
	t.Run("ReimburseCaller Error", func(t *testing.T) {
		result := handler.ReimburseCaller(ctx, gas)
		assert.Error(t, result.Err, "Expected error from ReimburseCaller")
	})

	// Test RewardBeneficiary with an error
	t.Run("RewardBeneficiary Error", func(t *testing.T) {
		result := handler.RewardBeneficiary(ctx, gas)
		assert.Error(t, result.Err, "Expected error from RewardBeneficiary")
	})

	// Test Output with a mock error
	t.Run("Output Error", func(t *testing.T) {
		_, err := handler.Output(ctx, FrameResult{})
		assert.Error(t, err, "Expected error from Output")
	})

	// Test End with a mock error
	t.Run("End Error", func(t *testing.T) {
		_, err := handler.End(ctx, EVMResultGeneric[ResultAndState, EvmError]{Value: ResultAndState{}})
		assert.Error(t, err, "Expected error from End")
	})
}

