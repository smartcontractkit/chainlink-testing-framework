package fake

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// NewDockerFakeDataProvider creates new fake data provider in Docker using testcontainers-go
func NewDockerFakeDataProvider(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

// NewWithContext creates new fake data provider in Docker using testcontainers-go
func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	bindPort := fmt.Sprintf("%d/tcp", in.Port)
	containerName := framework.DefaultTCName("fake")
	req := tc.ContainerRequest{
		Name:     containerName,
		Image:    in.Image,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: []string{bindPort},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort)
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(fmt.Sprintf("%d/tcp", in.Port))),
		),
	}
	_, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	out := &Output{
		BaseURLHost:   fmt.Sprintf("http://localhost:%d", in.Port),
		BaseURLDocker: fmt.Sprintf("http://%s:%d", containerName, in.Port),
	}
	in.Out = out
	return out, nil
}
