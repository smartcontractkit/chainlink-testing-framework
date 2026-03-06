package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
						Name:    "env",
						Aliases: []string{"e"},
						Usage:   "Generate a Chainlink Node developer environment",
						Description: `🔗 Chainlink's Developer Environment Generator 🔗

Prerequisites:
	Just will be automatically installed if not available (via Homebrew on macOS).
	For other platforms, please install it manually: https://github.com/casey/just

Usage:

	⚙️ Generate basic environment:
		ctf gen env --cli myenv --output-dir devenv --product-name Knilniahc --nodes 4

	📜 Read the docs in devenv/README.md

	🔧 Address all TODO comments and customize it
`,
						ArgsUsage: "",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "cli",
								Aliases: []string{"c"},
								Usage:   "Your devenv CLI binary name",
							},
							&cli.StringFlag{
								Name:    "output-dir",
								Aliases: []string{"o"},
								Value:   "devenv",
								Usage:   "Your devenv directory",
							},
							&cli.StringFlag{
								Name:    "product-name",
								Aliases: []string{"r"},
								Usage:   "Your product name",
							},
							&cli.IntFlag{
								Name:    "nodes",
								Aliases: []string{"n"},
								Value:   4,
								Usage:   "Chainlink Nodes",
							},
						},
						Action: func(c *cli.Context) error {
							outputDir := c.String("output-dir")
							nodes := c.Int("nodes")
							cliName := c.String("cli")
							if cliName == "" {
								return fmt.Errorf("CLI name can't be empty, choose your CLI name")
							}
							productName := c.String("product-name")
							if productName == "" {
								return fmt.Errorf("Product name must be specified, call your product somehow, any name")
							}
							framework.L.Info().
								Str("OutputDir", outputDir).
								Str("Name", cliName).
								Int("CLNodes", nodes).
								Msg("Generating developer environment")

							cg, err := framework.NewEnvBuilder(cliName, nodes, productName).
								OutputDir(outputDir).
								Build()
							if err != nil {
								return fmt.Errorf("failed to create codegen: %w", err)
							}
							if err := cg.Write(); err != nil {
								return fmt.Errorf("failed to generate environment: %w", err)
							}
							if err := cg.WriteServices(); err != nil {
								return fmt.Errorf("failed to generate services: %w", err)
							}
							if err := cg.WriteFakes(); err != nil {
								return fmt.Errorf("failed to generate fakes: %w", err)
							}
							if err := cg.WriteProducts(); err != nil {
								return fmt.Errorf("failed to generate products: %w", err)
							}

							fmt.Println()
							fmt.Printf("📁 Your environment directory is: %s\n", outputDir)
							fmt.Printf("💻 Your CLI name is: %s\n", cliName)
							fmt.Printf("📜 More docs can be found in %s/README.md\n", outputDir)
							fmt.Printf("⬛ Entering the shell..\n")
							fmt.Println()

							// Ensure 'just' is installed before proceeding
							if err := ensureJustInstalled(); err != nil {
								return fmt.Errorf("failed to ensure 'just' is installed: %w", err)
							}

							cmd := exec.Command("just", "cli")
							cmd.Env = os.Environ()
							cmd.Dir = outputDir
							out, err := cmd.CombinedOutput()
							if err != nil {
								return fmt.Errorf("failed to build CLI via Justfile: %w, output: %s", err, string(out))
							}
							if err := os.Chdir(outputDir); err != nil {
								return err
							}
							cmd = exec.Command(cliName, "sh")
							cmd.Env = os.Environ()
							cmd.Stdin = os.Stdin
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr
							err = cmd.Run()
							if err != nil {
								return fmt.Errorf("failed to enter devenv shell: %w", err)
							}
							return nil
						},
					},
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
		ctf gen load my-namespace
	With workload:
		ctf gen load -w my-namespace
	With workload and name:
		ctf gen load -w -n TestSomething my-namespace

Be aware that any TODO requires your attention before your run the final test!
`,
						ArgsUsage: "--workload --name $name --output-dir $dir --module $go_mod_name [NAMESPACE]",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Value:   "TestGeneratedLoadChaos",
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
							&cli.StringFlag{
								Name:    "latency-ms",
								Aliases: []string{"l"},
								Value:   "300",
								Usage:   "Default latency for delay experiments in milliseconds",
							},
							&cli.StringFlag{
								Name:    "jitter-ms",
								Aliases: []string{"j"},
								Value:   "100",
								Usage:   "Default jitter for delay experiments in milliseconds",
							},
						},
						Action: func(c *cli.Context) error {
							if c.Args().Len() == 0 {
								return fmt.Errorf("Kubernetes namespace argument is required")
							}
							ns := c.Args().First()
							testSuiteName := c.String("name")
							podLabelKey := c.String("pod-label-key")
							latencyMs := c.Int("latency-ms")
							jitterMs := c.Int("jitter-ms")
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
								Latency(latencyMs).
								Jitter(jitterMs).
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
				Name:    "compat",
				Aliases: []string{"c"},
				Usage:   "Performs cluster compatibility testing",
				Subcommands: []*cli.Command{
					{
						Name:    "restore",
						Aliases: []string{"r"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "base_branch",
								Usage: "Base branch on which to restore",
								Value: "develop",
							},
						},
						Usage: "Restores back to develop",
						Action: func(c *cli.Context) error {
							return framework.CheckOut(c.String("base_branch"))
						},
					},
					{
						Name:    "backward",
						Aliases: []string{"b"},
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "versions-back",
								Aliases: []string{"v"},
								Usage:   "How many versions back to test",
								Value:   1,
							},
							&cli.IntFlag{
								Name:    "nodes",
								Aliases: []string{"n"},
								Usage:   "How many nodes to upgrade",
								Value:   3,
							},
							&cli.StringFlag{
								Name:    "registry",
								Aliases: []string{"r"},
								Usage:   "Docker Image registry for Chainlink node",
								Value:   "smartcontract/chainlink",
							},
							&cli.StringFlag{
								Name:  "strip-image-suffix",
								Usage: "Stripts image suffix from ref to map it to registry images",
							},
							&cli.StringFlag{
								Name:    "buildcmd",
								Aliases: []string{"b"},
								Usage:   "Environment build command",
								Value:   "just cli",
							},
							&cli.StringFlag{
								Name:    "envcmd",
								Aliases: []string{"e"},
								Usage:   "Environment bootstrap command",
							},
							&cli.StringFlag{
								Name:    "testcmd",
								Aliases: []string{"t"},
								Usage:   "Test verification command",
							},
							&cli.StringFlag{
								Name:  "nop",
								Usage: "Find specific NOPs ref upgrade sequence",
							},
							&cli.StringFlag{
								Name:  "node-name-template",
								Usage: "CL node Docker container name template",
								Value: "don-node%d",
							},
							&cli.StringFlag{
								Name:  "sot-url",
								Usage: "RANE SOT snapshot API URL",
								Value: "https://rane-sot-app.main.prod.cldev.sh/v1/snapshot",
							},
							&cli.StringSliceFlag{
								Name:  "refs",
								Usage: "Refs to test, can be tag or commit. The corresponding image with ref should be present in registry",
							},
							&cli.StringSliceFlag{
								Name:  "include-refs",
								Usage: "Patterns to include specific refs (e.g., beta,rc,v0,v1), also used in CI to verify on test refs(tags)",
							},
							&cli.StringSliceFlag{
								Name:  "exclude-refs",
								Usage: "Patterns to exclude specific refs (e.g., beta,rc,v0,v1)",
							},
							&cli.BoolFlag{
								Name:  "skip-pull",
								Usage: "Skip docker pull; use locally built images (e.g. for local testing)",
							},
						},
						Usage: "Rollbacks N versions back, runs the test the upgrades CL nodes with new versions",
						Action: func(c *cli.Context) error {
							versionsBack := c.Int("versions-back")
							registry := c.String("registry")
							refs := c.StringSlice("refs")
							include := c.StringSlice("include-refs")
							exclude := c.StringSlice("exclude-refs")

							buildcmd := c.String("buildcmd")
							envcmd := c.String("envcmd")
							testcmd := c.String("testcmd")
							nodes := c.Int("nodes")
							nodeNameTemplate := c.String("node-name-template")

							nop := c.String("nop")
							sotURL := c.String("sot-url")
							skipPull := c.Bool("skip-pull")

							// test logic is:
							// - rollback to selected ref
							// - spin up the env and perform the initial smoke test
							// - upgrade N CL nodes with preversing DB volume (shared database)
							// - perform the test again
							// - repeat until all the new versions are validated

							// if no refs provided find refs (tags) sequence SemVer sequence for last N versions_back
							// else, use refs param slice
							//
							var err error
							if len(refs) == 0 && nop == "" {
								refs, err = framework.FindSemVerRefSequence(versionsBack, include, exclude)
								if err != nil {
									return err
								}
							}
							if nop != "" {
								refs, err = framework.FindNOPRefs(sotURL, nop, exclude)
								if err != nil {
									return err
								}
							}
							framework.L.Info().
								Int("TotalSemVerRefs", len(refs)).
								Strs("SelectedRefs", refs).
								Str("EarliestRefs", refs[0]).
								Msg("Formed upgrade sequence")
							// if no commands just show the tags and return
							if buildcmd == "" || envcmd == "" || testcmd == "" {
								framework.L.Info().Msg("No envcmd or testcmd or buildcmd provided, skipping")
								return nil
							}
							// checkout the oldest ref
							if err := framework.CheckOut(refs[0]); err != nil {
								return err
							}

							// this is a hack allowing us to match what we have in Git to registry or NOP version
							// it'd exist until we stabilize tagging strategy
							testImgSuffix := c.String("strip-image-suffix")
							for i := range refs {
								refs[i] = strings.ReplaceAll(refs[i], testImgSuffix, "")
							}

							// setup the env and verify with test command
							framework.L.Info().Strs("Sequence", refs).Msg("Running upgrade sequence")
							// first env var is used by devenv, second by simple nodeset to override the image
							os.Setenv("CHAINLINK_IMAGE", fmt.Sprintf("%s:%s", registry, refs[0]))
							os.Setenv("CTF_CHAINLINK_IMAGE", fmt.Sprintf("%s:%s", registry, refs[0]))
							if _, err := framework.ExecCmdWithContext(c.Context, buildcmd); err != nil {
								return err
							}
							if _, err := framework.ExecCmdWithContext(c.Context, envcmd); err != nil {
								return err
							}
							if _, err := framework.ExecCmdWithContext(c.Context, testcmd); err != nil {
								return err
							}
							// start upgrading nodes and verifying with tests, all the versions except the oldest one
							for _, tag := range refs[1:] {
								framework.L.Info().
									Int("Nodes", nodes).
									Str("Version", tag).
									Msg("Upgrading nodes")
								img := fmt.Sprintf("%s:%s", registry, tag)
								if !skipPull {
									if _, err := framework.ExecCmdWithContext(c.Context, fmt.Sprintf("docker pull %s", img)); err != nil {
										return fmt.Errorf("failed to pull image %s: %w", img, err)
									}
								}
								for i := range nodes {
									if err := framework.UpgradeContainer(c.Context, fmt.Sprintf(nodeNameTemplate, i), img); err != nil {
										return err
									}
								}
								if _, err := framework.ExecCmd(testcmd); err != nil {
									return err
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "config",
				Usage: "Shapes your test config, removes outputs, formatting ,etc",
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

// ensureJustInstalled checks if 'just' is available in PATH, and if not, attempts to install it.
// On macOS, it tries to install via Homebrew. On other platforms, it provides installation instructions.
func ensureJustInstalled() error {
	// Check if just is already available
	if _, err := exec.LookPath("just"); err == nil {
		return nil
	}

	fmt.Println("⚠️  'just' command not found in PATH")
	fmt.Println("📦 Attempting to install 'just'...")

	// Try to install via Homebrew on macOS
	if runtime.GOOS == "darwin" {
		// Check if Homebrew is available
		if _, err := exec.LookPath("brew"); err == nil {
			fmt.Println("🍺 Installing 'just' via Homebrew...")
			cmd := exec.Command("brew", "install", "just")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install 'just' via Homebrew: %w. Please install manually: brew install just", err)
			}
			fmt.Println("✅ Successfully installed 'just'")
			return nil
		}
		// Homebrew not available, provide instructions
		return fmt.Errorf("'just' is not installed and Homebrew is not available. Please install 'just' manually:\n  brew install just\n  Or visit: https://github.com/casey/just")
	}

	// For non-macOS platforms, provide installation instructions
	return fmt.Errorf("'just' is not installed. Please install it manually:\n  Visit: https://github.com/casey/just\n  Or use your package manager (e.g., apt install just, pacman -S just)")
}
