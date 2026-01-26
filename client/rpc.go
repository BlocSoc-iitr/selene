package client

import (
	"encoding/json"
	"fmt"
	"github.com/BlocSoc-iitr/selene/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"log"
	"math/big"
	"net"
	"net/http"
	"sync"
)

type Rpc struct {
	Node    *Node
	Handle  *Server
	Address string
}

func (r Rpc) New(node *Node, ip *string, port uint16) Rpc {
	address := "127.0.0.1" //Default LocalHost value
	return Rpc{
		Node:    node,
		Handle:  nil,
		Address: address,
	}
}

func (r Rpc) Start() (string, error) {
    // Create channels to send the address and error back asynchronously
    addrChan := make(chan string)
    errChan := make(chan error)

    // Run the function in a separate goroutine
    go func() {
        rpcInner := RpcInner{
            Node:    r.Node,
            Address: r.Address,
        }
        handle, addr, err := rpcInner.StartServer()
        if err != nil {
            // Send the error on the error channel and close the channels
            errChan <- err
            close(addrChan)
            close(errChan)
            return
        }
        r.Handle = handle
        logrus.WithField("target", "selene::rpc").Infof("rpc server started at %s", addr)

        // Send the address on the address channel and close the channels
        addrChan <- addr
        close(addrChan)
        close(errChan)
    }()

    return <-addrChan, <-errChan
}

type RpcInner struct {
	Node    *Node
	Address string
}

func (r *RpcInner) GetBalance(address string, block common.BlockTag) (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)
    go func() {
        balance, err := r.Node.GetBalance(address, block)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- balance
    }()
    select {
    case balance := <-resultChan:
        return balance, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetTransactionCount(address string, block common.BlockTag) (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)

    go func() {
        txCount, err := r.Node.GetNonce(address, block)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- uint64(txCount)
    }()
    select {
    case txCount := <-resultChan:
        return txCount, nil
    case err := <-errorChan:
        return 0, err
    }
}
func (r *RpcInner) GetBlockTransactionCountByHash(hash [32]byte) (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)

    go func() {
        txCount, err := r.Node.GetBlockTransactionCountByHash(hash)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txCount
    }()
    select {
    case txCount := <-resultChan:
        return txCount, nil
    case err := <-errorChan:
        return 0, err
    }
}
func (r *RpcInner) GetBlockTransactionCountByNumber(block common.BlockTag) (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)

    go func() {
        txCount, err := r.Node.GetBlockTransactionCountByNumber(block)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txCount
 
	}()
    select {
    case txCount := <-resultChan:
        return txCount, nil
    case err := <-errorChan:
        return 0, err
    }
}
func (r *RpcInner) GetCode(address string, block common.BlockTag) ([]byte, error) {
    resultChan := make(chan []byte)
    errorChan := make(chan error)

    go func() {
        code, err := r.Node.GetCode(address, block)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- code
    	}()
    select {
    case code := <-resultChan:
        return code, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) Call(tx TransactionRequest, block common.BlockTag) ([]byte, error) {
    resultChan := make(chan []byte)
    errorChan := make(chan error)
 
	go func() {
        result, err := r.Node.Call(&tx, block)
        if err != nil {
            errorChan <- err.cause
            return
        }
        resultChan <- result
    
		}()
    select {
    case result := <-resultChan:
        return result, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) EstimateGas(tx *TransactionRequest) (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)
    go func() {
        gas, err := r.Node.EstimateGas(tx)
        if err != nil {
            errorChan <- err.cause
            return
        }
        resultChan <- gas
    }()
    select {
    case gas := <-resultChan:
        return gas, nil
    case err := <-errorChan:
        return 0, err
    }
}
func (r *RpcInner) ChainId() (uint64, error) {
	return r.Node.ChainId(), nil
}
func (r *RpcInner) GasPrice() (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        gasPrice, err := r.Node.GetGasPrice()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- gasPrice
    }()

    select {
    case gasPrice := <-resultChan:
        return gasPrice, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) MaxPriorityFeePerGas() *big.Int {
	return r.Node.GetPriorityFee()
}
func (r *RpcInner) BlockNumber() (uint64, error) {
    resultChan := make(chan uint64)
    errorChan := make(chan error)

    go func() {
        number, err := r.Node.GetBlockNumber()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- number.Uint64()
    }()

    select {
    case number := <-resultChan:
        return number, nil
    case err := <-errorChan:
        return 0, err
    }
}
func (r *RpcInner) GetBlockByNumber(blockTag common.BlockTag, fullTx bool) (*common.Block, error) {
    resultChan := make(chan *common.Block)
    errorChan := make(chan error)

    go func() {
        block, err := r.Node.GetBlockByNumber(blockTag, fullTx)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- &block
    }()

    select {
    case block := <-resultChan:
        return block, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetBlockByHash(hash [32]byte, fullTx bool) (*common.Block, error) {
    resultChan := make(chan *common.Block)
    errorChan := make(chan error)

    go func() {
        block, err := r.Node.GetBlockByHash(hash, fullTx)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- &block
    }()

    select {
    case block := <-resultChan:
        return block, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) SendRawTransaction(bytes []uint8) ([32]byte, error) {
    resultChan := make(chan [32]byte)
    errorChan := make(chan error)

    go func() {
        hash, err := r.Node.SendRawTransaction(&bytes)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- hash
    }()

    select {
    case hash := <-resultChan:
        return hash, nil
    case err := <-errorChan:
        return [32]byte{}, err
    }
}
func (r *RpcInner) GetTransactionReceipt(hash [32]byte) (*types.Receipt, error) {
    resultChan := make(chan *types.Receipt)
    errorChan := make(chan error)

    go func() {
        txnReceipt, err := r.Node.GetTransactionReceipt(hash)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txnReceipt
    }()

    select {
    case txnReceipt := <-resultChan:
        return txnReceipt, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetTransactionByHash(hash [32]byte) (*types.Transaction, error) {
    resultChan := make(chan *types.Transaction)
    errorChan := make(chan error)

    go func() {
        txn, err := r.Node.GetTransactionByHash(hash)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txn
    }()

    select {
    case txn := <-resultChan:
        return txn, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetTransactionByBlockHashAndIndex(hash [32]byte, index uint64) (*types.Transaction, error) {
    resultChan := make(chan *types.Transaction)
    errorChan := make(chan error)

    go func() {
        txn, err := r.Node.GetTransactionByBlockHashAndIndex(hash, index)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- txn
    }()

    select {
    case txn := <-resultChan:
        return txn, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) Coinbase() ([20]byte, error) {
    resultChan := make(chan [20]byte)
    errorChan := make(chan error)

    go func() {
        coinbase, err := r.Node.GetCoinbase()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- coinbase
    }()

    select {
    case coinbase := <-resultChan:
        return coinbase, nil
    case err := <-errorChan:
        return [20]byte{}, err
    }
}
func (r *RpcInner) Syncing() (SyncStatus, error) {
    resultChan := make(chan SyncStatus)
    errorChan := make(chan error)

    go func() {
        sync, err := r.Node.Syncing()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- sync
    }()
    select {
    case sync := <-resultChan:
        return sync, nil
    case err := <-errorChan:
        return SyncStatus{}, err
    }
}
func (r *RpcInner) GetLogs(filter ethereum.FilterQuery) ([]types.Log, error) {
    resultChan := make(chan []types.Log)
    errorChan := make(chan error)

    go func() {
        logs, err := r.Node.GetLogs(&filter)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- logs
    }()
    select {
    case logs := <-resultChan:
        return logs, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetFilterChanges(filterId *big.Int) ([]types.Log, error) {
    resultChan := make(chan []types.Log)
    errorChan := make(chan error)

    go func() {
        logs, err := r.Node.GetFilterChanges(filterId)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- logs
    }()
    select {
    case logs := <-resultChan:
        return logs, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) UninstallFilter(filterId *big.Int) (bool, error) {
    resultChan := make(chan bool)
    errorChan := make(chan error)

    go func() {
        boolean, err := r.Node.UninstallFilter(filterId)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- boolean
    }()
    select {
    case boolean := <-resultChan:
        return boolean, nil
    case err := <-errorChan:
        return false, err
    }
}
func (r *RpcInner) GetNewFilter(filter ethereum.FilterQuery) (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        filterId, err := r.Node.GetNewFilter(&filter)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- filterId
    }()
    select {
    case filterId := <-resultChan:
        return filterId, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetNewBlockFilter() (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        filterId, err := r.Node.GetNewBlockFilter()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- filterId
    }()
    select {
    case filterId := <-resultChan:
        return filterId, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetNewPendingTransactionFilter() (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)

    go func() {
        filterId, err := r.Node.GetNewPendingTransactionFilter()
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- filterId
    }()
    select {
    case filterId := <-resultChan:
        return filterId, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) GetStorageAt(address string, slot [32]byte, block common.BlockTag) (*big.Int, error) {
    resultChan := make(chan *big.Int)
    errorChan := make(chan error)
    go func() {
        value, err := r.Node.GetStorageAt(address, slot, block)
        if err != nil {
            errorChan <- err
            return
        }
        resultChan <- value
    }()
    select {
    case value := <-resultChan:
        return value, nil
    case err := <-errorChan:
        return nil, err
    }
}
func (r *RpcInner) Version() (uint64, error) {
	return r.Node.ChainId(), nil
}

// Define a struct to simulate JSON-RPC request/response
type JsonRpcRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     int             `json:"id"`
}
type JsonRpcResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	ID     int         `json:"id"`
}

// Server struct simulating server handle
type Server struct {
	Address string
	Methods map[string]http.HandlerFunc
	mu      sync.Mutex
}

// Start the server and listen for requests
func (r *RpcInner) StartServer() (*Server, string, error) {
	address := r.Address
	server := &Server{
		Address: address,
		Methods: map[string]http.HandlerFunc{},
	}
	// Use a listener to get the local address (port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, "", err
	}
	// Start the server asynchronously
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			rpcHandler(w, req)
		})
		if err := http.Serve(listener, nil); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return server, listener.Addr().String(), nil
}

type RPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSON-RPC response structure
type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   string      `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// Handle JSON-RPC requests
// Defines methods to call functions in the RPC Server
func rpcHandler(w http.ResponseWriter, r *http.Request) {
	var req RPCRequest
	var rpc RpcInner
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Handle the RPC methods
	switch req.Method {
	case "eth_getBalance":
		resultChan := make(chan *big.Int)
		errorChan := make(chan error)

		go func() {
			address := req.Params[0].(string)
			block := req.Params[1].(common.BlockTag)
			balance, err := rpc.GetBalance(address, block)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- balance
		}()

		select {
		case balance := <-resultChan:
			writeRPCResponse(w, req.ID, balance.String(), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getTransactionCount":
		resultChan := make(chan uint64)
		errorChan := make(chan error)

		go func() {
			address := req.Params[0].(string)
			block := req.Params[1].(common.BlockTag)
			nonce, err := rpc.GetTransactionCount(address, block)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- nonce
		}()

		select {
		case nonce := <-resultChan:
			writeRPCResponse(w, req.ID, fmt.Sprintf("%d", nonce), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getBlockTransactionCountByHash":
		resultChan := make(chan uint64)
		errorChan := make(chan error)

		go func() {
			hash := req.Params[0].([32]byte)
			count, err := rpc.GetBlockTransactionCountByHash(hash)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- count
		}()

		select {
		case count := <-resultChan:
			writeRPCResponse(w, req.ID, fmt.Sprintf("%d", count), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getBlockTransactionCountByNumber":
		resultChan := make(chan uint64)
		errorChan := make(chan error)

		go func() {
			block := req.Params[0].(common.BlockTag)
			nonce, err := rpc.GetBlockTransactionCountByNumber(block)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- nonce
		}()

		select {
		case nonce := <-resultChan:
			writeRPCResponse(w, req.ID, fmt.Sprintf("%d", nonce), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getCode":
		resultChan := make(chan []byte)
		errorChan := make(chan error)

		go func() {
			address := req.Params[0].(string)
			block := req.Params[1].(common.BlockTag)
			code, err := rpc.GetCode(address, block)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- code
		}()

		select {
		case code := <-resultChan:
			writeRPCResponse(w, req.ID, string(code), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_call":
		resultChan := make(chan []byte)
		errorChan := make(chan error)

		go func() {
			tx := req.Params[0].(TransactionRequest)
			block := req.Params[1].(common.BlockTag)
			result, err := rpc.Call(tx, block)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			writeRPCResponse(w, req.ID, string(result), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_estimateGas":
		resultChan := make(chan uint64)
		errorChan := make(chan error)

		go func() {
			tx := req.Params[0].(TransactionRequest)
			gas, err := rpc.EstimateGas(&tx)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- gas
		}()

		select {
		case gas := <-resultChan:
			writeRPCResponse(w, req.ID, fmt.Sprintf("%d", gas), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_chainId":
		resultChan := make(chan uint64)
		errorChan := make(chan error)

		go func() {
			id, err := rpc.ChainId()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- id
		}()

		select {
		case id := <-resultChan:
			writeRPCResponse(w, req.ID, fmt.Sprintf("%d", id), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_gasPrice":
		resultChan := make(chan *big.Int)
		errorChan := make(chan error)

		go func() {
			result, err := rpc.GasPrice()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			writeRPCResponse(w, req.ID, result.String(), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_maxPriorityFeePerGas":
		result := rpc.MaxPriorityFeePerGas()
		writeRPCResponse(w, req.ID, result.String(), "")
	case "eth_blockNumber":
		resultChan := make(chan uint64)
		errorChan := make(chan error)

		go func() {
			result, err := rpc.BlockNumber()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			writeRPCResponse(w, req.ID, fmt.Sprintf("%d", result), "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getBlockByNumber":
		resultChan := make(chan *common.Block)
		errorChan := make(chan error)

		go func() {
			block := req.Params[0].(common.BlockTag)
			fullTx := req.Params[1].(bool)
			result, err := rpc.GetBlockByNumber(block, fullTx)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			writeRPCResponse(w, req.ID, result, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getBlockByHash":
		resultChan := make(chan *common.Block)
		errorChan := make(chan error)

		go func() {
			hash := req.Params[0].([32]byte)
			fullTx := req.Params[1].(bool)
			result, err := rpc.GetBlockByHash(hash, fullTx)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			writeRPCResponse(w, req.ID, result, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_sendRawTransaction":
		resultChan := make(chan [32]byte)
		errorChan := make(chan error)

		go func() {
			rawTransaction := req.Params[0].([]uint8)
			hash, err := rpc.SendRawTransaction(rawTransaction)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- hash
		}()

		select {
		case hash := <-resultChan:
			writeRPCResponse(w, req.ID, hash, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getTransactionReceipt":
		resultChan := make(chan *types.Receipt)
		errorChan := make(chan error)

		go func() {
			hash := req.Params[0].([32]byte)
			receipt, err := rpc.GetTransactionReceipt(hash)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- receipt
		}()

		select {
		case receipt := <-resultChan:
			writeRPCResponse(w, req.ID, receipt, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getTransactionByHash":
		resultChan := make(chan *types.Transaction)
		errorChan := make(chan error)

		go func() {
			hash := req.Params[0].([32]byte)
			transaction, err := rpc.GetTransactionByHash(hash)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- transaction
		}()

		select {
		case transaction := <-resultChan:
			writeRPCResponse(w, req.ID, transaction, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getTransactionByBlockHashAndIndex":
		resultChan := make(chan *types.Transaction)
		errorChan := make(chan error)

		go func() {
			blockHash := req.Params[0].([32]byte)
			index := req.Params[1].(uint64)
			transaction, err := rpc.GetTransactionByBlockHashAndIndex(blockHash, index)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- transaction
		}()

		select {
		case transaction := <-resultChan:
			writeRPCResponse(w, req.ID, transaction, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getLogs":
		resultChan := make(chan []types.Log)
		errorChan := make(chan error)

		go func() {
			filter := req.Params[0].(ethereum.FilterQuery)
			logs, err := rpc.GetLogs(filter)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- logs
		}()

		select {
		case logs := <-resultChan:
			writeRPCResponse(w, req.ID, logs, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getFilterChanges":
		resultChan := make(chan []types.Log)
		errorChan := make(chan error)

		go func() {
			filterID := req.Params[0].(*big.Int)
			logs, err := rpc.GetFilterChanges(filterID)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- logs
		}()

		select {
		case logs := <-resultChan:
			writeRPCResponse(w, req.ID, logs, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_uninstallFilter":
		resultChan := make(chan bool)
		errorChan := make(chan error)

		go func() {
			filterID := req.Params[0].(*big.Int)
			result, err := rpc.UninstallFilter(filterID)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case result := <-resultChan:
			writeRPCResponse(w, req.ID, result, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_newFilter":
		resultChan := make(chan *big.Int)
		errorChan := make(chan error)

		go func() {
			filter := req.Params[0].(ethereum.FilterQuery)
			filterID, err := rpc.GetNewFilter(filter)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- filterID
		}()

		select {
		case filterID := <-resultChan:
			writeRPCResponse(w, req.ID, filterID, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_newBlockFilter":
		resultChan := make(chan *big.Int)
		errorChan := make(chan error)

		go func() {
			filterID, err := rpc.GetNewBlockFilter()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- filterID
		}()

		select {
		case filterID := <-resultChan:
			writeRPCResponse(w, req.ID, filterID, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_newPendingTransactionFilter":
		resultChan := make(chan *big.Int)
		errorChan := make(chan error)

		go func() {
			filterID, err := rpc.GetNewPendingTransactionFilter()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- filterID
		}()

		select {
		case filterID := <-resultChan:
			writeRPCResponse(w, req.ID, filterID, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_getStorageAt":
		resultChan := make(chan *big.Int)
		errorChan := make(chan error)

		go func() {
			address := req.Params[0].(string)
			slot := req.Params[1].([32]byte)
			blockTag := req.Params[2].(common.BlockTag)
			storage, err := rpc.GetStorageAt(address, slot, blockTag)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- storage
		}()

		select {
		case storage := <-resultChan:
			writeRPCResponse(w, req.ID, storage, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_coinbase":
		resultChan := make(chan [20]byte)
		errorChan := make(chan error)

		go func() {
			coinbaseAddress, err := rpc.Coinbase()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- coinbaseAddress
		}()

		select {
		case coinbaseAddress := <-resultChan:
			writeRPCResponse(w, req.ID, coinbaseAddress, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "eth_syncing":
		resultChan := make(chan SyncStatus)
		errorChan := make(chan error)

		go func() {
			syncStatus, err := rpc.Syncing()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- syncStatus
		}()

		select {
		case syncStatus := <-resultChan:
			writeRPCResponse(w, req.ID, syncStatus, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}
	case "net_version":
		resultChan := make(chan uint64)
		errorChan := make(chan error)

		go func() {
			result, err := rpc.Version()
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()

		select {
		case version := <-resultChan:
			writeRPCResponse(w, req.ID, version, "")
		case err := <-errorChan:
			writeRPCResponse(w, req.ID, nil, err.Error())
		}

	default:
		writeRPCResponse(w, req.ID, nil, "Method not found")
	}
}

// Write a JSON-RPC response
func writeRPCResponse(w http.ResponseWriter, id int, result interface{}, err string) {
	resp := RPCResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
		Error:   err,
	}
	json.NewEncoder(w).Encode(resp)
}