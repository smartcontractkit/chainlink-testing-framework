package blockchain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DefaultAnvilPrivateKey = `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
)

func defaultAnvil(in *Input) {
	if in.Image == "" {
		in.Image = "f4hrenh9it/foundry:latest"
	}
	if in.ChainID == "" {
		in.ChainID = "1337"
	}
	if in.Port == "" {
		in.Port = "8545"
	}
}

// newAnvil deploy foundry anvil node
func newAnvil(in *Input) (*Output, error) {
	defaultAnvil(in)
	ctx := context.Background()
	entryPoint := []string{"anvil"}
	defaultCmd := []string{"--host", "0.0.0.0", "--port", in.Port, "--chain-id", in.ChainID}
	entryPoint = append(entryPoint, defaultCmd...)
	entryPoint = append(entryPoint, in.DockerCmdParamsOverrides...)
	framework.L.Info().Any("Cmd", strings.Join(entryPoint, " ")).Msg("Creating anvil with command")
	bindPort := fmt.Sprintf("%s/tcp", in.Port)
	containerName := framework.DefaultTCName("blockchain-node")

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           in.Image,
		Labels:          framework.DefaultTCLabels(),
		Name:            containerName,
		ExposedPorts:    []string{bindPort},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort)
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		WaitingFor: wait.ForListeningPort(nat.Port(in.Port)).WithStartupTimeout(10 * time.Second).WithPollInterval(200 * time.Millisecond),
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
	}, nil
}
