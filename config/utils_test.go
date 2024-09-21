package config

import (
	"testing"
)

func TestBytesSerialise(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{nil, "null"},
		{[]byte{}, "\"\""},
		{[]byte{0x01, 0x02, 0x03}, "\"010203\""},
	}

	for _, test := range tests {
		result, err := BytesSerialise(test.input)
		if err != nil {
			t.Fatalf("BytesSerialise(%v) returned an error: %v", test.input, err)
		}

		resultStr := string(result)

		if resultStr != test.expected {
			t.Errorf("BytesSerialise(%v) = %v, expected %v", test.input, resultStr, test.expected)
		}
	}
}

func TestBytesDeserialise(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"null", nil},
		{"\"\"", []byte{}},
		{"\"010203\"", []byte{0x01, 0x02, 0x03}},
	}

	for _, test := range tests {
		data := []byte(test.input)
		result, err := BytesDeserialise(data)
		if err != nil {
			t.Fatalf("BytesDeserialise(%v) returned an error: %v", test.input, err)
		}

		if !equal(result, test.expected) {
			t.Errorf("BytesDeserialise(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// Helper function to compare byte slices
func equal(a, b []byte) bool {
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
