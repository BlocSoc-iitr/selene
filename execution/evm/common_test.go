package evm

import (
	"math/big"
	"testing"

	// "github.com/stretchr/testify"
	Common "github.com/ethereum/go-ethereum/common"
	"github.com/magiconair/properties/assert"
)

func TestU256Initialization(t *testing.T) {
	var zero U256 = big.NewInt(0)
	if zero.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected zero U256, got %v", zero)
	}
	var positive U256 = big.NewInt(17)
	if positive.Cmp(big.NewInt(17)) != 0 {
		t.Errorf("Expected U256 to be 17, got %v", positive)
	}
	largeNum := new(big.Int)
	largeNum.SetString("123456789012345678901234567890", 10)
	var uLarge U256 = largeNum
	if uLarge.Cmp(largeNum) != 0 {
		t.Errorf("Expected U256 to be large number %v, got %v", largeNum, uLarge)
	}
	if positive.Cmp(largeNum) == 0 {
		t.Errorf("Expected U256 %v to be different from large number %v", positive, largeNum)
	}

}
func TestU256Equality(t *testing.T) {
	var a U256 = big.NewInt(17)
	var b U256 = big.NewInt(17)
	var c U256 = big.NewInt(100)
	if a.Cmp(b) != 0 {
		t.Errorf("Expected U256 instances to be equal, got a: %v, b: %v", a, b)
	}
	if a.Cmp(c) == 0 {
		t.Errorf("Expected U256 instances to be different, got a: %v, c: %v", a, c)
	}
}
func TestB256Initialization(t *testing.T) {
	var hash B256 = Common.BytesToHash([]byte("selene"))
	expectedHash := Common.BytesToHash([]byte("selene"))
	assert.Equal(t, expectedHash, hash, "Expected hash and expectedHash to be same")

}
func TestAddressInitialization(t *testing.T) {
    // Test Address initialization
    addr := Common.HexToAddress("0x0000000000000000000000000000000000000001")
    if addr.String() != "0x0000000000000000000000000000000000000001" {
        t.Errorf("Expected Address to be 0x0000000000000000000000000000000000000001, got %v", addr)
    }

    // Test Address validity with a valid input
    validAddr := Common.HexToAddress("0x0000000000000000000000000000000000000001")
    if validAddr == (Common.Address{}) {
        t.Errorf("Expected a valid Address, got zero address")
    }

    // Test Address validity with an invalid input
    invalidAddr := Common.HexToAddress("invalid-address")
    if invalidAddr != (Common.Address{}) {
        t.Errorf("Expected zero Address for invalid input, got %v", invalidAddr)
    }
}


func TestBytesInitialization(t *testing.T) {
	// Test Bytes initialization
	data := []byte{1, 2, 3, 4}
	var byteData Bytes = data

	if len(byteData) != 4 {
		t.Errorf("Expected Bytes length to be 4, got %v", len(byteData))
	}

	if byteData[0] != 1 {
		t.Errorf("Expected first byte to be 1, got %v", byteData[0])
	}
}