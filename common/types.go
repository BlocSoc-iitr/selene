package common

//other option is to import a package of go-ethereum but that was weird
import (
	"github.com/holiman/uint256"
   "encoding/json"
   "strconv"
   "fmt"
   "github.com/ethereum/go-ethereum/core/types"
)

// need to confirm how such primitive types will be imported,
//
//Transaction  https://docs.rs/alloy/latest/alloy/rpc/types/struct.Transaction.html
// address:  20 bytes
// B256: 32 bytes  https://bluealloy.github.io/revm/docs/revm/precompile/primitives/type.B256.html
type Address struct{
   addr [20]byte
}

type Block struct {
   number uint64
   base_fee_per_gas uint256.Int
   difficulty uint256.Int
   extra_data []byte
   gas_limit uint64
   gas_used uint64
   hash [32]byte
   logs_bloom []byte
   miner Address
   mix_hash [32]byte
   nonce string
   parent_hash [32]byte
   receipts_root [32]byte
   sha3_uncles [32]byte
   size uint64
   state_root [32]byte
   timestamp uint64
   total_difficulty uint64
   transactions Transactions
   transactions_root [32]byte
   uncles [][32]byte
}
// an enum having 2 types- how to implement??
type Transactions struct {
   Hashes [][32]byte
   Full []types.Transaction //transaction needs to be defined 
}
func Default() *Transactions{
   return &Transactions{
      Full: []types.Transaction{},
   }
}
func (t *Transactions) hashes() [][32]byte{
   if len(t.Hashes)>0{//if Transactions struct contains hashes then return them directly
      return t.Hashes
   }
   hashes := make([][32]byte, len(t.Full))
	for i, tx := range t.Full {
		hashes[i] = tx.Hash()
	}
	return hashes
}
func (t Transactions) MarshalJSON()([]byte,error){
   if len(t.Hashes)>0{
      return json.Marshal(t.Hashes)
   }
   return json.Marshal(t.Full)
}

// can be an enum
type BlockTag struct {
   Latest bool
   Finalized bool
   Number uint64
}

func (b BlockTag) String() string {
   if b.Latest{
      return "latest"
   }
   if b.Finalized{
      return "finalized"
   }
   return fmt.Sprintf("%d", b.Number)
}

func (b *BlockTag)UnmarshalJSON(data []byte) error{
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
