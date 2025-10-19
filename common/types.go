package common

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

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
	BlobGasUsed      *uint64
	ExcessBlobGas    *uint64
}

type Transactions struct {
	Hashes [][32]byte
	Full   []Transaction // transaction needs to be defined
}
// Updated as earlier, txn data fetched from rpc was not able to unmarshal
// into the struct
type Transaction struct {
	AccessList           types.AccessList `json:"accessList"`
	Hash                 common.Hash      `json:"hash"`
	Nonce                hexutil.Uint64   `json:"nonce"`
	BlockHash            common.Hash           `json:"blockHash"`   // Pointer because it's nullable
	BlockNumber          hexutil.Uint64   `json:"blockNumber"` // Pointer because it's nullable
	TransactionIndex     hexutil.Uint64   `json:"transactionIndex"`
	From                 *common.Address           `json:"from"`
	To                   *common.Address  `json:"to"` // Pointer because 'to' can be null for contract creation
	Value                hexutil.Big      `json:"value"`
	GasPrice             hexutil.Big      `json:"gasPrice"`
	Gas                  hexutil.Uint64   `json:"gas"`
	Input                hexutil.Bytes    `json:"input"`
	ChainID              hexutil.Big      `json:"chainId"`
	TransactionType      hexutil.Uint     `json:"type"`
	Signature            *Signature       `json:"signature"`
	MaxFeePerGas         hexutil.Big      `json:"maxFeePerGas"`
	MaxPriorityFeePerGas hexutil.Big      `json:"maxPriorityFeePerGas"`
	MaxFeePerBlobGas     hexutil.Big      `json:"maxFeePerBlobGas"`
	BlobVersionedHashes  []common.Hash    `json:"blobVersionedHashes"`
}

type Signature struct {
	R       string
	S       string
	V       uint64
	YParity Parity
}

type Parity struct {
	Value bool
}

func Default() *Transactions {
	return &Transactions{
		Full: []Transaction{},
	}
}

func (t *Transactions) HashesFunc() [][32]byte {
	if len(t.Hashes) > 0 {
		return t.Hashes
	}
	hashes := make([][32]byte, len(t.Full))
	for i := range t.Full {
		hashes[i] = t.Full[i].Hash // Use the Hash field directly
	}
	return hashes
}

func (t Transactions) MarshalJSON() ([]byte, error) {
	if len(t.Hashes) > 0 {
		return json.Marshal(t.Hashes)
	}
	return json.Marshal(t.Full)
}

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
	return fmt.Sprintf("0x%x", b.Number)
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

// Example error structs can be defined here
// type BlockNotFoundError struct {}
