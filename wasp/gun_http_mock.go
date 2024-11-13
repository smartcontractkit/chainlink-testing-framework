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

// NewHTTPMockGun creates a new instance of MockHTTPGun with the provided configuration.
// It initializes a Resty client and an empty data slice. The function returns a pointer
// to the newly created MockHTTPGun.
func NewHTTPMockGun(cfg *MockHTTPGunConfig) *MockHTTPGun {
	return &MockHTTPGun{
		client: resty.New(),
		cfg:    cfg,
		Data:   make([]string, 0),
	}
}

// Call sends an HTTP GET request to the target URL specified in the MockHTTPGun's configuration.
// It returns a Response containing the result of the request or an error if the request fails or the status is not "200 OK".
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
