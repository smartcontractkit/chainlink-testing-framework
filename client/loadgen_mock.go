package client

import (
	"errors"
	"time"
)

type MockGunConfig struct {
	Fail        bool
	CallSleep   time.Duration
	VerifySleep time.Duration
}

type MockGun struct {
	cfg  *MockGunConfig
	Data []string
}

func NewMockGun(cfg *MockGunConfig) *MockGun {
	return &MockGun{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

func (m *MockGun) Call(data interface{}) CallResult {
	time.Sleep(m.cfg.CallSleep)
	if m.cfg.Fail {
		return CallResult{Data: "failedCallData", Error: errors.New("error")}
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
