package blockchain

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultTonSimpleServerPort = "8000"
	// NOTE: Prefunded high-load wallet from MyLocalTon pre-funded wallet, that can send up to 254 messages per 1 external message
	// https://docs.ton.org/v3/documentation/smart-contracts/contracts-specs/highload-wallet#highload-wallet-v2
	DefaultTonHlWalletAddress  = "-1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b"
	DefaultTonHlWalletMnemonic = "twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin"
)

type hostPortMapping struct {
	SimpleServer string
	LiteServer   string
	DHTServer    string
	Console      string
	ValidatorUDP string
}

func generateUniquePortsFromBase(basePort string) (*hostPortMapping, error) {
	base, err := strconv.Atoi(basePort)
	if err != nil {
		return nil, fmt.Errorf("invalid base port %s: %w", basePort, err)
	}
	return &hostPortMapping{
		SimpleServer: basePort, // external HTTP → internal 8000
		LiteServer:   strconv.Itoa(base + 10),
		DHTServer:    strconv.Itoa(base + 20),
		Console:      strconv.Itoa(base + 30),
		ValidatorUDP: strconv.Itoa(base + 40),
	}, nil
}

func defaultTon(in *Input) {
	if in.Image == "" {
		in.Image = "ghcr.io/neodix42/mylocalton-docker:latest"
	}
	if in.Port == "" {
		in.Port = DefaultTonSimpleServerPort
	}
}

// newTon starts only the genesis node and nothing else.
func newTon(in *Input) (*Output, error) {
	defaultTon(in)

	hostPorts, err := generateUniquePortsFromBase(in.Port)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	net, err := network.New(ctx,
		network.WithAttachable(),
		network.WithLabels(framework.DefaultTCLabels()),
	)
	if err != nil {
		return nil, err
	}
	networkName := net.Name

	bindPorts := []string{
		fmt.Sprintf("%s:%s/tcp", hostPorts.SimpleServer, DefaultTonSimpleServerPort),
		fmt.Sprintf("%s:%s/tcp", hostPorts.LiteServer, hostPorts.LiteServer),
		fmt.Sprintf("%s:40003/udp", hostPorts.DHTServer),
		fmt.Sprintf("%s:40002/tcp", hostPorts.Console),
		fmt.Sprintf("%s:40001/udp", hostPorts.ValidatorUDP),
	}

	req := testcontainers.ContainerRequest{
		Image:           in.Image,
		AlwaysPullImage: in.PullImage,
		Name:            framework.DefaultTCName("ton-genesis"),
		ExposedPorts:    bindPorts,
		Networks:        []string{networkName},
		NetworkAliases:  map[string][]string{networkName: {"genesis"}},
		Labels:          framework.DefaultTCLabels(),
		Env: map[string]string{
			"GENESIS":           "true",
			"NAME":              "genesis",
			"LITE_PORT":         hostPorts.LiteServer,
			"CUSTOM_PARAMETERS": "--state-ttl 315360000 --archive-ttl 315360000",
		},
		WaitingFor: wait.ForExec([]string{
			"/usr/local/bin/lite-client",
			"-a", fmt.Sprintf("127.0.0.1:%s", hostPorts.LiteServer),
			"-b", "E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY=",
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
		Nodes: []*Node{{
			// Note: define if we need more access other than the global config(tonutils-go only uses liteclients defined in the config)
			ExternalHTTPUrl: fmt.Sprintf("%s:%s", "localhost", hostPorts.SimpleServer),
			InternalHTTPUrl: fmt.Sprintf("%s:%s", name, DefaultTonSimpleServerPort),
		}},
	}, nil
}
