package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
)

const (
	adminRequestTimeout = 5 * time.Second
)

type RegisterSubscriberRequest struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type RegisterSubscriberResponse struct {
	ID string `json:"id"`
}

type RealthResponse struct {
	AdminURL string `json:"admin_url"`
	GRPCURL  string `json:"grpc_url"`
}

type Client struct {
	httpClient *http.Client
	adminURL   string
	grpcURL    string
}

func New(ctx context.Context, adminURL, grpcURL string) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: adminRequestTimeout},
		adminURL:   adminURL,
		grpcURL:    grpcURL,
	}

	if !isHTTPReady(ctx, c.adminURL) {
		return nil, fmt.Errorf("chip ingress router admin endpoint is not reachable: %s", c.adminURL)
	}
	if !isTCPReady(c.grpcURL) {
		return nil, fmt.Errorf("chip ingress router grpc endpoint is not reachable: %s", c.grpcURL)
	}
	return c, nil
}

func (c *Client) RegisterSubscriber(ctx context.Context, name, endpoint string) (string, error) {
	body, err := json.Marshal(RegisterSubscriberRequest{Name: name, Endpoint: endpoint})
	if err != nil {
		return "", pkgerrors.Wrap(err, "marshal chip router register request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.adminURL, "/")+"/subscribers", bytes.NewReader(body))
	if err != nil {
		return "", pkgerrors.Wrap(err, "create chip router register request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", pkgerrors.Wrap(err, "perform chip router register request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chip router register request failed with status %s", resp.Status)
	}

	var out RegisterSubscriberResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", pkgerrors.Wrap(err, "decode chip router register response")
	}
	if out.ID == "" {
		return "", pkgerrors.New("chip router register response missing subscriber id")
	}

	return out.ID, nil
}

func (c *Client) UnregisterSubscriber(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, strings.TrimRight(c.adminURL, "/")+"/subscribers/"+id, nil)
	if err != nil {
		return pkgerrors.Wrap(err, "create chip router unregister request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return pkgerrors.Wrap(err, "perform chip router unregister request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("chip router unregister request failed with status %s", resp.Status)
	}
	return nil
}

func isHTTPReady(ctx context.Context, adminURL string) bool {
	_, err := fetchHealth(ctx, adminURL)
	return err == nil
}

func fetchHealth(ctx context.Context, adminURL string) (*RealthResponse, error) {
	if strings.TrimSpace(adminURL) == "" {
		return nil, pkgerrors.New("admin url is empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(adminURL, "/")+"/health", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected health status %s", resp.Status)
	}
	var health RealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}
	return &health, nil
}

func isTCPReady(addr string) bool {
	dialer := &net.Dialer{Timeout: time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
