package evm

import (
	"testing"
	"math/big"
	"github.com/stretchr/testify/assert"
	// "github.com/BlocSoc-iitr/selene/common"
)
// TestInnerEvmContextWithJournalState ensures InnerEvmContext integrates correctly with JournaledState.
func TestInnerEvmContextWithJournalState(t *testing.T) {
	db := NewEmptyDB()
	ctx := NewInnerEvmContext(db)

	assert.NotNil(t, ctx.JournaledState, "JournaledState should be initialized")
	assert.Nil(t, ctx.JournaledState.State, "State in JournaledState should be nil on initialization")
	assert.Equal(t, ctx.JournaledState.Spec, DefaultSpecId(), "Spec in JournaledState should match DefaultSpecId")
}

// TestTransientStorageInitialization ensures TransientStorage is initialized correctly in JournaledState.
func TestTransientStorageInitialization(t *testing.T) {
	spec := DefaultSpecId()
	journalState := NewJournalState(spec, nil)
	
	// Initialize a TransientStorage entry
	key := Key{
		Account: Address{Addr:[20]byte{0x1}},
		Slot:    U256(big.NewInt(0)),
	}
	value := U256(big.NewInt(12345))
	journalState.TransientStorage = make(TransientStorage)
	journalState.TransientStorage[key] = value
	
	assert.Equal(t, value, journalState.TransientStorage[key], "TransientStorage should store the correct value for the key")
}

// TestLogInitialization verifies that Logs are initialized and added to JournaledState.
func TestLogInitialization(t *testing.T) {
	spec := DefaultSpecId()
	journalState := NewJournalState(spec, nil)

	log := Log[LogData]{
		Address: Address{Addr:[20]byte{0x1}},
		Data: LogData{
			Topics: []B256{{0x1, 0x2, 0x3, 0x4}},
			Data:   Bytes{0x10, 0x20, 0x30},
		},
	}
	journalState.Logs = append(journalState.Logs, log)

	assert.Len(t, journalState.Logs, 1, "Logs should contain one entry")
	assert.Equal(t, log, journalState.Logs[0], "Log entry should match the added log")
}
