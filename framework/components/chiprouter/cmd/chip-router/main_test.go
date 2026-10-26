package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

// requireMessageContains asserts the gRPC status message carries every part, so
// a caller reading only the error can tell what happened and how to fix it.
func requireMessageContains(t *testing.T, err error, parts ...string) {
	t.Helper()
	msg := status.Convert(err).Message()
	for _, p := range parts {
		if !strings.Contains(msg, p) {
			t.Errorf("error message %q does not contain %q", msg, p)
		}
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

func TestPublishBatchResultsMismatchFromOnlySubscriberFailsWithInternal(t *testing.T) {
	// The write-service returns one result per event, so a short list cannot
	// be aligned with the batch: it counts as a failed forward, not as an
	// acknowledgement of the events it does not cover.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", name: "short-sink", endpoint: "10.0.0.5:50051", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{{EventId: "e1"}},
		}}},
	}}

	_, err := r.PublishBatch(t.Context(), batchOf("e1", "e2"))
	requireCode(t, err, codes.Internal)
	// The caller must be able to act on the error alone: which subscriber, what
	// it returned, and what a correct response looks like.
	requireMessageContains(t, err,
		`subscriber "short-sink" (10.0.0.5:50051)`,
		"returned 1 results for 2 events",
		"exactly one PublishResult per event, in input order",
		"or no results to accept the whole batch",
	)
}

func TestPublishBatchResultsMismatchIsIgnoredInFavorOfFullDetail(t *testing.T) {
	// A returns a short list (contract violation) and is dropped; B returns
	// full detail rejecting e2, which must be the outcome.
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

func TestPublishDetaillessSuccessAcceptsEvent(t *testing.T) {
	// chainlink's chip-testsink subscriber returns an empty PublishResponse on
	// success; the router must still acknowledge the event and echo its ID.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{}}},
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
	if got, want := status.Convert(err).Message(), "event rejected by subscribers: invalid event"; got != want {
		t.Errorf("error message = %q, want %q", got, want)
	}
}

func TestPublishRejectionWithoutReasonStillHasMessage(t *testing.T) {
	// A subscriber may reject without a Reason; the caller must not see an
	// empty gRPC error message.
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{rejected("e1", "")},
		}}},
	}}

	_, err := r.Publish(t.Context(), eventOf("e1"))
	requireCode(t, err, codes.Internal)
	if got := status.Convert(err).Message(); got == "" {
		t.Error("error message is empty, want a non-empty message")
	}
}

func TestPublishDetaillessAcceptanceOverridesARejection(t *testing.T) {
	// A rejects the event per-event; B accepts with an empty response (no
	// per-event detail) — the event was delivered to B, so Publish
	// acknowledges it.
	r := &router{subscribers: map[string]*subscriber{
		"a": {id: "a", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{rejected("e1", "rejected by a")},
		}}},
		"b": {id: "b", client: &stubPublishClient{response: &chippb.PublishResponse{}}},
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

func TestPublishBatchNoSubscribersMessageSaysHowToRegister(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{}}

	_, err := r.PublishBatch(t.Context(), batchOf("e1"))
	requireCode(t, err, codes.Internal)
	requireMessageContains(t, err, "no subscribers registered", "POST /subscribers", "admin address")
}

func TestPublishBatchAllForwardsFailedMessageNamesEachSubscriber(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"b": {id: "b", name: "sub-b", endpoint: "b:1", client: &stubPublishClient{err: errors.New("connection refused")}},
		"a": {id: "a", name: "sub-a", endpoint: "a:1", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{{EventId: "e1"}},
		}}},
	}}

	_, err := r.PublishBatch(t.Context(), batchOf("e1", "e2"))
	requireCode(t, err, codes.Internal)
	requireMessageContains(t, err,
		"no subscriber accepted the batch (2 registered)",
		`subscriber "sub-a" (a:1) returned 1 results for 2 events`,
		`subscriber "sub-b" (b:1) PublishBatch RPC failed: connection refused`,
	)
	// Order must not depend on which forward finished first.
	msg := status.Convert(err).Message()
	if strings.Index(msg, "sub-a") > strings.Index(msg, "sub-b") {
		t.Errorf("subscribers not sorted by name in %q", msg)
	}
}

func TestPublishBatchFailureMessageIsBounded(t *testing.T) {
	long := strings.Repeat("x", 500)
	subs := map[string]*subscriber{}
	for i := 1; i <= 7; i++ {
		id := fmt.Sprintf("s%d", i)
		subs[id] = &subscriber{id: id, name: id, client: &stubPublishClient{err: errors.New(long)}}
	}
	r := &router{subscribers: subs}

	_, err := r.PublishBatch(t.Context(), batchOf("e1"))
	requireCode(t, err, codes.Internal)
	msg := status.Convert(err).Message()
	requireMessageContains(t, err, `"s1"`, `"s5"`, "and 2 more")
	if strings.Contains(msg, `"s6"`) || strings.Contains(msg, `"s7"`) {
		t.Errorf("message lists more than %d subscribers: %q", maxFailureDetails, msg)
	}
	if strings.Contains(msg, strings.Repeat("x", maxErrTextLen+1)) {
		t.Errorf("downstream error text not truncated to %d chars", maxErrTextLen)
	}
	if !strings.Contains(msg, strings.Repeat("x", maxErrTextLen)+"...") {
		t.Errorf("truncated downstream error text should end in an ellipsis: %q", msg)
	}
}

func TestPublishBatchTransactionalRejectionMessageSaysWhichEventsAndWhatToDo(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", client: &stubPublishClient{response: &chippb.PublishResponse{
			Results: []*chippb.PublishResult{
				{EventId: "e1"},
				rejected("e2", "invalid event"),
				rejected("e3", ""),
			},
		}}},
	}}

	batch := batchOf("e1", "e2", "e3")
	batch.Options = &chippb.PublishOptions{TransactionEnabled: ptr(true)}
	_, err := r.PublishBatch(t.Context(), batch)
	requireCode(t, err, codes.Internal)
	requireMessageContains(t, err,
		"transaction_enabled=true and 2 of 3 events were rejected",
		`event[1] id="e2": invalid event`,
		`event[2] id="e3": rejected (`,
		"no reason given",
		"send the batch without transaction_enabled",
	)
}

func TestPublishAllForwardsFailedMessageNamesEachSubscriber(t *testing.T) {
	r := &router{subscribers: map[string]*subscriber{
		"sink": {id: "sink", name: "sink-1", endpoint: "sink:50051", client: &stubPublishClient{err: errors.New("connection refused")}},
	}}

	_, err := r.Publish(t.Context(), eventOf("e1"))
	requireCode(t, err, codes.Internal)
	requireMessageContains(t, err,
		"no subscriber accepted the event (1 registered)",
		`subscriber "sink-1" (sink:50051) Publish RPC failed: connection refused`,
	)
}
