package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJournalState(t *testing.T) {
	var adr Address = Address{
		Addr: [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78},
	}
	// Define test spec and preloaded addresses for initialization
	spec := SpecId(1)
	preloadedAddresses := map[Address]struct{}{
		adr: {},
	}

	// Initialize a new JournaledState
	journalState := NewJournalState(spec, preloadedAddresses)

	// Test assertions
	assert.Nil(t, journalState.State, "State should be nil on initialization")
	assert.Nil(t, journalState.TransientStorage, "TransientStorage should be nil on initialization")
	assert.Empty(t, journalState.Logs, "Logs should be empty on initialization")
	assert.Equal(t, uint(0), journalState.Depth, "Depth should be initialized to 0")
	assert.Empty(t, journalState.Journal, "Journal should be empty on initialization")
	assert.Equal(t, spec, journalState.Spec, "Spec ID should match the initialized value")
	assert.Equal(t, preloadedAddresses, journalState.WarmPreloadedAddresses, "WarmPreloadedAddresses should match the provided map")
}

func TestSetSpecId(t *testing.T) {
	// Define initial and new Spec IDs
	initialSpec := SpecId(1)
	newSpec := SpecId(2)

	// Initialize JournaledState and set Spec ID
	journalState := NewJournalState(initialSpec, nil)
	assert.Equal(t, initialSpec, journalState.Spec, "Initial Spec ID should match")

	// Call setSpecId to change the Spec ID
	journalState.setSpecId(newSpec)

	// Verify that Spec ID has been updated
	assert.Equal(t, newSpec, journalState.Spec, "Spec ID should be updated to new value")
}
