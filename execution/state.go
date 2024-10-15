package execution

import (
	"github.com/BlocSoc-iitr/selene/common"
	"github.com/holiman/uint256"
	"sync"
)

type State struct {
	mu             sync.RWMutex
	blocks         map[uint64]*common.Block
	finalizedBlock *common.Block
	hashes         map[[32]byte]uint64
	txs            map[[32]byte]TransactionLocation
	historyLength  uint64
}
type TransactionLocation struct {
	Block uint64
	Index int
}

func NewState(historyLength uint64, blockChan <-chan *common.Block, finalizedBlockChan <-chan *common.Block) *State {
	s := &State{
		blocks:        make(map[uint64]*common.Block),
		hashes:        make(map[[32]byte]uint64),
		txs:           make(map[[32]byte]TransactionLocation),
		historyLength: historyLength,
	}
	go func() {
		for {
			select {
			case block := <-blockChan:
				if block != nil {
					s.PushBlock(block)
				}
			case block := <-finalizedBlockChan:
				if block != nil {
					s.PushFinalizedBlock(block)
				}
			}
		}
	}()

	return s
}
func (s *State) PushBlock(block *common.Block) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashes[block.Hash] = block.Number
	for i, txHash := range block.Transactions.Hashes {
		loc := TransactionLocation{
			Block: block.Number,
			Index: i,
		}
		s.txs[txHash] = loc
	}

	s.blocks[block.Number] = block

	for len(s.blocks) > int(s.historyLength) {
		var oldestNumber uint64 = ^uint64(0)
		for number := range s.blocks {
			if number < oldestNumber {
				oldestNumber = number
			}
		}
		s.removeBlock(oldestNumber)
	}
}
func (s *State) PushFinalizedBlock(block *common.Block) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.finalizedBlock = block

	if oldBlock, exists := s.blocks[block.Number]; exists {
		if oldBlock.Hash != block.Hash {
			s.removeBlock(oldBlock.Number)
			s.PushBlock(block)
		}
	} else {
		s.PushBlock(block)
	}
}
func (s *State) removeBlock(number uint64) {
	if block, exists := s.blocks[number]; exists {
		delete(s.blocks, number)
		delete(s.hashes, block.Hash)
		for _, txHash := range block.Transactions.Hashes {
			delete(s.txs, txHash)
		}
	}
}
func (s *State) GetBlock(tag common.BlockTag) *common.Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if tag.Latest {
		var latestNumber uint64
		var latestBlock *common.Block
		for number, block := range s.blocks {
			if number > latestNumber {
				latestNumber = number
				latestBlock = block
			}
		}
		return latestBlock
	} else if tag.Finalized {
		return s.finalizedBlock
	} else {
		return s.blocks[tag.Number]
	}
}
func (s *State) GetBlockByHash(hash [32]byte) *common.Block {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if number, exists := s.hashes[hash]; exists {
		return s.blocks[number]
	}
	return nil
}
func (s *State) GetTransaction(hash [32]byte) *common.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if loc, exists := s.txs[hash]; exists {
		if block, exists := s.blocks[loc.Block]; exists {
			if len(block.Transactions.Full) > loc.Index {
				return &block.Transactions.Full[loc.Index]
			}
		}
	}
	return nil
}
func (s *State) GetTransactionByBlockAndIndex(blockHash [32]byte, index uint64) *common.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if number, exists := s.hashes[blockHash]; exists {
		if block, exists := s.blocks[number]; exists {
			if int(index) < len(block.Transactions.Full) {
				return &block.Transactions.Full[index]
			}
		}
	}
	return nil
}
func (s *State) GetStateRoot(tag common.BlockTag) *[32]byte {
	if block := s.GetBlock(tag); block != nil {
		return &block.StateRoot
	}
	return nil
}
func (s *State) GetReceiptsRoot(tag common.BlockTag) *[32]byte {
	if block := s.GetBlock(tag); block != nil {
		return &block.ReceiptsRoot
	}
	return nil
}
func (s *State) GetBaseFee(tag common.BlockTag) *uint256.Int {
	if block := s.GetBlock(tag); block != nil {
		return &block.BaseFeePerGas
	}
	return nil
}
func (s *State) GetCoinbase(tag common.BlockTag) *common.Address {
	if block := s.GetBlock(tag); block != nil {
		return &block.Miner
	}
	return nil
}
func (s *State) LatestBlockNumber() *uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var latestNumber uint64
	for number := range s.blocks {
		if number > latestNumber {
			latestNumber = number
		}
	}
	if latestNumber > 0 {
		return &latestNumber
	}
	return nil
}
func (s *State) OldestBlockNumber() *uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var oldestNumber uint64 = ^uint64(0)
	for number := range s.blocks {
		if number < oldestNumber {
			oldestNumber = number
		}
	}
	if oldestNumber < ^uint64(0) {
		return &oldestNumber
	}
	return nil
}
