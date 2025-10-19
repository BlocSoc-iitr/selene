package execution

import (
	"encoding/hex"
	"sync"
	"testing"

	"github.com/BlocSoc-iitr/selene/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)
func TestNewState(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	// Simulate blocks to push
	block1 := &common.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
		Transactions: common.Transactions{
			Hashes: [][32]byte{{0x11}, {0x12}},
		},
	}
	block2 := &common.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: common.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
		},
	}

	// Push blocks through channel
	go func() {
		blockChan <- block1
		blockChan <- block2
		close(blockChan)
	}()

	// Allow goroutine to process the blocks
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for len(state.blocks) < 2 {
			// wait for blocks to be processed
		}
	}()

	wg.Wait()

	// Verify blocks are in state
	assert.Equal(t, block1, state.blocks[1], "common.Block 1 should be present in state")
	assert.Equal(t, block2, state.blocks[2], "common.Block 2 should be present in state")

	// Check transactions are mapped correctly
	assert.Equal(t, TransactionLocation{Block: 1, Index: 0}, state.txs[[32]byte{0x11}], "Transaction 0x11 should map to common.Block 1")
	assert.Equal(t, TransactionLocation{Block: 2, Index: 1}, state.txs[[32]byte{0x22}], "Transaction 0x22 should map to common.Block 2")

	// Simulate finalized block
	finalizedBlock := &common.Block{
		Number: 2,
		Hash:   [32]byte{0x2}, // Same as block2
	}
	go func() {
		finalizedBlockChan <- finalizedBlock
		close(finalizedBlockChan)
	}()

	// Wait for finalized block to be processed
	wg.Add(1)
	go func() {
		defer wg.Done()
		for state.finalizedBlock == nil {
			_=0
			// wait for finalized block to be processed
		}
	}()
	wg.Wait()

	// Verify finalized block
	assert.Equal(t, finalizedBlock, state.finalizedBlock, "Finalized block should be correctly stored")
}

func TestRemoveBlock(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	block1 := &common.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
		Transactions: common.Transactions{
			Hashes: [][32]byte{{0x11}, {0x12}},
		},
	}
	block2 := &common.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: common.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
		},
	}

	// Add blocks
	state.PushBlock(block1)
	state.PushBlock(block2)

	// Ensure block1 is present before removal
	assert.NotNil(t, state.blocks[1], "common.Block 1 should exist before removal")
	assert.NotNil(t, state.hashes[[32]byte{0x1}], "common.Block 1's hash should exist before removal")
	assert.NotNil(t, state.txs[[32]byte{0x11}], "Transaction 0x11 should exist before removal")

	// Remove block1
	state.removeBlock(1)

	// Ensure block1 is removed
	assert.Nil(t, state.blocks[1], "common.Block 1 should not exist after removal")
	assert.Equal(t, state.hashes[[32]byte{0x1}], uint64(0), "common.Block 1's hash should not exist after removal")
	assert.Equal(t, state.txs[[32]byte{0x11}], TransactionLocation{}, "Transaction 0x11 should not exist after removal")
}

func TestGetStateBlock(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	block1 := &common.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
	}
	block2 := &common.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
	}

	state.PushBlock(block1)
	state.PushBlock(block2)

	// Test getting the latest block
	tag := common.BlockTag{Latest: true}
	latestBlock := state.GetBlock(tag)
	assert.Equal(t, block2, latestBlock, "The latest block should be block2")

	// Test getting a specific block by number
	tag = common.BlockTag{Number: 1}
	blockByNumber := state.GetBlock(tag)
	assert.Equal(t, block1, blockByNumber, "Should return block1 for number 1")

	// Test getting the finalized block (set a finalized block first)
	state.PushFinalizedBlock(block2)
	tag = common.BlockTag{Finalized: true}
	finalizedBlock := state.GetBlock(tag)
	assert.Equal(t, block2, finalizedBlock, "The finalized block should be block2")
}

func TestGetBlockByHash(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	block1 := &common.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
	}
	block2 := &common.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
	}

	state.PushBlock(block1)
	state.PushBlock(block2)

	// Test getting block by hash
	blockByHash := state.GetBlockByHash([32]byte{0x1})
	assert.Equal(t, block1, blockByHash, "Should return block1 for hash 0x1")

	blockByHash = state.GetBlockByHash([32]byte{0x2})
	assert.Equal(t, block2, blockByHash, "Should return block2 for hash 0x2")

	// Test when block is not found
	blockByHash = state.GetBlockByHash([32]byte{0x99})
	assert.Nil(t, blockByHash, "Should return nil for non-existing block hash")
}

func TestStateGetTransaction(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	tx1 := common.Transaction{Hash: [32]byte{0x11}}
	tx2 := common.Transaction{Hash: [32]byte{0x12}}

	block := &common.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
		Transactions: common.Transactions{
			Full: []common.Transaction{tx1, tx2},
			Hashes: [][32]byte{
				{0x11}, {0x12},
			},
		},
	}

	state.PushBlock(block)

	// Test for valid transaction
	foundTx := state.GetTransaction([32]byte{0x11})
	assert.NotNil(t, foundTx, "Transaction 0x11 should exist")
	assert.Equal(t, tx1, *foundTx, "Should return correct transaction")

	// Test for non-existent transaction
	nonExistentTx := state.GetTransaction([32]byte{0x99})
	assert.Nil(t, nonExistentTx, "Should return nil for non-existent transaction")
}

func TestGetTransactionByBlockAndIndex(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	tx1 := common.Transaction{Hash: [32]byte{0x11}}
	tx2 := common.Transaction{Hash: [32]byte{0x12}}

	block := &common.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
		Transactions: common.Transactions{
			Full: []common.Transaction{tx1, tx2},
		},
	}

	state.PushBlock(block)

	// Test valid transaction by block hash and index
	foundTx := state.GetTransactionByBlockAndIndex([32]byte{0x1}, 0)
	assert.NotNil(t, foundTx, "Transaction at index 0 should exist in block 0x1")
	assert.Equal(t, tx1, *foundTx, "Should return correct transaction for index 0")

	// Test for out of bounds index
	outOfBoundsTx := state.GetTransactionByBlockAndIndex([32]byte{0x1}, 2)
	assert.Nil(t, outOfBoundsTx, "Should return nil for out of bounds index")

	// Test for non-existent block hash
	nonExistentBlockTx := state.GetTransactionByBlockAndIndex([32]byte{0x99}, 0)
	assert.Nil(t, nonExistentBlockTx, "Should return nil for non-existent block hash")
}

func TestGetStateRoot(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	stateRoot := [32]byte{0xaa}
	block := &common.Block{
		Number:    1,
		Hash:      [32]byte{0x1},
		StateRoot: stateRoot,
	}

	state.PushBlock(block)

	// Test getting state root by block tag
	tag := common.BlockTag{Number: 1}
	foundStateRoot := state.GetStateRoot(tag)
	assert.NotNil(t, foundStateRoot, "State root for block 1 should exist")
	assert.Equal(t, &stateRoot, foundStateRoot, "Should return correct state root")
}

func TestGetReceiptsRoot(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	receiptsRoot := [32]byte{0xbb}
	block := &common.Block{
		Number:       1,
		Hash:         [32]byte{0x1},
		ReceiptsRoot: receiptsRoot,
	}

	state.PushBlock(block)

	// Test getting receipts root by block tag
	tag := common.BlockTag{Number: 1}
	foundReceiptsRoot := state.GetReceiptsRoot(tag)
	assert.NotNil(t, foundReceiptsRoot, "Receipts root for block 1 should exist")
	assert.Equal(t, &receiptsRoot, foundReceiptsRoot, "Should return correct receipts root")
}

func TestGetBaseFee(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	baseFee := uint256.NewInt(1000)
	block := &common.Block{
		Number:        1,
		Hash:          [32]byte{0x1},
		BaseFeePerGas: *baseFee,
	}

	state.PushBlock(block)

	// Test getting base fee by block tag
	tag := common.BlockTag{Number: 1}
	foundBaseFee := state.GetBaseFee(tag)
	assert.NotNil(t, foundBaseFee, "Base fee for block 1 should exist")
	assert.Equal(t, baseFee, foundBaseFee, "Should return correct base fee")
}

func TestGetCoinbase(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	address, err := hex.DecodeString("eccf26e9F5474882a671D6136B32BE1DF8b2CDda")
	if err != nil {
		t.Errorf("Error in converting address string to bytes")
	}
	coinbase := common.Address{Addr: [20]byte(address)}
	block := &common.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
		Miner:  coinbase,
	}

	state.PushBlock(block)

	// Test getting coinbase by block tag
	tag := common.BlockTag{Number: 1}
	foundCoinbase := state.GetCoinbase(tag)
	assert.NotNil(t, foundCoinbase, "Coinbase for block 1 should exist")
	assert.Equal(t, &coinbase, foundCoinbase, "Should return correct coinbase")
}

func TestLatestBlockNumber(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	block1 := &common.Block{Number: 1}
	block2 := &common.Block{Number: 2}

	state.PushBlock(block1)
	state.PushBlock(block2)

	latestBlockNumber := state.LatestBlockNumber()
	assert.NotNil(t, latestBlockNumber, "Latest block number should exist")
	assert.Equal(t, uint64(2), *latestBlockNumber, "Should return the latest block number 2")
}

func TestOldestBlockNumber(t *testing.T) {
	blockChan := make(chan *common.Block)
	finalizedBlockChan := make(chan *common.Block)
	state := NewState(5, blockChan, finalizedBlockChan)

	block1 := &common.Block{Number: 1}
	block2 := &common.Block{Number: 2}

	state.PushBlock(block1)
	state.PushBlock(block2)

	oldestBlockNumber := state.OldestBlockNumber()
	assert.NotNil(t, oldestBlockNumber, "Oldest block number should exist")
	assert.Equal(t, uint64(1), *oldestBlockNumber, "Should return the oldest block number 1")
}