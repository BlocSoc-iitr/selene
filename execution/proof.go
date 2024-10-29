package execution

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

func VerifyProof(proof [][]byte, root []byte, path []byte, value []byte) (bool, error) {
	expectedHash := root
	pathOffset := 0

	for i, node := range proof {
		if !bytes.Equal(expectedHash, keccak256(node)) {
			return false, nil
		}

		var nodeList [][]byte
		if err := rlp.DecodeBytes(node, &nodeList); err != nil {
			fmt.Println("Error decoding node:", err)
			return false, err
		}

		if len(nodeList) == 17 {
			if i == len(proof)-1 {
				// exclusion proof
				nibble := getNibble(path, pathOffset)
				if len(nodeList[nibble]) == 0 && isEmptyValue(value) {
					return true, nil
				}
			} else {
				nibble := getNibble(path, pathOffset)
				expectedHash = nodeList[nibble]
				pathOffset++
			}
		} else if len(nodeList) == 2 {
			if i == len(proof)-1 {
				// exclusion proof
				if !pathsMatch(nodeList[0], skipLength(nodeList[0]), path, pathOffset) && isEmptyValue(value) {
					return true, nil
				}

				// inclusion proof
				if bytes.Equal(nodeList[1], value) {
					return pathsMatch(nodeList[0], skipLength(nodeList[0]), path, pathOffset), nil
				}
			} else {
				nodePath := nodeList[0]
				prefixLength := sharedPrefixLength(path, pathOffset, nodePath)
				if prefixLength < len(nodePath)*2-skipLength(nodePath) {
					// Proof shows a divergent path , but we're not at the leaf yet
					return false, nil
				}
				pathOffset += prefixLength
				expectedHash = nodeList[1]
			}
		} else {
			return false, nil
		}
	}

	return false, nil
}

func pathsMatch(p1 []byte, s1 int, p2 []byte, s2 int) bool {
	len1 := len(p1)*2 - s1
	len2 := len(p2)*2 - s2

	if len1 != len2 {
		return false
	}

	for offset := 0; offset < len1; offset++ {
		n1 := getNibble(p1, s1+offset)
		n2 := getNibble(p2, s2+offset)
		if n1 != n2 {
			return false
		}
	}

	return true
}

// dead code
func GetRestPath(p []byte, s int) string {
	var ret string
	for i := s; i < len(p)*2; i++ {
		n := getNibble(p, i)
		ret += fmt.Sprintf("%01x", n)
	}
	return ret
}

func isEmptyValue(value []byte) bool {
	emptyAccount := Account{
		Nonce:       0,
		Balance:     uint256.NewInt(0).ToBig(),
		StorageHash: [32]byte{0x56, 0xe8, 0x1f, 0x17, 0x1b, 0xcc, 0x55, 0xa6, 0xff, 0x83, 0x45, 0xe6, 0x92, 0xc0, 0xf8, 0x6e, 0x5b, 0x48, 0xe0, 0x1b, 0x99, 0x6c, 0xad, 0xc0, 0x01, 0x62, 0x2f, 0xb5, 0xe3, 0x63, 0xb4, 0x21},
		CodeHash:    [32]byte{0xc5, 0xd2, 0x46, 0x01, 0x86, 0xf7, 0x23, 0x3c, 0x92, 0x7e, 0x7d, 0xb2, 0xdc, 0xc7, 0x03, 0xc0, 0xe5, 0x00, 0xb6, 0x53, 0xca, 0x82, 0x27, 0x3b, 0x7b, 0xfa, 0xd8, 0x04, 0x5d, 0x85, 0xa4, 0x70},
	}

	encodedEmptyAccount, _ := rlp.EncodeToBytes(emptyAccount)

	isEmptySlot := len(value) == 1 && value[0] == 0x80
	isEmptyAccount := bytes.Equal(value, encodedEmptyAccount)

	return isEmptySlot || isEmptyAccount
}

func sharedPrefixLength(path []byte, pathOffset int, nodePath []byte) int {
	skipLength := skipLength(nodePath)

	len1 := min(len(nodePath)*2-skipLength, len(path)*2-pathOffset)
	prefixLen := 0

	for i := 0; i < len1; i++ {
		pathNibble := getNibble(path, i+pathOffset)
		nodePathNibble := getNibble(nodePath, i+skipLength)
		if pathNibble != nodePathNibble {
			break
		}
		prefixLen++
	}

	return prefixLen
}

func skipLength(node []byte) int {
	if len(node) == 0 {
		return 0
	}

	nibble := getNibble(node, 0)
	switch nibble {
	case 0, 2:
		return 2
	case 1, 3:
		return 1
	default:
		return 0
	}
}

func getNibble(path []byte, offset int) byte {
	byteVal := path[offset/2]
	if offset%2 == 0 {
		return byteVal >> 4
	}
	return byteVal & 0xF
}

func keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}
// Updated as EIP1186ProofResponse was updated
func EncodeAccount(proof *EIP1186ProofResponse) ([]byte, error) {
	account := Account{
		Nonce:       uint64(proof.Nonce),
		Balance:     proof.Balance.ToBig(),
		StorageHash: proof.StorageHash,
		CodeHash:    proof.CodeHash,
	}

	return rlp.EncodeToBytes(account)
}

// Make a generic function for it
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}