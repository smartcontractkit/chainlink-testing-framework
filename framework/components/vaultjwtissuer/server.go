package vaultjwtissuer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"crypto/rsa"
)

type Config struct {
	HTTPPort int
}

type Server struct {
	httpPort     int
	privateKey   *rsa.PrivateKey
	httpServer   *http.Server
	httpListener net.Listener
}

func NewServer(cfg Config) (*Server, error) {
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = DefaultHTTPPort
	}

	privateKey, err := parseDefaultJWTSigningKey()
	if err != nil {
		return nil, err
	}

	return &Server{
		httpPort:   cfg.HTTPPort,
		privateKey: privateKey,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	httpListener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", fmt.Sprintf(":%d", s.httpPort))
	if err != nil {
		return fmt.Errorf("failed to listen on HTTP port %d: %w", s.httpPort, err)
	}

	s.httpListener = httpListener

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", s.handleJWKS)
	mux.HandleFunc("/.well-known/openid-configuration", s.handleOpenIDConfiguration)
	mux.HandleFunc("/admin/healthz", s.handleHealth)
	s.httpServer = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		_ = s.httpServer.Serve(httpListener)
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	if s.httpListener != nil {
		_ = s.httpListener.Close()
	}

	return nil
}

func (s *Server) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"keys": []jwtWebKey{rsaPublicKeyToJWK(DefaultJWTIssuerKeyID, &s.privateKey.PublicKey)},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	issuerURL := requestBaseURL(r)
	resp := map[string]string{
		"issuer":   issuerURL,
		"jwks_uri": strings.TrimSuffix(issuerURL, "/") + "/.well-known/jwks.json",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
