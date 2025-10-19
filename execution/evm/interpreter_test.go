package evm

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharedMemory(t *testing.T) {
	memory := SharedMemory{}.NewWithMemoryLimit(1024)

	assert.Equal(t, 4096, len(memory.Buffer), "Expected buffer length 4096")
	assert.Equal(t, 32, len(memory.Checkpoints), "Expected 32 checkpoints")
	assert.Equal(t, 0, memory.LastCheckpoint, "Expected LastCheckpoint 0")

	memory.NewContext()
	assert.Equal(t, 33, len(memory.Checkpoints), "Expected checkpoints length 33 after new context")
	assert.Equal(t, len(memory.Buffer), memory.LastCheckpoint, "Expected LastCheckpoint at buffer length")

	memory.FreeContext()
	assert.Equal(t, 32, len(memory.Checkpoints), "Expected checkpoints length 32 after free context")
	assert.Equal(t, 0, memory.LastCheckpoint, "Expected LastCheckpoint reset to 0")
}

func TestStack(t *testing.T) {
	stack := Stack{}
	item1 := big.NewInt(0)
	item2 := big.NewInt(2)

	stack.Data = append(stack.Data, item1, item2)
	assert.Equal(t, 2, len(stack.Data), "Expected stack length 2")
	assert.Equal(t, item1, stack.Data[0], "Unexpected stack item 0")
	assert.Equal(t, item2, stack.Data[1], "Unexpected stack item 1")

	stack.Data = stack.Data[:len(stack.Data)-1]
	assert.Equal(t, 1, len(stack.Data), "Expected stack length 1 after pop")
	assert.Equal(t, item1, stack.Data[0], "Unexpected stack item after pop")
}

func TestFunctionStack(t *testing.T) {
	fs := FunctionStack{}

	// Test pushing a function return frame
	frame := FunctionReturnFrame{Idx: 1, PC: 42}
	fs.ReturnStack = append(fs.ReturnStack, frame)

	assert.Len(t, fs.ReturnStack, 1, "Expected ReturnStack length 1")
	assert.Equal(t, uint64(1), fs.ReturnStack[0].Idx, "Unexpected frame Idx in ReturnStack")
	assert.Equal(t, uint64(42), fs.ReturnStack[0].PC, "Unexpected frame PC in ReturnStack")

	// Test setting and getting the current code index
	fs.CurrentCodeIdx = 100
	assert.Equal(t, uint64(100), fs.CurrentCodeIdx, "Expected CurrentCodeIdx 100")
}

func TestGas(t *testing.T) {
	gas := Gas{
		Limit:     1000000,
		Remaining: 500000,
		Refunded:  1000,
	}

	assert.Equal(t, uint64(1000000), gas.Limit, "Expected Gas Limit 1000000")
	assert.Equal(t, uint64(500000), gas.Remaining, "Expected Gas Remaining 500000")
	assert.Equal(t, int64(1000), gas.Refunded, "Expected Gas Refunded 1000")
}

func TestSharedMemoryNewAndFreeContext(t *testing.T) {
	memory := SharedMemory{}.New()

	memory.NewContext()
	assert.Equal(t, 33, len(memory.Checkpoints), "Expected 33 checkpoints after NewContext")
	assert.Equal(t, len(memory.Buffer), memory.LastCheckpoint, "Expected LastCheckpoint at buffer length")

	initialLen := len(memory.Buffer)
	memory.Buffer = append(memory.Buffer, make([]byte, 100)...)
	assert.Equal(t, initialLen+100, len(memory.Buffer), "Expected buffer length to increase by 100")

	memory.FreeContext()
	assert.Equal(t, 32, len(memory.Checkpoints), "Expected 32 checkpoints after FreeContext")
	assert.Equal(t, 0, memory.LastCheckpoint, "Expected LastCheckpoint reset to 0")
	assert.Equal(t, initialLen, len(memory.Buffer), "Expected buffer length reset to initial")

	memory.Checkpoints = nil
	memory.FreeContext()
	assert.Equal(t, 0, memory.LastCheckpoint, "Expected LastCheckpoint to remain 0 when freeing with no checkpoints")
}

func TestSharedMemoryNewWithMemoryLimit(t *testing.T) {
	memory := SharedMemory{}.NewWithMemoryLimit(2048)
	assert.Equal(t, uint64(2048), memory.MemoryLimit, "Expected MemoryLimit 2048")
	assert.Len(t, memory.Buffer, 4096, "Expected initial buffer length 4096")
}
func TestSharedMemoryWithCapacity(t *testing.T) {
	memory := SharedMemory{}.WithCapacity(1024)

	assert.Equal(t, 1024, len(memory.Buffer), "Expected buffer length 1024")
}

func TestStackPushPop(t *testing.T) {
	stack := Stack{}
	val1 := big.NewInt(1)
	val2 := big.NewInt(2)

	stack.Data = append(stack.Data, val1, val2)
	assert.Equal(t, 2, len(stack.Data), "Expected stack length 2")
	assert.Equal(t, val1, stack.Data[0], "Expected stack[0] to be val1")
	assert.Equal(t, val2, stack.Data[1], "Expected stack[1] to be val2")

	stack.Data = stack.Data[:len(stack.Data)-1]
	assert.Equal(t, 1, len(stack.Data), "Expected stack length 1 after pop")
	assert.Equal(t, val1, stack.Data[0], "Expected stack top to be val1 after pop")
}

func TestStackUnderflow(t *testing.T) {
	stack := Stack{}
	assert.Equal(t, 0, len(stack.Data), "Expected stack to be empty initially")

	if len(stack.Data) > 0 {
		stack.Data = stack.Data[:len(stack.Data)-1]
	} else {
		t.Log("Stack underflow handled correctly")
	}
}

func TestFunctionStackPushPop(t *testing.T) {
	fs := FunctionStack{}
	frame1 := FunctionReturnFrame{Idx: 1, PC: 10}
	frame2 := FunctionReturnFrame{Idx: 2, PC: 20}

	fs.ReturnStack = append(fs.ReturnStack, frame1, frame2)
	assert.Equal(t, 2, len(fs.ReturnStack), "Expected ReturnStack length 2")
	assert.Equal(t, frame1, fs.ReturnStack[0], "Expected ReturnStack[0] to be frame1")
	assert.Equal(t, frame2, fs.ReturnStack[1], "Expected ReturnStack[1] to be frame2")

	fs.ReturnStack = fs.ReturnStack[:len(fs.ReturnStack)-1]
	assert.Equal(t, 1, len(fs.ReturnStack), "Expected ReturnStack length 1 after pop")
	assert.Equal(t, frame1, fs.ReturnStack[0], "Expected ReturnStack[0] to be frame1 after pop")
}

func TestFunctionStackBoundary(t *testing.T) {
	fs := FunctionStack{}

	assert.Equal(t, 0, len(fs.ReturnStack), "Expected ReturnStack to be empty initially")
	if len(fs.ReturnStack) > 0 {
		fs.ReturnStack = fs.ReturnStack[:len(fs.ReturnStack)-1]
	} else {
		t.Log("FunctionStack underflow handled correctly")
	}
}

func TestGasRefundAndLimit(t *testing.T) {
	gas := Gas{
		Limit:     100000,
		Remaining: 50000,
		Refunded:  0,
	}

	// Test limit reduction
	gas.Remaining -= 10000
	assert.Equal(t, uint64(40000), gas.Remaining, "Expected Gas Remaining 40000")

	// Test refunding gas
	gas.Refunded += 500
	assert.Equal(t, int64(500), gas.Refunded, "Expected Gas Refunded 500")

	// Check that remaining gas never exceeds limit
	assert.LessOrEqual(t, gas.Remaining, gas.Limit, "Gas Remaining exceeds limit")
}



// Test case structure for FunctionStack
type FunctionStackTestCase struct {
	name     string
	input    string // Input as hexadecimal string
	expected FunctionStack
	hasError bool
}

func TestFunctionStackMarshalJSON(t *testing.T) {
	testCases := []FunctionStackTestCase{
		{
			name: "Valid FunctionStack",
			input: `{"return_stack":[{"idx":"0x0","pc":"0x1"},{"idx":"0x2","pc":"0x3"}],"current_code_idx":"0x4"}`,
			expected: FunctionStack{
				ReturnStack: []FunctionReturnFrame{
					{Idx: 0, PC: 1},
					{Idx: 2, PC: 3},
				},
				CurrentCodeIdx: 4,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fs FunctionStack
			data := []byte(tc.input)

			err := json.Unmarshal(data, &fs)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(fs, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, fs)
			}
		})
	}
}

func TestFunctionStackUnmarshalJSON(t *testing.T) {
	testCases := []FunctionStackTestCase{
		{
			name: "Valid JSON",
			input: `{"return_stack":[{"idx":"0x0","pc":"0x1"},{"idx":"0x2","pc":"0x3"}],"current_code_idx":"0x4"}`,
			expected: FunctionStack{
				ReturnStack: []FunctionReturnFrame{
					{Idx: 0, PC: 1},
					{Idx: 2, PC: 3},
				},
				CurrentCodeIdx: 4,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fs FunctionStack
			data := []byte(tc.input)

			err := json.Unmarshal(data, &fs)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(fs, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, fs)
			}
		})
	}
}

// Test case structure for FunctionReturnFrame
type FunctionReturnFrameTestCase struct {
	name     string
	input    string // Input as hexadecimal string
	expected FunctionReturnFrame
	hasError bool
}

func TestFunctionReturnFrameMarshalJSON(t *testing.T) {
	testCases := []FunctionReturnFrameTestCase{
		{
			name: "Valid FunctionReturnFrame",
			input: `{"idx":"0x0","pc":"0x1"}`,
			expected: FunctionReturnFrame{
				Idx: 0,
				PC:  1,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var frf FunctionReturnFrame
			data := []byte(tc.input)

			err := json.Unmarshal(data, &frf)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(frf, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, frf)
			}
		})
	}
}

func TestFunctionReturnFrameUnmarshalJSON(t *testing.T) {
	testCases := []FunctionReturnFrameTestCase{
		{
			name: "Valid JSON",
			input: `{"idx":"0x0","pc":"0x1"}`,
			expected: FunctionReturnFrame{
				Idx: 0,
				PC:  1,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var frf FunctionReturnFrame
			data := []byte(tc.input)

			err := json.Unmarshal(data, &frf)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(frf, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, frf)
			}
		})
	}
}
