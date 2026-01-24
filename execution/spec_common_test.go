package evm

import (
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
