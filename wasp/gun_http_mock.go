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

// NewHTTPMockGun initializes a new instance of MockHTTPGun using the provided configuration.
// It sets up an HTTP client and prepares an empty slice to store data. 
// The returned MockHTTPGun can be used for testing purposes, simulating HTTP interactions without making real network calls.
func NewHTTPMockGun(cfg *MockHTTPGunConfig) *MockHTTPGun {
	return &MockHTTPGun{
		client: resty.New(),
		cfg:    cfg,
		Data:   make([]string, 0),
	}
}

// Call executes an HTTP GET request to the target URL configured in the MockHTTPGun. 
// It returns a Response containing the result of the request. If the request encounters 
// an error or does not return a status of "200 OK", the Response will include the error 
// message. If the request is successful, the Response will contain the data retrieved 
// from the target URL.
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
