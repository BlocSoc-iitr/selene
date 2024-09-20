package execution
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)
type Account struct {
	Balance     *big.Int
	Nonce       uint64
	CodeHash    common.Hash
	Code        []byte
	StorageHash common.Hash
	Slots       map[common.Hash]*big.Int
}
type CallOpts struct {
	From     *common.Address `json:"from,omitempty"`
	To       *common.Address `json:"to,omitempty"`
	Gas      *big.Int        `json:"gas,omitempty"`
	GasPrice *big.Int        `json:"gasPrice,omitempty"`
	Value    *big.Int        `json:"value,omitempty"`
	Data     []byte          `json:"data,omitempty"`
}
func (c *CallOpts) MarshalJSON() ([]byte, error) {
	type Alias CallOpts
	return json.Marshal(&struct {
		*Alias
		Data string `json:"data,omitempty"`
	}{
		Alias: (*Alias)(c),
		Data:  hex.EncodeToString(c.Data),
	})
}
func (c *CallOpts) UnmarshalJSON(data []byte) error {
	type Alias CallOpts
	aux := &struct {
		*Alias
		Data string `json:"data,omitempty"`
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	decodedData, err := hex.DecodeString(aux.Data)
	if err != nil {
		return err
	}
	c.Data = decodedData
	return nil
}
func (c *CallOpts) String() string {
	return fmt.Sprintf("CallOpts{From: %v, To: %v, Value: %v, Data: %s}",
		c.From, c.To, c.Value, hex.EncodeToString(c.Data))
}
