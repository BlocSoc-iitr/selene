package rpc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
)

func TestNewMockRpc(t *testing.T) {
	path := "/tmp/testdata"
	mockRpc := NewMockRpc(path)
	if mockRpc.testdata != path {
		t.Errorf("Expected testdata path to be %s, got %s", path, mockRpc.testdata)
	}
}
func TestGetBootstrap(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mock_rpc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	mockBootstrap := BootstrapResponse{
		Data: consensus_core.Bootstrap{
			Header: consensus_core.Header{
				Slot: 1000,
			},
		},
	}
	bootstrapJSON, _ := json.Marshal(mockBootstrap)
	err = os.WriteFile(filepath.Join(tempDir, "bootstrap.json"), bootstrapJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock bootstrap file: %v", err)
	}
	mockRpc := NewMockRpc(tempDir)
	bootstrap, err := mockRpc.GetBootstrap([32]byte{})
	if err != nil {
		t.Fatalf("GetBootstrap failed: %v", err)
	}
	if bootstrap.Header.Slot != 1000 {
		t.Errorf("Expected bootstrap slot to be 1000, got %d", bootstrap.Header.Slot)
	}
}
func TestGetUpdates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mock_rpc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	mockUpdates := UpdateResponse{
		{Data: consensus_core.Update{SignatureSlot: 1}},
		{Data: consensus_core.Update{SignatureSlot: 2}},
	}
	updatesJSON, _ := json.Marshal(mockUpdates)
	err = os.WriteFile(filepath.Join(tempDir, "updates.json"), updatesJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock updates file: %v", err)
	}
	mockRpc := NewMockRpc(tempDir)
	updates, err := mockRpc.GetUpdates(1, 2)
	if err != nil {
		t.Fatalf("GetUpdates failed: %v", err)
	}
	if len(updates) != 2 {
		t.Errorf("Expected 2 updates, got %d", len(updates))
	}
	if updates[0].SignatureSlot != 1 || updates[1].SignatureSlot != 2 {
		t.Errorf("Unexpected update signature slots")
	}
}
func TestGetFinalityUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mock_rpc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	mockFinality := FinalityUpdateResponse{
		Data: consensus_core.FinalityUpdate{
			FinalizedHeader: consensus_core.Header{
				Slot: 2000,
			},
		},
	}
	finalityJSON, _ := json.Marshal(mockFinality)
	err = os.WriteFile(filepath.Join(tempDir, "finality.json"), finalityJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock finality file: %v", err)
	}
	mockRpc := NewMockRpc(tempDir)
	finality, err := mockRpc.GetFinalityUpdate()
	if err != nil {
		t.Fatalf("GetFinalityUpdate failed: %v", err)
	}
	if finality.FinalizedHeader.Slot != 2000 {
		t.Errorf("Expected finality slot to be 2000, got %d", finality.FinalizedHeader.Slot)
	}
}
func TestGetOptimisticUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mock_rpc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	mockOptimistic := OptimisticUpdateResponse{
		Data: consensus_core.OptimisticUpdate{
			SignatureSlot: 3000,
		},
	}
	optimisticJSON, _ := json.Marshal(mockOptimistic)
	err = os.WriteFile(filepath.Join(tempDir, "optimistic.json"), optimisticJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock optimistic file: %v", err)
	}
	mockRpc := NewMockRpc(tempDir)
	optimistic, err := mockRpc.GetOptimisticUpdate()
	if err != nil {
		t.Fatalf("GetOptimisticUpdate failed: %v", err)
	}
	if optimistic.SignatureSlot != 3000 {
		t.Errorf("Expected optimistic signature slot to be 3000, got %d", optimistic.SignatureSlot)
	}
}
func TestGetBlock(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mock_rpc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	blocksDir := filepath.Join(tempDir, "blocks")
	err = os.Mkdir(blocksDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create blocks directory: %v", err)
	}
	mockBlock := BeaconBlockResponse{
		Data: struct {
			Message consensus_core.BeaconBlock
		}{
			Message: consensus_core.BeaconBlock{
				Slot: 4000,
			},
		},
	}
	blockJSON, _ := json.Marshal(mockBlock)
	err = os.WriteFile(filepath.Join(blocksDir, "4000.json"), blockJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write mock block file: %v", err)
	}
	mockRpc := NewMockRpc(tempDir)
	block, err := mockRpc.GetBlock(4000)
	if err != nil {
		t.Fatalf("GetBlock failed: %v", err)
	}
	if block.Slot != 4000 {
		t.Errorf("Expected block slot to be 4000, got %d", block.Slot)
	}
}
func TestChainId(t *testing.T) {
	mockRpc := NewMockRpc("/tmp/testdata")
	_, err := mockRpc.ChainId()
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("Expected 'not implemented' error, got %v", err)
	}
}
