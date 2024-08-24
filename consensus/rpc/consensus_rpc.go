package rpc

// return types not mention and oarameters as well
type ConsensusRpc interface{
	new()
	get_bootstrap()
	get_updates()
	get_finality_update();
	get_optimistic_update();
	get_block()
	chain_id()
}    