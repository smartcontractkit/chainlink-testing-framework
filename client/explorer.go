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

func (em *ExplorerClient) PostAdminNodes(nodeName string) (*NodeCreated, error) {
	em.Header = map[string][]string{
		"x-explore-admin-password": {em.Config.AdminPassword},
		"x-explore-admin-username": {em.Config.AdminUsername},
		"Content-Type": {"application/json"},
	}
	requestBody := &Name{Name: nodeName}
	responseBody := &NodeCreated{}
	log.Info().Str("Explorer URL", em.Config.URL).Msg("Creating node credentials")
	_, err := em.do(http.MethodPost, "/api/v1/admin/nodes", &requestBody, &responseBody, http.StatusCreated)
	return responseBody, err
}

type Name struct {
	Name string `json:"name"`
}

type NodeCreated struct {
	Id        string `json:"id"`
	AccessKey string `json:"accessKey"`
	Secret    string `json:"secret"`
}
