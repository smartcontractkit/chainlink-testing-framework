package main

import (
	"embed"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog"
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
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Shapes your test config, removes outputs, formatting ,etc",
				Subcommands: []*cli.Command{
					{
						Name:    "fmt",
						Aliases: []string{"f"},
						Usage:   "Formats TOML config",
						Action: func(c *cli.Context) error {
							in := c.Args().Get(0)
							return PrettyPrintTOML(in, in)
						},
					},
				},
			},
			{
				Name:    "docker",
				Aliases: []string{"d"},
				Usage:   "Control docker containers marked with 'framework=ctf' label",
				Subcommands: []*cli.Command{
					{
						Name:    "remove",
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
	framework.L.Info().Str("label", "framework=ctf").Msg("Cleaning up docker containers")
	// Bash command for removing Docker containers and networks with "framework=ctf" label
	cmd := exec.Command("bash", "-c", `
		docker ps -aq --filter "label=framework=ctf" | xargs -r docker rm -f && \
		docker network ls --filter "label=framework=ctf" -q | xargs -r docker network rm
	`)
	framework.L.Debug().Msg("Running command")
	if framework.L.GetLevel() == zerolog.DebugLevel {
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	framework.L.Info().Msgf("Done")
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
	framework.L.Info().Msg("Creating local observability stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, "compose"))
	framework.L.Debug().Msg("Running command")
	if framework.L.GetLevel() == zerolog.DebugLevel {
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	framework.L.Info().Msg("Done")
	fmt.Println()
	framework.L.Info().Msgf("Loki: %s", LocalLogsURL)
	framework.L.Info().Msgf("All logs: %s", "{job=\"ctf\"}")
	framework.L.Info().Msgf("By log level: %s", "{job=\"ctf\", container=~\"node.*\"} |= \"WARN|INFO|DEBUG\"")
	framework.L.Info().Msgf("Pyroscope: %s", LocalPyroScopeURL)
	return nil
}

func observabilityDown() error {
	framework.L.Info().Msg("Removing local observability stack")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, "compose"))
	framework.L.Debug().Msg("Running command")
	if framework.L.GetLevel() == zerolog.DebugLevel {
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	framework.L.Info().Msg("Done")
	return nil
}

// PrettyPrintTOML pretty prints TOML
func PrettyPrintTOML(inputFile string, outputFile string) error {
	tomlData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	tree, err := toml.Load(string(tomlData))
	if err != nil {
		return fmt.Errorf("error parsing TOML: %v", err)
	}

	// Write the result to a new file
	dumpData, err := tree.ToTomlString()
	if err != nil {
		return fmt.Errorf("error converting to TOML string: %v", err)
	}

	err = os.WriteFile(outputFile, []byte(dumpData), 0644)
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}
	framework.L.Info().Str("File", outputFile).Msg("File cleaned up and saved")
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
