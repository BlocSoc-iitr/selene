package execution

import (
	"bytes"
	"fmt"
	"math/big"

	"errors"

	seleneCommon "github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

const MAX_SUPPORTED_LOGS_NUMBER = 5
const KECCAK_EMPTY = "0x"

type ExecutionClient struct {
	Rpc   ExecutionRpc
	state *State
}

func (e *ExecutionClient) New(rpc string, state *State) (*ExecutionClient, error) {
	// This is updated as earlier, when ExecutionRpc.New() was called, it was giving an
	// invalid address or nil pointer dereference error as there wasn't a concrete type that implemented ExecutionRpc
	var r ExecutionRpc
	r, err := (&HttpRpc{}).New(&rpc)
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
	go func() {
		rpcChainID, err := e.Rpc.ChainId()
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
		return NewIncorrectRpcNetworkError()
	}
	return nil
}

// GetAccount retrieves the account information
func (e *ExecutionClient) GetAccount(address *seleneCommon.Address, slots *[]common.Hash, tag seleneCommon.BlockTag) (Account, error) { //Account from execution/types.go
	// block := e.state.GetBlock(tag)
	// Error Handling
	proof, err := e.Rpc.GetProof(address, slots, 21093292) // block.Number instead of hardcoded value
	fmt.Printf("%v", proof)
	if err != nil {
		return Account{}, err
	}
	accountPath := crypto.Keccak256(address.Addr[:])
	accountEncoded, err := EncodeAccount(&proof)
	if err != nil {
		return Account{}, err
	}
	accountProofBytes := make([][]byte, len(proof.AccountProof))
	for i, hexByte := range proof.AccountProof {
		accountProofBytes[i] = hexByte
	}
	fmt.Printf("%v", accountProofBytes)
	isValid, err := VerifyProof(accountProofBytes, common.Hex2Bytes("0xaba5664183bc9b0bbb9d092bc85cf895ab729b9d5ff98f4055e8869e8d948ad4"), accountPath, accountEncoded) //block.StateRoot[:] instead of hardcoded value
	if err != nil {
		return Account{}, err
	}
	if !isValid {
		return Account{}, NewInvalidAccountProofError(address.Addr)
	}
	// modify
	slotMap := []Slot{}
	for _, storageProof := range proof.StorageProof {
		key, err := utils.Hex_str_to_bytes(storageProof.Key.Hex())
		if err != nil {
			return Account{}, err
		}
		value, err := rlp.EncodeToBytes(storageProof.Value)
		if err != nil {
			return Account{}, err
		}
		keyHash := crypto.Keccak256(key)
		proofBytes := make([][]byte, len(storageProof.Proof))
		for i, hexByte := range storageProof.Proof {
			proofBytes[i] = hexByte
		}
		isValid, err := VerifyProof(
			proofBytes,
			proof.StorageHash.Bytes(),
			keyHash,
			value,
		)
		if err != nil {
			return Account{}, err
		}
		if !isValid {
			return Account{}, fmt.Errorf("invalid storage proof for address: %v, key: %v", *address, storageProof.Key)
		}
		slotMap = append(slotMap, Slot{
			Key:   storageProof.Key,
			Value: storageProof.Value.ToBig(),
		})
	}
	var code []byte
	if proof.CodeHash == [32]byte{} {
		code = []byte{}
	} else {
		code, err := e.Rpc.GetCode(address, 21093292) // block.Number instead of hardcoded value
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
		Balance:     proof.Balance.ToBig(),
		Nonce:       uint64(proof.Nonce),
		Code:        code,
		CodeHash:    proof.CodeHash,
		StorageHash: proof.StorageHash,
		Slots:       slotMap,
	}
	return account, nil
}
func (e *ExecutionClient) SendRawTransaction(bytes []byte) (common.Hash, error) {
	var txHash common.Hash
	var err error
	done := make(chan bool)
	go func() {
		txHash, err = e.Rpc.SendRawTransaction(&bytes)
		done <- true
	}()
	<-done
	return txHash, err
}
func (e *ExecutionClient) GetBlock(tag seleneCommon.BlockTag, full_tx bool) (seleneCommon.Block, error) {
	blockChan := make(chan seleneCommon.Block)
	errChan := make(chan error)
	go func() {
		block := e.state.GetBlock(tag)
		blockChan <- *block
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
func (e *ExecutionClient) GetBlockByHash(hash common.Hash, full_tx bool) (seleneCommon.Block, error) {
	blockChan := make(chan seleneCommon.Block)
	errChan := make(chan error)
	go func() {
		block := e.state.GetBlockByHash(hash)
		blockChan <- *block
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
func (e *ExecutionClient) GetTransactionByBlockHashAndIndex(blockHash common.Hash, index uint64) (seleneCommon.Transaction, error) {
	txChan := make(chan seleneCommon.Transaction)
	errChan := make(chan error)
	go func() {
		tx := e.state.GetTransactionByBlockAndIndex(blockHash, index)
		txChan <- *tx
	}()
	select {
	case tx := <-txChan:
		return tx, nil
	case err := <-errChan:
		return seleneCommon.Transaction{}, err
	}

}
func (e *ExecutionClient) GetTransactionReceipt(txHash common.Hash) (types.Receipt, error) {
	receiptChan := make(chan types.Receipt)
	errChan := make(chan error)
	// var receipt types.Receipt
	go func() {
		receipt, err := e.Rpc.GetTransactionReceipt(&txHash)
		if err != nil {
			errChan <- err
			return
		}
		receiptChan <- receipt
	}()
	select {
	case receipt := <-receiptChan:
		blocknumber := receipt.BlockNumber
		blockChan := make(chan seleneCommon.Block)
		errChan := make(chan error)
		go func() {
			block := e.state.GetBlock(seleneCommon.BlockTag{Number: blocknumber.Uint64()})
			blockChan <- *block
		}()
		select {
		case block := <-blockChan:
			txHashes := block.Transactions.Hashes
			receiptsChan := make(chan types.Receipt)
			receiptsErrChan := make(chan error)
			for _, hash := range txHashes {
				go func(hash common.Hash) {
					receipt, err := e.Rpc.GetTransactionReceipt(&hash)
					if err != nil {
						receiptsErrChan <- err
						return
					}
					receiptsChan <- receipt
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
				encodedReceipt, err := encodeReceipt(&receipt)
				if err != nil {
					receiptsErrChan <- err
					return types.Receipt{}, err
				}
				receiptsEncoded = append(receiptsEncoded, encodedReceipt)
			}
			expectedReceiptRoot, err := CalculateReceiptRoot(receiptsEncoded)
			if err != nil {
				return types.Receipt{}, err
			}

			if [32]byte(expectedReceiptRoot.Bytes()) != block.ReceiptsRoot || !contains(receipts, receipt) {
				return types.Receipt{}, fmt.Errorf("receipt root mismatch: %s", txHash.String())
			}

			return receipt, nil

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
		tx := e.state.GetTransaction(hash)
		txChan <- *tx
	}()
	select {
	case tx := <-txChan:
		return tx, nil
	case err := <-errChan:
		return seleneCommon.Transaction{}, err
	}
}
func (e *ExecutionClient) GetLogs(filter ethereum.FilterQuery) ([]types.Log, error) {
	if filter.ToBlock == nil && filter.BlockHash == nil {
		block := e.state.LatestBlockNumber()
		filter.ToBlock = new(big.Int).SetUint64(*block)
		if filter.FromBlock == nil {
			filter.FromBlock = new(big.Int).SetUint64(*block)
		}
	}
	logsChan := make(chan []types.Log)
	errChan := make(chan error)
	go func() {
		logs, err := e.Rpc.GetLogs(&filter)
		if err != nil {
			errChan <- err
			return
		}
		logsChan <- logs
	}()
	select {
	case logs := <-logsChan:
		if len(logs) > MAX_SUPPORTED_LOGS_NUMBER {
			// The earlier error was not returning properly
			return nil, errors.New("logs exceed max supported logs number")
		}
		logPtrs := make([]*types.Log, len(logs))
		for i := range logs {
			logPtrs[i] = &logs[i]
		}
		if err := e.verifyLogs(logPtrs); err != nil {
			return nil, err
		}

		return logs, nil
	case err := <-errChan:
		return nil, err
	}
}
func (e *ExecutionClient) GetFilterChanges(filterID *uint256.Int) ([]types.Log, error) {
	logsChan := make(chan []types.Log)
	errChan := make(chan error)
	go func() {
		logs, err := e.Rpc.GetFilterChanges(filterID)
		if err != nil {
			errChan <- err
			return
		}
		logsChan <- logs
	}()
	select {
	case logs := <-logsChan:
		if len(logs) > MAX_SUPPORTED_LOGS_NUMBER {
			return nil, &ExecutionError{
				Kind:    "TooManyLogs",
				Details: fmt.Sprintf("Too many logs to prove: %d, max: %d", len(logs), MAX_SUPPORTED_LOGS_NUMBER),
			}
		}
		logPtrs := make([]*types.Log, len(logs))
		for i := range logs {
			logPtrs[i] = &logs[i]
		}
		if err := e.verifyLogs(logPtrs); err != nil {
			return nil, err
		}
		return logs, nil
	case err := <-errChan:
		return nil, err
	}
}
func (e *ExecutionClient) UninstallFilter(filterID *uint256.Int) (bool, error) {
	resultChan := make(chan struct {
		result bool
		err    error
	})
	go func() {
		result, err := e.Rpc.UninstallFilter(filterID)
		resultChan <- struct {
			result bool
			err    error
		}{result, err}
	}()
	result := <-resultChan
	return result.result, result.err
}
func (e *ExecutionClient) GetNewFilter(filter ethereum.FilterQuery) (uint256.Int, error) {
	if filter.ToBlock == nil && filter.BlockHash == nil {
		block := e.state.LatestBlockNumber()
		filter.ToBlock = new(big.Int).SetUint64(*block)
		if filter.FromBlock == nil {
			filter.FromBlock = new(big.Int).SetUint64(*block)
		}
	}
	filterIDChan := make(chan uint256.Int)
	errChan := make(chan error)
	go func() {
		filterID, err := e.Rpc.GetNewFilter(&filter)
		if err != nil {
			errChan <- err
			return
		}
		filterIDChan <- filterID
	}()
	select {
	case filterID := <-filterIDChan:
		return filterID, nil
	case err := <-errChan:
		return uint256.Int{}, err
	}
}
func (e *ExecutionClient) GetNewBlockFilter() (uint256.Int, error) {
	filterIDChan := make(chan uint256.Int)
	errChan := make(chan error)
	go func() {
		filterID, err := e.Rpc.GetNewBlockFilter()
		if err != nil {
			errChan <- err
			return
		}
		filterIDChan <- filterID
	}()
	select {
	case filterID := <-filterIDChan:
		return filterID, nil
	case err := <-errChan:
		return uint256.Int{}, err
	}
}
func (e *ExecutionClient) GetNewPendingTransactionFilter() (uint256.Int, error) {
	filterIDChan := make(chan uint256.Int)
	errChan := make(chan error)
	go func() {
		filterID, err := e.Rpc.GetNewPendingTransactionFilter()
		if err != nil {
			errChan <- err
			return
		}
		filterIDChan <- filterID
	}()
	select {
	case filterID := <-filterIDChan:
		return filterID, nil
	case err := <-errChan:
		return uint256.Int{}, err
	}
}
func (e *ExecutionClient) verifyLogs(logs []*types.Log) error {
	errChan := make(chan error, len(logs))
	for _, log := range logs {
		go func(log *types.Log) {
			receiptSubChan := make(chan *types.Receipt)
			go func() {
				receipt, err := e.Rpc.GetTransactionReceipt(&log.TxHash)
				if err != nil {
					errChan <- err
					return
				}
				receiptSubChan <- &receipt
			}()
			select {
			case receipt := <-receiptSubChan:
				receiptLogsEncoded := make([][]byte, len(receipt.Logs))
				for i, receiptLog := range receipt.Logs {
					receiptLogsEncoded[i] = receiptLog.Data
				}
				logEncoded := log.Data
				found := false
				for _, encoded := range receiptLogsEncoded {
					if string(encoded) == string(logEncoded) {
						found = true
						break
					}
				}
				if !found {
					errChan <- fmt.Errorf("missing log for transaction %s", log.TxHash.Hex())
					return
				}
			case err := <-errChan:
				errChan <- err
				return
			}
			errChan <- nil
		}(log)
	}
	for range logs {
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}
func encodeReceipt(receipt *types.Receipt) ([]byte, error) {
	var stream []interface{}
	stream = append(stream, receipt.Status, receipt.CumulativeGasUsed, receipt.Bloom, receipt.Logs)
	legacyReceiptEncoded, err := rlp.EncodeToBytes(stream)
	if err != nil {
		return nil, err
	}
	txType := &receipt.Type
	if *txType == 0 {
		return legacyReceiptEncoded, nil
	}
	txTypeBytes := []byte{*txType}
	return append(txTypeBytes, legacyReceiptEncoded...), nil
}

// need to confirm if TxHash is actually used as the key to calculate the receipt root or not
func CalculateReceiptRoot(receipts [][]byte) (common.Hash, error) {
	if len(receipts) == 0 {
		return common.Hash{}, errors.New("no receipts to calculate root")
	}

	var receiptHashes []common.Hash
	for _, receipt := range receipts {
		receiptHash, err := rlpHash(receipt)
		if err != nil {
			return common.Hash{}, err
		}
		receiptHashes = append(receiptHashes, receiptHash)
	}
	return calculateMerkleRoot(receiptHashes), nil
}
func rlpHash(obj interface{}) (common.Hash, error) {
	encoded, err := rlp.EncodeToBytes(obj)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(encoded), nil
}

// This function is updated as it was going in an infinite loop
func calculateMerkleRoot(hashes []common.Hash) common.Hash {
	switch len(hashes) {
	case 0:
		return common.Hash{} // Return empty hash for empty slice
	case 1:
		return hashes[0]
	default:
		if len(hashes)%2 != 0 {
			hashes = append(hashes, hashes[len(hashes)-1])
		}
		var newLevel []common.Hash
		for i := 0; i < len(hashes); i += 2 {
			combinedHash := crypto.Keccak256(append(hashes[i].Bytes(), hashes[i+1].Bytes()...))
			newLevel = append(newLevel, common.BytesToHash(combinedHash))
		}
		return calculateMerkleRoot(newLevel)
	}
}

// contains checks if a receipt is in the list of receipts
func contains(receipts []types.Receipt, receipt types.Receipt) bool {
	for _, r := range receipts {
		if r.TxHash == receipt.TxHash {
			return true
		}
	}
	return false
}
