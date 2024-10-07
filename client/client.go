package client

import (
	"errors"
	"log"
	"math/big"
	"runtime"
	"time"
	"github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)
type ClientBuilder struct {
	Network              *config.Network
	ConsensusRpc         *string
	ExecutionRpc         *string
	Checkpoint           *[32]byte
	//If architecture is Web Assembly (wasm32), this will be nil
	RpcBindIp            *string
	//If architecture is Web Assembly (wasm32), this will be nil
	RpcPort              *uint16
	//If architecture is Web Assembly (wasm32), this will be nil
	DataDir              *string
	Config               *config.Config
	Fallback             *string
	LoadExternalFallback bool
	StrictCheckpointAge  bool
}
func (c Client) Build(b ClientBuilder) (*Client, error) {
	var baseConfig config.BaseConfig
	var err error
	_ = err
	if b.Network != nil {
		baseConfig, err = b.Network.BaseConfig(string(*b.Network))
	} else {
		if b.Config == nil {
			return nil, errors.New("missing network config")
		}
		baseConfig = b.Config.ToBaseConfig()
	}
	consensusRPC := b.ConsensusRpc
	if consensusRPC == nil {
		if b.Config == nil {
			return nil, errors.New("missing consensus rpc")
		}
		consensusRPC = &b.Config.ConsensusRpc
	}
	executionRPC := b.ExecutionRpc
	if executionRPC == nil {
		if b.Config == nil {
			return nil, errors.New("missing execution rpc")
		}
		executionRPC = &b.Config.ExecutionRpc
	}
	checkpoint := b.Checkpoint
	if checkpoint == nil {
		if b.Config != nil {
			checkpoint = b.Config.Checkpoint
		}
	}
	var defaultCheckpoint [32]byte
	if b.Config != nil {
		defaultCheckpoint = b.Config.DefaultCheckpoint
	} else {
		defaultCheckpoint = baseConfig.DefaultCheckpoint
	}
	var rpcBindIP, dataDir *string
	var rpcPort *uint16
	if runtime.GOARCH == "wasm" {
		rpcBindIP = nil
	} else {
		if b.RpcBindIp != nil {
			rpcBindIP = b.RpcBindIp
		} else if b.Config != nil {
			rpcBindIP = b.Config.RpcBindIp
		}
	}
	if runtime.GOARCH == "wasm" {
		rpcPort = nil
	} else {
		if b.RpcPort != nil {
			rpcPort = b.RpcPort
		} else if b.Config != nil {
			rpcPort = b.Config.RpcPort
		}
	}
	if runtime.GOARCH == "wasm" {
		dataDir = nil
	} else {
		if b.DataDir != nil {
			dataDir = b.DataDir
		} else if b.Config != nil {
			dataDir = b.Config.DataDir
		}
	}
	fallback := b.Fallback
	if fallback == nil && b.Config != nil {
		fallback = b.Config.Fallback
	}
	loadExternalFallback := b.LoadExternalFallback
	if b.Config != nil {
		loadExternalFallback = loadExternalFallback || b.Config.LoadExternalFallback
	}
	strictCheckpointAge := b.StrictCheckpointAge
	if b.Config != nil {
		strictCheckpointAge = strictCheckpointAge || b.Config.StrictCheckpointAge
	}
	config := config.Config{
		ConsensusRpc:         *consensusRPC,
		ExecutionRpc:         *executionRPC,
		Checkpoint:           checkpoint,
		DefaultCheckpoint:    defaultCheckpoint,
		RpcBindIp:            rpcBindIP,
		RpcPort:              rpcPort,
		DataDir:              dataDir,
		Chain:                baseConfig.Chain,
		Forks:                baseConfig.Forks,
		MaxCheckpointAge:     baseConfig.MaxCheckpointAge,
		Fallback:             fallback,
		LoadExternalFallback: loadExternalFallback,
		StrictCheckpointAge:  strictCheckpointAge,
		DatabaseType:         nil,
	}
	client, err := Client{}.New(config)
	return &client, nil
}
type Client struct {
	Node *Node
	Rpc  *Rpc
}
func (c Client) New(clientConfig config.Config) (Client, error) {
	config := &clientConfig
	node, err := NewNode(config)
	var rpc Rpc
	if err != nil {
		return Client{}, err
	}
	if config.RpcBindIp != nil || config.RpcPort != nil {
		rpc = Rpc{}.New(node, config.RpcBindIp, *config.RpcPort)
	}
	return Client{
		Node: node,
		Rpc: func() *Rpc {
			if runtime.GOARCH == "wasm" {
				return nil
			} else {
				return &rpc
			}
		}(),
	}, nil
}
func (c *Client) Start() error {
	errorChan := make(chan error, 1)
	go func() {
		defer close(errorChan)
		if runtime.GOARCH != "wasm" {
			if c.Rpc == nil {
				errorChan <- errors.New("Rpc not found.")
				return
			}
			_, err := c.Rpc.Start()
			if err != nil {
				errorChan <- err
				return
			}
		}
		errorChan <- nil
	}()
	err := <-errorChan
	return err
}
func (c *Client) Shutdown() {
	log.Println("selene::client - shutting down")
	go func() {
		err := c.Node.Consensus.Shutdown()
		if err != nil {
			log.Printf("selene::client - graceful shutdown failed: %v\n", err)
		}
	}()
}
func (c *Client) Call(tx *TransactionRequest, block common.BlockTag) ([]byte, error) {
	resultChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.Call(tx, block)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) EstimateGas(tx *TransactionRequest) (uint64, error) {
	resultChan := make(chan uint64, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.EstimateGas(tx)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return 0, err
	}
}
func (c *Client) GetBalance(address string, block common.BlockTag) (*big.Int, error) {
	resultChan := make(chan *big.Int, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetBalance(address, block)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetNonce(address string, block common.BlockTag) (uint64, error) {
	resultChan := make(chan uint64, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetNonce(address, block)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return 0, err
	}
}
func (c *Client) GetBlockTransactionCountByHash(hash [32]byte) (uint64, error) {
	resultChan := make(chan uint64, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetBlockTransactionCountByHash(hash)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return 0, err
	}
}
func (c *Client) GetBlockTransactionCountByNumber(block common.BlockTag) (uint64, error) {
	resultChan := make(chan uint64, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetBlockTransactionCountByNumber(block)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return 0, err
	}
}
func (c *Client) GetCode(address string, block common.BlockTag) ([]byte, error) {
	resultChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetCode(address, block)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetStorageAt(address string, slot [32]byte, block common.BlockTag) (*big.Int, error) {
	resultChan := make(chan *big.Int, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetStorageAt(address, slot, block)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) SendRawTransaction(bytes *[]uint8) ([32]byte, error) {
	resultChan := make(chan [32]byte, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.SendRawTransaction(bytes)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return [32]byte{}, err
	}
}
func (c *Client) GetTransactionReceipt(txHash [32]byte) (*types.Receipt, error) {
	resultChan := make(chan *types.Receipt, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetTransactionReceipt(txHash)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetTransactionByHash(txHash [32]byte) (*types.Transaction, error) {
	resultChan := make(chan *types.Transaction, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetTransactionByHash(txHash)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetLogs(filter *ethereum.FilterQuery) ([]types.Log, error) {
	resultChan := make(chan []types.Log, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetLogs(filter)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetFilterChanges(filterId *big.Int) ([]types.Log, error) {
	resultChan := make(chan []types.Log, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetFilterChanges(filterId)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) UninstallFilter(filterId *big.Int) (bool, error) {
	resultChan := make(chan bool, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.UninstallFilter(filterId)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return false, err
	}
}
func (c *Client) GetNewFilter(filter *ethereum.FilterQuery) (*big.Int, error) {
	resultChan := make(chan *big.Int, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetNewFilter(filter)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetNewBlockFilter() (*big.Int, error) {
	resultChan := make(chan *big.Int, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetNewBlockFilter()
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetNewPendingTransactionFilter() (*big.Int, error) {
	resultChan := make(chan *big.Int, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetNewPendingTransactionFilter()
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetGasPrice() (*big.Int, error) {
	resultChan := make(chan *big.Int, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetGasPrice()
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetPriorityFee() *big.Int {
	return c.Node.GetPriorityFee()
}
func (c *Client) GetBlockNumber() (*big.Int, error) {
	resultChan := make(chan *big.Int, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetBlockNumber()
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetBlockByNumber(block common.BlockTag, fullTx bool) (*common.Block, error) {
	resultChan := make(chan *common.Block, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetBlockByNumber(block, fullTx)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- &result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetBlockByHash(hash [32]byte, fullTx bool) (*common.Block, error) {
	resultChan := make(chan *common.Block, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetBlockByHash(hash, fullTx)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- &result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetTransactionByBlockHashAndIndex(blockHash [32]byte, index uint64) (*types.Transaction, error) {
	resultChan := make(chan *types.Transaction, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetTransactionByBlockHashAndIndex(blockHash, index)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) ChainID() uint64 {
	return c.Node.ChainId()
}
func (c *Client) Syncing() (*SyncStatus, error) {
	resultChan := make(chan *SyncStatus, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.Syncing()
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- &result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}
func (c *Client) GetCoinbase() ([20]byte, error) {
	resultChan := make(chan [20]byte, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := c.Node.GetCoinbase()
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return [20]byte{}, err
	}
}
func (c *Client) WaitSynced() error {
	for {
		statusChan := make(chan *SyncStatus)
		errChan := make(chan error)
		go func() {
			status, err := c.Syncing()
			if err != nil {
				errChan <- err
			} else {
				statusChan <- status
			}
		}()
		select {
		case status := <-statusChan:
			if status.Status == "None" {
				return nil
			}
		case err := <-errChan:
			return err
		case <-time.After(100 * time.Millisecond):
			// Periodic checking of SyncStatus or error
		}
	}
}