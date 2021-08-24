package client

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/explorer"
	"net/http"
)

type ExplorerClient struct {
	*BasicHTTPClient
	Config *config.ExplorerMockConfig
}

// NewExplorerMockClient creates a new explorer mock client
func NewExplorerMockClient(cfg *config.ExplorerMockConfig) *ExplorerClient {
	return &ExplorerClient{
		Config:          cfg,
		BasicHTTPClient: NewBasicHTTPClient(&http.Client{}, cfg.URL),
	}
}

// Count get explorer messages and telemetry counts
func (em *ExplorerClient) Count() (*explorer.MessagesCount, error) {
	mc := &explorer.MessagesCount{}
	log.Info().Str("Explorer URL", em.Config.URL).Msg("Checking explorer telemetry messages")
	_, err := em.do(http.MethodGet, "/count", nil, &mc, http.StatusOK)
	return mc, err
}

// Messages get explorer messages and telemetry data
func (em *ExplorerClient) Messages() (*explorer.Messages, error) {
	ms := &explorer.Messages{}
	log.Info().Str("Explorer URL", em.Config.URL).Msg("Checking explorer telemetry messages")
	_, err := em.do(http.MethodGet, "/messages", nil, &ms, http.StatusOK)
	return ms, err
}
