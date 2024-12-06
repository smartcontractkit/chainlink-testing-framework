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

// NewMockVU creates a new MockVirtualUser with the provided configuration.
// It initializes control structures and prepares data storage.
// Use this function to simulate virtual users for testing decentralized services.
func NewMockVU(cfg *MockVirtualUserConfig) *MockVirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

// Clone returns a copy of the MockVirtualUser with a new VUControl and duplicated configuration.
// It is used to create independent virtual user instances for load testing.
func (m *MockVirtualUser) Clone(_ *Generator) VirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

// Setup initializes the VirtualUser using the provided Generator.
// It prepares necessary resources and returns an error if the setup process fails.
func (m *MockVirtualUser) Setup(_ *Generator) error {
	if m.cfg.SetupFailure {
		return errors.New("setup failure")
	}
	time.Sleep(m.cfg.SetupSleep)
	return nil
}

// Teardown cleans up the VirtualUser by releasing resources and performing necessary shutdown procedures.
// It returns an error if the teardown process fails, allowing callers to handle cleanup failures appropriately.
func (m *MockVirtualUser) Teardown(_ *Generator) error {
	if m.cfg.TeardownFailure {
		return errors.New("teardown failure")
	}
	time.Sleep(m.cfg.TeardownSleep)
	return nil
}

// Call simulates a virtual user's call to the Generator.
// It sends a Response to the Generator's ResponsesChan, which may indicate success, failure, or timeout based on the mock configuration.
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
