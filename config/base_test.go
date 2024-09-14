package config

import (
	"testing"
)
func TestCorrectDefaultBaseConfig(t *testing.T) {
	baseConfig := BaseConfig{}

	baseConfig = baseConfig.Default()

	if baseConfig.RpcBindIp != "127.0.0.1" {
		t.Errorf("Expected RpcBindIP to be %s, but got %s", "127.0.0.1", baseConfig.RpcBindIp)
	}
	if baseConfig.RpcPort != 0 {
		t.Errorf("Expected RpcPort to be %v, but got %v", 0, baseConfig.RpcPort)
	}
	if baseConfig.ConsensusRpc != nil {
		t.Errorf("Expected ConsensusRpc to be %v, but got %v", nil, baseConfig.ConsensusRpc)
	}
}