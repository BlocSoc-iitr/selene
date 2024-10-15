package config
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
)
func TestNetwork_BaseConfig(t *testing.T) {
	tests := []struct {
		name            string
		inputNetwork    string
		expectedChainID uint64
		expectedGenesis uint64
		expectedRPCPort uint16
		checkConsensusRPC func(*testing.T, *string)
		wantErr         bool
	}{
		{
			name:            "Mainnet",
			inputNetwork:    "MAINNET",
			expectedChainID: 1,
			expectedGenesis: 1606824023,
			expectedRPCPort: 8545,
			checkConsensusRPC: func(t *testing.T, rpc *string) {
				assert.NotNil(t, rpc)
				assert.Equal(t, "https://www.lightclientdata.org", *rpc)
			},
			wantErr:         false,
		},
		{
			name:            "Goerli",
			inputNetwork:    "GOERLI",
			expectedChainID: 5,
			expectedGenesis: 1616508000,
			expectedRPCPort: 8545,
			checkConsensusRPC: func(t *testing.T, rpc *string) {
				assert.Nil(t, rpc)
			},
			wantErr:         false,
		},
		{
			name:            "Sepolia",
			inputNetwork:    "SEPOLIA",
			expectedChainID: 11155111,
			expectedGenesis: 1655733600,
			expectedRPCPort: 8545,
			checkConsensusRPC: func(t *testing.T, rpc *string) {
				assert.Nil(t, rpc)
			},
			wantErr:         false,
		},
		{
			name:            "Invalid",
			inputNetwork:    "INVALID",
			expectedChainID: 0,
			expectedGenesis: 0,
			expectedRPCPort: 0,
			checkConsensusRPC: func(t *testing.T, rpc *string) {},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Network("") // The receiver doesn't matter, we're testing the input
			config, err := n.BaseConfig(tt.inputNetwork)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedChainID, config.Chain.ChainID)
				assert.Equal(t, tt.expectedGenesis, config.Chain.GenesisTime)
				assert.Equal(t, tt.expectedRPCPort, config.RpcPort)
				tt.checkConsensusRPC(t, config.ConsensusRpc)
				
				// Check Forks
				assert.NotEmpty(t, config.Forks.Genesis)
				assert.NotEmpty(t, config.Forks.Altair)
				assert.NotEmpty(t, config.Forks.Bellatrix)
				assert.NotEmpty(t, config.Forks.Capella)
				assert.NotEmpty(t, config.Forks.Deneb)
				
				// Check MaxCheckpointAge
				assert.Equal(t, uint64(1_209_600), config.MaxCheckpointAge)
				
				// Check DataDir
				assert.NotNil(t, config.DataDir)
				assert.Contains(t, strings.ToLower(*config.DataDir), strings.ToLower(tt.inputNetwork))
			}
		})
	}
}
func TestNetwork_ChainID(t *testing.T) {
	tests := []struct {
		name            string
		inputChainID    uint64
		expectedChainID uint64
		expectedGenesis uint64
		expectedRPCPort uint16
		checkConsensusRPC func(*testing.T, *string)
		wantErr         bool
	}{
		{
			name:            "Mainnet",
			inputChainID:    1,
			expectedChainID: 1,
			expectedGenesis: 1606824023,
			expectedRPCPort: 8545,
			checkConsensusRPC: func(t *testing.T, rpc *string) {
				assert.NotNil(t, rpc)
				assert.Equal(t, "https://www.lightclientdata.org", *rpc)
			},
			wantErr:         false,
		},
		{
			name:            "Goerli",
			inputChainID:    5,
			expectedChainID: 5,
			expectedGenesis: 1616508000,
			expectedRPCPort: 8545,
			checkConsensusRPC: func(t *testing.T, rpc *string) {
				assert.Nil(t, rpc)
			},
			wantErr:         false,
		},
		{
			name:            "Sepolia",
			inputChainID:    11155111,
			expectedChainID: 11155111,
			expectedGenesis: 1655733600,
			expectedRPCPort: 8545,
			checkConsensusRPC: func(t *testing.T, rpc *string) {
				assert.Nil(t, rpc)
			},
			wantErr:         false,
		},
		{
			name:            "Invalid ChainID",
			inputChainID:    9999, // Non-existent ChainID
			expectedChainID: 0,
			expectedGenesis: 0,
			expectedRPCPort: 0,
			checkConsensusRPC: func(t *testing.T, rpc *string) {},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Network("") // The receiver doesn't matter, we're testing the input
			config, err := n.ChainID(tt.inputChainID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedChainID, config.Chain.ChainID)
				assert.Equal(t, tt.expectedGenesis, config.Chain.GenesisTime)
				assert.Equal(t, tt.expectedRPCPort, config.RpcPort)
				tt.checkConsensusRPC(t, config.ConsensusRpc)
				
				// Check Forks
				assert.NotEmpty(t, config.Forks.Genesis)
				assert.NotEmpty(t, config.Forks.Altair)
				assert.NotEmpty(t, config.Forks.Bellatrix)
				assert.NotEmpty(t, config.Forks.Capella)
				assert.NotEmpty(t, config.Forks.Deneb)
				
				// Check MaxCheckpointAge
				assert.Equal(t, uint64(1_209_600), config.MaxCheckpointAge)
				
				// Check DataDir
				assert.NotNil(t, config.DataDir)
				assert.Contains(t, strings.ToLower(*config.DataDir), strings.ToLower(tt.name))
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
