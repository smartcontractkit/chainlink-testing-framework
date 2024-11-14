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

// NewHTTPMockGun initializes and returns a new MockHTTPGun instance.
// It sets up the REST client using resty and applies the provided configuration.
// The Data field is initialized as an empty slice of strings.
func NewHTTPMockGun(cfg *MockHTTPGunConfig) *MockHTTPGun {
	return &MockHTTPGun{
		client: resty.New(),
		cfg:    cfg,
		Data:   make([]string, 0),
	}
}

// Call makes an HTTP GET request to the target URL using the provided Generator.
// It returns a Response containing the retrieved data and any error encountered during the request.
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
