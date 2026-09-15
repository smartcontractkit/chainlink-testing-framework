package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"net/netip"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
)

const (
	// DefaultStellarImage is the official Stellar quickstart image for local development.
	// Pinned to a specific multi-arch (amd64+arm64) per-commit tag instead of :latest so test
	// runs are reproducible and a quickstart release cannot silently change flags/ports.
	// https://github.com/stellar/quickstart
	DefaultStellarImage = "stellar/quickstart:v667-b1428.1-latest"

	// DefaultStellarRPCPort is the port the quickstart unified HTTP gateway listens on.
	// The gateway multiplexes by path: Horizon at "/", Soroban RPC at "/rpc", Friendbot at "/friendbot".
	DefaultStellarRPCPort = "8000"

	// DefaultStellarNetworkPassphrase is the network passphrase for local standalone network
	// https://stellar.org/developers/guides/concepts/networks
	DefaultStellarNetworkPassphrase = "Standalone Network ; February 2017"
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

	// The quickstart image publishes multi-arch (amd64+arm64) manifests, so select the
	// native platform to avoid amd64 emulation on arm64 hosts (e.g. Apple Silicon).
	imagePlatform := "linux/amd64"
	if runtime.GOARCH == "arm64" {
		imagePlatform = "linux/arm64"
	}
	if in.ImagePlatform != nil {
		imagePlatform = *in.ImagePlatform
	}

	// Build the command arguments. In --local mode the quickstart image runs all services by
	// default (core, horizon, Soroban RPC, Friendbot, Lab), so no service-enable flag is needed.
	// The older "--enable-soroban-rpc" flag is no longer valid; the current form is "--enable"
	// with a comma-separated service list, which is only used to run a subset.
	// https://github.com/stellar/quickstart#usage
	cmd := []string{
		"--local",
	}

	// Allow additional command overrides
	if len(in.DockerCmdParamsOverrides) > 0 {
		cmd = append(cmd, in.DockerCmdParamsOverrides...)
	}

	if pods.K8sEnabled() {
		return nil, fmt.Errorf("K8s support is not yet implemented")
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
			h.PortBindings = network.PortMap{
				network.MustParsePort(containerPort): []network.PortBinding{
					{
						HostIP:   netip.MustParseAddr("0.0.0.0"),
						HostPort: in.Port,
					},
				},
			}
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		ImagePlatform: imagePlatform,
		Cmd:           cmd,
		// Cheap TCP-listening gate on the gateway port. The real readiness check is the
		// app-level getHealth poll in waitForStellarRPC below (RPC is only healthy after
		// core+horizon bootstrap), which is why we don't gate on the Horizon "/" root here.
		WaitingFor: wait.ForListeningPort(containerPort).
			WithStartupTimeout(1 * time.Minute).
			WithPollInterval(500 * time.Millisecond),
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
			if checkStellarHealth(ctx, rpcURL) {
				return nil
			}
			framework.L.Debug().Str("url", rpcURL).Msg("Waiting for Stellar RPC to be ready...")
		}
	}
}

// checkStellarHealth checks if Stellar RPC reports a healthy getHealth result.
// Readiness, per the stellar-rpc spec, is getHealth result.status == "healthy" (the RPC only
// becomes healthy after core+horizon have bootstrapped), so we parse the JSON-RPC envelope
// rather than substring-matching the body. The request is bound to ctx so an in-flight probe
// is aborted when the caller's context (or the waitForStellarRPC deadline) is cancelled.
func checkStellarHealth(ctx context.Context, rpcURL string) bool {
	client := &http.Client{Timeout: 5 * time.Second}

	reqBody := `{"jsonrpc":"2.0","id":1,"method":"getHealth"}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcURL, strings.NewReader(reqBody))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var rpcResp struct {
		Result struct {
			Status string `json:"status"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return false
	}
	return rpcResp.Error == nil && rpcResp.Result.Status == "healthy"
}
