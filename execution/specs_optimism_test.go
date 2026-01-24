package evm

import (
	"encoding/json"
	"fmt"
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
