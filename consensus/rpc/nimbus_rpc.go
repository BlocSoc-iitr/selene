package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"io"
	"net/http"
	"strconv"
	"time"
)

// uses types package
const (
	MAX_REQUEST_LIGHT_CLIENT_UPDATES uint8 = 128
	initialBackoff                         = 50 * time.Millisecond // Initial backoff of 50ms
	maxRetries                             = 5                     // Number of retry attempts
)

func get[R any](req string, result *R) error {
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(req) //sends get request to the url
		if err != nil {            // error occurs
			time.Sleep(initialBackoff)
			continue //moves to next iteration
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
		return nil // no error occurs and result has been populated with the JSON data.
	}
	return fmt.Errorf("failed to fetch data from %s after %d retries", req, maxRetries)
}
func min(a uint8, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}
type NimbusRpc struct {
	//ConsensusRpc
	rpc string
}
func NewNimbusRpc(rpc string) *NimbusRpc {
	return &NimbusRpc{
		rpc: rpc}
}
func (n *NimbusRpc) GetBootstrap(block_root []byte) (consensus_core.Bootstrap, error) {
	root_hex := fmt.Sprintf("%x", block_root)
	req := fmt.Sprintf("%s/eth/v1/beacon/light_client/bootstrap/0x%s", n.rpc, root_hex)
	var res BootstrapResponse
	err := get(req, &res)
	if err != nil {
		return consensus_core.Bootstrap{}, fmt.Errorf("bootstrap error: %w", err)
	}
	return res.Data, nil
}
func (n *NimbusRpc) GetUpdates(period uint64, count uint8) ([]consensus_core.Update, error) {
	count = min(count, MAX_REQUEST_LIGHT_CLIENT_UPDATES)
	req := fmt.Sprintf("%s/eth/v1/beacon/light_client/updates?start_period=%s&count=%s", n.rpc, strconv.FormatUint(period, 10), strconv.FormatUint(uint64(count), 10))
	var res UpdateResponse
	err := get(req, &res)
	if err != nil {
		return nil, fmt.Errorf("updates error: %w", err)
	}
	updates := make([]consensus_core.Update, len(res))
	for i, update := range res {
		updates[i] = update.Data
	}
	return updates, nil
}
func (n *NimbusRpc) GetFinalityUpdate() (consensus_core.FinalityUpdate, error) {
	req := fmt.Sprintf("%s/eth/v1/beacon/light_client/finality_update", n.rpc)
	var res FinalityUpdateResponse
	err := get(req, &res)
	if err != nil {
		return consensus_core.FinalityUpdate{}, fmt.Errorf("finality update error: %w", err)
	}
	return res.Data, nil
}
func (n *NimbusRpc) GetOptimisticUpdate() (consensus_core.OptimisticUpdate, error) {
	req := fmt.Sprintf("%s/eth/v1/beacon/light_client/optimistic_update", n.rpc)
	var res OptimisticUpdateResponse
	err := get(req, &res)
	if err != nil {
		return consensus_core.OptimisticUpdate{}, fmt.Errorf("finality update error: %w", err)
	}
	return res.Data, nil
}
func (n *NimbusRpc) GetBlock(slot uint64) (consensus_core.BeaconBlock, error) {
	req := fmt.Sprintf("%s/eth/v2/beacon/blocks/%s", n.rpc, strconv.FormatUint(slot, 10))
	var res BeaconBlockResponse
	err := get(req, &res)
	if err != nil {
		return consensus_core.BeaconBlock{}, fmt.Errorf("block error: %w", err)
	}
	return res.Data.Message, nil
}
func (n *NimbusRpc) ChainId() (uint64, error) {
	req := fmt.Sprintf("%s/eth/v1/config/spec", n.rpc)
	var res SpecResponse
	err := get(req, &res)
	if err != nil {
		return 0, fmt.Errorf("spec error: %w", err)
	}
	return res.Data.ChainId, nil
}
// BeaconBlock, Update,FinalityUpdate ,OptimisticUpdate,Bootstrap yet to be defined in consensus-core/src/types/mod.go
// For now defined in consensus/consensus_core.go
type BeaconBlockResponse struct {
	Data BeaconBlockData
}
type BeaconBlockData struct {
	Message consensus_core.BeaconBlock
}
type UpdateResponse = []UpdateData
type UpdateData struct {
	Data consensus_core.Update
}
type FinalityUpdateResponse struct {
	Data consensus_core.FinalityUpdate
}
type OptimisticUpdateResponse struct {
	Data consensus_core.OptimisticUpdate
}
type SpecResponse struct {
	Data Spec
}
type Spec struct {
	ChainId uint64
}
type BootstrapResponse struct {
	Data consensus_core.Bootstrap
}

