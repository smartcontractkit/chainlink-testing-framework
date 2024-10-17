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
	containerName := framework.DefaultTCName("anvil")

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", in.Image, in.Tag),
		Labels:       framework.DefaultTCLabels(),
		Name:         containerName,
		ExposedPorts: []string{bindPort},
		Networks:     []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		WaitingFor: wait.ForListeningPort(nat.Port(in.Port)).WithStartupTimeout(10 * time.Second),
		Entrypoint: entryPoint,
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}
	return &Output{
		ChainID: in.ChainID,
		Nodes: []*Node{
			{
				WSUrl:                 fmt.Sprintf("ws://%s:%s", host, mp),
				HTTPUrl:               fmt.Sprintf("http://%s:%s", host, mp),
				DockerInternalWSUrl:   fmt.Sprintf("ws://%s:%s", containerName, in.Port),
				DockerInternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}
