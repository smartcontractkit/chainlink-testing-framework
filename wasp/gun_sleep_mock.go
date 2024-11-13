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

// NewMockGun creates and returns a new instance of MockGun initialized with the provided configuration. 
// It sets up the internal data structure to hold string data, starting with an empty slice. 
// The returned MockGun can be used to simulate gun-related operations in a testing environment.
func NewMockGun(cfg *MockGunConfig) *MockGun {
	return &MockGun{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

// Call executes a request using the provided Generator and returns a Response. 
// It may simulate a failure or a timeout based on the configuration settings of the MockGun. 
// If the InternalStop configuration is enabled, it will stop the Generator before proceeding. 
// The function will also respect the CallSleep duration specified in the configuration. 
// The returned Response contains data indicating success or failure, along with an error message if applicable. 
// If a timeout occurs, the Response will reflect that with a Timeout flag set to true.
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

// convertResponsesData processes the generator's response data and returns three values: 
// a slice of strings containing the successfully processed data, 
// a slice of pointers to Response objects representing successful responses, 
// and a slice of pointers to Response objects representing failed responses. 
// The function ensures thread safety by locking the necessary mutexes during data access.
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
