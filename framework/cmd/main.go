package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
)

func main() {
	app := &cli.App{
		Name:  "ctf",
		Usage: "Manage Docker containers, networks, and TOML files for CTF framework",
		Commands: []*cli.Command{
			{
				Name:  "clean",
				Usage: "Remove Docker containers and networks with 'framework=ctf' label",
				Action: func(c *cli.Context) error {
					// Execute the bash command
					err := cleanDockerResources()
					if err != nil {
						return fmt.Errorf("failed to clean Docker resources: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "observability",
				Usage: "Process a TOML file, remove fields with '.out' keys",
				Subcommands: []*cli.Command{
					{
						Name:        "up",
						Usage:       "",
						UsageText:   "",
						Description: "",
						Action:      func(c *cli.Context) error { return observabilityUp() },
					},
					{
						Name:        "down",
						Usage:       "",
						UsageText:   "",
						Description: "",
						Action:      func(c *cli.Context) error { return observabilityDown() },
					},
				},
				Action: func(c *cli.Context) error {
					_ = c.String("file")
					// TODO: might be useful?
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanDockerResources() error {
	// Bash command for removing Docker containers and networks with "framework=ctf" label
	cmd := exec.Command("bash", "-c", `
		docker ps -aq --filter "label=framework=ctf" | xargs -r docker rm -f && \
		docker network ls --filter "label=framework=ctf" -q | xargs -r docker network rm
	`)
	framework.L.Info().Str("Cmd", cmd.String()).Msg("Running command")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}

func observabilityUp() error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up
	`, framework.ObservabilityPath))
	framework.L.Info().Str("Cmd", cmd.String()).Msg("Running command")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}

func observabilityDown() error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, framework.ObservabilityPath))
	framework.L.Info().Str("Cmd", cmd.String()).Msg("Running command")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}
