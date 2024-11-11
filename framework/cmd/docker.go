package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"os"
	"os/exec"
	"strings"
)

func rmTestContainers() error {
	framework.L.Info().Str("label", "framework=ctf").Msg("Cleaning up docker containers")
	// Bash command for removing Docker containers and networks with "framework=ctf" label
	cmd := exec.Command("bash", "-c", `
		docker ps -aq --filter "label=framework=ctf" | xargs -r docker rm -f && \
		docker network ls --filter "label=framework=ctf" -q | xargs -r docker network rm && \
		docker volume rm postgresql_data || true
	`)
	framework.L.Debug().Msg("Running command")
	if framework.L.GetLevel() == zerolog.DebugLevel {
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}

func cleanUpDockerResources() error {
	framework.L.Info().Msg("Cleaning up docker resources")
	// Bash command for removing Docker containers and networks with "framework=ctf" label
	cmd := exec.Command("bash", "-c", `
		docker system prune -a --volumes -f
	`)
	framework.L.Debug().Msg("Running command")
	if framework.L.GetLevel() == zerolog.DebugLevel {
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}

// BuildDocker runs Docker commands to set up a local registry, build an image, and push it.
func BuildDocker(dockerfile string, buildContext string, imageName string) error {
	registryRunning := isContainerRunning("local-registry")
	if registryRunning {
		fmt.Println("Local registry container is already running.")
	} else {
		framework.L.Info().Msg("Starting local registry container...")
		err := runCommand("docker", "run", "-d", "-p", "5050:5000", "--name", "local-registry", "registry:2")
		if err != nil {
			return fmt.Errorf("failed to start local registry: %w", err)
		}
		framework.L.Info().Msg("Local registry started")
	}

	img := fmt.Sprintf("localhost:5050/%s:latest", imageName)
	framework.L.Info().Str("DockerFile", dockerfile).Str("Context", buildContext).Msg("Building Docker image")
	err := runCommand("docker", "build", "-t", fmt.Sprintf("localhost:5050/%s:latest", imageName), "-f", dockerfile, buildContext)
	if err != nil {
		return fmt.Errorf("failed to build Docker image: %w", err)
	}
	framework.L.Info().Msg("Docker image built successfully")

	framework.L.Info().Str("Image", img).Msg("Pushing Docker image to local registry")
	fmt.Println("Pushing Docker image to local registry...")
	err = runCommand("docker", "push", img)
	if err != nil {
		return fmt.Errorf("failed to push Docker image: %w", err)
	}
	framework.L.Info().Msg("Docker image pushed successfully")
	return nil
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
