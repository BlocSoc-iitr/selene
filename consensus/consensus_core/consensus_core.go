// go implementation of consensus_core/src/types/mod.rs
// structs not defined yet:
// LightClientStore,genericUpdate
package consensus_core

import "github.com/BlocSoc-iitr/selene/consensus/types"

type BeaconBlock struct {
	slot           uint64
	proposer_index uint64
	parent_root    [32]byte
	state_root     [32]byte
	Body           BeaconBlockBody
}
type Bytes32 [32]byte
type BLSPubKey [48]byte
type Address [20]byte
type LogsBloom = [256]byte
type SignatureBytes [96]byte
type Eth1Data struct {
	deposit_root  Bytes32
	deposit_count uint64
	block_hash    Bytes32
}
type ProposerSlashing struct {
	signed_header_1 SignedBeaconBlockHeader
	signed_header_2 SignedBeaconBlockHeader
}
type SignedBeaconBlockHeader struct {
	message   BeaconBlockHeader
	signature SignatureBytes
}
type BeaconBlockHeader struct {
	slot           uint64
	proposer_index uint64
	Parent_root    Bytes32
	state_root     Bytes32
	body_root      Bytes32
}
type AttesterSlashing struct {
	attestation_1 IndexedAttestation
	attestation_2 IndexedAttestation
}
type IndexedAttestation struct {
	attesting_indices []uint64 //max length 2048 to be ensured
	data              AttestationData
	signature         SignatureBytes
}
type AttestationData struct {
	slot              uint64
	index             uint64
	beacon_block_root Bytes32
	source            Checkpoint
	target            Checkpoint
}
type Checkpoint struct {
	epoch uint64
	root  Bytes32
}
type Bitlist []bool
type Attestation struct {
	aggregation_bits Bitlist `ssz-max:"2048"`
	data             AttestationData
	signature        SignatureBytes
}
type Deposit struct {
	proof [33]Bytes32 //fixed size array
	data  DepositData
}
type DepositData struct {
	pubkey                 [48]byte
	withdrawal_credentials Bytes32
	amount                 uint64
	signature              SignatureBytes
}
type SignedVoluntaryExit struct {
	message   VoluntaryExit
	signature SignatureBytes
}
type VoluntaryExit struct {
	epoch           uint64
	validator_index uint64
}
type SyncAggregate struct {
	Sync_committee_bits      [64]byte
	Sync_committee_signature SignatureBytes
}
type Withdrawal struct {
	index           uint64
	validator_index uint64
	address         Address
	amount          uint64
}
type ExecutionPayload struct { //not implemented
	Parent_hash      Bytes32
	Fee_recipient    Address
	State_root       Bytes32
	Receipts_root    Bytes32
	Logs_bloom       LogsBloom
	Prev_randao      Bytes32
	Block_number     uint64
	Gas_limit        uint64
	Gas_used         uint64
	Timestamp        uint64
	Extra_data       [32]byte
	Base_fee_per_gas uint64
	Block_hash       Bytes32
	Transactions     types.Transaction
	withdrawals      types.Withdrawal
	blob_gas_used    uint64
	excess_blob_gas  uint64
}
type SignedBlsToExecutionChange struct {
	message   BlsToExecutionChange
	signature SignatureBytes
}
type BlsToExecutionChange struct {
	validator_index uint64
	from_bls_pubkey [48]byte
}
type BeaconBlockBody struct { //not implemented
	randao_reveal            SignatureBytes
	eth1_data                Eth1Data
	graffiti                 Bytes32
	proposer_slashings       []ProposerSlashing //max length 16 to be insured how?
	attester_slashings       []AttesterSlashing //max length 2 to be ensured
	attestations             []Attestation      //max length 128 to be ensured
	deposits                 []Deposit          //max length 16 to be ensured
	voluntary_exits          SignedVoluntaryExit
	sync_aggregate           SyncAggregate
	Execution_payload        ExecutionPayload
	bls_to_execution_changes []SignedBlsToExecutionChange //max length 16 to be ensured
	blob_kzg_commitments     [][48]byte                   //max length 4096 to be ensured
}
type Header struct {
	Slot           uint64
	Proposer_index uint64
	Parent_root    Bytes32
	State_root     Bytes32
	Body_root      Bytes32
}
type SyncComittee struct {
	Pubkeys          [512]BLSPubKey
	Aggregate_pubkey BLSPubKey
}
type Update struct {
	AttestedHeader         Header
	NextSyncCommitee       SyncComittee
	NextSyncCommiteeBranch []Bytes32
	FinalizedHeader        Header
	FinalityBranch         []Bytes32
	SyncAggregate          SyncAggregate
	SignatureSlot          uint64
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
	CurrentSyncCommittee       SyncComittee
	CurrentSyncCommitteeBranch []Bytes32
}
