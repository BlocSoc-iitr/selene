package config

type Network string

const (
	MAINNET Network = "MAINNET"
	GOERLI  Network = "GOERLI"
	SEPOLIA Network = "SEPOLIA"
)

// hardcode the base configurations for each of the networks
// check s and write switch case
func (n Network) base_config(s string) (BaseConfig, error) {

}

func (n Network) chain_id(s string) (uint64, error) {

}
func data_dir() {}
