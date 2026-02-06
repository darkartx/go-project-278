package internal

import (
	"strings"
	"testing"
)

func TestGenerateShortName(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	tests := []struct {
		min, max uint
	}{
		{1, 1},
		{5, 10},
		{0, 5},
		{10, 15},
	}

	for _, tt := range tests {
		name := GenerateShortName(tt.min, tt.max)
		length := uint(len(name))

		if length < tt.min || length > tt.max {
			t.Errorf("GenerateShortName(%d, %d) = length %d; want length in [%d, %d]", tt.min, tt.max, length, tt.min, tt.max)
		}

		for _, char := range name {
			if !strings.ContainsRune(alphabet, char) {
				t.Errorf("GenerateShortName(%d, %d) contains invalid char: %c", tt.min, tt.max, char)
			}
		}
	}
}
