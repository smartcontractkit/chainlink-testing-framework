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

// NewMockVU create a mock virtual user
func NewMockVU(cfg *MockVirtualUserConfig) *MockVirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

func (m *MockVirtualUser) Clone(_ *Generator) VirtualUser {
	return &MockVirtualUser{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

func (m *MockVirtualUser) Setup(_ *Generator) error {
	if m.cfg.SetupFailure {
		return errors.New("setup failure")
	}
	time.Sleep(m.cfg.SetupSleep)
	return nil
}

func (m *MockVirtualUser) Teardown(_ *Generator) error {
	if m.cfg.TeardownFailure {
		return errors.New("teardown failure")
	}
	time.Sleep(m.cfg.TeardownSleep)
	return nil
}

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
