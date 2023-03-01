package loadgen

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

func (m *MockInstance) Run(l *Generator) {
	l.ResponsesWaitGroup.Add(1)
	go func() {
		for {
			select {
			// TODO: this is mandatory, we should stop the instance when test is done
			// TODO: wrap this in closure to simplify setup
			case <-l.ResponsesCtx.Done():
				l.ResponsesWaitGroup.Done()
				return
			default:
				startedAt := time.Now()
				time.Sleep(m.cfg.CallSleep)
				if m.cfg.FailRatio > 0 && m.cfg.FailRatio <= 100 {
					//nolint
					r := rand.Intn(100)
					if r <= m.cfg.FailRatio {
						l.ResponsesChan <- CallResult{StartedAt: &startedAt, Data: "failedCallData", Error: "error", Failed: true}
					}
				}
				if m.cfg.TimeoutRatio > 0 && m.cfg.TimeoutRatio <= 100 {
					//nolint
					r := rand.Intn(100)
					if r <= m.cfg.TimeoutRatio {
						time.Sleep(m.cfg.CallSleep + 100*time.Millisecond)
						l.ResponsesChan <- CallResult{StartedAt: &startedAt, Data: "timeoutData", Timeout: true}
					}
				}
				l.ResponsesChan <- CallResult{StartedAt: &startedAt, Data: "successCallData"}
			}
		}
	}()
}
