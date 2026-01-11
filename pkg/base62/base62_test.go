package base62

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"Zero", 0, "0"},
		{"Um", 1, "1"},
		{"Base", 61, "z"},
		{"Base+1", 62, "10"},
		{"Número Grande", 999, "G7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.input)
			if got != tt.expected {
				t.Errorf("Encode(%d) = %s; expected %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
	}{
		{"Zero", "0", 0},
		{"Um", "1", 1},
		{"Base", "z", 61},
		{"Base+1", "10", 62},
		{"Número Grande", "G7", 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("Decode(%s) = %d; expected %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIntegrity(t *testing.T) {
	var initial uint64 = 123456789
	encoded := Encode(initial)
	decoded, _ := Decode(encoded)

	if initial != decoded {
		t.Errorf("Falha de integridade: %d -> %s -> %d", initial, encoded, decoded)
	}
}
