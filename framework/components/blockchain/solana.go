package blockchain

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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
[11,2,35,236,230,251,215,68,220,208,166,157,229,181,164,26,150,230,218,229,41,20,235,80,183,97,20,117,191,159,228,243,130,101,145,43,51,163,139,142,11,174,113,54,206,213,188,127,131,147,154,31,176,81,181,147,78,226,25,216,193,243,136,149]
`

func defaultSolana(in *Input) {
	ci := os.Getenv("CI") == "true"
	if in.Image == "" && !ci {
		in.Image = "f4hrenh9it/solana"
	}
	if in.Image == "" && ci {
		in.Image = "anzaxyz/agave:v2.1.13"
	}
	if in.Port == "" {
		in.Port = "8999"
	}
}

func newSolana(in *Input) (*Output, error) {
	defaultSolana(in)
	ctx := context.Background()

	containerName := framework.DefaultTCName("blockchain-node")
	// Solana do not allow to set ws port, it just uses --rpc-port=N and sets WS as N+1 automatically
	bindPort := fmt.Sprintf("%s/tcp", in.Port)
	pp, err := strconv.Atoi(in.Port)
	if err != nil {
		return nil, fmt.Errorf("in.Port is not a number")
	}
	in.WSPort = strconv.Itoa(pp + 1)
	wsBindPort := fmt.Sprintf("%s/tcp", in.WSPort)

	configYml, err := os.CreateTemp("", "config.yml")
	if err != nil {
		return nil, err
	}
	configYmlRaw = fmt.Sprintf(configYmlRaw, in.Port, in.WSPort)
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

	contractsDir, err := filepath.Abs(in.ContractsDir)
	if err != nil {
		return nil, err
	}

	flags := []string{}
	for k, v := range in.SolanaPrograms {
		flags = append(flags, "--upgradeable-program", v, filepath.Join("/programs", k+".so"), in.PublicKey)
	}
	args := append([]string{
		"--reset",
		"--rpc-port", in.Port,
		"--mint", in.PublicKey,
	}, flags...)
	args = append(args, in.DockerCmdParamsOverrides...)

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
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort, wsBindPort)
			framework.ResourceLimitsFunc(h, in.ContainerResources)
			h.Mounts = append(h.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   contractsDir,
				Target:   "/programs",
				ReadOnly: false,
			})
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      configYml.Name(),
				ContainerFilePath: "/root/.config/solana/cli/config.yml",
				FileMode:          0644,
			},
			{
				HostFilePath:      idJSON.Name(),
				ContainerFilePath: "/root/.config/solana/cli/id.json",
				FileMode:          0644,
			},
		},
		Entrypoint: []string{"sh", "-c", fmt.Sprintf("mkdir -p /root/.config/solana/cli && solana-test-validator %s", strings.Join(args, " "))},
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

	return &Output{
		UseCache:      true,
		Type:          in.Type,
		Family:        FamilySolana,
		ContainerName: containerName,
		Container:     c,
		Nodes: []*Node{
			{
				ExternalWSUrl:   fmt.Sprintf("ws://%s:%s", host, in.WSPort),
				ExternalHTTPUrl: fmt.Sprintf("http://%s:%s", host, in.Port),
				InternalWSUrl:   fmt.Sprintf("ws://%s:%s", containerName, in.WSPort),
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}
