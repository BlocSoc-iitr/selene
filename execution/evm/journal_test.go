package evm

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
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

type logUnmarshalTestCase struct {
	name     string
	jsonData string
	expected Log[LogData]
}

// Helper function to decode hex strings to a fixed byte slice of given length.
func decodeHexString(t *testing.T, hexStr string, expectedLen int) []byte {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("Failed to decode hex string %s: %v", hexStr, err)
	}
	if len(bytes) != expectedLen {
		t.Fatalf("Decoded hex string %s does not match expected length %d; got %d", hexStr, expectedLen, len(bytes))
	}
	return bytes
}

// TestLogUnmarshal tests JSON unmarshalling for Log[LogData] using hex-encoded strings.
func TestLogUnmarshal(t *testing.T) {
	// Define test cases
	testCases := []logUnmarshalTestCase{
		{
			name: "Valid Log with two topics",
			jsonData: `{
				"address": "0x1234567890abcdef1234567890abcdef12345678",
				"data": {
					"topics": ["0x0000000000000000000000000000000012345678000000000000000000000000", "0x000000000000000000000000000000009abcdef0000000000000000000000000"],
					"data": "0x01020304"
				}
			}`,
			expected: Log[LogData]{
				Address: Address{Addr: [20]byte(decodeHexString(t, "1234567890abcdef1234567890abcdef12345678", 20))},
				Data: LogData{
					Topics: []B256{
						B256(decodeHexString(t, "0000000000000000000000000000000012345678000000000000000000000000", 32)),
						B256(decodeHexString(t, "000000000000000000000000000000009abcdef0000000000000000000000000", 32)),
					},
					Data: Bytes(decodeHexString(t, "01020304", 4)),
				},
			},
		},
		// Additional test cases can be added here
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal JSON data into Log[LogData]
			var logInstance Log[LogData]
			if err := json.Unmarshal([]byte(tc.jsonData), &logInstance); err != nil {
				t.Fatalf("Failed to unmarshal Log[LogData]: %v", err)
			}

			// Compare the unmarshalled result with the expected result
			if !reflect.DeepEqual(logInstance, tc.expected) {
				t.Errorf("Test %s failed.\nExpected: %+v\nGot: %+v", tc.name, tc.expected, logInstance)
			}
		})
	}
}
