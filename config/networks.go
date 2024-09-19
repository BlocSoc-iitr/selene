package config
import (
	"fmt"
	"github.com/BlocSoc-iitr/selene/common"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)
type Network string
const (
	MAINNET Network = "MAINNET"
	GOERLI  Network = "GOERLI"
	SEPOLIA Network = "SEPOLIA"
)
func (n Network) baseConfig(s string) (BaseConfig, error) {
	switch strings.ToUpper(s) {
	case "MAINNET":
		config, err := Mainnet()
		if err != nil {
			return BaseConfig{}, err
		}
		return config, nil
	case "GOERLI":
		config, err := Goerli()
		if err != nil {
			return BaseConfig{}, err
		}
		return config, nil
	case "SEPOLIA":
		config, err := Sepolia()
		if err != nil {
			return BaseConfig{}, err
		}
		return config, nil
	default:
		return BaseConfig{}, errors.New("network not recognized")
	}
}
func (n Network) chainID (id uint64) (BaseConfig, error) {
	switch id {
	case 1:
		config, err := Mainnet()
		if err != nil {
			return BaseConfig{}, err
		}
		return config, nil
	case 5:
		config, err := Goerli()
		if err != nil {
			return BaseConfig{}, err
		}
		return config, nil
	case 11155111:
		
		config, err := Sepolia()
		if err != nil {
			return BaseConfig{}, err
		}
		return config, nil
	default:
		return BaseConfig{}, errors.New("ChainID not recognized")
	}
}
func dataDir(network Network) (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get user home directory: %w", err)
    }
    path := filepath.Join(homeDir, fmt.Sprintf("selene/data/%s", strings.ToLower(string(network))))
    return path, nil
}
func Mainnet() (BaseConfig, error) {
	defaultCheckpoint, err := common.Hex_str_to_bytes("c7fc7b2f4b548bfc9305fa80bc1865ddc6eea4557f0a80507af5dc34db7bd9ce")
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to parse default checkpoint: %w", err)
	}
	genesisRoot, err := common.Hex_str_to_bytes("4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95")
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to parse genesis root: %w", err)
	}
	consensusRPC := "https://www.lightclientdata.org"
	dataDir, err := dataDir(MAINNET)
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to get data directory: %w", err)
	}
	return BaseConfig{
		DefaultCheckpoint: [32]byte(defaultCheckpoint),
		RpcPort:           8545,
		ConsensusRpc:      &consensusRPC,
		Chain: ChainConfig{
			ChainID:     1,
			GenesisTime: 1606824023,
			GenesisRoot: genesisRoot,
		},
		Forks: Forks{
			Genesis: Fork{
				Epoch:       0,
				ForkVersion: []byte{0x00, 0x00, 0x00, 0x00}},
			Altair: Fork{
				Epoch:       74240,
				ForkVersion: []byte{0x01, 0x00, 0x00, 0x00}},
			Bellatrix: Fork{
				Epoch:       144896,
				ForkVersion: []byte{0x02, 0x00, 0x00, 0x00}},
			Capella: Fork{
				Epoch:       194048,
				ForkVersion: []byte{0x03, 0x00, 0x00, 0x00}},
			Deneb: Fork{
				Epoch:       269568,
				ForkVersion: []byte{0x04, 0x00, 0x00, 0x00}},
		},
		MaxCheckpointAge: 1_209_600, // 14 days
		DataDir:        &dataDir,
	}, nil
}
func Goerli() (BaseConfig, error) {
	defaultCheckpoint, err := common.Hex_str_to_bytes("f6e9d5fdd7c406834e888961beab07b2443b64703c36acc1274ae1ce8bb48839")
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to parse default checkpoint: %w", err)
	}
	genesisRoot, err := common.Hex_str_to_bytes("043db0d9a83813551ee2f33450d23797757d430911a9320530ad8a0eabc43efb")
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to parse genesis root: %w", err)
	}
	dataDir, err := dataDir(GOERLI)
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to get data directory: %w", err)
	}
	return BaseConfig{
		DefaultCheckpoint: [32]byte(defaultCheckpoint),
		RpcPort:           8545,
		ConsensusRpc:      nil,
		Chain: ChainConfig{
			ChainID:     5,
			GenesisTime: 1616508000,
			GenesisRoot: genesisRoot,
		},
		Forks: Forks{
			Genesis: Fork{
				Epoch:       0,
				ForkVersion: []byte{0x00, 0x10, 0x20, 0x00}},
			Altair: Fork{
				Epoch:       36660,
				ForkVersion: []byte{0x01, 0x10, 0x20, 0x00}},
			Bellatrix: Fork{
				Epoch:       112260,
				ForkVersion: []byte{0x02, 0x10, 0x20, 0x00}},
			Capella: Fork{
				Epoch:       162304,
				ForkVersion: []byte{0x03, 0x10, 0x20, 0x00}},
			Deneb: Fork{
				Epoch:       231680,
				ForkVersion: []byte{0x04, 0x10, 0x20, 0x00}},
		},
		MaxCheckpointAge: 1_209_600, // 14 days
		DataDir:          &dataDir,
	}, nil
}
func Sepolia() (BaseConfig, error) {
	defaultCheckpoint, err := common.Hex_str_to_bytes("4135bf01bddcfadac11143ba911f1c7f0772fdd6b87742b0bc229887bbf62b48")
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to parse default checkpoint: %w", err)
	}
	genesisRoot, err := common.Hex_str_to_bytes("d8ea171f3c94aea21ebc42a1ed61052acf3f9209c00e4efbaaddac09ed9b8078")
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to parse genesis root: %w", err)
	}
	dataDir, err := dataDir(SEPOLIA)
	if err != nil {
		return BaseConfig{}, fmt.Errorf("failed to get data directory: %w", err)
	}
	return BaseConfig{
		DefaultCheckpoint: [32]byte(defaultCheckpoint),
		RpcPort:           8545,
		ConsensusRpc:      nil,
		Chain: ChainConfig{
			ChainID:     11155111,
			GenesisTime: 1655733600,
			GenesisRoot: genesisRoot,
		},
		Forks: Forks{
			Genesis: Fork{
				Epoch:       0,
				ForkVersion: []byte{0x90, 0x00, 0x00, 0x69}},
			Altair: Fork{
				Epoch:       50,
				ForkVersion: []byte{0x90, 0x00, 0x00, 0x70}},
			Bellatrix: Fork{
				Epoch:       100,
				ForkVersion: []byte{0x90, 0x00, 0x00, 0x71}},
			Capella: Fork{
				Epoch:       56832,
				ForkVersion: []byte{0x90, 0x00, 0x00, 0x72}},
			Deneb: Fork{
				Epoch:       132608,
				ForkVersion: []byte{0x90, 0x00, 0x00, 0x73}},
		},
		MaxCheckpointAge: 1_209_600, // 14 days
		DataDir:          &dataDir,
	}, nil
}