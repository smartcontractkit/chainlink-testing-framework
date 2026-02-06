package blockchain

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	// DefaultStellarImage is the official Stellar quickstart image for local development
	// https://github.com/stellar/quickstart
	DefaultStellarImage = "stellar/quickstart:latest"

	// DefaultStellarRPCPort is the port Stellar RPC listens on
	DefaultStellarRPCPort = "8000"

	// DefaultStellarNetworkPassphrase is the network passphrase for local standalone network
	// https://stellar.org/developers/guides/concepts/networks
	DefaultStellarNetworkPassphrase = "Standalone Network ; February 2017"

	// DefaultStellarFriendbotPort is the port for the Friendbot faucet service
	DefaultStellarFriendbotPort = "8000"
)

// StellarNetworkInfo contains Stellar network-specific configuration
type StellarNetworkInfo struct {
	NetworkPassphrase string `toml:"network_passphrase" json:"networkPassphrase" comment:"Stellar network passphrase"`
	FriendbotURL      string `toml:"friendbot_url" json:"friendbotUrl" comment:"Friendbot faucet URL for funding accounts"`
}

func defaultStellar(in *Input) {
	if in.Image == "" {
		in.Image = DefaultStellarImage
	}
	if in.Port == "" {
		in.Port = DefaultStellarRPCPort
	}
}

func newStellar(ctx context.Context, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	defaultStellar(in)

	containerName := framework.DefaultTCName("stellar-node")

	// Stellar RPC container always listens on port 8000 internally
	containerPort := fmt.Sprintf("%s/tcp", DefaultStellarRPCPort)

	// default to amd64
	imagePlatform := "linux/amd64"
	if in.ImagePlatform != nil {
		imagePlatform = *in.ImagePlatform
	}

	// Build the command arguments
	cmd := []string{
		"--local",
		"--enable-soroban-rpc",
	}

	// Allow additional command overrides
	if len(in.DockerCmdParamsOverrides) > 0 {
		cmd = append(cmd, in.DockerCmdParamsOverrides...)
	}

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           in.Image,
		ExposedPorts:    []string{containerPort},
		Name:            containerName,
		Labels:          framework.DefaultTCLabels(),
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		HostConfigModifier: func(h *container.HostConfig) {
			// Map user-provided host port to container's default port (8000)
			h.PortBindings = nat.PortMap{
				nat.Port(containerPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: in.Port,
					},
				},
			}
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		ImagePlatform: imagePlatform,
		Cmd:           cmd,
		// Wait for passing health check
		WaitingFor: wait.ForHTTP("/").
			WithPort(nat.Port(containerPort)).
			WithStatusCodeMatcher(func(status int) bool {
				return status >= 200 && status < 500
			}).
			WithStartupTimeout(3 * time.Minute).
			WithPollInterval(2 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Stellar container: %w", err)
	}

	host, err := framework.GetHostWithContext(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	// Verify the RPC is actually responding
	if err := waitForStellarRPC(ctx, host, in.Port); err != nil {
		return nil, fmt.Errorf("stellar RPC failed to become ready: %w", err)
	}

	framework.L.Info().
		Str("host", host).
		Str("port", in.Port).
		Str("network_passphrase", DefaultStellarNetworkPassphrase).
		Msg("Stellar node is ready")

	return &Output{
		ChainID:       in.ChainID,
		UseCache:      true,
		Type:          in.Type,
		Family:        FamilyStellar,
		ContainerName: containerName,
		Container:     c,
		NetworkSpecificData: &NetworkSpecificData{
			StellarNetwork: &StellarNetworkInfo{
				NetworkPassphrase: DefaultStellarNetworkPassphrase,
				FriendbotURL:      fmt.Sprintf("http://%s:%s/friendbot", host, in.Port),
			},
		},
		Nodes: []*Node{
			{
				// RPC endpoint for JSON-RPC calls
				ExternalHTTPUrl: fmt.Sprintf("http://%s:%s/rpc", host, in.Port),
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s/rpc", containerName, DefaultStellarRPCPort),
			},
		},
	}, nil
}

// waitForStellarRPC polls the Stellar RPC endpoint until it responds to getHealth
func waitForStellarRPC(ctx context.Context, host, port string) error {
	rpcURL := fmt.Sprintf("http://%s:%s/rpc", host, port)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(3 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for Stellar RPC at %s", rpcURL)
		case <-ticker.C:
			if checkStellarHealth(rpcURL) {
				return nil
			}
			framework.L.Debug().Str("url", rpcURL).Msg("Waiting for Stellar RPC to be ready...")
		}
	}
}

// checkStellarHealth checks if Stellar RPC responds to getHealth method
func checkStellarHealth(rpcURL string) bool {
	client := &http.Client{Timeout: 5 * time.Second}

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"getHealth"}`
	resp, err := client.Post(rpcURL, "application/json", strings.NewReader(reqBody))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Read response body to check for valid JSON-RPC response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// Check if we got a valid JSON-RPC response (not an error)
	return resp.StatusCode == 200 && len(body) > 0 && strings.Contains(string(body), "result")
}
