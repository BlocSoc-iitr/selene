package consensus_core

import (
	"bytes"

	"github.com/ugorji/go/codec"
)

type Bytes32 [32]byte
type BLSPubKey [48]byte
type Address [20]byte
type LogsBloom [256]byte
type SignatureBytes [96]byte
type Transaction = [1073741824]byte //1073741824

type BeaconBlock struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    Bytes32
	StateRoot     Bytes32
	Body          BeaconBlockBody
}

type Eth1Data struct {
	DepositRoot  Bytes32
	DepositCount uint64
	BlockHash    Bytes32
}

type ProposerSlashing struct {
	SignedHeader1 SignedBeaconBlockHeader
	SignedHeader2 SignedBeaconBlockHeader
}

type SignedBeaconBlockHeader struct {
	Message   BeaconBlockHeader
	Signature SignatureBytes
}

type BeaconBlockHeader struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    Bytes32
	StateRoot     Bytes32
	BodyRoot      Bytes32
}

type AttesterSlashing struct {
	Attestation1 IndexedAttestation
	Attestation2 IndexedAttestation
}

type IndexedAttestation struct {
	AttestingIndices [2048]uint64
	Data             AttestationData
	Signature        SignatureBytes
}

type AttestationData struct {
	Slot            uint64
	Index           uint64
	BeaconBlockRoot Bytes32
	Source          Checkpoint
	Target          Checkpoint
}

type Checkpoint struct {
	Epoch uint64
	Root  Bytes32
}

type Bitlist []bool

type Attestation struct {
	AggregationBits Bitlist `ssz-max:"2048"`
	Data            AttestationData
	Signature       SignatureBytes
}

type Deposit struct {
	Proof [33]Bytes32 // fixed size array
	Data  DepositData
}

type DepositData struct {
	Pubkey                [48]byte
	WithdrawalCredentials Bytes32
	Amount                uint64
	Signature             SignatureBytes
}

type SignedVoluntaryExit struct {
	Message   VoluntaryExit
	Signature SignatureBytes
}

type VoluntaryExit struct {
	Epoch          uint64
	ValidatorIndex uint64
}

type SyncAggregate struct {
	SyncCommitteeBits      [64]byte
	SyncCommitteeSignature SignatureBytes
}

type Withdrawal struct {
	Index          uint64
	ValidatorIndex uint64
	Address        Address
	Amount         uint64
}

type ExecutionPayload struct {
	ParentHash    Bytes32
	FeeRecipient  Address
	StateRoot     Bytes32
	ReceiptsRoot  Bytes32
	LogsBloom     LogsBloom
	PrevRandao    Bytes32
	BlockNumber   uint64
	GasLimit      uint64
	GasUsed       uint64
	Timestamp     uint64
	ExtraData     [32]byte
	BaseFeePerGas uint64
	BlockHash     Bytes32
	Transactions  []Transaction `ssz-max:"1048576"`
	Withdrawals   []Withdrawal  `ssz-max:"16"`
	BlobGasUsed   *uint64       // Deneb-specific field, use pointer for optionality
	ExcessBlobGas *uint64       // Deneb-specific field, use pointer for optionality
}

type SignedBlsToExecutionChange struct {
	Message   BlsToExecutionChange
	Signature SignatureBytes
}

type BlsToExecutionChange struct {
	ValidatorIndex uint64
	FromBlsPubkey  [48]byte
}

// BeaconBlockBody represents the body of a beacon block.
type BeaconBlockBody struct {
	RandaoReveal          SignatureBytes
	Eth1Data              Eth1Data
	Graffiti              Bytes32
	ProposerSlashings     []ProposerSlashing    `ssz-max:"16"`
	AttesterSlashings     []AttesterSlashing    `ssz-max:"2"`
	Attestations          []Attestation         `ssz-max:"128"`
	Deposits              []Deposit             `ssz-max:"16"`
	VoluntaryExits        []SignedVoluntaryExit `ssz-max:"16"`
	SyncAggregate         SyncAggregate
	ExecutionPayload      ExecutionPayload
	BlsToExecutionChanges []SignedBlsToExecutionChange `ssz-max:"16"`
	BlobKzgCommitments    [][48]byte                   `ssz-max:"4096"`
}

type Header struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    Bytes32
	StateRoot     Bytes32
	BodyRoot      Bytes32
}

type SyncCommittee struct {
	Pubkeys         [512]BLSPubKey
	AggregatePubkey BLSPubKey
}

type Update struct {
	AttestedHeader          Header
	NextSyncCommittee       SyncCommittee
	NextSyncCommitteeBranch []Bytes32
	FinalizedHeader         Header
	FinalityBranch          []Bytes32
	SyncAggregate           SyncAggregate
	SignatureSlot           uint64
}

type FinalityUpdate struct {
	AttestedHeader  Header
	FinalizedHeader Header
	FinalityBranch  []Bytes32
	SyncAggregate   SyncAggregate
	SignatureSlot   uint64
}

type OptimisticUpdate struct {
	AttestedHeader Header
	SyncAggregate  SyncAggregate
	SignatureSlot  uint64
}

type Bootstrap struct {
	Header                     Header
	CurrentSyncCommittee       SyncCommittee
	CurrentSyncCommitteeBranch []Bytes32
}

type Fork struct {
	Epoch       uint64 `json:"epoch"`
	ForkVersion []byte `json:"fork_version"`
}
type Forks struct {
	Genesis   Fork `json:"genesis"`
	Altair    Fork `json:"altair"`
	Bellatrix Fork `json:"bellatrix"`
	Capella   Fork `json:"capella"`
	Deneb     Fork `json:"deneb"`
}

// ToBytes serializes the Header struct to a byte slice.
func (h *Header) ToBytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, new(codec.JsonHandle))
	_ = enc.Encode(h) // Ignore error
	return buf.Bytes()
}

func (b *BeaconBlockBody) ToBytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, new(codec.JsonHandle))
	_ = enc.Encode(b) // Ignore error
	return buf.Bytes()
}

// ToBytes serializes the SyncCommittee struct to a byte slice.
func (sc *SyncCommittee) ToBytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, new(codec.JsonHandle))
	_ = enc.Encode(sc) // Ignore error
	return buf.Bytes()
}
