package main

import (
	"fmt"
	"github.com/BlocSoc-iitr/selene/config"
	"github.com/BlocSoc-iitr/selene/consensus"
	Common "github.com/ethereum/go-ethereum/common"
)

func main() {
	var n config.Network
	baseConfig, err := n.BaseConfig("MAINNET")
	if err != nil {
		fmt.Printf("Error in base config creation: %v", err)
	}

	checkpoint := "b21924031f38635d45297d68e7b7a408d40b194d435b25eeccad41c522841bd5"
	consensusRpcUrl := "http://testing.mainnet.beacon-api.nimbus.team"
	executionRpcUrl := "https://eth-mainnet.g.alchemy.com/v2/j28GcevSYukh-GvSeBOYcwHOfIggF1Gt"

	config := &config.Config{
		ConsensusRpc:         consensusRpcUrl,
		ExecutionRpc:         executionRpcUrl,
		RpcBindIp:            &baseConfig.RpcBindIp,
		RpcPort:              &baseConfig.RpcPort,
		DefaultCheckpoint:    baseConfig.DefaultCheckpoint,
		Checkpoint:           (*[32]byte)(Common.Hex2Bytes(checkpoint)),
		Chain:                baseConfig.Chain,
		Forks:                baseConfig.Forks,
		StrictCheckpointAge:  false,
		DataDir:              baseConfig.DataDir,
		DatabaseType:         nil,
		MaxCheckpointAge:     baseConfig.MaxCheckpointAge,
		LoadExternalFallback: false,
	}
	consensusClient := consensus.ConsensusClient{}.New(&consensusRpcUrl, *config)

	err = consensusClient.Shutdown()
	if err != nil {
		fmt.Println("Found error in shutting down: ", err)
	}
}