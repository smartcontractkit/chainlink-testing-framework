package blockchain

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	AccountsFile = `{
	 "hdPath": "m/44'/195'/0'/0/",
	 "mnemonic": "resemble birth wool happy sun burger fatal trumpet globe purity health ritual",
	 "privateKeys": [
	       "cf36898af3c63e13537063ae165ec8262fafac188f09200647b4b76b6f212b90",
	       "7d5110f81cc6b2c65a532066e81fe813edf781e24f4e0fa42d22a3003dae7a54",
	       "ae0de6fb5450622bfc96ec0c25a8a5cb85256d1f9d6cbbe5fd9de073d22f0060",
	       "e972c3c213f8ba8cfe9a75e5d0b48310f4e35715f70986edd1eade904dd03437",
	       "23d81a4d6c85661b58922e68db09cca0ebe77c787beb0e12c8d29da111568855"
	  ],
	 "more": [
	   {
	     "hdPath": "m/44'/195'/0'/0/",
	     "mnemonic": "resemble birth wool happy sun burger fatal trumpet globe purity health ritual",
	     "privateKeys": [
	       "cf36898af3c63e13537063ae165ec8262fafac188f09200647b4b76b6f212b90",
	       "7d5110f81cc6b2c65a532066e81fe813edf781e24f4e0fa42d22a3003dae7a54",
	       "ae0de6fb5450622bfc96ec0c25a8a5cb85256d1f9d6cbbe5fd9de073d22f0060",
	       "e972c3c213f8ba8cfe9a75e5d0b48310f4e35715f70986edd1eade904dd03437",
	       "23d81a4d6c85661b58922e68db09cca0ebe77c787beb0e12c8d29da111568855"
	      ]
		}
	 ]
	}
	`
	DefaultTronPort         = "9090"
	DefaultTronSolidityPort = "8091"
)

func defaultTron(in *Input) {
	if in.Image == "" {
		in.Image = "trontools/quickstart:2.1.1"
	}
	if in.Port == "" {
		in.Port = DefaultTronPort
	}
}

func newTron(in *Input) (*Output, error) {
	defaultTron(in)
	ctx := context.Background()

	containerName := framework.DefaultTCName("blockchain-node")
	bindPort := fmt.Sprintf("%s/tcp", in.Port)

	accounts, err := os.CreateTemp("", "accounts.json")
	if err != nil {
		return nil, err
	}
	_, err = accounts.WriteString(AccountsFile)
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           in.Image,
		Name:            containerName,
		ExposedPorts:    []string{bindPort, "18190/tcp", "18191/tcp"},
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Env: map[string]string{
			"accounts": "10",
		},
		Labels: framework.DefaultTCLabels(),
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort, "18190/tcp", "19191/tcp")
		},
		WaitingFor: wait.ForListeningPort(nat.Port(in.Port)).WithStartupTimeout(60 * time.Second).WithPollInterval(200 * time.Millisecond),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      accounts.Name(),
				ContainerFilePath: "/config/accounts.json",
				FileMode:          0644,
			},
		},
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

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Family:        "tron",
		ContainerName: containerName,
		Nodes: []*Node{
			{
				HostHTTPUrl:           fmt.Sprintf("http://%s:%s", host, in.Port),
				DockerInternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}
