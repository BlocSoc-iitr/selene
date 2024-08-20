package config

import (
    "encoding/json"
    "encoding/hex"
)
type ChainConfig struct{
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
func (c ChainConfig) Copy() ChainConfig{
    genesisrootcopy:= make([]byte,len(c.GenesisRoot))
    copy(genesisrootcopy, c.GenesisRoot)
    return ChainConfig{
        ChainID: c.ChainID,
        GenesisTime: c.GenesisTime,
        GenesisRoot: genesisrootcopy,
    }

}

func (f Fork) Copy() Fork {
    forkVersionCopy := make([]byte, len(f.ForkVersion))
    copy(forkVersionCopy, f.ForkVersion)
    return Fork{
        Epoch:       f.Epoch,
        ForkVersion: forkVersionCopy,
    }
}


func (f Forks) Copy() Forks {
    return Forks{
        Genesis:   f.Genesis.Copy(),
        Altair:    f.Altair.Copy(),
        Bellatrix: f.Bellatrix.Copy(),
        Capella:   f.Capella.Copy(),
        Deneb:     f.Deneb.Copy(),
    }
}