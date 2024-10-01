package utils

import (
	"encoding/hex"
        "net/url"
	"strings"

	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/BlocSoc-iitr/selene/utils/bls"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/pkg/errors"
	merkletree "github.com/wealdtech/go-merkletree"
)

// if we need to export the functions , just make their first letter capitalised
func Hex_str_to_bytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	bytesArray, err := hex.DecodeString(s)

	if err != nil {
		return nil, err
	}

	return bytesArray, nil
}

func Address_to_hex_string(addr common.Address) string {
	bytesArray := addr.Bytes()
	return fmt.Sprintf("0x%x", hex.EncodeToString(bytesArray))
}

func U64_to_hex_string(val uint64) string {
	return fmt.Sprintf("0x%x", val)
}

func Bytes_deserialize(data []byte) ([]byte, error) {
	var hexString string
	if err := json.Unmarshal(data, &hexString); err != nil {
		return nil, err
	}

	return Hex_str_to_bytes(hexString)
}

func Bytes_serialize(bytes []byte) ([]byte, error) {
	if bytes == nil {
		return json.Marshal(nil)
	}
	hexString := hex.EncodeToString(bytes)
	return json.Marshal(hexString)
}

// TreeHashRoot computes the Merkle root from the provided leaves in a flat []byte slice.
// TreeHashRoot calculates the root hash from the input data.
func TreeHashRoot(data []byte) ([]byte, error) {
	// Convert the input data into a slice of leaves
	leaves, err := bytesToLeaves(data)
	if err != nil {
		return nil, fmt.Errorf("error converting bytes to leaves: %w", err)
	}

	// Create the Merkle tree using the leaves
	tree, errorCreatingMerkleTree := merkletree.New(leaves)
	if errorCreatingMerkleTree != nil {
		return nil, fmt.Errorf("error creating Merkle tree: %w", err)
	}

	// Fetch the root hash of the tree
	root := tree.Root()
	if root == nil {
		return nil, errors.New("failed to calculate the Merkle root: root is nil")
	}

	return root, nil
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
func IsAggregateValid(sigBytes consensus_core.SignatureBytes, msg [32]byte, pks []*bls.G2Point) bool {
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

func IsProofValid(
	attestedHeader *consensus_core.Header,
	leafObject []byte, // Single byte slice of the leaf object
	branch [][]byte, // Slice of byte slices for the branch
	depth, index int, // Depth of the Merkle proof and index of the leaf
) bool {
	// If the branch length is not equal to the depth, return false
	if len(branch) != depth {
		return false
	}

	// Initialize the derived root as the leaf object's hash
	derivedRoot, errFetchingTreeHashRoot := TreeHashRoot(leafObject)

	if errFetchingTreeHashRoot != nil {
		fmt.Printf("Error fetching tree hash root: %v", errFetchingTreeHashRoot)
		return false
	}

	// Iterate over the proof's hashes
	for i, hash := range branch {
		hasher := sha256.New()

		// Use the index to determine how to combine the hashes
		if (index>>i)&1 == 1 {
			// If index is odd, hash node || derivedRoot
			hasher.Write(hash)
			hasher.Write(derivedRoot)
		} else {
			// If index is even, hash derivedRoot || node
			hasher.Write(derivedRoot)
			hasher.Write(hash)
		}
		// Update derivedRoot for the next level
		derivedRoot = hasher.Sum(nil)
	}

	// Compare the final derived root with the attested header's state root
	return bytes.Equal(derivedRoot, attestedHeader.StateRoot[:])
}

func CalculateForkVersion(forks *consensus_core.Forks, slot uint64) [4]byte {
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

func Bytes32ToNode(bytes consensus_core.Bytes32) []byte {
	return []byte(bytes[:])
}

// branchToNodes converts a slice of Bytes32 to a slice of Node
func BranchToNodes(branch []consensus_core.Bytes32) ([][]byte, error) {
	nodes := make([][]byte, len(branch))
	for i, b32 := range branch {
		nodes[i] = Bytes32ToNode(b32)
	}
	return nodes, nil
}
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
