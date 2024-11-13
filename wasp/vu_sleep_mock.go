package wasp

import (
	"errors"
	"math/rand"
	"time"
)

// MockVirtualUserConfig configures a mock virtual user
type MockVirtualUserConfig struct {
	// FailRatio in percentage, 0-100
	FailRatio int
	// TimeoutRatio in percentage, 0-100
	TimeoutRatio int
	// CallSleep time spent waiting inside a call
	CallSleep       time.Duration
	SetupSleep      time.Duration
	SetupFailure    bool
	TeardownSleep   time.Duration
	TeardownFailure bool
}

// MockVirtualUser is a mock virtual user
type MockVirtualUser struct {
	*VUControl
	cfg  *MockVirtualUserConfig
	Data []string
}

// NewMockVU creates a new instance of MockVirtualUser with the provided configuration.
// It initializes the VUControl using NewVUControl and sets up an empty data slice.
func NewMockVU(cfg *MockVirtualUserConfig) *MockVirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

// Clone creates a new instance of MockVirtualUser with a fresh VUControl and an empty data slice.
// It retains the configuration from the original MockVirtualUser. This function is typically used
// to generate multiple virtual users in scenarios where load testing requires concurrent execution.
func (m *MockVirtualUser) Clone(_ *Generator) VirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

// Setup initializes the MockVirtualUser using the provided Generator.
// It simulates a setup process that may fail based on the configuration.
// If the setup fails, it returns an error. Otherwise, it completes after a specified sleep duration.
func (m *MockVirtualUser) Setup(_ *Generator) error {
	if m.cfg.SetupFailure {
		return errors.New("setup failure")
	}
	time.Sleep(m.cfg.SetupSleep)
	return nil
}

// Teardown performs the teardown process for a MockVirtualUser. 
// It returns an error if the teardown is configured to fail. 
// Otherwise, it waits for a specified duration before completing successfully.
func (m *MockVirtualUser) Teardown(_ *Generator) error {
	if m.cfg.TeardownFailure {
		return errors.New("teardown failure")
	}
	time.Sleep(m.cfg.TeardownSleep)
	return nil
}

// Call simulates a virtual user making a request to the Generator. 
// It introduces a delay based on the configured CallSleep duration. 
// The function can simulate failures and timeouts based on the configured 
// FailRatio and TimeoutRatio, respectively. It sends a Response to the 
// Generator's ResponsesChan, indicating success or failure of the call.
func (m *MockVirtualUser) Call(l *Generator) {
	startedAt := time.Now()
	time.Sleep(m.cfg.CallSleep)
	if m.cfg.FailRatio > 0 && m.cfg.FailRatio <= 100 {
		//nolint
		r := rand.Intn(100)
		if r <= m.cfg.FailRatio {
			l.ResponsesChan <- &Response{StartedAt: &startedAt, Data: "failedCallData", Error: "error", Failed: true}
		}
	}
	if m.cfg.TimeoutRatio > 0 && m.cfg.TimeoutRatio <= 100 {
		//nolint
		r := rand.Intn(100)
		if r <= m.cfg.TimeoutRatio {
			time.Sleep(m.cfg.CallSleep + 20*time.Millisecond)
		}
	}
	l.ResponsesChan <- &Response{StartedAt: &startedAt, Data: "successCallData"}
}
