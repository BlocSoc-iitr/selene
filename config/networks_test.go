package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBaseConfig(t *testing.T) {
	tests := []struct {
		name          string
		networkString string
		expectedError bool
		expectedChainID uint64
	}{
		{"Mainnet", "MAINNET", false, 1},
		{"Goerli", "GOERLI", false, 5},
		{"Sepolia", "SEPOLIA", false, 11155111},
		{"Unknown Network", "INVALID", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network := Network(tt.networkString)
			config, err := network.BaseConfig(tt.networkString)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedChainID, config.Chain.ChainID)
			}
		})
	}
}

func TestChainID(t *testing.T) {
	tests := []struct {
		name          string
		chainID       uint64
		expectedError bool
		expectedChainID uint64
	}{
		{"Mainnet ChainID", 1, false, 1},
		{"Goerli ChainID", 5, false, 5},
		{"Sepolia ChainID", 11155111, false, 11155111},
		{"Unknown ChainID", 9999, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network := Network("")
			config, err := network.ChainID(tt.chainID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedChainID, config.Chain.ChainID)
			}
		})
	}
}

func TestDataDir(t *testing.T) {
	tests := []struct {
		name     string
		network  Network
		expected string
	}{
		{"Mainnet DataDir", MAINNET, "selene/data/mainnet"},
		{"Goerli DataDir", GOERLI, "selene/data/goerli"},
		{"Sepolia DataDir", SEPOLIA, "selene/data/sepolia"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := dataDir(tt.network)
			assert.NoError(t, err)
			assert.Contains(t, path, tt.expected)
		})
	}
}
