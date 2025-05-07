package blockchain

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
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
	if in.Image == "" {
		// Note: mylocalton is a compose file, not a single image. Reusing common image field
		in.Image = "https://raw.githubusercontent.com/neodix42/mylocalton-docker/main/docker-compose.yaml"
		// Note: mylocalton-docker's essential services, excluded explorer, restarter, faucet app,
		in.CoreServices = []string{
			"genesis", "tonhttpapi", "event-cache",
			"index-postgres", "index-worker", "index-api",
		}
	}
}

func newTon(in *Input) (*Output, error) {
	defaultTon(in)
	containerName := framework.DefaultTCName("blockchain-node")

	resp, err := http.Get(in.Image)
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
	services := []string{}
	// Note: in local env having all services could be useful(explorer, faucet), in CI we need only core services
	if os.Getenv("CI") == "true" && len(services) == 0 {
		services = in.CoreServices
	}

	if len(services) > 0 {
		upOpts = append(upOpts, compose.RunServices(services...))
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
	cfgCtr, _ := stack.ServiceContainer(ctx, "genesis")
	cfgHost, _ := cfgCtr.Host(ctx)
	cfgPort, _ := cfgCtr.MappedPort(ctx, nat.Port("8000/tcp"))

	// discover lite‐server addr
	liteCtr, _ := stack.ServiceContainer(ctx, "genesis")
	liteHost, _ := liteCtr.Host(ctx)
	litePort, _ := liteCtr.MappedPort(ctx, nat.Port("40004/tcp"))

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Type:          in.Type,
		Family:        FamilyTon,
		ContainerName: containerName,
		// Note: in case we need 1+ validators, we need to modify the compose file
		Nodes: []*Node{{
			// todo: define if we need more access other than lite client(tonutils-go only needs lite client)
			ExternalHTTPUrl: fmt.Sprintf("%s:%s", liteHost, litePort.Port()),
		}},
		NetworkSpecificData: &NetworkSpecificData{
			TonGlobalConfigURL: fmt.Sprintf("http://%s:%s/localhost.global.config.json", cfgHost, cfgPort.Port()),
		},
	}, nil
}
