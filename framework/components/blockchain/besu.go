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

func defaultBesu(in *Input) {
	if in.Image == "" {
		in.Image = "hyperledger/besu:22.1.0"
	}
	if in.ChainID == "" {
		in.ChainID = "1337"
	}
	if in.Port == "" {
		in.Port = "8545"
	}
	if in.WSPort == "" {
		in.WSPort = "8546"
	}
}

func newBesu(in *Input) (*Output, error) {
	defaultBesu(in)
	ctx := context.Background()
	defaultCmd := []string{
		"--network=dev",
		"--miner-enabled",
		"--miner-coinbase=0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
		"--rpc-http-cors-origins=all",
		"--host-allowlist=*",
		"--rpc-ws-enabled",
		"--rpc-http-enabled",
		"--rpc-http-host", "0.0.0.0",
		"--rpc-ws-host", "0.0.0.0",
		"--rpc-http-port", in.Port,
		"--rpc-ws-port", in.WSPort,
		"--data-path=/tmp/tmpDatdir",
	}
	entryPoint := append(defaultCmd, in.DockerCmdParamsOverrides...)

	containerName := framework.DefaultTCName("blockchain-node")
	bindPort := fmt.Sprintf("%s/tcp", in.Port)
	bindPortWs := fmt.Sprintf("%s/tcp", in.WSPort)

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           in.Image,
		Name:            containerName,
		ExposedPorts:    []string{bindPort, bindPortWs},
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Labels: framework.DefaultTCLabels(),
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = nat.PortMap{
				nat.Port(bindPortWs): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: bindPortWs,
					},
				},
				nat.Port(bindPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: bindPort,
					},
				},
			}
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
	mpWs, err := c.MappedPort(ctx, nat.Port(bindPortWs))
	if err != nil {
		return nil, err
	}

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Family:        "evm",
		ContainerName: containerName,
		Nodes: []*Node{
			{
				HostHTTPUrl:           fmt.Sprintf("http://%s:%s", host, mp.Port()),
				HostWSUrl:             fmt.Sprintf("ws://%s:%s", host, mpWs.Port()),
				DockerInternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
				DockerInternalWSUrl:   fmt.Sprintf("ws://%s:%s", containerName, in.WSPort),
			},
		},
	}, nil
}
