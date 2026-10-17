package main

import (
	"context"
	"testing"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// stubPublishClient is immutable after construction, so it is safe for the
// router's concurrent forwards. A nil response makes Publish return the
// write-service default for a single event (one accepted result) and
// PublishBatch return a detail-less success (empty results).
type stubPublishClient struct {
	err      error
	response *chippb.PublishResponse
}

func (s *stubPublishClient) Publish(_ context.Context, in *cepb.CloudEvent, _ ...grpc.CallOption) (*chippb.PublishResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.response != nil {
		return s.response, nil
	}
	return &chippb.PublishResponse{Results: []*chippb.PublishResult{{EventId: in.GetId()}}}, nil
}

func (s *stubPublishClient) PublishBatch(context.Context, *chippb.CloudEventBatch, ...grpc.CallOption) (*chippb.PublishResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.response != nil {
		return s.response, nil
	}
	return &chippb.PublishResponse{}, nil
}

func ptr(b bool) *bool { return &b }

func eventOf(id string) *cepb.CloudEvent {
	return &cepb.CloudEvent{Id: id, Source: "s", SpecVersion: "1.0", Type: "t"}
}

func batchOf(ids ...string) *chippb.CloudEventBatch {
	events := make([]*cepb.CloudEvent, 0, len(ids))
	for _, id := range ids {
		events = append(events, eventOf(id))
	}
	return &chippb.CloudEventBatch{Events: events}
}

func rejected(id, reason string) *chippb.PublishResult {
	return &chippb.PublishResult{EventId: id, Error: &chippb.PublishError{Reason: reason}}
}

func requireCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("want %s error, got nil", want)
	}
	if status.Code(err) != want {
		t.Errorf("want %s, got %v", want, status.Code(err))
	}
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

func TestPublishBatchEmptyBatchSucceedsWithoutSubscribers(t *testing.T) {
	// The write-service treats an empty batch as a no-op success before any
	// delivery attempt.
	r := &router{subscribers: map[string]*subscriber{}}

	resp, err := r.PublishBatch(t.Context(), batchOf())
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	if len(resp.Results) != 0 {
		t.Errorf("want 0 results, got %d", len(resp.Results))
	}
}

func TestPublishBatchNoSubscribersFailsWithInternal(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{}}

	_, err := r.PublishBatch(t.Context(), batchOf("e1"))
	requireCode(t, err, codes.Internal)
}

func TestPublishBatchAllForwardsFailedFailsWithInternal(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{err: context.DeadlineExceeded}},
	}}

	_, err := r.PublishBatch(t.Context(), batchOf("e1", "e2"))
	requireCode(t, err, codes.Internal)
}

func TestPublishBatchAggregatesDownstreamRejections(t *testing.T) {
	// One subscriber rejects e2 per-event (nil RPC error); the router must
	// pass that rejection through rather than acknowledging the whole batch.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1"},
				rejected("e2", "invalid event"),
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

func TestPublishBatchRejectionSurvivesAnotherSubscribersRPCFailure(t *testing.T) {
	// A rejects e1 per-event; B's forward fails. e1 is still reported
	// rejected in partial-delivery mode — it was accepted nowhere.
	r := &router{subscribers: map[string]*subscriber{
		"a": {id: "a", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{rejected("e1", "rejected by a")},
		}}},
		"b": {id: "b", client: &stubPublishClient{err: context.DeadlineExceeded}},
	}}

	resp, err := r.PublishBatch(t.Context(), batchOf("e1"))
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].Error == nil {
		t.Errorf("want e1 rejected, got results=%+v", resp.Results)
	}
}

func TestPublishBatchAnAcceptingSubscriberOverridesARejection(t *testing.T) {
	// Two detail-providing subscribers disagree: the event is rejected by one
	// and accepted by the other — acceptance wins.
	r := &router{subscribers: map[string]*subscriber{
		"a": {id: "a", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{rejected("e1", "rejected by a")},
		}}},
		"b": {id: "b", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{{EventId: "e1"}},
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

func TestPublishBatchDetaillessAcceptanceOverridesRejections(t *testing.T) {
	// A rejects both events per-event; B accepts the whole batch without
	// per-event detail. The batch was delivered to B, so both events are
	// acknowledged regardless of completion order.
	r := &router{subscribers: map[string]*subscriber{
		"a": {id: "a", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				rejected("e1", "rejected by a"),
				rejected("e2", "rejected by a"),
			},
		}}},
		"b": {id: "b", client: &stubPublishClient{}},
	}}

	resp, err := r.PublishBatch(t.Context(), batchOf("e1", "e2"))
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("want 2 results, got %d", len(resp.Results))
	}
	for i, res := range resp.Results {
		if res.Error != nil {
			t.Errorf("results[%d].Error = %v, want nil (batch accepted by b)", i, res.Error)
		}
	}
}

func TestPublishBatchMergeExtendsBeyondFirstSubscriberDetail(t *testing.T) {
	// A returns detail for only one of two events (contract violation); B
	// returns full detail rejecting e2. e2's rejection must survive: detail
	// is sized to the batch, not to the first subscriber's response.
	r := &router{subscribers: map[string]*subscriber{
		"a": {id: "a", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{{EventId: "e1"}},
		}}},
		"b": {id: "b", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1"},
				rejected("e2", "rejected by b"),
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
		t.Errorf("results[0].Error = %v, want nil", resp.Results[0].Error)
	}
	if resp.Results[1].Error == nil {
		t.Error("results[1].Error = nil, want the rejection from b passed through")
	}
}

func TestPublishBatchTransactionalPartialRejectionFailsRPC(t *testing.T) {
	// Transactional callers resolve the whole batch from the RPC outcome
	// alone, so a partial rejection must fail the RPC, not ack everything.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1"},
				rejected("e2", "rejected"),
			},
		}}},
	}}

	batch := batchOf("e1", "e2")
	batch.Options = &chippb.PublishOptions{TransactionEnabled: ptr(true)}
	_, err := r.PublishBatch(t.Context(), batch)
	requireCode(t, err, codes.Internal)
}

func TestPublishBatchTransactionalAllRejectedByDownstreamFailsRPC(t *testing.T) {
	// A nil downstream RPC error with every event rejected must NOT count as
	// accepted: a transactional caller would otherwise delete them undelivered.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				rejected("e1", "rejected"),
				rejected("e2", "rejected"),
			},
		}}},
	}}

	batch := batchOf("e1", "e2")
	batch.Options = &chippb.PublishOptions{TransactionEnabled: ptr(true)}
	_, err := r.PublishBatch(t.Context(), batch)
	requireCode(t, err, codes.Internal)
}

func TestPublishBatchTransactionalAllAcceptedPasses(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1"},
				{EventId: "e2"},
			},
		}}},
	}}

	batch := batchOf("e1", "e2")
	batch.Options = &chippb.PublishOptions{TransactionEnabled: ptr(true)}
	resp, err := r.PublishBatch(t.Context(), batch)
	if err != nil {
		t.Fatalf("PublishBatch: %v", err)
	}
	for i, res := range resp.Results {
		if res.Error != nil {
			t.Errorf("results[%d].Error = %v, want nil", i, res.Error)
		}
	}
}

func TestPublishReturnsOneResult(t *testing.T) {
	// The write-service returns exactly one result for a single event.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{}},
	}}

	resp, err := r.Publish(t.Context(), eventOf("e1"))
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("want 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].EventId != "e1" {
		t.Errorf("results[0].EventId = %q, want %q", resp.Results[0].EventId, "e1")
	}
	if resp.Results[0].Error != nil {
		t.Errorf("results[0].Error = %v, want nil", resp.Results[0].Error)
	}
}

func TestPublishNoSubscribersFailsWithInternal(t *testing.T) {
	// Acks without delivery would let a durable store delete the event undelivered.
	r := &router{subscribers: map[string]*subscriber{}}

	_, err := r.Publish(t.Context(), eventOf("e1"))
	requireCode(t, err, codes.Internal)
}

func TestPublishAllForwardsFailedFailsWithInternal(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{err: context.DeadlineExceeded}},
	}}

	_, err := r.Publish(t.Context(), eventOf("e1"))
	requireCode(t, err, codes.Internal)
}

func TestPublishAllSubscribersRejectedFailsWithInternal(t *testing.T) {
	// Single-event Publish is all-or-nothing (the write-service treats it as
	// transactional): a per-event rejection from every subscriber fails the RPC.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{rejected("e1", "invalid event")},
		}}},
	}}

	_, err := r.Publish(t.Context(), eventOf("e1"))
	requireCode(t, err, codes.Internal)
}

func TestPublishDetaillessAcceptanceOverridesARejection(t *testing.T) {
	// A rejects the event per-event; B accepts without per-event detail —
	// the event was delivered to B, so Publish acknowledges it.
	r := &router{subscribers: map[string]*subscriber{
		"a": {id: "a", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{rejected("e1", "rejected by a")},
		}}},
		"b": {id: "b", client: &stubPublishClient{}},
	}}

	resp, err := r.Publish(t.Context(), eventOf("e1"))
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].Error != nil {
		t.Errorf("want e1 accepted, got results=%+v", resp.Results)
	}
}

func TestPingMatchesWriteService(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{}}

	resp, err := r.Ping(context.Background(), &chippb.EmptyRequest{})
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if resp.GetMessage() != "pong" {
		t.Errorf("Ping message = %q, want %q", resp.GetMessage(), "pong")
	}
}
