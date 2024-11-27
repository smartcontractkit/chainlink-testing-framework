package main

import (
	"embed"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/urfave/cli/v2"
	"io/fs"
	"log"
	"os"
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
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Build an environment interactively, suitable for non-technical users",
				Subcommands: []*cli.Command{
					{
						Name:    "node_set",
						Aliases: []string{"ns"},
						Usage:   "Builds a NodeSet and connect it to some networks",
						Action: func(c *cli.Context) error {
							return runSetupForm()
						},
					},
				},
			},
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
						Name:    "clean",
						Aliases: []string{"c"},
						Usage:   "Cleanup all docker resources: volumes, images, build caches",
						Action: func(c *cli.Context) error {
							err := cleanUpDockerResources()
							if err != nil {
								return fmt.Errorf("failed to clean Docker resources: %w", err)
							}
							return nil
						},
					},
					{
						Name:    "remove",
						Aliases: []string{"rm"},
						Usage:   "Remove Docker containers and networks with 'framework=ctf' label",
						Action: func(c *cli.Context) error {
							err := rmTestContainers()
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
						Aliases:     []string{"u"},
						Description: "Spins up a local observability stack: Grafana, Loki, Pyroscope",
						Action:      func(c *cli.Context) error { return observabilityUp() },
					},
					{
						Name:        "down",
						Usage:       "ctf obs down",
						Aliases:     []string{"d"},
						Description: "Removes local observability stack",
						Action:      func(c *cli.Context) error { return observabilityDown() },
					},
				},
			},
			{
				Name:    "blockscout",
				Aliases: []string{"bs"},
				Usage:   "Controls local Blockscout stack",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "rpc",
						Aliases: []string{"r"},
						Usage:   "RPC URL for blockchain node to index",
						Value:   "http://host.docker.internal:8545",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:        "up",
						Usage:       "ctf bs up",
						Aliases:     []string{"u"},
						Description: "Spins up Blockscout stack",
						Action: func(c *cli.Context) error {
							return blockscoutUp(c.String("rpc"))
						},
					},
					{
						Name:        "down",
						Usage:       "ctf bs down",
						Aliases:     []string{"d"},
						Description: "Removes Blockscout stack, wipes all Blockscout databases data",
						Action: func(c *cli.Context) error {
							return blockscoutDown(c.String("rpc"))
						},
					},
					{
						Name:        "reboot",
						Usage:       "ctf bs reboot or ctf bs r",
						Aliases:     []string{"r"},
						Description: "Reboots Blockscout stack",
						Action: func(c *cli.Context) error {
							rpc := c.String("rpc")
							if err := blockscoutDown(rpc); err != nil {
								return err
							}
							return blockscoutUp(rpc)
						},
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
		err = os.WriteFile(targetPath, content, 0777)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
		return nil
	})

	return err
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
