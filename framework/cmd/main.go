package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/urfave/cli/v2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func main() {
	app := &cli.App{
		Name:      "ctf",
		Usage:     "Chainlink Testing Framework CLI",
		UsageText: "'ctf' is a useful utility that can:\n- clean up test docker containers\n- modify test files\n- create a local observability stack with Grafana/Loki/Pyroscope",
		Commands: []*cli.Command{
			{
				Name:    "gen",
				Aliases: []string{"g"},
				Usage:   "Generates various test templates",
				Subcommands: []*cli.Command{
					{
						Name:    "load",
						Aliases: []string{"l"},
						Usage:   "Generates a load/chaos test template for Kubernetes namespace",
						Description: `Scans a Kubernetes namespace and generates load testing templates for discovered services.

Prerequisites:

	Connect to K8s and don't forget to switch context first:
		kubectl config use-context <your_ctx>
	By default test sends data to a local CTF stack, see //TODO comments to change that, spin up the stack:
		ctf obs up

Usage:

	Generate basic kill/latency tests:
		ctf gen k8s-load my-namespace
	With workload:
		ctf gen k8s-load -w my-namespace
	With workload and name:
		ctf gen k8s-load -w -n TestSomething my-namespace

Be aware that any TODO requires your attention before your run the final test!
`,
						ArgsUsage: "--workload --name $name --output-dir $dir --module $go_mod_name [NAMESPACE]",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Value:   "TestLoadChaos",
								Usage:   "Test suite name",
							},
							&cli.StringFlag{
								Name:    "output-dir",
								Aliases: []string{"o"},
								Value:   "wasp-test",
								Usage:   "Output directory for generated files",
							},
							&cli.StringFlag{
								Name:    "module",
								Aliases: []string{"m"},
								Value:   "github.com/smartcontractkit/chainlink-testing-framework/wasp-test",
								Usage:   "Go module name for generated project",
							},
							&cli.BoolFlag{
								Name:    "workload",
								Aliases: []string{"w"},
								Value:   false,
								Usage:   "Include workload generation in tests",
							},
							&cli.StringFlag{
								Name:    "pod-label-key",
								Aliases: []string{"k"},
								Value:   "app.kubernetes.io/instance",
								Usage:   "Default unique pod key, read more here: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/",
							},
						},
						Action: func(c *cli.Context) error {
							if c.Args().Len() == 0 {
								return fmt.Errorf("Kubernetes namespace argument is required")
							}
							ns := c.Args().First()
							testSuiteName := c.String("name")
							podLabelKey := c.String("pod-label-key")
							outputDir := c.String("output-dir")
							moduleName := c.String("module")
							includeWorkload := c.Bool("workload")
							framework.L.Info().
								Str("SuiteName", testSuiteName).
								Str("OutputDir", outputDir).
								Str("GoModuleName", moduleName).
								Bool("Workload", includeWorkload).
								Msg("Generating load&chaos test template")

							k8sClient, err := wasp.NewK8s()
							if err != nil {
								return fmt.Errorf("failed to create K8s client")
							}

							cg, err := wasp.NewLoadTestGenBuilder(k8sClient, ns).
								TestSuiteName(testSuiteName).
								UniqPodLabelKey(podLabelKey).
								Workload(includeWorkload).
								OutputDir(outputDir).
								GoModName(moduleName).
								Build()
							if err != nil {
								return fmt.Errorf("failed to create codegen: %w", err)
							}
							if err := cg.Read(); err != nil {
								return fmt.Errorf("failed to scan namespace: %w", err)
							}
							if err := cg.Write(); err != nil {
								return fmt.Errorf("failed to generate module: %w", err)
							}
							return nil
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
						Name:    "up",
						Usage:   "ctf obs up",
						Aliases: []string{"u"},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "full",
								Aliases: []string{"f"},
								Usage:   "Spin up all the observability services",
								Value:   false,
							},
						},
						Description: "Spins up a local observability stack. Has two modes, standard (Loki, Prometheus, Grafana and OTEL) and full including also Tempo, Cadvisor and PostgreSQL metrics",
						Action: func(c *cli.Context) error {
							if c.Bool("full") {
								return framework.ObservabilityUpFull()
							}
							return framework.ObservabilityUp()
						},
					},
					{
						Name:    "down",
						Usage:   "ctf obs down",
						Aliases: []string{"d"},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "full",
								Aliases: []string{"f"},
								Usage:   "Removes all the observability services (this flag exists for compatibility, all the services are always removed with 'down')",
								Value:   false,
							},
						},
						Description: "Removes local observability stack",
						Action:      func(c *cli.Context) error { return framework.ObservabilityDown() },
					},
					{
						Name:    "restart",
						Usage:   "ctf obs r",
						Aliases: []string{"r"},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "full",
								Aliases: []string{"f"},
								Usage:   "Restart all observability services (this flag exists for compatibility, all the services are always removed with 'down')",
								Value:   false,
							},
						},
						Description: "Restart a local observability stack",
						Action: func(c *cli.Context) error {
							// always remove all the containers and volumes to clean up the data
							if err := framework.ObservabilityDown(); err != nil {
								return err
							}
							if c.Bool("full") {
								return framework.ObservabilityUpFull()
							}
							return framework.ObservabilityUp()
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
							return framework.BlockScoutUp(
								c.String("rpc"),
								c.String("chain-id"),
							)
						},
					},
					{
						Name:        "down",
						Usage:       "ctf bs down",
						Aliases:     []string{"d"},
						Description: "Removes Blockscout stack, wipes all Blockscout databases data",
						Action: func(c *cli.Context) error {
							return framework.BlockScoutDown(c.String("rpc"))
						},
					},
					{
						Name:        "reboot",
						Usage:       "ctf bs reboot or ctf bs r",
						Aliases:     []string{"r"},
						Description: "Reboots Blockscout stack",
						Action: func(c *cli.Context) error {
							rpc := c.String("rpc")
							if err := framework.BlockScoutDown(rpc); err != nil {
								return err
							}
							return framework.BlockScoutUp(
								c.String("rpc"),
								c.String("chain-id"),
							)
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
	err = os.WriteFile(outputFile, []byte(dumpData), 0o644)
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
