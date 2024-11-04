package rpc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetBootstrap(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mock_rpc_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mockBootstrap := map[string]interface{}{
		"data": map[string]interface{}{
			"header": map[string]interface{}{
				"beacon": map[string]interface{}{
					"slot":           "1000",
					"proposer_index": "12345",
					"parent_root":    "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
					"state_root":     "0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40",
					"body_root":      "0x4142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f60",
				},
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

	mockUpdates := []map[string]interface{}{
		{
			"data": map[string]interface{}{
				"attested_header": map[string]interface{}{
					"beacon": map[string]interface{}{
						"slot":           "1",
						"proposer_index": "12345",
						"parent_root":    "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
						"state_root":     "0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40",
						"body_root":      "0x4142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f60",
					},
				},
				"signature_slot": "1",
			},
		},
		{
			"data": map[string]interface{}{
				"attested_header": map[string]interface{}{
					"beacon": map[string]interface{}{
						"slot": "2",
					},
				},
				"signature_slot": "2",
			},
		},
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

	mockFinality := map[string]interface{}{
		"data": map[string]interface{}{
			"attested_header": map[string]interface{}{
				"beacon": map[string]interface{}{
					"slot": "1000",
				},
			},
			"finalized_header": map[string]interface{}{
				"beacon": map[string]interface{}{
					"slot": "2000",
				},
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

	mockOptimistic := map[string]interface{}{
		"data": map[string]interface{}{
			"attested_header": map[string]interface{}{
				"beacon": map[string]interface{}{
					"slot": "3000",
				},
			},
			"signature_slot": "3000",
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

	mockBlock := map[string]interface{}{
		"data": map[string]interface{}{
			"message": map[string]interface{}{
				"slot": "4000",
			},
		},
	}

	blockJSON, err := json.Marshal(mockBlock)
	if err != nil {
		t.Fatalf("Failed to marshal mock block: %v", err)
	}

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
