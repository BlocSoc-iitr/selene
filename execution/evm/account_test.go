package evm

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"
)

// Test case structure for unmarshaling JSON data
type TestCase struct {
	name     string
	input    string
	expected Opcode
	hasError bool
}

func TestOpcodeUnmarshalJSON(t *testing.T) {
	testCases := []TestCase{
		{
			name: "Valid JSON",
			input: `{
  "initcode": {
    "header": {
      "types_size": "0xa",
      "code_sizes": ["0x14", "0x1e"],
      "container_sizes": ["0x28", "0x32"],
      "data_size": "0x3c",
      "sum_code_sizes": "0x46",
      "sum_container_sizes": "0x50"
    },
    "body": {
      "types_section": [
        {
          "inputs": "0x1",
          "outputs": "0x2",
          "max_stack_size": "0x200"
        }
      ],
      "code_section": ["0x010203"],
      "container_section": ["0x040506"],
      "data_section": "0x070809",
      "is_data_filled": true
    },
    "raw": "0x010203"
  },
  "created_address": "0x1234567890abcdef1234567890abcdef12345678",
  "input": "0x1234567890abcdef1234567890abcdef12345678"
}`,
			expected: Opcode{
				InitCode: Eof{
					Header: EofHeader{
						TypesSize:         10,
						CodeSizes:         []uint16{20, 30},
						ContainerSizes:    []uint16{40, 50},
						DataSize:          60,
						SumCodeSizes:      70,
						SumContainerSizes: 80,
					},
					Body: EofBody{
						TypesSection: []TypesSection{
							{Inputs: 1, Outputs: 2, MaxStackSize: 512},
						},
						CodeSection:      [][]byte{hexToBytes("0x010203")},
						ContainerSection: [][]byte{hexToBytes("0x040506")},
						DataSection:      hexToBytes("0x070809"),
						IsDataFilled:     true,
					},
					Raw: hexToBytes("0x010203"),
				},
				CreatedAddress: Address{
					Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678"),
				},
				Input: hexToBytes("0x1234567890abcdef1234567890abcdef12345678"),
			},
			hasError: false,
		},
		{
			name: "Invalid JSON",
			input: `{
				"initcode": {
					"header": {
						"types_size": "not_a_number",
						"code_sizes": ["20", "30"],
						"container_sizes": ["40", "50"],
						"data_size": "60",
						"sum_code_sizes": "70",
						"sum_container_sizes": "80"
					},
					"body": {
						"types_section": [
							{"inputs": 1, "outputs": 2, "max_stack_size": "512"}
						],
						"code_section": ["0x010203"],
						"container_section": ["0x040506"],
						"data_section": "0x070809",
						"is_data_filled": true
					},
					"raw": "0x010203"
				},
				"created_address": "0x1234567890abcdef1234567890abcdef12345678",
				"input": "0x1234567890abcdef1234567890abcdef12345678"
			}`,
			expected: Opcode{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var opcode Opcode
			err := json.Unmarshal([]byte(tc.input), &opcode)

			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !compareOpcodes(opcode, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, opcode)
			}
		})
	}
}

func hexToBytes(hexStr string) []byte {
	bytes, err := hex.DecodeString(hexStr[2:])
	if err != nil {
		return nil
	}
	return bytes
}
func parseHexAddress(hexStr string) [20]byte {
	var addr [20]byte
	bytes, _ := hex.DecodeString(hexStr[2:])
	copy(addr[:], bytes)
	return addr
}
func compareOpcodes(a, b Opcode) bool {
	return compareEof(a.InitCode, b.InitCode) && a.CreatedAddress == b.CreatedAddress && reflect.DeepEqual(a.Input, b.Input)
}
func compareEof(a, b Eof) bool {
	return compareEofHeader(a.Header, b.Header) &&
		compareEofBody(a.Body, b.Body) &&
		reflect.DeepEqual(a.Raw, b.Raw)
}

func compareEofHeader(a, b EofHeader) bool {
	return a.TypesSize == b.TypesSize &&
		reflect.DeepEqual(a.CodeSizes, b.CodeSizes) &&
		reflect.DeepEqual(a.ContainerSizes, b.ContainerSizes) &&
		a.DataSize == b.DataSize &&
		a.SumCodeSizes == b.SumCodeSizes &&
		a.SumContainerSizes == b.SumContainerSizes
}

// Comparison function for EofBody struct
func compareEofBody(a, b EofBody) bool {
	return reflect.DeepEqual(a.TypesSection, b.TypesSection) &&
		reflect.DeepEqual(a.CodeSection, b.CodeSection) &&
		reflect.DeepEqual(a.ContainerSection, b.ContainerSection) &&
		reflect.DeepEqual(a.DataSection, b.DataSection) &&
		a.IsDataFilled == b.IsDataFilled
}
