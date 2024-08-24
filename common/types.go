package common

//other option is to import a package of go-ethereum but that was weird
import (
	"github.com/holiman/uint256"
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
}

// can be an enum
type BlockTag struct {
}

// need some error structs and enums as well
// Example BlockNotFoundError
