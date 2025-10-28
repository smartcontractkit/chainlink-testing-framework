package jd

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// GRPCHealthStrategy implements a wait strategy for gRPC health checks
type GRPCHealthStrategy struct {
	Port         nat.Port
	PollInterval time.Duration
	timeout      time.Duration
}

// NewGRPCHealthStrategy creates a new gRPC health check wait strategy
func NewGRPCHealthStrategy(port nat.Port) *GRPCHealthStrategy {
	return &GRPCHealthStrategy{
		Port:         port,
		PollInterval: 200 * time.Millisecond,
		timeout:      3 * time.Minute,
	}
}

// WithTimeout sets the timeout for the gRPC health check strategy
func (g *GRPCHealthStrategy) WithTimeout(timeout time.Duration) *GRPCHealthStrategy {
	g.timeout = timeout
	return g
}

// WithPollInterval sets the poll interval for the gRPC health check strategy
func (g *GRPCHealthStrategy) WithPollInterval(interval time.Duration) *GRPCHealthStrategy {
	g.PollInterval = interval
	return g
}

// WaitUntilReady implements Strategy.WaitUntilReady
func (g *GRPCHealthStrategy) WaitUntilReady(ctx context.Context, target tcwait.StrategyTarget) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(g.PollInterval):
			// Check if container is still running
			state, err := target.State(ctx)
			if err != nil {
				return err
			}
			if !state.Running {
				return fmt.Errorf("container is not running: %s", state.Status)
			}

			// Get host and port
			host, err := framework.GetHostWithContext(ctx, target.(tc.Container))
			if err != nil {
				continue
			}

			mappedPort, err := target.MappedPort(ctx, g.Port)
			if err != nil {
				continue
			}

			// Attempt gRPC health check
			address := fmt.Sprintf("%s:%s", host, mappedPort.Port())
			if err := g.checkHealth(ctx, address); err == nil {
				return nil
			}
		}
	}
}

// checkHealth performs the actual gRPC health check
func (g *GRPCHealthStrategy) checkHealth(ctx context.Context, address string) error {
	// Create a short timeout for the individual check
	checkCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Use plaintext/insecure connection (standard for local testing and health checks)
	return g.tryHealthCheck(checkCtx, address, insecure.NewCredentials())
}

// tryHealthCheck attempts a health check with specific credentials
func (g *GRPCHealthStrategy) tryHealthCheck(ctx context.Context, address string, creds credentials.TransportCredentials) error {
	// Build dial options similar to the working JD connection code
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Create the gRPC client connection
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	// Create health check client
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Perform health check
	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return err
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("service not serving, status: %v", resp.Status)
	}

	return nil
}
