package wasp

import "github.com/go-resty/resty/v2"

// MockHTTPGunConfig configures a mock HTTP gun
type MockHTTPGunConfig struct {
	TargetURL string
}

// MockHTTPGun is a mock gun
type MockHTTPGun struct {
	client *resty.Client
	cfg    *MockHTTPGunConfig
	Data   []string
}

// NewHTTPMockGun initializes a MockHTTPGun with the given configuration.
// It sets up the HTTP client and data storage, enabling simulated HTTP interactions for testing.
func NewHTTPMockGun(cfg *MockHTTPGunConfig) *MockHTTPGun {
	return &MockHTTPGun{
		client: resty.New(),
		cfg:    cfg,
		Data:   make([]string, 0),
	}
}

// Call sends an HTTP GET request to the configured target URL and returns the response data.
// It is used to simulate HTTP calls for testing or load generation purposes.
func (m *MockHTTPGun) Call(l *Generator) *Response {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.cfg.TargetURL)
	if err != nil {
		return &Response{Data: result, Error: err.Error()}
	}
	if r.Status() != "200 OK" {
		return &Response{Data: result, Error: "not 200"}
	}
	return &Response{Data: result}
}
