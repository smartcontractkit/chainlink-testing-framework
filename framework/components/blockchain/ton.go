package blockchain

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	// NOTE: Prefunded high-load wallet from MyLocalTon pre-funded wallet, that can send up to 254 messages per 1 external message
	// https://docs.ton.org/v3/documentation/smart-contracts/contracts-specs/highload-wallet#highload-wallet-v2
	DefaultTonHlWalletAddress  = "-1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b"
	DefaultTonHlWalletMnemonic = "twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin"
	// internals
	defaultTonHTTPServerPort   = "8000"
	defaultLiteServerPublicKey = "E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY="
	liteServerPortOffset       = 100 // arbitrary offset for lite server port
)

type portMapping struct {
	HTTPServer string
	LiteServer string
}

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

	base, err := strconv.Atoi(in.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid base port %s: %w", in.Port, err)
	}

	ports := &portMapping{
		HTTPServer: in.Port,
		LiteServer: strconv.Itoa(base + liteServerPortOffset),
	}

	baseEnv := map[string]string{
		"GENESIS": "true",
		"NAME":    "genesis",

		"EMBEDDED_FILE_HTTP_SERVER":      "true",
		"EMBEDDED_FILE_HTTP_SERVER_PORT": ports.HTTPServer,
		"LITE_PORT":                      ports.LiteServer,

		"CUSTOM_PARAMETERS": "--state-ttl 315360000 --archive-ttl 315360000",
	}

	// merge with additional environment variables from input
	finalEnv := baseEnv
	if in.CustomEnv != nil {
		for key, value := range in.CustomEnv {
			finalEnv[key] = value
		}
	}
	networkName := framework.DefaultNetworkName

	req := testcontainers.ContainerRequest{
		Image:           in.Image,
		AlwaysPullImage: in.PullImage,
		Name:            framework.DefaultTCName("ton-genesis"),
		ExposedPorts: []string{
			fmt.Sprintf("%s:%s/tcp", ports.HTTPServer, ports.HTTPServer),
			fmt.Sprintf("%s:%s/tcp", ports.LiteServer, ports.LiteServer),
			"40003/udp",
			"40002/tcp",
			"40001/udp",
		},
		Networks:       []string{networkName},
		NetworkAliases: map[string][]string{networkName: {"genesis"}},
		Labels:         framework.DefaultTCLabels(),
		Env:            finalEnv,
		WaitingFor: wait.ForExec([]string{
			"/usr/local/bin/lite-client",
			"-a", fmt.Sprintf("127.0.0.1:%s", ports.LiteServer),
			"-b", defaultLiteServerPublicKey,
			"-t", "3", "-c", "last",
		}).WithStartupTimeout(2 * time.Minute),
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("shared-data-%s", networkName)},
				Target: "/usr/share/data",
			},
			{
				Source: testcontainers.GenericVolumeMountSource{Name: fmt.Sprintf("ton-db-%s", networkName)},
				Target: "/var/ton-work/db",
			},
		},
		HostConfigModifier: func(h *container.HostConfig) {
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

	host, err := c.Host(ctx)
	if err != nil {
		return nil, err
	}

	name, err := c.Name(ctx)
	if err != nil {
		return nil, err
	}

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Type:          in.Type,
		Family:        FamilyTon,
		ContainerName: name,
		Container:     c,
		Nodes: []*Node{{
			// URLs now contain liteserver://publickey@host:port
			ExternalHTTPUrl: fmt.Sprintf("liteserver://%s@%s:%s", defaultLiteServerPublicKey, host, ports.LiteServer),
			InternalHTTPUrl: fmt.Sprintf("liteserver://%s@%s:%s", defaultLiteServerPublicKey, name, ports.LiteServer),
		}},
	}, nil
}
