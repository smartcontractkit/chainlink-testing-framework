package chiprouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	DefaultGRPCPort         = 50051
	DefaultAdminPort        = 50050
	DefaultBeholderGRPCPort = 50053
	adminPathHealth         = "/health"

	ImageOverrideEnvVar = "CTF_CHIP_ROUTER_IMAGE"
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
	if in.ContainerName == "" {
		in.ContainerName = framework.DefaultTCName("chip-router")
	}
	if strings.TrimSpace(os.Getenv(ImageOverrideEnvVar)) != "" {
		in.Image = os.Getenv(ImageOverrideEnvVar)
	}
}

func New(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}

	defaults(in)

	if strings.TrimSpace(in.Image) == "" {
		return nil, fmt.Errorf("chip router image must be provided")
	}

	var internalGRPCNatPort string
	var internalGRPCPort int
	if in.GRPCPort == 0 {
		internalGRPCNatPort = fmt.Sprintf("%d/tcp", DefaultGRPCPort)
		internalGRPCPort = DefaultGRPCPort
	} else {
		internalGRPCNatPort = fmt.Sprintf("%d/tcp", in.GRPCPort)
		internalGRPCPort = in.GRPCPort
	}

	var internalAdminNatPort string
	var internalAdminPort int
	if in.AdminPort == 0 {
		internalAdminNatPort = fmt.Sprintf("%d/tcp", DefaultAdminPort)
		internalAdminPort = DefaultAdminPort
	} else {
		internalAdminNatPort = fmt.Sprintf("%d/tcp", in.AdminPort)
		internalAdminPort = in.AdminPort
	}

	req := tc.ContainerRequest{
		Name:            in.ContainerName,
		Image:           in.Image,
		AlwaysPullImage: in.PullImage,
		Labels:          framework.DefaultTCLabels(),
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {in.ContainerName},
		},
		ExposedPorts: []string{internalGRPCNatPort, internalAdminNatPort},
		Env: map[string]string{
			"CHIP_ROUTER_GRPC_ADDR":  fmt.Sprintf("0.0.0.0:%d", internalGRPCPort),
			"CHIP_ROUTER_ADMIN_ADDR": fmt.Sprintf("0.0.0.0:%d", internalAdminPort),
			"CTF_LOG_LEVEL":          in.LogLevel,
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(internalGRPCNatPort)).WithPollInterval(200*time.Millisecond),
			tcwait.ForHTTP(adminPathHealth).
				WithPort(nat.Port(internalAdminNatPort)).
				WithStartupTimeout(1*time.Minute).
				WithPollInterval(200*time.Millisecond),
		),
	}

	staticPortBindings := []string{}
	if in.GRPCPort != 0 {
		staticPortBindings = append(staticPortBindings, fmt.Sprintf("%d/tcp", in.GRPCPort))
	}
	if in.AdminPort != 0 {
		staticPortBindings = append(staticPortBindings, fmt.Sprintf("%d/tcp", in.AdminPort))
	}

	if len(staticPortBindings) > 0 {
		req.HostConfigModifier = func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(staticPortBindings...)
			h.ExtraHosts = append(h.ExtraHosts, "host.docker.internal:host-gateway")
		}
	} else {
		req.HostConfigModifier = func(h *container.HostConfig) {
			h.ExtraHosts = append(h.ExtraHosts, "host.docker.internal:host-gateway")
		}
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
		InternalAdminURL: fmt.Sprintf("http://%s:%d", in.ContainerName, internalAdminPort),
		InternalGRPCURL:  fmt.Sprintf("%s:%d", in.ContainerName, internalGRPCPort),
	}

	if in.GRPCPort != 0 {
		out.ExternalGRPCURL = fmt.Sprintf("%s:%d", host, in.GRPCPort)
	} else {
		if p, err := c.MappedPort(ctx, nat.Port(internalGRPCNatPort)); err != nil {
			return nil, err
		} else {
			out.ExternalGRPCURL = fmt.Sprintf("%s:%d", host, p.Int())
		}
	}

	if in.AdminPort != 0 {
		out.ExternalAdminURL = fmt.Sprintf("http://%s:%d", host, in.AdminPort)
	} else {
		if p, err := c.MappedPort(ctx, nat.Port(internalAdminNatPort)); err != nil {
			return nil, err
		} else {
			out.ExternalAdminURL = fmt.Sprintf("http://%s:%d", host, p.Int())
		}
	}

	in.Out = out

	fmt.Println("[ctf] chip router internal grpc url", out.InternalGRPCURL)
	fmt.Println("[ctf] chip router internal grpc url", in.Out.InternalGRPCURL)

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
