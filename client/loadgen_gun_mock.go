package client

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
}

// MockGun is a mock gun
type MockGun struct {
	cfg  *MockGunConfig
	Data []string
}

// NewMockGun create a mock gun
func NewMockGun(cfg *MockGunConfig) *MockGun {
	return &MockGun{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

func (m *MockGun) Call(data interface{}) CallResult {
	time.Sleep(m.cfg.CallSleep)
	if m.cfg.FailRatio > 0 && m.cfg.FailRatio <= 100 {
		//nolint
		r := rand.Intn(100)
		if r <= m.cfg.FailRatio {
			return CallResult{Data: "failedCallData", Error: "error", Failed: true}
		}
	}
	if m.cfg.TimeoutRatio > 0 && m.cfg.TimeoutRatio <= 100 {
		//nolint
		r := rand.Intn(100)
		if r <= m.cfg.TimeoutRatio {
			time.Sleep(m.cfg.CallSleep + 100*time.Millisecond)
			return CallResult{}
		}
	}
	return CallResult{Data: "successCallData"}
}

func convertResponsesData(rd *ResponseData) ([]string, []CallResult, []CallResult) {
	ok := make([]string, 0)
	for _, d := range rd.OKData {
		ok = append(ok, d.(string))
	}
	return ok, rd.OKResponses, rd.FailResponses
}
