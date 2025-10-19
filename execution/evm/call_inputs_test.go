package evm

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
)

type TestCase1 struct {
	name     string
	input    string
	expected CallInputs
	hasError bool
}

func TestCallInputsJSON(t *testing.T) {
	testCases := []TestCase1{
		{
			name: "Valid CallInputs",
			input: `{
				"input": "0x1234",
				"return_memory_offset": {"start": 0, "end": 32},
				"gas_limit": "0x186a0",
				"bytecode_address": "0x1234567890abcdef1234567890abcdef12345678",
				"target_address": "0x1234567890abcdef1234567890abcdef12345678",
				"caller": "0x9876543210fedcba9876543210fedcba98765432",
				"value": {"value_type": "transfer", "amount": "0xde0b6b3a7640000"},
				"scheme": "0x00",
				"is_static": false,
				"is_eof": false
			}`,

			expected: CallInputs{
				Input:              hexToBytes("0x1234"),
				ReturnMemoryOffset: Range{Start: 0, End: 32},
				GasLimit:           100000,
				BytecodeAddress:    Address{Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678")},
				TargetAddress:      Address{Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678")},
				Caller:             Address{Addr: parseHexAddress("0x9876543210fedcba9876543210fedcba98765432")},
				Value: CallValue{
					ValueType: "transfer",
					Amount:    parseU256("0xde0b6b3a7640000"),
				},
				Scheme:   ICall,
				IsStatic: false,
				IsEof:    false,
			},
			hasError: false,
		},
		{
			name: "Invalid CallInputs",
			input: `{
				"input": "not_a_hex",
				"gas_limit": "not_a_number"
			}`,
			expected: CallInputs{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var callInputs CallInputs
			err := json.Unmarshal([]byte(tc.input), &callInputs)

			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError {
				expected := tc.expected
				if !compareCallInputs(&callInputs, &expected) {
					t.Errorf("expected: %+v, got: %+v", expected, callInputs)
				}
			}
		})

	}
}

type TestCase2 struct {
	name     string
	input    string
	expected CallValue
	hasError bool
}

func TestCallValueJSON(t *testing.T) {
	testCases := []TestCase2{
		{
			name: "Valid Transfer CallValue",
			input: `{
				"value_type": "transfer", 
				"amount": "0xde0b6b3a7640000"
			}`,
			expected: CallValue{
				ValueType: "transfer",
				Amount:    parseU256("0xde0b6b3a7640000"),
			},
			hasError: false,
		},
		{
			name: "Valid Apparent CallValue",
			input: `{
				"value_type": "apparent", 
				"amount": "0x6f05b59d3b20000"
			}`,
			expected: CallValue{
				ValueType: "apparent",
				Amount:    parseU256("0x6f05b59d3b20000"),
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var callValue CallValue
			err := json.Unmarshal([]byte(tc.input), &callValue)

			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError {
				expected := tc.expected
				if !compareCallValues(&callValue, &expected) {
					t.Errorf("expected: %+v, got: %+v", expected, callValue)
				}
			}
		})
	}
}

type TestCase3 struct {
	name     string
	input    string
	expected CreateInputs
	hasError bool
}

func TestCreateInputsJSON(t *testing.T) {
	testCases := []TestCase3{
		{
			name: "Valid CreateInputs",
			input: `{
				"caller": "0x1234567890abcdef1234567890abcdef12345678",
				"scheme": {
					"scheme_type": "0x00",
					"salt": "0x1234567890abcdef1234567890abcdef12345678"
				},
				"value": "0xde0b6b3a7640000",
				"init_code": "0x010203",
				"gas_limit": "0x186a0"
			}`,
			expected: CreateInputs{
				Caller: Address{Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678")},
				Scheme: CreateScheme{
					SchemeType: SchemeTypeCreate,
					Salt:       parsePtrU256("0x1234567890abcdef1234567890abcdef12345678"),
				},
				Value:    parseU256("0xde0b6b3a7640000"),
				InitCode: hexToBytes("0x010203"),
				GasLimit: 100000,
			},
			hasError: false,
		},
		{
			name: "Invalid CreateInputs",
			input: `{
				"caller": "not_an_address",
				"scheme": {
					"scheme_type": "invalid"
				}
			}`,
			expected: CreateInputs{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var createInputs CreateInputs
			err := json.Unmarshal([]byte(tc.input), &createInputs)

			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError {
				expected := tc.expected
				if !compareCreateInputs(&createInputs, &expected) {
					t.Errorf("expected: %+v, got: %+v", expected, createInputs)
				}
			}
		})
	}
}

type TestCase4 struct {
	name     string
	input    string
	expected CreateScheme
	hasError bool
}

func TestCreateSchemeJSON(t *testing.T) {
	testCases := []TestCase4{
		{
			name: "Valid Create Scheme",
			input: `{
				"scheme_type": "0x00"
			}`,
			expected: CreateScheme{
				SchemeType: SchemeTypeCreate,
				Salt:       nil,
			},
			hasError: false,
		},
		{
			name: "Valid Create2 Scheme",
			input: `{
				"scheme_type": "0x01",
				"salt": "0x1234567890abcdef1234567890abcdef12345678"
			}`,
			expected: CreateScheme{
				SchemeType: SchemeTypeCreate2,
				Salt:       parsePtrU256("0x1234567890abcdef1234567890abcdef12345678"),
			},
			hasError: false,
		},
		{
			name: "Invalid CreateScheme",
			input: `{
				"scheme_type": "invalid"
			}`,
			expected: CreateScheme{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var createScheme CreateScheme
			err := json.Unmarshal([]byte(tc.input), &createScheme)

			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError {
				expected := tc.expected
				if !compareCreateSchemes(&createScheme, &expected) {
					t.Errorf("expected: %+v, got: %+v", expected, createScheme)
				}
			}
		})
	}
}

// Helper functions from the original test file
// func hexToBytes(hexStr string) []byte {
// 	bytes, err := hex.DecodeString(hexStr[2:])
// 	if err != nil {
// 		return nil
// 	}
// 	return bytes
// }

// func parseHexAddress(hexStr string) [20]byte {
// 	var addr [20]byte
// 	bytes, _ := hex.DecodeString(hexStr[2:])
// 	copy(addr[:], bytes)
// 	return addr
// }

func parseU256(hexStr string) U256 {
	value, ok := new(big.Int).SetString(hexStr[2:], 16)
	if !ok {
		return U256(new(big.Int)) // Return a zero value or another default
	}
	return U256(value)
}
func parsePtrU256(hexStr string) *U256 {
	value, _ := new(big.Int).SetString(hexStr[2:], 16)
	u256 := U256(value)
	return &u256
}

// Comparison functions for our structs
func compareCallInputs(a, b *CallInputs) bool {
	return reflect.DeepEqual(a.Input, b.Input) &&
		reflect.DeepEqual(a.ReturnMemoryOffset, b.ReturnMemoryOffset) &&
		a.GasLimit == b.GasLimit &&
		a.BytecodeAddress == b.BytecodeAddress &&
		a.TargetAddress == b.TargetAddress &&
		a.Caller == b.Caller &&
		compareCallValues(&a.Value, &b.Value) &&
		a.Scheme == b.Scheme &&
		a.IsStatic == b.IsStatic &&
		a.IsEof == b.IsEof
}

func compareCallValues(a, b *CallValue) bool {
	return a.ValueType == b.ValueType &&
		reflect.DeepEqual(a.Amount, b.Amount)
}

func compareCreateInputs(a, b *CreateInputs) bool {
	return a.Caller == b.Caller &&
		compareCreateSchemes(&a.Scheme, &b.Scheme) &&
		reflect.DeepEqual(a.Value, b.Value) &&
		reflect.DeepEqual(a.InitCode, b.InitCode) &&
		a.GasLimit == b.GasLimit
}

// Comparison function for CreateSchemes
func compareCreateSchemes(a, b *CreateScheme) bool {
	// Compare SchemeType
	if a.SchemeType != b.SchemeType {
		return false
	}

	// Special handling for nil/non-nil salt
	if (a.Salt == nil) != (b.Salt == nil) {
		return false
	}

	// If both salts are non-nil, compare their values
	if a.Salt != nil && b.Salt != nil {
		return (*a.Salt).Cmp(*b.Salt) == 0
	}

	return true
}
