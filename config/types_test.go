package config

import (
	"encoding/json"
	"testing"
)

func TestChainConfigMarshalUnmarshal(t *testing.T) {
	originalConfig := ChainConfig{
		ChainID:     1,
		GenesisTime: 1606824023,
		GenesisRoot: []byte{0x4b, 0x36, 0x3d, 0xb9},
	}

	// Marshal ChainConfig to JSON
	marshaledData, err := json.Marshal(originalConfig)
	if err != nil {
		t.Fatalf("Error marshaling ChainConfig: %v", err)
	}

	// Unmarshal the JSON back to ChainConfig
	var unmarshaledConfig ChainConfig
	err = json.Unmarshal(marshaledData, &unmarshaledConfig)
	if err != nil {
		t.Fatalf("Error unmarshaling ChainConfig: %v", err)
	}

	// Verify that the original and unmarshaled configs are the same
	if originalConfig.ChainID != unmarshaledConfig.ChainID {
		t.Errorf("ChainID mismatch. Got %d, expected %d", unmarshaledConfig.ChainID, originalConfig.ChainID)
	}
	if originalConfig.GenesisTime != unmarshaledConfig.GenesisTime {
		t.Errorf("GenesisTime mismatch. Got %d, expected %d", unmarshaledConfig.GenesisTime, originalConfig.GenesisTime)
	}
	if string(originalConfig.GenesisRoot) != string(unmarshaledConfig.GenesisRoot) {
		t.Errorf("GenesisRoot mismatch. Got %x, expected %x", unmarshaledConfig.GenesisRoot, originalConfig.GenesisRoot)
	}
}

func TestForkMarshalUnmarshal(t *testing.T) {
	originalFork := Fork{
		Epoch:       0,
		ForkVersion: []byte{0x01, 0x00, 0x00, 0x00},
	}

	// Marshal Fork to JSON
	marshaledData, err := json.Marshal(originalFork)
	if err != nil {
		t.Fatalf("Error marshaling Fork: %v", err)
	}

	// Unmarshal the JSON back to Fork
	var unmarshaledFork Fork
	err = json.Unmarshal(marshaledData, &unmarshaledFork)
	if err != nil {
		t.Fatalf("Error unmarshaling Fork: %v", err)
	}

	// Verify that the original and unmarshaled Fork are the same
	if originalFork.Epoch != unmarshaledFork.Epoch {
		t.Errorf("Epoch mismatch. Got %d, expected %d", unmarshaledFork.Epoch, originalFork.Epoch)
	}
	if string(originalFork.ForkVersion) != string(unmarshaledFork.ForkVersion) {
		t.Errorf("ForkVersion mismatch. Got %x, expected %x", unmarshaledFork.ForkVersion, originalFork.ForkVersion)
	}
}

func TestUnmarshalInvalidHex(t *testing.T) {
	invalidJSON := `{
		"epoch": 0,
		"fork_version": "invalid_hex_string"
	}`

	var fork Fork
	err := json.Unmarshal([]byte(invalidJSON), &fork)
	if err == nil {
		t.Fatal("Expected error unmarshaling invalid hex string, but got nil")
	}
}
