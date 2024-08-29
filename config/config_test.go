package config

func main() {
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
	
	config.from_file(&path, &network, &cliConfig)
}

func TestMainnetBaseConfig() {
	network := "MAINNET"
	path := ""
	
	var cliConfig CliConfig
	var config Config

	config.from_file(&path, &network, &cliConfig)

	// config should have BaseConfig values for MAINNET as cliConfig and TomlConfig are uninitialised
	if config.Chain !=  
}