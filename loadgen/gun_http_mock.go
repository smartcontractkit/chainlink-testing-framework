package loadgen

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

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

// NewHTTPMockGun create an HTTP mock gun
func NewHTTPMockGun(cfg *MockHTTPGunConfig) *MockHTTPGun {
	return &MockHTTPGun{
		client: resty.New(),
		cfg:    cfg,
		Data:   make([]string, 0),
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *MockHTTPGun) Call(l *Generator) CallResult {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.cfg.TargetURL)
	if err != nil {
		return CallResult{Data: result, Error: err}
	}
	if r.Status() != "200 OK" {
		return CallResult{Data: result, Error: errors.New("not 200")}
	}
	return CallResult{Data: result}
}
