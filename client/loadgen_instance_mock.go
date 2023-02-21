package client

import (
	"math/rand"
	"time"
)

// MockInstanceConfig configures a mock instance
type MockInstanceConfig struct {
	// FailRatio in percentage, 0-100
	FailRatio int
	// TimeoutRatio in percentage, 0-100
	TimeoutRatio int
	// CallSleep time spent waiting inside a call
	CallSleep time.Duration
}

// MockInstance is a mock instance
type MockInstance struct {
	cfg  *MockInstanceConfig
	Data []string
}

// NewMockInstance create a mock instance
func NewMockInstance(cfg *MockInstanceConfig) *MockInstance {
	return &MockInstance{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

func (m *MockInstance) Run(data interface{}, ch chan CallResult) {
	go func() {
		for {
			startedAt := time.Now()
			time.Sleep(m.cfg.CallSleep)
			if m.cfg.FailRatio > 0 && m.cfg.FailRatio <= 100 {
				//nolint
				r := rand.Intn(100)
				if r <= m.cfg.FailRatio {
					ch <- CallResult{StartedAt: startedAt, Data: "failedCallData", Error: "error", Failed: true}
				}
			}
			if m.cfg.TimeoutRatio > 0 && m.cfg.TimeoutRatio <= 100 {
				//nolint
				r := rand.Intn(100)
				if r <= m.cfg.TimeoutRatio {
					time.Sleep(m.cfg.CallSleep + 100*time.Millisecond)
					ch <- CallResult{}
				}
			}
			ch <- CallResult{StartedAt: startedAt, Data: "successCallData"}
		}
	}()
}
