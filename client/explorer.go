package client

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"net/http"
)

type ExplorerClient struct {
	*BasicHTTPClient
	Config *config.ExplorerConfig
}

// NewExplorerClient creates a new explorer mock client
func NewExplorerClient(cfg *config.ExplorerConfig) *ExplorerClient {
	return &ExplorerClient{
		Config:          cfg,
		BasicHTTPClient: NewBasicHTTPClient(&http.Client{}, cfg.URL),
	}
}

// PostAdminNodes is used to exercise the POST /api/v1/admin/nodes endpoint
// This endpoint is used to create access keys for nodes
func (em *ExplorerClient) PostAdminNodes(nodeName string) (NodeAccessKeys, error) {
	em.Header = map[string][]string{
		"x-explore-admin-password": {em.Config.AdminPassword},
		"x-explore-admin-username": {em.Config.AdminUsername},
		"Content-Type":             {"application/json"},
	}
	requestBody := &Name{Name: nodeName}
	responseBody := NodeAccessKeys{}
	log.Info().Str("Explorer URL", em.Config.URL).Msg("Creating node credentials")
	_, err := em.do(http.MethodPost, "/api/v1/admin/nodes", &requestBody, &responseBody, http.StatusCreated)
	return responseBody, err
}

type Name struct {
	Name string `json:"name"`
}

type NodeAccessKeys struct {
	ID        string `mapstructure:"id" yaml:"id"`
	AccessKey string `mapstructure:"accesKey" yaml:"accessKey"`
	Secret    string `mapstructure:"secret" yaml:"secret"`
}
