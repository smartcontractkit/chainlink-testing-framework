package chiprouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	DefaultGRPCPort         = 50051
	DefaultAdminPort        = 50050
	DefaultBeholderGRPCPort = 50053
	adminPathHealth         = "/health"
)

type Input struct {
	Image         string  `toml:"image" comment:"Chip router Docker image"`
	GRPCPort      int     `toml:"grpc_port" comment:"Chip router gRPC host/container port"`
	AdminPort     int     `toml:"admin_port" comment:"Chip router admin HTTP host/container port"`
	ContainerName string  `toml:"container_name" comment:"Docker container name"`
	PullImage     bool    `toml:"pull_image" comment:"Whether to pull Chip router image or not"`
	LogLevel      string  `toml:"log_level" comment:"Chip router log level (trace, debug, info, warn, error)"`
	Out           *Output `toml:"out" comment:"Chip router output"`
}

type Output struct {
	UseCache         bool   `toml:"use_cache" comment:"Whether to reuse cached output"`
	ContainerName    string `toml:"container_name" comment:"Docker container name"`
	ExternalGRPCURL  string `toml:"grpc_external_url" comment:"Host-reachable gRPC endpoint"`
	InternalGRPCURL  string `toml:"grpc_internal_url" comment:"Docker-network gRPC endpoint"`
	ExternalAdminURL string `toml:"admin_external_url" comment:"Host-reachable admin endpoint"`
	InternalAdminURL string `toml:"admin_internal_url" comment:"Docker-network admin endpoint"`
}

type registerSubscriberRequest struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type registerSubscriberResponse struct {
	ID string `json:"id"`
}

type HealthResponse struct {
}

func defaults(in *Input) {
	if in.GRPCPort == 0 {
		in.GRPCPort = DefaultGRPCPort
	}
	if in.AdminPort == 0 {
		in.AdminPort = DefaultAdminPort
	}
	if in.ContainerName == "" {
		in.ContainerName = framework.DefaultTCName("chip-router")
	}
}

func New(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}

	if strings.TrimSpace(in.Image) == "" {
		return nil, fmt.Errorf("chip router image must be provided")
	}

	defaults(in)

	grpcPort := fmt.Sprintf("%d/tcp", in.GRPCPort)
	adminPort := fmt.Sprintf("%d/tcp", in.AdminPort)

	req := tc.ContainerRequest{
		Name:            in.ContainerName,
		Image:           in.Image,
		AlwaysPullImage: in.PullImage,
		Labels:          framework.DefaultTCLabels(),
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {in.ContainerName},
		},
		ExposedPorts: []string{grpcPort, adminPort},
		Env: map[string]string{
			"CHIP_ROUTER_GRPC_ADDR":  fmt.Sprintf("0.0.0.0:%d", in.GRPCPort),
			"CHIP_ROUTER_ADMIN_ADDR": fmt.Sprintf("0.0.0.0:%d", in.AdminPort),
			"CTF_LOG_LEVEL":          in.LogLevel,
		},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(grpcPort, adminPort)
			h.ExtraHosts = append(h.ExtraHosts, "host.docker.internal:host-gateway")
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(grpcPort).WithPollInterval(200*time.Millisecond),
			tcwait.ForHTTP(adminPathHealth).
				WithPort(adminPort).
				WithStartupTimeout(1*time.Minute).
				WithPollInterval(200*time.Millisecond),
		),
	}

	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := framework.GetHostWithContext(ctx, c)
	if err != nil {
		return nil, err
	}

	out := &Output{
		UseCache:         true,
		ContainerName:    in.ContainerName,
		ExternalGRPCURL:  fmt.Sprintf("%s:%d", host, in.GRPCPort),
		InternalGRPCURL:  fmt.Sprintf("%s:%d", in.ContainerName, in.GRPCPort),
		ExternalAdminURL: fmt.Sprintf("http://%s:%d", host, in.AdminPort),
		InternalAdminURL: fmt.Sprintf("http://%s:%d", in.ContainerName, in.AdminPort),
	}
	in.Out = out
	return out, nil
}

func Health(ctx context.Context, adminURL string) (*HealthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(adminURL, "/")+adminPathHealth, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chip router health request failed with status %s", resp.Status)
	}
	var out HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func RegisterSubscriber(ctx context.Context, adminURL, name, endpoint string) (string, error) {
	body, err := json.Marshal(registerSubscriberRequest{Name: name, Endpoint: endpoint})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(adminURL, "/")+"/subscribers", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chip router register request failed with status %s", resp.Status)
	}
	var out registerSubscriberResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if strings.TrimSpace(out.ID) == "" {
		return "", fmt.Errorf("chip router register response missing subscriber id")
	}
	return out.ID, nil
}

func UnregisterSubscriber(ctx context.Context, adminURL, id string) error {
	if strings.TrimSpace(id) == "" {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, strings.TrimRight(adminURL, "/")+"/subscribers/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("chip router unregister request failed with status %s", resp.Status)
	}
	return nil
}
