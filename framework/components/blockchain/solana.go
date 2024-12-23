package blockchain

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var configYmlRaw = `
json_rpc_url: http://0.0.0.0:%s
websocket_url: ws://0.0.0.0:%s
keypair_path: /root/.config/solana/cli/id.json
address_labels:
  "11111111111111111111111111111111": ""
commitment: finalized
`

var idJSONRaw = `
[94,214,238,83,144,226,75,151,226,20,5,188,42,110,64,180,196,244,6,199,29,231,108,112,67,175,110,182,3,242,102,83,103,72,221,132,137,219,215,192,224,17,146,227,94,4,173,67,173,207,11,239,127,174,101,204,65,225,90,88,224,45,205,117]
`

func defaultSolana(in *Input) {
	if in.Image == "" {
		in.Image = "solanalabs/solana:v1.18.26"
	}
	if in.Port == "" {
		in.Port = "8899"
	}
}

func newSolana(in *Input) (*Output, error) {
	defaultSolana(in)
	ctx := context.Background()
	wsPortNumberStr, err := wsPort(in.Port)
	if err != nil {
		return nil, err
	}
	framework.L.Info().Msg("Creating solana container")
	bindPort := fmt.Sprintf("%s/tcp", in.Port)
	containerName := framework.DefaultTCName("blockchain-solana-node")
	wsBindPort := fmt.Sprintf("%s/tcp", wsPortNumberStr)

	configYml, err := os.CreateTemp("", "config.yml")
	if err != nil {
		return nil, err
	}
	configYmlRaw = fmt.Sprintf(configYmlRaw, in.Port, wsPortNumberStr)
	_, err = configYml.WriteString(configYmlRaw)
	if err != nil {
		return nil, err
	}

	idJSON, err := os.CreateTemp("", "id.json")
	if err != nil {
		return nil, err
	}
	_, err = idJSON.WriteString(idJSONRaw)
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           in.Image,
		Labels:          framework.DefaultTCLabels(),
		Name:            containerName,
		ExposedPorts:    []string{bindPort, wsBindPort},
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		WaitingFor: wait.ForLog("Processed Slot: 1").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
		Env: map[string]string{
			"SERVER_PORT": "1080", // what is this?
		},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort)
			h.PortBindings[nat.Port(wsBindPort)] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: wsBindPort,
				},
			}
			h.Mounts = append(h.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   in.ContractsDir,
				Target:   "/programs",
				ReadOnly: false,
			})
		},
		LifecycleHooks: []testcontainers.ContainerLifecycleHooks{
			{
				PostStarts: []testcontainers.ContainerHook{
					func(ctx context.Context, container testcontainers.Container) error {
						err = container.CopyFileToContainer(ctx, configYml.Name(), "/root/.config/solana/cli/config.yml", 0644)
						if err != nil {
							return err
						}
						err = container.CopyFileToContainer(ctx, idJSON.Name(), "/root/.config/solana/cli/id.json", 0644)
						return err
					},
				},
			},
		},
		Entrypoint: []string{"sh", "-c", "mkdir -p /root/.config/solana/cli && solana-test-validator -r --mint=" + in.PublicKey},
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
	wsmp, err := c.MappedPort(ctx, nat.Port(wsBindPort))
	if err != nil {
		return nil, err
	}

	framework.L.Info().Msg("Started Solana container")

	return &Output{
		UseCache:      true,
		Family:        "non-evm",
		ChainID:       in.ChainID,
		ContainerName: containerName,
		Nodes: []*Node{
			{
				HostWSUrl:             fmt.Sprintf("ws://%s:%s", host, wsmp.Port()),
				HostHTTPUrl:           fmt.Sprintf("http://%s:%s", host, mp.Port()),
				DockerInternalWSUrl:   fmt.Sprintf("ws://%s:%s", containerName, wsPortNumberStr),
				DockerInternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}

func wsPort(rpcPort string) (string, error) {
	wsPortNumber, err := strconv.Atoi(rpcPort)
	if err != nil {
		return "", fmt.Errorf("failed to convert port to integer: %w", err)
	}
	wsPortNumber += 1 // Increment by 1
	wsPortNumberStr := strconv.Itoa(wsPortNumber)
	return wsPortNumberStr, nil
}
