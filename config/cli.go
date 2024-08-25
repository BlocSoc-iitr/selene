package config

import (
    "encoding/hex"
)
// The format of configuration to be stored in the configuratin file is map[string]interface{}
type CliConfig struct {
    ExecutionRPC         *string `mapstructure:"execution_rpc"`
    ConsensusRPC         *string `mapstructure:"consensus_rpc"`
    Checkpoint           *[]byte `mapstructure:"checkpoint"`
    RPCBindIP            *string `mapstructure:"rpc_bind_ip"`
    RPCPort              *uint16 `mapstructure:"rpc_port"`
    DataDir              *string `mapstructure:"data_dir"`
    Fallback             *string `mapstructure:"fallback"`
    LoadExternalFallback *bool   `mapstructure:"load_external_fallback"`
    StrictCheckpointAge  *bool   `mapstructure:"strict_checkpoint_age"`
}
func (cfg *CliConfig) as_provider() map[string]interface{} {
    // Create a map to hold the configuration data
    userDict := make(map[string]interface{})
    // Populate the map with values from the CliConfig struct
    if cfg.ExecutionRPC != nil {
        userDict["execution_rpc"] = *cfg.ExecutionRPC
    }
    if cfg.ConsensusRPC != nil {
        userDict["consensus_rpc"] = *cfg.ConsensusRPC
    }
    if cfg.Checkpoint != nil {
        userDict["checkpoint"] = hex.EncodeToString(*cfg.Checkpoint)
    }
    if cfg.RPCBindIP != nil {
        userDict["rpc_bind_ip"] = *cfg.RPCBindIP
    }
    if cfg.RPCPort != nil {
        userDict["rpc_port"] = *cfg.RPCPort
    }
    if cfg.DataDir != nil {
        userDict["data_dir"] = *cfg.DataDir
    }
    if cfg.Fallback != nil {
        userDict["fallback"] = *cfg.Fallback
    }
    if cfg.LoadExternalFallback != nil {
        userDict["load_external_fallback"] = *cfg.LoadExternalFallback
    }
    if cfg.StrictCheckpointAge != nil {
        userDict["strict_checkpoint_age"] = *cfg.StrictCheckpointAge
    }
    return userDict
}
