package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
)

func ComposeFilePath(rawFilePath string, serviceName string) (string, error) {
	// if it's not a URL, return it as is and assume it's a local file
	if !strings.HasPrefix(rawFilePath, "http") {
		return rawFilePath, nil
	}

	resp, respErr := http.Get(rawFilePath)
	if respErr != nil {
		return "", errors.Wrap(respErr, "failed to download docker-compose file")
	}
	defer resp.Body.Close()

	tempFile, tempErr := os.CreateTemp("", serviceName+"-docker-compose-*.yml")
	if tempErr != nil {
		return "", errors.Wrap(tempErr, "failed to create temp file")
	}
	defer tempFile.Close()

	_, copyErr := io.Copy(tempFile, resp.Body)
	if copyErr != nil {
		tempFile.Close()
		return "", errors.Wrap(copyErr, "failed to write compose file")
	}

	return tempFile.Name(), nil
}

func GetContainerHost(ctx context.Context, container *testcontainers.DockerContainer) (string, error) {
	return container.Host(ctx)
}

func FindMappedPort(ctx context.Context, timeout time.Duration, container *testcontainers.DockerContainer, port nat.Port) (nat.Port, error) {
	forCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tickerInterval := 5 * time.Second
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-forCtx.Done():
			return "", fmt.Errorf("timeout while waiting for mapped port for %s", port)
		case <-ticker.C:
			portCtx, portCancel := context.WithTimeout(ctx, tickerInterval)
			defer portCancel()
			mappedPort, mappedPortErr := container.MappedPort(portCtx, port)
			if mappedPortErr != nil {
				return "", errors.Wrapf(mappedPortErr, "failed to get mapped port for %s", port)
			}
			if mappedPort.Port() == "" {
				return "", fmt.Errorf("mapped port for %s is empty", port)
			}
			return mappedPort, nil
		}
	}
}

func ConnectNetwork(connCtx context.Context, timeout time.Duration, dockerClient *client.Client, containerID, networkName, stackIdentifier string) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	networkCtx, networkCancel := context.WithTimeout(connCtx, timeout)
	defer networkCancel()

	for {
		select {
		case <-networkCtx.Done():
			return fmt.Errorf("timeout while trying to connect billing-platform-service to default network")
		case <-ticker.C:
			if networkErr := dockerClient.NetworkConnect(
				connCtx,
				networkName,
				containerID,
				&networkTypes.EndpointSettings{
					Aliases: []string{stackIdentifier},
				},
			); networkErr != nil && !strings.Contains(networkErr.Error(), "already exists in network") {
				framework.L.Trace().Msgf("failed to connect to default network: %v", networkErr)
				continue
			}
			framework.L.Trace().Msgf("connected to %s network", networkName)
			return nil
		}
	}
}
