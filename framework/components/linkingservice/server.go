package linkingservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	linkingclient "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
)

const (
	DefaultGRPCPort  = 18124
	DefaultAdminPort = 18125
)

type Config struct {
	GRPCPort  int
	AdminPort int
}

type Server struct {
	linkingclient.UnimplementedLinkingServiceServer

	grpcPort  int
	adminPort int

	grpcServer    *grpc.Server
	grpcListener  net.Listener
	adminServer   *http.Server
	adminListener net.Listener

	mu         sync.RWMutex
	ownerToOrg map[string]string
}

func NewServer(cfg Config) *Server {
	if cfg.GRPCPort == 0 {
		cfg.GRPCPort = DefaultGRPCPort
	}
	if cfg.AdminPort == 0 {
		cfg.AdminPort = DefaultAdminPort
	}

	return &Server{
		grpcPort:   cfg.GRPCPort,
		adminPort:  cfg.AdminPort,
		ownerToOrg: make(map[string]string),
	}
}

func (s *Server) Start(ctx context.Context) error {
	grpcListener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", fmt.Sprintf(":%d", s.grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen on linking gRPC port %d: %w", s.grpcPort, err)
	}

	adminListener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", fmt.Sprintf(":%d", s.adminPort))
	if err != nil {
		_ = grpcListener.Close()
		return fmt.Errorf("failed to listen on linking admin port %d: %w", s.adminPort, err)
	}

	s.grpcListener = grpcListener
	s.adminListener = adminListener

	s.grpcServer = grpc.NewServer()
	linkingclient.RegisterLinkingServiceServer(s.grpcServer, s)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/healthz", s.handleHealth)
	adminMux.HandleFunc("/admin/link", s.handleLink)
	s.adminServer = &http.Server{
		Handler:           adminMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		_ = s.grpcServer.Serve(grpcListener)
	}()
	go func() {
		_ = s.adminServer.Serve(adminListener)
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	if s.grpcListener != nil {
		_ = s.grpcListener.Close()
	}

	if s.adminServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.adminServer.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	if s.adminListener != nil {
		_ = s.adminListener.Close()
	}

	return nil
}

func (s *Server) GetOrganizationFromWorkflowOwner(_ context.Context, req *linkingclient.GetOrganizationFromWorkflowOwnerRequest) (*linkingclient.GetOrganizationFromWorkflowOwnerResponse, error) {
	owner := normalizeWorkflowOwner(req.GetWorkflowOwner())

	s.mu.RLock()
	orgID, ok := s.ownerToOrg[owner]
	s.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "workflow owner %q not linked", req.GetWorkflowOwner())
	}

	return &linkingclient.GetOrganizationFromWorkflowOwnerResponse{OrganizationId: orgID}, nil
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		WorkflowOwner string `json:"workflowOwner"`
		OrgID         string `json:"orgID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request: %v", err), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.WorkflowOwner) == "" || strings.TrimSpace(req.OrgID) == "" {
		http.Error(w, "workflowOwner and orgID are required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.ownerToOrg[normalizeWorkflowOwner(req.WorkflowOwner)] = req.OrgID
	s.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func normalizeWorkflowOwner(owner string) string {
	owner = strings.ToLower(strings.TrimSpace(owner))
	return strings.TrimPrefix(owner, "0x")
}
