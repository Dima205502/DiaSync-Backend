package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	var testCases = []struct {
		in  string
		out string
	}{
		{"AsddF", "ae6a139b7d39cf004bbac3b8bb25acf420b654d61cf2296d603d148daadc1ff0"},
		{"1231", "52a6eb687cd22e80d3342eac6fcc7f2e19209e8f83eb9b82e81c6f3e6f30743b"},
		{"123jh3H", "c0a89ad3417d97b330bb4f62b4932f5e060baea22d1e2ce4c43096402f7e746a"},
	}

	for _, tt := range testCases {
		result := HashPassword(tt.in)

		if result != tt.out {
			t.Errorf("got %s, want %s", result, tt.out)
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	var testCases = []struct {
		password string
		hash     string
		out      bool
	}{
		{"AsddF", "ae6a139b7d39cf004bbac3b8bb25acf420b654d61cf2296d603d148daadc1ff0", true},
		{"12315", "52a6eb687cd22e80d3342eac6fcc7f2e19209e8f83eb9b82e81c6f3e6f30743b", false},
		{"123jh3H", "c0a89ad3417d97b330bb4f62b4932f5e060baea22d1e2ce4c43096402f7e74f6a", false},
	}

	for _, tt := range testCases {
		result := CheckPasswordHash(tt.password, tt.hash)

		if result != tt.out {
			t.Errorf("got %t, want %t", result, tt.out)
		}
	}
}
