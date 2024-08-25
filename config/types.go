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
type Fork struct {
	Epoch       uint64 `json:"epoch"`
	ForkVersion []byte `json:"fork_version"`
}
type Forks struct {
	Genesis   Fork `json:"genesis"`
	Altair    Fork `json:"altair"`
	Bellatrix Fork `json:"bellatrix"`
	Capella   Fork `json:"capella"`
	Deneb     Fork `json:"deneb"`
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
func (f Fork) MarshalJSON() ([]byte, error) {
	type Alias Fork
	return json.Marshal(&struct {
		Alias
		ForkVersion string `json:"fork_version"`
	}{
		Alias:       (Alias)(f),
		ForkVersion: hex.EncodeToString(f.ForkVersion),
	})
}
func (f *Fork) UnmarshalJSON(data []byte) error {
	type Alias Fork
	aux := &struct {
		*Alias
		ForkVersion string `json:"fork_version"`
	}{
		Alias: (*Alias)(f),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	f.ForkVersion, err = hex.DecodeString(aux.ForkVersion)
	return err
}
