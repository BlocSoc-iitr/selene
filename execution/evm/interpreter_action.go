package evm
import(

)
type InterpreterAction struct {
	// We use a type field to distinguish between different actions
	actionType      ActionType
	callInputs      *CallInputs
	createInputs    *CreateInputs
	eofCreateInputs *EOFCreateInputs
	result          *InterpreterResult
}
type InterpreterResult struct {
	// The result of the execution.
	Result InstructionResult
	Output Bytes
	Gas    Gas
}
//Check marshal unmarshal remove if not required
// type jsonInterpreterAction struct {
// 	Type         string             `json:"type"`
// 	CallInputs   *CallInputs        `json:"call_inputs,omitempty"`
// 	CreateInputs *CreateInputs      `json:"create_inputs,omitempty"`
// 	EOFInputs    *EOFCreateInputs   `json:"eof_inputs,omitempty"`
// 	Result       *InterpreterResult `json:"result,omitempty"`
// }
// func(ja *jsonInterpreterAction) MarshalJSON()([]byte,error){
// 	return marshalJSON(ja)
// }
// // MarshalJSON implements the json.Marshaler interface

// func (ja *jsonInterpreterAction) UnmarshalJSON(data []byte) error {
// 	return unmarshalJSON(data,ja)
// }

// UnmarshalJSON implements the json.Unmarshaler interface

// ActionType represents the type of interpreter action
type ActionType int

const (
	ActionTypeNone ActionType = iota
	ActionTypeCall
	ActionTypeCreate
	ActionTypeEOFCreate
	ActionTypeReturn
)

// NewCallAction creates a new Call action
func NewCallAction(inputs *CallInputs) *InterpreterAction {
	return &InterpreterAction{
		actionType: ActionTypeCall,
		callInputs: inputs,
	}
}

// NewCreateAction creates a new Create action
func NewCreateAction(inputs *CreateInputs) *InterpreterAction {
	return &InterpreterAction{
		actionType:   ActionTypeCreate,
		createInputs: inputs,
	}
}

// NewEOFCreateAction creates a new EOFCreate action
func NewEOFCreateAction(inputs *EOFCreateInputs) *InterpreterAction {
	return &InterpreterAction{
		actionType:      ActionTypeEOFCreate,
		eofCreateInputs: inputs,
	}
}

// NewReturnAction creates a new Return action
func NewReturnAction(result *InterpreterResult) *InterpreterAction {
	return &InterpreterAction{
		actionType: ActionTypeReturn,
		result:     result,
	}
}

// NewNoneAction creates a new None action
func NewNoneAction() *InterpreterAction {
	return &InterpreterAction{
		actionType: ActionTypeNone,
	}
}
