package common

//other option is to import a package of go-ethereum but that was weird
import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"strconv"
)

// need to confirm how such primitive types will be imported,
//
// Transaction  https://docs.rs/alloy/latest/alloy/rpc/types/struct.Transaction.html
// address:  20 bytes
// B256: 32 bytes  https://bluealloy.github.io/revm/docs/revm/precompile/primitives/type.B256.html
type Address struct {
	Addr [20]byte
}

type Block struct {
	Number           uint64
	BaseFeePerGas    uint256.Int
	Difficulty       uint256.Int
	ExtraData        []byte
	GasLimit         uint64
	GasUsed          uint64
	Hash             [32]byte
	LogsBloom        []byte
	Miner            Address
	MixHash          [32]byte
	Nonce            string
	ParentHash       [32]byte
	ReceiptsRoot     [32]byte
	Sha3Uncles       [32]byte
	Size             uint64
	StateRoot        [32]byte
	Timestamp        uint64
	TotalDifficulty  uint64
	Transactions     Transactions
	TransactionsRoot [32]byte
	Uncles           [][32]byte
}

// an enum having 2 types- how to implement??
type Transactions struct {
	Hashes [][32]byte
	Full   []types.Transaction //transaction needs to be defined
}

func Default() *Transactions {
	return &Transactions{
		Full: []types.Transaction{},
	}
}
func (t *Transactions) HashesFunc() [][32]byte {
	if len(t.Hashes) > 0 { //if Transactions struct contains hashes then return them directly
		return t.Hashes
	}
	hashes := make([][32]byte, len(t.Full))
	for i := range t.Full {
		hashes[i] = t.Full[i].Hash()
	}
	return hashes
}
func (t Transactions) MarshalJSON() ([]byte, error) {
	if len(t.Hashes) > 0 {
		return json.Marshal(t.Hashes)
	}
	return json.Marshal(t.Full)
}

// can be an enum
type BlockTag struct {
	Latest    bool
	Finalized bool
	Number    uint64
}

func (b BlockTag) String() string {
	if b.Latest {
		return "latest"
	}
	if b.Finalized {
		return "finalized"
	}
	return fmt.Sprintf("%d", b.Number)
}

func (b *BlockTag) UnmarshalJSON(data []byte) error {
	var block string
	if err := json.Unmarshal(data, &block); err != nil {
		return err
	}
	switch block {
	case "latest":
		b.Latest = true
	case "finalized":
		b.Finalized = true
	default:
		var err error
		b.Number, err = parseBlockNumber(block)
		if err != nil {
			return err
		}
	}
	return nil
}
func parseBlockNumber(block string) (uint64, error) {
	if len(block) > 2 && block[:2] == "0x" {
		return parseHexUint64(block[2:])
	}
	return parseDecimalUint64(block)
}
func parseHexUint64(hexStr string) (uint64, error) {
	return strconv.ParseUint(hexStr, 16, 64)
}
func parseDecimalUint64(decStr string) (uint64, error) {
	return strconv.ParseUint(decStr, 10, 64)
}

// need some error structs and enums as well
// Example BlockNotFoundError
