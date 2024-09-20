package consensus_core

import "github.com/BlocSoc-iitr/selene/consensus/types"

type BeaconBlock struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    [32]byte
	StateRoot     [32]byte
	Body          BeaconBlockBody
}

type Bytes32 [32]byte
type BLSPubKey [48]byte
type Address [20]byte
type LogsBloom = [256]byte
type SignatureBytes [96]byte

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
	AttestingIndices []uint64 //max length 2048 to be ensured
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
	Proof [33]Bytes32 //fixed size array
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

type ExecutionPayload struct { //not implemented
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
	Transactions  types.Transaction
	Withdrawals   types.Withdrawal
	BlobGasUsed   uint64
	ExcessBlobGas uint64
}

type SignedBlsToExecutionChange struct {
	Message   BlsToExecutionChange
	Signature SignatureBytes
}

type BlsToExecutionChange struct {
	ValidatorIndex uint64
	FromBlsPubkey  [48]byte
}

type BeaconBlockBody struct { //not implemented
	RandaoReveal          SignatureBytes
	Eth1Data              Eth1Data
	Graffiti              Bytes32
	ProposerSlashings     []ProposerSlashing //max length 16 to be insured how?
	AttesterSlashings     []AttesterSlashing //max length 2 to be ensured
	Attestations          []Attestation      //max length 128 to be ensured
	Deposits              []Deposit          //max length 16 to be ensured
	VoluntaryExits        SignedVoluntaryExit
	SyncAggregate         SyncAggregate
	ExecutionPayload      ExecutionPayload
	BlsToExecutionChanges []SignedBlsToExecutionChange //max length 16 to be ensured
	BlobKzgCommitments    [][48]byte                   //max length 4096 to be ensured
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
