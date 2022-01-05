package client

import (
	"fmt"
	"net/http"
	"strings"
)

// MockserverClient mockserver client
type MockserverClient struct {
	*BasicHTTPClient
	Config *MockserverConfig
}

// MockserverConfig holds config information for MockserverClient
type MockserverConfig struct {
	LocalURL   string
	ClusterURL string
}

// NewMockserverClient returns a mockserver client
func NewMockserverClient(cfg *MockserverConfig) *MockserverClient {
	return &MockserverClient{
		Config:          cfg,
		BasicHTTPClient: NewBasicHTTPClient(&http.Client{}, cfg.LocalURL),
	}
}

// SetValuePath sets an int for a path
func (em *MockserverClient) SetValuePath(path string, v int) error {
	sanitizedPath := strings.ReplaceAll(path, "/", "_")
	initializer := HttpInitializer{
		Id:      fmt.Sprintf("%s_mock_id", sanitizedPath),
		Request: HttpRequest{Path: path},
		Response: HttpResponse{Body: AdapterResponse{
			Id:    "",
			Data:  AdapterResult{Result: v},
			Error: nil,
		}},
	}
	initializers := []HttpInitializer{initializer}
	_, err := em.Do(http.MethodPut, "/expectation", &initializers, nil, http.StatusCreated)
	return err
}

// PathSelector represents the json object used to find expectations by path
type PathSelector struct {
	Path string `json:"path"`
}

// HttpRequest represents the httpRequest json object used in the mockserver initializer
type HttpRequest struct {
	Path string `json:"path"`
}

// HttpResponse represents the httpResponse json object used in the mockserver initializer
type HttpResponse struct {
	Body interface{} `json:"body"`
}

// HttpInitializer represents an element of the initializer array used in the mockserver initializer
type HttpInitializer struct {
	Id       string       `json:"id"`
	Request  HttpRequest  `json:"httpRequest"`
	Response HttpResponse `json:"httpResponse"`
}

// For OTPE - weiwatchers

// NodeInfoJSON represents an element of the nodes array used to deliver configs to otpe
type NodeInfoJSON struct {
	ID          string   `json:"id"`
	NodeAddress []string `json:"nodeAddress"`
}

// ContractInfoJSON represents an element of the contracts array used to deliver configs to otpe
type ContractInfoJSON struct {
	ContractAddress string `json:"contractAddress"`
	ContractVersion int    `json:"contractVersion"`
	Path            string `json:"path"`
	Status          string `json:"status"`
}

// For Adapter endpoints

// AdapterResult represents an int result for an adapter
type AdapterResult struct {
	Result int `json:"result"`
}

// AdapterResponse represents a response from an adapter
type AdapterResponse struct {
	Id    string        `json:"id"`
	Data  AdapterResult `json:"data"`
	Error interface{}   `json:"error"`
}
