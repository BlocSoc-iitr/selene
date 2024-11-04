package consensus_core

type Transaction = [1073741824]byte
type Bytes32 [32]byte
type BLSPubKey [48]byte
type Address [20]byte
type LogsBloom [256]byte
type SignatureBytes [96]byte

type BeaconBlock struct {
	Slot          uint64          `json:"slot"`
	ProposerIndex uint64          `json:"proposer_index"`
	ParentRoot    Bytes32         `json:"parent_root"`
	StateRoot     Bytes32         `json:"state_root"`
	Body          BeaconBlockBody `json:"body"`
	Hash          []byte          
}

type Eth1Data struct {
	DepositRoot  Bytes32 `json:"deposit_root"`
	DepositCount uint64  `json:"deposit_count"`
	BlockHash    Bytes32 `json:"block_hash"`
}

type ProposerSlashing struct {
	SignedHeader1 SignedBeaconBlockHeader `json:"signed_header_1"`
	SignedHeader2 SignedBeaconBlockHeader `json:"signed_header_2"`
}

type SignedBeaconBlockHeader struct {
	Message   BeaconBlockHeader `json:"message"`
	Signature SignatureBytes    `json:"signature"`
}

type BeaconBlockHeader struct {
	Slot          uint64  `json:"slot"`
	ProposerIndex uint64  `json:"proposer_index"`
	ParentRoot    Bytes32 `json:"parent_root"`
	StateRoot     Bytes32 `json:"state_root"`
	BodyRoot      Bytes32 `json:"body_root"`
}

type AttesterSlashing struct {
	Attestation1 IndexedAttestation `json:"attestation_1"`
	Attestation2 IndexedAttestation `json:"attestation_2"`
}

type IndexedAttestation struct {
	AttestingIndices [2048]uint64    `json:"attesting_indices"` // Dynamic slice
	Data             AttestationData `json:"data"`
	Signature        SignatureBytes  `json:"signature"`
}

type AttestationData struct {
	Slot            uint64     `json:"slot"`
	Index           uint64     `json:"index"`
	BeaconBlockRoot Bytes32    `json:"beacon_block_root"`
	Source          Checkpoint `json:"source"`
	Target          Checkpoint `json:"target"`
}

type Checkpoint struct {
	Epoch uint64  `json:"epoch"`
	Root  Bytes32 `json:"root"`
}

type Bitlist []bool

type Attestation struct {
	AggregationBits Bitlist         `json:"aggregation_bits"`
	Data            AttestationData `json:"data"`
	Signature       SignatureBytes  `json:"signature"`
}

type Deposit struct {
	Proof [33]Bytes32 `json:"proof"` // Dynamic slice
	Data  DepositData `json:"data"`
}

type DepositData struct {
	Pubkey                BLSPubKey      `json:"pubkey"`
	WithdrawalCredentials Bytes32        `json:"withdrawal_credentials"`
	Amount                uint64         `json:"amount"`
	Signature             SignatureBytes `json:"signature"`
}

type SignedVoluntaryExit struct {
	Message   VoluntaryExit  `json:"message"`
	Signature SignatureBytes `json:"signature"`
}

type VoluntaryExit struct {
	Epoch          uint64 `json:"epoch"`
	ValidatorIndex uint64 `json:"validator_index"`
}

type SyncAggregate struct {
	SyncCommitteeBits      []byte         `json:"sync_committee_bits"`
	SyncCommitteeSignature SignatureBytes `json:"sync_committee_signature"`
}

type Withdrawal struct {
	Index          uint64  `json:"index"`
	ValidatorIndex uint64  `json:"validator_index"`
	Address        Address `json:"address"`
	Amount         uint64  `json:"amount"`
}

type ExecutionPayload struct {
	ParentHash    Bytes32       `json:"parent_hash"`
	FeeRecipient  Address       `json:"fee_recipient"`
	StateRoot     Bytes32       `json:"state_root"`
	ReceiptsRoot  Bytes32       `json:"receipts_root"`
	LogsBloom     LogsBloom     `json:"logs_bloom"`
	PrevRandao    Bytes32       `json:"prev_randao"`
	BlockNumber   uint64        `json:"block_number"`
	GasLimit      uint64        `json:"gas_limit"`
	GasUsed       uint64        `json:"gas_used"`
	Timestamp     uint64        `json:"timestamp"`
	ExtraData     []byte        `json:"extra_data"`
	BaseFeePerGas uint64        `json:"base_fee_per_gas"`
	BlockHash     Bytes32       `json:"block_hash"`
	Transactions  [][]byte      `json:"transactions"`
	Withdrawals   *[]Withdrawal `json:"withdrawals"`     //Only capella and deneb
	BlobGasUsed   *uint64       `json:"blob_gas_used"`   // Only deneb
	ExcessBlobGas *uint64       `json:"excess_blob_gas"` // Only deneb
}

type SignedBlsToExecutionChange struct {
	Message   BlsToExecutionChange `json:"message"`
	Signature SignatureBytes       `json:"signature"`
}

type BlsToExecutionChange struct {
	ValidatorIndex uint64    `json:"validator_index"`
	FromBlsPubkey  BLSPubKey `json:"from_bls_pubkey"`
}

// BeaconBlockBody represents the body of a beacon block.
type BeaconBlockBody struct {
	RandaoReveal          SignatureBytes                `json:"randao_reveal"`
	Eth1Data              Eth1Data                      `json:"eth1_data"`
	Graffiti              Bytes32                       `json:"graffiti"`
	ProposerSlashings     []ProposerSlashing            `json:"proposer_slashings"`
	AttesterSlashings     []AttesterSlashing            `json:"attester_slashings"`
	Attestations          []Attestation                 `json:"attestations"`
	Deposits              []Deposit                     `json:"deposits"`
	VoluntaryExits        []SignedVoluntaryExit         `json:"voluntary_exits"`
	SyncAggregate         SyncAggregate                 `json:"sync_aggregate"`
	ExecutionPayload      ExecutionPayload              `json:"execution_payload"`
	BlsToExecutionChanges *[]SignedBlsToExecutionChange `json:"bls_to_execution_changes"` //Only capella and deneb
	BlobKzgCommitments    *[][]byte                     `json:"blob_kzg_commitments"`     // Dynamic slice
}

type Header struct {
	Slot          uint64  `json:"slot"`
	ProposerIndex uint64  `json:"proposer_index"`
	ParentRoot    Bytes32 `json:"parent_root"`
	StateRoot     Bytes32 `json:"state_root"`
	BodyRoot      Bytes32 `json:"body_root"`
}

type SyncCommittee struct {
	Pubkeys         []BLSPubKey `json:"pubkeys"`
	AggregatePubkey BLSPubKey   `json:"aggregate_pubkey"`
}

type Update struct {
	AttestedHeader          Header        `json:"attested_header"`
	NextSyncCommittee       SyncCommittee `json:"next_sync_committee"`
	NextSyncCommitteeBranch []Bytes32     `json:"next_sync_committee_branch"`
	FinalizedHeader         Header        `json:"finalized_header"`
	FinalityBranch          []Bytes32     `json:"finality_branch"`
	SyncAggregate           SyncAggregate `json:"sync_aggregate"`
	SignatureSlot           uint64        `json:"signature_slot"`
}

type FinalityUpdate struct {
	AttestedHeader  Header        `json:"attested_header"`
	FinalizedHeader Header        `json:"finalized_header"`
	FinalityBranch  []Bytes32     `json:"finality_branch"`
	SyncAggregate   SyncAggregate `json:"sync_aggregate"`
	SignatureSlot   uint64        `json:"signature_slot"`
}

type OptimisticUpdate struct {
	AttestedHeader Header        `json:"attested_header"`
	SyncAggregate  SyncAggregate `json:"sync_aggregate"`
	SignatureSlot  uint64        `json:"signature_slot"`
}

type Bootstrap struct {
	Header                     Header        `json:"header"`
	CurrentSyncCommittee       SyncCommittee `json:"current_sync_committee"`
	CurrentSyncCommitteeBranch []Bytes32     `json:"current_sync_committee_branch"`
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
