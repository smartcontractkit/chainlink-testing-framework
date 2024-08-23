package seth

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

// LoggingTransport is a custom transport to log requests and responses
type LoggingTransport struct {
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface
func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		L.Error().Err(err).Msg("Error dumping request")
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
		L.Error().Err(err).Msg("Error dumping response")
	} else {
		fmt.Printf("Response:\n%s\n", string(respDump))
	}

	fmt.Printf("Request took %s\n", time.Since(start))
	return resp, nil
}

// NewLoggingTransport creates a new logging transport for GAP or default transport
// controlled by SETH_LOG_LEVEL
func NewLoggingTransport() http.RoundTripper {
	if os.Getenv(LogLevelEnvVar) == "debug" {
		return &LoggingTransport{
			// TODO: GAP, add proper certificates
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		return &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}
