package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewPlainInstructionTable verifies that NewPlainInstructionTable initializes correctly.
func TestNewPlainInstructionTable(t *testing.T) {
	// Define a sample InstructionTable
	var sampleTable InstructionTable[any]

	// Initialize with NewPlainInstructionTable
	instructionTables := NewPlainInstructionTable(sampleTable)

	// Ensure PlainTable is set and BoxedTable is nil
	assert.NotNil(t, instructionTables.PlainTable)
	assert.Nil(t, instructionTables.BoxedTable)
	assert.Equal(t, PlainTableMode, instructionTables.Mode)
}

// TestNewBoxedInstructionTable verifies that NewBoxedInstructionTable initializes correctly.
func TestNewBoxedInstructionTable(t *testing.T) {
	// Define a sample BoxedInstructionTable
	var sampleBoxedTable BoxedInstructionTable[any]

	// Initialize with NewBoxedInstructionTable
	instructionTables := NewBoxedInstructionTable(sampleBoxedTable)

	// Ensure BoxedTable is set and PlainTable is nil
	assert.NotNil(t, instructionTables.BoxedTable)
	assert.Nil(t, instructionTables.PlainTable)
	assert.Equal(t, BoxedTableMode, instructionTables.Mode)
}

// TestInstructionTablesMode verifies the correct mode is set for each table type.
func TestInstructionTablesMode(t *testing.T) {
	// Test with PlainTable
	var plainTable InstructionTable[any]
	instructionTables := NewPlainInstructionTable(plainTable)
	assert.Equal(t, PlainTableMode, instructionTables.Mode)

	// Test with BoxedTable
	var boxedTable BoxedInstructionTable[any]
	instructionTables = NewBoxedInstructionTable(boxedTable)
	assert.Equal(t, BoxedTableMode, instructionTables.Mode)
}