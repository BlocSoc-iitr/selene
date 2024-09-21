package consensus

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/BlocSoc-iitr/selene/config"
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/BlocSoc-iitr/selene/utils/bls"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/ethereum/go-ethereum/crypto"
)

// TreeHashRoot computes the Merkle root from the provided leaves in a flat []byte slice.
func TreeHashRoot(data []byte) ([]byte, error) {
	// Convert the input data into a slice of leaves
	leaves, err := bytesToLeaves(data)
	if err != nil {
		return nil, err
	}

	nodes := leaves // Start with the leaf nodes

	for len(nodes) > 1 {
		var newLevel [][]byte

		// Pair nodes and hash them
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				// Hash pair of nodes
				nodeHash := crypto.Keccak256(append(nodes[i], nodes[i+1]...))
				newLevel = append(newLevel, nodeHash)
			} else {
				// Handle odd number of nodes (carry last node up)
				newLevel = append(newLevel, nodes[i])
			}
		}

		nodes = newLevel
	}

	// Return the root hash
	return nodes[0], nil
}

func bytesToLeaves(data []byte) ([][]byte, error) {
	var leaves [][]byte
	if err := json.Unmarshal(data, &leaves); err != nil {
		return nil, err
	}
	return leaves, nil
}

func CalcSyncPeriod(slot uint64) uint64 {
	epoch := slot / 32
	return epoch / 256
}

// isAggregateValid checks if the provided signature is valid for the given message and public keys.
func isAggregateValid(sigBytes consensus_core.SignatureBytes, msg [32]byte, pks []*bls.G2Point) bool {
	var sigInBytes [96]byte
	copy(sigInBytes[:], sigBytes[:])
	// Deserialize the signature from bytes
	var sig bls12381.G1Affine
	if err := sig.Unmarshal(sigInBytes[:]); err != nil {
		return false
	}

	// Map the message to a point on the curve
	msgPoint := bls.MapToCurve(msg)

	// Aggregate the public keys
	aggPubKey := bls.AggregatePublicKeys(pks)

	// Prepare the pairing check inputs
	P := [2]bls12381.G1Affine{*msgPoint, sig}
	Q := [2]bls12381.G2Affine{*aggPubKey.G2Affine, *bls.GetG2Generator()}

	// Perform the pairing check
	ok, err := bls12381.PairingCheck(P[:], Q[:])
	if err != nil {
		return false
	}
	return ok
}

func isProofValid(
	attestedHeader *consensus_core.Header,
	leafObject []byte, // Single byte slice of the leaf object
	branch [][]byte, // Slice of byte slices for the branch
	depth, index int, // Depth of the Merkle proof and index of the leaf
) bool {
	// If the branch length is not equal to the depth, return false
	if len(branch) != depth {
		return false
	}

	// Compute the root hash of the leaf object
	derivedRoot, err := TreeHashRoot(leafObject)
	if err != nil {
		return false
	}

	// Iterate through the branch and compute the Merkle root
	for i, node := range branch {
		hasher := sha256.New()

		// Check if index / 2^i is odd or even
		if (index/int(math.Pow(2, float64(i))))%2 != 0 {
			// If odd, hash(node || derived_root)
			hasher.Write(node)
			hasher.Write(derivedRoot[:])
		} else {
			// If even, hash(derived_root || node)
			hasher.Write(derivedRoot[:])
			hasher.Write(node)
		}

		// Update the derived root
		derivedRootNew := sha256.Sum256(hasher.Sum(nil))
		derivedRoot = derivedRootNew[:]
	}

	// Compare the final derived root with the attested header's state root
	return bytes.Equal(derivedRoot[:], attestedHeader.StateRoot[:])
}

func CalculateForkVersion(forks *config.Forks, slot uint64) [4]byte {
	epoch := slot / 32

	switch {
	case epoch >= forks.Deneb.Epoch:
		return [4]byte(forks.Deneb.ForkVersion)
	case epoch >= forks.Capella.Epoch:
		return [4]byte(forks.Capella.ForkVersion)
	case epoch >= forks.Bellatrix.Epoch:
		return [4]byte(forks.Bellatrix.ForkVersion)
	case epoch >= forks.Altair.Epoch:
		return [4]byte(forks.Altair.ForkVersion)
	default:
		return [4]byte(forks.Genesis.ForkVersion)
	}
}

func ComputeForkDataRoot(currentVersion [4]byte, genesisValidatorRoot consensus_core.Bytes32) consensus_core.Bytes32 {
	forkData := ForkData{
		CurrentVersion:       currentVersion,
		GenesisValidatorRoot: genesisValidatorRoot,
	}

	hash, err := TreeHashRoot(forkData.ToBytes())
	if err != nil {
		return consensus_core.Bytes32{}
	}
	return consensus_core.Bytes32(hash)
}

// GetParticipatingKeys retrieves the participating public keys from the committee based on the bitfield represented as a byte array.
func GetParticipatingKeys(committee *consensus_core.SyncCommittee, bitfield [64]byte) ([]consensus_core.BLSPubKey, error) {
	var pks []consensus_core.BLSPubKey
	numBits := len(bitfield) * 8 // Total number of bits

	if len(committee.Pubkeys) > numBits {
		return nil, fmt.Errorf("bitfield is too short for the number of public keys")
	}

	for i := 0; i < len(bitfield); i++ {
		byteVal := bitfield[i]
		for bit := 0; bit < 8; bit++ {
			if (byteVal & (1 << bit)) != 0 {
				index := i*8 + bit
				if index >= len(committee.Pubkeys) {
					break
				}
				pks = append(pks, committee.Pubkeys[index])
			}
		}
	}

	return pks, nil
}

func ComputeSigningRoot(objectRoot, domain consensus_core.Bytes32) consensus_core.Bytes32 {
	signingData := SigningData{
		ObjectRoot: objectRoot,
		Domain:     domain,
	}
	hash, err := TreeHashRoot(signingData.ToBytes())
	if err != nil {
		return consensus_core.Bytes32{}
	}
	return consensus_core.Bytes32(hash)
}

func ComputeDomain(domainType [4]byte, forkDataRoot consensus_core.Bytes32) consensus_core.Bytes32 {
	data := append(domainType[:], forkDataRoot[:28]...)
	return sha256.Sum256(data)
}

type SigningData struct {
	ObjectRoot consensus_core.Bytes32
	Domain     consensus_core.Bytes32
}

type ForkData struct {
	CurrentVersion       [4]byte
	GenesisValidatorRoot consensus_core.Bytes32
}

func (fd *ForkData) ToBytes() []byte {
	data, err := json.Marshal(fd)
	if err != nil {
		log.Println("Error marshaling ForkData:", err)
		return nil // Or return an empty slice, based on your preference
	}
	return data
}

func (sd *SigningData) ToBytes() []byte {
	data, err := json.Marshal(sd)
	if err != nil {
		log.Println("Error marshaling SigningData:", err)
		return nil // Or return an empty slice, based on your preference
	}
	return data
}

func bytes32ToNode(bytes consensus_core.Bytes32) []byte {
	return []byte(bytes[:])
}

// branchToNodes converts a slice of Bytes32 to a slice of Node
func branchToNodes(branch []consensus_core.Bytes32) ([][]byte, error) {
	nodes := make([][]byte, len(branch))
	for i, b32 := range branch {
		nodes[i] = bytes32ToNode(b32)
	}
	return nodes, nil
}
