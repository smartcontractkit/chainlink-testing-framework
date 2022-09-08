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
		return CallResult{"failedCallData", errors.New("error")}
	}
	return CallResult{"successCallData", nil}
}

func (m *MockGun) CollectData() interface{} {
	return m.Data
}

func convertResponsesData(rd *ResponseData) ([]string, []string) {
	ok, fail := make([]string, 0), make([]string, 0)
	for _, d := range rd.OKData {
		ok = append(ok, d.(string))
	}
	for _, d := range rd.FailData {
		if d.Data != nil {
			fail = append(fail, d.Data.(string))
		}
		fail = append(fail, d.Error.Error())
	}
	return ok, fail
}
