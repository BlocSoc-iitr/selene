package config

type ChainConfig struct{

}
type Fork struct{
	epoch uint64
	fork_version  []uint8
}
type Forks struct{
	genesis Fork
	altair Fork
    bellatrix Fork
    capella Fork
    deneb Fork
}