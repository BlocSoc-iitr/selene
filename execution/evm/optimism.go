//go:build optimism_default_handler && !negate_optimism_default_handler
// +build optimism_default_handler,!negate_optimism_default_handler

package evm

func getDefaultOptimismSetting() bool {
    return true
}