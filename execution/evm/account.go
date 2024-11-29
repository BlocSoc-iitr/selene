package evm
import(
	"sync"
	"fmt"
	"encoding/json"
	"reflect")
type Account struct {
	Info    AccountInfo
	Storage EvmStorage
	Status  AccountStatus
}
type AccountInfo struct {
	Balance  U256
	Nonce    uint64
	CodeHash B256
	Code     *Bytecode
}
func NewAccountInfo(balance U256, nonce uint64, codeHash B256, code Bytecode) AccountInfo {
	return AccountInfo{
		Balance:  balance,
		Nonce:    nonce,
		CodeHash: codeHash,
		Code:     &code,
	}
}
type EvmStorage map[U256]EvmStorageSlot
type EvmStorageSlot struct {
	OriginalValue U256
	PresentValue  U256
	IsCold        bool
}
type AccountStatus uint8
const (
	Loaded              AccountStatus = 0b00000000 // Account is loaded but not interacted with
	Created             AccountStatus = 0b00000001 // Account is newly created, no DB fetch required
	SelfDestructed      AccountStatus = 0b00000010 // Account is marked for self-destruction
	Touched             AccountStatus = 0b00000100 // Account is touched and will be saved to DB
	LoadedAsNotExisting AccountStatus = 0b0001000  // Pre-spurious state; account loaded but marked as non-existing
	Cold                AccountStatus = 0b0010000  // Account is marked as cold (not frequently accessed)
)
type BytecodeKind int
const (
	LegacyRawKind BytecodeKind = iota
	LegacyAnalyzedKind
	EofKind
)
// BytecodeKind serves as discriminator for Bytecode to identify which variant of enum is being used
// 1 , 2 or 3 in this case
type Bytecode struct {
	Kind           BytecodeKind
	LegacyRaw      []byte                  // For LegacyRaw variant
	LegacyAnalyzed *LegacyAnalyzedBytecode // For LegacyAnalyzed variant
	Eof            *Eof                    // For Eof variant
}
type LegacyAnalyzedBytecode struct {
	Bytecode    []byte
	OriginalLen uint64
	JumpTable   JumpTable
}

// JumpTable equivalent in Go, using a sync.Once pointer to simulate Arc and BitVec<u8>.
type JumpTable struct {
	BitVector *Bitvector // Simulating BitVec<u8> as []byte
	once      sync.Once  // Lazy initialization if needed
}
type Bitvector struct {
	bits []uint8
	size int // Total number of bits represented
}
type Opcode struct {
	InitCode       Eof     `json:"initcode"`
	Input          Bytes   `json:"input"`
	CreatedAddress Address `json:"created_address"`
}
type Eof struct {
	Header EofHeader `json:"header"`
	Body   EofBody   `json:"body"`
	Raw    Bytes     `json:"raw"`
}

// EofHeader represents the header of the EOF.
type EofHeader struct {
	TypesSize         uint16   `json:"types_size"`
	CodeSizes         []uint16 `json:"code_sizes"`
	ContainerSizes    []uint16 `json:"container_sizes"`
	DataSize          uint16   `json:"data_size"`
	SumCodeSizes      int      `json:"sum_code_sizes"`
	SumContainerSizes int      `json:"sum_container_sizes"`
}

// Marshaler interface for custom marshaling
type Marshaler interface {
	MarshalType(val reflect.Value, b []byte, lastWrittenIdx uint64) (nextIdx uint64, err error)
}

// Unmarshaler interface for custom unmarshaling
type Unmarshaler interface {
	UnmarshalType(target reflect.Value, b []byte, lastReadIdx uint64) (nextIdx uint64, err error)
}

// Implement Marshaler for EOFCreateInputs
func (eci EOFCreateInputs) MarshalType(val reflect.Value, b []byte, lastWrittenIdx uint64) (nextIdx uint64, err error) {
	jsonData, err := json.Marshal(eci)
	if err != nil {
		return lastWrittenIdx, err
	}

	return uint64(len(jsonData)), nil
}

// Implement Unmarshaler for EOFCreateInputs
func (eci *EOFCreateInputs) UnmarshalType(target reflect.Value, b []byte, lastReadIdx uint64) (nextIdx uint64, err error) {
	err = json.Unmarshal(b, eci)
	if err != nil {
		return lastReadIdx, err
	}

	return lastReadIdx + uint64(len(b)), nil
}

type EofBody struct {
	TypesSection     []TypesSection `json:"types_section"`
	CodeSection      []Bytes        `json:"code_section"`      // Using [][]byte for Vec<Bytes>
	ContainerSection []Bytes        `json:"container_section"` // Using [][]byte for Vec<Bytes>
	DataSection      Bytes          `json:"data_section"`      // Using []byte for Bytes
	IsDataFilled     bool           `json:"is_data_filled"`
}

// TypesSection represents a section describing the types used in the EOF body.
type TypesSection struct {
	Inputs       uint8  `json:"inputs"`         // 1 byte
	Outputs      uint8  `json:"outputs"`        // 1 byte
	MaxStackSize uint16 `json:"max_stack_size"` // 2 bytes
}

// MarshalJSON customizes the JSON marshalling for EOFCreateKind.
func (e *EOFCreateKind) MarshalJSON() ([]byte, error) {
	switch e.Kind {
	case TxK:
		return json.Marshal(map[string]interface{}{
			"type": "Tx",
			"data": e.Data,
		})
	case OpcodeK:
		return json.Marshal(map[string]interface{}{
			"type": "Opcode",
			"data": e.Data,
		})
	default:
		return nil, fmt.Errorf("unknown EOFCreateKind type")
	}
}

// UnmarshalJSON customizes the JSON unmarshalling for EOFCreateKind.
func (e *EOFCreateKind) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Type {
	case "Tx":
		var tx Tx
		if err := json.Unmarshal(aux.Data, &tx); err != nil {
			return err
		}
		e.Kind = TxK
		e.Data = tx
	case "Opcode":
		var opcode Opcode
		if err := json.Unmarshal(aux.Data, &opcode); err != nil {
			return err
		}
		e.Kind = OpcodeK
		e.Data = opcode
	default:
		return fmt.Errorf("unknown EOFCreateKind type: %s", aux.Type)
	}

	return nil
}
