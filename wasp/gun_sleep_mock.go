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

// NewMockGun initializes a new MockGun with the provided configuration.
// It sets the cfg field and initializes Data as an empty slice of strings.
// Returns a pointer to the MockGun instance.
func NewMockGun(cfg *MockGunConfig) *MockGun {
	return &MockGun{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

// Call performs a simulated service call using the provided Generator.
// Depending on the MockGun's configuration, it may induce delays, fail the call, or cause a timeout.
// It returns a Response that reflects the outcome of the simulated call.
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

// convertResponsesData safely retrieves the OK data as a slice of strings along with the corresponding OK and failed responses.
// It locks the necessary mutexes to ensure thread-safe access to the responses data.
// The function returns the extracted OK data, the OK responses, and the failed responses.
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
