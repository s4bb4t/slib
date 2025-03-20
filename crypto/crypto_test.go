package crypto

import (
	"fmt"
	"testing"
)

func TestGenSepKey(t *testing.T) {
	tests := []struct {
		name      string
		base      int
		parts     int
		wantErr   bool
		expectLen int
	}{
		{
			name:      "Valid parts and base",
			base:      123,
			parts:     5,
			wantErr:   false,
			expectLen: 5*(lenOfPart+1) - 1, // length includes separators
		},
		{
			name:    "Parts exceed maximum limit",
			base:    123,
			parts:   65,
			wantErr: true,
		},
		{
			name:      "Parts at maximum allowed limit",
			base:      123,
			parts:     64,
			wantErr:   false,
			expectLen: 64*(lenOfPart+1) - 1,
		},
		{
			name:      "Minimum valid parts",
			base:      123,
			parts:     1,
			wantErr:   false,
			expectLen: lenOfPart + 1 - 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenSepKey(tt.base, tt.parts)
			fmt.Println(string(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("GenSepKey() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.expectLen {
				t.Errorf("GenSepKey() generated length = %d, expected = %d", len(got), tt.expectLen)
			}
		})
	}
}

func TestGenKey(t *testing.T) {
	tests := []struct {
		name      string
		base      int
		length    int
		wantErr   bool
		expectLen int
	}{
		{
			name:      "Valid key length",
			base:      123,
			length:    32,
			wantErr:   false,
			expectLen: 32,
		},
		{
			name:      "Key length less than base58 encoding limit",
			base:      123,
			length:    16,
			wantErr:   false,
			expectLen: 16,
		},
		{
			name:    "Key length exceeds maximum limit",
			base:    123,
			length:  65,
			wantErr: true,
		},
		{
			name:      "Key length at maximum allowed limit",
			base:      123,
			length:    64,
			wantErr:   false,
			expectLen: 64,
		},
		{
			name:      "Minimum valid key length",
			base:      123,
			length:    1,
			wantErr:   false,
			expectLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenKey(tt.base, tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenKey() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.expectLen {
				t.Errorf("GenKey() generated length = %d, expected = %d", len(got), tt.expectLen)
			}
		})
	}
}
