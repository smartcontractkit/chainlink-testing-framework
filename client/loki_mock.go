package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/loki/clients/pkg/promtail/api"
	lokiClient "github.com/grafana/loki/clients/pkg/promtail/client"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/prometheus/common/model"
)

type PromtailSendResult struct {
	Labels model.LabelSet
	Time   time.Time
	Entry  string
}

type MockPromtailClient struct {
	Results       []PromtailSendResult
	OnHandleEntry api.EntryHandler
}

// ExtendedLokiClient an extended Loki/Promtail client used for testing last results in batch
type ExtendedLokiClient interface {
	lokiClient.Client
	api.EntryHandler
	HandleStruct(ls model.LabelSet, t time.Time, st interface{}) error
	LastHandleResult() PromtailSendResult
	AllHandleResults() []PromtailSendResult
}

func (m *LokiClient) LastHandleResult() PromtailSendResult {
	panic("implement me")
}

func (m *LokiClient) AllHandleResults() []PromtailSendResult {
	panic("implement me")
}

func NewMockPromtailClient() ExtendedLokiClient {
	mc := &MockPromtailClient{
		Results: make([]PromtailSendResult, 0),
	}
	entries := make(chan api.Entry)
	done := make(chan struct{})
	stop := func() {
		close(entries)
		<-done
	}
	go func() {
		defer close(done)
		for e := range entries {
			mc.Results = append(mc.Results, PromtailSendResult{Labels: e.Labels, Time: e.Timestamp, Entry: e.Line})
		}
	}()
	mc.OnHandleEntry = api.NewEntryHandler(entries, stop)
	return mc
}

// Name implements api.EntryHandler
func (c *MockPromtailClient) Name() string { return "" }

// Chan implements api.EntryHandler
func (c *MockPromtailClient) Chan() chan<- api.Entry {
	return c.OnHandleEntry.Chan()
}

// Stop implements api.EntryHandler
func (c *MockPromtailClient) Stop() {}

// StopNow implements api.EntryHandler
func (c *MockPromtailClient) StopNow() {}

func (c *MockPromtailClient) HandleStruct(ls model.LabelSet, t time.Time, st interface{}) error {
	d, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("failed to marshal struct in response: %v", st)
	}
	c.Chan() <- api.Entry{
		Labels: ls,
		Entry: logproto.Entry{
			Timestamp: t,
			Line:      string(d),
		},
	}
	return nil
}

func (c *MockPromtailClient) LastHandleResult() PromtailSendResult {
	time.Sleep(2 * time.Second)
	return c.Results[len(c.Results)-1]
}

func (c *MockPromtailClient) AllHandleResults() []PromtailSendResult {
	return c.Results
}
