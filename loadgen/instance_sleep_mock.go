package loadgen

import (
	"math/rand"
	"time"
)

// MockInstanceConfig configures a mock instanceTemplate
type MockInstanceConfig struct {
	// FailRatio in percentage, 0-100
	FailRatio int
	// TimeoutRatio in percentage, 0-100
	TimeoutRatio int
	// CallSleep time spent waiting inside a call
	CallSleep time.Duration
}

// MockInstance is a mock instanceTemplate
type MockInstance struct {
	cfg  MockInstanceConfig
	stop chan struct{}
	Data []string
}

// NewMockInstance create a mock instanceTemplate
func NewMockInstance(cfg MockInstanceConfig) MockInstance {
	return MockInstance{
		cfg:  cfg,
		stop: make(chan struct{}, 1),
		Data: make([]string, 0),
	}
}

func (m MockInstance) Clone(l *Generator) Instance {
	return MockInstance{
		cfg:  m.cfg,
		stop: make(chan struct{}, 1),
		Data: make([]string, 0),
	}
}

func (m MockInstance) Run(l *Generator) {
	l.ResponsesWaitGroup.Add(1)
	go func() {
		defer l.ResponsesWaitGroup.Done()
		for {
			select {
			case <-l.ResponsesCtx.Done():
				return
			case <-m.stop:
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

func (m MockInstance) Stop(l *Generator) {
	m.stop <- struct{}{}
}
