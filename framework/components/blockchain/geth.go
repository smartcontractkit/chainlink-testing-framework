package blockchain

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func defaultGeth(in *Input) {
	if in.Image == "" {
		in.Image = "ethereum/client-go:v1.13.8"
	}
	if in.ChainID == "" {
		in.ChainID = "1337"
	}
	if in.Port == "" {
		in.Port = "8545"
	}
}

func newGeth(in *Input) (*Output, error) {
	defaultGeth(in)
	ctx := context.Background()
	defaultCmd := []string{
		"--http.corsdomain=*",
		"--http.vhosts=*",
		"--http",
		"--http.addr",
		"0.0.0.0",
		"--http.port",
		in.Port,
		"--http.api",
		"eth,net,web3",
		"--ws",
		"--ws.addr",
		"0.0.0.0",
		"--ws.port",
		in.Port,
		"--ws.api",
		"eth,net,web3",
		fmt.Sprintf("--networkid=%s", in.ChainID),
		"--ipcdisable",
		"--graphql",
		"-graphql.corsdomain", "*",
		"--allow-insecure-unlock",
		"--vmdebug",
		"--mine",
		"--dev",
		"--dev.period",
		"1",
	}
	entryPoint := append(defaultCmd, in.DockerCmdParamsOverrides...)

	containerName := framework.DefaultTCName("blockchain-node")
	bindPort := fmt.Sprintf("%s/tcp", in.Port)

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           in.Image,
		Labels:          framework.DefaultTCLabels(),
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Name:         containerName,
		ExposedPorts: []string{bindPort},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort)
		},
		WaitingFor: wait.ForListeningPort(nat.Port(in.Port)).WithStartupTimeout(15 * time.Second),
		Cmd:        entryPoint,
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := c.Host(ctx)
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
		Nodes: []*Node{
			{
				HostHTTPUrl:           fmt.Sprintf("http://%s:%s", host, mp.Port()),
				HostWSUrl:             fmt.Sprintf("ws://%s:%s", host, mp.Port()),
				DockerInternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
				DockerInternalWSUrl:   fmt.Sprintf("ws://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}
