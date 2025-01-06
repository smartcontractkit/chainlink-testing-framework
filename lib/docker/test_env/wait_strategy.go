package test_env

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

type HTTPStrategy struct {
	Path               string
	Port               nat.Port
	RetryDelay         time.Duration
	ExpectedStatusCode int
	timeout            time.Duration
}

// NewHTTPStrategy initializes a new HTTP strategy for waiting on a service to become available.
// It sets the path, port, retry delay, expected status code, and timeout, allowing for flexible service readiness checks.
func NewHTTPStrategy(path string, port nat.Port) *HTTPStrategy {
	return &HTTPStrategy{
		Path:               path,
		Port:               port,
		RetryDelay:         10 * time.Second,
		ExpectedStatusCode: 200,
		timeout:            2 * time.Minute,
	}
}

// WithTimeout sets the timeout duration for HTTP requests.
// It returns the updated HTTPStrategy instance, allowing for method chaining.
func (w *HTTPStrategy) WithTimeout(timeout time.Duration) *HTTPStrategy {
	w.timeout = timeout
	return w
}

// WithStatusCode sets the expected HTTP status code for the HTTP strategy.
// This allows users to specify the desired response code to validate during service startup.
func (w *HTTPStrategy) WithStatusCode(statusCode int) *HTTPStrategy {
	w.ExpectedStatusCode = statusCode
	return w
}

// WaitUntilReady implements Strategy.WaitUntilReady
func (w *HTTPStrategy) WaitUntilReady(ctx context.Context, target tcwait.StrategyTarget) (err error) {

	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	host, err := GetHost(ctx, target.(tc.Container))
	if err != nil {
		return
	}

	var mappedPort nat.Port
	mappedPort, err = target.MappedPort(ctx, w.Port)
	if err != nil {
		return err
	}

	tripper := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{Transport: tripper, Timeout: time.Second}
	address := net.JoinHostPort(host, strconv.Itoa(mappedPort.Int()))

	endpoint := url.URL{
		Scheme: "http",
		Host:   address,
		Path:   w.Path,
	}

	var body []byte
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			state, err := target.State(ctx)
			if err != nil {
				return err
			}
			if !state.Running {
				return fmt.Errorf("container is not running %s", state.Status)
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), bytes.NewReader(body))
			if err != nil {
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			if resp.StatusCode != w.ExpectedStatusCode {
				_ = resp.Body.Close()
				continue
			}
			if err := resp.Body.Close(); err != nil {
				continue
			}
			return nil
		}
	}
}

type WebSocketStrategy struct {
	Port       nat.Port
	RetryDelay time.Duration
	timeout    time.Duration
	l          zerolog.Logger
}

// NewWebSocketStrategy initializes a WebSocket strategy for monitoring service readiness.
// It sets the port and defines retry behavior, making it useful for ensuring services are operational before proceeding.
func NewWebSocketStrategy(port nat.Port, l zerolog.Logger) *WebSocketStrategy {
	return &WebSocketStrategy{
		Port:       port,
		RetryDelay: 10 * time.Second,
		timeout:    2 * time.Minute,
	}
}

// WithTimeout sets the timeout duration for the WebSocket strategy.
// It allows users to specify how long to wait for a response before timing out,
// enhancing control over connection behavior in network operations.
func (w *WebSocketStrategy) WithTimeout(timeout time.Duration) *WebSocketStrategy {
	w.timeout = timeout
	return w
}

// WaitUntilReady waits for the WebSocket service to become available by repeatedly attempting to connect.
// It returns an error if the connection cannot be established within the specified timeout.
func (w *WebSocketStrategy) WaitUntilReady(ctx context.Context, target tcwait.StrategyTarget) (err error) {
	var client *rpc.Client
	var host string
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()
	i := 0
	for {
		host, err = GetHost(ctx, target.(tc.Container))
		if err != nil {
			w.l.Error().Msg("Failed to get the target host")
			return err
		}
		wsPort, err := target.MappedPort(ctx, w.Port)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
		w.l.Info().Msgf("Attempting to dial %s", url)
		client, err = rpc.DialContext(ctx, url)
		if err == nil {
			client.Close()
			w.l.Info().Msg("WebSocket rpc port is ready")
			return nil
		}
		if client != nil {
			client.Close() // Close client if DialContext failed
			client = nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(w.RetryDelay):
			i++
			w.l.Info().Msgf("WebSocket attempt %d failed: %s. Retrying...", i, err)
		}
	}
}
