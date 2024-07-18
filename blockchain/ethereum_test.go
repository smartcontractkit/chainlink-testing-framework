/*
This should be removed when we migrate all Ethereum client code to Seth
*/
package blockchain

import (
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRPCCustomHeadersFromEnv(t *testing.T) {
	l := logging.GetTestLogger(t)
	tests := []struct {
		name            string
		headerEnvString string
		expected        http.Header
		expectedErr     error
	}{
		{
			name:            "single k-v",
			headerEnvString: "Host=http.com",
			expected:        http.Header{"Host": []string{"http.com"}},
		},
		{
			name:            "multiple k-v",
			headerEnvString: "Host=http.com,Accept=application/json",
			expected: http.Header{
				"Host":   []string{"http.com"},
				"Accept": []string{"application/json"},
			}},
		{
			name:            "empty value",
			headerEnvString: "Host=a,Accept",
			expectedErr:     InvalidHeadersErr,
		},
		{
			name:            "invalid value",
			headerEnvString: "Host=a,Accept",
			expectedErr:     InvalidHeadersErr,
		},
		{
			name:            "invalid k-v, multiple =",
			headerEnvString: "Host=a=b,Host=c",
			expectedErr:     InvalidHeadersErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ETH_RPC_HEADERS", tt.headerEnvString)
			result, err := ReadEnvRPCHeaders(l)
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
