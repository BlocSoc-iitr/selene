package evm

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the NewSpec function
func TestNewSpec(t *testing.T) {
	specID := SpecId(42)
	spec := NewSpec(specID)

	assert.Equal(t, specID, spec.SpecID(), "SpecID should match the initialized value")
}

// Test the Enabled method in BaseSpec
func TestBaseSpec_Enabled(t *testing.T) {
	spec1 := NewSpec(SpecId(50))
	spec2 := NewSpec(SpecId(30))

	assert.True(t, spec1.Enabled(SpecId(30)), "Spec1 should be enabled for SpecId 30")
	assert.False(t, spec2.Enabled(SpecId(50)), "Spec2 should not be enabled for SpecId 50")
}

// Test TryFromUint8 with valid and invalid values
func TestTryFromUint8(t *testing.T) {
	// Test valid values
	validIDs := []uint8{0, 10, 19, uint8(LATEST)} // 0, 50, 100 and 255 (LATEST) should be valid
	for _, id := range validIDs {
		specID, ok := TryFromUint8(id)
		assert.True(t, ok, "Expected valid SpecId for uint8 value %d", id)
		assert.Equal(t, SpecId(id), specID, "Expected SpecId %d for uint8 value %d", id, id)
	}

}

// Test the IsEnabledIn method
func TestIsEnabledIn(t *testing.T) {
	spec1 := SpecId(50)
	spec2 := SpecId(30)

	assert.True(t, spec1.IsEnabledIn(spec2), "Spec1 should be enabled in Spec2")
	assert.False(t, spec2.IsEnabledIn(spec1), "Spec2 should not be enabled in Spec1")
}

type optimismFieldsTestCase struct {
	name     string
	jsonData string
	expected OptimismFields
}

// TestOptimismFieldsUnmarshal tests JSON unmarshalling for OptimismFields.
func TestOptimismFieldsUnmarshal(t *testing.T) {
	// Define test cases
	mintValue := uint64(1000)
	isSystemTx := true
	testCases := []optimismFieldsTestCase{
		{
			name: "Full OptimismFields with all fields present",
			jsonData: `{
				"source_hash": "0x0000000000000000000000000000000012345678000000000000000000000000",
				"mint": 1000,
				"is_system_transaction": true,
				"enveloped_tx": "0x01020304"
			}`,
			expected: OptimismFields{
				SourceHash:          (*B256)(decodeHexString(t, "0000000000000000000000000000000012345678000000000000000000000000", 32)),
				Mint:                &mintValue,
				IsSystemTransaction: &isSystemTx,
				EnvelopedTx:         Bytes(decodeHexString(t, "01020304", 4)),
			},
		},
		{
			name: "OptimismFields with only EnvelopedTx field",
			jsonData: `{
				"enveloped_tx": "0x0a0b0c0d"
			}`,
			expected: OptimismFields{
				EnvelopedTx: Bytes(decodeHexString(t, "0a0b0c0d", 4)),
			},
		},
		// Additional test cases for other field combinations can be added here
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal JSON data into OptimismFields
			var optimismInstance OptimismFields
			if err := json.Unmarshal([]byte(tc.jsonData), &optimismInstance); err != nil {
				t.Fatalf("Failed to unmarshal OptimismFields: %v", err)
			}

			// Compare the unmarshalled result with the expected result
			if !reflect.DeepEqual(optimismInstance, tc.expected) {
				t.Errorf("Test %s failed.\nExpected: %+v\nGot: %+v", tc.name, tc.expected, optimismInstance)
			}
		})
	}
}
