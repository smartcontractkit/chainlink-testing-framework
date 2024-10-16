package blockchain

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

func deployAnvil(in *Input) (*Output, error) {
	entryPoint := []string{"anvil", "--host", "0.0.0.0", "--port", in.Port, "--chain-id", in.ChainID}
	entryPoint = append(entryPoint, in.DockerCmdParamsOverrides...)
	bindPort := fmt.Sprintf("%s/tcp", in.Port)

	req := testcontainers.ContainerRequest{
		Image:      "ghcr.io/foundry-rs/foundry",
		Labels:     framework.DefaultTCLabels(),
		Name:       framework.DefaultTCName("anvil"),
		WaitingFor: wait.ForListeningPort(nat.Port(in.Port)).WithStartupTimeout(10 * time.Second),
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.NetworkMode = "host"
			hc.PortBindings = framework.MapTheSamePort(bindPort)
		},
		Entrypoint: entryPoint,
	}
	_, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	return &Output{
		ChainID: in.ChainID,
		Nodes: []*Node{
			{
				WSUrl:   fmt.Sprintf("ws://localhost:%s", in.Port),
				HTTPUrl: fmt.Sprintf("http://localhost:%s", in.Port),
			},
		},
	}, nil
}
