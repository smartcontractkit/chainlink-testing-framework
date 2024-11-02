package framework

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func IsDockerRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Ping(context.Background())
	return err == nil
}

func GetHost(container tc.Container) (string, error) {
	host, err := container.Host(context.Background())
	if err != nil {
		return "", err
	}
	// if localhost then force it to ipv4 localhost
	if host == "localhost" {
		host = "127.0.0.1"
	}
	return host, nil
}

func MapTheSamePort(port string) nat.PortMap {
	return nat.PortMap{
		nat.Port(port): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: port,
			},
		},
	}
}

func DefaultTCLabels() map[string]string {
	return map[string]string{
		"framework": "ctf",
		"logging":   "promtail",
	}
}

func DefaultTCName(name string) string {
	return fmt.Sprintf("%s-%s", name, uuid.NewString()[0:5])
}

// BuildAndPublishLocalDockerImage runs Docker commands to set up a local registry, build an image, and push it.
func BuildAndPublishLocalDockerImage(once *sync.Once, dockerfile string, buildContext string, imageName string) error {
	var retErr error
	once.Do(func() {
		L.Info().
			Str("Dockerfile", dockerfile).
			Str("Ctx", buildContext).
			Str("ImageName", imageName).
			Msg("Building local docker file")
		registryRunning := isContainerRunning("local-registry")
		if registryRunning {
			fmt.Println("Local registry container is already running.")
		} else {
			L.Info().Msg("Starting local registry container...")
			err := runCommand("docker", "run", "-d", "-p", "5050:5000", "--name", "local-registry", "registry:2")
			if err != nil {
				retErr = fmt.Errorf("failed to start local registry: %w", err)
			}
			L.Info().Msg("Local registry started")
		}

		img := fmt.Sprintf("localhost:5050/%s:latest", imageName)
		err := runCommand("docker", "build", "-t", fmt.Sprintf("localhost:5050/%s:latest", imageName), "-f", dockerfile, buildContext)
		if err != nil {
			retErr = fmt.Errorf("failed to build Docker image: %w", err)
		}
		L.Info().Msg("Docker image built successfully")

		L.Info().Str("Image", img).Msg("Pushing Docker image to local registry")
		fmt.Println("Pushing Docker image to local registry...")
		err = runCommand("docker", "push", img)
		if err != nil {
			retErr = fmt.Errorf("failed to push Docker image: %w", err)
		}
		L.Info().Msg("Docker image pushed successfully")
	})
	return retErr
}

// isContainerRunning checks if a Docker container with the given name is running.
func isContainerRunning(containerName string) bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), containerName)
}

// runCommand executes a command and prints the output.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RebuildDockerImage rebuilds docker image if necessary
func RebuildDockerImage(once *sync.Once, dockerfile string, buildContext string, imageName string) (string, error) {
	if dockerfile == "" {
		return "", errors.New("docker_file path must be provided")
	}
	if buildContext == "" {
		return "", errors.New("docker_ctx path must be provided")
	}
	if imageName == "" {
		imageName = "ctftmp"
	}
	if err := BuildAndPublishLocalDockerImage(once, dockerfile, buildContext, imageName); err != nil {
		return "", err
	}
	return fmt.Sprintf("localhost:5050/%s:latest", imageName), nil
}
