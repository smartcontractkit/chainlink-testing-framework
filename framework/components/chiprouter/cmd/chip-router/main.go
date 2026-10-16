package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/google/uuid"
	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	defaultGRPCAddr  = "0.0.0.0:50051"
	defaultAdminAddr = "0.0.0.0:50050"
	forwardTimeout   = 10 * time.Second
	forwardParallel  = 8
)

type publishClient interface {
	Publish(context.Context, *cepb.CloudEvent, ...grpc.CallOption) (*chippb.PublishResponse, error)
	PublishBatch(context.Context, *chippb.CloudEventBatch, ...grpc.CallOption) (*chippb.PublishResponse, error)
}

type subscriber struct {
	id       string
	name     string
	endpoint string
	conn     *grpc.ClientConn
	client   publishClient
}

type router struct {
	chippb.UnimplementedChipIngressServer

	mu          sync.RWMutex
	subscribers map[string]*subscriber
}

type registerSubscriberRequest struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type registerSubscriberResponse struct {
	ID string `json:"id"`
}

type healthResponse struct {
}

func main() {
	if err := run(); err != nil {
		framework.L.Error().Msgf("chip router failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	grpcAddr := envOrDefault("CHIP_ROUTER_GRPC_ADDR", defaultGRPCAddr)
	adminAddr := envOrDefault("CHIP_ROUTER_ADMIN_ADDR", defaultAdminAddr)

	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}
	defer grpcLis.Close()

	adminLis, err := net.Listen("tcp", adminAddr)
	if err != nil {
		return fmt.Errorf("listen admin: %w", err)
	}
	defer adminLis.Close()

	r := &router{subscribers: make(map[string]*subscriber)}

	grpcServer := grpc.NewServer()
	chippb.RegisterChipIngressServer(grpcServer, r)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", r.handleHealth)
	mux.HandleFunc("/subscribers", r.handleSubscribers)
	mux.HandleFunc("/subscribers/", r.handleSubscriberByID)

	adminServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 2)
	go func() { errCh <- grpcServer.Serve(grpcLis) }()
	go func() { errCh <- adminServer.Serve(adminLis) }()

	framework.L.Info().Msgf("chip router started: grpc=%s admin=%s", grpcAddr, adminAddr)

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		done := make(chan struct{})
		go func() {
			defer close(done)
			grpcServer.GracefulStop()
		}()
		select {
		case <-done:
		case <-shutdownCtx.Done():
			grpcServer.Stop()
		}
		_ = adminServer.Shutdown(shutdownCtx)
		r.closeSubscribers()
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		framework.L.Error().Msgf("chip router server error: %v", err)
		return err
	}
}

// Publish matches the write-service's single-event contract (atlas
// chip-ingress internal/server): the RPC is all-or-nothing — it fails unless
// at least one subscriber accepted the event, and returns exactly one
// PublishResult on success.
func (r *router) Publish(ctx context.Context, event *cepb.CloudEvent) (*chippb.PublishResponse, error) {
	snapshot := r.snapshotSubscribers()
	if len(snapshot) == 0 {
		// Acking an event nothing carried would let a durable store delete it undelivered.
		return nil, status.Error(codes.Internal, "no subscribers registered")
	}

	// anyForwarded: at least one subscriber's RPC succeeded. accepted: some
	// subscriber accepted the event (an accepted result or a detail-less
	// success) — sticky, so a later rejection cannot unset it. rejection:
	// recorded only while no subscriber has accepted.
	var mu sync.Mutex
	anyForwarded := false
	accepted := false
	var rejection *chippb.PublishError

	var group errgroup.Group
	group.SetLimit(forwardParallel)
	for _, sub := range snapshot {
		group.Go(func() error {
			framework.L.Debug().Msgf("chip router forwarding event to subscriber id=%s name=%s endpoint=%s", sub.id, sub.name, sub.endpoint)
			forwardCtx, cancel := context.WithTimeout(ctx, forwardTimeout)
			defer cancel()
			resp, err := sub.client.Publish(forwardCtx, event)
			if err != nil {
				framework.L.Error().Msgf("chip router failed to forward event to subscriber id=%s name=%s endpoint=%s err=%v", sub.id, sub.name, sub.endpoint, err)
				return nil
			}
			framework.L.Debug().Msgf("chip router forwarded event to subscriber id=%s", sub.id)
			mu.Lock()
			defer mu.Unlock()
			anyForwarded = true
			if results := resp.GetResults(); len(results) > 0 && results[0].GetError() != nil {
				if !accepted && rejection == nil {
					rejection = results[0].GetError()
				}
				return nil
			}
			accepted = true
			return nil
		})
	}
	_ = group.Wait()

	if !anyForwarded {
		return nil, status.Error(codes.Internal, fmt.Sprintf("all %d subscribers failed to accept the event", len(snapshot)))
	}
	if !accepted {
		return nil, status.Error(codes.Internal, "event was not accepted by any subscriber")
	}
	return &chippb.PublishResponse{Results: []*chippb.PublishResult{{EventId: event.GetId()}}}, nil
}

func (r *router) PublishBatch(ctx context.Context, batch *chippb.CloudEventBatch) (*chippb.PublishResponse, error) {
	// Matches the write-service: an empty batch is a no-op success and never
	// reaches the subscribers.
	if len(batch.GetEvents()) == 0 {
		return &chippb.PublishResponse{Results: []*chippb.PublishResult{}}, nil
	}

	snapshot := r.snapshotSubscribers()
	if len(snapshot) == 0 {
		// Acking events nothing carried would let a durable store delete them undelivered.
		return nil, status.Error(codes.Internal, "no subscribers registered")
	}

	// anyForwarded: at least one subscriber's RPC succeeded. anyBatchAccepted:
	// a subscriber accepted the whole batch without per-event detail. detail:
	// merged per-event results, index-aligned with the batch — an acceptance
	// always wins over a rejection, and a rejection only fills an empty slot,
	// so the outcome does not depend on subscriber completion order.
	var mu sync.Mutex
	anyForwarded := false
	anyBatchAccepted := false
	detail := make([]*chippb.PublishResult, len(batch.GetEvents()))

	var group errgroup.Group
	group.SetLimit(forwardParallel)
	for _, sub := range snapshot {
		group.Go(func() error {
			framework.L.Debug().Msgf("chip router forwarding batch to subscriber id=%s name=%s endpoint=%s", sub.id, sub.name, sub.endpoint)
			forwardCtx, cancel := context.WithTimeout(ctx, forwardTimeout)
			defer cancel()
			resp, err := sub.client.PublishBatch(forwardCtx, batch)
			if err != nil {
				framework.L.Error().Msgf("chip router failed to forward batch to subscriber id=%s name=%s endpoint=%s err=%v", sub.id, sub.name, sub.endpoint, err)
				return nil
			}
			framework.L.Debug().Msgf("chip router forwarded batch to subscriber id=%s", sub.id)
			mu.Lock()
			defer mu.Unlock()
			anyForwarded = true
			results := resp.GetResults()
			if len(results) == 0 {
				// A detail-less success accepts the whole batch.
				anyBatchAccepted = true
				return nil
			}
			for i, res := range results {
				if i >= len(detail) {
					break
				}
				if res.GetError() == nil {
					detail[i] = res
				} else if detail[i] == nil {
					detail[i] = res
				}
			}
			return nil
		})
	}
	_ = group.Wait()

	// No subscriber accepted the batch: fail the RPC so every caller mode
	// retains and retries.
	if !anyForwarded {
		return nil, status.Error(codes.Internal, fmt.Sprintf("all %d subscribers failed to accept the batch", len(snapshot)))
	}

	anyRejected := false
	results := make([]*chippb.PublishResult, 0, len(batch.GetEvents()))
	for i, ev := range batch.GetEvents() {
		result := &chippb.PublishResult{EventId: ev.GetId()}
		if !anyBatchAccepted && detail[i] != nil && detail[i].GetError() != nil {
			result.Error = detail[i].GetError()
			anyRejected = true
		}
		results = append(results, result)
	}

	// Transactional callers resolve the whole batch from the RPC outcome
	// alone, so any unaccepted event must fail the RPC (all-or-nothing).
	if batch.GetOptions().GetTransactionEnabled() && anyRejected {
		return nil, status.Error(codes.Internal, "not all events were accepted by a subscriber")
	}
	return &chippb.PublishResponse{Results: results}, nil
}

func (r *router) handleHealth(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(healthResponse{})
}

func (r *router) handleSubscribers(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body registerSubscriberRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("decode request: %v", err), http.StatusBadRequest)
		return
	}
	body.Endpoint = strings.TrimSpace(body.Endpoint)
	if body.Endpoint == "" {
		http.Error(w, "endpoint is required", http.StatusBadRequest)
		return
	}
	body.Name = strings.TrimSpace(body.Name)

	r.mu.Lock()
	if existing := r.findSubscriberLocked(body.Name, body.Endpoint); existing != nil {
		r.mu.Unlock()
		framework.L.Info().Msgf("chip router reused subscriber id=%s name=%s endpoint=%s", existing.id, existing.name, existing.endpoint)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(registerSubscriberResponse{ID: existing.id})
		return
	}

	framework.L.Info().Msgf("chip router attempting to register subscriber name=%s endpoint=%s", body.Name, body.Endpoint)
	conn, err := grpc.NewClient(body.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		r.mu.Unlock()
		http.Error(w, fmt.Sprintf("dial subscriber: %v", err), http.StatusBadRequest)
		return
	}
	id := uuid.NewString()
	r.subscribers[id] = &subscriber{
		id:       id,
		name:     body.Name,
		endpoint: body.Endpoint,
		conn:     conn,
		client:   chippb.NewChipIngressClient(conn),
	}
	r.mu.Unlock()
	framework.L.Info().Msgf("chip router registered subscriber id=%s name=%s endpoint=%s", id, body.Name, body.Endpoint)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(registerSubscriberResponse{ID: id})
}

func (r *router) findSubscriberLocked(name, endpoint string) *subscriber {
	for _, sub := range r.subscribers {
		if sub.name == name && sub.endpoint == endpoint {
			return sub
		}
	}
	return nil
}

func (r *router) handleSubscriberByID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/subscribers/")
	if id == "" {
		http.Error(w, "subscriber id is required", http.StatusBadRequest)
		return
	}
	r.mu.Lock()
	sub, ok := r.subscribers[id]
	if ok {
		delete(r.subscribers, id)
	}
	r.mu.Unlock()
	if ok && sub.conn != nil {
		_ = sub.conn.Close()
		framework.L.Info().Msgf("chip router unregistered subscriber id=%s name=%s endpoint=%s", id, sub.name, sub.endpoint)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (r *router) snapshotSubscribers() []*subscriber {
	r.mu.RLock()
	defer r.mu.RUnlock()
	snapshot := make([]*subscriber, 0, len(r.subscribers))
	for _, sub := range r.subscribers {
		snapshot = append(snapshot, sub)
	}
	return snapshot
}

func (r *router) closeSubscribers() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, sub := range r.subscribers {
		if sub.conn != nil {
			_ = sub.conn.Close()
		}
		delete(r.subscribers, id)
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
