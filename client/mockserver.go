package client

import (
	"net/http"
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

// PutExpectations sets the expectations (i.e. mocked responses)
func (em *MockserverClient) PutExpectations(body interface{}) error {
	_, err := em.do(http.MethodPut, "/expectation", &body, nil, http.StatusCreated)
	return err
}

// ClearExpectation clears expectations
func (em *MockserverClient) ClearExpectation(body interface{}) error {
	_, err := em.do(http.MethodPut, "/clear", &body, nil, http.StatusOK)
	return err
}

// SetVariable sets an int for /variable
func (em *MockserverClient) SetVariable(v int) error {
	pathSelector := PathSelector{Path: "/variable"}
	if err := em.ClearExpectation(pathSelector); err != nil {
		return err
	}
	initializer := HttpInitializer{
		Request: HttpRequest{Path: "/variable"},
		Response: HttpResponse{Body: AdapterResponse{
			Id:    "",
			Data:  AdapterResult{Result: v},
			Error: nil,
		}},
	}
	initializers := []HttpInitializer{initializer}
	if err := em.PutExpectations(initializers); err != nil {
		return err
	}
	return nil
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
