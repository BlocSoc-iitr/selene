package execution

import (
	"bytes"
	"encoding/json"
	"math/big"
	"sync"
	"testing"

	seleneCommon "github.com/BlocSoc-iitr/selene/common"
	// "github.com/BlocSoc-iitr/selene/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func CreateNewState() *State {
	blockChan := make(chan *seleneCommon.Block)
	finalizedBlockChan := make(chan *seleneCommon.Block)

	state := NewState(5, blockChan, finalizedBlockChan)

	// Simulate blocks to push
	block1 := &seleneCommon.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x11}, {0x12}},
		},
	}
	block2 := &seleneCommon.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
			Full: []seleneCommon.Transaction{
				{
					Hash:         common.Hash([32]byte{0x21}),
					GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
					Gas:          hexutil.Uint64(5),
					MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
				},
				{Hash: common.Hash([32]byte{0x22})},
			},
		},
	}

	block3 := &seleneCommon.Block{
		Number: 1000000,
		Hash:   common.HexToHash("0x8e38b4dbf6b11fcc3b9dee84fb7986e29ca0a02cecd8977c161ff7333329681e"),
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{
				common.HexToHash("0xe9e91f1ee4b56c0df2e9f06c2b8c27c6076195a88a7b8537ba8313d80e6f124e"),
				common.HexToHash("0xea1093d492a1dcb1bef708f771a99a96ff05dcab81ca76c31940300177fcf49f"),
			},
			Full: []seleneCommon.Transaction{
				{
					Hash:         common.HexToHash("0xe9e91f1ee4b56c0df2e9f06c2b8c27c6076195a88a7b8537ba8313d80e6f124e"),
					GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
					Gas:          hexutil.Uint64(60),
					MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
				},
				{
					Hash:         common.HexToHash("0xea1093d492a1dcb1bef708f771a99a96ff05dcab81ca76c31940300177fcf49f"),
					GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
					Gas:          hexutil.Uint64(60),
					MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
				},
			},
		},
	}

	// Push blocks through channel
	go func() {
		blockChan <- block1
		blockChan <- block2
		blockChan <- block3
		close(blockChan)
	}()

	// Allow goroutine to process the blocks
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for len(state.blocks) < 3 {
			// wait for blocks to be processed
		}
	}()

	wg.Wait()

	// Simulate finalized block
	finalizedBlock := &seleneCommon.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
			Full: []seleneCommon.Transaction{
				{
					Hash:         common.Hash([32]byte{0x21}),
					GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
					Gas:          hexutil.Uint64(5),
					MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
				},
				{Hash: common.Hash([32]byte{0x22})},
			},
		},
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
			// wait for finalized block to be processed
		}
	}()
	wg.Wait()

	return state
}

func CreateNewExecutionClient() *ExecutionClient {
	rpc := "https://eth-mainnet.g.alchemy.com/v2/j28GcevSYukh-GvSeBOYcwHOfIggF1Gt"
	state := CreateNewState()
	var executionClient *ExecutionClient
	executionClient, _ = executionClient.New(rpc, state)

	return executionClient
}

func TestCheckRpc(t *testing.T) {
	executionClient := CreateNewExecutionClient()
	chainId := uint64(1)
	err := executionClient.CheckRpc(chainId)

	assert.NoError(t, err, "Error Found")

	chainId = uint64(2)
	err = executionClient.CheckRpc(chainId)

	assert.Equal(t, NewIncorrectRpcNetworkError(), err, "Error didn't match")
}

// Both GetAccount() and GetTransactionReceipt() depend on state

// func TestGetAccount(t *testing.T) {
// 	executionClient := CreateNewExecutionClient()
// 	addressBytes, _ := utils.Hex_str_to_bytes("0x457a22804cf255ee8c1b7628601c5682b3d70c71")

// 	address := seleneCommon.Address{Addr: [20]byte(addressBytes)}

// 	slots := []common.Hash{}
// 	tag := seleneCommon.BlockTag{
// 		Finalized: true,
// 	}
// 	print("Check\n")
// 	account, err := executionClient.GetAccount(&address, &slots, tag)

// 	assert.NoError(t, err, "Error found")
// 	assert.Equal(t, Account{}, account, "Account didn't match")
// }

func TestToBlockNumberArg(t *testing.T) {
	blockNumber := uint64(5050)
	assert.Equal(t, "0x13ba", toBlockNumArg(blockNumber), "Block Number didn't match")

	blockNumber = uint64(0)
	assert.Equal(t, "latest", toBlockNumArg(blockNumber), "Block Number didn't match")
}

//* Not possible to test as we would have to send actual txn on mainnet
// func TestSendRawTransaction(t *testing.T) {
// 	executionClient := CreateNewExecutionClient()
// 	transaction := common.Hex2Bytes("02f8720113840a436fe4850749a01900825208942ce3384fcaea81a0f10b2599ffb2f0603e6169f1878e1bc9bf04000080c080a097f6540a48025bd28dd3c43f33aa0269a29b40d852396fab1ab7c2f95a3930e7a03f69a6bca9ef4be6ce60735e76133670617286e15e18af96b7e5e0afcdc240c6")
// 	fmt.Printf("Transaction Bytes: %v ", transaction)
// 	hash, err := executionClient.SendRawTransaction(transaction)

// 	assert.NoError(t, err, "Found Error")
// 	assert.Equal(t, common.Hash{}, hash, "Transaction Hash didn't match")
// }

func TestGetBlock(t *testing.T) {
	executionClient := CreateNewExecutionClient()
	blockTag := seleneCommon.BlockTag{
		Number: 1,
	}

	block, err := executionClient.GetBlock(blockTag, false)
	expected := seleneCommon.Block{
		Number: 1,
		Hash:   [32]byte{0x1},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x11}, {0x12}},
		},
	}

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expected, block, "Value didn't match expected")

	blockTag = seleneCommon.BlockTag{
		Finalized: true,
	}

	block, err = executionClient.GetBlock(blockTag, false)
	expected = seleneCommon.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
		},
	}

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expected, block, "Value didn't match expected")

	block, err = executionClient.GetBlock(blockTag, true)
	expected = seleneCommon.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
			Full: []seleneCommon.Transaction{
				{
					Hash:         common.Hash([32]byte{0x21}),
					GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
					Gas:          hexutil.Uint64(5),
					MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
				},
				{Hash: common.Hash([32]byte{0x22})},
			},
		},
	}

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expected, block, "Value didn't match expected")
}

func TestExecutionGetBlockByHash(t *testing.T) {
	executionClient := CreateNewExecutionClient()
	hash := common.Hash([32]byte{0x2})

	block, err := executionClient.GetBlockByHash(hash, false)
	expected := seleneCommon.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
		},
	}

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expected, block, "Value didn't match expected")

	block, err = executionClient.GetBlockByHash(hash, true)
	expected = seleneCommon.Block{
		Number: 2,
		Hash:   [32]byte{0x2},
		Transactions: seleneCommon.Transactions{
			Hashes: [][32]byte{{0x21}, {0x22}},
			Full: []seleneCommon.Transaction{
				{
					Hash:         common.Hash([32]byte{0x21}),
					GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
					Gas:          hexutil.Uint64(5),
					MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
				},
				{Hash: common.Hash([32]byte{0x22})},
			},
		},
	}

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expected, block, "Value didn't match expected")
}

func TestGetTransactionByBlockHashAndIndex(t *testing.T) {
	executionClient := CreateNewExecutionClient()
	hash := common.Hash([32]byte{0x2})

	txn, err := executionClient.GetTransactionByBlockHashAndIndex(hash, 0)

	expected := seleneCommon.Transaction{
		Hash:         common.Hash([32]byte{0x21}),
		GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
		Gas:          hexutil.Uint64(5),
		MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
	}

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expected, txn, "Value didn't match expected")
}

// func TestExecutionGetTransactionReceipt(t *testing.T) {
// 	executionClient := CreateNewExecutionClient()
// 	txHash := common.HexToHash("0xea1093d492a1dcb1bef708f771a99a96ff05dcab81ca76c31940300177fcf49f")
// 	txnReceipt, err := executionClient.GetTransactionReceipt(txHash)
// 	expected := types.Receipt{}

// 	assert.NoError(t, err, "Found Error")
// 	assert.Equal(t, expected, txnReceipt, "Receipt didn't match")
// }

func TestExecutionGetTransaction(t *testing.T) {
	executionClient := CreateNewExecutionClient()
	txHash := common.Hash([32]byte{0x21})

	txn, err := executionClient.GetTransaction(txHash)
	expected := seleneCommon.Transaction{
		Hash:         common.Hash([32]byte{0x21}),
		GasPrice:     hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
		Gas:          hexutil.Uint64(5),
		MaxFeePerGas: hexutil.Big(*hexutil.MustDecodeBig("0x12345")),
	}

	assert.NoError(t, err, "Found Error")
	assert.Equal(t, expected, txn, "Txn didn't match")
}

func TestExecutionGetLogs(t *testing.T) {
	executionClient := CreateNewExecutionClient()

	// Test case 1: Basic log retrieval
	filter := ethereum.FilterQuery{
		FromBlock: hexutil.MustDecodeBig("0x14057d0"),
		ToBlock:   hexutil.MustDecodeBig("0x1405814"),
		Addresses: []common.Address{
			common.HexToAddress("0x6F9116ea572a207e4267f92b1D7b6F9b23536b07"),
		},
	}

	_, err := executionClient.GetLogs(filter)
	assert.NoError(t, err, "GetLogs should not return error for valid filter")
}

func TestExecutionGetFilterChanges(t *testing.T) {
	executionClient := CreateNewExecutionClient()
	// Test case 1: Basic filter creation
	filter := ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(2),
		Addresses: []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
	}

	filterID, _ := executionClient.GetNewFilter(filter)

	_, err := executionClient.GetFilterChanges(&filterID)
	assert.NoError(t, err, "GetFilterChanges should not return error for valid filter ID")
}

func TestExecutionUninstallFilter(t *testing.T) {
	executionClient := CreateNewExecutionClient()

	// Test case 1: Basic filter creation
	filter := ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(2),
		Addresses: []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
	}
	filterID, _ := executionClient.GetNewFilter(filter)

	result, err := executionClient.UninstallFilter(&filterID)
	assert.NoError(t, err, "UninstallFilter should not return error for valid filter ID")
	assert.True(t, result, "UninstallFilter should return true on success")

	// Test case 2: Invalid filter ID
	invalidFilterID := uint256.NewInt(999)
	result, err = executionClient.UninstallFilter(invalidFilterID)
	assert.NoError(t, err, "UninstallFilter should not return error for invalid filter ID")
	assert.False(t, result, "UninstallFilter should return false for invalid filter ID")
}

func TestExecutionGetNewFilter(t *testing.T) {
	executionClient := CreateNewExecutionClient()

	// Test case 1: Basic filter creation
	filter := ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(2),
		Addresses: []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
	}

	filterID, err := executionClient.GetNewFilter(filter)
	assert.NoError(t, err, "GetNewFilter should not return error for valid filter")
	assert.NotEqual(t, uint256.Int{}, filterID, "FilterID should not be empty")

	// Test case 2: Null blocks defaults to latest
	filterNullBlocks := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
	}

	filterID, err = executionClient.GetNewFilter(filterNullBlocks)
	assert.NoError(t, err, "GetNewFilter should handle null blocks")
	assert.NotEqual(t, uint256.Int{}, filterID, "FilterID should not be empty")
}

func TestExecutionGetNewBlockFilter(t *testing.T) {
	executionClient := CreateNewExecutionClient()

	filterID, err := executionClient.GetNewBlockFilter()
	assert.NoError(t, err, "GetNewBlockFilter should not return error")
	assert.NotEqual(t, uint256.Int{}, filterID, "FilterID should not be empty")
}

func TestExecutionGetNewPendingTransactionFilter(t *testing.T) {
	executionClient := CreateNewExecutionClient()

	filterID, err := executionClient.GetNewPendingTransactionFilter()
	assert.NoError(t, err, "GetNewPendingTransactionFilter should not return error")
	assert.NotEqual(t, uint256.Int{}, filterID, "FilterID should not be empty")
}

func TestCalculateReceiptRoot(t *testing.T) {
	// Test case 1: Empty receipts
	_, err := CalculateReceiptRoot([][]byte{})
	assert.Error(t, err, "CalculateReceiptRoot should return error for empty receipts")
	assert.Contains(t, err.Error(), "no receipts to calculate root")

	// Test case 2: Single receipt
	receipt1 := []byte{1, 2, 3}
	root, err := CalculateReceiptRoot([][]byte{receipt1})
	assert.NoError(t, err, "CalculateReceiptRoot should not return error for single receipt")
	assert.NotEqual(t, common.Hash{}, root, "Root should not be empty")

	// Test case 3: Multiple receipts
	receipt2 := []byte{4, 5, 6}
	root, err = CalculateReceiptRoot([][]byte{receipt1, receipt2})
	assert.NoError(t, err, "CalculateReceiptRoot should not return error for multiple receipts")
	assert.NotEqual(t, common.Hash{}, root, "Root should not be empty")
}

func TestVerifyLogs(t *testing.T) {
	executionClient := CreateNewExecutionClient()
	logData := `{
	  "address": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
	  "blockHash": "0xa495e3d87cf01430e9838341e1d36929fc9c8dfdb89133b372cae4019b5a534a",
	  "blockNumber": "0xC166C3",
	  "data": "0x0000000000000000000000000000000000000000000000003a395039b420f1bd",
	  "logIndex": "0x0",
	  "removed": false,
	  "topics": [
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"0x00000000000000000000000043cc25b1fb6435d8d893fcf308de5c300a568be2",
		"0x00000000000000000000000044d34985826578e5ba24ec78c93be968549bb918"],
	  "transactionHash": "0xac80bab0940f061e184b0dda380d994e6fc14ab5d0c6f689035631c81bfe220b",
	  "transactionIndex": "0x0"
	}`

	var log types.Log
	err := json.Unmarshal([]byte(logData), &log)
	assert.NoError(t, err, "VerifyLogs should not give an error")

	// fmt.Printf("Expected Log: %v\n", log)
	err = executionClient.verifyLogs([]*types.Log{&log})
	assert.NoError(t, err, "VerifyLogs should not give an error")
}

func TestContains(t *testing.T) {
	// Create test receipts
	receipt1 := types.Receipt{
		TxHash: common.HexToHash("0x1"),
	}
	receipt2 := types.Receipt{
		TxHash: common.HexToHash("0x2"),
	}
	receipts := []types.Receipt{receipt1}

	// Test case 1: Receipt exists
	assert.True(t, contains(receipts, receipt1), "Contains should return true for existing receipt")

	// Test case 2: Receipt doesn't exist
	assert.False(t, contains(receipts, receipt2), "Contains should return false for non-existing receipt")
}

// TestRlpHash tests the rlpHash function with various inputs
func TestRlpHash(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "Simple string",
			input:   "test",
			wantErr: false,
		},
		{
			name:    "Empty slice",
			input:   []byte{},
			wantErr: false,
		},
		{
			name: "Complex struct",
			input: struct {
				A string
				B uint64
			}{"test", 123},
			wantErr: false,
		},
		{
			name:    "Invalid input",
			input:   make(chan int), // channels cannot be RLP encoded
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := rlpHash(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("rlpHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify the hash is not empty
				if hash == (common.Hash{}) {
					t.Error("rlpHash() returned empty hash")
				}

				// Verify deterministic behavior
				hash2, _ := rlpHash(tt.input)
				if hash != hash2 {
					t.Error("rlpHash() is not deterministic")
				}
			}
		})
	}
}

// TestCalculateMerkleRoot tests the calculateMerkleRoot function
func TestCalculateMerkleRoot(t *testing.T) {
	tests := []struct {
		name   string
		hashes []common.Hash
		want   common.Hash
	}{
		{
			name: "Two hashes",
			hashes: []common.Hash{
				common.HexToHash("0x1234"),
				common.HexToHash("0x5678"),
			},
			want: crypto.Keccak256Hash(
				append(
					common.HexToHash("0x1234").Bytes(),
					common.HexToHash("0x5678").Bytes()...,
				),
			),
		},
		{
			name: "Odd number of hashes",
			hashes: []common.Hash{
				common.HexToHash("0x1234"),
				common.HexToHash("0x5678"),
				common.HexToHash("0x9abc"),
			},
			want: func() common.Hash {
				// Duplicate last hash since odd number
				hashes := []common.Hash{
					common.HexToHash("0x1234"),
					common.HexToHash("0x5678"),
					common.HexToHash("0x9abc"),
					common.HexToHash("0x9abc"),
				}
				// Calculate first level
				h1 := crypto.Keccak256Hash(append(hashes[0].Bytes(), hashes[1].Bytes()...))
				h2 := crypto.Keccak256Hash(append(hashes[2].Bytes(), hashes[3].Bytes()...))
				// Calculate root
				return crypto.Keccak256Hash(append(h1.Bytes(), h2.Bytes()...))
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateMerkleRoot(tt.hashes)
			if got != tt.want {
				t.Errorf("calculateMerkleRoot() = %v, want %v", got, tt.want)
			}

			// Verify deterministic behavior
			got2 := calculateMerkleRoot(tt.hashes)
			if got != got2 {
				t.Error("calculateMerkleRoot() is not deterministic")
			}
		})
	}
}

// TestEncodeReceipt tests the encodeReceipt function
func TestEncodeReceipt(t *testing.T) {
	bloom := types.Bloom{}
	logs := []*types.Log{}

	tests := []struct {
		name    string
		receipt *types.Receipt
		wantErr bool
	}{
		{
			name: "Legacy receipt (type 0)",
			receipt: &types.Receipt{
				Type:              0,
				Status:            1,
				CumulativeGasUsed: 21000,
				Bloom:             bloom,
				Logs:              logs,
			},
			wantErr: false,
		},
		{
			name: "EIP-2930 receipt (type 1)",
			receipt: &types.Receipt{
				Type:              1,
				Status:            1,
				CumulativeGasUsed: 21000,
				Bloom:             bloom,
				Logs:              logs,
			},
			wantErr: false,
		},
		{
			name: "EIP-1559 receipt (type 2)",
			receipt: &types.Receipt{
				Type:              2,
				Status:            1,
				CumulativeGasUsed: 21000,
				Bloom:             bloom,
				Logs:              logs,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeReceipt(tt.receipt)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeReceipt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the encoding starts with the correct type byte for non-legacy transactions
				if tt.receipt.Type > 0 {
					if got[0] != tt.receipt.Type {
						t.Errorf("encodeReceipt() first byte = %d, want %d", got[0], tt.receipt.Type)
					}
				}

				// Verify the encoding is deterministic
				got2, _ := encodeReceipt(tt.receipt)
				if !bytes.Equal(got, got2) {
					t.Error("encodeReceipt() is not deterministic")
				}

				// Verify the encoded data can be decoded back
				var decoded []interface{}
				var decodedData []byte
				if tt.receipt.Type > 0 {
					decodedData = got[1:] // Skip type byte
				} else {
					decodedData = got
				}

				err = rlp.DecodeBytes(decodedData, &decoded)
				if err != nil {
					t.Errorf("Failed to decode receipt: %v", err)
				}

				// Verify decoded fields match original
				if len(decoded) != 4 {
					t.Errorf("Decoded receipt has wrong number of fields: got %d, want 4", len(decoded))
				}
			}
		})
	}
}
