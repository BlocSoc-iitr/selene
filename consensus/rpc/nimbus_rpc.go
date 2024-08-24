package rpc

// uses types package
type NimbusRpc struct {
	ConsensusRpc
	rpc string
}

func get() {}

func (n NimbusRpc) new()                   {}
func (n NimbusRpc) get_bootstrap()         {}
func (n NimbusRpc) get_updates()           {}
func (n NimbusRpc) get_finality_update()   {}
func (n NimbusRpc) get_optimistic_update() {}
func (n NimbusRpc) get_block()             {}
func (n NimbusRpc) chain_id()              {}

type BeaconBlockResponse struct{}
type BeaconBlockData struct{}
type UpdateResponse = []UpdateData
type UpdateData struct{}
type FinalityUpdateResponse struct{}
type OptimisticUpdateResponse struct{}
type SpecResponse struct{}
type Spec struct{}
type BootstrapResponse struct{}
