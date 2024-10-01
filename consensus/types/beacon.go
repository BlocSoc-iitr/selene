package types

import (
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
)

// Serialization and Deserialization for ExecutionPayload and BeaconBlockBody can be done by importing from prewritten functions in utils wherever needed.

type GenericUpdate struct {
	AttestedHeader          consensus_core.Header
	SyncAggregate           consensus_core.SyncAggregate
	SignatureSlot           uint64
	NextSyncCommittee       consensus_core.SyncCommittee
	NextSyncCommitteeBranch []consensus_core.Bytes32
	FinalizedHeader         consensus_core.Header
	FinalityBranch          []consensus_core.Bytes32
}

func (g *GenericUpdate) From() {
	g.AttestedHeader = consensus_core.Header{}
	g.SyncAggregate = consensus_core.SyncAggregate{}
	g.SignatureSlot = 0
	g.NextSyncCommittee = consensus_core.SyncCommittee{}
	g.NextSyncCommitteeBranch = []consensus_core.Bytes32{}
	g.FinalizedHeader = consensus_core.Header{}
	g.FinalityBranch = []consensus_core.Bytes32{}
}
