package evm
import(
	"fmt"
	"encoding/json"

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
type jsonInterpreterAction struct {
	Type         string             `json:"type"`
	CallInputs   *CallInputs        `json:"call_inputs,omitempty"`
	CreateInputs *CreateInputs      `json:"create_inputs,omitempty"`
	EOFInputs    *EOFCreateInputs   `json:"eof_inputs,omitempty"`
	Result       *InterpreterResult `json:"result,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface
func (a *InterpreterAction) MarshalJSON() ([]byte, error) {
	var jsonAction jsonInterpreterAction

	switch a.actionType {
	case ActionTypeCall:
		jsonAction.Type = "Call"
		jsonAction.CallInputs = a.callInputs
	case ActionTypeCreate:
		jsonAction.Type = "Create"
		jsonAction.CreateInputs = a.createInputs
	case ActionTypeEOFCreate:
		jsonAction.Type = "EOFCreate"
		jsonAction.EOFInputs = a.eofCreateInputs
	case ActionTypeReturn:
		jsonAction.Type = "Return"
		jsonAction.Result = a.result
	case ActionTypeNone:
		jsonAction.Type = "None"
	default:
		return nil, fmt.Errorf("unknown action type: %v", a.actionType)
	}

	return json.Marshal(jsonAction)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *InterpreterAction) UnmarshalJSON(data []byte) error {
	var jsonAction jsonInterpreterAction
	if err := json.Unmarshal(data, &jsonAction); err != nil {
		return err
	}

	switch jsonAction.Type {
	case "Call":
		if jsonAction.CallInputs == nil {
			return fmt.Errorf("Call action missing inputs")
		}
		*a = *NewCallAction(jsonAction.CallInputs)
	case "Create":
		if jsonAction.CreateInputs == nil {
			return fmt.Errorf("Create action missing inputs")
		}
		*a = *NewCreateAction(jsonAction.CreateInputs)
	case "EOFCreate":
		if jsonAction.EOFInputs == nil {
			return fmt.Errorf("EOFCreate action missing inputs")
		}
		*a = *NewEOFCreateAction(jsonAction.EOFInputs)
	case "Return":
		if jsonAction.Result == nil {
			return fmt.Errorf("Return action missing result")
		}
		*a = *NewReturnAction(jsonAction.Result)
	case "None":
		*a = *NewNoneAction()
	default:
		return fmt.Errorf("unknown action type: %s", jsonAction.Type)
	}

	return nil
}

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