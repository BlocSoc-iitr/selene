package rpc

import (
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"

)

// return types not mention and oarameters as well
type ConsensusRpc interface {
	GetBootstrap(block_root [32]byte) (consensus_core.Bootstrap, error)
	GetUpdates(period uint64, count uint8) ([]consensus_core.Update, error)
	GetFinalityUpdate() (consensus_core.FinalityUpdate, error)
	GetOptimisticUpdate() (consensus_core.OptimisticUpdate, error)
	GetBlock(slot uint64) (consensus_core.BeaconBlock, error)
	ChainId() (uint64, error)
}


func NewConsensusRpc(rpc string) ConsensusRpc {
	if rpc == "testdata/" {
		return NewMockRpc(rpc)
	} else {
		return NewNimbusRpc(rpc)
	}
}
