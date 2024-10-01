package execution

import (
	"bytes"
	"fmt"
	"math/big"
	seleneCommon "github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/rlp"
)

const MAX_SUPPORTED_LOGS_NUMBER = 5

type ExecutionClient struct {
	Rpc   ExecutionRpc
	state State
}

type ExecutionRpc interface{}

// State struct (assuming it has necessary methods)

// New creates a new ExecutionClient
func (e *ExecutionClient) New(rpc string, state State) (*ExecutionClient, error) {
	r, err := ExecutionRpc.New(rpc)
	if err != nil {
		return nil, err
	}
	return &ExecutionClient{
		Rpc:   r,
		state: state,
	}, nil
}

// CheckRpc checks the chain ID against the expected value
func (e *ExecutionClient) CheckRpc(chainID uint64) error {
	resultChan := make(chan struct {
		id  uint64
		err error
	})
	// Make sure the .ChainID() function is implemented in the rpc package
	go func() {
		rpcChainID, err := e.Rpc.ChainID()
		resultChan <- struct {
			id  uint64
			err error
		}{rpcChainID, err}
	}()
	result := <-resultChan
	if result.err != nil {
		return result.err
	}
	if result.id != chainID {
		return ErrIncorrectRpcNetwork // return IncorrectRpcNetworkError from the error file in execution package
	}
	return nil
}

// GetAccount retrieves the account information
func (e *ExecutionClient) GetAccount(address *seleneCommon.Address, slots H256, tag BlockTag) (Account, error) { //Account from execution/types.go
	block, err := e.state.GetBlock(tag) // GetBlock function is not defined yet
	if err != nil {
		return Account{}, ErrBlockNotFound // return BlockNotFoundError
	}
	proof, err := e.Rpc.GetProof(address, slots,block.Number) // GetProof function is not defined yet

	accountPath := crypto.Keccak256(address.Addr[:]) // Convert Address to bytes
	accountEncoded, err := EncodeAccount(proof) // Encode the proof
	isValid := VerifyProof(proof.AccountProof, block.StateRoot[:], accountPath, accountEncoded)
	if !isValid {
		return Account{}, InvalidAccountProof(address)
	}
	// modify 
	slotMap := make(map[common.Hash]*big.Int)
	for _, storageProof := range proof.StorageProof {
    	key, err := utils.Hex_str_to_bytes(storageProof.Key.Hex())
    	if err != nil {
        	return Account{}, err
    	}
		value := encode(storageProof.Value).Bytes()
		keyHash := crypto.Keccak256(key)
		isValid := VerifyProof(    // VerifyProof from execution/proof.go
			storageProof.Proof,
			proof.StorageHash.Bytes(),
			keyHash,
			value,
		)
		if !isValid {
			return Account{}, fmt.Errorf("invalid storage proof for address: %v, key: %v", *address, storageProof.Key)
		}
		slotMap[storageProof.Key] = storageProof.Value
	}
	var code []byte
	if bytes.Equal(proof.CodeHash.Bytes(), crypto.Keccak256(KECCAK_EMPTY)) {
    	code = []byte{}
	} else {
    	code, err := e.Rpc.GetCode(address, block.Number.Uint64())
    	if err != nil {
       		return Account{}, err
    	}
    	codeHash := crypto.Keccak256(code)
    	if !bytes.Equal(proof.CodeHash.Bytes(), codeHash) {
        	return Account{}, fmt.Errorf("code hash mismatch for address: %v, expected: %v, got: %v", 
            *address, common.BytesToHash(codeHash).String(), proof.CodeHash.String())
   		}
	}
	account := Account{
		Balance:     proof.Balance,
		Nonce:       proof.Nonce.Uint64(),
		Code:        code,
		CodeHash:    proof.CodeHash,
		StorageHash: proof.StorageHash,
		Slots:       slotMap,
	}
	return account, nil
}
//need to confirm the return type, whether it should return the channel directly or wait for the completion of the SendRawTransaction method
func (e *ExecutionClient) SendRawTransaction(bytes []byte) (common.Hash, error) {
	var txHash common.Hash
    var err error
    done := make(chan bool)
    go func() {
        txHash, err = e.Rpc.SendRawTransaction(bytes)
        done <- true
    }()
    <-done
    return txHash, err
}

func (e *ExecutionClient) GetBlock(tag seleneCommon.BlockTag,full_tx bool)(seleneCommon.Block, error) {
	blockChan := make(chan seleneCommon.Block)
    errChan := make(chan error)
    go func() {
        block, err := e.state.GetBlock(tag)
        if err != nil {
			errChan <- seleneCommon.BlockNotFoundError{Block: tag}
            return
        }
        blockChan <- block
    }()
    select {
    case block := <-blockChan:
        if !full_tx {
            block.Transactions = seleneCommon.Transactions{Hashes: block.Transactions.HashesFunc()}
        }
        return block, nil
    case err := <-errChan:
        return seleneCommon.Block{}, err
    }
}

func (e *ExecutionClient) GetBlockByHash(hash common.Hash,full_tx bool) (seleneCommon.Block,error){	
	blockChan := make(chan seleneCommon.Block)
    errChan := make(chan error)
    go func() {
        block, err := e.state.GetBlockByHash(hash)
        if err != nil {
			errChan <- err
            return
        }
        blockChan <- block
    }()
    select {
    case block := <-blockChan:
        if !full_tx {
            block.Transactions = seleneCommon.Transactions{Hashes: block.Transactions.HashesFunc()}
        }
        return block, nil
    case err := <-errChan:
        return seleneCommon.Block{}, err
    }
}
// don't know what to do with the `future` implementation in helios, right now i have kept it simple  
func (e *ExecutionClient) GetTransactionByBlockHashAndIndex(blockHash common.Hash,index uint64) (seleneCommon.Transaction,error){
	txChan := make(chan seleneCommon.Transaction)
	errChan := make(chan error)
	go func() {
		tx, err := e.state.GetTransactionByBlockHashAndIndex(blockHash,index)
		if err != nil {
			errChan <- err
			return
		}
		txChan <- tx
	}()
	select {
	case tx := <-txChan:
		return tx, nil
	case err := <-errChan:
		return seleneCommon.Transaction{}, err
	}

}

// comparing the fetched transaction with the expected transaction from the receipt is left to be impelemented
func (e *ExecutionClient) GetTransactionReceipt(txHash common.Hash) (types.Receipt,error) {
	receiptChan := make(chan types.Receipt)
	errChan := make(chan error)
	// var receipt types.Receipt
	go func() {
		receipt, err := e.state.GetTransactionReceipt(txHash)
		if err != nil {
			errChan <- err
			return
		}
		receiptChan <- receipt
	}()
	select {
	case receipt := <-receiptChan:
		blocknumber:=receipt.BlockNumber
		blockChan := make(chan seleneCommon.Block)
    	errChan := make(chan error)
    	go func() {
			block, err := e.state.GetBlock(seleneCommon.BlockTag{Number: blocknumber.Uint64()})
			if err != nil {
				errChan <- seleneCommon.BlockNotFoundError{Block: seleneCommon.BlockTag{Number: blocknumber.Uint64()}}
				return
			}
			blockChan <- block
    	}()
    	select {
    	case block := <-blockChan:
        	txHashes := block.Transactions.Hashes
			receiptsChan := make(chan types.Receipt)
			receiptsErrChan := make(chan error)
			for _, hash := range txHashes {
				go func(hash common.Hash) {
					receipt, err := e.Rpc.GetTransactionReceipt(hash)
					if err != nil {
						receiptsErrChan <- err
						return
					}
					if receipt == nil {
						receiptsErrChan <- fmt.Errorf("not reachable")
						return
					}
					receiptsChan <- *receipt
				}(hash)
			}
			var receipts []types.Receipt
			for range txHashes {
				select {
				case receipt := <-receiptsChan:
					receipts = append(receipts, receipt)
				case err := <-receiptsErrChan:
					return types.Receipt{}, err
				}
			}
			var receiptsEncoded [][]byte
			for _, receipt := range receipts {
				receiptsEncoded = append(receiptsEncoded, encodeReceipt(receipt))
			}
			expectedReceiptRoot := trie.(receiptsEncoded)
			expectedReceiptRoot = common.BytesToHash(expectedReceiptRoot)

			if expectedReceiptRoot != block.ReceiptsRoot || !contains(receipts, receipt) {
				return types.Receipt{}, fmt.Errorf("receipt root mismatch: %s", txHash.String())
			}

			return receipt, nil // Return the found receipt

    	case err := <-errChan:
        	return types.Receipt{}, err
    	}		
	case err := <-errChan:
		return types.Receipt{}, err
	}


}

func (e *ExecutionClient) GetTransaction(hash common.Hash) (seleneCommon.Transaction, error) {
    txChan := make(chan seleneCommon.Transaction)
    errChan := make(chan error)
    go func() {
        tx, err := e.state.GetTransaction(hash)
        if err != nil {
            errChan <- err
            return
        }
        txChan <- tx
    }()

    // Use select to wait for the result from the channel
    select {
    case tx := <-txChan:
        return tx, nil
    case err := <-errChan:
        return seleneCommon.Transaction{}, err
    }
}
//what to do with filter, should i define it on my own as it is not defined in geth
func (e *ExecutionClient) GetLogs(filter ??) ([]Log, error) {
    clonedFilter := filter.Clone()
    if clonedFilter.GetToBlock() == nil && clonedFilter.GetBlockHash() == nil {
        block, err := e.state.LatestBlockNumber()
        if err != nil {
            return nil, err
        }
        clonedFilter = clonedFilter.ToBlock(block)
        if clonedFilter.GetFromBlock() == nil {
            clonedFilter = clonedFilter.FromBlock(block)
        }
    }
    logsChan := make(chan []Log)
    errChan := make(chan error)
    go func() {
        logs, err := e.rpc.GetLogs(clonedFilter)
        if err != nil {
            errChan <- err
            return
        }
        logsChan <- logs
    }()
    select {
    case logs := <-logsChan:
        if len(logs) > MAX_SUPPORTED_LOGS_NUMBER {
            return nil, ExecutionError{Message: fmt.Sprintf("Too many logs to prove: %d, max: %d", len(logs), MAX_SUPPORTED_LOGS_NUMBER)}
        }
        if err := e.VerifyLogs(logs); err != nil {
            return nil, err
        }

        return logs, nil
    case err := <-errChan:
        return nil, err
    }
}


func (e *ExecutionClient) GetFilterChanges() {}

func (e *ExecutionClient) UninstallFilter() {}

func (e *ExecutionClient) GetNewFilter() {}

func (e *ExecutionClient) GetNewBlockFilter() {}

func (e *ExecutionClient) GetNewPendingTransactionFilter() {}

// verifyLogs verifies that each log is present in the corresponding receipt
func (e *ExecutionClient) verifyLogs(logs []*types.Log) error {
    errChan := make(chan error, len(logs))            
    receiptChan := make(chan *types.Receipt, len(logs))  
    for _, log := range logs {
        go func(log *types.Log) {
            if log.TxHash == (common.Hash{}) {
                errChan <- errors.New("tx hash not found in log")
                return
            }
            receiptSubChan := make(chan *types.Receipt)
            go func() {
                receipt, err := e.Rpc.GetTransactionReceipt(log.TxHash) 
                if err != nil {
                    errChan <- err
                    return
                }
                receiptSubChan <- receipt
            }()
            select {
            case receipt := <-receiptSubChan:
                if receipt == nil {
                    errChan <- errors.New("no receipt for transaction " + log.TxHash.Hex())
                    return
                }
                // Check if the receipt contains the desired log
                receiptLogsEncoded := make([][]byte, len(receipt.Logs))
                for i, receiptLog := range receipt.Logs {
                    receiptLogsEncoded[i] = receiptLog.Data // Use log.Data for encoded representation
                }
                logEncoded := log.Data // Check against the log's data
                found := false
                for _, encoded := range receiptLogsEncoded {
                    if string(encoded) == string(logEncoded) { // Compare as byte slices
                        found = true
                        break
                    }
                }
                if !found {
                    errChan <- errors.New("missing log for transaction " + log.TxHash.Hex())
                }
            case err := <-errChan:
                errChan <- err
            }
        }(log)
    }
    for range logs {
        select {
        case err := <-errChan:
            return err
        }
    }
    return nil
}

//not tested yet
func encodeReceipt(receipt *types.Receipt) ([]byte, error) {
    var stream []interface{}
    stream = append(stream, receipt.Status, receipt.CumulativeGasUsed, receipt.LogsBloom, receipt.Logs)
    legacyReceiptEncoded, err := rlp.EncodeToBytes(stream)
    if err != nil {
        return nil, err
    }
    txType := receipt.TransactionType
    if txType == nil {
        return legacyReceiptEncoded, nil
    }
    if *txType == 0 {
        return legacyReceiptEncoded, nil
    }
    txTypeBytes := []byte{*txType}
    return append(txTypeBytes, legacyReceiptEncoded...), nil
}
// need to confirm if TxHash is actually used as the key to calculate the receipt root or not
func CalculateReceiptRoot(receipts []types.Receipt) (common.Hash, error) {
    t,err := trie.New()
	if(err!=nil){
		return common.Hash{},err
	}
    for _, receipt := range receipts {
        key := receipt.TxHash.Bytes()
        value := encodeReceipt(receipt) 
        if err := t.Update(key, value); err != nil {
            return common.Hash{}, err
        }
    }

    return t.Hash(), nil
}
