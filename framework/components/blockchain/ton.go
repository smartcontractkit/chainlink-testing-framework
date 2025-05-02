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
)

func defaultTon(in *Input) {
	if in.Image == "" {
		in.Image = "neodix42/mylocalton-docker:latest"
	}
	if in.Port != "" {
		framework.L.Warn().Msgf("'port' field is set but only default port can be used: %s", DefaultTonHTTPAPIPort)
	}
	in.Port = DefaultTonHTTPAPIPort
}

func newTon(in *Input) (*Output, error) {
	defaultTon(in)
	containerName := framework.DefaultTCName("blockchain-node")

	resp, err := http.Get("https://raw.githubusercontent.com/neodix42/mylocalton-docker/main/docker-compose.yaml")
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

	// always wait for healthy
	upOpts = append(upOpts, compose.Wait(true))
	services := in.CoreServices
	if os.Getenv("CI") == "true" && len(services) == 0 {
		services = []string{
			"genesis", "tonhttpapi", "event-cache",
			"index-postgres", "index-worker", "index-api",
		}
	}

	if len(services) > 0 {
		upOpts = append(upOpts, compose.RunServices(services...))
	}

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
	globalCfgURL := fmt.Sprintf("http://%s:%s/localhost.global.config.json", cfgHost, cfgPort.Port())

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
		Nodes: []*Node{{
			// todo: do we need more access?
			ExternalHTTPUrl: fmt.Sprintf("%s:%s", liteHost, litePort.Port()),
		}},
		NetworkSpecificData: &NetworkSpecificData{
			TonGlobalConfigURL: globalCfgURL,
		},
	}, nil
}
