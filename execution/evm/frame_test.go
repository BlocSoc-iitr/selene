package evm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallFrame(t *testing.T) {
    // Arrange
    expectedRange := Range{Start: 100, End: 200}
    expectedFrameData := FrameData{
        Checkpoint: JournalCheckpoint{
			Log_i: 10,
			Journal_i: 20,
		}, 
        Interpreter: Interpreter{
			InstructionPointer: func(b byte) *byte { return &b }(0x10),
			Gas: Gas{
				Limit: 100,
				Remaining: 50,
				Refunded: 50,
			},
			IsEof: true,
			IsEofInit: true,
			SharedMemory: SharedMemory{
				Buffer: []byte{0x1, 0x2, 0x3, 0x4, 0x5},
				Checkpoints: []int{0, 2, 4},
				LastCheckpoint: 4,
				MemoryLimit: 1024,
			},
			Stack: Stack{
				Data: []*big.Int{big.NewInt(15)},
			},
			IsStatic: true,
			Contract: Contract{
				Input: []byte{0xde, 0xad, 0xbe, 0xef},
				Bytecode: Bytecode{},
                Hash: &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
				TargetAddress: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
                Caller: Address{Addr: [20]byte{0x1b, 0x3d, 0x5f, 0x7a, 0x9c}},
                CallValue: U256(big.NewInt(1000)),
			},
			
			},              
		}
	
		// Act
		frame := CallFrame{
			ReturnMemoryRange: expectedRange,
			FrameData:         expectedFrameData,
		}
	
		assert.Equal(t, expectedRange, frame.ReturnMemoryRange)
		assert.Equal(t, expectedFrameData, frame.FrameData)
}


func TestCreateFrame(t *testing.T) {
    // Arrange
    expectedAddress :=  Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}}
    expectedFrameData := FrameData{
        Checkpoint: JournalCheckpoint{
			Log_i: 10,
			Journal_i: 20,
		}, // Initialize with your checkpoint data
        Interpreter: Interpreter{
			InstructionPointer: func(b byte) *byte { return &b }(0x10),
			Gas: Gas{
				Limit: 100,
				Remaining: 50,
				Refunded: 50,
			},
			IsEof: true,
			IsEofInit: true,
			SharedMemory: SharedMemory{
				Buffer: []byte{0x1, 0x2, 0x3, 0x4, 0x5},
				Checkpoints: []int{0, 2, 4},
				LastCheckpoint: 4,
				MemoryLimit: 1024,
			},
			Stack: Stack{
				Data: []*big.Int{big.NewInt(15)},
			},
			IsStatic: true,
			Contract: Contract{
				Input: []byte{0xde, 0xad, 0xbe, 0xef},
				Bytecode: Bytecode{},
                Hash: &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
				TargetAddress: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
                Caller: Address{Addr: [20]byte{0x1b, 0x3d, 0x5f, 0x7a, 0x9c}},
                CallValue: U256(big.NewInt(1000)),
			},
			
		},              
    }

    // Act
    frame := CreateFrame{
        CreatedAddress: expectedAddress,
        FrameData:      expectedFrameData,
    }

    // Assert
    assert.Equal(t, expectedAddress, frame.CreatedAddress)
    assert.Equal(t, expectedFrameData, frame.FrameData)
}
func TestEOFCreatedFrame(t *testing.T){
	expectedAddress :=  Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}}
    expectedFrameData := FrameData{
        Checkpoint: JournalCheckpoint{
			Log_i: 10,
			Journal_i: 20,
		}, // Initialize with your checkpoint data
        Interpreter: Interpreter{
			InstructionPointer: func(b byte) *byte { return &b }(0x10),
			Gas: Gas{
				Limit: 100,
				Remaining: 50,
				Refunded: 50,
			},
			IsEof: true,
			IsEofInit: true,
			SharedMemory: SharedMemory{
				Buffer: []byte{0x1, 0x2, 0x3, 0x4, 0x5},
				Checkpoints: []int{0, 2, 4},
				LastCheckpoint: 4,
				MemoryLimit: 1024,
			},
			Stack: Stack{
				Data: []*big.Int{big.NewInt(15)},
			},
			IsStatic: true,
			Contract: Contract{
				Input: []byte{0xde, 0xad, 0xbe, 0xef},
				Bytecode: Bytecode{},
                Hash: &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
				TargetAddress: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
                Caller: Address{Addr: [20]byte{0x1b, 0x3d, 0x5f, 0x7a, 0x9c}},
                CallValue: U256(big.NewInt(1000)),
			},
			
		},              
    }

    // Act
    frame := EOFCreateFrame{
        CreatedAddress: expectedAddress,
        FrameData:      expectedFrameData,
    }

    // Assert
    assert.Equal(t, expectedAddress, frame.CreatedAddress)
    assert.Equal(t, expectedFrameData, frame.FrameData)

}
func TestFrameData(t *testing.T) {
	expectedCheckpoint := JournalCheckpoint{
		Log_i: 10,
		Journal_i: 50,
	}
	expectedInterpreter := Interpreter{
		InstructionPointer: func(b byte) *byte { return &b }(0x10),
		Gas: Gas{
			Limit:     100,
			Remaining: 50,
			Refunded:  50,
		},
		IsEof:      true,
		IsEofInit:  true,
		SharedMemory: SharedMemory{
			Buffer:         []byte{0x1, 0x2, 0x3, 0x4, 0x5},
			Checkpoints:    []int{0, 2, 4},
			LastCheckpoint: 4,
			MemoryLimit:    1024,
		},
		Stack: Stack{
			Data: []*big.Int{big.NewInt(15)},
		},
		IsStatic: true,
		Contract: Contract{
			Input:         []byte{0xde, 0xad, 0xbe, 0xef},
			Bytecode:      Bytecode{},
			Hash:          &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
			TargetAddress: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
			Caller:        Address{Addr: [20]byte{0x1b, 0x3d, 0x5f, 0x7a, 0x9c}},
			CallValue:     U256(big.NewInt(1000)),
		},
	}

	// Act
	frame := FrameData{
		Checkpoint:  expectedCheckpoint,
		Interpreter: expectedInterpreter,
	}

	// Assert
	assert.Equal(t, expectedCheckpoint, frame.Checkpoint)
	assert.Equal(t, expectedInterpreter, frame.Interpreter)
}
	
func TestCallOutcome(t *testing.T){
	expectedResult :=InterpreterResult{
		Result: 8,
		Output: []byte{0xde, 0xad, 0xbe, 0xef},
		Gas: Gas{
			Limit: 100,
			Remaining: 50,
			Refunded: 50,
		},
	
		}
	expectedMemoryOffset :=   Range{Start: 100, End: 200}

	frame := CallOutcome{
		Result: expectedResult,
		MemoryOffset: expectedMemoryOffset,
	}

	// Use the frame variable to avoid the "declared and not used" error
	assert.Equal(t, expectedResult, frame.Result)
	assert.Equal(t, expectedMemoryOffset, frame.MemoryOffset)
}
func TestCreateOutcome(t *testing.T){
	expectedResult:=InterpreterResult{
		Result: 8,
		Output: []byte{0xde, 0xad, 0xbe, 0xef},
		Gas: Gas{
			Limit: 100,
			Remaining: 50,
			Refunded: 50,
		},
	
		}
	expectedAddress:=Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}}

	frame := CreateOutcome{
		Result:  expectedResult,
		Address: &expectedAddress,
	}

	// Use the frame variable to avoid the "declared and not used" error
	assert.Equal(t, expectedResult, frame.Result)
	assert.Equal(t, &expectedAddress, frame.Address)
	
}
func TestFrame(t *testing.T) {
    callFrame := &CallFrame{
        ReturnMemoryRange: Range{Start: 100, End: 200},
        FrameData: FrameData{
            Checkpoint: JournalCheckpoint{
                Log_i:    10,
                Journal_i: 20,
            },
            Interpreter: Interpreter{
                InstructionPointer: func(b byte) *byte { return &b }(0x10),
                Gas: Gas{
                    Limit:     100,
                    Remaining: 50,
                    Refunded:  50,
                },
                IsEof:      true,
                IsEofInit:  true,
                SharedMemory: SharedMemory{
                    Buffer:         []byte{0x1, 0x2, 0x3, 0x4, 0x5},
                    Checkpoints:    []int{0, 2, 4},
                    LastCheckpoint: 4,
                    MemoryLimit:    1024,
                },
                Stack: Stack{
                    Data: []*big.Int{big.NewInt(15)},
                },
                IsStatic: true,
                Contract: Contract{
                    Input:         []byte{0xde, 0xad, 0xbe, 0xef},
                    Bytecode:      Bytecode{},
                    Hash:          &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
                    TargetAddress: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
                    Caller:        Address{Addr: [20]byte{0x1b, 0x3d, 0x5f, 0x7a, 0x9c}},
                    CallValue:     U256(big.NewInt(1000)),
                },
            },
        },
    }

    frame := Frame{
        Type: Call_Frame,
        Call: callFrame,
    }

    // Assert
    assert.Equal(t, Call_Frame, frame.Type)
    assert.Equal(t, callFrame, frame.Call)
}

func TestFrameResult(t *testing.T) {
    // Arrange
	expectedResultType := FrameResultType(10)
	expectedCall := CallOutcome{
		InterpreterResult{
			Result: 8,
			Output: []byte{0xde, 0xad, 0xbe, 0xef},
			Gas: Gas{
				Limit: 100,
				Remaining: 50,
				Refunded: 50,
			},
		},
		Range{Start: 100, End: 200},
	}
	expectedCreateOutcome := CreateOutcome{
		Result: InterpreterResult{
			Result: 8,
			Output: []byte{0xde, 0xad, 0xbe, 0xef},
			Gas: Gas{
				Limit: 100,
				Remaining: 50,
				Refunded: 50,
			},
		},
		Address: &Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
	}
	frame := FrameResult{
		ResultType: expectedResultType,
		Call:       &expectedCall,
		Create:     &expectedCreateOutcome,
	}

	// Use the frame variable to avoid the "declared and not used" error
	assert.Equal(t, expectedResultType, frame.ResultType)
	assert.Equal(t, &expectedCall, frame.Call)
	assert.Equal(t, &expectedCreateOutcome, frame.Create)
		
}
	



func TestFrameOrResult(t *testing.T) {
    callFrame := Frame{
        Type: Call_Frame,
        Call: &CallFrame{
            ReturnMemoryRange: Range{Start: 100, End: 200},
            FrameData: FrameData{
                Checkpoint: JournalCheckpoint{
                    Log_i:    10,
                    Journal_i: 20,
                },
                Interpreter: Interpreter{
                    InstructionPointer: func(b byte) *byte { return &b }(0x10),
                    Gas: Gas{
                        Limit:     100,
                        Remaining: 50,
                        Refunded:  50,
                    },
                    IsEof: true,
                    Stack: Stack{
                        Data: []*big.Int{big.NewInt(15)},
                    },
                    Contract: Contract{
                        Input:         []byte{0xde, 0xad, 0xbe, 0xef},
                        Bytecode:      Bytecode{},
                        Hash:          &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
                        TargetAddress: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
                        Caller:        Address{Addr: [20]byte{0x1b, 0x3d, 0x5f, 0x7a, 0x9c}},
                        CallValue:     U256(big.NewInt(1000)),
                    },
                },
            },
        },
    }

    frameOrResult := FrameOrResult{
        Type:  Frame_FrameOrResult,
        Frame: callFrame,
    }

    // Assert
    assert.Equal(t, Frame_FrameOrResult, frameOrResult.Type)
    assert.Equal(t, callFrame, frameOrResult.Frame)
}
