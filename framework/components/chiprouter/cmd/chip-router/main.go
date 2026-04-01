package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	adminServer := &http.Server{Handler: mux}

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
		if err == nil || err == http.ErrServerClosed {
			return nil
		}
		framework.L.Error().Msgf("chip router server error: %v", err)
		return err
	}
}

func (r *router) Publish(_ context.Context, event *cepb.CloudEvent) (*chippb.PublishResponse, error) {
	snapshot := r.snapshotSubscribers()
	if len(snapshot) == 0 {
		return &chippb.PublishResponse{}, nil
	}

	var group errgroup.Group
	group.SetLimit(forwardParallel)
	for _, sub := range snapshot {
		group.Go(func() error {
			framework.L.Debug().Msgf("chip router forwarding event to subscriber id=%s name=%s endpoint=%s", sub.id, sub.name, sub.endpoint)
			forwardCtx, cancel := context.WithTimeout(context.Background(), forwardTimeout)
			defer cancel()
			_, err := sub.client.Publish(forwardCtx, event)
			if err != nil {
				framework.L.Error().Msgf("chip router failed to forward event to subscriber id=%s name=%s endpoint=%s err=%v", sub.id, sub.name, sub.endpoint, err)
			}
			framework.L.Debug().Msgf("chip router forwarded event to subscriber id=%s", sub.id)
			return nil
		})
	}
	_ = group.Wait()
	return &chippb.PublishResponse{}, nil
}

func (r *router) PublishBatch(_ context.Context, batch *chippb.CloudEventBatch) (*chippb.PublishResponse, error) {
	snapshot := r.snapshotSubscribers()
	if len(snapshot) == 0 {
		return &chippb.PublishResponse{}, nil
	}

	var group errgroup.Group
	group.SetLimit(forwardParallel)
	for _, sub := range snapshot {
		group.Go(func() error {
			framework.L.Debug().Msgf("chip router forwarding batch to subscriber id=%s name=%s endpoint=%s", sub.id, sub.name, sub.endpoint)
			forwardCtx, cancel := context.WithTimeout(context.Background(), forwardTimeout)
			defer cancel()
			_, err := sub.client.PublishBatch(forwardCtx, batch)
			if err != nil {
				log.Printf("chip router failed to forward batch to subscriber id=%s name=%s endpoint=%s err=%v", sub.id, sub.name, sub.endpoint, err)
			}
			framework.L.Debug().Msgf("chip router forwarded batch to subscriber id=%s", sub.id)
			return nil
		})
	}
	_ = group.Wait()
	return &chippb.PublishResponse{}, nil
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
	framework.L.Info().Msgf("chip router attempting to register subscriber name=%s endpoint=%s", body.Name, body.Endpoint)
	conn, err := grpc.NewClient(body.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, fmt.Sprintf("dial subscriber: %v", err), http.StatusBadRequest)
		return
	}
	id := uuid.NewString()
	r.mu.Lock()
	r.subscribers[id] = &subscriber{
		id:       id,
		name:     strings.TrimSpace(body.Name),
		endpoint: body.Endpoint,
		conn:     conn,
		client:   chippb.NewChipIngressClient(conn),
	}
	r.mu.Unlock()
	framework.L.Info().Msgf("chip router registered subscriber id=%s name=%s endpoint=%s", id, body.Name, body.Endpoint)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(registerSubscriberResponse{ID: id})
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
	framework.L.Info().Msgf("chip router unregistered subscriber id=%s name=%s endpoint=%s", id, sub.name, sub.endpoint)
	if ok && sub.conn != nil {
		_ = sub.conn.Close()
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
