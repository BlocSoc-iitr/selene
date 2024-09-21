package proof

import (
	"fmt"

	merkletree "github.com/wealdtech/go-merkletree"
)

func validateMerkleProof(
	root []byte,
	leaf []byte,
	proof *merkletree.Proof,
) (bool, error) {
	is_valid, err := merkletree.VerifyProof(leaf, proof, root)
	if err != nil {
		fmt.Printf("Error in validating merkle proof: %v", err)
	}
	return is_valid, nil
}
