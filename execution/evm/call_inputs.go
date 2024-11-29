package evm
import(
	"encoding/json"
)
func (c *CallInputs) New(txEnv *TxEnv, gasLimit uint64) *CallInputs {
	// Check if the transaction kind is Call and extract the target address.
	if txEnv.TransactTo.Type != Call2 {
		return nil
	}
	targetAddress := txEnv.TransactTo.Address
	// Create and return the CallInputs instance.
	return &CallInputs{
		Input:              txEnv.Data,
		GasLimit:           gasLimit,
		TargetAddress:      *targetAddress,
		BytecodeAddress:    *targetAddress, // Set to target_address as in Rust code.
		Caller:             txEnv.Caller,
		Value:              CallValue{ValueType: "transfer", Amount: txEnv.Value},
		Scheme:             ICall, // Assuming CallScheme.Call is represented by Call.
		IsStatic:           false,
		IsEof:              false,
		ReturnMemoryOffset: Range{Start: 0, End: 0},
	}
}
func (c *CallInputs) NewBoxed(txEnv *TxEnv, gasLimit uint64) *CallInputs {
	return c.New(txEnv, gasLimit) // Returns a pointer or nil, similar to Option<Box<Self>> in Rust.
}

func (e EOFCreateInputs) NewTx(tx *TxEnv, gasLimit uint64) EOFCreateInputs {
	return EOFCreateInputs{
		Caller:   tx.Caller,
		Value:    tx.Value,
		GasLimit: gasLimit,
		Kind: EOFCreateKind{
			Kind:     TxK,
			Data: tx.Data,
		},
	}
}

func (c *CreateInputs) New(txEnv *TxEnv, gasLimit uint64) *CreateInputs {
	if txEnv.TransactTo.Type != Create2 {
		return nil
	}
	return &CreateInputs{
		Caller: txEnv.Caller,
		Scheme: CreateScheme{
			SchemeType: SchemeTypeCreate,
		},
		Value: txEnv.Value,
		InitCode: txEnv.Data,
		GasLimit: gasLimit,
	}
}

func (c *CreateInputs) NewBoxed(txEnv *TxEnv, gasLimit uint64) *CreateInputs {
	return c.New(txEnv, gasLimit)
}

type CallValue struct {
	ValueType string `json:"value_type"`
	Amount    U256   `json:"amount,omitempty"`
}

// Transfer creates a CallValue representing a transfer.
func Transfer(amount U256) CallValue {
	return CallValue{ValueType: "transfer", Amount: amount}
}

// Apparent creates a CallValue representing an apparent value.
func Apparent(amount U256) CallValue {
	return CallValue{ValueType: "apparent", Amount: amount}
}

type CallScheme int
const (
    ICall CallScheme = iota
    ICallCode
    IDelegateCall
    IStaticCall
    IExtCall	
    IExtStaticCall
    IExtDelegateCall
)
type CallInputs struct {
	Input              Bytes      `json:"input"`
	ReturnMemoryOffset Range      `json:"return_memory_offset"`
	GasLimit           uint64     `json:"gas_limit"`
	BytecodeAddress    Address    `json:"bytecode_address"`
	TargetAddress      Address    `json:"target_address"`
	Caller             Address    `json:"caller"`
	Value              CallValue  `json:"value"`
	Scheme             CallScheme `json:"scheme"`
	IsStatic           bool       `json:"is_static"`
	IsEof              bool       `json:"is_eof"`
}

// JSON serialization for CallValue
func (cv CallValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ValueType string `json:"value_type"`
		Amount    *U256  `json:"amount,omitempty"`
	}{
		ValueType: cv.ValueType,
		Amount:    &cv.Amount,
	})
}

// JSON deserialization for CallValue
func (cv *CallValue) UnmarshalJSON(data []byte) error {
	var aux struct {
		ValueType string `json:"value_type"`
		Amount    *U256  `json:"amount,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	cv.ValueType = aux.ValueType
	if aux.ValueType == "transfer" {
		cv.Amount = *aux.Amount
	} else if aux.ValueType == "apparent" {
		cv.Amount = *aux.Amount
	}
	return nil
}

type CreateInputs struct {
	Caller   Address      `json:"caller"`
	Scheme   CreateScheme `json:"scheme"`
	Value    U256         `json:"value"`
	InitCode Bytes        `json:"init_code"`
	GasLimit uint64       `json:"gas_limit"`
}
type CreateScheme struct {
	SchemeType SchemeType `json:"scheme_type"`    // can be "create" or "create2"
	Salt       *U256      `json:"salt,omitempty"` // salt is optional for Create
}

// SchemeType represents the type of creation scheme
type SchemeType int

const (
	SchemeTypeCreate SchemeType = iota
	SchemeTypeCreate2
)

type EOFCreateKindType int

const (
	TxK EOFCreateKindType = iota
	OpcodeK
)

// / Inputs for EOF create call.
type EOFCreateInputs struct {
	Caller   Address       `json:"caller"`
	Value    U256          `json:"value"`
	GasLimit uint64        `json:"gas_limit"`
	Kind     EOFCreateKind `json:"kind"`
}

type EOFCreateKind struct {
	Kind EOFCreateKindType
	Data interface{} // Use an interface to hold the associated data (Tx or Opcode)
}

// TxKind represents a transaction-based EOF create kind.
type Tx struct {
	InitData Bytes `json:"initdata"`
}
