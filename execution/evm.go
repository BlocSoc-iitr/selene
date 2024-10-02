package execution

import (
	"context"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"github.com/ethereum/go-ethereum/core"
	"github.com/BlocSoc-iitr/selene/common"
	Common "github.com/ethereum/go-ethereum/common" //geth common imported as Common
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/BlocSoc-iitr/selene/execution/logging"
	"github.com/ethereum/go-ethereum/params"
)

const PARALLEL_QUERY_BATCH_SIZE = 20 // Same value as used in helios
type BlockTag = common.BlockTag
type U256 = *big.Int
type B256 = Common.Hash
type Address = Common.Address


type Evm struct {
	execution *ExecutionClient
	chainID   uint64
	tag       BlockTag
}
func NewEvm(execution *ExecutionClient, chainID uint64, tag BlockTag) *Evm {
	return &Evm{
		execution: execution,
		chainID:   chainID,
		tag:       tag,
	}
}
type ExecutionResult struct {
	Success  bool
	Output   []byte
	GasUsed  uint64
	Reverted bool
}
func (e *Evm) Call(ctx context.Context, opts *CallOpts) ([]byte, error) {
	tx, err := e.callInner(ctx, opts)
	if err != nil {
		return nil, err
	}
	if tx.Success {
		return tx.Output, nil
	} else if tx.Reverted {
		return nil, &EvmError{Kind: "Revert", Details: tx.Output}
	}
	return nil, &EvmError{Kind: "Halt", Details: nil}
}
func (e *Evm) EstimateGas(ctx context.Context,opts *CallOpts) (uint64, error) {
	tx, err := e.callInner(ctx,opts)
	if err != nil {
		return 0, err
	}

	switch result := tx.result.(type) {
	case Success:
		return result.GasUsed, nil
	case Revert:
		return result.GasUsed, nil
	case Halt:
		return result.GasUsed, nil
	default:
		return 0, fmt.Errorf("unexpected execution result")
	}
}

func (e *Evm) callInner(ctx context.Context, opts *CallOpts) (*ExecutionResult, error) {
	db, err := NewProofDB(ctx, e.tag, e.execution)
	if err != nil {
		return nil, err
	}

	if err := db.State.PrefetchState(ctx, opts); err != nil {
		return nil, err
	}
	env, err := e.getEnv(ctx, opts, e.tag)
    if err != nil {
        return nil, err
    }

	evm := vm.NewEVM(env.BlockContext, env.TxContext, db.State, params.MainnetChainConfig, vm.Config{})

	chainConfig := params.MainnetChainConfig
	chainConfig.ChainID = new(big.Int).SetUint64(e.chainID)
	evm := vm.NewEVM(env.BlockContext, env.TxContext, db.State, params.MainnetChainConfig, vm.Config{})

	for {
		if db.State.NeedsUpdate() {
			if err := db.State.UpdateState(ctx); err != nil {
				return nil, err
			}
		}
		// Create a new types.Transaction
		tx := types.NewTransaction(
			0, // Nonce is not provided in CallOpts, so we use 0
			*opts.To,
			opts.Value,
			opts.Gas.Uint64(),
			opts.GasPrice,
			opts.Data,
		)

		// Create a Message from the transaction
		signer := types.NewEIP155Signer(chainConfig.ChainID)
		msg, err := core.TransactionToMessage(tx, signer, blockContext.BaseFee)
		if err != nil {
			return nil, err
		}

		// Override the From address with the one provided in CallOpts
		msg.From = *opts.From

		// Create a new EVM and apply the message
		evm := vm.NewEVM(blockContext, txContext, db.State, chainConfig, vm.Config{})
		result, err := core.ApplyMessage(evm, msg, new(core.GasPool).AddGas(opts.Gas.Uint64()))
		if err != nil {
			return nil, err
		}

		if result.Err == nil {
			return &ExecutionResult{
				Success: true,
				Output:  result.ReturnData,
				GasUsed: result.UsedGas,
			}, nil
		}

		if _, ok := result.Err.(vm.ErrExecutionReverted); ok {
			return &ExecutionResult{
				Success:  false,
				Output:   result.ReturnData,
				GasUsed:  result.UsedGas,
				Reverted: true,
			}, nil
		}

		// If the error is not a revert, we break the loop and return the error
		return nil, result.Err
	}
}
func (e *Evm) getEnv(ctx context.Context, opts *CallOpts, tag BlockTag) (vm.BlockContext, vm.TxContext, error) {
	block, err := e.execution.GetBlock(ctx, tag, false)
	if err != nil {
		return vm.BlockContext{}, vm.TxContext{}, err
	}
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash: func(n uint64) B256 {
			return B256{} // You might want to implement this properly
		},
		Coinbase:    block.Miner,
		BlockNumber: new(big.Int).SetUint64(block.Number),
		Time:        block.Timestamp,
		Difficulty:  block.Difficulty.ToBig(),
		GasLimit:    block.GasLimit,
		BaseFee:     block.BaseFeePerGas.ToBig(),
	}

	txContext := vm.TxContext{
		Origin:   *opts.From,
		GasPrice: opts.GasPrice,
	}

	return blockContext, txContext, nil
}
type ProofDB struct {
	State *EvmState
}

func NewProofDB(ctx context.Context, tag BlockTag, execution *ExecutionClient) (*ProofDB, error) {
	state := NewEvmState(execution, tag)
	return &ProofDB{
		State: state,
	}, nil
}

type StateAccess struct {
	Basic     *Address
	BlockHash *uint64
	Storage   map[Address]U256
}
type AccountInfo struct {
	Balance  U256
	Nonce    uint64
	CodeHash B256
	Code     hexutil.Bytes //Doubtful
}

func NewAccountInfo(balance U256, nonce uint64, codeHash B256, code hexutil.Bytes) AccountInfo {
	return AccountInfo{
		Balance:  balance,
		Nonce:    nonce,
		CodeHash: codeHash,
		Code:     code,
	}
}

type EvmState struct {
	Basic      map[Address]AccountInfo
	BlockHash  map[uint64]B256
	Storage    map[Address]map[U256]U256
	Block      BlockTag
	Access     *StateAccess
	Execution  *ExecutionClient
	AccessList map[Address]struct{}
} //added just now : UPDATE
func NewEvmState(execution *ExecutionClient, block BlockTag) *EvmState {
	return &EvmState{
		Basic:      make(map[Address]AccountInfo),
		BlockHash:  make(map[uint64]B256),
		Storage:    make(map[Address]map[U256]U256),
		Block:      block,
		Access:     nil,
		Execution:  execution,
		AccessList: make(map[Address]struct{}), //added just now : UPDATE
	}
}
func (e *EvmState) UpdateState(ctx context.Context) error {
	if e.Access == nil {
		return nil
	}
	access := e.Access
	e.Access = nil // Equivalent to Rust's self.access.take()
	switch {
	case access.Basic != nil:
		account, err := e.Execution.GetAccount(ctx, access.Basic, nil, e.Block)
		if err != nil {
			return err
		}
		bytecode := NewBytecodeRaw(account.Code)
		codeHash := B256FromSlice(account.CodeHash[:])
		balance := ConvertU256(account.Balance)
		accountInfo := AccountInfo{
			Balance:  balance,
			Nonce:    account.Nonce,
			CodeHash: codeHash,
			Code:     bytecode,
		}
		e.Basic[*access.Basic] = accountInfo
	case access.Storage != nil:
		for address, slotValue := range access.Storage {
			slot := Common.BigToHash(slotValue) // Use slotValue directly
			slots := []B256{slot}
			account, err := e.Execution.GetAccount(ctx, &address, slots, e.Block)
			if err != nil {
				return err
			}
			storage, ok := e.Storage[address]
			if !ok {
				storage = make(map[U256]U256) // Initialize with *big.Int
				e.Storage[address] = storage
			}
			value, ok := account.Slots[slot]
			if !ok {
				return errors.New("slot not found in account")
			}
			storage[slotValue] = value
		}
	case access.BlockHash != nil:
		block, err := e.Execution.GetBlock(ctx, BlockTag{Number: *access.BlockHash}, false)
		if err != nil {
			return err
		}
		hash := B256FromSlice(block.Hash[:]) // Convert [32]byte to []byte
		e.BlockHash[*access.BlockHash] = hash
	default:
		return errors.New("invalid access type")
	}
	return nil
}
func (e *EvmState) NeedsUpdate() bool {
	return e.Access != nil //Checks if access Field is non zero
}
func (e *EvmState) GetBasic(address Address) (AccountInfo, error) {
	if account, exists := e.Basic[address]; exists {
		return account, nil
	} else {
		e.Access = &StateAccess{Basic: &address}
		return AccountInfo{}, errors.New("state missing")
	}
}
func (e *EvmState) GetStorage(address Address, slot U256) (U256, error) {
	storage := e.Storage[address]
	if value, exists := storage[slot]; exists {
		return value, nil
	} else {
		e.Access = &StateAccess{Storage: map[Address]U256{address: slot}}
		return &big.Int{}, errors.New("state missing") // Return an empty U256 and the error
	}
}
func (e *EvmState) GetBlockHash(block uint64) (B256, error) {
	if hash, exists := e.BlockHash[block]; exists {
		return hash, nil
	} else {
		e.Access = &StateAccess{BlockHash: &block}
		return B256{}, errors.New("state missing")
	}
}
func (e *EvmState) PrefetchState(ctx context.Context, opts *CallOpts) error {
	list, err := e.Execution.Rpc.CreateAccessList(ctx, opts, e.Block)
	if err != nil {
		return err
	}
	fromAccessEntry := AccessListItem{
		Address:     *opts.From,
		StorageKeys: []B256{},
	}
	toAccessEntry := AccessListItem{
		Address:     *opts.To,
		StorageKeys: []B256{},
	}
	coinbase, err := e.Execution.GetBlock(ctx, e.Block, false)
	if err != nil {
		return err
	}
	producerAccessEntry := AccessListItem{
		Address:     coinbase.Miner,
		StorageKeys: []B256{},
	}
	listAddresses := make(map[Address]bool)
	for _, item := range list {
		listAddresses[item.Address] = true
	}

	if !listAddresses[fromAccessEntry.Address] {
		list = append(list, fromAccessEntry)
	}
	if !listAddresses[toAccessEntry.Address] {
		list = append(list, toAccessEntry)
	}
	if !listAddresses[producerAccessEntry.Address] {
		list = append(list, producerAccessEntry)
	}
	accountMap := make(map[Address]Account)
	var wg sync.WaitGroup
	sem := make(chan struct{}, PARALLEL_QUERY_BATCH_SIZE)
	var mu sync.Mutex
	for _, account := range list {
		wg.Add(1)
		sem <- struct{}{}
		go func(account AccessListItem) {
			defer wg.Done()
			defer func() { <-sem }()

			acc, err := e.Execution.GetAccount(ctx, &account.Address, account.StorageKeys, e.Block)
			if err == nil {
				mu.Lock()
				accountMap[account.Address] = *acc
				mu.Unlock()
			}
		}(account)
	}
	wg.Wait()
	for address, account := range accountMap {
		bytecode := NewBytecodeRaw(account.Code)
		codeHash := Common.BytesToHash(account.CodeHash[:])
		balance := ConvertU256(account.Balance)
		info := NewAccountInfo(balance, account.Nonce, codeHash, bytecode)
		e.Basic[address] = info
		for slot, value := range account.Slots {
			slotHash := B256FromSlice(slot[:])
			valueU256 := ConvertU256(value)
			storage, exists := e.Storage[address]
			if !exists {
				storage = make(map[U256]U256)
				e.Storage[address] = storage
			}
			slotBigInt := new(big.Int).SetBytes(slotHash.Bytes())
			storage[slotBigInt] = valueU256
		}
	}
	return nil
}

type AccessListItem struct {
	Address     Address //I used Common here instead of common
	StorageKeys []B256
}
type Bytecode []byte

func NewBytecodeRaw(code []byte) hexutil.Bytes {
	return hexutil.Bytes(code)
}
func B256FromSlice(slice []byte) Common.Hash {
	return Common.BytesToHash(slice)
}
func ConvertU256(value *big.Int) *big.Int {
	valueSlice := make([]byte, 32)
	value.FillBytes(valueSlice)
	result := new(big.Int).SetBytes(valueSlice)
	return result
}

type Database interface {
	Basic(address Address) (AccountInfo, error)
	BlockHash(number uint64) (B256, error)
	Storage(address Address, slot *big.Int) (*big.Int, error)
	CodeByHash(codeHash B256) (Bytecode, error)
}

func (db *ProofDB) Basic(address Address) (AccountInfo, error) {
	if isPrecompile(address) {
		return AccountInfo{}, nil // Return a default AccountInfo
	}
	logging.Trace("fetch basic evm state for address", zap.String("address", address.Hex()))
	return db.State.GetBasic(address)
}
func (db *ProofDB) BlockHash(number uint64) (B256, error) {
	logging.Trace("fetch block hash for block number", zap.Uint64("number", number))
	return db.State.GetBlockHash(number)
}
func (db *ProofDB) Storage(address Address, slot *big.Int) (*big.Int, error) {
	logging.Trace("fetch storage for address and slot",
		zap.String("address", address.Hex()),
		zap.String("slot", slot.String()))
	return db.State.GetStorage(address, slot)
}
func (db *ProofDB) CodeByHash(codeHash B256) (Bytecode, error) {
	logging.Trace("fetch code by hash", zap.String("codeHash", codeHash.Hex()))
	return nil, nil
}
func isPrecompile(address Address) bool {
	precompileAddress := Address{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09}
	return address.Cmp(precompileAddress) <= 0 && address.Cmp(Address{}) > 0
}


//Skipped the testing for now
//Proposal: We should be using geth Address instead of locally defined address in common/types.go in the entire codebase
