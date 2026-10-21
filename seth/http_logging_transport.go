package seth

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// InsecureSkipVerifyEnvVar disables TLS certificate verification for RPC connections when set to "true".
// Use only in testing environments with self-signed certificates.
const InsecureSkipVerifyEnvVar = "SETH_INSECURE_SKIP_VERIFY"

// LoggingTransport is a custom transport to log requests and responses
type LoggingTransport struct {
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface
func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		L.Warn().Err(err).Msg("Error dumping request")
	} else {
		fmt.Printf("Request:\n%s\n", string(reqDump))
	}

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		L.Warn().Err(err).Msg("Error dumping response")
	} else {
		fmt.Printf("Response:\n%s\n", string(respDump))
	}

	fmt.Printf("Request took %s\n", time.Since(start))
	return resp, nil
}

// NewLoggingTransport creates a new transport for RPC HTTP connections. TLS certificates are verified
// unless InsecureSkipVerifyEnvVar is set to "true". When SETH_LOG_LEVEL is "trace" requests and responses are logged.
func NewLoggingTransport() http.RoundTripper {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if strings.EqualFold(os.Getenv(InsecureSkipVerifyEnvVar), "true") {
		L.Warn().
			Str("EnvVar", InsecureSkipVerifyEnvVar).
			Msg("TLS certificate verification is disabled for RPC connections. Do not use this setting outside of testing environments")
		transport.TLSClientConfig = &tls.Config{
			//nolint:gosec // explicit opt-in for testing environments with self-signed certificates
			InsecureSkipVerify: true,
		}
	}

	if strings.EqualFold(os.Getenv(LogLevelEnvVar), "trace") {
		return &LoggingTransport{Transport: transport}
	}

	return transport
}
