package evm

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecIdIsEnabledIn(t *testing.T) {
	tests := []struct {
		name   string
		specId SpecId
		other  SpecId
		want   bool
	}{
		{"FRONTIER >= FRONTIER", FRONTIER, FRONTIER, true},
		{"HOMESTEAD >= FRONTIER", HOMESTEAD, FRONTIER, true},
		{"FRONTIER < HOMESTEAD", FRONTIER, HOMESTEAD, false},
		{"LATEST >= FRONTIER", LATEST, FRONTIER, true},
		{"FRONTIER < LATEST", FRONTIER, LATEST, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.specId.IsEnabledIn(tt.other))
		})
	}
}

func TestTryFromUint8Optimism(t *testing.T) {
	tests := []struct {
		specID uint8
		want   SpecId
		ok     bool
	}{
		{0, FRONTIER, true},
		{1, FRONTIER_THAWING, true},
		{19, PRAGUE_EOF, true},
		{20, 0, false},
		{255, LATEST, true},
	}

	for _, tt := range tests {
		got, ok := TryFromUint8(tt.specID)
		assert.Equal(t, tt.ok, ok)
		assert.Equal(t, tt.want, got)
	}

}

func TestSpecToGeneric(t *testing.T) {
	tests := []struct {
		specID   SpecId
		want     string
		wantType string
	}{
		{FRONTIER, "FRONTIER", "evm.FrontierSpec"},
		{HOMESTEAD, "HOMESTEAD", "evm.HomesteadSpec"},
		{TANGERINE, "TANGERINE", "evm.TangerineSpec"},
		{SPURIOUS_DRAGON, "SPURIOUS_DRAGON", "evm.SpuriousDragonSpec"},
		{BYZANTIUM, "BYZANTIUM", "evm.ByzantiumSpec"},
		{PETERSBURG, "PETERSBURG", "evm.PetersburgSpec"},
		{ISTANBUL, "ISTANBUL", "evm.IstanbulSpec"},
		{BERLIN, "BERLIN", "evm.BerlinSpec"},
		{LONDON, "LONDON", "evm.LondonSpec"},
		{MERGE, "MERGE", "evm.MergeSpec"},
		{SHANGHAI, "SHANGHAI", "evm.ShanghaiSpec"},
		{CANCUN, "CANCUN", "evm.CancunSpec"},
		{PRAGUE, "PRAGUE", "evm.PragueSpec"},
		{PRAGUE_EOF, "PRAGUE_EOF", "evm.PragueEofSpec"},
		{LATEST, "LATEST", "evm.LatestSpec"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			spec := SpecToGeneric(tt.specID)

			// Check if the type matches
			if typeName := fmt.Sprintf("%T", spec); typeName != tt.wantType {
				t.Errorf("expected type %s, got %s", tt.wantType, typeName)
			}

			// Type assertion for SpecID method
			switch specTyped := spec.(type) {
			case interface{ SpecID() SpecId }:
				if got := specTyped.SpecID().String(); got != tt.want {
					t.Errorf("expected %s, got %s", tt.want, got)
				}
			default:
				t.Errorf("spec does not implement SpecID method for specID %d", tt.specID)
			}
		})
	}
}

func TestSpecIdMarshalJSON(t *testing.T) {
	tests := []struct {
		specId SpecId
		want   string
	}{
		{FRONTIER, `"FRONTIER"`},
		{HOMESTEAD, `"HOMESTEAD"`},
		{LATEST, `"LATEST"`},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			data, err := json.Marshal(tt.specId)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, string(data))
		})
	}
}

func TestSpecIdUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    SpecId
		wantErr bool
	}{
		{"Valid FRONTIER", `"FRONTIER"`, FRONTIER, false},
		{"Valid LATEST", `"LATEST"`, LATEST, false},
		{"Invalid SpecId", `"INVALID"`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var specId SpecId
			err := json.Unmarshal([]byte(tt.input), &specId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, specId)
			}
		})
	}
}
func TestTxEnvUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    TxEnv
		wantErr bool
	}{
		{
			name: "basic transaction",
			input: `{
				"caller": "0x1234567890123456789012345678901234567890",
				"gas_limit": "0x5208",
				"gas_price": "0x4a817c800",
				"transact_to": {"type": "1", "address": "0x2345678901234567890123456789012345678901"},
				"value": "0xde0b6b3a7640000",
				"data": "0x",
				"access_list": []
			}`,
			want: TxEnv{
				Caller:   Address{Addr: [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}},
				GasLimit: 21000,
				GasPrice: big.NewInt(20000000000),
				TransactTo: TxKind{
					Type:    Call2,
					Address: &Address{Addr: [20]byte{0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01}},
				},
				Value: big.NewInt(1000000000000000000),
				Data:  []byte{},
				// ChainID:    pointer(uint64(1)),
				AccessList: []AccessListItem{},
			},
			wantErr: false,
		},
		{
			name: "full transaction with all fields",
			input: `{
				"caller": "0x1234567890123456789012345678901234567890",
				"gas_limit": "0x5208",
				"gas_price": "0x4a817c800",
				"transact_to": {"type": "0"},
				"value": "0x0",
				"data": "0x1234",
				
				"gas_priority_fee": "0x1234",
				"blob_hashes": [
					"0x1234567890123456789012345678901234567890123456789012345678901234"
				],
				"max_fee_per_blob_gas": "0x5678",
				"authorization_list": {},
				"access_list": [
					{
						"address": "0x3456789012345678901234567890123456789012",
						"storage_keys": [
							"0x1234567890123456789012345678901234567890123456789012345678901234"
						]
					}
				]
			}`,
			want: TxEnv{
				Caller:   Address{Addr: [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}},
				GasLimit: 21000,
				GasPrice: big.NewInt(20000000000),
				TransactTo: TxKind{
					Type: Create2,
				},
				Value: big.NewInt(0),
				Data:  []byte{0x12, 0x34},
				// Nonce:   pointer(uint64(1)),
				// ChainID: pointer(uint64(1)),

				GasPriorityFee:    big.NewInt(4660),
				BlobHashes:        []B256{{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34}},
				MaxFeePerBlobGas:  big.NewInt(22136),
				AuthorizationList: &AuthorizationList{},
				AccessList: []AccessListItem{
					{
						Address: Address{Addr: [20]byte{0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12}},
						StorageKeys: []B256{
							{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid caller address",
			input: `{
				"caller": "invalid",
				"gas_limit": "0x5208",
				"gas_price": "0x4a817c800",
				"transact_to": {"type": "1", "address": "0x2345678901234567890123456789012345678901"},
				"value": "0x0",
				"data": "0x",
				"access_list": []
			}`,
			wantErr: true,
		},
		{
			name: "invalid gas limit",
			input: `{
				"caller": "0x1234567890123456789012345678901234567890",
				"gas_limit": "invalid",
				"gas_price": "0x4a817c800",
				"transact_to": {"type": "1", "address": "0x2345678901234567890123456789012345678901"},
				"value": "0x0",
				"data": "0x",
				"access_list": []
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TxEnv
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("TxEnv.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("TxEnv.UnmarshalJSON() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}


func TestTxEnvUnmarshalJSONEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    TxEnv
		wantErr bool
	}{
		{
			name: "empty access list",
			input: `{
				"caller": "0x1234567890123456789012345678901234567890",
				"gas_limit": "0x5208",
				"gas_price": "0x4a817c800",
				"transact_to": {"type": "1", "address": "0x2345678901234567890123456789012345678901"},
				"value": "0x0",
				"data": "0x",
				"access_list": []
			}`,
			want: TxEnv{
				Caller:   Address{Addr: [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}},
				GasLimit: 21000,
				GasPrice: big.NewInt(20000000000),
				TransactTo: TxKind{
					Type:    Call2,
					Address: &Address{Addr: [20]byte{0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01}},
				},
				Value:      big.NewInt(0),
				Data:       []byte{},
				AccessList: []AccessListItem{},
			},
			wantErr: false,
		},
		{
			name: "zero values",
			input: `{
				"caller": "0x0000000000000000000000000000000000000000",
				"gas_limit": "0x0",
				"gas_price": "0x0",
				"transact_to": {"type": "1", "address": "0x0000000000000000000000000000000000000000"},
				"value": "0x0",
				"data": "0x",
				"access_list": []
			}`,
			want: TxEnv{
				Caller:   Address{},
				GasLimit: 0,
				GasPrice: big.NewInt(0),
				TransactTo: TxKind{
					Type:    Call2,
					Address: &Address{},
				},
				Value:      big.NewInt(0),
				Data:       []byte{},
				AccessList: []AccessListItem{},
			},
			wantErr: false,
		},
		{
			name: "very large numbers",
			input: `{
				"caller": "0x1234567890123456789012345678901234567890",
				"gas_limit": "0xffffffffffffffff",
				"gas_price": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				"transact_to": {"type": "1", "address": "0x2345678901234567890123456789012345678901"},
				"value": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				"data": "0x",
				"access_list": []
			}`,
			want: TxEnv{
				Caller:   Address{Addr: [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}},
				GasLimit: ^uint64(0),
				GasPrice: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)),
				TransactTo: TxKind{
					Type:    Call2,
					Address: &Address{Addr: [20]byte{0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, 0x67, 0x89, 0x01}},
				},
				Value:      new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)),
				Data:       []byte{},
				AccessList: []AccessListItem{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TxEnv
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("TxEnv.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("TxEnv.UnmarshalJSON() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
