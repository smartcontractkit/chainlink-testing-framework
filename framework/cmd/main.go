package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/urfave/cli/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

//go:embed observability/*
var embeddedObservabilityFiles embed.FS

const (
	LocalCLNodeErrorsURL   = "http://localhost:3000/d/a7de535b-3e0f-4066-bed7-d505b6ec9ef1/cl-node-errors?orgId=1"
	LocalWorkflowEngineURL = "http://localhost:3000/d/ce589a98-b4be-4f80-bed1-bc62f3e4414a/workflow-engine?orgId=1&refresh=30s"
	LocalLogsURL           = "http://localhost:3000/explore?panes=%7B%22qZw%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-6h%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1"
	LocalPrometheusURL     = "http://localhost:3000/explore?panes=%7B%22qZw%22:%7B%22datasource%22:%22PBFA97CFB590B2093%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%22,%22range%22:true,%22datasource%22:%7B%22type%22:%22prometheus%22,%22uid%22:%22PBFA97CFB590B2093%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-6h%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1"
	LocalPostgresDebugURL  = "http://localhost:3000/d/000000039/postgresql-database?orgId=1&refresh=10s&var-DS_PROMETHEUS=PBFA97CFB590B2093&var-interval=$__auto_interval_interval&var-namespace=&var-release=&var-instance=postgres_exporter_0:9187&var-datname=All&var-mode=All&from=now-5m&to=now"
	LocalPyroScopeURL      = "http://localhost:4040"
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
					{
						Name:    "clean",
						Aliases: []string{"c"},
						Usage:   "Removes all cache files",
						Action: func(c *cli.Context) error {
							return RemoveCacheFiles()
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
							err := framework.RemoveTestContainers()
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
					{
						Name:        "load",
						Usage:       "ctf obs l",
						Aliases:     []string{"l"},
						Description: "Loads logs to Loki",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "raw-url",
								Aliases: []string{"u"},
								Usage:   "URL to GitHub raw log data",
							},
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to logs, output of 'gh run download $run_id'",
							},
							&cli.IntFlag{
								Name:    "rps",
								Aliases: []string{"r"},
								Usage:   "RPS for uploading log chunks",
								Value:   30,
							},
							&cli.IntFlag{
								Name:    "chunk",
								Aliases: []string{"c"},
								Usage:   "Amount of chunks the files will be split in",
								Value:   100,
							},
						},
						Action: func(c *cli.Context) error {
							return loadLogs(
								c.String("raw-url"),
								c.String("dir"),
								c.Int("rps"),
								c.Int("chunk"),
							)
						},
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

			{
				Name:  "ci",
				Usage: "Analyze CI job durations and statistics",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "repository",
						Aliases:  []string{"r"},
						Usage:    "GitHub repository in format owner/repo",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "workflow",
						Aliases:  []string{"w"},
						Usage:    "Name of GitHub workflow to analyze",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "start",
						Aliases: []string{"s"},
						Value:   "1",
						Usage:   "How many days to analyze",
					},
					&cli.StringFlag{
						Name:    "end",
						Aliases: []string{"e"},
						Value:   "0",
						Usage:   "How many days to analyze",
					},
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "Analytics type: jobs or steps",
					},
					&cli.BoolFlag{
						Name:  "debug",
						Usage: "Dumps all the workflow/jobs files for debugging purposes",
					},
				},
				Action: func(c *cli.Context) error {
					repo := c.String("repository")
					parts := strings.Split(repo, "/")
					if len(parts) != 2 {
						return fmt.Errorf("repository must be in format owner/repo, got: %s", repo)
					}
					typ := c.String("type")
					if typ != "jobs" && typ != "steps" {
						return fmt.Errorf("type must be 'jobs' or 'steps'")
					}
					_, err := AnalyzeCIRuns(&AnalysisConfig{
						Debug:               c.Bool("debug"),
						Owner:               parts[0],
						Repo:                parts[1],
						WorkflowName:        c.String("workflow"),
						TimeDaysBeforeStart: c.Int("start"),
						TimeDaysBeforeEnd:   c.Int("end"),
						Typ:                 typ,
						ResultsFile:         "ctf-ci.json",
					})
					return err
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
		//nolint
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

	//nolint
	err = os.WriteFile(outputFile, []byte(dumpData), 0644)
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}
	framework.L.Info().Str("File", outputFile).Msg("File cleaned up and saved")
	return nil
}

func RemoveCacheFiles() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "-cache.toml") {
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to remove file %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	framework.L.Info().Msg("All cache files has been removed")
	return nil
}
