package test_env

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
)

// NatPortFormat formats a given port string to include the TCP protocol suffix.
// This is useful for standardizing port representations in container configurations.
func NatPortFormat(port string) string {
	return fmt.Sprintf("%s/tcp", port)
}

// NatPort converts a string representation of a port into a nat.Port type.
// This is useful for ensuring that the port is formatted correctly for networking operations.
func NatPort(port string) nat.Port {
	return nat.Port(NatPortFormat(port))
}

// GetHost returns the host of a container, if localhost then force ipv4 localhost
// to avoid ipv6 docker bugs https://github.com/moby/moby/issues/42442 https://github.com/moby/moby/issues/42375
func GetHost(ctx context.Context, container tc.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}
	// if localhost then force it to ipv4 localhost
	if host == "localhost" {
		host = "127.0.0.1"
	}
	return host, nil
}

// GetEndpointFromPort returns the endpoint of a container associated with a port,
// if localhost then force ipv4 localhost
// to avoid ipv6 docker bugs https://github.com/moby/moby/issues/42442 https://github.com/moby/moby/issues/42375
func GetEndpointFromPort(ctx context.Context, container tc.Container, endpointType string, portStr string) (string, error) {
	port, err := nat.NewPort("tcp", portStr)
	if err != nil {
		return "", err
	}
	endpoint, err := container.PortEndpoint(ctx, port, endpointType)
	if err != nil {
		return "", err
	}
	return strings.Replace(endpoint, "localhost", "127.0.0.1", 1), nil
}

// GetEndpoint returns the endpoint of a container, if localhost then force ipv4 localhost
// to avoid ipv6 docker bugs https://github.com/moby/moby/issues/42442 https://github.com/moby/moby/issues/42375
func GetEndpoint(ctx context.Context, container tc.Container, endpointType string) (string, error) {
	endpoint, err := container.Endpoint(ctx, endpointType)
	if err != nil {
		return "", err
	}
	return strings.Replace(endpoint, "localhost", "127.0.0.1", 1), nil
}

// FormatHttpUrl constructs a standard HTTP URL using the provided host and port.
// This function is useful for generating URLs for services running in a containerized environment.
func FormatHttpUrl(host string, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}

// FormatWsUrl constructs a WebSocket URL using the provided host and port.
// This function is useful for establishing WebSocket connections to Ethereum nodes.
func FormatWsUrl(host string, port string) string {
	return fmt.Sprintf("ws://%s:%s", host, port)
}

// UniqueStringSlice returns a deduplicated slice of strings
func UniqueStringSlice(slice []string) []string {
	stringSet := make(map[string]struct{})
	deduplicated := make([]string, 0)

	for _, el := range slice {
		if _, exists := stringSet[el]; exists {
			continue
		}

		stringSet[el] = struct{}{}
		deduplicated = append(deduplicated, el)
	}

	return deduplicated
}
