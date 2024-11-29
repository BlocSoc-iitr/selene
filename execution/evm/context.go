package evm
import (
	"math/big"
)
var BLOCK_HASH_HISTORY = big.NewInt(256)
type Context[EXT interface{}, DB Database] struct {
	Evm      EvmContext[DB]
	External interface{}
}
func DefaultContext[EXT interface{}]() Context[EXT, *EmptyDB] {
    return Context[EXT, *EmptyDB]{
        Evm:      NewEvmContext(new(EmptyDB)),
        External: nil,
    }	
}
func NewContext[EXT interface{}, DB Database](evm EvmContext[DB], external EXT) Context[EXT, DB] {
	return Context[EXT, DB]{
		Evm:      evm,
		External: external,
	}
}
type ContextWithHandlerCfg[EXT interface{}, DB Database] struct {
	Context Context[EXT, DB]
	Cfg     HandlerCfg
}
func (js JournaledState) SetSpecId(spec SpecId) {
	js.Spec = spec
}
//Not being used
func (c EvmContext[DB]) WithDB(db DB) EvmContext[DB] {
	return EvmContext[DB]{
		Inner:       c.Inner.WithDB(db),
		Precompiles: DefaultContextPrecompiles[DB](),
	}
}
////////////////////////////////////
///  impl Host for Context /////////
////////////////////////////////////

// func (c *Context[EXT, DB]) Env() *Env {
// 	return c.Evm.Inner.Env
// }

// func (c *Context[EXT, DB]) EnvMut() *Env {
// 	return c.Evm.Inner.Env
// }

// func (c *Context[EXT, DB]) LoadAccount(address Address) (LoadAccountResult, bool) {
// 	res, err := c.Evm.Inner.LoadAccountExist(address)
// 	if err != nil {
// 		c.Evm.Inner.Error = err
// 		return LoadAccountResult{}, false
// 	}
// 	return res, true
// }

// //Get the block hash of the given block `number`.
// func (c *Context[EXT, DB]) BlockHash(number uint64) (B256, bool) {
// 	blockNumber := c.Evm.Inner.Env.Block.Number
// 	requestedNumber := big.NewInt(int64(number))

// 	diff := new(big.Int).Sub(blockNumber, requestedNumber)

// 	if diff.Cmp(big.NewInt(0)) == 0 {
// 		return common.Hash{}, true
// 	}

// 	if diff.Cmp(BLOCK_HASH_HISTORY) <= 0 {
// 		hash, err := c.Evm.Inner.DB.BlockHash(number)
// 		if err != nil {
// 			c.Evm.Inner.Error = err
// 			return common.Hash{}, false
// 		}
// 		return hash, true
// 	}

// 	return common.Hash{}, true
// }

// // Get balance of `address` and if the account is cold.
// func (c *Context[EXT, DB]) Balance(address Address) (U256, bool, bool)

// // Get code of `address` and if the account is cold.
// func (c *Context[EXT, DB]) Code(address Address) (Bytes, bool, bool)

// // Get code hash of `address` and if the account is cold.
// func (c *Context[EXT, DB]) CodeHash(address Address) (B256, bool, bool)

// // Get storage value of `address` at `index` and if the account is cold.
// func (c *Context[EXT, DB]) SLoad(address Address, index U256) (U256, bool, bool)

// // Set storage value of account address at index.
// // Returns (original, present, new, is_cold).
// func (c *Context[EXT, DB]) SStore(address Address, index U256, value U256) (SStoreResult, bool)

// // Get the transient storage value of `address` at `index`.
// func (c *Context[EXT, DB]) TLoad(address Address, index U256) U256

// // Set the transient storage value of `address` at `index`.
// func (c *Context[EXT, DB]) TStore(address Address, index U256, value U256)

// // Emit a log owned by `address` with given `LogData`.
// func (c *Context[EXT, DB]) Log(log Log [any])

// // Mark `address` to be deleted, with funds transferred to `target`.
// func (c *Context[EXT, DB]) SelfDestruct(address Address, target Address) (SelfDestructResult, bool)

////////////////////////////////////
////////////////////////////////////
////////////////////////////////////

/*
func (c EvmContext[DB1]) WithNewEvmDB[DB2 Database](db DB2) *EvmContext[ DB2] {
    return &EvmContext[DB2]{
		Inner:	   c.Inner.WithDB(db),
        context: NewContext[EXT, DB2](
            eb.context.Evm.WithNewDB(db),
            eb.context.External.(EXT),
        ),
        handler: Handler[Context[EXT, DB2], EXT, DB2]{Cfg: eb.handler.Cfg},
        phantom: struct{}{},
    }
}*/













