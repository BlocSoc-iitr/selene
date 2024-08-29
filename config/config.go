package config

import (
	// "github.com/BlocSoc-iitr/selene/consensus_core"
	"net"
	"github.com/spf13/viper"
)
type Config struct {
	ConsensusRpc         string      `json:"consensus_rpc"`
	ExecutionRpc         string      `json:"execution_rpc"`
	RpcBindIp            *net.IP     `json:"rpc_bind_ip"`
	RpcPort              *uint16     `json:"rpc_port"`
	DefaultCheckpoint    [32]byte    `json:"default_checkpoint"`
	Checkpoint           *[32]byte   `json:"checkpoint"`
	DataDir              *string     `json:"data_dir"`
	Chain                ChainConfig `json:"chain"`
	Forks                Forks       `json:"forks"`
	MaxCheckpointAge     uint64      `json:"max_checkpoint_age"`
	Fallback             *string     `json:"fallback"`
	LoadExternalFallback bool        `json:"load_external_fallback"`
	StrictCheckpointAge  bool        `json:"strict_checkpoint_age"`
	DatabaseType         *string     `json:"database_type"`
}

// only if we are using CLI
func (c Config) from_file(configPath *string, network *string, cliConfig *CliConfig) Config {
	n := Network(*network)
	baseConfig, err := n.baseConfig(*network)
	if err != nil {
		baseConfig = BaseConfig{}.Default()
	}
	v := viper.New()
	v.SetConfigFile(*configPath)
	err = v.ReadInConfig()
	tomlProvider := v.AllSettings()     // returns the config as a map[string]any
	cliProvider := cliConfig.as_provider()
	finalConfig := Config{
		ConsensusRpc: *baseConfig.ConsensusRpc,
		RpcBindIp: &baseConfig.RpcBindIp,
		RpcPort: &baseConfig.RpcPort,
		DefaultCheckpoint: baseConfig.DefaultCheckpoint,
		Chain: baseConfig.Chain,
		Forks: baseConfig.Forks,
		MaxCheckpointAge: baseConfig.MaxCheckpointAge,
		DataDir: baseConfig.DataDir,
		LoadExternalFallback: baseConfig.LoadExternalFallback,
		StrictCheckpointAge: baseConfig.StrictCheckpointAge,
	}
	for key, value := range tomlProvider {
		v.Set(key, value)
	}
	for key, value := range cliProvider {
		v.Set(key, value)
	}
	err = v.Unmarshal(&finalConfig)
	return finalConfig
}	
func (c Config) to_base_config() BaseConfig {
	return BaseConfig{
		RpcBindIp: func() net.IP {
			if c.RpcBindIp != nil {
				return *c.RpcBindIp
			}
			return net.IPv4(127, 0, 0, 1)
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
// func (c Config) fork_version(slot uint64) []uint8 {
// 	// Importing this from consensus_core package. This function is written in consensus_core.go
// 	// Calculate_fork_version(c.Forks, slot)
// }