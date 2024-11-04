package consensus

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/BlocSoc-iitr/selene/common"
	"github.com/BlocSoc-iitr/selene/config"
	"github.com/BlocSoc-iitr/selene/consensus/consensus_core"
	"github.com/BlocSoc-iitr/selene/utils"
)

func GetClient(strictCheckpointAge bool, sync bool) (*Inner, error) {
	var n config.Network
	baseConfig, err := n.BaseConfig("MAINNET")
	if err != nil {
		return nil, err
	}

	config := &config.Config{
		ConsensusRpc:        "",
		ExecutionRpc:        "",
		Chain:               baseConfig.Chain,
		Forks:               baseConfig.Forks,
		StrictCheckpointAge: strictCheckpointAge,
	}

	checkpoint := "5afc212a7924789b2bc86acad3ab3a6ffb1f6e97253ea50bee7f4f51422c9275"

	//Decode the hex string into a byte slice
	checkpointBytes, err := hex.DecodeString(checkpoint)
	checkpointBytes32 := [32]byte{}
	copy(checkpointBytes32[:], checkpointBytes)
	if err != nil {
		log.Fatalf("failed to decode checkpoint: %v", err)
	}

	blockSend := make(chan *common.Block, 256)
	finalizedBlockSend := make(chan *common.Block)
	channelSend := make(chan *[]byte)

	In := Inner{}
	client := In.New(
		"testdata/",
		blockSend,
		finalizedBlockSend,
		channelSend,
		config,
	)

	if sync {
		err := client.sync(checkpointBytes32)
		if err != nil {
			return nil, err
		}
	} else {
		client.bootstrap(checkpointBytes32)
	}

	return client, nil
}

// testVerifyUpdate runs the test and returns its result via a channel (no inputs required).
func TestVerifyUpdate(t *testing.T) {
	//Get the client
	client, err := GetClient(false, false)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	//Calculate the sync period
	period := utils.CalcSyncPeriod(client.Store.FinalizedHeader.Slot)

	//Fetch updates
	updates, err := client.RPC.GetUpdates(period, MAX_REQUEST_LIGHT_CLIENT_UPDATES)
	if err != nil {
		t.Fatalf("failed to get updates: %v", err)
	}

	//Ensure we have updates to verify
	if len(updates) == 0 {
		t.Fatalf("no updates fetched")
	}

	//Verify the first update
	update := updates[0]
	err = client.verify_update(&update)
	if err != nil {
		t.Fatalf("failed to verify update: %v", err)
	}

}

func TestVerifyUpdateInvalidCommittee(t *testing.T) {
	client, err := GetClient(false, false)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	period := utils.CalcSyncPeriod(client.Store.FinalizedHeader.Slot)
	updates, err := client.RPC.GetUpdates(period, MAX_REQUEST_LIGHT_CLIENT_UPDATES)
	if err != nil {
		t.Fatalf("failed to get updates: %v", err)
	}

	if len(updates) == 0 {
		t.Fatalf("no updates fetched")
	}

	update := updates[0]
	update.NextSyncCommittee.Pubkeys[0] = consensus_core.BLSPubKey{} // Invalid public key

	err = client.verify_update(&update)
	if err == nil || err.Error() != "invalid next sync committee proof" {
		t.Fatalf("expected 'invalid next sync committee proof', got %v", err)
	}
}

func TestVerifyUpdateInvalidFinality(t *testing.T) {
	client, err := GetClient(false, false)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	period := utils.CalcSyncPeriod(client.Store.FinalizedHeader.Slot)
	updates, err := client.RPC.GetUpdates(period, MAX_REQUEST_LIGHT_CLIENT_UPDATES)
	if err != nil {
		t.Fatalf("failed to get updates: %v", err)
	}

	if len(updates) == 0 {
		t.Fatalf("no updates fetched")
	}

	update := updates[0]
	update.FinalizedHeader = consensus_core.Header{} // Assuming an empty header is invalid

	err = client.verify_update(&update)
	if err == nil || err.Error() != "invalid finality proof" {
		t.Fatalf("expected 'invalid finality proof', got %v", err)
	}
}

func TestVerifyUpdateInvalidSig(t *testing.T) {
	client, err := GetClient(false, false)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	period := utils.CalcSyncPeriod(client.Store.FinalizedHeader.Slot)
	updates, err := client.RPC.GetUpdates(period, MAX_REQUEST_LIGHT_CLIENT_UPDATES)
	if err != nil {
		t.Fatalf("failed to get updates: %v", err)
	}

	if len(updates) == 0 {
		t.Fatalf("no updates fetched")
	}

	update := updates[0]
	update.SyncAggregate.SyncCommitteeSignature = consensus_core.SignatureBytes{} // Assuming an empty signature is invalid

	err = client.verify_update(&update)
	if err == nil || err.Error() != "invalid signature" {
		t.Fatalf("expected 'invalid signature', got %v", err)
	}
}

func TestVerifyFinality(t *testing.T) {
	//Get the client
	client, err := GetClient(false, true)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	//Fetch the finality update
	update, err := client.RPC.GetFinalityUpdate()
	if err != nil {
		t.Fatalf("failed to get finality update: %v", err)
	}

	//Verify the finality update
	err = client.verify_finality_update(&update)
	if err != nil {
		t.Fatalf("finality verification failed: %v", err)
	}
}

func TestVerifyFinalityInvalidFinality(t *testing.T) {
	//Get the client
	client, err := GetClient(false, true)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	//Fetch the finality update
	update, err := client.RPC.GetFinalityUpdate()
	if err != nil {
		t.Fatalf("failed to get finality update: %v", err)
	}

	//Modify the finalized header to be invalid
	update.FinalizedHeader = consensus_core.Header{} //Assuming an empty header is invalid

	//Verify the finality update and expect an error
	err = client.verify_finality_update(&update)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	//Check if the error matches the expected error message
	expectedErr := "invalid finality proof"
	if err.Error() != expectedErr {
		t.Errorf("expected %s, got %v", expectedErr, err)
	}
}

func TestVerifyFinalityInvalidSignature(t *testing.T) {
	//Get the client
	client, err := GetClient(false, true)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	//Fetch the finality update
	update, err := client.RPC.GetFinalityUpdate()
	if err != nil {
		t.Fatalf("failed to get finality update: %v", err)
	}

	//Modify the sync aggregate signature to be invalid
	update.SyncAggregate.SyncCommitteeSignature = consensus_core.SignatureBytes{} //Assuming an empty signature is invalid

	//Verify the finality update and expect an error
	err = client.verify_finality_update(&update)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	//Check if the error matches the expected error message
	expectedErr := "invalid signature"
	if err.Error() != expectedErr {
		t.Errorf("expected %s, got %v", expectedErr, err)
	}
}

func TestVerifyOptimistic(t *testing.T) {
	//Get the client
	client, err := GetClient(false, true)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	//Fetch the optimistic update
	update, err := client.RPC.GetOptimisticUpdate()
	if err != nil {
		t.Fatalf("failed to get optimistic update: %v", err)
	}

	//Verify the optimistic update
	err = client.verify_optimistic_update(&update)
	if err != nil {
		t.Fatalf("optimistic verification failed: %v", err)
	}
}
func TestVerifyOptimisticInvalidSignature(t *testing.T) {
	//Get the client
	client, err := GetClient(false, true)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	//Fetch the optimistic update
	update, err := client.RPC.GetOptimisticUpdate()
	if err != nil {
		t.Fatalf("failed to get optimistic update: %v", err)
	}

	//Modify the sync aggregate signature to be invalid
	update.SyncAggregate.SyncCommitteeSignature = consensus_core.SignatureBytes{} //Assuming an empty signature is invalid

	//Verify the optimistic update and expect an error
	err = client.verify_optimistic_update(&update)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	//Check if the error matches the expected error message
	expectedErr := "invalid signature"
	if err.Error() != expectedErr {
		t.Errorf("expected %s, got %v", expectedErr, err)
	}
}
func TestVerifyCheckpointAgeInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic due to invalid checkpoint age, but no panic occurred")
		} else {
			expectedPanicMessage := "checkpoint too old, consider using a more recent checkpoint"
			if msg, ok := r.(string); ok && msg != expectedPanicMessage {
				t.Errorf("expected panic message '%s', got '%s'", expectedPanicMessage, msg)
			}
		}
	}()

	// This should trigger a panic due to the invalid checkpoint age
	_, err := GetClient(true, false)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}
}

func TestSendBlocks(t *testing.T) {
	//Get the client
	client, err := GetClient(false, true)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	errAdvancongClient := client.advance()
	if errAdvancongClient != nil {
		t.Fatalf("failed to advance client: %v", errAdvancongClient)
	}
	//Send the blocks
	errSendingBlock := client.send_blocks()
	if errSendingBlock != nil {
		t.Fatalf("failed to send blocks: %v", errSendingBlock)
	}

}

func TestGetPayloads(t *testing.T) {
	//Get the client
	client, err := GetClient(false, true)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}

	//Fetch the payloads
	payloads, err := client.Get_payloads(7109344, 7109344)
	if err != nil {
		t.Fatalf("failed to get payloads: %v", err)
	}

	//Ensure we have payloads
	if len(payloads) == 0 {
		t.Fatalf("no payloads fetched")
	}
}
