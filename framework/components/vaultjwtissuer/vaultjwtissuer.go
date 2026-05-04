package vaultjwtissuer

import (
	"context"
	"fmt"
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
	DefaultContainerName = "local-cre-vault-jwt-issuer"
	DefaultHTTPPort      = 18123
	ImageEnvVar          = "CTF_VAULT_JWT_ISSUER_IMAGE"
)

type Input struct {
	Image         string  `toml:"image" comment:"Vault JWT issuer Docker image"`
	HTTPPort      int     `toml:"http_port" comment:"Vault JWT issuer host/container HTTP port"`
	ContainerName string  `toml:"container_name" comment:"Docker container name"`
	PullImage     bool    `toml:"pull_image" comment:"Whether to pull the image or not"`
	Out           *Output `toml:"out" comment:"Vault JWT issuer output"`
}

type Output struct {
	UseCache      bool   `toml:"use_cache" comment:"Whether to reuse cached output"`
	ContainerName string `toml:"container_name" comment:"Docker container name"`
	LocalHTTPURL  string `toml:"local_http_url" comment:"Host-reachable HTTP endpoint"`
	DockerHTTPURL string `toml:"docker_http_url" comment:"Docker-network HTTP endpoint"`
}

func defaults(in *Input) {
	if in.HTTPPort == 0 {
		in.HTTPPort = DefaultHTTPPort
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
		return nil, fmt.Errorf("vault JWT issuer input is nil")
	}
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}

	defaults(in)
	if strings.TrimSpace(in.Image) == "" {
		return nil, fmt.Errorf("vault JWT issuer image must be provided")
	}
	httpPort := fmt.Sprintf("%d/tcp", in.HTTPPort)

	req := tc.ContainerRequest{
		Name:            in.ContainerName,
		Image:           in.Image,
		AlwaysPullImage: in.PullImage,
		Labels:          framework.DefaultTCLabels(),
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {in.ContainerName},
		},
		ExposedPorts: []string{httpPort},
		Env: map[string]string{
			"VAULT_JWT_ISSUER_HTTP_PORT": fmt.Sprintf("%d", in.HTTPPort),
		},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(httpPort)
			h.ExtraHosts = append(h.ExtraHosts, "host.docker.internal:host-gateway")
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(httpPort)).WithPollInterval(200*time.Millisecond),
			tcwait.ForHTTP("/admin/healthz").
				WithPort(nat.Port(httpPort)).
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
		UseCache:      true,
		ContainerName: in.ContainerName,
		LocalHTTPURL:  fmt.Sprintf("http://%s:%d", host, in.HTTPPort),
		DockerHTTPURL: fmt.Sprintf("http://%s:%d", in.ContainerName, in.HTTPPort),
	}
	in.Out = out
	return out, nil
}
