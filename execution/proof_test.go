package execution

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (a *Account) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{a.Nonce, a.Balance, a.StorageHash, a.CodeHash})
}

func (a *Account) DecodeRLP(s *rlp.Stream) error {
	var accountRLP struct {
		Nonce       uint64
		Balance     *big.Int
		StorageHash common.Hash
		CodeHash    common.Hash
	}
	if err := s.Decode(&accountRLP); err != nil {
		return err
	}
	a.Nonce = accountRLP.Nonce
	a.Balance = accountRLP.Balance
	a.StorageHash = accountRLP.StorageHash
	a.CodeHash = accountRLP.CodeHash

	return nil
}

func TestSharedPrefixLength(t *testing.T) {
	t.Run("Match first 5 nibbles", func(t *testing.T) {
		path := []byte{0x12, 0x13, 0x14, 0x6f, 0x6c, 0x64, 0x21}
		pathOffset := 6
		nodePath := []byte{0x6f, 0x6c, 0x63, 0x21}

		sharedLen := sharedPrefixLength(path, pathOffset, nodePath)
		assert.Equal(t, 5, sharedLen)
	})

	t.Run("Match first 7 nibbles with skip", func(t *testing.T) {
		path := []byte{0x12, 0x13, 0x14, 0x6f, 0x6c, 0x64, 0x21}
		pathOffset := 5
		nodePath := []byte{0x14, 0x6f, 0x6c, 0x64, 0x11}

		sharedLen := sharedPrefixLength(path, pathOffset, nodePath)
		assert.Equal(t, 7, sharedLen)
	})

	t.Run("No match", func(t *testing.T) {
		path := []byte{0x12, 0x34, 0x56}
		pathOffset := 0
		nodePath := []byte{0x78, 0x9a, 0xbc}

		sharedLen := sharedPrefixLength(path, pathOffset, nodePath)
		assert.Equal(t, 0, sharedLen)
	})

	t.Run("Complete match", func(t *testing.T) {
		path := []byte{0x12, 0x34, 0x56}
		pathOffset := 0
		nodePath := []byte{0x00, 0x12, 0x34, 0x56}

		sharedLen := sharedPrefixLength(path, pathOffset, nodePath)
		assert.Equal(t, 6, sharedLen)
	})
}

func TestSkipLength(t *testing.T) {
	testCases := []struct {
		name     string
		node     []byte
		expected int
	}{
		{"Empty node", []byte{}, 0},
		{"Nibble 0", []byte{0x00}, 2},
		{"Nibble 1", []byte{0x10}, 1},
		{"Nibble 2", []byte{0x20}, 2},
		{"Nibble 3", []byte{0x30}, 1},
		{"Nibble 4", []byte{0x40}, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := skipLength(tc.node)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetNibble(t *testing.T) {
	path := []byte{0x12, 0x34, 0x56}

	testCases := []struct {
		offset   int
		expected byte
	}{
		{0, 0x1},
		{1, 0x2},
		{2, 0x3},
		{3, 0x4},
		{4, 0x5},
		{5, 0x6},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Offset %d", tc.offset), func(t *testing.T) {
			result := getNibble(path, tc.offset)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsEmptyValue(t *testing.T) {
	t.Run("Empty slot", func(t *testing.T) {
		value := []byte{0x80}
		assert.True(t, isEmptyValue(value))
	})

	t.Run("Empty account", func(t *testing.T) {
		emptyAccount := Account{
			Nonce:       0,
			Balance:     big.NewInt(0),
			StorageHash: [32]byte{0x56, 0xe8, 0x1f, 0x17, 0x1b, 0xcc, 0x55, 0xa6, 0xff, 0x83, 0x45, 0xe6, 0x92, 0xc0, 0xf8, 0x6e, 0x5b, 0x48, 0xe0, 0x1b, 0x99, 0x6c, 0xad, 0xc0, 0x01, 0x62, 0x2f, 0xb5, 0xe3, 0x63, 0xb4, 0x21},
			CodeHash:    [32]byte{0xc5, 0xd2, 0x46, 0x01, 0x86, 0xf7, 0x23, 0x3c, 0x92, 0x7e, 0x7d, 0xb2, 0xdc, 0xc7, 0x03, 0xc0, 0xe5, 0x00, 0xb6, 0x53, 0xca, 0x82, 0x27, 0x3b, 0x7b, 0xfa, 0xd8, 0x04, 0x5d, 0x85, 0xa4, 0x70},
		}
		encodedEmptyAccount, err := rlp.EncodeToBytes(&emptyAccount)
		require.NoError(t, err)
		assert.True(t, isEmptyValue(encodedEmptyAccount))
	})

	t.Run("Non-empty value", func(t *testing.T) {
		value := []byte{0x01, 0x23, 0x45}
		assert.False(t, isEmptyValue(value))
	})
}

func TestEncodeAccount(t *testing.T) {
	proof := &EIP1186AccountProofResponse{
		Nonce:       1,
		Balance:     big.NewInt(1000),
		StorageHash: [32]byte{1, 2, 3},
		CodeHash:    [32]byte{4, 5, 6},
	}

	encoded, err := EncodeAccount(proof)
	require.NoError(t, err)

	var decodedAccount Account
	err = rlp.DecodeBytes(encoded, &decodedAccount)
	require.NoError(t, err)

	// Assert the decoded values match the original
	assert.Equal(t, uint64(1), decodedAccount.Nonce)
	assert.Equal(t, big.NewInt(1000), decodedAccount.Balance)
	assert.Equal(t, proof.StorageHash, decodedAccount.StorageHash)
	assert.Equal(t, proof.CodeHash, decodedAccount.CodeHash)
}

func TestKeccak256(t *testing.T) {
	input := []byte("hello world")
	expected, _ := hex.DecodeString("47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad")

	result := keccak256(input)
	assert.Equal(t, expected, result)
}
