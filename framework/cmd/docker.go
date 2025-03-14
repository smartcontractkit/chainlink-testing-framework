package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func removeTestContainers() error {
	framework.L.Info().Str("label", "framework=ctf").Msg("Cleaning up docker containers")
	// Bash command for removing Docker containers and networks with "framework=ctf" label
	cmd := exec.Command("bash", "-c", `
		docker ps -aq --filter "label=framework=ctf" | xargs -r docker rm -f && \
		docker network ls --filter "label=framework=ctf" -q | xargs -r docker network rm && \
		docker volume ls -q | xargs -r docker volume rm || true
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

// runCommand executes a command and prints the output.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
