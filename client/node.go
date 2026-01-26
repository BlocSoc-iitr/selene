package client

import (
	"errors"
	"github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/config"
	"github.com/BlocSoc-iitr/selene/consensus"
	"github.com/BlocSoc-iitr/selene/execution"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

// Used instead of Alloy::TransactionRequest
type TransactionRequest struct {
	From     string
	To       string
	Value    *big.Int
	Gas      uint64
	GasPrice *big.Int
	Data     []byte
}

// Here [32]byte is used for B256 type in Rust and *big.Int is used for U256 type
type Node struct {
	Consensus   *consensus.ConsensusClient
	Execution   *execution.ExecutionClient
	Config      *config.Config
	HistorySize int
}

func NewNode(config *config.Config) (*Node, error) {
	consensusRPC := config.ConsensusRpc
	executionRPC := config.ExecutionRpc
	// Initialize ConsensusClient
	consensus := consensus.ConsensusClient{}.New(&consensusRPC, *config)
	// Extract block receivers
	blockRecv := consensus.BlockRecv
	if blockRecv == nil {
		return nil, errors.New("blockRecv is nil")
	}
	finalizedBlockRecv := consensus.FinalizedBlockRecv
	if finalizedBlockRecv == nil {
		return nil, errors.New("finalizedBlockRecv is nil")
	}
	// Initialize State
	state := execution.State{}.New(blockRecv, finalizedBlockRecv, 256)
	// Initialize ExecutionClient
	execution, err := execution.ExecutionClient{}.New(executionRPC, state)
	if err != nil {
		return nil, errors.New("ExecutionClient creation error")
	}
	// Return the constructed Node
	return &Node{
		Consensus:   &consensus,
		Execution:   execution,
		Config:      config,
		HistorySize: 64,
	}, nil
}

func (n *Node) Call(tx *TransactionRequest, block common.BlockTag) ([]byte, *NodeError) {
    resultChan := make(chan []byte)
    errorChan := make(chan *NodeError)

    go func() {
        n.CheckBlocktagAge(block)
        evm := (&execution.Evm{}).New(n.Execution, n.ChainId, block)
        result, err := evm.Call(tx)
        if err != nil {
            errorChan <- NewExecutionClientCreationError(err)
            return
        }
        resultChan <- result
        errorChan <- nil
    }()

    select {
    case result := <-resultChan:
        return result, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) EstimateGas(tx *TransactionRequest) (uint64, *NodeError) {
    resultChan := make(chan uint64)
    errorChan := make(chan *NodeError)

    go func() {
        n.CheckHeadAge()
        evm := (&execution.Evm{}).New(n.Execution, n.ChainId, common.BlockTag{Latest: true})
        result, err := evm.EstimateGas(tx)
        if err != nil {
            errorChan <- NewExecutionEvmError(err)
            return
        }
        resultChan <- result
        errorChan <- nil
    }()

    select {
    case result := <-resultChan:
        return result, nil
    case err := <-errorChan:
        return 0, err
    }
}

func (n *Node) GetBalance(address string, tag common.BlockTag) (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        n.CheckBlocktagAge(tag)
        account, err := n.Execution.GetAccount(address, nil, tag)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- account.Balance
        errorChan <- nil
    }()

    select {
    case balance := <-resultChan:
        return balance, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetNonce(address string, tag common.BlockTag) (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)

    go func() {
        n.CheckBlocktagAge(tag)
        account, err := n.Execution.GetAccount(address, nil, tag)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- account.Nonce
        errorChan <- nil
    }()

    select {
    case nonce := <-resultChan:
        return nonce, nil
    case err := <-errorChan:
        return 0, err
    }
}

func (n *Node) GetBlockTransactionCountByHash(hash [32]byte) (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)

    go func() {
        block, err := n.Execution.GetBlockByHash(hash, false)
        if err != nil {
            errorChan <- err
            return
        }
        transactionCount := len(block.Transactions.HashesFunc())
        resultChan <- uint64(transactionCount)
        errorChan <- nil
    }()

    select {
    case txCount := <-resultChan:
        return txCount, nil
    case err := <-errorChan:
        return 0, err
    }
}

func (n *Node) GetBlockTransactionCountByNumber(tag common.BlockTag) (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)

    go func() {
        block, err := n.Execution.GetBlock(tag, false)
        if err != nil {
            errorChan <- err
            return
        }
        transactionCount := len(block.Transactions.HashesFunc())
        resultChan <- uint64(transactionCount)
        errorChan <- nil
    }()

    select {
    case txCount := <-resultChan:
        return txCount, nil
    case err := <-errorChan:
        return 0, err
    }
}

func (n *Node) GetCode(address string, tag common.BlockTag) ([]byte, error) {
    resultChan := make(chan []byte)
    errorChan := make(chan error)

    go func() {
        n.CheckBlocktagAge(tag)
        account, err := n.Execution.GetAccount(address, nil, tag)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- account.Code
        errorChan <- nil
    }()

    select {
    case code := <-resultChan:
        return code, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetStorageAt(address string, slot [32]byte, tag common.BlockTag) (*big.Int, error) {
	resultChan := make(chan *big.Int)
	errorChan := make(chan error)

	go func() {
		n.CheckHeadAge()
		account := n.Execution.GetAccount(address, &slot, tag)
		value := account.Slots.Get(&slot)
		if value == nil {
			resultChan <- nil
			errorChan <- errors.New("slot not found")
			return
		}
		resultChan <- value
		errorChan <- nil
	}()

	return <-resultChan, <-errorChan
}

func (n *Node) SendRawTransaction(bytes *[]uint8) ([32]byte, error) {
    resultChan := make(chan [32]byte)
    errorChan := make(chan error)

    go func() {
        txHash, err := n.Execution.SendRawTransaction(bytes)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txHash
        errorChan <- nil
    }()

    select {
    case txHash := <-resultChan:
        return txHash, nil
    case err := <-errorChan:
        return [32]byte{}, err
    }
}

func (n *Node) GetTransactionReceipt(txHash [32]byte) (*types.Receipt, error) {
    resultChan := make(chan *types.Receipt)
    errorChan := make(chan error)

    go func() {
        txnReceipt, err := n.Execution.GetTransactionReceipt(txHash)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txnReceipt
        errorChan <- nil
    }()

    select {
    case receipt := <-resultChan:
        return receipt, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetTransactionByHash(txHash [32]byte) (*types.Transaction, error) {
    resultChan := make(chan *types.Transaction)
    errorChan := make(chan error)

    go func() {
        txn, err := n.Execution.GetTransactionByHash(txHash)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txn
        errorChan <- nil
    }()

    select {
    case txn := <-resultChan:
        return txn, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetTransactionByBlockHashAndIndex(hash [32]byte, index uint64) (*types.Transaction, error) {
    resultChan := make(chan *types.Transaction)
    errorChan := make(chan error)

    go func() {
        txn, err := n.Execution.GetTransactionByBlockHashAndIndex(hash, index)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txn
        errorChan <- nil
    }()

    select {
    case txn := <-resultChan:
        return txn, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetLogs(filter *ethereum.FilterQuery) ([]types.Log, error) {
    resultChan := make(chan []types.Log)
    errorChan := make(chan error)

    go func() {
        logs, err := n.Execution.GetLogs(filter)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- logs
        errorChan <- nil
    }()

    select {
    case logs := <-resultChan:
        return logs, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetFilterChanges(filterId *big.Int) ([]types.Log, error) {
    resultChan := make(chan []types.Log)
    errorChan := make(chan error)

    go func() {
        logs, err := n.Execution.GetFilterChanges(filterId)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- logs
        errorChan <- nil
    }()

    select {
    case logs := <-resultChan:
        return logs, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) UninstallFilter(filterId *big.Int) (bool, error) {
    resultChan := make(chan bool)
    errorChan := make(chan error)

    go func() {
        success, err := n.Execution.UninstallFilter(filterId)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- success
        errorChan <- nil
    }()

    select {
    case success := <-resultChan:
        return success, nil
    case err := <-errorChan:
        return false, err
    }
}

func (n *Node) GetNewFilter(filter *ethereum.FilterQuery) (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        filterId, err := n.Execution.GetNewFilter(filter)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- filterId
        errorChan <- nil
    }()

    select {
    case filterId := <-resultChan:
        return filterId, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetNewBlockFilter() (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        filterId, err := n.Execution.GetNewBlockFilter()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- filterId
        errorChan <- nil
    }()

    select {
    case filterId := <-resultChan:
        return filterId, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetNewPendingTransactionFilter() (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        filterId, err := n.Execution.GetNewPendingTransactionFilter()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- filterId
        errorChan <- nil
    }()

    select {
    case filterId := <-resultChan:
        return filterId, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetGasPrice() (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        n.CheckHeadAge()
        block, err := n.Execution.GetBlock(common.BlockTag{}.Latest, false)
        if err != nil {
            errorChan <- err
            return
        }
        baseFee := block.BaseFeePerGas
        tip := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)  // 1 Gwei
        resultChan <- new(big.Int).Add(baseFee.ToBig(), tip)
        errorChan <- nil
    }()

    select {
    case gasPrice := <-resultChan:
        return gasPrice, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetPriorityFee() *big.Int {
	tip := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	return tip
}

func (n *Node) GetBlockNumber() (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        n.CheckHeadAge()
        block, err := n.Execution.GetBlock(common.BlockTag{}.Latest, false)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- big.NewInt(int64(block.Number))
        errorChan <- nil
    }()

    select {
    case blockNumber := <-resultChan:
        return blockNumber, nil
    case err := <-errorChan:
        return nil, err
    }
}

func (n *Node) GetBlockByNumber(tag common.BlockTag, fullTx bool) (common.Block, error) {
    resultChan := make(chan common.Block)
    errorChan := make(chan error)

    go func() {
        n.CheckBlocktagAge(tag)
        block, err := n.Execution.GetBlock(tag, fullTx)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- block
        errorChan <- nil
    }()

    select {
    case block := <-resultChan:
        return block, nil
    case err := <-errorChan:
        return common.Block{}, err
    }
}

func (n *Node) GetBlockByHash(hash [32]byte, fullTx bool) (common.Block, error) {
    resultChan := make(chan common.Block)
    errorChan := make(chan error)

    go func() {
        block, err := n.Execution.GetBlockByHash(hash, fullTx)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- block
        errorChan <- nil
    }()

    select {
    case block := <-resultChan:
        return block, nil
    case err := <-errorChan:
        return common.Block{}, err
    }
}

func (n *Node) ChainId() uint64 {
	return n.Config.Chain.ChainID
}

type SyncInfo struct {
	CurrentBlock  *big.Int
	HighestBlock  *big.Int
	StartingBlock *big.Int
}
type SyncStatus struct {
	Status string
	Info   SyncInfo
}

func (n *Node) Syncing() (SyncStatus, error) {
    resultChan := make(chan SyncStatus)
    errorChan := make(chan error)

    go func() {
        headErrChan := make(chan error)
        blockNumberChan := make(chan *big.Int)
        highestBlockChan := make(chan uint64)

        go func() {
            headErrChan <- n.CheckHeadAge()
        }()

        go func() {
            blockNumber, err := n.GetBlockNumber()
            if err != nil {
                blockNumberChan <- big.NewInt(0)
            } else {
                blockNumberChan <- blockNumber
            }
        }()

        go func() {
            highestBlock := n.Consensus.Expected_current_slot()
            highestBlockChan <- highestBlock
        }()
        headErr := <-headErrChan
        if headErr == nil {
            resultChan <- SyncStatus{
                Status: "None",
            }
            errorChan <- nil
            return
        }
        latestSyncedBlock := <-blockNumberChan
        highestBlock := <-highestBlockChan

        resultChan <- SyncStatus{
            Info: SyncInfo{
                CurrentBlock:  latestSyncedBlock,
                HighestBlock:  big.NewInt(int64(highestBlock)),
                StartingBlock: big.NewInt(0),
            },
        }
        errorChan <- headErr
    }()

    select {
    case result := <-resultChan:
        return result, nil
    case err := <-errorChan:
        return SyncStatus{}, err
    }
}

func (n *Node) GetCoinbase() ([20]byte, error) {
    resultChan := make(chan [20]byte)
    errorChan := make(chan error)

    go func() {
        headErrChan := make(chan error)
        blockChan := make(chan *common.Block)

        go func() {
            headErrChan <- n.CheckHeadAge()
        }()

        go func() {
            block, err := n.Execution.GetBlock(common.BlockTag{}.Latest, false)
            if err != nil {
                errorChan <- err
                resultChan <- [20]byte{}
                return
            }
            blockChan <- block
        }()

        if headErr := <-headErrChan; headErr != nil {
            resultChan <- [20]byte{}
            errorChan <- headErr
            return
        }

        block := <-blockChan
        resultChan <- block.Miner.Addr
        errorChan <- nil
    }()

    select {
    case coinbase := <-resultChan:
        return coinbase, nil
    case err := <-errorChan:
        return [20]byte{}, err
    }
}

func (n *Node) CheckHeadAge() *NodeError {
    resultChan := make(chan *NodeError)

    go func() {
        currentTime := time.Now()
        currentTimestamp := currentTime.Unix()

        block, err := n.Execution.GetBlock(common.BlockTag{Latest: true}, false)
        if err != nil {
            resultChan <- &NodeError{}
            return
        }
        blockTimestamp := block.Timestamp.Unix()
        delay := currentTimestamp - blockTimestamp
        if delay > 60 {
            resultChan <- NewOutOfSyncError(delay)
            return
        }

        resultChan <- nil
    }()

    return <-resultChan
}

func (n *Node) CheckBlocktagAge(block common.BlockTag) *NodeError {
    errorChan := make(chan *NodeError)

    go func() {
        if block.Latest {
            headErr := n.CheckHeadAge()
            errorChan <- headErr
            return
        }
        if block.Finalized {
            errorChan <- nil
            return
        }
        errorChan <- nil
    }()
    return <-errorChan
}