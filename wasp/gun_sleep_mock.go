package wasp

import (
	"math/rand"
	"time"
)

// MockGunConfig configures a mock gun
type MockGunConfig struct {
	// FailRatio in percentage, 0-100
	FailRatio int
	// TimeoutRatio in percentage, 0-100
	TimeoutRatio int
	// CallSleep time spent waiting inside a call
	CallSleep time.Duration
	// InternalStop break the test immediately
	InternalStop bool
}

// MockGun is a mock gun
type MockGun struct {
	cfg  *MockGunConfig
	Data []string
}

// NewMockGun creates a new instance of MockGun with the provided configuration.
// It initializes the Data slice as an empty slice of strings.
// The function returns a pointer to the newly created MockGun instance.
func NewMockGun(cfg *MockGunConfig) *MockGun {
	return &MockGun{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

// Call simulates a request to a service using the provided Generator. It may
// stop the Generator if configured to do so. The function introduces a delay
// based on the configuration and can simulate failures or timeouts with
// specified probabilities. It returns a Response indicating the result of the
// call, which may include success, failure, or timeout information.
func (m *MockGun) Call(l *Generator) *Response {
	if m.cfg.InternalStop {
		l.Stop()
	}
	time.Sleep(m.cfg.CallSleep)
	if m.cfg.FailRatio > 0 && m.cfg.FailRatio <= 100 {
		//nolint
		r := rand.Intn(100)
		if r <= m.cfg.FailRatio {
			return &Response{Data: "failedCallData", Error: "error", Failed: true}
		}
	}
	if m.cfg.TimeoutRatio > 0 && m.cfg.TimeoutRatio <= 100 {
		//nolint
		r := rand.Intn(100)
		if r <= m.cfg.TimeoutRatio {
			time.Sleep(m.cfg.CallSleep + 20*time.Millisecond)
		}
	}
	return &Response{Data: "successCallData"}
}

// convertResponsesData retrieves and converts response data from the Generator.
// It locks and unlocks the necessary mutexes to ensure thread safety while accessing
// the response data. The function returns a slice of strings representing successful
// data, a slice of pointers to Response for successful responses, and a slice of
// pointers to Response for failed responses.
func convertResponsesData(g *Generator) ([]string, []*Response, []*Response) {
	g.responsesData.okDataMu.Lock()
	defer g.responsesData.okDataMu.Unlock()
	g.responsesData.failResponsesMu.Lock()
	defer g.responsesData.failResponsesMu.Unlock()
	ok := make([]string, 0)
	for _, d := range g.responsesData.OKData.Data {
		ok = append(ok, d.(string))
	}
	return ok, g.responsesData.OKResponses.Data, g.responsesData.FailResponses.Data
}
