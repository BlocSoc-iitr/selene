package rpc

import (
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockConsensusRpc is a mock implementation of the ConsensusRpc interface
type MockConsensusRpc struct {
	mock.Mock
}

func (m *MockConsensusRpc) GetBootstrap(block_root [32]byte) (consensus_core.Bootstrap, error) {
	args := m.Called(block_root)
	return args.Get(0).(consensus_core.Bootstrap), args.Error(1)
}
func (m *MockConsensusRpc) GetUpdates(period uint64, count uint8) ([]consensus_core.Update, error) {
	args := m.Called(period, count)
	return args.Get(0).([]consensus_core.Update), args.Error(1)
}
func (m *MockConsensusRpc) GetFinalityUpdate() (consensus_core.FinalityUpdate, error) {
	args := m.Called()
	return args.Get(0).(consensus_core.FinalityUpdate), args.Error(1)
}
func (m *MockConsensusRpc) GetOptimisticUpdate() (consensus_core.OptimisticUpdate, error) {
	args := m.Called()
	return args.Get(0).(consensus_core.OptimisticUpdate), args.Error(1)
}
func (m *MockConsensusRpc) GetBlock(slot uint64) (consensus_core.BeaconBlock, error) {
	args := m.Called(slot)
	return args.Get(0).(consensus_core.BeaconBlock), args.Error(1)
}
func (m *MockConsensusRpc) ChainId() (uint64, error) {
	args := m.Called()
	return args.Get(0).(uint64), args.Error(1)
}
func TestNewConsensusRpc(t *testing.T) {
	rpcURL := "http://example.com"
	consensusRpc := NewConsensusRpc(rpcURL)
	assert.Implements(t, (*ConsensusRpc)(nil), consensusRpc)
	_, ok := consensusRpc.(*NimbusRpc)
	assert.True(t, ok, "NewConsensusRpc should return a *NimbusRpc")
}

func TestConsensusRpcInterface(t *testing.T) {
	mockRpc := new(MockConsensusRpc)
	// Test GetBootstrap
	mockBootstrap := consensus_core.Bootstrap{Header: consensus_core.Header{Slot: 1000}}
	mockRpc.On("GetBootstrap", mock.Anything).Return(mockBootstrap, nil)
	bootstrap, err := mockRpc.GetBootstrap([32]byte{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000), bootstrap.Header.Slot)
	// Test GetUpdates
	mockUpdates := []consensus_core.Update{{AttestedHeader: consensus_core.Header{Slot: 2000}}}
	mockRpc.On("GetUpdates", uint64(1), uint8(5)).Return(mockUpdates, nil)
	updates, err := mockRpc.GetUpdates(1, 5)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2000), updates[0].AttestedHeader.Slot)
	// Test GetFinalityUpdate
	mockFinalityUpdate := consensus_core.FinalityUpdate{FinalizedHeader: consensus_core.Header{Slot: 3000}}
	mockRpc.On("GetFinalityUpdate").Return(mockFinalityUpdate, nil)
	finalityUpdate, err := mockRpc.GetFinalityUpdate()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3000), finalityUpdate.FinalizedHeader.Slot)
	// Test GetOptimisticUpdate
	mockOptimisticUpdate := consensus_core.OptimisticUpdate{SignatureSlot: 4000}
	mockRpc.On("GetOptimisticUpdate").Return(mockOptimisticUpdate, nil)
	optimisticUpdate, err := mockRpc.GetOptimisticUpdate()
	assert.NoError(t, err)
	assert.Equal(t, uint64(4000), optimisticUpdate.SignatureSlot)
	// Test GetBlock
	mockBlock := consensus_core.BeaconBlock{Slot: 5000}
	mockRpc.On("GetBlock", uint64(5000)).Return(mockBlock, nil)
	block, err := mockRpc.GetBlock(5000)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5000), block.Slot)
	// Test ChainId
	mockChainId := uint64(1)
	mockRpc.On("ChainId").Return(mockChainId, nil)
	chainId, err := mockRpc.ChainId()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), chainId)
	// Assert that all expected mock calls were made
	mockRpc.AssertExpectations(t)
}
