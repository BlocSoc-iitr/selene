package evm

import (

	"testing"
)

// TestNewCallAction verifies that NewCallAction initializes a Call action with the correct inputs.
func TestNewCallAction(t *testing.T) {
	callInputs := &CallInputs{}
	action := NewCallAction(callInputs)

	if action.actionType != ActionTypeCall {
		t.Errorf("Expected action type %v, got %v", ActionTypeCall, action.actionType)
	}
	if action.callInputs != callInputs {
		t.Errorf("Expected callInputs to be %v, got %v", callInputs, action.callInputs)
	}
}

// TestNewCreateAction verifies that NewCreateAction initializes a Create action with the correct inputs.
func TestNewCreateAction(t *testing.T) {
	createInputs := &CreateInputs{}
	action := NewCreateAction(createInputs)

	if action.actionType != ActionTypeCreate {
		t.Errorf("Expected action type %v, got %v", ActionTypeCreate, action.actionType)
	}
	if action.createInputs != createInputs {
		t.Errorf("Expected createInputs to be %v, got %v", createInputs, action.createInputs)
	}
}

// TestNewEOFCreateAction verifies that NewEOFCreateAction initializes an EOFCreate action with the correct inputs.
func TestNewEOFCreateAction(t *testing.T) {
	eofCreateInputs := &EOFCreateInputs{}
	action := NewEOFCreateAction(eofCreateInputs)

	if action.actionType != ActionTypeEOFCreate {
		t.Errorf("Expected action type %v, got %v", ActionTypeEOFCreate, action.actionType)
	}
	if action.eofCreateInputs != eofCreateInputs {
		t.Errorf("Expected eofCreateInputs to be %v, got %v", eofCreateInputs, action.eofCreateInputs)
	}
}

// TestNewReturnAction verifies that NewReturnAction initializes a Return action with the correct result.
func TestNewReturnAction(t *testing.T) {
	result := &InterpreterResult{}
	action := NewReturnAction(result)

	if action.actionType != ActionTypeReturn {
		t.Errorf("Expected action type %v, got %v", ActionTypeReturn, action.actionType)
	}
	if action.result != result {
		t.Errorf("Expected result to be %v, got %v", result, action.result)
	}
}

// TestNewNoneAction verifies that NewNoneAction initializes a None action.
func TestNewNoneAction(t *testing.T) {
	action := NewNoneAction()

	if action.actionType != ActionTypeNone {
		t.Errorf("Expected action type %v, got %v", ActionTypeNone, action.actionType)
	}
	if action.callInputs != nil || action.createInputs != nil || action.eofCreateInputs != nil || action.result != nil {
		t.Error("Expected inputs and result fields to be nil")
	}
}

// TestActionTypeEnum checks that action types are assigned correctly.
func TestActionTypeEnum(t *testing.T) {
	expected := []ActionType{
		ActionTypeNone,
		ActionTypeCall,
		ActionTypeCreate,
		ActionTypeEOFCreate,
		ActionTypeReturn,
	}

	for i, v := range expected {
		if ActionType(i) != v {
			t.Errorf("Expected ActionType %v at index %d, but got %v", v, i, ActionType(i))
		}
	}
}

// Helper structs for testing (placeholders to allow test compilation).
// Replace these with real struct definitions as needed.
