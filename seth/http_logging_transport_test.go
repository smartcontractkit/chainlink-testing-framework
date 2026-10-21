package seth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoggingTransport_TLS(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	for _, logLevel := range []string{"info", "trace"} {
		t.Run("untrusted certificate is rejected with log level "+logLevel, func(t *testing.T) {
			t.Setenv(LogLevelEnvVar, logLevel)
			client := &http.Client{Transport: NewLoggingTransport()}
			resp, err := client.Get(server.URL)
			if resp != nil {
				_ = resp.Body.Close()
			}
			require.Error(t, err, "expected TLS verification to fail for untrusted certificate")
			require.Contains(t, err.Error(), "certificate")
		})

		t.Run("untrusted certificate is accepted when verification is disabled with log level "+logLevel, func(t *testing.T) {
			t.Setenv(LogLevelEnvVar, logLevel)
			t.Setenv(InsecureSkipVerifyEnvVar, "true")
			client := &http.Client{Transport: NewLoggingTransport()}
			resp, err := client.Get(server.URL)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
