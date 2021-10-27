package client

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

// EIService is external initiator service
type EIService interface {
	URL() string
	Keys() (*Response, error)
}

// EIClient is external initiator client
type EIClient struct {
	*BasicHTTPClient
	Config *EIServiceConfig
}

// URL returns EIService URL
func (t *EIClient) URL() string {
	return t.Config.URL
}

// Keys returns EI public keys
func (t *EIClient) Keys() (*Response, error) {
	specObj := &Response{}
	log.Info().Str("EIService URL", t.Config.URL).Msg("Reading EI public key")
	_, err := t.do(http.MethodGet, "/keys", nil, specObj, http.StatusOK)
	return specObj, err
}

// NewEIServiceClient creates new EIService client
func NewEIServiceClient(c *EIServiceConfig, httpClient *http.Client) (*EIClient, error) {
	return &EIClient{
		BasicHTTPClient: NewBasicHTTPClient(httpClient, c.URL),
	}, nil
}
