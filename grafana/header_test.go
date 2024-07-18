package grafana

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRPCCustomHeadersFromEnv(t *testing.T) {
	l := defaultLogger()
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
			expectedErr:     ErrInvalidHeaders,
		},
		{
			name:            "invalid value",
			headerEnvString: "Host=a,Accept",
			expectedErr:     ErrInvalidHeaders,
		},
		{
			name:            "invalid k-v, multiple =",
			headerEnvString: "Host=a=b,Host=c",
			expectedErr:     ErrInvalidHeaders,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CTF_HTTP_HEADERS", tt.headerEnvString)
			result, err := ReadEnvHTTPHeaders(l)
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
