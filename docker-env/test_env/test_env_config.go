package test_env

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink-testing-framework/docker-env/types/envcommon"
)

type TestEnvConfig struct {
	Networks   []string         `json:"networks"`
	Geth       GethConfig       `json:"geth"`
	MockServer MockServerConfig `json:"mockserver"`
	Nodes      []ClNodeConfig   `json:"nodes"`
}

type MockServerConfig struct {
	ContainerName string   `json:"container_name"`
	EAMockUrls    []string `json:"external_adapters_mock_urls"`
}

type GethConfig struct {
	ContainerName string `json:"container_name"`
}

type ClNodeConfig struct {
	NodeContainerName string `json:"container_name"`
	DbContainerName   string `json:"db_container_name"`
}

func NewTestEnvConfigFromFile(path string) (*TestEnvConfig, error) {
	c := &TestEnvConfig{}
	err := envcommon.ParseJSONFile(path, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *TestEnvConfig) Json() string {
	b, _ := json.Marshal(c)
	return string(b)
}
