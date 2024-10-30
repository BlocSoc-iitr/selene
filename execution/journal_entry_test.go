package evm

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HexToAddress converts a hex string to an Address type
func HexToAddress(hexStr string) Address {
	var addr Address
	bytes, err := hex.DecodeString(hexStr[2:]) // Remove '0x' prefix
	if err == nil && len(bytes) == 20 {
		var bytesNew [20]byte
		copy(bytesNew[:], bytes)
		addr = Address{
			Addr: bytesNew,
		}
	}
	return addr
}

// Helper function to create U256 values
func newU256(value int64) *big.Int {
	return big.NewInt(value)
}

func TestNewAccountWarmedEntry(t *testing.T) {
	address := Address{}
	entry := NewAccountWarmedEntry(address)

	assert.Equal(t, AccountWarmedType, entry.Type, "Type should be AccountWarmedType")
	assert.Equal(t, address, entry.Address, "Address should match")
}

func TestNewAccountDestroyedEntry(t *testing.T) {
	address := HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	target := HexToAddress("0x5678567856785678567856785678567856785678")
	hadBalance := newU256(100)
	entry := NewAccountDestroyedEntry(address, target, true, hadBalance)

	assert.Equal(t, AccountDestroyedType, entry.Type, "Type should be AccountDestroyedType")
	assert.Equal(t, address, entry.Address, "Address should match")
	assert.Equal(t, target, entry.Target, "Target should match")
	assert.True(t, entry.WasDestroyed, "WasDestroyed should be true")
	assert.Equal(t, hadBalance, entry.HadBalance, "HadBalance should match")
}

func TestNewAccountTouchedEntry(t *testing.T) {
	address := HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	entry := NewAccountTouchedEntry(address)

	assert.Equal(t, AccountTouchedType, entry.Type, "Type should be AccountTouchedType")
	assert.Equal(t, address, entry.Address, "Address should match")
}

func TestNewBalanceTransferEntry(t *testing.T) {
	from := HexToAddress("0x1111111111111111111111111111111111111111")
	to := HexToAddress("0x2222222222222222222222222222222222222222")
	balance := newU256(50)
	entry := NewBalanceTransferEntry(from, to, balance)

	assert.Equal(t, BalanceTransferType, entry.Type, "Type should be BalanceTransferType")
	assert.Equal(t, from, entry.From, "From address should match")
	assert.Equal(t, to, entry.To, "To address should match")
	assert.Equal(t, balance, entry.Balance, "Balance should match")
}

func TestNewNonceChangeEntry(t *testing.T) {
	address := HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	entry := NewNonceChangeEntry(address)

	assert.Equal(t, NonceChangeType, entry.Type, "Type should be NonceChangeType")
	assert.Equal(t, address, entry.Address, "Address should match")
}

func TestNewAccountCreatedEntry(t *testing.T) {
	address := HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	entry := NewAccountCreatedEntry(address)

	assert.Equal(t, AccountCreatedType, entry.Type, "Type should be AccountCreatedType")
	assert.Equal(t, address, entry.Address, "Address should match")
}

func TestNewStorageChangedEntry(t *testing.T) {
	address := HexToAddress("0x3333333333333333333333333333333333333333")
	key := newU256(1)
	hadValue := newU256(10)
	entry := NewStorageChangedEntry(address, key, hadValue)

	assert.Equal(t, StorageChangedType, entry.Type, "Type should be StorageChangedType")
	assert.Equal(t, address, entry.Address, "Address should match")
	assert.Equal(t, key, entry.Key, "Key should match")
	assert.Equal(t, hadValue, entry.HadValue, "HadValue should match")
}

func TestNewStorageWarmedEntry(t *testing.T) {
	address := HexToAddress("0x3333333333333333333333333333333333333333")
	key := newU256(1)
	entry := NewStorageWarmedEntry(address, key)

	assert.Equal(t, StorageWarmedType, entry.Type, "Type should be StorageWarmedType")
	assert.Equal(t, address, entry.Address, "Address should match")
	assert.Equal(t, key, entry.Key, "Key should match")
}

func TestNewTransientStorageChangeEntry(t *testing.T) {
	address := HexToAddress("0x4444444444444444444444444444444444444444")
	key := newU256(5)
	hadValue := newU256(20)
	entry := NewTransientStorageChangeEntry(address, key, hadValue)

	assert.Equal(t, TransientStorageChangeType, entry.Type, "Type should be TransientStorageChangeType")
	assert.Equal(t, address, entry.Address, "Address should match")
	assert.Equal(t, key, entry.Key, "Key should match")
	assert.Equal(t, hadValue, entry.HadValue, "HadValue should match")
}

func TestNewCodeChangeEntry(t *testing.T) {
	address := HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	entry := NewCodeChangeEntry(address)

	assert.Equal(t, CodeChangeType, entry.Type, "Type should be CodeChangeType")
	assert.Equal(t, address, entry.Address, "Address should match")
}

func TestMarshalJSON(t *testing.T) {
	// Create an example JournalEntry object
	entry := JournalEntry{
		Type:         1,
		Address:      Address{Addr: [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}},
		Target:       Address{Addr: [20]byte{0x56, 0x78, 0x56, 0x78, 0x56, 0x78, 0x56, 0x78}},
		WasDestroyed: true,
		HadBalance:   big.NewInt(100),
		Balance:      big.NewInt(200),
		From:         Address{Addr: [20]byte{0x00, 0x00, 0x00, 0x00}},
		To:           Address{Addr: [20]byte{0x00, 0x00, 0x00, 0x00}},
		Key:          big.NewInt(300),
		HadValue:     big.NewInt(400),
	}

	// Marshal the entry to JSON
	actualJSON, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Define the expected JSON output
	expectedJSON := `{
		"type": 1,
		"address": {"Addr": [18, 52, 86, 120, 144, 171, 205, 239, 0 , 0 , 0 , 0 , 0 ,0 , 0 ,0 ,0, 0, 0, 0]},
		"target": {"Addr": [86, 120, 86, 120, 86, 120, 86, 120, 0 ,0 ,0,0,0,0,0,0,0,0,0,0]},
		"was_destroyed": true,
		"had_balance": 100,
		"balance": 200,
		"from": {"Addr": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
		"to": {"Addr": [0, 0, 0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]},
		"key": 300,
		"had_value": 400
	}`

	// Compare the actual JSON output to the expected JSON
	var expected, actual map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(actualJSON, &actual); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	if len(expected) != len(actual) || !compareJSON(expected, actual) {
		t.Errorf("Expected JSON does not match actual JSON.\nExpected: %s\nActual: %s", expectedJSON, string(actualJSON))
	}
}

// Helper function to compare two JSON objects
func compareJSON(a, b map[string]interface{}) bool {
	return reflect.DeepEqual(a, b)
}
