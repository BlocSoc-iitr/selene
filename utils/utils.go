package utils

import (
	"encoding/hex"
	"strconv"
	"strings"

	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/beacon/merkle"
	"github.com/ethereum/go-ethereum/common"

	consensus_core "github.com/BlocSoc-iitr/selene/consensus/consensus_core"

	beacon "github.com/ethereum/go-ethereum/beacon/types"
	kbls "github.com/kilic/bls12-381"
	bls "github.com/protolambda/bls12-381-util"
)
var domain = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")


// if we need to export the functions , just make their first letter capitalised
func Hex_str_to_bytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	bytesArray, err := hex.DecodeString(s)

	if err != nil {
		return nil, err
	}

	return bytesArray, nil
}

func Hex_str_to_Bytes32(s string) (consensus_core.Bytes32, error) {
	bytesArray, err := Hex_str_to_bytes(s)
	if err != nil {
		return consensus_core.Bytes32{}, err
	}

	var bytes32 consensus_core.Bytes32
	copy(bytes32[:], bytesArray)
	return bytes32, nil
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

func Str_to_uint64(s string) (uint64, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func CalcSyncPeriod(slot uint64) uint64 {
	epoch := slot / 32
	return epoch / 256
}

// isAggregateValid checks if the provided signature is valid for the given message and public keys.
//NO NEED OF THIS FUNCTION HERE AS IT IS NOW DIRECTLY IMPORTED FROM GETH IN CONSENSUS

func IsProofValid(
	attestedHeader *consensus_core.Header,
	leafObject common.Hash, // Single byte slice of the leaf object
	branch []consensus_core.Bytes32, // Slice of byte slices for the branch
	depth, index int, // Depth of the Merkle proof and index of the leaf
) bool {
	var branchInMerkle merkle.Values
	for _, node := range branch {
		branchInMerkle = append(branchInMerkle, merkle.Value(node))
	}
	fmt.Print(attestedHeader.StateRoot)
	err := merkle.VerifyProof(common.Hash(attestedHeader.StateRoot), uint64(index), branchInMerkle, merkle.Value(leafObject))
	if err != nil {
		log.Println("Error in verifying Merkle proof:", err)
		return false
	}
	return true
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

// computed not as helios but as per go-ethereum/beacon/consensus.go
func ComputeForkDataRoot(currentVersion [4]byte, genesisValidatorRoot consensus_core.Bytes32) consensus_core.Bytes32 {
	var (
		hasher        = sha256.New()
		forkVersion32 merkle.Value
		forkDataRoot  merkle.Value
	)
	copy(forkVersion32[:], currentVersion[:])
	hasher.Write(forkVersion32[:])
	hasher.Write(genesisValidatorRoot[:])
	hasher.Sum(forkDataRoot[:0])

	return consensus_core.Bytes32(forkDataRoot)
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

// implemented as per go-ethereum/beacon/types/config.go
func ComputeSigningRoot(header *beacon.Header, domain consensus_core.Bytes32) consensus_core.Bytes32 {
	var (
		signingRoot common.Hash
		headerHash  = header.Hash()
		hasher      = sha256.New()
	)
	hasher.Write(headerHash[:])
	hasher.Write(domain[:])
	hasher.Sum(signingRoot[:0])

	return consensus_core.Bytes32(signingRoot.Bytes())
}

// implemented as per go-ethereum/beacon/types/config.go
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

func FastAggregateVerify(pubkeys []*bls.Pubkey, message []byte, signature *bls.Signature) bool {
	n := uint64(len(pubkeys))
	if n == 0 {
		return false
	}

	g1 := kbls.NewG1()
	// Procedure:
	// 1. aggregate = pubkey_to_point(PK_1)
	// copy the first pubkey
	aggregate := *(*kbls.PointG1)(pubkeys[0])
	// check identity pubkey
	if (*kbls.G1)(nil).IsZero(&aggregate) {
		fmt.Println("Identity pubkey")
		return false
	}
	// 2. for i in 2, ..., n:
	for i := uint64(1); i < n; i++ {
		// 3. next = pubkey_to_point(PK_i)
		next := (*kbls.PointG1)(pubkeys[i])
		// check identity pubkey
		if (*kbls.G1)(nil).IsZero(next) {
			fmt.Println("Identity pubkey")
			return false
		}
		// 4. aggregate = aggregate + next
		g1.Add(&aggregate, &aggregate, next)
	}
	// 5. PK = point_to_pubkey(aggregate)
	PK := (*bls.Pubkey)(&aggregate)
	return coreVerify(PK, message, signature)
}

func coreVerify(pk *bls.Pubkey, message []byte, signature *bls.Signature) bool {
	// 1. R = signature_to_point(signature)
	R := (*kbls.PointG2)(signature)
	// 2. If R is INVALID, return INVALID
	// 3. If signature_subgroup_check(R) is INVALID, return INVALID
	// 4. If KeyValidate(PK) is INVALID, return INVALID
	// steps 2-4 are part of bytes -> *Signature deserialization
	if (*kbls.G2)(nil).IsZero(R) {
		// KeyValidate is assumed through deserialization of Pubkey and Signature,
		// but the identity pubkey/signature case is not part of that, thus verify here.
		fmt.Print("Identity signature")
		return false
	}

	// 5. xP = pubkey_to_point(PK)
	xP := (*kbls.PointG1)(pk)
	// 6. Q = hash_to_point(message)
	Q, err := kbls.NewG2().HashToCurve(message, domain)
	if err != nil {
		// e.g. when the domain is too long. Maybe change to panic if never due to a usage error?
		fmt.Print("Failed to hash message to point")
		return false
	}

	// 7. C1 = pairing(Q, xP)
	eng := kbls.NewEngine()
	eng.AddPair(xP, Q)
	// 8. C2 = pairing(R, P)
	P := &kbls.G1One
	eng.AddPairInv(P, R)
	// 9. If C1 != C2, return INVALID
	return !eng.Check()
}
