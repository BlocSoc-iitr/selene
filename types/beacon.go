package types

type BeaconBlockHeader struct {
	slot          uint64
	proposerIndex uint64
	parentRoot    Bytes32
	stateRoot     Bytes32
	bodyRoot      Bytes32
}

type SyncCommittee struct {
	pubkeys   [512]Bytes48
	aggPubKey Bytes48
}

type SyncAggregate struct {
	syncCommitteeBits      [512]Bitvector
	syncCommitteeSignature Bytes48
}
