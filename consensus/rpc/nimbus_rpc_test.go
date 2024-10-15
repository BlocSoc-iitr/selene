package rpc
import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
func TestNewNimbusRpc(t *testing.T) {
	rpcURL := "http://example.com"
	nimbusRpc := NewNimbusRpc(rpcURL)
	assert.Equal(t, rpcURL, nimbusRpc.rpc)
}
func TestNimbusGetBootstrap(t *testing.T) {
	blockRoot := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	expectedPath := fmt.Sprintf("/eth/v1/beacon/light_client/bootstrap/0x%x", blockRoot)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedPath, r.URL.Path)
		response := BootstrapResponse{
			Data: consensus_core.Bootstrap{
				Header: consensus_core.Header{
					Slot: 1000,
				},
			},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()
	nimbusRpc := NewNimbusRpc(server.URL)
	bootstrap, err := nimbusRpc.GetBootstrap(blockRoot)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000), bootstrap.Header.Slot)
}
func TestNimbusGetUpdates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v1/beacon/light_client/updates", r.URL.Path)
		assert.Equal(t, "start_period=1000&count=5", r.URL.RawQuery)
		response := UpdateResponse{
			{Data: consensus_core.Update{AttestedHeader: consensus_core.Header{Slot: 1000}}},
			{Data: consensus_core.Update{AttestedHeader: consensus_core.Header{Slot: 1001}}},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()
	nimbusRpc := NewNimbusRpc(server.URL)
	updates, err := nimbusRpc.GetUpdates(1000, 5)
	assert.NoError(t, err)
	assert.Len(t, updates, 2)
	assert.Equal(t, uint64(1000), updates[0].AttestedHeader.Slot)
	assert.Equal(t, uint64(1001), updates[1].AttestedHeader.Slot)
}
func TestNimbusGetFinalityUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v1/beacon/light_client/finality_update", r.URL.Path)
		response := FinalityUpdateResponse{
			Data: consensus_core.FinalityUpdate{
				FinalizedHeader: consensus_core.Header{
					Slot: 2000,
				},
			},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()
	nimbusRpc := NewNimbusRpc(server.URL)
	finalityUpdate, err := nimbusRpc.GetFinalityUpdate()
	assert.NoError(t, err)
	assert.Equal(t, uint64(2000), finalityUpdate.FinalizedHeader.Slot)
}
func TestNimbusGetOptimisticUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v1/beacon/light_client/optimistic_update", r.URL.Path)
		response := OptimisticUpdateResponse{
			Data: consensus_core.OptimisticUpdate{
				SignatureSlot: 3000,
			},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()
	nimbusRpc := NewNimbusRpc(server.URL)
	optimisticUpdate, err := nimbusRpc.GetOptimisticUpdate()
	assert.NoError(t, err)
	assert.Equal(t, uint64(3000), optimisticUpdate.SignatureSlot)
}
func TestNimbusGetBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v2/beacon/blocks/4000", r.URL.Path)
		response := BeaconBlockResponse{
			Data: BeaconBlockData{
				Message: consensus_core.BeaconBlock{
					Slot: 4000,
				},
			},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()
	nimbusRpc := NewNimbusRpc(server.URL)
	block, err := nimbusRpc.GetBlock(4000)
	assert.NoError(t, err)
	assert.Equal(t, uint64(4000), block.Slot)
}
func TestNimbusChainId(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v1/config/spec", r.URL.Path)
		response := SpecResponse{
			Data: Spec{
				ChainId: 5000,
			},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()
	nimbusRpc := NewNimbusRpc(server.URL)
	chainId, err := nimbusRpc.ChainId()
	assert.NoError(t, err)
	assert.Equal(t, uint64(5000), chainId)
}