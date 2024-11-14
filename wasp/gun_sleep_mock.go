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

// NewMockGun creates a new MockGun instance using the provided configuration.
// It is used to simulate gun behavior for testing or development purposes.
func NewMockGun(cfg *MockGunConfig) *MockGun {
	return &MockGun{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

// Call simulates a request to the Generator.
// Depending on MockGun's configuration, it may succeed, fail, or timeout,
// allowing testing of various response scenarios.
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

// convertResponsesData extracts successful and failed response data from the Generator.
// It returns a slice of successful response strings, OK responses, and failed responses.
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
