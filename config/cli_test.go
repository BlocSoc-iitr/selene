package config

import (
    "encoding/hex"
    "testing"

    "github.com/spf13/viper"
)
// for the test I have used viper, but ant other configuration library can be used. Just the format of the configuration should be the same i.e. map[string]interface{}
func TestCliConfigAsProvider(t *testing.T) {
	// used some random values for the test
    executionRpc := "http://localhost:8545"
    consensusRpc := "http://localhost:5052"
    checkpoint := []byte{0x01, 0x02, 0x03}
    rpcBindIp := "127.0.0.1"
    rpcPort := uint16(8080)
    dataDir := "/data"
    fallback := "http://fallback.example.com"
    loadExternalFallback := true
    strictCheckpointAge := false

    config := CliConfig{
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
    configMap := config.as_provider()
    for key, value := range configMap {
        v.Set(key, value)
    }
    // Assertions to check that Viper has the correct values
    if v.GetString("execution_rpc") != executionRpc {
        t.Errorf("Expected execution_rpc to be %s, but got %s", executionRpc, v.GetString("execution_rpc"))
    }
    if v.GetString("consensus_rpc") != consensusRpc {
        t.Errorf("Expected consensus_rpc to be %s, but got %s", consensusRpc, v.GetString("consensus_rpc"))
    }
    if v.GetString("checkpoint") != hex.EncodeToString(checkpoint) {
        t.Errorf("Expected checkpoint to be %s, but got %s", hex.EncodeToString(checkpoint), v.GetString("checkpoint"))
    }
    if v.GetString("rpc_bind_ip") != rpcBindIp {
        t.Errorf("Expected rpc_bind_ip to be %s, but got %s", rpcBindIp, v.GetString("rpc_bind_ip"))
    }
    if v.GetUint16("rpc_port") != rpcPort {
        t.Errorf("Expected rpc_port to be %d, but got %d", rpcPort, v.GetUint16("rpc_port"))
    }
    if v.GetString("data_dir") != dataDir {
        t.Errorf("Expected data_dir to be %s, but got %s", dataDir, v.GetString("data_dir"))
    }
    if v.GetString("fallback") != fallback {
        t.Errorf("Expected fallback to be %s, but got %s", fallback, v.GetString("fallback"))
    }
    if v.GetBool("load_external_fallback") != loadExternalFallback {
        t.Errorf("Expected load_external_fallback to be %v, but got %v", loadExternalFallback, v.GetBool("load_external_fallback"))
    }
    if v.GetBool("strict_checkpoint_age") != strictCheckpointAge {
        t.Errorf("Expected strict_checkpoint_age to be %v, but got %v", strictCheckpointAge, v.GetBool("strict_checkpoint_age"))
    }
}
