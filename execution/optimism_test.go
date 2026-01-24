//go:build optimism_default_handler && !negate_optimism_default_handler
// +build optimism_default_handler,!negate_optimism_default_handler

package evm

import "testing"

func TestShouldReturnTrue(t *testing.T) {
	// Test logic
	result := getDefaultOptimismSetting()
	if result != true {
		t.Fatalf("Test failed: expected true but got %v", result)
	}
}