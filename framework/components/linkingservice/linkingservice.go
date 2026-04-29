package linkingservice

import (
	"context"
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
	DefaultContainerName = "local-cre-linking-service"
	ImageEnvVar          = "CTF_LINKING_SERVICE_IMAGE"
)

type Input struct {
	Image         string  `toml:"image" comment:"Linking service Docker image"`
	GRPCPort      int     `toml:"grpc_port" comment:"Linking service gRPC host/container port"`
	AdminPort     int     `toml:"admin_port" comment:"Linking service admin HTTP host/container port"`
	ContainerName string  `toml:"container_name" comment:"Docker container name"`
	PullImage     bool    `toml:"pull_image" comment:"Whether to pull the image or not"`
	Out           *Output `toml:"out" comment:"Linking service output"`
}

type Output struct {
	UseCache       bool   `toml:"use_cache" comment:"Whether to reuse cached output"`
	ContainerName  string `toml:"container_name" comment:"Docker container name"`
	LocalGRPCURL   string `toml:"local_grpc_url" comment:"Host-reachable gRPC endpoint"`
	DockerGRPCURL  string `toml:"docker_grpc_url" comment:"Docker-network gRPC endpoint"`
	LocalAdminURL  string `toml:"local_admin_url" comment:"Host-reachable admin endpoint"`
	DockerAdminURL string `toml:"docker_admin_url" comment:"Docker-network admin endpoint"`
}

func defaults(in *Input) {
	if in.GRPCPort == 0 {
		in.GRPCPort = DefaultGRPCPort
	}
	if in.AdminPort == 0 {
		in.AdminPort = DefaultAdminPort
	}
	ApplyImageOverride(in)
	if in.ContainerName == "" {
		in.ContainerName = DefaultContainerName
	}
}

func ApplyImageOverride(in *Input) string {
	if in == nil {
		return ""
	}

	override := strings.TrimSpace(os.Getenv(ImageEnvVar))
	if override == "" {
		return ""
	}

	in.Image = override
	return override
}

func New(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in == nil {
		return nil, fmt.Errorf("linking service input is nil")
	}
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}

	defaults(in)
	if strings.TrimSpace(in.Image) == "" {
		return nil, fmt.Errorf("linking service image must be provided")
	}

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
			"LINKING_SERVICE_GRPC_PORT":  fmt.Sprintf("%d", in.GRPCPort),
			"LINKING_SERVICE_ADMIN_PORT": fmt.Sprintf("%d", in.AdminPort),
		},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(grpcPort, adminPort)
			h.ExtraHosts = append(h.ExtraHosts, "host.docker.internal:host-gateway")
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(grpcPort)).WithPollInterval(200*time.Millisecond),
			tcwait.ForHTTP("/admin/healthz").
				WithPort(nat.Port(adminPort)).
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
		UseCache:       true,
		ContainerName:  in.ContainerName,
		LocalGRPCURL:   fmt.Sprintf("%s:%d", host, in.GRPCPort),
		DockerGRPCURL:  fmt.Sprintf("%s:%d", in.ContainerName, in.GRPCPort),
		LocalAdminURL:  fmt.Sprintf("http://%s:%d", host, in.AdminPort),
		DockerAdminURL: fmt.Sprintf("http://%s:%d", in.ContainerName, in.AdminPort),
	}
	in.Out = out
	return out, nil
}

func Health(ctx context.Context, adminURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(adminURL, "/")+"/admin/healthz", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("linking service health request failed with status %s", resp.Status)
	}
	return nil
}
