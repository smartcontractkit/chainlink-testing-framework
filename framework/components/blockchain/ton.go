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
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
)

const (
	// NOTE: Prefunded high-load wallet from MyLocalTon pre-funded wallet, that can send up to 254 messages per 1 external message
	// https://docs.ton.org/v3/documentation/smart-contracts/contracts-specs/highload-wallet#highload-wallet-v2
	DefaultTonHlWalletAddress  = "-1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b"
	DefaultTonHlWalletMnemonic = "twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin"
	// internals
	defaultTonHTTPServerPort   = "8000"
	defaultLiteServerPort      = "40000"
	defaultLiteServerPublicKey = "E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY="
)

func defaultTon(in *Input) {
	if in.Image == "" {
		in.Image = "ghcr.io/neodix42/mylocalton-docker:latest"
	}
	if in.Port == "" {
		in.Port = defaultTonHTTPServerPort
	}
}

func newTon(ctx context.Context, in *Input) (*Output, error) {
	defaultTon(in)

	containerName := framework.DefaultTCName("ton-genesis")

	baseEnv := map[string]string{
		"GENESIS": "true",
		"NAME":    "genesis",

		"EMBEDDED_FILE_HTTP_SERVER":      "true",
		"EMBEDDED_FILE_HTTP_SERVER_PORT": defaultTonHTTPServerPort,
		"LITE_PORT":                      defaultLiteServerPort,

		"CUSTOM_PARAMETERS": "--state-ttl 315360000 --archive-ttl 315360000",
	}

	// merge with additional environment variables from input
	finalEnv := baseEnv
	if in.CustomEnv != nil {
		for key, value := range in.CustomEnv {
			finalEnv[key] = value
		}
	}

	if pods.K8sEnabled() {
		return nil, fmt.Errorf("K8s support is not yet implemented")
	}

	req := testcontainers.ContainerRequest{
		Image:           in.Image,
		AlwaysPullImage: in.PullImage,
		Name:            containerName,
		ExposedPorts: []string{
			fmt.Sprintf("%s/tcp", defaultTonHTTPServerPort),
			fmt.Sprintf("%s/tcp", defaultLiteServerPort),
		},
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Labels: framework.DefaultTCLabels(),
		Env:    finalEnv,
		WaitingFor: wait.ForExec([]string{
			"/usr/local/bin/lite-client",
			"-a", fmt.Sprintf("127.0.0.1:%s", defaultLiteServerPort),
			"-b", defaultLiteServerPublicKey,
			"-t", "3", "-c", "last",
		}).WithStartupTimeout(2 * time.Minute),
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("ton-data-%s", containerName)},
				Target: "/usr/share/data",
			},
			{
				Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("ton-db-%s", containerName)},
				Target: "/var/ton-work/db",
			},
		},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = nat.PortMap{
				nat.Port(fmt.Sprintf("%s/tcp", defaultTonHTTPServerPort)): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: in.Port,
					},
				},
				nat.Port(fmt.Sprintf("%s/tcp", defaultLiteServerPort)): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "", // Docker assigns a dynamic available port
					},
				},
			}
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
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

	httpMappedPort, err := c.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", defaultTonHTTPServerPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped HTTP port: %w", err)
	}
	lsMappedPort, err := c.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", defaultLiteServerPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped LiteServer port: %w", err)
	}

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Type:          in.Type,
		Family:        FamilyTon,
		ContainerName: containerName,
		Container:     c,
		Nodes: []*Node{{
			ExternalHTTPUrl: fmt.Sprintf("liteserver://%s@%s:%s", defaultLiteServerPublicKey, host, lsMappedPort.Port()),
			InternalHTTPUrl: fmt.Sprintf("liteserver://%s@%s:%s", defaultLiteServerPublicKey, containerName, defaultLiteServerPort),
			ExternalWSUrl:   fmt.Sprintf("http://%s:%s", host, httpMappedPort.Port()),
			InternalWSUrl:   fmt.Sprintf("http://%s:%s", containerName, defaultTonHTTPServerPort),
		}},
	}, nil
}
