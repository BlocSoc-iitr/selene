package config

import (
	"encoding/hex"
	"encoding/json"
)

type ChainConfig struct {
	ChainID     uint64 `json:"chain_id"`
	GenesisTime uint64 `json:"genesis_time"`
	GenesisRoot []byte `json:"genesis_root"`
}

func (c ChainConfig) MarshalJSON() ([]byte, error) {
	type Alias ChainConfig
	return json.Marshal(&struct {
		Alias
		GenesisRoot string `json:"genesis_root"`
	}{
		Alias:       (Alias)(c),
		GenesisRoot: hex.EncodeToString(c.GenesisRoot),
	})
}
func (c *ChainConfig) UnmarshalJSON(data []byte) error {
	type Alias ChainConfig
	aux := &struct {
		*Alias
		GenesisRoot string `json:"genesis_root"`
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	c.GenesisRoot, err = hex.DecodeString(aux.GenesisRoot)
	return err
}
