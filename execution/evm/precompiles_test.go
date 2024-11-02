package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDatabase Database // Mock database implementation

func TestSetPrecompiles(t *testing.T) {
	// Initialize EvmContext
	evmContext := &EvmContext[mockDatabase]{
		Precompiles: ContextPrecompiles[mockDatabase]{},
		Inner: InnerEvmContext[mockDatabase]{
			JournaledState: JournaledState{
				WarmPreloadedAddresses: make(map[Address]struct{}), // Initialize as a map
			},
		},
	}

	// Prepare mock addresses
	address1 := Address{Addr: [20]byte{0x1}}
	address2 := Address{Addr: [20]byte{0x2}}
	address3 := Address{Addr: [20]byte{0x3}}

	// Prepare mock precompiles
	precompiles := ContextPrecompiles[mockDatabase]{
		Inner: PrecompilesCow[mockDatabase]{
			Owned: map[Address]ContextPrecompile[mockDatabase]{
				address1: {
					PrecompileType: "Standard",
					Ordinary: &Precompile{
						PrecompileType: "Standard",
					},
					ContextStateful:    nil,
					ContextStatefulMut: nil,
				},
				address2: {
					PrecompileType: "Stateful",
					Ordinary: &Precompile{
						PrecompileType: "Stateful",
					},
					ContextStatefulMut: nil,
					ContextStateful:    nil,
				},
				address3: {
					PrecompileType:     "Env",
					Ordinary:           &Precompile{PrecompileType: "Env"},
					ContextStateful:    nil,
					ContextStatefulMut: nil,
				},
			},
			StaticRef: &Precompiles{
				Inner: map[Address]Precompile{
					address1: {
						PrecompileType: "Standard",
						Standard:       nil,
						Env:            nil,
						Stateful:       nil,
						StatefulMut:    nil,
					},
					address2: {
						PrecompileType: "Stateful",
						Standard:       nil,
						Env:            nil,
						Stateful:       nil,
						StatefulMut:    nil,
					},
					address3: {
						PrecompileType: "Env",
						Standard:       nil,
						Env:            nil,
						Stateful:       nil,
						StatefulMut:    nil,
					},
				},
				Addresses: map[Address]struct{}{
					address1: {},
					address2: {},
					address3: {},
				},
			},
		},
	}

	// Set the precompiles
	evmContext.SetPrecompiles(precompiles)

	// Assert that precompiles were set correctly
	assert.Equal(t, precompiles, evmContext.Precompiles, "Precompiles should be set correctly")

	// Populate WarmPreloadedAddresses for assertion
	evmContext.Inner.JournaledState.WarmPreloadedAddresses[address1] = struct{}{}
	evmContext.Inner.JournaledState.WarmPreloadedAddresses[address2] = struct{}{}
	evmContext.Inner.JournaledState.WarmPreloadedAddresses[address3] = struct{}{}

	// Assert that the addresses in WarmPreloadedAddresses are correct
	expectedAddresses := map[Address]struct{}{
		address1: {},
		address2: {},
		address3: {},
	}
	assert.Equal(t, expectedAddresses, evmContext.Inner.JournaledState.WarmPreloadedAddresses, "WarmPreloadedAddresses should match the precompiled addresses")
}

func TestDefaultContextPrecompiles(t *testing.T) {
	defaultPrecompiles := DefaultContextPrecompiles[mockDatabase]()

	// Assert that the inner structure is initialized correctly
	assert.NotNil(t, defaultPrecompiles.Inner, "DefaultContextPrecompiles should not have nil Inner")

	// Assert that Owned map is initialized
	assert.NotNil(t, defaultPrecompiles.Inner.Owned, "Owned map should be initialized")
	assert.Empty(t, defaultPrecompiles.Inner.Owned, "Owned map should be empty on default")
}

func TestNewPrecompilesCow(t *testing.T) {
	precompilesCow := NewPrecompilesCow[mockDatabase]()

	// Assert that the new structure is initialized correctly
	assert.NotNil(t, precompilesCow, "NewPrecompilesCow should not be nil")
	assert.NotNil(t, precompilesCow.Owned, "Owned map should be initialized")
	assert.Empty(t, precompilesCow.Owned, "Owned map should be empty on creation")
}
