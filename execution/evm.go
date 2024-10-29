package execution
import (
	"math/big"
	"fmt"
	"bytes"
	"sync"
	"encoding/hex"
	Common "github.com/ethereum/go-ethereum/common" //geth common imported as Common
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	Gevm "github.com/BlocSoc-iitr/selene/execution/evm"
	"github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/execution/logging"
	"go.uber.org/zap"
)
type BlockTag = common.BlockTag
type U256 = *big.Int
type B256 = Common.Hash
type Address = common.Address
type ProofDB struct {
    State *EvmState
}
func NewProofDB( tag BlockTag, execution *ExecutionClient) (*ProofDB, error) {
    state := NewEvmState(execution, tag)
    return &ProofDB{
        State: state,
    }, nil
}
type StateAccess struct {
    Basic     *Address
    BlockHash *uint64
    Storage   *struct {
        Address Address
        Slot    U256
    }
}
type EvmState struct {
    Basic      map[Address]Gevm.AccountInfo
    BlockHash  map[uint64]B256
    Storage    map[Address]map[U256]U256
    Block      BlockTag
    Access     *StateAccess
    Execution  *ExecutionClient
    mu         sync.RWMutex
}
func NewEvmState(execution *ExecutionClient, block BlockTag) *EvmState {
	return &EvmState{
		Basic:      make(map[Address]Gevm.AccountInfo),
		BlockHash:  make(map[uint64]B256),
		Storage:    make(map[Address]map[U256]U256),
		Block:      block,
		Access:     nil,
		Execution:  execution,
	}
}
func (e *EvmState) UpdateState() error {
    e.mu.Lock()
    if e.Access == nil {
        e.mu.Unlock()
        return nil
    }
    access := e.Access
    e.Access = nil
    e.mu.Unlock()

    switch {
    case access.Basic != nil:
        account, err := e.Execution.GetAccount(access.Basic, []Common.Hash{}, e.Block)
        if err != nil {
            return err
        }
        e.mu.Lock()
        bytecode := NewRawBytecode(account.Code)
        codeHash := B256FromSlice(account.CodeHash[:])
        balance := ConvertU256(account.Balance)
        accountInfo := Gevm.AccountInfo{
            Balance:  balance,
            Nonce:    account.Nonce,
            CodeHash: codeHash,
            Code:     &bytecode,
        }
        e.Basic[*access.Basic] = accountInfo
        e.mu.Unlock()
    case access.Storage != nil:
        slotHash := Common.BytesToHash(access.Storage.Slot.Bytes())
        slots := []Common.Hash{slotHash}
        account, err := e.Execution.GetAccount(&access.Storage.Address, slots, e.Block)
        if err != nil {
            return err
        }
        e.mu.Lock()
        storage, ok := e.Storage[access.Storage.Address]
        if !ok {
            storage = make(map[U256]U256)
            e.Storage[access.Storage.Address] = storage
        }
        
        slotValue, exists := account.Slots[slotHash]
        if !exists {
            e.mu.Unlock()
            return fmt.Errorf("storage slot %v not found in account", slotHash)
        }
        
        value := U256FromBigEndian(slotValue.Bytes())
        storage[access.Storage.Slot] = value
        e.mu.Unlock()

    case access.BlockHash != nil:
        block, err := e.Execution.GetBlock(BlockTag{Number: *access.BlockHash}, false)
        if err != nil {
            return err
        }
        
        e.mu.Lock()
        hash := B256FromSlice(block.Hash[:])
        e.BlockHash[*access.BlockHash] = hash
        e.mu.Unlock()

    default:
        return errors.New("invalid access type")
    }
    return nil
}

func (e *EvmState) NeedsUpdate() bool {
    e.mu.RLock()
    defer e.mu.RUnlock()
    return e.Access != nil
}
func (e *EvmState) GetBasic(address Address) (Gevm.AccountInfo, error) {
    e.mu.RLock()
    account, exists := e.Basic[address]
    e.mu.RUnlock()

    if exists {
        return account, nil
    }

    e.mu.Lock()
    e.Access = &StateAccess{Basic: &address}
    e.mu.Unlock()
    
    return Gevm.AccountInfo{}, fmt.Errorf("state missing")
}

func (e *EvmState) GetStorage(address Address, slot U256) (U256, error) {
    // Lock for reading
    e.mu.RLock()
    storage, exists := e.Storage[address]
    if exists {
        if value, exists := storage[slot]; exists {
            e.mu.RUnlock()
            return value, nil
        }
    }
    e.mu.RUnlock()

    // If we need to update state, we need a write lock
    e.mu.Lock()
    // Set the access pattern for state update
    e.Access = &StateAccess{
        Storage: &struct {
            Address Address
            Slot    U256
        }{
            Address: address,
            Slot:    slot,
        },
    }
    e.mu.Unlock()

    // Return zero value and error to indicate state needs to be updated
    return &big.Int{}, errors.New("state missing")
}
func (e *EvmState) GetBlockHash(block uint64) (B256, error) {
    e.mu.RLock()
    hash, exists := e.BlockHash[block]
    e.mu.RUnlock()

    if exists {
        return hash, nil
    }

    e.mu.Lock()
    e.Access = &StateAccess{BlockHash: &block}
    e.mu.Unlock()

    return B256{}, fmt.Errorf("state missing")
}

// PrefetchState prefetches state data
func (e *EvmState) PrefetchState(opts *CallOpts) error {
    list, err := e.Execution.Rpc.CreateAccessList(*opts, e.Block)
    if err != nil {
        return fmt.Errorf("create access list: %w", err)
    }

    // Create access entries
    fromAccessEntry := AccessListItem{
        Address:     *opts.From,
        StorageKeys: make([]Common.Hash, 0),
    }
    toAccessEntry := AccessListItem{
        Address:     *opts.To,
        StorageKeys: make([]Common.Hash, 0),
    }

    block, err := e.Execution.GetBlock(e.Block, false)
    if err != nil {
        return fmt.Errorf("get block: %w", err)
    }

    producerAccessEntry := AccessListItem{
        Address:     block.Miner,
        StorageKeys: make([]Common.Hash, 0),
    }

    // Use a map for O(1) lookup of addresses
    listAddresses := make(map[Address]struct{})
    for _, item := range list {
        listAddresses[item.Address] = struct{}{}
    }

    // Add missing entries
    if _, exists := listAddresses[fromAccessEntry.Address]; !exists {
        list = append(list, fromAccessEntry)
    }
    if _, exists := listAddresses[toAccessEntry.Address]; !exists {
        list = append(list, toAccessEntry)
    }
    if _, exists := listAddresses[producerAccessEntry.Address]; !exists {
        list = append(list, producerAccessEntry)
    }

    // Process accounts in parallel with bounded concurrency
    type accountResult struct {
        address Address
        account Account
        err     error
    }

    batchSize := PARALLEL_QUERY_BATCH_SIZE
    resultChan := make(chan accountResult, len(list))
    semaphore := make(chan struct{}, batchSize)

    var wg sync.WaitGroup
    for _, item := range list {
        wg.Add(1)
        go func(item AccessListItem) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release

            account, err := e.Execution.GetAccount(&item.Address, item.StorageKeys, e.Block)
            resultChan <- accountResult{
                address: item.Address,
                account: account,
                err:     err,
            }
        }(item)
    }

    // Close result channel when all goroutines complete
    go func() {
        wg.Wait()
        close(resultChan)
    }()

    // Process results and update state
    e.mu.Lock()
    defer e.mu.Unlock()

    for result := range resultChan {
        if result.err != nil {
            continue
        }

        account := result.account
        address := result.address

        // Update basic account info
        bytecode := NewRawBytecode(account.Code)
        codeHash := B256FromSlice(account.CodeHash[:])
        balance := ConvertU256(account.Balance)
        info := Gevm.NewAccountInfo(balance, account.Nonce, codeHash, bytecode)
        e.Basic[address] = info

        // Update storage
        storage := e.Storage[address]
        if storage == nil {
            storage = make(map[U256]U256)
            e.Storage[address] = storage
        }

        for slot, value := range account.Slots {
            slotHash := B256FromSlice(slot[:])
            valueU256 := ConvertU256(value)
            storage[U256FromBytes(slotHash.Bytes())] = valueU256
        }
    }

    return nil
}
func U256FromBytes(b []byte) U256 {
    return new(big.Int).SetBytes(b)
}
func U256FromBigEndian(b []byte) *big.Int {
    if len(b) != 32 {
        return nil // or handle the error appropriately
    }
    return new(big.Int).SetBytes(b)
}
func (db *ProofDB) Basic(address Address) (Gevm.AccountInfo, error) {
	if isPrecompile(address) {
		return Gevm.AccountInfo{}, nil // Return a default AccountInfo
	}
	//logging.Trace("fetch basic evm state for address", zap.String("address", address.Hex()))
	logging.Trace("fetch basic evm state for address", zap.String("address",hex.EncodeToString(address.Addr[:]) ))
	return db.State.GetBasic(address)
}

func (db *ProofDB) BlockHash(number uint64) (B256, error) {
	logging.Trace("fetch block hash for block number", zap.Uint64("number", number))
	return db.State.GetBlockHash(number)
}
func (db *ProofDB) Storage(address Address, slot *big.Int) (*big.Int, error) {
	logging.Trace("fetch storage for address and slot",
		zap.String("address",hex.EncodeToString(address.Addr[:]) ),
		zap.String("slot", slot.String()))
	return db.State.GetStorage(address, slot)
}
func (db *ProofDB) CodeByHash(_ B256) (Gevm.Bytecode, error) {
    return Gevm  .Bytecode{}, errors.New("should never be called")
}
func isPrecompile(address Address) bool {
    precompileAddress := Common.BytesToAddress([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09})
	zeroAddress := common.Address{}
	return bytes.Compare(address.Addr[:], precompileAddress[:]) <= 0 && bytes.Compare(address.Addr[:], zeroAddress.Addr[:]) > 0
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
func NewRawBytecode(raw []byte) Gevm.Bytecode {
	return Gevm.Bytecode{
		Kind:      Gevm.LegacyRawKind,
		LegacyRaw: raw,
	}
}
