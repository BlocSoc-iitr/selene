package config

import (
	"fmt"

	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/BlocSoc-iitr/selene/utils"
	"github.com/spf13/viper"
)

type Config struct {
	ConsensusRpc         string               `json:"consensus_rpc"`
	ExecutionRpc         string               `json:"execution_rpc"`
	RpcBindIp            *string              `json:"rpc_bind_ip"`
	RpcPort              *uint16              `json:"rpc_port"`
	DefaultCheckpoint    [32]byte             `json:"default_checkpoint"` // In cli.go, checkpoint is currently taken as []byte{}
	Checkpoint           *[32]byte            `json:"checkpoint"`         // but it should be of 32 bytes or [32]byte{}
	DataDir              *string              `json:"data_dir"`
	Chain                ChainConfig          `json:"chain"`
	Forks                consensus_core.Forks `json:"forks"`
	MaxCheckpointAge     uint64               `json:"max_checkpoint_age"`
	Fallback             *string              `json:"fallback"`
	LoadExternalFallback bool                 `json:"load_external_fallback"`
	StrictCheckpointAge  bool                 `json:"strict_checkpoint_age"`
	DatabaseType         *string              `json:"database_type"`
}

// only if we are using CLI
func (c Config) FromFile(configPath *string, network *string, cliConfig *CliConfig) Config {
	n := Network(*network)
	baseConfig, err := n.BaseConfig(*network)
	if err != nil {
		baseConfig = BaseConfig{}.Default()
	}
	v := viper.New()
	v.SetConfigFile(*configPath)
	err = v.ReadInConfig()
	if err != nil {
		fmt.Printf("Error in reading config file: %v", err)
	}
	tomlProvider := v.AllSettings() // returns the config as a map[string]any
	finalConfig := Config{
		ConsensusRpc:         *baseConfig.ConsensusRpc,
		RpcBindIp:            &baseConfig.RpcBindIp,
		RpcPort:              &baseConfig.RpcPort,
		DefaultCheckpoint:    baseConfig.DefaultCheckpoint,
		Chain:                baseConfig.Chain,
		Forks:                baseConfig.Forks,
		MaxCheckpointAge:     baseConfig.MaxCheckpointAge,
		DataDir:              baseConfig.DataDir,
		LoadExternalFallback: baseConfig.LoadExternalFallback,
		StrictCheckpointAge:  baseConfig.StrictCheckpointAge,
	}
	if cliConfig.ConsensusRpc != nil && *cliConfig.ConsensusRpc != "" {
		finalConfig.ConsensusRpc = *cliConfig.ConsensusRpc
	} else if tomlProvider["consensus_rpc"] != nil && tomlProvider["consensus_rpc"].(string) != "" {
		finalConfig.ConsensusRpc = tomlProvider["consensus_rpc"].(string)
	}
	if cliConfig.ExecutionRpc != nil && *cliConfig.ExecutionRpc != "" {
		finalConfig.ExecutionRpc = *cliConfig.ExecutionRpc
	} else if tomlProvider["execution_rpc"] != nil && tomlProvider["execution_rpc"].(string) != "" {
		finalConfig.ExecutionRpc = tomlProvider["execution_rpc"].(string)
	} else {
		fmt.Println("Need an execution_rpc value but got none")
	}
	if cliConfig.RpcBindIp != nil && *cliConfig.RpcBindIp != "" {
		finalConfig.RpcBindIp = cliConfig.RpcBindIp
	} else if tomlProvider["rpc_bind_ip"] != nil && tomlProvider["rpc_bind_ip"].(string) != "" {
		rpcBindIp, _ := tomlProvider["rpc_bind_ip"].(string)
		finalConfig.RpcBindIp = &rpcBindIp
	}
	if cliConfig.RpcPort != nil && *cliConfig.RpcPort != 0 {
		finalConfig.RpcPort = cliConfig.RpcPort
	} else if tomlProvider["rpc_port"] != nil && uint16(tomlProvider["rpc_port"].(int64)) != 0 {
		rpcPort := uint16(tomlProvider["rpc_port"].(int64))
		finalConfig.RpcPort = &rpcPort
	}
	if cliConfig.Checkpoint != nil {
		finalConfig.Checkpoint = (*[32]byte)(*cliConfig.Checkpoint)
	} else if tomlProvider["checkpoint"] != nil && tomlProvider["checkpoint"].(string) != "" {
		checkpoint, _ := tomlProvider["checkpoint"].(string)
		checkpointBytes, err := utils.Hex_str_to_bytes(checkpoint)
		if err != nil {
			fmt.Printf("Failed to convert checkpoint value to byte slice, %v", err)
		}
		finalConfig.Checkpoint = (*[32]byte)(checkpointBytes)
	}
	if cliConfig.DataDir != nil && *cliConfig.DataDir != "" {
		finalConfig.DataDir = cliConfig.DataDir
	} else if tomlProvider["data_dir"] != nil && tomlProvider["data_dir"].(string) != "" {
		dataDir, _ := tomlProvider["data_dir"].(string)
		finalConfig.DataDir = &dataDir
	}
	if tomlProvider["max_checkpoint_age"] != nil && uint64(tomlProvider["max_checkpoint_age"].(int64)) != 0 {
		maxCheckpointAge, err := tomlProvider["max_checkpoint_age"].(int64)
		if !err {
			fmt.Printf("Failed to convert %s value to %s", "max_checkpoint_age", "uint64")
		}
		finalConfig.MaxCheckpointAge = uint64(maxCheckpointAge)
	}
	if cliConfig.Fallback != nil && *cliConfig.Fallback != "" {
		finalConfig.Fallback = cliConfig.Fallback
	} else if tomlProvider["fallback"] != nil && tomlProvider["fallback"].(string) != "" {
		fallback, _ := tomlProvider["fallback"].(string)
		finalConfig.Fallback = &fallback
	}
	if cliConfig.LoadExternalFallback != nil {
		finalConfig.LoadExternalFallback = *cliConfig.LoadExternalFallback
	} else if tomlProvider["load_external_fallback"] != nil {
		load_external_fallback, err := tomlProvider["load_external_fallback"].(bool)
		if !err {
			fmt.Printf("Failed to convert %s value to %s", "load_external_fallback", "bool")
		}
		finalConfig.LoadExternalFallback = load_external_fallback
	}
	if cliConfig.StrictCheckpointAge != nil {
		finalConfig.StrictCheckpointAge = *cliConfig.StrictCheckpointAge
	}
	return finalConfig
}
func (c Config) ToBaseConfig() BaseConfig {
	return BaseConfig{
		RpcBindIp: func() string {
			if c.RpcBindIp != nil {
				return *c.RpcBindIp
			}
			return "127.0.0.1"
		}(),
		RpcPort: func() uint16 {
			if c.RpcPort != nil {
				return *c.RpcPort
			}
			return 8545
		}(),
		ConsensusRpc:         &c.ConsensusRpc,
		DefaultCheckpoint:    c.DefaultCheckpoint,
		Chain:                c.Chain,
		Forks:                c.Forks,
		MaxCheckpointAge:     c.MaxCheckpointAge,
		DataDir:              c.DataDir,
		LoadExternalFallback: c.LoadExternalFallback,
		StrictCheckpointAge:  c.StrictCheckpointAge,
	}
}
