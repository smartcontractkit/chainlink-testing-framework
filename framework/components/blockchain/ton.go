package blockchain

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	// default ports from mylocalton-docker
	DefaultTonHTTPAPIPort      = "8081"
	DefaultTonSimpleServerPort = "8000"
	DefaultTonTONExplorerPort  = "8080"
	DefaultTonLiteServerPort   = "40004"

	// NOTE: Prefunded high-load wallet from MyLocalTon pre-funded wallet, that can send up to 254 messages per 1 external message
	// https://docs.ton.org/v3/documentation/smart-contracts/contracts-specs/highload-wallet#highload-wallet-v2
	DefaultTonHlWalletAddress  = "-1:5ee77ced0b7ae6ef88ab3f4350d8872c64667ffbe76073455215d3cdfab3294b"
	DefaultTonHlWalletMnemonic = "twenty unfair stay entry during please water april fabric morning length lumber style tomorrow melody similar forum width ride render void rather custom coin"
)

func defaultTon(in *Input) {
	if in.DockerComposeFileURL == "" {
		in.DockerComposeFileURL = "https://raw.githubusercontent.com/neodix42/mylocalton-docker/main/docker-compose.yaml"
	}
	// Note: in local env having all services could be useful(explorer, faucet), in CI we need only core services
	if os.Getenv("CI") == "true" && len(in.TonCoreServices) == 0 {
		// Note: mylocalton-docker's essential services, excluded explorer, restarter, faucet app,
		in.TonCoreServices = []string{
			"genesis", "tonhttpapi", "event-cache",
			"index-postgres", "index-worker", "index-api",
		}
	}
}

func newTon(in *Input) (*Output, error) {
	defaultTon(in)
	containerName := framework.DefaultTCName("blockchain-node")

	resp, err := http.Get(in.DockerComposeFileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download docker-compose file: %v", err)
	}
	defer resp.Body.Close()

	tempDir, err := os.MkdirTemp(".", "ton-mylocalton-docker")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	defer func() {
		// delete the folder whether it was successful or not
		_ = os.RemoveAll(tempDir)
	}()

	composeFile := filepath.Join(tempDir, "docker-compose.yaml")
	file, err := os.Create(composeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create compose file: %v", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to write compose file: %v", err)
	}
	file.Close()

	ctx := context.Background()

	var stack compose.ComposeStack
	stack, err = compose.NewDockerComposeWith(
		compose.WithStackFiles(composeFile),
		compose.StackIdentifier(containerName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create compose stack: %v", err)
	}

	var upOpts []compose.StackUpOption
	upOpts = append(upOpts, compose.Wait(true))

	if len(in.TonCoreServices) > 0 {
		upOpts = append(upOpts, compose.RunServices(in.TonCoreServices...))
	}

	// always wait for healthy
	const genesisBlockID = "E7XwFSQzNkcRepUC23J2nRpASXpnsEKmyyHYV4u/FZY="
	execStrat := wait.ForExec([]string{
		"/usr/local/bin/lite-client",
		"-a", "127.0.0.1:" + DefaultTonLiteServerPort,
		"-b", genesisBlockID,
		"-t", "3",
		"-c", "last",
	}).
		WithPollInterval(5 * time.Second).
		WithStartupTimeout(180 * time.Second)

	stack = stack.
		WaitForService("genesis", execStrat).
		WaitForService("tonhttpapi", wait.ForListeningPort(DefaultTonHTTPAPIPort+"/tcp"))

	if err := stack.Up(ctx, upOpts...); err != nil {
		return nil, fmt.Errorf("failed to start compose stack: %w", err)
	}

	// node container is started, now we need to connect it to the network
	genesisCtr, _ := stack.ServiceContainer(ctx, "genesis")

	// grab and connect to the network
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err := cli.NetworkConnect(
		ctx,
		framework.DefaultNetworkName,
		genesisCtr.ID,
		&networkTypes.EndpointSettings{
			Aliases: []string{containerName},
		},
	); err != nil {
		return nil, fmt.Errorf("failed to connect to network: %v", err)
	}

	// verify that the container is connected to the network
	inspected, err := cli.ContainerInspect(ctx, genesisCtr.ID)
	if err != nil {
		return nil, fmt.Errorf("inspect error: %w", err)
	}

	ns, ok := inspected.NetworkSettings.Networks[framework.DefaultNetworkName]
	if !ok {
		return nil, fmt.Errorf("container %s is NOT on network %s", genesisCtr.ID, framework.DefaultNetworkName)
	}

	fmt.Printf("✅ TON genesis '%s' is on network %s with IP %s and Aliases %v\n",
		genesisCtr.ID, framework.DefaultNetworkName, ns.IPAddress, ns.Aliases)

	httpHost, _ := genesisCtr.Host(ctx)
	httpPort, _ := genesisCtr.MappedPort(ctx, nat.Port(fmt.Sprintf("%s/tcp", DefaultTonSimpleServerPort)))

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Type:          in.Type,
		Family:        FamilyTon,
		ContainerName: containerName,
		// Note: in case we need 1+ validators, we need to modify the compose file
		Nodes: []*Node{{
			// Note: define if we need more access other than the global config(tonutils-go only uses liteclients defined in the config)
			ExternalHTTPUrl: fmt.Sprintf("%s:%s", httpHost, httpPort.Port()),
			InternalHTTPUrl: fmt.Sprintf("%s:%s", containerName, httpPort.Port()),
		}},
	}, nil
}
