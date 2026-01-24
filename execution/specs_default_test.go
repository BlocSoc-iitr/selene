package evm

import (
	"bytes"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

const (
	TxTypeCreate = iota
)

func TestTxEnvSerialization(t *testing.T) {
	gasPrice := big.NewInt(1000000000)
	value := big.NewInt(5000000000000000000)
	accessList := []AccessListItem{
		{
			Address: Address{
				Addr: [20]byte{0x01},
			},
			StorageKeys: []common.Hash{{0x1}, {0x2}},
		},
	}
	txEnv := TxEnv{
		Caller: Address{
			Addr: [20]byte{0xab},
		},
		GasLimit:          21000,
		GasPrice:          gasPrice,
		TransactTo:        TxKind{Address: nil, Type: TxTypeCreate},
		Value:             value,
		Data:              []byte("transaction data"),
		Nonce:             new(uint64),
		ChainID:           new(uint64),
		AccessList:        accessList,
		GasPriorityFee:    big.NewInt(2000000000),
		BlobHashes:        []B256{{0xef}},
		MaxFeePerBlobGas:  big.NewInt(3000000000),
		AuthorizationList: nil,
	}

	*txEnv.Nonce = 1
	*txEnv.ChainID = 1

	data, err := json.Marshal(txEnv)
	if err != nil {
		t.Fatalf("failed to marshal TxEnv: %v", err)
	}

	var unmarshaledTxEnv TxEnv
	if err := json.Unmarshal(data, &unmarshaledTxEnv); err != nil {
		t.Fatalf("failed to unmarshal TxEnv: %v", err)
	}

	if !compareTxEnv(txEnv, unmarshaledTxEnv) {
		t.Fatalf("expected %+v, got %+v", txEnv, unmarshaledTxEnv)
	}
}

func compareTxEnv(a, b TxEnv) bool {
	return a.Caller == b.Caller &&
		a.GasLimit == b.GasLimit &&
		a.GasPrice.Cmp(b.GasPrice) == 0 &&
		a.TransactTo == b.TransactTo &&
		a.Value.Cmp(b.Value) == 0 &&
		bytes.Equal(a.Data, b.Data) &&
		*a.Nonce == *b.Nonce &&
		*a.ChainID == *b.ChainID &&
		a.GasPriorityFee.Cmp(b.GasPriorityFee) == 0 &&
		compareAccessList(a.AccessList, b.AccessList) &&
		compareB256Slice(a.BlobHashes, b.BlobHashes) &&
		a.MaxFeePerBlobGas.Cmp(b.MaxFeePerBlobGas) == 0
}

func compareAccessList(a, b []AccessListItem) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Address != b[i].Address ||
			!compareHashSlice(a[i].StorageKeys, b[i].StorageKeys) {
			return false
		}
	}
	return true
}

func compareHashSlice(a, b []common.Hash) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func compareB256Slice(a, b []B256) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSpecIdJSONMarshalling(t *testing.T) {
	tests := []struct {
		specID   SpecId
		expected string
	}{
		{FRONTIER, `"FRONTIER"`},
		{HOMESTEAD, `"HOMESTEAD"`},
		{TANGERINE, `"TANGERINE"`},
		{SPURIOUS_DRAGON, `"SPURIOUS_DRAGON"`},
		{BYZANTIUM, `"BYZANTIUM"`},
		{PETERSBURG, `"PETERSBURG"`},
		{ISTANBUL, `"ISTANBUL"`},
		{MERGE, `"MERGE"`},
		{LATEST, `"LATEST"`},
	}

	for _, tt := range tests {
		t.Run(tt.specID.String(), func(t *testing.T) {
			data, err := json.Marshal(tt.specID)
			if err != nil {
				t.Fatalf("failed to marshal SpecId %d: %v", tt.specID, err)
			}

			if string(data) != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestSpecIdJSONUnmarshalling(t *testing.T) {
	tests := []struct {
		input    string
		expected SpecId
	}{
		{`"FRONTIER"`, FRONTIER},
		{`"HOMESTEAD"`, HOMESTEAD},
		{`"TANGERINE"`, TANGERINE},
		{`"SPURIOUS_DRAGON"`, SPURIOUS_DRAGON},
		{`"BYZANTIUM"`, BYZANTIUM},
		{`"PETERSBURG"`, PETERSBURG},
		{`"ISTANBUL"`, ISTANBUL},
		{`"MERGE"`, MERGE},
		{`"LATEST"`, LATEST},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var specID SpecId
			if err := json.Unmarshal([]byte(tt.input), &specID); err != nil {
				t.Fatalf("failed to unmarshal SpecId from %s: %v", tt.input, err)
			}

			if specID != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, specID)
			}
		})
	}

	t.Run("invalid input", func(t *testing.T) {
		var specID SpecId
		err := json.Unmarshal([]byte(`"INVALID"`), &specID)
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
	})
}

func TestSpecToGenericDefault(t *testing.T) {
	tests := []struct {
		specID   SpecId
		expected string
	}{
		{FRONTIER, "FRONTIER"},
		{HOMESTEAD, "HOMESTEAD"},
		{TANGERINE, "TANGERINE"},
		{SPURIOUS_DRAGON, "SPURIOUS_DRAGON"},
		{BYZANTIUM, "BYZANTIUM"},
		{PETERSBURG, "PETERSBURG"},
		{ISTANBUL, "ISTANBUL"},
		{MERGE, "MERGE"},
		{LATEST, "LATEST"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			spec := SpecToGeneric(tt.specID)
			if spec.SpecID().String() != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, spec.SpecID().String())
			}
		})
	}
}
