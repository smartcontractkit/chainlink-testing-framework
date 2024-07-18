package grafana

import (
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"strings"
)

const (
	EnvVarHTTPHeaders = "CTF_HTTP_HEADERS"
)

var (
	ErrInvalidHeaders = errors.New("invalid RPC headers, format should be 'k=v,k=v', no trailing comma")
)

// ReadEnvHTTPHeaders reads custom RPC headers from env vars
func ReadEnvHTTPHeaders(logger zerolog.Logger) (http.Header, error) {
	hm := http.Header{}
	customHeader := os.Getenv(EnvVarHTTPHeaders)
	if customHeader == "" {
		return nil, nil
	}
	headers := strings.Split(customHeader, ",")
	for _, h := range headers {
		headerKV := strings.Split(h, "=")
		if len(headerKV) != 2 {
			return nil, ErrInvalidHeaders
		}
		hm.Set(headerKV[0], headerKV[1])
	}
	logger.Debug().Msgf("Using custom RPC headers: %s", hm)
	return hm, nil
}

// HeaderToMultiValueFormat converts header to multi-value format
func HeaderToMultiValueFormat(h http.Header) map[string]string {
	hm := make(map[string]string)
	for k, v := range h {
		hm[k] = strings.Join(v, ",")
	}
	return hm
}
