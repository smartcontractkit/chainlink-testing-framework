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
	mu       sync.Mutex
	batches  int
	err      error
	response *chippb.PublishResponse
}

func (s *stubPublishClient) Publish(context.Context, *cepb.CloudEvent, ...grpc.CallOption) (*chippb.PublishResponse, error) {
	return &chippb.PublishResponse{}, nil
}

func (s *stubPublishClient) PublishBatch(context.Context, *chippb.CloudEventBatch, ...grpc.CallOption) (*chippb.PublishResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batches++
	if s.err != nil {
		return nil, s.err
	}
	if s.response != nil {
		return s.response, nil
	}
	return &chippb.PublishResponse{}, nil
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

func TestPublishBatchAllForwardsFailedIsUnavailable(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{err: context.DeadlineExceeded}},
	}}

	// A nil RPC error with per-event errors would still acknowledge the whole
	// batch for transactional callers (they resolve delivery from the RPC
	// outcome alone), so the router must fail the RPC instead.
	_, err := r.PublishBatch(t.Context(), batchOf("e1", "e2"))
	if err == nil {
		t.Fatal("want error when every forward fails, got nil")
	}
	if status.Code(err) != codes.Unavailable {
		t.Errorf("want Unavailable, got %v", status.Code(err))
	}
}

func TestPublishBatchAggregatesDownstreamRejections(t *testing.T) {
	// One subscriber rejects e2 per-event (nil RPC error); the router must
	// pass that rejection through rather than acknowledging the whole batch.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1"},
				{EventId: "e2", Error: &chippb.PublishError{
					ErrorCode: chippb.PublishErrorCode_PUBLISH_ERROR_CODE_VALIDATION_FAILED,
					Reason:    "invalid event",
				}},
			},
		}}},
	}}

	resp, err := r.PublishBatch(t.Context(), batchOf("e1", "e2"))
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("want 2 results, got %d", len(resp.Results))
	}
	if resp.Results[0].Error != nil {
		t.Errorf("results[0].Error = %v, want nil (accepted)", resp.Results[0].Error)
	}
	if resp.Results[1].Error == nil {
		t.Error("results[1].Error = nil, want the downstream rejection passed through")
	}
}

func TestPublishBatchAnAcceptingSubscriberOverridesADetaillessRejection(t *testing.T) {
	// Two detail-providing subscribers disagree: the event is rejected by one
	// and accepted by the other — acceptance wins.
	r := &router{subscribers: map[string]*subscriber{
		"a": {id: "a", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1", Error: &chippb.PublishError{Reason: "rejected by a"}},
			},
		}}},
		"b": {id: "b", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1"},
			},
		}}},
	}}

	resp, err := r.PublishBatch(t.Context(), batchOf("e1"))
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].Error != nil {
		t.Errorf("want e1 accepted, got results=%+v", resp.Results)
	}
}
