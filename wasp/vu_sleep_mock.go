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

// NewMockVU creates and returns a new instance of MockVirtualUser initialized with the provided configuration. 
// It sets up a VUControl for managing virtual user operations and initializes an empty slice for storing data. 
// The returned MockVirtualUser can be used to simulate virtual user behavior in a testing environment.
func NewMockVU(cfg *MockVirtualUserConfig) *MockVirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

// Clone creates a new instance of MockVirtualUser, initializing it with a new VUControl and copying the configuration from the original instance. 
// The Data field is also initialized as an empty slice of strings. 
// This function is primarily used to generate multiple virtual user instances for load testing scenarios. 
// The returned VirtualUser can be utilized independently of the original instance.
func (m *MockVirtualUser) Clone(_ *Generator) VirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

// Setup initializes the virtual user with the provided generator. 
// It returns an error if the setup fails, which can occur based on the configuration settings. 
// If the setup is successful, it may introduce a delay as specified in the configuration. 
// The function is intended to be called in a context where a timeout may be necessary, 
// allowing for graceful handling of setup failures or timeouts.
func (m *MockVirtualUser) Setup(_ *Generator) error {
	if m.cfg.SetupFailure {
		return errors.New("setup failure")
	}
	time.Sleep(m.cfg.SetupSleep)
	return nil
}

// Teardown performs cleanup operations for the virtual user after the generator has completed its tasks. 
// It returns an error if the teardown process fails, which can be configured through the virtual user's settings. 
// If the teardown is successful, it will return nil after a specified sleep duration, allowing for any necessary 
// delays before the function completes. This function is typically called in a separate goroutine to avoid 
// blocking the main execution flow, and it can be subject to a timeout to ensure it does not hang indefinitely.
func (m *MockVirtualUser) Teardown(_ *Generator) error {
	if m.cfg.TeardownFailure {
		return errors.New("teardown failure")
	}
	time.Sleep(m.cfg.TeardownSleep)
	return nil
}

// Call simulates a virtual user making a call to the provided generator. 
// It introduces a delay based on the configuration and may randomly fail or timeout 
// according to the specified fail and timeout ratios. 
// The function sends a response to the generator's ResponsesChan, indicating 
// whether the call was successful or failed, along with the timestamp of when 
// the call started. If the call fails, it includes an error message in the response. 
// If a timeout occurs, the function will also handle that by sending a timeout response.
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
