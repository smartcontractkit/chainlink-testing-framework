package slice_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/slice"
)

func TestValidateAndDeduplicateAddresses(t *testing.T) {
	tests := []struct {
		name          string
		addresses     []string
		expected      []string
		expectError   bool
		hasDuplicates bool
	}{
		{
			name:          "no duplicates",
			addresses:     []string{"0xAbC1230000000000000000000000000000000000", "0x00000000000000000000000000000000000000Ab"},
			expected:      []string{"0xAbC1230000000000000000000000000000000000", "0x00000000000000000000000000000000000000Ab"},
			expectError:   false,
			hasDuplicates: false,
		},
		{
			name:          "with duplicates",
			addresses:     []string{"0xAbC1230000000000000000000000000000000000", "0xAbC1230000000000000000000000000000000000"},
			expected:      []string{"0xAbC1230000000000000000000000000000000000"},
			expectError:   false,
			hasDuplicates: true,
		},
		{
			name:          "invalid address",
			addresses:     []string{"0xZZZ0000000000000000000000000000000000000"},
			expected:      nil,
			expectError:   true,
			hasDuplicates: false,
		},
		{
			name:          "mixed valid and invalid addresses",
			addresses:     []string{"0xAbC1230000000000000000000000000000000000", "not_a_hex_address"},
			expected:      nil,
			expectError:   true,
			hasDuplicates: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function with test case addresses
			result, hadDuplicates, err := slice.ValidateAndDeduplicateAddresses(tt.addresses)

			// Check for expected error
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
				require.Equal(t, tt.hasDuplicates, hadDuplicates)
			}
		})
	}
}
