package seth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pelletier/go-toml/v2"
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

		t.Run("untrusted certificate is accepted when verification is disabled via config with log level "+logLevel, func(t *testing.T) {
			t.Setenv(LogLevelEnvVar, logLevel)
			client := &http.Client{Transport: newRPCTransport(true)}
			resp, err := client.Get(server.URL)
			require.NoError(t, err)
			_ = resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}

	t.Run("custom http.DefaultTransport does not panic and still verifies certificates", func(t *testing.T) {
		original := http.DefaultTransport
		http.DefaultTransport = &LoggingTransport{Transport: original}
		t.Cleanup(func() { http.DefaultTransport = original })

		client := &http.Client{Transport: NewLoggingTransport()}
		resp, err := client.Get(server.URL)
		if resp != nil {
			_ = resp.Body.Close()
		}
		require.Error(t, err, "expected TLS verification to fail for untrusted certificate")
		require.Contains(t, err.Error(), "certificate")
	})
}

func TestNetworkInsecureSkipVerify(t *testing.T) {
	t.Run("is read from TOML network config", func(t *testing.T) {
		var cfg Config
		err := toml.Unmarshal([]byte(`
[[networks]]
name = "Secure"
urls_secret = ["https://rpc.example.com"]

[[networks]]
name = "SelfSigned"
urls_secret = ["https://localhost:8545"]
insecure_skip_verify = true
`), &cfg)
		require.NoError(t, err)
		require.Len(t, cfg.Networks, 2)
		require.False(t, cfg.Networks[0].InsecureSkipVerify, "expected TLS verification to be enabled by default")
		require.True(t, cfg.Networks[1].InsecureSkipVerify, "expected TLS verification to be disabled for SelfSigned network")
	})

	t.Run("is set on selected network by builder", func(t *testing.T) {
		cfg, err := NewClientBuilder().
			WithRpcUrl("https://localhost:8545").
			WithInsecureSkipVerify(true).
			WithReadOnlyMode().
			BuildConfig()
		require.NoError(t, err)
		require.True(t, cfg.Network.InsecureSkipVerify)
		require.Len(t, cfg.Networks, 1)
		require.True(t, cfg.Networks[0].InsecureSkipVerify)
	})
}
