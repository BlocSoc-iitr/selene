package config

import (
    "testing"
    "github.com/spf13/viper"
    "fmt"
    "os"
    "reflect"
)

func TestMain(m *testing.M) {
	var config Config
	
	executionRpc := "http://localhost:8545"
    consensusRpc := "http://localhost:5052"
    checkpoint := []byte{0x01, 0x02, 0x03}
    rpcBindIp := "127.0.0.1"
    rpcPort := uint16(8080)
    dataDir := "/data"
    fallback := "http://fallback.example.com"
    loadExternalFallback := true
    strictCheckpointAge := false

	path := ""
	network := "MAINNET"

	cliConfig := CliConfig{
        ExecutionRpc:         &executionRpc,
        ConsensusRpc:         &consensusRpc,
        Checkpoint:           &checkpoint,
        RpcBindIp:            &rpcBindIp,
        RpcPort:              &rpcPort,
        DataDir:              &dataDir,
        Fallback:             &fallback,
        LoadExternalFallback: &loadExternalFallback,
        StrictCheckpointAge:  &strictCheckpointAge,
    }
	
    v := viper.New()

	// Set default values for the configuration
	v.SetDefault("consensus_rpc", "http://localhost:8545")
	v.SetDefault("execution_rpc", "http://localhost:8546")
	v.SetDefault("rpc_bind_ip", "127.0.0.1")
	v.SetDefault("rpc_port", 8545)
	v.SetDefault("default_checkpoint", make([]byte, 32))
	v.SetDefault("data_dir", "./data")
	v.SetDefault("max_checkpoint_age", 1000)
	v.SetDefault("fallback", "fallback_value")
	v.SetDefault("load_external_fallback", false)
	v.SetDefault("strict_checkpoint_age", true)
	v.SetDefault("database_type", "sqlite")

    createConfigFile(v)

	// Specify the configuration file name and type
	v.SetConfigName("config") // Name of the file without extension
	v.SetConfigType("toml")   // Configuration file format (toml, json, yaml, etc.)
	v.AddConfigPath(".")      // Optionally, specify where to look for the config file (current directory)

    config.from_file(&path, &network, &cliConfig)

    code := m.Run()
    os.Exit(code)
}

func TestMainnetBaseConfig(t *testing.T) {
	network := "MAINNET"
	path := "./config.toml"

    v := viper.New()
    createConfigFile(v)
	
	var cliConfig CliConfig
	var config Config

	config.from_file(&path, &network, &cliConfig)

    mainnetConfig, _ := Mainnet()

	// config should have BaseConfig values for MAINNET as cliConfig and TomlConfig are uninitialised
	if !reflect.DeepEqual(config.Chain, mainnetConfig.Chain) {
        t.Errorf("Expected Chain to be %v, but got %v", mainnetConfig.Chain, config.Chain)
    }
    if !reflect.DeepEqual(config.Forks, mainnetConfig.Forks) {
        t.Errorf("Expected Chain to be %v, but got %v", mainnetConfig.Forks, config.Forks)
    }
    if !(*config.RpcBindIp).Equal(mainnetConfig.RpcBindIp) {
        t.Errorf("Expected Chain to be %s, but got %s", mainnetConfig.RpcBindIp.String(), (*config.RpcBindIp).String())
    }
    if *config.RpcPort != mainnetConfig.RpcPort {
        t.Errorf("Expected Chain to be %v, but got %v", mainnetConfig.RpcPort ,*config.RpcPort)
    }
}

func createConfigFile(v *viper.Viper) {
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
