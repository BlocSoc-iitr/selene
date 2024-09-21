package types

// Serialization and Deserialization for ExecutionPayload and BeaconBlockBody can be done by importing from prewritten functions in utils wherever needed.



type BeaconBlock struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    [32]byte
	StateRoot     [32]byte
	Body          BeaconBlockBody
}

// just implemented the deneb variant
type BeaconBlockBody struct {
	RandaoReveal          SignatureBytes                 `json:"randao_reveal"`
	Eth1Data              Eth1Data                       `json:"eth1_data"`
	Graffiti              Bytes32                        `json:"graffiti"`
	ProposerSlashings     [16]ProposerSlashing           `json:"proposer_slashings"` // max length 16 to be ensured
	AttesterSlashings     [2]AttesterSlashing            `json:"attester_slashings"` // max length 2 to be ensured
	Attestations          [128]Attestation               `json:"attestations"`       // max length 128 to be ensured
	Deposits              [16]Deposit                    `json:"deposits"`           // max length 16 to be ensured
	VoluntaryExits        [16]SignedVoluntaryExit        `json:"voluntary_exits"`    // max length 16 to be ensured
	SyncAggregate         SyncAggregate                  `json:"sync_aggregate"`
	ExecutionPayload      ExecutionPayload               `json:"execution_payload"`
	BlsToExecutionChanges [16]SignedBlsToExecutionChange `json:"bls_to_execution_changes"` // max length 16 to be ensured
	BlobKzgCommitments    [4096][48]byte                 `json:"blob_kzg_commitments"`     // max length 4096 to be ensured
} // 1 method

func (b *BeaconBlockBody) Def() {
	b.RandaoReveal = SignatureBytes{}
	b.Eth1Data = Eth1Data{}
	b.Graffiti = Bytes32{}
	b.ProposerSlashings = [16]ProposerSlashing{}
	b.AttesterSlashings = [2]AttesterSlashing{}
	b.Attestations = [128]Attestation{}
	b.Deposits = [16]Deposit{}
	b.VoluntaryExits = [16]SignedVoluntaryExit{}
	b.SyncAggregate = SyncAggregate{}
	b.ExecutionPayload = ExecutionPayload{}
	b.BlsToExecutionChanges = [16]SignedBlsToExecutionChange{}
	b.BlobKzgCommitments = [4096][48]byte{}
}

type SignedBlsToExecutionChange struct {
	Message   BlsToExecutionChange
	Signature SignatureBytes
}
type BlsToExecutionChange struct {
	ValidatorIndex uint64
	FromBlsPubkey  [48]byte
}

type ExecutionPayload struct {
	ParentHash    Bytes32          `json:"parent_hash"`
	FeeRecipient  Address          `json:"fee_recipient"`
	StateRoot     Bytes32          `json:"state_root"`
	ReceiptsRoot  Bytes32          `json:"receipts_root"`
	LogsBloom     LogsBloom        `json:"logs_bloom"`
	PrevRandao    Bytes32          `json:"prev_randao"`
	BlockNumber   uint64           `json:"block_number"`
	GasLimit      uint64           `json:"gas_limit"`
	GasUsed       uint64           `json:"gas_used"`
	Timestamp     uint64           `json:"timestamp"`
	ExtraData     [32]byte         `json:"extra_data"`
	BaseFeePerGas uint64           `json:"base_fee_per_gas"`
	BlockHash     Bytes32          `json:"block_hash"`
	Transactions  [1073741824]byte `json:"transactions"`    // max length 1073741824 to be implemented
	Withdrawals   [16]Withdrawal   `json:"withdrawals"`     // max length 16 to be implemented
	BlobGasUsed   uint64           `json:"blob_gas_used"`   // Deneb variant only
	ExcessBlobGas uint64           `json:"excess_blob_gas"` // Deneb variant only
} // 1 method

func (exe *ExecutionPayload) Def() {
	exe.ParentHash = Bytes32{}
	exe.FeeRecipient = Address{}
	exe.StateRoot = Bytes32{}
	exe.ReceiptsRoot = Bytes32{}
	exe.LogsBloom = LogsBloom{}
	exe.PrevRandao = Bytes32{}
	exe.BlockNumber = 0
	exe.GasLimit = 0
	exe.GasUsed = 0
	exe.Timestamp = 0
	exe.ExtraData = [32]byte{}
	exe.BaseFeePerGas = 0
	exe.BlockHash = Bytes32{}
	exe.Transactions = [1073741824]byte{} // max length 1073741824 to be implemented
	exe.Withdrawals = [16]Withdrawal{}    // max length 16 to be implemented
	exe.BlobGasUsed = 0                   // Only for Deneb
	exe.ExcessBlobGas = 0                 // Only for Deneb
} // default



type Withdrawal struct {
	Index          uint64
	Amount         uint64
	Address        Address
	ValidatorIndex uint64
}
type ProposalSlashing struct {
	SignedHeader1 SignedBeaconBlockHeader
	SignedHeader2 SignedBeaconBlockHeader
}
type SignedBeaconBlockHeader struct {
	Message   BeaconBlockHeader
	Signature SignatureBytes
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
type Attestation struct {
	AggregateBits [2048]bool
	Data          AttestationData
	Signature     SignatureBytes
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
type SignedVoluntaryExit struct {
	Message   VoluntaryExit
	Signature SignatureBytes
}
type VoluntaryExit struct {
	Epoch          uint64
	ValidatorIndex uint64
}
type Deposit struct {
	Proof [33]byte
	Data  DepositData
}
type DepositData struct {
	Pubkey               [48]byte
	WitdrawalCredentials [32]byte
	Amount               uint64
	Signature            SignatureBytes
}
type Eth1Data struct {
	DepositRoot  Bytes32
	DepositCount uint64
	BlockHash    Bytes32
}
type BeaconBlockHeader struct {
	Slot         uint64
	ProposeIndex uint64
	ParentRoot   Bytes32
	StateRoot    Bytes32
	BodyRoot     Bytes32
}
type Update struct {
	AttestedHeader         Header
	NextSyncCommittee      SyncCommittee
	NextSyncCommitteBranch []Bytes32
	FinalizedHeader        Header
	FinalizedBranch        []Bytes32
	SyncAggregate          SyncAggregate
	SignatureSlot          uint64
}
type FinalityUpdate struct {
	AttestedHeader  Header
	FinalizedHeader Header
	FinalizedBranch []Bytes32
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
type SyncAggregate struct {
	SyncCommitteeBits      [512]bool
	SyncCommitteeSignature SignatureBytes
}
type ProposerSlashing struct {
	SignedHeader1 SignedBeaconBlockHeader
	SignedHeader2 SignedBeaconBlockHeader
}
type GenericUpdate struct {
	AttestedHeader          Header
	SyncAggregate           SyncAggregate
	SignatureSlot           uint64
	NextSyncCommittee       SyncCommittee
	NextSyncCommitteeBranch []Bytes32
	FinalizedHeader         Header
	FinalityBranch          []Bytes32
}

func (g *GenericUpdate) From() {
	g.AttestedHeader = Header{}
	g.SyncAggregate = SyncAggregate{}
	g.SignatureSlot = 0
	g.NextSyncCommittee = SyncCommittee{}
	g.NextSyncCommitteeBranch = [][32]byte{}
	g.FinalizedHeader = Header{}
	g.FinalityBranch = [][32]byte{}
}
