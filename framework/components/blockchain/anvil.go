package blockchain

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

func deployAnvil(in *Input) (*Output, error) {
	ctx := context.Background()
	entryPoint := []string{"anvil", "--host", "0.0.0.0", "--port", in.Port, "--chain-id", in.ChainID}
	entryPoint = append(entryPoint, in.DockerCmdParamsOverrides...)
	bindPort := fmt.Sprintf("%s/tcp", in.Port)

	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/foundry-rs/foundry",
		Labels:       framework.DefaultTCLabels(),
		Name:         framework.DefaultTCName("anvil"),
		ExposedPorts: []string{bindPort},
		NetworkAliases: map[string][]string{
			"bridge": {"anvil"},
		},
		WaitingFor: wait.ForListeningPort(nat.Port(in.Port)).WithStartupTimeout(10 * time.Second),
		//HostConfigModifier: func(hc *container.HostConfig) {
		//	hc.NetworkMode = "host"
		//	hc.PortBindings = framework.MapTheSamePort(bindPort)
		//},
		Entrypoint: entryPoint,
	}
	c, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Output{
		ChainID: in.ChainID,
		Nodes: []*Node{
			{
				WSUrl:   fmt.Sprintf("ws://%s:%s", host, mp.Port()),
				HTTPUrl: fmt.Sprintf("http://%s:%s", host, mp.Port()),
			},
		},
	}, nil
}
