package main

import (
	"embed"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/urfave/cli/v2"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed observability/*
var embeddedObservabilityFiles embed.FS

const (
	LocalLogsURL      = "http://localhost:3000/explore"
	LocalPyroScopeURL = "http://localhost:4040"
)

func main() {
	app := &cli.App{
		Name:      "ctf",
		Usage:     "Chainlink Testing Framework CLI",
		UsageText: "'ctf' is a useful utility that can:\n- clean up test docker containers\n- modify test files\n- create a local observability stack with Grafana/Loki/Pyroscope",
		Commands: []*cli.Command{
			{
				Name:    "docker",
				Aliases: []string{"d"},
				Usage:   "Control docker containers marked with 'framework=ctf' label",
				Subcommands: []*cli.Command{
					{
						Name:    "clean",
						Aliases: []string{"rm"},
						Usage:   "Remove Docker containers and networks with 'framework=ctf' label",
						Action: func(c *cli.Context) error {
							err := cleanDockerResources()
							if err != nil {
								return fmt.Errorf("failed to clean Docker resources: %w", err)
							}
							return nil
						},
					},
				},
			},
			{
				Name:    "observability",
				Aliases: []string{"obs"},
				Usage:   "Spins up a local observability stack: Grafana, Loki, Pyroscope",
				Subcommands: []*cli.Command{
					{
						Name:        "up",
						Usage:       "ctf obs up",
						Description: "Spins up a local observability stack: Grafana, Loki, Pyroscope",
						Action:      func(c *cli.Context) error { return observabilityUp() },
					},
					{
						Name:        "down",
						Usage:       "ctf obs down",
						Description: "Removes local observability stack",
						Action:      func(c *cli.Context) error { return observabilityDown() },
					},
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

// extractAllFiles goes through the embedded directory and extracts all files to the current directory
func extractAllFiles(embeddedDir string) error {
	// Get current working directory where CLI is running
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Walk through the embedded files
	err = fs.WalkDir(embeddedObservabilityFiles, embeddedDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the directory: %w", err)
		}
		if strings.Contains(path, "README.md") {
			return nil
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file content from embedded file system
		content, err := embeddedObservabilityFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Determine the target path (strip out the `embeddedDir` part)
		relativePath, err := filepath.Rel(embeddedDir, path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", path, err)
		}
		targetPath := filepath.Join(currentDir, relativePath)

		// Create target directories if necessary
		targetDir := filepath.Dir(targetPath)
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		// Write the file content to the target path
		err = ioutil.WriteFile(targetPath, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
		return nil
	})

	return err
}

func observabilityUp() error {
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, "compose"))
	framework.L.Info().Str("Cmd", cmd.String()).Msg("Running command")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	fmt.Printf("Check services logs here: %s\nUse '{job=\"ctf\"}' query to select all logs\nPyroscope traces are available at %s\n", LocalLogsURL, LocalPyroScopeURL)
	return nil
}

func observabilityDown() error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, "compose"))
	framework.L.Info().Str("Cmd", cmd.String()).Msg("Running command")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}
