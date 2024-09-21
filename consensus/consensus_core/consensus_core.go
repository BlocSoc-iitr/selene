package consensus_core

import (
	"encoding/binary"

	"github.com/BlocSoc-iitr/selene/consensus/types"
)

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
	Transactions  []types.Transaction `ssz-max:"1048576"`
	Withdrawals   []types.Withdrawal  `ssz-max:"16"`
	BlobGasUsed   *uint64             // Deneb-specific field, use pointer for optionality
	ExcessBlobGas *uint64             // Deneb-specific field, use pointer for optionality
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

func (h *Header) ToBytes() [][]byte {
	data := make([][]byte, 5)

	slotBytes := make([]byte, 8)
	proposerIndexBytes := make([]byte, 8)

	binary.LittleEndian.PutUint64(slotBytes, h.Slot)
	binary.LittleEndian.PutUint64(proposerIndexBytes, h.ProposerIndex)

	data[0] = slotBytes
	data[1] = proposerIndexBytes
	data[2] = h.ParentRoot[:]
	data[3] = h.StateRoot[:]
	data[4] = h.BodyRoot[:]

	return data
}

func (sc *SyncCommittee) ToBytes() [][]byte {
	data := make([][]byte, 513) // 512 Pubkeys + 1 AggregatePubkey

	// Convert each BLS pubkey in Pubkeys array to []byte
	for i, pubkey := range sc.Pubkeys {
		data[i] = pubkey[:] // Convert BLSPubKey (assumed to be a fixed size array) to []byte
	}

	// Convert the AggregatePubkey to []byte
	data[512] = sc.AggregatePubkey[:]

	return data
}

func (body *BeaconBlockBody) ToBytes() [][]byte {
    var data [][]byte

    // RandaoReveal
    data = append(data, body.RandaoReveal[:])

    // Eth1Data
    data = append(data, body.Eth1Data.DepositRoot[:])
    depositCountBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(depositCountBytes, body.Eth1Data.DepositCount)
    data = append(data, depositCountBytes)
    data = append(data, body.Eth1Data.BlockHash[:])

    // Graffiti
    data = append(data, body.Graffiti[:])

    // ProposerSlashings
    for _, proposerSlashing := range body.ProposerSlashings {
        // SignedHeader1
        slotBytes := make([]byte, 8)
        proposerIndexBytes := make([]byte, 8)
        binary.LittleEndian.PutUint64(slotBytes, proposerSlashing.SignedHeader1.Message.Slot)
        binary.LittleEndian.PutUint64(proposerIndexBytes, proposerSlashing.SignedHeader1.Message.ProposerIndex)
        data = append(data, slotBytes)
        data = append(data, proposerIndexBytes)
        data = append(data, proposerSlashing.SignedHeader1.Message.ParentRoot[:])
        data = append(data, proposerSlashing.SignedHeader1.Message.StateRoot[:])
        data = append(data, proposerSlashing.SignedHeader1.Message.BodyRoot[:])
        data = append(data, proposerSlashing.SignedHeader1.Signature[:])

        // SignedHeader2
        slotBytes = make([]byte, 8)
        proposerIndexBytes = make([]byte, 8)
        binary.LittleEndian.PutUint64(slotBytes, proposerSlashing.SignedHeader2.Message.Slot)
        binary.LittleEndian.PutUint64(proposerIndexBytes, proposerSlashing.SignedHeader2.Message.ProposerIndex)
        data = append(data, slotBytes)
        data = append(data, proposerIndexBytes)
        data = append(data, proposerSlashing.SignedHeader2.Message.ParentRoot[:])
        data = append(data, proposerSlashing.SignedHeader2.Message.StateRoot[:])
        data = append(data, proposerSlashing.SignedHeader2.Message.BodyRoot[:])
        data = append(data, proposerSlashing.SignedHeader2.Signature[:])
    }

    // AttesterSlashings
    for _, attesterSlashing := range body.AttesterSlashings {
        // Attestation1
        for _, index := range attesterSlashing.Attestation1.AttestingIndices {
            indexBytes := make([]byte, 8)
            binary.LittleEndian.PutUint64(indexBytes, index)
            data = append(data, indexBytes)
        }
        data = append(data, attesterSlashing.Attestation1.Data.BeaconBlockRoot[:])
        data = append(data, attesterSlashing.Attestation1.Signature[:])

        // Attestation2
        for _, index := range attesterSlashing.Attestation2.AttestingIndices {
            indexBytes := make([]byte, 8)
            binary.LittleEndian.PutUint64(indexBytes, index)
            data = append(data, indexBytes)
        }
        data = append(data, attesterSlashing.Attestation2.Data.BeaconBlockRoot[:])
        data = append(data, attesterSlashing.Attestation2.Signature[:])
    }

    // Attestations
    for _, attestation := range body.Attestations {
        // AggregationBits (Bitlist)
        byteLen := (len(attestation.AggregationBits) + 7) / 8
        aggregationBitsBytes := make([]byte, byteLen)
        for i, bit := range attestation.AggregationBits {
            if bit {
                aggregationBitsBytes[i/8] |= 1 << (i % 8)
            }
        }
        data = append(data, aggregationBitsBytes)

        // AttestationData
        slotBytes := make([]byte, 8)
        indexBytes := make([]byte, 8)
        binary.LittleEndian.PutUint64(slotBytes, attestation.Data.Slot)
        binary.LittleEndian.PutUint64(indexBytes, attestation.Data.Index)
        data = append(data, slotBytes)
        data = append(data, indexBytes)
        data = append(data, attestation.Data.BeaconBlockRoot[:])
        data = append(data, attestation.Data.Source.Root[:])
        data = append(data, attestation.Data.Target.Root[:])

        // Signature
        data = append(data, attestation.Signature[:])
    }

    // Deposits
    for _, deposit := range body.Deposits {
        // Proof
        for _, proof := range deposit.Proof {
            data = append(data, proof[:])
        }

        // DepositData
        data = append(data, deposit.Data.Pubkey[:])
        data = append(data, deposit.Data.WithdrawalCredentials[:])
        amountBytes := make([]byte, 8)
        binary.LittleEndian.PutUint64(amountBytes, deposit.Data.Amount)
        data = append(data, amountBytes)
        data = append(data, deposit.Data.Signature[:])
    }

    // VoluntaryExits
    for _, exit := range body.VoluntaryExits {
        epochBytes := make([]byte, 8)
        validatorIndexBytes := make([]byte, 8)
        binary.LittleEndian.PutUint64(epochBytes, exit.Message.Epoch)
        binary.LittleEndian.PutUint64(validatorIndexBytes, exit.Message.ValidatorIndex)
        data = append(data, epochBytes)
        data = append(data, validatorIndexBytes)
        data = append(data, exit.Signature[:])
    }

    // SyncAggregate
    data = append(data, body.SyncAggregate.SyncCommitteeBits[:])
    data = append(data, body.SyncAggregate.SyncCommitteeSignature[:])

    // ExecutionPayload
    data = append(data, body.ExecutionPayload.ParentHash[:])
    data = append(data, body.ExecutionPayload.FeeRecipient[:])
    data = append(data, body.ExecutionPayload.StateRoot[:])
    data = append(data, body.ExecutionPayload.ReceiptsRoot[:])
    data = append(data, body.ExecutionPayload.LogsBloom[:])
    data = append(data, body.ExecutionPayload.PrevRandao[:])
    blockNumberBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(blockNumberBytes, body.ExecutionPayload.BlockNumber)
    data = append(data, blockNumberBytes)
    gasLimitBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(gasLimitBytes, body.ExecutionPayload.GasLimit)
    data = append(data, gasLimitBytes)
    gasUsedBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(gasUsedBytes, body.ExecutionPayload.GasUsed)
    data = append(data, gasUsedBytes)
    timestampBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(timestampBytes, body.ExecutionPayload.Timestamp)
    data = append(data, timestampBytes)
    data = append(data, body.ExecutionPayload.ExtraData[:])
    baseFeePerGasBytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(baseFeePerGasBytes, body.ExecutionPayload.BaseFeePerGas)
    data = append(data, baseFeePerGasBytes)
    data = append(data, body.ExecutionPayload.BlockHash[:])

    // Transactions (list of transactions)
    for _, tx := range body.ExecutionPayload.Transactions {
        data = append(data, tx[:]) // Assuming tx.Data() returns the transaction bytes
    }

    // Withdrawals
    for _, withdrawal := range body.ExecutionPayload.Withdrawals {
        indexBytes := make([]byte, 8)
        validatorIndexBytes := make([]byte, 8)
        binary.LittleEndian.PutUint64(indexBytes, withdrawal.Index)
        binary.LittleEndian.PutUint64(validatorIndexBytes, withdrawal.ValidatorIndex)
        data = append(data, indexBytes)
        data = append(data, validatorIndexBytes)
        data = append(data, withdrawal.Address[:])
        amountBytes := make([]byte, 8)
        binary.LittleEndian.PutUint64(amountBytes, withdrawal.Amount)
        data = append(data, amountBytes)
    }

    // BlobKzgCommitments
    for _, commitment := range body.BlobKzgCommitments {
        data = append(data, commitment[:])
    }

    return data
}

