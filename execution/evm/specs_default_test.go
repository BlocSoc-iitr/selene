package evm

import (
	"encoding/json"
	"testing"

	
)

const (
	TxTypeCreate = iota
)




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
