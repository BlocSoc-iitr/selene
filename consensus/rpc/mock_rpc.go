package rpc

//if we need to test the working we can by adding the whole code of nimbus_rpc here
// and add a path to some local testdata folder
import (
	"encoding/json"
	"fmt"
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"os"
	"path/filepath"
)

type MockRpc struct {
	testdata string
}

func NewMockRpc(path string) *MockRpc {
	return &MockRpc{
		testdata: path,
	}
}
func (m *MockRpc) GetBootstrap(block_root [32]byte) (consensus_core.Bootstrap, error) {
	path := filepath.Join(m.testdata, "bootstrap.json")
	res, err := os.ReadFile(path)
	if err != nil {
		return consensus_core.Bootstrap{}, fmt.Errorf("failed to read file: %w", err)
	}
	var bootstrap BootstrapResponse
	err = json.Unmarshal(res, &bootstrap)
	if err != nil {
		return consensus_core.Bootstrap{}, fmt.Errorf("bootstrap error: %w", err)
	}
	return bootstrap.Data, nil
}
func (m *MockRpc) GetUpdates(period uint64, count uint8) ([]consensus_core.Update, error) {
	path := filepath.Join(m.testdata, "updates.json")
	res, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var updatesResponse UpdateResponse

	err = json.Unmarshal(res, &updatesResponse)
	if err != nil {
		return nil, fmt.Errorf("updates error: %w", err)
	}

	var updates = make([]consensus_core.Update, len(updatesResponse))
	for i, update := range updatesResponse {
		updates[i] = update.Data
	}

	//fmt.Print("Data", updates)
	return updates, nil
}
func (m *MockRpc) GetFinalityUpdate() (consensus_core.FinalityUpdate, error) {
	path := filepath.Join(m.testdata, "finality.json")
	res, err := os.ReadFile(path)
	if err != nil {
		return consensus_core.FinalityUpdate{}, fmt.Errorf("failed to read file: %w", err)
	}
	var finality FinalityUpdateResponse
	err = json.Unmarshal(res, &finality)
	if err != nil {
		return consensus_core.FinalityUpdate{}, fmt.Errorf("finality update error: %w", err)
	}
	return finality.Data, nil
}
func (m *MockRpc) GetOptimisticUpdate() (consensus_core.OptimisticUpdate, error) {
	path := filepath.Join(m.testdata, "optimistic.json")
	res, err := os.ReadFile(path)
	if err != nil {
		return consensus_core.OptimisticUpdate{}, fmt.Errorf("failed to read file: %w", err)
	}
	var optimistic OptimisticUpdateResponse
	err = json.Unmarshal(res, &optimistic)
	if err != nil {
		return consensus_core.OptimisticUpdate{}, fmt.Errorf("optimistic update error: %w", err)
	}
	return optimistic.Data, nil
}
func (m *MockRpc) GetBlock(slot uint64) (consensus_core.BeaconBlock, error) {
	path := filepath.Join(m.testdata, fmt.Sprintf("blocks/%d.json", slot))
	res, err := os.ReadFile(path)
	if err != nil {
		return consensus_core.BeaconBlock{}, fmt.Errorf("failed to read file: %w", err)
	}
	var block BeaconBlockResponse
	err = json.Unmarshal(res, &block)
	if err != nil {
		return consensus_core.BeaconBlock{}, err
	}
	return block.Data.Message, nil
}
func (m *MockRpc) ChainId() (uint64, error) {
	return 0, fmt.Errorf("not implemented")
}
