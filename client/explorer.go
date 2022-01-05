package client

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// ExplorerClient is used to call Explorer API endpoints
type ExplorerClient struct {
	*BasicHTTPClient
	Config *ExplorerConfig
}

// NewExplorerClient creates a new explorer mock client
func NewExplorerClient(cfg *ExplorerConfig) *ExplorerClient {
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
	_, err := em.Do(http.MethodPost, "/api/v1/admin/nodes", &requestBody, &responseBody, http.StatusCreated)
	return responseBody, err
}

// Name is the body of the request
type Name struct {
	Name string `json:"name"`
}

// NodeAccessKeys is the body of the response
type NodeAccessKeys struct {
	ID        string `mapstructure:"id" yaml:"id"`
	AccessKey string `mapstructure:"accessKey" yaml:"accessKey"`
	Secret    string `mapstructure:"secret" yaml:"secret"`
}

// ExplorerConfig holds config information for ExplorerClient
type ExplorerConfig struct {
	URL           string
	AdminUsername string
	AdminPassword string
}
