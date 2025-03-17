package blockchain

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func baseRequest(in *Input, useWS bool) testcontainers.ContainerRequest {
	containerName := framework.DefaultTCName("blockchain-node")
	bindPort := fmt.Sprintf("%s/tcp", in.Port)
	exposedPorts := []string{bindPort}
	if useWS {
		exposedPorts = append(exposedPorts, fmt.Sprintf("%s/tcp", in.WSPort))
	}

	return testcontainers.ContainerRequest{
		Name:     containerName,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: exposedPorts,
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(exposedPorts...)
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		WaitingFor: wait.ForListeningPort(nat.Port(in.Port)).WithStartupTimeout(15 * time.Second).WithPollInterval(200 * time.Millisecond),
	}
}

func createGenericEvmContainer(in *Input, req testcontainers.ContainerRequest, useWS bool) (*Output, error) {
	ctx := context.Background()
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

	bindPort := req.ExposedPorts[0]
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}

	containerName := req.Name

	output := Output{
		UseCache:      true,
		Family:        "evm",
		ChainID:       in.ChainID,
		ContainerName: containerName,
		Container:     c,
		Nodes: []*Node{
			{
				HostWSUrl:             fmt.Sprintf("ws://%s:%s", host, mp.Port()),
				HostHTTPUrl:           fmt.Sprintf("http://%s:%s", host, mp.Port()),
				DockerInternalWSUrl:   fmt.Sprintf("ws://%s:%s", containerName, in.Port),
				DockerInternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}

	if useWS {
		mp, err := c.MappedPort(ctx, nat.Port(req.ExposedPorts[1]))
		if err != nil {
			return nil, err
		}
		output.Nodes[0].HostWSUrl = fmt.Sprintf("ws://%s:%s", host, mp.Port())
	}

	return &output, nil
}
