package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"header": map[string]interface{}{
					"beacon": map[string]interface{}{
						"slot":           "1000",
						"proposer_index": "12345",
						"parent_root":    "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
						"state_root":     "0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40",
						"body_root":      "0x4142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f60",
					},
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
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Test panicked: %v", r)
		}
	}()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v1/beacon/light_client/updates", r.URL.Path)
		assert.Equal(t, "start_period=1000&count=5", r.URL.RawQuery)
		response := []map[string]interface{}{
			{
				"data": map[string]interface{}{
					"attested_header": map[string]interface{}{
						"beacon": map[string]interface{}{
							"slot":           "1",
							"proposer_index": "12345",
							"parent_root":    "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
							"state_root":     "0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40",
							"body_root":      "0x4142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f60",
						},
					},
					"signature_slot": "1",
				},
			},
			{
				"data": map[string]interface{}{
					"attested_header": map[string]interface{}{
						"beacon": map[string]interface{}{
							"slot": "2",
						},
					},
					"signature_slot": "2",
				},
			},
		}
		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	}))
	defer server.Close()
	nimbusRpc := NewNimbusRpc(server.URL)
	updates, err := nimbusRpc.GetUpdates(1000, 5)
	assert.NoError(t, err)
	assert.Len(t, updates, 2)
	assert.Equal(t, uint64(1), updates[0].AttestedHeader.Slot)
	assert.Equal(t, uint64(2), updates[1].AttestedHeader.Slot)
}
func TestNimbusGetFinalityUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v1/beacon/light_client/finality_update", r.URL.Path)
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"attested_header": map[string]interface{}{
					"beacon": map[string]interface{}{
						"slot": "2000",
					},
				},
				"finalized_header": map[string]interface{}{
					"beacon": map[string]interface{}{
						"slot": "2000",
					},
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
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"attested_header": map[string]interface{}{
					"beacon": map[string]interface{}{
						"slot": "3000",
					},
				},
				"signature_slot": "3000",
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
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"message": map[string]interface{}{
					"slot": "4000",
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
