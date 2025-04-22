package blockchain

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultAptosAPIPort    = "8080"
	DefaultAptosFaucetPort = "8081"
)

var (
	DefaultAptosAccount    = "0xa337b42bd0eecf8fb59ee5929ea4541904b3c35a642040223f3d26ab57f59d6e"
	DefaultAptosPrivateKey = "0xd477c65f88ed9e6d4ec6e2014755c3cfa3e0c44e521d0111a02868c5f04c41d4"
)

func defaultAptos(in *Input) {
	if in.Image == "" {
		in.Image = "aptoslabs/tools:aptos-node-v1.27.2"
	}
	framework.L.Warn().Msgf("Aptos node API can only be exposed on port %s!", DefaultAptosAPIPort)
	if in.Port == "" {
		// enable default API exposed port
		in.Port = DefaultAptosAPIPort
	}
	if in.CustomPorts == nil {
		// enable default API and faucet forwarding
		in.CustomPorts = append(in.CustomPorts, fmt.Sprintf("%s:%s", in.Port, DefaultAptosAPIPort), fmt.Sprintf("%s:%s", DefaultAptosFaucetPort, DefaultAptosFaucetPort))
	}
}

func newAptos(in *Input) (*Output, error) {
	defaultAptos(in)
	ctx := context.Background()
	containerName := framework.DefaultTCName("blockchain-node")

	absPath, err := filepath.Abs(in.ContractsDir)
	if err != nil {
		return nil, err
	}

	exposedPorts, bindings, err := framework.GenerateCustomPortsData(in.CustomPorts)
	if err != nil {
		return nil, err
	}
	exposedPorts = append(exposedPorts, in.Port)

	cmd := []string{
		"aptos",
		"node",
		"run-local-testnet",
		"--with-faucet",
		"--force-restart",
		"--bind-to",
		"0.0.0.0",
	}

	if len(in.DockerCmdParamsOverrides) > 0 {
		cmd = append(cmd, in.DockerCmdParamsOverrides...)
	}

	req := testcontainers.ContainerRequest{
		Image:        in.Image,
		ExposedPorts: exposedPorts,
		WaitingFor:   wait.ForLog("Faucet is ready"),
		Name:         containerName,
		Labels:       framework.DefaultTCLabels(),
		Networks:     []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = bindings
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		ImagePlatform: "linux/amd64",
		Cmd:           cmd,
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      absPath,
				ContainerFilePath: "/",
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

	dc, err := framework.NewDockerClient()
	if err != nil {
		return nil, err
	}
	cmdStr := []string{"aptos", "init", "--network=local", "--assume-yes", fmt.Sprintf("--private-key=%s", DefaultAptosPrivateKey)}
	_, err = dc.ExecContainer(containerName, cmdStr)
	if err != nil {
		return nil, err
	}
	fundCmd := []string{"aptos", "account", "fund-with-faucet", "--account", DefaultAptosAccount, "--amount", "1000000000000"}
	_, err = dc.ExecContainer(containerName, fundCmd)
	if err != nil {
		return nil, err
	}
	// expose default API port if remapped
	var exposedAPIPort string
	for _, portPair := range in.CustomPorts {
		if strings.Contains(portPair, fmt.Sprintf(":%s", DefaultAptosAPIPort)) {
			exposedAPIPort = strings.Split(portPair, ":")[0]
		}
	}
	return &Output{
		UseCache:      true,
		Type:          in.Type,
		Family:        FamilyAptos,
		ContainerName: containerName,
		Nodes: []*Node{
			{
				ExternalHTTPUrl: fmt.Sprintf("http://%s:%s", host, exposedAPIPort),
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}
