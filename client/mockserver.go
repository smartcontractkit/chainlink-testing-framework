package client

import (
	"fmt"
	resty2 "github.com/smartcontractkit/chainlink-testing-framework/utils/resty"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/mockserver"
)

// MockserverClient mockserver client
type MockserverClient struct {
	APIClient *resty.Client
	Config    *MockserverConfig
}

// MockserverConfig holds config information for MockserverClient
type MockserverConfig struct {
	LocalURL   string
	ClusterURL string
}

// ConnectMockServer creates a connection to a deployed mockserver in the environment
func ConnectMockServer(e *environment.Environment) (*MockserverClient, error) {
	c, err := NewMockserverClient(&MockserverConfig{
		LocalURL:   e.URLs[mockserver.LocalURLsKey][0],
		ClusterURL: e.URLs[mockserver.InternalURLsKey][0],
	})
	return c, err
}

// ConnectMockServerURL creates a connection to a mockserver at a given url, should only be used for inside K8s tests
func ConnectMockServerURL(url string) (*MockserverClient, error) {
	c, err := NewMockserverClient(&MockserverConfig{
		LocalURL:   url,
		ClusterURL: url,
	})
	return c, err
}

// NewMockserverClient returns a mockserver client
func NewMockserverClient(cfg *MockserverConfig) (*MockserverClient, error) {
	log.Debug().Str("Local URL", cfg.LocalURL).Str("Remote URL", cfg.ClusterURL).Msg("Connected to MockServer")
	r, err := resty2.NewDefaultResty()
	if err != nil {
		return nil, err
	}
	return &MockserverClient{
		Config:    cfg,
		APIClient: r.SetBaseURL(cfg.LocalURL),
	}, nil
}

// PutExpectations sets the expectations (i.e. mocked responses)
func (em *MockserverClient) PutExpectations(body interface{}) error {
	resp, err := em.APIClient.R().SetBody(body).Put("/expectation")
	if resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("Unexpected Status Code. Expected %d; Got %d", http.StatusCreated, resp.StatusCode())
	}
	return err
}

// ClearExpectation clears expectations
func (em *MockserverClient) ClearExpectation(body interface{}) error {
	resp, err := em.APIClient.R().SetBody(body).Put("/clear")
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("Unexpected Status Code. Expected %d; Got %d", http.StatusOK, resp.StatusCode())
	}
	return err
}

// SetRandomValuePath sets a random int value for a path
func (em *MockserverClient) SetRandomValuePath(path string) error {
	sanitizedPath := strings.ReplaceAll(path, "/", "_")
	log.Debug().Str("ID", fmt.Sprintf("%s_mock_id", sanitizedPath)).
		Str("Path", path).
		Msg("Setting Random Value Mock Server Path")
	initializer := HttpInitializerTemplate{
		Id:      fmt.Sprintf("%s_mock_id", sanitizedPath),
		Request: HttpRequest{Path: path},
		Response: HttpResponseTemplate{
			Template:     "return { statusCode: 200, body: JSON.stringify({id: '', error: null, data: { result: Math.floor(Math.random() * (1000 - 900) + 900) } }) }",
			TemplateType: "JAVASCRIPT",
		},
	}
	initializers := []HttpInitializerTemplate{initializer}
	resp, err := em.APIClient.R().SetBody(&initializers).Put("/expectation")
	if resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("status code expected %d got %d", http.StatusCreated, resp.StatusCode())
	}
	return err
}

// SetValuePath sets an int for a path
func (em *MockserverClient) SetValuePath(path string, v int) error {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	sanitizedPath := strings.ReplaceAll(path, "/", "_")
	log.Debug().Str("ID", fmt.Sprintf("%s_mock_id", sanitizedPath)).
		Str("URL", em.APIClient.BaseURL).
		Str("Path", path).
		Int("Value", v).
		Msg("Setting Mock Server Path")
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
	resp, err := em.APIClient.R().SetBody(&initializers).Put("/expectation")
	if resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("status code expected %d got %d", http.StatusCreated, resp.StatusCode())
	}
	return err
}

// SetAnyValuePath sets any type of value for a path
func (em *MockserverClient) SetAnyValuePath(path string, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	sanitizedPath := strings.ReplaceAll(path, "/", "_")
	id := fmt.Sprintf("%s_mock_id", sanitizedPath)
	log.Debug().Str("ID", id).
		Str("Path", path).
		Interface("Value", v).
		Msg("Setting Mock Server Path")
	initializer := HttpInitializer{
		Id:      id,
		Request: HttpRequest{Path: path},
		Response: HttpResponse{
			Body: AdapterResponse{
				Id: "",
				Data: AdapterResult{
					Result: v,
				},
				Error: nil,
			},
		},
	}
	initializers := []HttpInitializer{initializer}
	resp, err := em.APIClient.R().SetBody(&initializers).Put("/expectation")
	if resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("status code expected %d got %d", http.StatusCreated, resp.StatusCode())
	}
	return err
}

func (em *MockserverClient) SetAnyValueResponse(path string, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	sanitizedPath := strings.ReplaceAll(path, "/", "_")
	id := fmt.Sprintf("%s_mock_id", sanitizedPath)
	log.Debug().Str("ID", id).
		Str("Path", path).
		Interface("Value", v).
		Msg("Setting Mock Server Path")
	initializer := HttpInitializer{
		Id:      id,
		Request: HttpRequest{Path: path},
		Response: HttpResponse{
			Body: v,
		},
	}
	initializers := []HttpInitializer{initializer}
	resp, err := em.APIClient.R().SetBody(&initializers).Put("/expectation")
	if resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("status code expected %d got %d", http.StatusCreated, resp.StatusCode())
	}
	return err
}

// SetStringValuePath sets a string value for a path and returns it as a raw string
func (em *MockserverClient) SetStringValuePath(path string, stringValue string) error {
	sanitizedPath := strings.ReplaceAll(path, "/", "_")
	id := fmt.Sprintf("%s_mock_id", sanitizedPath)
	log.Debug().Str("ID", id).
		Str("Path", path).
		Msg("Setting Mock Server String Path")

	initializer := HttpInitializer{
		Id:      id,
		Request: HttpRequest{Path: path},
		Response: HttpResponse{
			Body: stringValue,
		},
	}

	initializers := []HttpInitializer{initializer}
	resp, err := em.APIClient.R().SetBody(&initializers).Put("/expectation")
	if resp.StatusCode() != http.StatusCreated {
		err = fmt.Errorf("status code expected %d got %d", http.StatusCreated, resp.StatusCode())
	}
	return err
}

// LocalURL returns the local url of the mockserver
func (em *MockserverClient) LocalURL() string {
	return em.Config.LocalURL
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

// HttpResponse represents the httpResponse json object used in the mockserver initializer
type HttpResponseTemplate struct {
	Template     string `json:"template"`
	TemplateType string `json:"templateType"`
}

// HttpInitializer represents an element of the initializer array used in the mockserver initializer
type HttpInitializerTemplate struct {
	Id       string               `json:"id"`
	Request  HttpRequest          `json:"httpRequest"`
	Response HttpResponseTemplate `json:"httpResponseTemplate"`
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
	Result interface{} `json:"result"`
}

// AdapterResponse represents a response from an adapter
type AdapterResponse struct {
	Id    string        `json:"id"`
	Data  AdapterResult `json:"data"`
	Error interface{}   `json:"error"`
}
