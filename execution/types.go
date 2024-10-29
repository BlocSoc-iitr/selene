package execution

import (
	"encoding/json"
	"fmt"
	"github.com/BlocSoc-iitr/selene/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
	"math/big"
	"reflect"
)

type FeeHistory struct {
	BaseFeePerGas []hexutil.Big
	GasUsedRatio  []float64
	OldestBlock   *hexutil.Big
	Reward        [][]hexutil.Big
}

// This is to help in unmarshaling values from rpc response
type StorageProof struct {
	Key   common.Hash     `json:"key"`
	Proof []hexutil.Bytes `json:"proof"`
	Value *uint256.Int    `json:"value"`
}
type EIP1186ProofResponse struct {
	Address      common.Address  `json:"address"`
	Balance      *uint256.Int    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	AccountProof []hexutil.Bytes `json:"accountProof"`
	StorageProof []StorageProof  `json:"storageProof"`
}
type Account struct {
	Balance     *big.Int
	Nonce       uint64
	CodeHash    common.Hash
	Code        []byte
	StorageHash common.Hash
	Slots       []Slot
}
// This is to help in unmarshaling values from rpc response
type Slot struct {
	Key   common.Hash // The key (slot)
	Value *big.Int    // The value (storage value)
}
type CallOpts struct {
	From     *common.Address `json:"from,omitempty"`
	To       *common.Address `json:"to,omitempty"`
	Gas      *big.Int        `json:"gas,omitempty"`
	GasPrice *big.Int        `json:"gasPrice,omitempty"`
	Value    *big.Int        `json:"value,omitempty"`
	Data     []byte          `json:"data,omitempty"`
}

func (c *CallOpts) String() string {
	return fmt.Sprintf("CallOpts{From: %v, To: %v, Gas: %v, GasPrice: %v, Value: %v, Data: 0x%x}",
		c.From, c.To, c.Gas, c.GasPrice, c.Value, c.Data)
}

func (c *CallOpts) Serialize() ([]byte, error) {
	serialized := make(map[string]interface{})
	v := reflect.ValueOf(*c)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		if !field.IsNil() {
			var value interface{}
			var err error

			switch field.Interface().(type) {
			case *common.Address:
				value = utils.Address_to_hex_string(*field.Interface().(*common.Address))
			case *big.Int:
				value = utils.U64_to_hex_string(field.Interface().(*big.Int).Uint64())
			case []byte:
				value, err = utils.Bytes_serialize(field.Interface().([]byte))
				if err != nil {
					return nil, fmt.Errorf("error serializing %s: %w", fieldName, err)
				}
			default:
				return nil, fmt.Errorf("unsupported type for field %s", fieldName)
			}

			serialized[fieldName] = value
		}
	}

	return json.Marshal(serialized)
}

func (c *CallOpts) Deserialize(data []byte) error {
	var serialized map[string]string
	if err := json.Unmarshal(data, &serialized); err != nil {
		return err
	}

	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		if value, ok := serialized[fieldName]; ok {
			switch field.Interface().(type) {
			case *common.Address:
				addressBytes, err := utils.Hex_str_to_bytes(value)
				if err != nil {
					return fmt.Errorf("error deserializing %s: %w", fieldName, err)
				}
				addr := common.BytesToAddress(addressBytes)
				field.Set(reflect.ValueOf(&addr))
			case *big.Int:
				intBytes, err := utils.Hex_str_to_bytes(value)
				if err != nil {
					return fmt.Errorf("error deserializing %s: %w", fieldName, err)
				}
				bigInt := new(big.Int).SetBytes(intBytes)
				field.Set(reflect.ValueOf(bigInt))
			case []byte:
				byteValue, err := utils.Bytes_deserialize([]byte(value))
				if err != nil {
					return fmt.Errorf("error deserializing %s: %w", fieldName, err)
				}
				field.SetBytes(byteValue)
			default:
				return fmt.Errorf("unsupported type for field %s", fieldName)
			}
		}
	}

	return nil
}