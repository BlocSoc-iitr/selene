package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/spf13/viper"
)

var (
	executionRpc         = "http://localhost:8545"
	consensusRpc         = "http://localhost:5052"
	checkpoint           = "0x85e6151a246e8fdba36db27a0c7678a575346272fe978c9281e13a8b26cdfa68"
	rpcBindIp            = "127.0.0.1"
	rpcPort              = uint16(8080)
	dataDirectory        = "/data"
	fallback             = "http://fallback.example.com"
	loadExternalFallback = true
	maxCheckpointAge     = 86400
	strictCheckpointAge  = false
	defaultCheckpoint    = [32]byte{}
)

// ///////////////////////////
// /// FromFile() tests /////
// ///////////////////////////
func TestMainnetBaseConfig(t *testing.T) {
	network := "MAINNET"
	path := "./config.toml"
	v := viper.New()
	// Set default values for the configuration
	v.SetDefault("execution_rpc", executionRpc)
	createConfigFile(v)
	var cliConfig CliConfig
	var config Config

	config = config.FromFile(&path, &network, &cliConfig)

	mainnetConfig, _ := Mainnet()

	// config should have BaseConfig values for MAINNET as cliConfig and TomlConfig are uninitialised
	if !reflect.DeepEqual(config.Chain, mainnetConfig.Chain) {
		t.Errorf("Expected Chain to be %v, but got %v", mainnetConfig.Chain, config.Chain)
	}
	if !reflect.DeepEqual(config.Forks, mainnetConfig.Forks) {
		t.Errorf("Expected Forks to be %v, but got %v", mainnetConfig.Forks, config.Forks)
	}
	if *config.RpcBindIp != mainnetConfig.RpcBindIp {
		t.Errorf("Expected RpcBindIP to be %s, but got %s", mainnetConfig.RpcBindIp, *config.RpcBindIp)
	}
	if *config.RpcPort != mainnetConfig.RpcPort {
		t.Errorf("Expected RpcPort to be %v, but got %v", mainnetConfig.RpcPort, *config.RpcPort)
	}
	if config.ConsensusRpc != *mainnetConfig.ConsensusRpc {
		t.Errorf("Expected ConsensusRpc to be %s, but got %s", *mainnetConfig.ConsensusRpc, config.ConsensusRpc)
	}
}
func TestConfigFileCreatedSuccessfully(t *testing.T) {
	network := "MAINNET"
	path := "./config.toml"
	v := viper.New()
	// Set default values for the configuration
	v.SetDefault("consensus_rpc", consensusRpc)
	v.SetDefault("execution_rpc", executionRpc)
	v.SetDefault("rpc_bind_ip", rpcBindIp)
	v.SetDefault("rpc_port", rpcPort)
	v.SetDefault("checkpoint", checkpoint)
	v.SetDefault("data_dir", dataDirectory)
	v.SetDefault("fallback", fallback)
	v.SetDefault("load_external_fallback", loadExternalFallback)
	v.SetDefault("strict_checkpoint_age", maxCheckpointAge)
	createConfigFile(v)
	var cliConfig CliConfig
	var config Config

	config = config.FromFile(&path, &network, &cliConfig)

	if config.ConsensusRpc != consensusRpc {
		t.Errorf("Expected ConsensusRpc to be %s, but got %s", consensusRpc, config.ConsensusRpc)
	}
	if config.ExecutionRpc != executionRpc {
		t.Errorf("Expected executionRpc to be %s, but got %s", executionRpc, config.ExecutionRpc)
	}
	if *config.DataDir != dataDirectory {
		t.Errorf("Expected data directory to be %s, but got %s", dataDirectory, *config.DataDir)
	}
	if config.LoadExternalFallback != loadExternalFallback {
		t.Errorf("Expected load external fallback to be %v, but got %v", loadExternalFallback, config.LoadExternalFallback)
	}
	if hex.EncodeToString((*config.Checkpoint)[:]) != checkpoint[2:] {
		t.Errorf("Expected checkpoint to be %s, but got %s", checkpoint[2:], hex.EncodeToString((*config.Checkpoint)[:]))
	}

}
func TestCliConfig(t *testing.T) {
	network := "MAINNET"
	path := "./config.toml"
	v := viper.New()
	// Set default values for the configuration
	v.SetDefault("execution_rpc", executionRpc)
	cliConfig := CliConfig{
		ExecutionRpc: &executionRpc,
		ConsensusRpc: &consensusRpc,
		RpcBindIp:    &rpcBindIp,
	}
	var config Config

	config = config.FromFile(&path, &network, &cliConfig)

	if config.ExecutionRpc != *cliConfig.ExecutionRpc {
		t.Errorf("Expected Execution rpc to be %s, but got %s", *cliConfig.ExecutionRpc, config.ExecutionRpc)
	}
	if config.ConsensusRpc != *cliConfig.ConsensusRpc {
		t.Errorf("Expected ConsensusRpc to be %s, but got %s", consensusRpc, config.ConsensusRpc)
	}
	if *config.RpcBindIp != *cliConfig.RpcBindIp {
		t.Errorf("Expected rpc bind ip to be %s, but got %s", *cliConfig.RpcBindIp, *config.RpcBindIp)
	}
}
func TestIfFieldNotInCliDefaultsToTomlThenBaseConfig(t *testing.T) {
	network := "MAINNET"
	path := "./config.toml"
	v := viper.New()
	// Set default values for the configuration
	v.SetDefault("consensus_rpc", consensusRpc)
	v.SetDefault("execution_rpc", executionRpc)
	v.SetDefault("rpc_bind_ip", rpcBindIp)
	v.SetDefault("rpc_port", rpcPort)
	v.SetDefault("checkpoint", checkpoint)
	v.SetDefault("data_dir", dataDirectory)
	v.SetDefault("fallback", fallback)
	v.SetDefault("load_external_fallback", loadExternalFallback)
	v.SetDefault("max_checkpoint_age", maxCheckpointAge)
	// Create file
	createConfigFile(v)
	cliConfig := CliConfig{
		ConsensusRpc: &consensusRpc,
		RpcBindIp:    &rpcBindIp,
	}
	var config Config
	config = config.FromFile(&path, &network, &cliConfig)
	mainnetConfig, _ := Mainnet()

	// Rpc Port defined in toml file not in cli config
	if config.ExecutionRpc != executionRpc {
		t.Errorf("Expected executionRpc to be %s, but got %s", executionRpc, config.ExecutionRpc)
	}
	// Rpc Port defined in toml file not in cli config
	if *config.RpcPort != rpcPort {
		t.Errorf("Expected rpc port to be %v, but got %v", rpcPort, config.RpcPort)
	}
	// Chain not defined in either toml or config
	if !reflect.DeepEqual(config.Chain, mainnetConfig.Chain) {
		t.Errorf("Expected Chain to be %v, but got %v", mainnetConfig.Chain, config.Chain)
	}
}
func createConfigFile(v *viper.Viper) {
	// Specify the configuration file name and type
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	// Write configuration to file
	configFile := "./config.toml"
	if err := v.WriteConfigAs(configFile); err != nil {
		fmt.Printf("Error creating config file: %v\n", err)

		// Create the file if it doesn't exist
		if os.IsNotExist(err) {
			if err := v.WriteConfigAs(configFile); err != nil {
				fmt.Printf("Error creating config file: %v\n", err)
			} else {
				fmt.Println("Config file created successfully.")
			}
		} else {
			fmt.Printf("Failed to write config file: %v\n", err)
		}
	}
}

// ////////////////////////////////
// /// to_base_config() tests /////
// ////////////////////////////////
func TestReturnsCorrectBaseConfig(t *testing.T) {
	config := Config{
		ConsensusRpc:         consensusRpc,
		RpcBindIp:            &rpcBindIp,
		RpcPort:              &rpcPort,
		DefaultCheckpoint:    defaultCheckpoint,
		Chain:                ChainConfig{},
		Forks:                consensus_core.Forks{},
		MaxCheckpointAge:     uint64(maxCheckpointAge),
		DataDir:              &dataDirectory,
		LoadExternalFallback: loadExternalFallback,
		StrictCheckpointAge:  strictCheckpointAge,
	}

	baseConfig := config.ToBaseConfig()

	if !reflect.DeepEqual(baseConfig.Chain, config.Chain) {
		t.Errorf("Expected Chain to be %v, got %v", config.Chain, baseConfig.Chain)
	}
	if !reflect.DeepEqual(baseConfig.Forks, config.Forks) {
		t.Errorf("Expected Forks to be %v, got %v", config.Forks, baseConfig.Forks)
	}
	if baseConfig.MaxCheckpointAge != config.MaxCheckpointAge {
		t.Errorf("Expected Max Checkpoint age to be %v, got %v", config.MaxCheckpointAge, baseConfig.MaxCheckpointAge)
	}
}
func TestReturnsCorrectDefaultValues(t *testing.T) {
	config := Config{
		ConsensusRpc:         consensusRpc,
		DefaultCheckpoint:    defaultCheckpoint,
		Chain:                ChainConfig{},
		Forks:                consensus_core.Forks{},
		MaxCheckpointAge:     uint64(maxCheckpointAge),
		DataDir:              &dataDirectory,
		LoadExternalFallback: loadExternalFallback,
		StrictCheckpointAge:  strictCheckpointAge,
	}

	baseConfig := config.ToBaseConfig()
	if baseConfig.MaxCheckpointAge != config.MaxCheckpointAge {
		t.Errorf("Expected Max Checkpoint age to be %v, got %v", config.MaxCheckpointAge, baseConfig.MaxCheckpointAge)
	}
	if baseConfig.RpcPort != 8545 {
		t.Errorf("Expected rpc Port to be %v, got %v", 8545, baseConfig.RpcPort)
	}
	if baseConfig.RpcBindIp != "127.0.0.1" {
		t.Errorf("Expected Max Checkpoint age to be %v, got %v", "127.0.0.1", baseConfig.RpcBindIp)
	}
}
