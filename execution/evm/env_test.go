package evm

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnv(t *testing.T) {
	env := NewEnv()
	assert.NotNil(t, env, "NewEnv should not return nil")
}

func TestNewSignature(t *testing.T) {
	r := big.NewInt(12345)
	s := big.NewInt(67890)
	sig := NewSignature(1, r, s)

	assert.NotNil(t, sig)
	assert.Equal(t, uint8(1), sig.V)
	assert.Equal(t, r, sig.R)
	assert.Equal(t, s, sig.S)
}

func TestToRawSignature(t *testing.T) {
	r := big.NewInt(12345)
	s := big.NewInt(67890)
	sig := NewSignature(1, r, s)

	rawSig := sig.ToRawSignature()
	assert.Equal(t, 64, len(rawSig))
	assert.Equal(t, r.Bytes(), rawSig[32-len(r.Bytes()):32])
	assert.Equal(t, s.Bytes(), rawSig[64-len(s.Bytes()):])
}

func TestFromRawSignature(t *testing.T) {
	r := big.NewInt(12345)
	s := big.NewInt(67890)
	sig := NewSignature(1, r, s)
	rawSig := sig.ToRawSignature()

	reconstructedSig, err := FromRawSignature(rawSig, sig.V)
	assert.NoError(t, err)
	assert.Equal(t, sig.V, reconstructedSig.V)
	assert.Equal(t, r, reconstructedSig.R)
	assert.Equal(t, s, reconstructedSig.S)
}

func TestCfgEnvStruct(t *testing.T) {
	cfg := CfgEnv{
		ChainID:        1234,
		DisableBaseFee: true,
		MemoryLimit:    2048,
	}
	assert.Equal(t, uint64(1234), cfg.ChainID)
	assert.True(t, cfg.DisableBaseFee)
	assert.Equal(t, uint64(2048), cfg.MemoryLimit)
}

func TestTxKind(t *testing.T) {
	txKind := TxKind{Type: Call2}
	assert.Equal(t, Call2, txKind.Type)
	assert.Nil(t, txKind.Address)
}

func TestOptionalNonce(t *testing.T) {
	var nonceValue uint64 = 42
	nonce := OptionalNonce{Nonce: &nonceValue}
	assert.Equal(t, uint64(42), *nonce.Nonce)

	nonce = OptionalNonce{}
	assert.Nil(t, nonce.Nonce)
}

func TestAuthorizationStructs(t *testing.T) {
	chainID := ChainID(12345)
	auth := Authorization{
		ChainID: chainID,
		Address: Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		Nonce:   OptionalNonce{},
	}
	assert.Equal(t, chainID, auth.ChainID)
	assert.Equal(t, Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}}, auth.Address)
	assert.Nil(t, auth.Nonce.Nonce)

	signedAuth := SignedAuthorization{Inner: auth}
	assert.Equal(t, chainID, signedAuth.Inner.ChainID)

	recoveredAuth := RecoveredAuthorization{Inner: auth}
	assert.Equal(t, chainID, recoveredAuth.Inner.ChainID)
	assert.Nil(t, recoveredAuth.Authority)
}

func TestBlobExcessGasAndPrice(t *testing.T) {
	blob := BlobExcessGasAndPrice{
		ExcessGas:    1000,
		BlobGasPrice: 2000,
	}
	assert.Equal(t, uint64(1000), blob.ExcessGas)
	assert.Equal(t, uint64(2000), blob.BlobGasPrice)
}

func TestBlockEnvStruct(t *testing.T) {
	blockEnv := BlockEnv{
		Number:     U256(big.NewInt(1234)),
		Timestamp:  U256(big.NewInt(5678)),
		Coinbase:   Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		GasLimit:   U256(big.NewInt(1234)),
		BaseFee:    U256(big.NewInt(1234)),
		Difficulty: U256(big.NewInt(1234)),
		Prevrandao: &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee},
		BlobExcessGasAndPrice: &BlobExcessGasAndPrice{
			ExcessGas:    90,
			BlobGasPrice: 80,
		},
	}
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.Number)
	assert.Equal(t, U256(big.NewInt(5678)), blockEnv.Timestamp)
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.GasLimit)
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.BaseFee)
	assert.Equal(t, U256(big.NewInt(1234)), blockEnv.Difficulty)
	assert.Equal(t, &B256{0xaa, 0xbb, 0xcc, 0xdd, 0xee}, blockEnv.Prevrandao)
	assert.Equal(t, &BlobExcessGasAndPrice{ExcessGas: 90, BlobGasPrice: 80}, blockEnv.BlobExcessGasAndPrice)
}

func TestAccessListItem(t *testing.T) {
	accessList := AccessListItem{
		Address:     Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}},
		StorageKeys: []B256{{0xaa, 0xbb, 0xcc, 0xdd, 0xee}},
	}
	assert.Equal(t, Address{Addr: [20]byte{0x1a, 0x2b, 0x3c, 0x4d, 0x5e}}, accessList.Address)
	assert.Equal(t, []B256{{0xaa, 0xbb, 0xcc, 0xdd, 0xee}}, accessList.StorageKeys)
}

// type CfgEnvTestCase struct {
// 	name     string
// 	input    CfgEnv
// 	expected string
// 	hasError bool
// }

// func TestCfgEnvMarshalJSON(t *testing.T) {
// 	var limit uint64 = 1024
// 	testCases := []CfgEnvTestCase{
// 		{
// 			name: "Valid CfgEnv",
// 			input: CfgEnv{
// 				ChainID:                     1,
// 				KzgSettings:                 nil,
// 				PerfAnalyseCreatedBytecodes: AnalysisKind(20),
// 				LimitContractCodeSize:       &limit, // Ensure this is set to a non-nil pointer
// 				MemoryLimit:                 1024,
// 				DisableBalanceCheck:         true,
// 				DisableBlockGasLimit:        false, // Ensure this is included
// 				DisableEIP3607:              false, // Ensure this is included
// 				DisableGasRefund:            true,
// 				DisableBaseFee:              false, // Ensure this is included
// 				DisableBeneficiaryReward:    true,
// 			},
// 			expected: `{"chain_id":1,"perf_analyse_created_bytecodes":20,"limit_contract_code_size":1024,"memory_limit":1024,"disable_balance_check":true,"disable_block_gas_limit":false,"disable_eip3607":false,"disable_gas_refund":true,"disable_base_fee":false,"disable_beneficiary_reward":true}`,
// 			hasError: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			data, err := json.Marshal(tc.input)
// 			if (err != nil) != tc.hasError {
// 				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 			}
// 			fmt.Printf("Marshalled data: %s\n", string(data)) // Debug line
// 			if !tc.hasError && string(data) != tc.expected {
// 				t.Errorf("expected: %s, got: %s", tc.expected, string(data))
// 			}
// 		})
// 	}
// }

// func TestCfgEnvUnmarshalJSON(t *testing.T) {
// 	var limit uint64 = 1024
// 	testCases := []CfgEnvTestCase{
// 		{
// 			name: "Valid JSON",
// 			input: CfgEnv{
// 				ChainID:                     1,
// 				KzgSettings:                 nil,
// 				PerfAnalyseCreatedBytecodes: AnalysisKind(20),
// 				LimitContractCodeSize:       &limit, // Ensure this is set to a non-nil pointer
// 				MemoryLimit:                 1024,
// 				DisableBalanceCheck:         true,
// 				DisableBlockGasLimit:        false, // Ensure this is included
// 				DisableEIP3607:              false, // Ensure this is included
// 				DisableGasRefund:            true,
// 				DisableBaseFee:              false, // Ensure this is included
// 				DisableBeneficiaryReward:    true,
// 			},
// 			expected: `{"chain_id":"0x1","perf_analyse_created_bytecodes":"0x14","limit_contract_code_size":"0x400","memory_limit":"0x400","disable_balance_check":true,"disable_block_gas_limit":false,"disable_eip3607":false,"disable_gas_refund":true,"disable_base_fee":false,"disable_beneficiary_reward":true}`,
// 			hasError: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var ce CfgEnv
// 			data := []byte(tc.expected)

// 			err := json.Unmarshal(data, &ce)
// 			if (err != nil) != tc.hasError {
// 				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 			}
// 			if !tc.hasError && !reflect.DeepEqual(ce, tc.input) {
// 				t.Errorf("expected: %+v, got: %+v", tc.input, ce)
// 			}
// 		})
// 	}
// }

// type BlockEnvTestCase struct {
// 	name     string
// 	input    BlockEnv
// 	expected string
// 	hasError bool
// }

// func TestBlockEnvMarshalJSON(t *testing.T) {
// 	testCases := []BlockEnvTestCase{
// 		{
// 			name: "Valid BlockEnv",
// 			input: BlockEnv{
// 				Number:                parseU256("0xde0b6b3a7640000"),
// 				Coinbase:              Address{Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678")}, // Replace with actual address as needed
// 				Timestamp:             parseU256("0xde0b6b3a7640000"),
// 				GasLimit:              parseU256("0xde0b6b3a7640000"),
// 				BaseFee:               parseU256("0xde0b6b3a7640000"),
// 				Difficulty:            parseU256("0xde0b6b3a7640000"),
// 				Prevrandao:            nil,
// 				BlobExcessGasAndPrice: nil,
// 			},
// 			expected: `{"number":1000000000000000000,"coinbase":{"Addr":[18,52,86,120,144,171,205,239,18,52,86,120,144,171,205,239,18,52,86,120]},"timestamp":1000000000000000000,"gas_limit":1000000000000000000,"basefee":1000000000000000000,"difficulty":1000000000000000000}`,
// 			hasError: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			data, err := json.Marshal(tc.input)
// 			if (err != nil) != tc.hasError {
// 				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 			}
// 			if !tc.hasError && string(data) != tc.expected {
// 				t.Errorf("expected: %s, got: %s", tc.expected, string(data))
// 			}
// 		})
// 	}
// }

// func TestBlockEnvUnmarshalJSON(t *testing.T) {
// 	testCases := []BlockEnvTestCase{
// 		{
// 			name: "Valid JSON",
// 			input: BlockEnv{
// 				Number:                parseU256("0xde0b6b3a7640000"),
// 				Coinbase:              Address{Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678")}, // Replace with actual address as needed
// 				Timestamp:             parseU256("0xde0b6b3a7640000"),
// 				GasLimit:              parseU256("0xde0b6b3a7640000"),
// 				BaseFee:               parseU256("0xde0b6b3a7640000"),
// 				Difficulty:            parseU256("0xde0b6b3a7640000"),
// 				Prevrandao:            nil,
// 				BlobExcessGasAndPrice: nil,
// 			},
// 			expected: `{"number":"0xde0b6b3a7640000","coinbase":"0x1234567890abcdef1234567890abcdef12345678","timestamp":"0xde0b6b3a7640000","gas_limit":"0xde0b6b3a7640000","basefee":"0xde0b6b3a7640000","difficulty":"0xde0b6b3a7640000"}`,
// 			hasError: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var be BlockEnv
// 			data := []byte(tc.expected)

// 			err := json.Unmarshal(data, &be)
// 			if (err != nil) != tc.hasError {
// 				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 			}
// 			if !tc.hasError && !reflect.DeepEqual(be, tc.input) {
// 				t.Errorf("expected: %+v, got: %+v", tc.input, be)
// 			}
// 		})
// 	}
// }

// type BlobExcessGasAndPriceTestCase struct {
// 	name     string
// 	input    BlobExcessGasAndPrice
// 	expected string
// 	hasError bool
// }

// func TestBlobExcessGasAndPriceMarshalJSON(t *testing.T) {
// 	testCases := []BlobExcessGasAndPriceTestCase{
// 		{
// 			name: "Valid BlobExcessGasAndPrice",
// 			input: BlobExcessGasAndPrice{
// 				ExcessGas:    1,
// 				BlobGasPrice: 2,
// 			},
// 			expected: `{"excess_gas":1,"blob_gas_price":2}`,
// 			hasError: false,
// 		},
// 	}
// }
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			data, err := json.Marshal(tc.input)
// 			if (err != nil) != tc.hasError {
// 				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 			}
// 			if !tc.hasError && string(data) != tc.expected {
// 				t.Errorf("expected: %s, got: %s", tc.expected, string(data))
// 			}
// 		})
// 	}
// }

// func TestBlobExcessGasAndPriceUnmarshalJSON(t *testing.T) {
// 	testCases := []BlobExcessGasAndPriceTestCase{
// 		{
// 			name: "Valid JSON",
// 			input: BlobExcessGasAndPrice{
// 				ExcessGas:    1,
// 				BlobGasPrice: 2,
// 			},
// 			expected: `{"excess_gas":"0x1","blob_gas_price":"0x2"}`,
// 			hasError: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var bgp BlobExcessGasAndPrice
// 			data := []byte(tc.expected)

// 			err := json.Unmarshal(data, &bgp)
// 			if (err != nil) != tc.hasError {
// 				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 			}
// 			if !tc.hasError && !reflect.DeepEqual(bgp, tc.input) {
// 				t.Errorf("expected: %+v, got: %+v", tc.input, bgp)
// 			}
// 		})
// 	}
// }

// type EnvKzgSettingsTestCase struct {
// 	name     string
// 	input    EnvKzgSettings
// 	expected string
// 	hasError bool
// }

// func TestEnvKzgSettingsMarshalJSON(t *testing.T) {
// 	testCases := []EnvKzgSettingsTestCase{
// 		{
// 			name: "Valid EnvKzgSettings",
// 			input: EnvKzgSettings{
// 				Mode:   "test",
// 				Custom: nil, // Populate if you add fields to KzgSettings
// 			},
// 			expected: `{"mode":"test"}`,
// 			hasError: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			data, err := json.Marshal(tc.input)
// 			if (err != nil) != tc.hasError {
// 				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 			}
// 			if !tc.hasError && string(data) != tc.expected {
// 				t.Errorf("expected: %s, got: %s", tc.expected, string(data))
// 			}
// 		})
// 	}
// }

// func TestEnvKzgSettingsUnmarshalJSON(t *testing.T) {
// 	testCases := []EnvKzgSettingsTestCase{
// 		{
// 			name: "Valid JSON",
// 			input: EnvKzgSettings{
// 				Mode:   "test",
// 				Custom: nil, // Populate if you add fields to KzgSettings
// 			},
// 			expected: `{"mode":"test"}`,
// 			hasError: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var eks EnvKzgSettings
// 			data := []byte(tc.expected)

// 				err := json.Unmarshal(data, &eks)
// 				if (err != nil) != tc.hasError {
// 					t.Errorf("expected error: %v, got: %v", tc.hasError, err)
// 				}
// 				if !tc.hasError && !reflect.DeepEqual(eks, tc.input) {
// 					t.Errorf("expected: %+v, got: %+v", tc.input, eks)
// 				}
// 			})
// 		}
// 	}
type CfgEnvTestCase struct {
	name     string
	input    string
	expected CfgEnv
	hasError bool
}

func TestCfgEnvMarshalJSON(t *testing.T) {
	var limit uint64 = 1024
	testCases := []CfgEnvTestCase{
		{
			name: "Valid CfgEnv",
			input: `{"chain_id":1,"perf_analyse_created_bytecodes":20,"limit_contract_code_size":1024,"memory_limit":1024,"disable_balance_check":true,"disable_block_gas_limit":false,"disable_eip3607":false,"disable_gas_refund":true,"disable_base_fee":false,"disable_beneficiary_reward":true}`,
			expected: CfgEnv{
				ChainID:                     1,
				KzgSettings:                 nil,
				PerfAnalyseCreatedBytecodes: AnalysisKind(20),
				LimitContractCodeSize:       &limit, // Ensure this is set to a non-nil pointer
				MemoryLimit:                 1024,
				DisableBalanceCheck:         true,
				DisableBlockGasLimit:        false, // Ensure this is included
				DisableEIP3607:              false, // Ensure this is included
				DisableGasRefund:            true,
				DisableBaseFee:              false, // Ensure this is included
				DisableBeneficiaryReward:    true,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.expected)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			fmt.Printf("Marshalled data: %s\n", string(data)) // Debug line
			if !tc.hasError && string(data) != tc.input {
				t.Errorf("expected: %s, got: %s", tc.input, string(data))
			}
		})
	}
}

func TestCfgEnvUnmarshalJSON(t *testing.T) {
	var limit uint64 = 1024
	testCases := []CfgEnvTestCase{
		{
			name: "Valid JSON",
			input: `{"chain_id":"0x1","perf_analyse_created_bytecodes":"0x14","limit_contract_code_size":"0x400","memory_limit":"0x400","disable_balance_check":true,"disable_block_gas_limit":false,"disable_eip3607":false,"disable_gas_refund":true,"disable_base_fee":false,"disable_beneficiary_reward":true}`,
			expected: CfgEnv{
				ChainID:                     1,
				KzgSettings:                 nil,
				PerfAnalyseCreatedBytecodes: AnalysisKind(20),
				LimitContractCodeSize:       &limit, // Ensure this is set to a non-nil pointer
				MemoryLimit:                 1024,
				DisableBalanceCheck:         true,
				DisableBlockGasLimit:        false, // Ensure this is included
				DisableEIP3607:              false, // Ensure this is included
				DisableGasRefund:            true,
				DisableBaseFee:              false, // Ensure this is included
				DisableBeneficiaryReward:    true,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ce CfgEnv
			data := []byte(tc.input)

			err := json.Unmarshal(data, &ce)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(ce, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, ce)
			}
		})
	}
}

type BlockEnvTestCase struct {
	name     string
	input    string
	expected BlockEnv
	hasError bool
}

func TestBlockEnvMarshalJSON(t *testing.T) {
	testCases := []BlockEnvTestCase{
		{
			name: "Valid BlockEnv",
			input: `{"number":1000000000000000000,"coinbase":{"Addr":[18,52,86,120,144,171,205,239,18,52,86,120,144,171,205,239,18,52,86,120]},"timestamp":1000000000000000000,"gas_limit":1000000000000000000,"basefee":1000000000000000000,"difficulty":1000000000000000000}`,
			expected: BlockEnv{
				Number:                parseU256("0xde0b6b3a7640000"),
				Coinbase:              Address{Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678")}, // Replace with actual address as needed
				Timestamp:             parseU256("0xde0b6b3a7640000"),
				GasLimit:              parseU256("0xde0b6b3a7640000"),
				BaseFee:               parseU256("0xde0b6b3a7640000"),
				Difficulty:            parseU256("0xde0b6b3a7640000"),
				Prevrandao:            nil,
				BlobExcessGasAndPrice: nil,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.expected)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && string(data) != tc.input {
				t.Errorf("expected: %s, got: %s", tc.input, string(data))
			}
		})
	}
}

func TestBlockEnvUnmarshalJSON(t *testing.T) {
	testCases := []BlockEnvTestCase{
		{
			name: "Valid JSON",
			input: `{"number":"0xde0b6b3a7640000","coinbase":"0x1234567890abcdef1234567890abcdef12345678","timestamp":"0xde0b6b3a7640000","gas_limit":"0xde0b6b3a7640000","basefee":"0xde0b6b3a7640000","difficulty":"0xde0b6b3a7640000"}`,
			expected: BlockEnv{
				Number:                parseU256("0xde0b6b3a7640000"),
				Coinbase:              Address{Addr: parseHexAddress("0x1234567890abcdef1234567890abcdef12345678")}, // Replace with actual address as needed
				Timestamp:             parseU256("0xde0b6b3a7640000"),
				GasLimit:              parseU256("0xde0b6b3a7640000"),
				BaseFee:               parseU256("0xde0b6b3a7640000"),
				Difficulty:            parseU256("0xde0b6b3a7640000"),
				Prevrandao:            nil,
				BlobExcessGasAndPrice: nil,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var be BlockEnv
			data := []byte(tc.input)

			err := json.Unmarshal(data, &be)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(be, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, be)
			}
		})
	}
}

type BlobExcessGasAndPriceTestCase struct {
	name     string
	input    string
	expected BlobExcessGasAndPrice
	hasError bool
}

func TestBlobExcessGasAndPriceMarshalJSON(t *testing.T) {
	testCases := []BlobExcessGasAndPriceTestCase{
		{
			name: "Valid BlobExcessGasAndPrice",
			input: `{"excess_gas":1,"blob_gas_price":2}`,
			expected: BlobExcessGasAndPrice{
				ExcessGas:    1,
				BlobGasPrice: 2,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.expected)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && string(data) != tc.input {
				t.Errorf("expected: %s, got: %s", tc.input, string(data))
			}
		})
	}
}

func TestBlobExcessGasAndPriceUnmarshalJSON(t *testing.T) {
	testCases := []BlobExcessGasAndPriceTestCase{
		{
			name: "Valid JSON",
			input: `{"excess_gas":"0x1","blob_gas_price":"0x2"}`,
			expected: BlobExcessGasAndPrice{
				ExcessGas:    1,
				BlobGasPrice: 2,
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bgp BlobExcessGasAndPrice
			data := []byte(tc.input)

			err := json.Unmarshal(data, &bgp)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(bgp, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, bgp)
			}
		})
	}
}

type EnvKzgSettingsTestCase struct {
	name     string
	input    string
	expected EnvKzgSettings
	hasError bool
}

func TestEnvKzgSettingsMarshalJSON(t *testing.T) {
	testCases := []EnvKzgSettingsTestCase{
		{
			name: "Valid EnvKzgSettings",
			input: `{"mode":"test"}`,
			expected: EnvKzgSettings{
				Mode:   "test",
				Custom: nil, // Populate if you add fields to KzgSettings
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.expected)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && string(data) != tc.input {
				t.Errorf("expected: %s, got: %s", tc.input, string(data))
			}
		})
	}
}

func TestEnvKzgSettingsUnmarshalJSON(t *testing.T) {
	testCases := []EnvKzgSettingsTestCase{
		{
			name: "Valid JSON",
			input: `{"mode":"test"}`,
			expected: EnvKzgSettings{
				Mode:   "test",
				Custom: nil, // Populate if you add fields to KzgSettings
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var eks EnvKzgSettings
			data := []byte(tc.input)

			err := json.Unmarshal(data, &eks)
			if (err != nil) != tc.hasError {
				t.Errorf("expected error: %v, got: %v", tc.hasError, err)
			}
			if !tc.hasError && !reflect.DeepEqual(eks, tc.expected) {
				t.Errorf("expected: %+v, got: %+v", tc.expected, eks)
			}
		})
	}
}