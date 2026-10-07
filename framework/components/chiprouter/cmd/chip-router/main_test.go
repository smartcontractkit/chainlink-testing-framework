package main

import (
	"context"
	"sync"
	"testing"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stubPublishClient struct {
	mu      sync.Mutex
	batches int
	err     error
}

func (s *stubPublishClient) Publish(context.Context, *cepb.CloudEvent, ...grpc.CallOption) (*chippb.PublishResponse, error) {
	return &chippb.PublishResponse{}, nil
}

func (s *stubPublishClient) PublishBatch(context.Context, *chippb.CloudEventBatch, ...grpc.CallOption) (*chippb.PublishResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batches++
	return &chippb.PublishResponse{}, s.err
}

func (s *stubPublishClient) batchCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.batches
}

func batchOf(ids ...string) *chippb.CloudEventBatch {
	events := make([]*cepb.CloudEvent, 0, len(ids))
	for _, id := range ids {
		events = append(events, &cepb.CloudEvent{Id: id, Source: "s", SpecVersion: "1.0", Type: "t"})
	}
	return &chippb.CloudEventBatch{Events: events}
}

func TestPublishBatchPopulatesPerEventResults(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{}},
	}}

	resp, err := r.PublishBatch(t.Context(), batchOf("e1", "e2", "e3"))
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	if len(resp.Results) != 3 {
		t.Fatalf("want 3 results, got %d", len(resp.Results))
	}
	for i, want := range []string{"e1", "e2", "e3"} {
		if resp.Results[i].EventId != want {
			t.Errorf("results[%d].EventId = %q, want %q", i, resp.Results[i].EventId, want)
		}
		if resp.Results[i].Error != nil {
			t.Errorf("results[%d].Error = %v, want nil", i, resp.Results[i].Error)
		}
	}
}

func TestPublishBatchNoSubscribersIsUnavailable(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{}}

	_, err := r.PublishBatch(t.Context(), batchOf("e1"))
	if err == nil {
		t.Fatal("want error with no subscribers, got nil")
	}
	if status.Code(err) != codes.Unavailable {
		t.Errorf("want Unavailable, got %v", status.Code(err))
	}
}

func TestPublishBatchAllForwardsFailedReportsPerEventErrors(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{err: context.DeadlineExceeded}},
	}}

	resp, err := r.PublishBatch(t.Context(), batchOf("e1", "e2"))
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("want 2 results, got %d", len(resp.Results))
	}
	for i, res := range resp.Results {
		if res.Error == nil {
			t.Errorf("results[%d].Error = nil, want a forward-failure error", i)
		}
	}
}
