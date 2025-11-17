package framework

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
)

const (
	// DefaultTestSuiteName default table test name
	DefaultTestSuiteName = "TestSmokeE2E"
)

/* Templates */

const (
	// GoModTemplate go module template
	GoModTemplate = `module {{.ModuleName}}

go {{.RuntimeVersion}}

require (
	github.com/smartcontractkit/chainlink-evm v0.0.0-20250709215002-07f34ab867df
	github.com/smartcontractkit/chainlink-deployments-framework v0.17.0
	github.com/smartcontractkit/chainlink-testing-framework/framework v0.11.10
)

replace github.com/fbsobreira/gotron-sdk => github.com/smartcontractkit/chainlink-tron/relayer/gotron-sdk v0.0.5-0.20250528121202-292529af39df

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/prometheus/common => github.com/prometheus/common v0.62.0
	github.com/smartcontractkit/chainlink-testing-framework/lib => github.com/smartcontractkit/chainlink-testing-framework/lib v1.54.4
)
`
	// ConfigTOMLTmpl is a default env.toml template for devenv describind components configuration
	ConfigTOMLTmpl = `[on_chain]
  link_contract_address = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
  cl_nodes_funding_eth = 50
  cl_nodes_funding_link = 50
  verification_timeout_sec = 400
  contracts_configuration_timeout_sec = 60
  verify = false

  [on_chain.gas_settings]
  fee_cap_multiplier = 2
  tip_cap_multiplier = 2


[[blockchains]]
  chain_id = "1337"
  docker_cmd_params = ["-b", "1", "--mixed-mining", "--slots-in-an-epoch", "1"]
  image = "ghcr.io/foundry-rs/foundry:stable"
  port = "8545"
  type = "anvil"

[[nodesets]]
  name = "don"
  nodes = {{ .Nodes }}
  override_mode = "each"

  [nodesets.db]
    image = "postgres:15.0"

	{{- range .NodeIndices }}
	[[nodesets.node_specs]]
	    [nodesets.node_specs.node]
	    image = "public.ecr.aws/chainlink/chainlink:2.26.0"
	{{- end }}
`
	// CompletionTmpl is a go-prompt library completion template providing interactive prompt
	CompletionTmpl = `package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

func getCommands() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "", Description: "Choose command, press <space> for more options after selecting command"},
		{Text: "up", Description: "Spin up the development environment"},
		{Text: "down", Description: "Tear down the development environment"},
		{Text: "restart", Description: "Restart the development environment"},
		{Text: "test", Description: "Perform smoke or load/chaos testing"},
		{Text: "bs", Description: "Manage the Blockscout EVM block explorer"},
		{Text: "obs", Description: "Manage the observability stack"},
		{Text: "db", Description: "Inspect Databases"},
		{Text: "exit", Description: "Exit the interactive shell"},
	}
}

func getSubCommands(parent string) []prompt.Suggest {
	switch parent {
	case "test":
		return []prompt.Suggest{
			{Text: "soak", Description: "Run {{ .CLIName }} soak test"},
			{Text: "rpc-latency", Description: "Default soak test + 400ms RPC latency (all chains)"},
			{Text: "gas-spikes", Description: "Default soak test + slow and fast gas spikes"},
			{Text: "reorgs", Description: "Default soak test + reorgs (Requires 'up env.toml,env-geth.toml' environment"},
			{Text: "chaos", Description: "Default soak test + chaos (restarts, latency, data loss between services)"},
		}
	case "bs":
		return []prompt.Suggest{
			{Text: "up", Description: "Spin up Blockscout and listen to dst chain (8555)"},
			{Text: "up -u http://host.docker.internal:8545 -c 1337", Description: "Spin up Blockscout and listen to src chain (8545)"},
			{Text: "down", Description: "Remove Blockscout stack"},
			{Text: "restart", Description: "Restart Blockscout and listen to dst chain (8555)"},
			{Text: "restart -u http://host.docker.internal:8545 -c 1337", Description: "Restart Blockscout and listen to src chain (8545)"},
		}
	case "obs":
		return []prompt.Suggest{
			{Text: "up", Description: "Spin up observability stack (Loki/Prometheus/Grafana)"},
			{Text: "up -f", Description: "Spin up full observability stack (Pyroscope, cadvisor, postgres exporter)"},
			{Text: "down", Description: "Spin down observability stack"},
			{Text: "restart", Description: "Restart observability stack"},
			{Text: "restart -f", Description: "Restart full observability stack"},
		}
	case "u":
		fallthrough
	case "up":
		fallthrough
	case "r":
		fallthrough
	case "restart":
		return []prompt.Suggest{
			{Text: "env.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes"},
			{Text: "env.toml,env-cl-rebuild.toml", Description: "Spin up Anvil <> Anvil local chains, all services, 4 CL nodes (custom build)"},
			{Text: "env.toml,env-geth.toml", Description: "Spin up Geth <> Geth local chains (clique), all services, 4 CL nodes"},
			{Text: "env.toml,env-fuji-fantom.toml", Description: "Spin up testnets: Fuji <> Fantom, all services, 4 CL nodes"},
		}
	default:
		return []prompt.Suggest{}
	}
}

func executor(in string) {
	checkDockerIsRunning()
	in = strings.TrimSpace(in)
	if in == "" {
		return
	}
	if in == "exit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	args := strings.Fields(in)
	os.Args = append([]string{"{{ .CLIName }}"}, args...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// completer provides autocomplete suggestions for multi-word commands.
func completer(in prompt.Document) []prompt.Suggest {
	text := in.TextBeforeCursor()
	words := strings.Fields(text)
	lastCharIsSpace := len(text) > 0 && text[len(text)-1] == ' '

	switch {
	case len(words) == 0:
		return getCommands()
	case len(words) == 1:
		if lastCharIsSpace {
			return getSubCommands(words[0])
		} else {
			return prompt.FilterHasPrefix(getCommands(), words[0], true)
		}

	case len(words) >= 2:
		if lastCharIsSpace {
			return []prompt.Suggest{}
		} else {
			parent := words[0]
			currentWord := words[len(words)-1]
			return prompt.FilterHasPrefix(getSubCommands(parent), currentWord, true)
		}
	default:
		return []prompt.Suggest{}
	}
}

// resetTerm resets terminal settings to Unix defaults.
func resetTerm() {
	cmd := exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func StartShell() {
	defer resetTerm()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("{{ .CLIName }}> "),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionTitle("CCIP Interactive Shell"),
		prompt.OptionMaxSuggestion(15),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionCompletionWordSeparator(" "),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.Green),
		prompt.OptionScrollbarThumbColor(prompt.DarkGray),
		prompt.OptionScrollbarBGColor(prompt.Black),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(buf *prompt.Buffer) {
				fmt.Println("Interrupted, exiting...")
				resetTerm()
				os.Exit(0)
			},
		}),
	)
	p.Run()
}
`
	// CLITmpl is a Cobra library CLI template with basic devenv commands
	CLITmpl = `package main

import (
		"context"
		"fmt"
		"os"
		"os/exec"
		"syscall"

		"github.com/docker/docker/client"
		"github.com/spf13/cobra"

		"github.com/smartcontractkit/chainlink-testing-framework/framework"
		"{{ .DevEnvPkgImport }}"
)

const (
		LocalWASPLoadDashboard = "http://localhost:3000/d/WASPLoadTests/wasp-load-test?orgId=1&from=now-5m&to=now&refresh=5s"
		Local{{ .CLIName}}Dashboard      = "http://localhost:3000/d/f8a04cef-653f-46d3-86df-87c532300672/svr-soak-test?orgId=1&refresh=5s"
)

var rootCmd = &cobra.Command{
	Use:   "{{ .CLIName }}",
	Short: "A {{ .CLIName }} local environment tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		if debug {
			framework.L.Info().Msg("Debug mode enabled, setting CTF_CLNODE_DLV=true")
			os.Setenv("CTF_CLNODE_DLV", "true")
		}
		return nil
	},
}

var restartCmd = &cobra.Command{
		Use:     "restart",
		Aliases: []string{"r"},
		Args:    cobra.RangeArgs(0, 1),
		Short:   "Restart development environment, remove apps and apply default configuration again",
		RunE: func(cmd *cobra.Command, args []string) error {
			var configFile string
			if len(args) > 0 {
				configFile = args[0]
			} else {
				configFile = "env.toml"
			}
			framework.L.Info().Str("Config", configFile).Msg("Reconfiguring development environment")
			_ = os.Setenv("CTF_CONFIGS", configFile)
			_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			framework.L.Info().Msg("Tearing down the development environment")
			err := framework.RemoveTestContainers()
			if err != nil {
				return fmt.Errorf("failed to clean Docker resources: %w", err)
			}
			_, err = devenv.NewEnvironment()
			return err
		},
}

var upCmd = &cobra.Command{
		Use:     "up",
		Aliases: []string{"u"},
		Short:   "Spin up the development environment",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var configFile string
			if len(args) > 0 {
				configFile = args[0]
			} else {
				configFile = "env.toml"
			}
			framework.L.Info().Str("Config", configFile).Msg("Creating development environment")
			_ = os.Setenv("CTF_CONFIGS", configFile)
			_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
			_, err := devenv.NewEnvironment()
			if err != nil {
				return err
			}
			return nil
		},
}

var downCmd = &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Tear down the development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			framework.L.Info().Msg("Tearing down the development environment")
			err := framework.RemoveTestContainers()
			if err != nil {
				return fmt.Errorf("failed to clean Docker resources: %w", err)
			}
			return nil
		},
}

var bsCmd = &cobra.Command{
		Use:   "bs",
		Short: "Manage the Blockscout EVM block explorer",
		Long:  "Spin up or down the Blockscout EVM block explorer",
}

var bsUpCmd = &cobra.Command{
		Use:     "up",
		Aliases: []string{"u"},
		Short:   "Spin up Blockscout EVM block explorer",
		RunE: func(cmd *cobra.Command, args []string) error {
			url, _ := bsCmd.Flags().GetString("url")
			chainID, _ := bsCmd.Flags().GetString("chain-id")
			return framework.BlockScoutUp(url, chainID)
		},
}

var bsDownCmd = &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Spin down Blockscout EVM block explorer",
		RunE: func(cmd *cobra.Command, args []string) error {
			url, _ := bsCmd.Flags().GetString("url")
			return framework.BlockScoutDown(url)
		},
}

var bsRestartCmd = &cobra.Command{
		Use:     "restart",
		Aliases: []string{"r"},
		Short:   "Restart the Blockscout EVM block explorer",
		RunE: func(cmd *cobra.Command, args []string) error {
			url, _ := bsCmd.Flags().GetString("url")
			chainID, _ := bsCmd.Flags().GetString("chain-id")
			if err := framework.BlockScoutDown(url); err != nil {
				return err
			}
			return framework.BlockScoutUp(url, chainID)
		},
}

var obsCmd = &cobra.Command{
		Use:   "obs",
		Short: "Manage the observability stack",
		Long:  "Spin up or down the observability stack with subcommands 'up' and 'down'",
}

var obsUpCmd = &cobra.Command{
		Use:     "up",
		Aliases: []string{"u"},
		Short:   "Spin up the observability stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			full, _ := cmd.Flags().GetBool("full")
			var err error
			if full {
				err = framework.ObservabilityUpFull()
			} else {
				err = framework.ObservabilityUp()
			}
			if err != nil {
				return fmt.Errorf("observability up failed: %w", err)
			}
			devenv.Plog.Info().Msgf("{{ .CLIName }} Dashboard: %s", Local{{ .CLIName }}Dashboard)
			devenv.Plog.Info().Msgf("{{ .CLIName }} Load Test Dashboard: %s", LocalWASPLoadDashboard)
			return nil
		},
}

var obsDownCmd = &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Spin down the observability stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return framework.ObservabilityDown()
		},
}

var obsRestartCmd = &cobra.Command{
		Use:     "restart",
		Aliases: []string{"r"},
		Short:   "Restart the observability stack (data wipe)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := framework.ObservabilityDown(); err != nil {
				return fmt.Errorf("observability down failed: %w", err)
			}
			full, _ := cmd.Flags().GetBool("full")
			var err error
			if full {
				err = framework.ObservabilityUpFull()
			} else {
				err = framework.ObservabilityUp()
			}
			if err != nil {
				return fmt.Errorf("observability up failed: %w", err)
			}
			devenv.Plog.Info().Msgf("{{ .CLIName }} Dashboard: %s", Local{{ .CLIName }}Dashboard)
			devenv.Plog.Info().Msgf("{{ .CLIName }} Load Test Dashboard: %s", LocalWASPLoadDashboard)
			return nil
		},
}

var testCmd = &cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Run the tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("specify the test suite: smoke or load")
			}
			var testPattern string
			switch args[0] {
			case "soak":
				testPattern = "TestSoak"
			case "rpc-latency":
				testPattern = "TestSoak/rpc_latency"
			case "gas-spikes":
				testPattern = "TestSoak/gas"
			case "reorg":
				testPattern = "TestSoak/reorg"
			case "chaos":
				testPattern = "TestSoak/chaos"
			default:
				return fmt.Errorf("test suite %s is unknown, choose between smoke or load", args[0])
			}

			testCmd := exec.Command("go", "test", "-v", "-run", testPattern)
			testCmd.Dir = "./tests"
			testCmd.Stdout = os.Stdout
			testCmd.Stderr = os.Stderr
			testCmd.Stdin = os.Stdin

			if err := testCmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
						os.Exit(status.ExitStatus())
					}
					os.Exit(1)
				}
				return fmt.Errorf("failed to run test command: %w", err)
			}
			return nil
		},
}

func init() {
		rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable running services with dlv to allow remote debugging.")

		rootCmd.AddCommand(testCmd)

		// Blockscout, on-chain debug
		bsCmd.PersistentFlags().StringP("url", "u", "http://host.docker.internal:8555", "EVM RPC node URL (default to dst chain on 8555")
		bsCmd.PersistentFlags().StringP("chain-id", "c", "2337", "RPC's Chain ID")
		bsCmd.AddCommand(bsUpCmd)
		bsCmd.AddCommand(bsDownCmd)
		bsCmd.AddCommand(bsRestartCmd)
		rootCmd.AddCommand(bsCmd)

		// observability
		obsCmd.PersistentFlags().BoolP("full", "f", false, "Enable full observability stack with additional components")
		obsCmd.AddCommand(obsRestartCmd)
		obsCmd.AddCommand(obsUpCmd)
		obsCmd.AddCommand(obsDownCmd)
		rootCmd.AddCommand(obsCmd)

		// main env commands
		rootCmd.AddCommand(upCmd)
		rootCmd.AddCommand(restartCmd)
		rootCmd.AddCommand(downCmd)
}

func checkDockerIsRunning() {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Println("Can't create Docker client, please check if Docker daemon is running!")
			os.Exit(1)
		}
		defer cli.Close()
		_, err = cli.Ping(context.Background())
		if err != nil {
			fmt.Println("Docker is not running, please start Docker daemon first!")
			os.Exit(1)
		}
}

func main() {
		checkDockerIsRunning()
		if len(os.Args) == 2 && (os.Args[1] == "shell" || os.Args[1] == "sh") {
			_ = os.Setenv("CTF_CONFIGS", "env.toml") // Set default config for shell

			StartShell()
			return
		}
		if err := rootCmd.Execute(); err != nil {
			devenv.Plog.Err(err).Send()
			os.Exit(1)
		}
}`
	// TestsTmpl is a an e2e table test template
	TestsTmpl = ``
	// CLDFTmpl is a Chainlink Deployments Framework template
	CLDFTmpl = `package {{ .PackageName }}

import (
		"context"
		"crypto/ecdsa"
		"errors"
		"fmt"
		"math/big"
		"strings"
		"time"

		"github.com/Masterminds/semver/v3"
		"github.com/ethereum/go-ethereum/accounts/abi/bind"
		"github.com/ethereum/go-ethereum/common"
		"github.com/ethereum/go-ethereum/core/types"
		"github.com/ethereum/go-ethereum/crypto"
		"github.com/ethereum/go-ethereum/ethclient"
		"github.com/smartcontractkit/chainlink-common/pkg/logger"
		"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
		"github.com/smartcontractkit/chainlink-deployments-framework/operations"
		"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
		"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
		"go.uber.org/zap"
		"golang.org/x/sync/errgroup"
		"google.golang.org/grpc"
		"google.golang.org/grpc/credentials/insecure"

		chainsel "github.com/smartcontractkit/chain-selectors"
		cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
		cldfevm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
		cldfevmprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
		cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
		csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
		jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
		nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

const (
		AnvilKey0                     = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		DefaultNativeTransferGasPrice = 21000
)

const LinkToken cldf.ContractType = "LinkToken"

var _ cldf.ChangeSet[[]uint64] = DeployLinkToken

type JobDistributor struct {
	nodev1.NodeServiceClient
	jobv1.JobServiceClient
	csav1.CSAServiceClient
	WSRPC string
}

type JDConfig struct {
	GRPC  string
	WSRPC string
}

// DeployLinkToken deploys a link token contract to the chain identified by the ChainSelector.
func DeployLinkToken(e cldf.Environment, chains []uint64) (cldf.ChangesetOutput, error) { //nolint:gocritic
	newAddresses := cldf.NewMemoryAddressBook()
	deployGrp := errgroup.Group{}
	for _, chain := range chains {
		family, err := chainsel.GetSelectorFamily(chain)
		if err != nil {
			return cldf.ChangesetOutput{AddressBook: newAddresses}, err
		}
		var deployFn func() error
		switch family {
		case chainsel.FamilyEVM:
			// Deploy EVM LINK token
			deployFn = func() error {
				_, err := deployLinkTokenContractEVM(
					e.Logger, e.BlockChains.EVMChains()[chain], newAddresses,
				)
				return err
			}
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("unsupported chain family %s", family)
		}
		deployGrp.Go(func() error {
			err := deployFn()
			if err != nil {
				e.Logger.Errorw("Failed to deploy link token", "chain", chain, "err", err)
				return fmt.Errorf("failed to deploy link token for chain %d: %w", chain, err)
			}
			return nil
		})
	}
	return cldf.ChangesetOutput{AddressBook: newAddresses}, deployGrp.Wait()
}

func deployLinkTokenContractEVM(
		lggr logger.Logger,
		chain cldfevm.Chain, //nolint:gocritic
		ab cldf.AddressBook,
) (*cldf.ContractDeploy[*link_token.LinkToken], error) {
	linkToken, err := cldf.DeployContract[*link_token.LinkToken](lggr, chain, ab,
		func(chain cldfevm.Chain) cldf.ContractDeploy[*link_token.LinkToken] {
			var (
				linkTokenAddr common.Address
				tx            *types.Transaction
				linkToken     *link_token.LinkToken
				err2          error
			)
			if !chain.IsZkSyncVM {
				linkTokenAddr, tx, linkToken, err2 = link_token.DeployLinkToken(
					chain.DeployerKey,
					chain.Client,
				)
			} else {
				linkTokenAddr, _, linkToken, err2 = link_token.DeployLinkTokenZk(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
				)
			}
			return cldf.ContractDeploy[*link_token.LinkToken]{
				Address:  linkTokenAddr,
				Contract: linkToken,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(LinkToken, *semver.MustParse("1.0.0")),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy link token", "chain", chain.String(), "err", err)
		return linkToken, err
	}
	return linkToken, nil
}

// LoadCLDFEnvironment loads CLDF environment with a memory data store and JD client.
func LoadCLDFEnvironment(in *Cfg) (cldf.Environment, error) {
	ctx := context.Background()

	getCtx := func() context.Context {
		return ctx
	}

	// This only generates a brand new datastore and does not load any existing data.
	// We will need to figure out how data will be persisted and loaded in the future.
	ds := datastore.NewMemoryDataStore().Seal()

	lggr, err := logger.NewWith(func(config *zap.Config) {
		config.Development = true
		config.Encoding = "console"
	})
	if err != nil {
		return cldf.Environment{}, fmt.Errorf("failed to create logger: %w", err)
	}

	blockchains, err := loadCLDFChains(in.Blockchains)
	if err != nil {
		return cldf.Environment{}, fmt.Errorf("failed to load CLDF chains: %w", err)
	}

	jd, err := NewJDClient(ctx, JDConfig{
		GRPC:  in.JD.Out.ExternalGRPCUrl,
		WSRPC: in.JD.Out.ExternalWSRPCUrl,
	})
	if err != nil {
		return cldf.Environment{},
			fmt.Errorf("failed to load offchain client: %w", err)
	}

	opBundle := operations.NewBundle(
		getCtx,
		lggr,
		operations.NewMemoryReporter(),
		operations.WithOperationRegistry(operations.NewOperationRegistry()),
	)

	return cldf.Environment{
		Name:              "local",
		Logger:            lggr,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         ds,
		Offchain:          jd,
		GetContext:        getCtx,
		OperationsBundle:  opBundle,
		BlockChains:       cldfchain.NewBlockChainsFromSlice(blockchains),
	}, nil
}

func loadCLDFChains(bcis []*blockchain.Input) ([]cldfchain.BlockChain, error) {
	blockchains := make([]cldfchain.BlockChain, 0)
	for _, bci := range bcis {
		switch bci.Type {
		case "anvil":
			bc, err := loadEVMChain(bci)
			if err != nil {
				return blockchains, fmt.Errorf("failed to load EVM chain %s: %w", bci.ChainID, err)
			}

			blockchains = append(blockchains, bc)
		default:
			return blockchains, fmt.Errorf("unsupported chain type %s", bci.Type)
		}
	}

	return blockchains, nil
}

func loadEVMChain(bci *blockchain.Input) (cldfchain.BlockChain, error) {
	if bci.Out == nil {
		return nil, fmt.Errorf("output configuration for %s blockchain %s is not set", bci.Type, bci.ChainID)
	}

	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(bci.ChainID, chainsel.FamilyEVM)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain details for %s: %w", bci.ChainID, err)
	}

	chain, err := cldfevmprovider.NewRPCChainProvider(
		chainDetails.ChainSelector,
		cldfevmprovider.RPCChainProviderConfig{
			DeployerTransactorGen: cldfevmprovider.TransactorFromRaw(
				// TODO: we need to figure out a reliable way to get secrets here that is
				// TODO: - easy for developers
				// TODO: - works the same way locally, in K8s and in CI
				// TODO: - do not require specific AWS access like AWSSecretsManager
				// TODO: for now it's just an Anvil 0 key
				AnvilKey0,
			),
			RPCs: []cldf.RPC{
				{
					Name:               "default",
					WSURL:              bci.Out.Nodes[0].ExternalWSUrl,
					HTTPURL:            bci.Out.Nodes[0].ExternalHTTPUrl,
					PreferredURLScheme: cldf.URLSchemePreferenceHTTP,
				},
			},
			ConfirmFunctor: cldfevmprovider.ConfirmFuncGeth(1 * time.Minute),
		},
	).Initialize(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize EVM chain %s: %w", bci.ChainID, err)
	}

	return chain, nil
}

// NewJDClient creates a new JobDistributor client.
func NewJDClient(ctx context.Context, cfg JDConfig) (cldf.OffchainClient, error) {
	conn, err := NewJDConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}
	jd := &JobDistributor{
		WSRPC:             cfg.WSRPC,
		NodeServiceClient: nodev1.NewNodeServiceClient(conn),
		JobServiceClient:  jobv1.NewJobServiceClient(conn),
		CSAServiceClient:  csav1.NewCSAServiceClient(conn),
	}

	return jd, err
}

func (jd JobDistributor) GetCSAPublicKey(ctx context.Context) (string, error) {
	keypairs, err := jd.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
	if err != nil {
		return "", err
	}
	if keypairs == nil || len(keypairs.Keypairs) == 0 {
		return "", errors.New("no keypairs found")
	}
	csakey := keypairs.Keypairs[0].PublicKey
	return csakey, nil
}

// ProposeJob proposes jobs through the jobService and accepts the proposed job on selected node based on ProposeJobRequest.NodeId.
func (jd JobDistributor) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	res, err := jd.JobServiceClient.ProposeJob(ctx, in, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to propose job. err: %w", err)
	}
	if res.Proposal == nil {
		return nil, errors.New("failed to propose job. err: proposal is nil")
	}

	return res, nil
}

// NewJDConnection creates new gRPC connection with JobDistributor.
func NewJDConnection(cfg JDConfig) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	interceptors := []grpc.UnaryClientInterceptor{}

	if len(interceptors) > 0 {
		opts = append(opts, grpc.WithChainUnaryInterceptor(interceptors...))
	}

	conn, err := grpc.NewClient(cfg.GRPC, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}

	return conn, nil
}

// FundNodeEIP1559 funds CL node using RPC URL, recipient address and amount of funds to send (ETH).
// Uses EIP-1559 transaction type.
func FundNodeEIP1559(c *ethclient.Client, pkey, recipientAddress string, amountOfFundsInETH float64) error {
	amount := new(big.Float).Mul(big.NewFloat(amountOfFundsInETH), big.NewFloat(1e18))
	amountWei, _ := amount.Int(nil)
	Plog.Info().Str("Addr", recipientAddress).Str("Wei", amountWei.String()).Msg("Funding Node")

	chainID, err := c.NetworkID(context.Background())
	if err != nil {
		return err
	}
	privateKeyStr := strings.TrimPrefix(pkey, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	feeCap, err := c.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	tipCap, err := c.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}
	recipient := common.HexToAddress(recipientAddress)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &recipient,
		Value:     amountWei,
		Gas:       DefaultNativeTransferGasPrice,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
	})
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return err
	}
	err = c.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	if _, err := bind.WaitMined(context.Background(), c, signedTx); err != nil {
		return err
	}
	Plog.Info().Str("Wei", amountWei.String()).Msg("Funded with ETH")
	return nil
}

/*
This is just a basic ETH client, CLDF should provide something like this
*/

// ETHClient creates a basic Ethereum client using PRIVATE_KEY env var and tip/cap gas settings
func ETHClient(in *Cfg) (*ethclient.Client, *bind.TransactOpts, string, error) {
	rpcURL := in.Blockchains[0].Out.Nodes[0].ExternalWSUrl
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not connect to eth client: %w", err)
	}
	privateKey, err := crypto.HexToECDSA(GetNetworkPrivateKey())
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not parse private key: %w", err)
	}
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey).String()
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not get chain ID: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not create transactor: %w", err)
	}
	gasSettings := in.OnChain.GasSettings
	fc, tc, err := MultiplyEIP1559GasPrices(client, gasSettings.FeeCapMultiplier, gasSettings.TipCapMultiplier)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not get bumped gas price: %w", err)
	}
	auth.GasFeeCap = fc
	auth.GasTipCap = tc
	Plog.Info().
		Str("GasFeeCap", fc.String()).
		Str("GasTipCap", tc.String()).
		Msg("Default gas prices set")
	return client, auth, address, nil
}

// MultiplyEIP1559GasPrices returns bumped EIP1159 gas prices increased by multiplier
func MultiplyEIP1559GasPrices(client *ethclient.Client, fcMult, tcMult int64) (*big.Int, *big.Int, error) {
	feeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, err
	}
	tipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return new(big.Int).Mul(feeCap, big.NewInt(fcMult)), new(big.Int).Mul(tipCap, big.NewInt(tcMult)), nil
}
`
	// DebugToolsTmpl is a template for various debug tools, tracing, tx debug, etc
	DebugToolsTmpl = `package {{ .PackageName }}

import (
		"os"
		"runtime/trace"
)


func tracing() func() {
		f, err := os.Create("trace.out")
		if err != nil {
			panic("can't create trace.out file")
		}
		if err := trace.Start(f); err != nil {
			panic("can't start tracing")
		}
		return func() {
			trace.Stop()
			f.Close()
		}
}
`
	// ConfigTmpl is a template for reading and writing devenv configuration (env.toml, env-out.toml)
	ConfigTmpl = `package {{ .PackageName }}

/*
This file provides a simple boilerplate for TOML configuration with overrides
It has 4 functions: Load[T], Store[T], LoadCache[T] and GetNetworkPrivateKey

To configure the environment we use a set of files we read from the env var CTF_CONFIGS=env.toml,overrides.toml (can be more than 2) in Load[T]
To store infra or product component outputs we use Store[T] that creates env-cache.toml file.
This file can be used in tests or in any other code that integrated with dev environment.
LoadCache[T] is used if you need to write outputs the second time.

GetNetworkPrivateKey is used to get your network private key from the env var we are using across all our environments, or fallback to default Anvil's key.
*/

import (
		"errors"
		"fmt"
		"os"
		"path/filepath"
		"strings"

		"github.com/davecgh/go-spew/spew"
		"github.com/pelletier/go-toml/v2"
		"github.com/rs/zerolog"
		"github.com/rs/zerolog/log"
)

const (
		// DefaultConfigDir is the default directory we are expecting TOML config to be.
		DefaultConfigDir = "."
		// EnvVarTestConfigs is the environment variable name to read config paths from, ex.: CTF_CONFIGS=env.toml,overrides.toml.
		EnvVarTestConfigs = "CTF_CONFIGS"
		// DefaultOverridesFilePath is the default overrides.toml file path.
		DefaultOverridesFilePath = "overrides.toml"
		// DefaultAnvilKey is a default, well-known Anvil first key
		DefaultAnvilKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel)

// Load loads TOML configurations from environment variable, ex.: CTF_CONFIGS=env.toml,overrides.toml
// and unmarshalls the files from left to right overriding keys.
func Load[T any]() (*T, error) {
		var config T
		paths := strings.Split(os.Getenv(EnvVarTestConfigs), ",")
		for _, path := range paths {
			L.Info().Str("Path", path).Msg("Loading configuration input")
			data, err := os.ReadFile(filepath.Join(DefaultConfigDir, path)) //nolint:gosec
			if err != nil {
				if path == DefaultOverridesFilePath {
					L.Info().Str("Path", path).Msg("Overrides file not found or empty")
					continue
				}
				return nil, fmt.Errorf("error reading config file %s: %w", path, err)
			}
			if L.GetLevel() == zerolog.TraceLevel {
				fmt.Println(string(data)) //nolint:forbidigo
			}

			decoder := toml.NewDecoder(strings.NewReader(string(data)))
			decoder.DisallowUnknownFields()

			if err := decoder.Decode(&config); err != nil {
				var details *toml.StrictMissingError
				if errors.As(err, &details) {
					fmt.Println(details.String()) //nolint:forbidigo
				}
				return nil, fmt.Errorf("failed to decode TOML config, strict mode: %s", err)
			}
		}
		if L.GetLevel() == zerolog.TraceLevel {
			L.Trace().Msg("Merged inputs")
			spew.Dump(config) //nolint:forbidigo
		}
		return &config, nil
}

// Store writes config to a file, adds -cache.toml suffix if it's an initial configuration.
func Store[T any](cfg *T) error {
		baseConfigPath, err := BaseConfigPath()
		if err != nil {
			return err
		}
		newCacheName := strings.ReplaceAll(baseConfigPath, ".toml", "")
		var outCacheName string
		if strings.Contains(newCacheName, "cache") {
			L.Info().Str("Cache", baseConfigPath).Msg("Cache file already exists, overriding")
			outCacheName = baseConfigPath
		} else {
			outCacheName = fmt.Sprintf("%s-out.toml", strings.ReplaceAll(baseConfigPath, ".toml", ""))
		}
		L.Info().Str("OutputFile", outCacheName).Msg("Storing configuration output")
		d, err := toml.Marshal(cfg)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(DefaultConfigDir, outCacheName), d, 0o600)
}

// LoadOutput loads config output file from path.
func LoadOutput[T any](path string) (*T, error) {
		_ = os.Setenv(EnvVarTestConfigs, path)
		return Load[T]()
}

// BaseConfigPath returns base config path, ex. env.toml,overrides.toml -> env.toml.
func BaseConfigPath() (string, error) {
		configs := os.Getenv(EnvVarTestConfigs)
		if configs == "" {
			return "", fmt.Errorf("no %s env var is provided, you should provide at least one test config in TOML", EnvVarTestConfigs)
		}
		L.Debug().Str("Configs", configs).Msg("Getting base config path")
		return strings.Split(configs, ",")[0], nil
}

// GetNetworkPrivateKey gets network private key or fallback to default simulator key (Anvil's first key)
func GetNetworkPrivateKey() string {
		pk := os.Getenv("PRIVATE_KEY")
		if pk == "" {
			// that's the first Anvil and Geth private key, serves as a fallback for local testing if not overridden
			return DefaultAnvilKey
		}
		return pk
}
`
	// EnvironmentTmpl is an environment.go template - main file for environment composition
	EnvironmentTmpl = `package {{ .PackageName }}

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type Cfg struct {
    OnChain         *OnChain                ` + "`" + `toml:"on_chain"` + "`" + `
    Blockchains []*blockchain.Input ` + "`" + `toml:"blockchains" validate:"required"` + "`" + `
    NodeSets    []*ns.Input         ` + "`" + `toml:"nodesets"    validate:"required"` + "`" + `
    JD          *jd.Input           ` + "`" + `toml:"jd"` + "`" + `
}

func NewEnvironment() (*Cfg, error) {
		endTracing := tracing()
		defer endTracing()

		if err := framework.DefaultNetwork(nil); err != nil {
			return nil, err
		}
		in, err := Load[Cfg]()
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration: %w", err)
		}
		_, err = blockchain.NewBlockchainNetwork(in.Blockchains[0])
		if err != nil {
			return nil, fmt.Errorf("failed to create blockchain network 1337: %w", err)
		}
		if err := DefaultProductConfiguration(in, ConfigureNodesNetwork); err != nil {
			return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
		}
		_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create new shared db node set: %w", err)
		}
		if err := DefaultProductConfiguration(in, ConfigureProductContractsJobs); err != nil {
			return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
		}
		return in, Store[Cfg](in)
}
`
	// SingleNetworkProductConfigurationTmpl is an single-network product configuration template
	SingleNetworkProductConfigurationTmpl = `package {{ .PackageName }}

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
)

const (
		ConfigureNodesNetwork ConfigPhase = iota
		ConfigureProductContractsJobs
)

var Plog = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "on_chain"}).Logger()

type OnChain struct {
		LinkContractAddress              string                 ` + "`" + `toml:"link_contract_address"` + "`" + `
		CLNodesFundingETH                float64                ` + "`" + `toml:"cl_nodes_funding_eth"` + "`" + `
		CLNodesFundingLink               float64                ` + "`" + `toml:"cl_nodes_funding_link"` + "`" + `
		VerificationTimeoutSec           time.Duration          ` + "`" + `toml:"verification_timeout_sec"` + "`" + `
		ContractsConfigurationTimeoutSec time.Duration          ` + "`" + `toml:"contracts_configuration_timeout_sec"` + "`" + `
		GasSettings                      *GasSettings           ` + "`" + `toml:"gas_settings"` + "`" + `
		Verify                           bool                   ` + "`" + `toml:"verify"` + "`" + `
		DeployedContracts                *DeployedContracts     ` + "`" + `toml:"deployed_contracts"` + "`" + `
}

type DeployedContracts struct {
	SomeContractAddr string ` + "`" + `toml:"some_contract_addr"` + "`" + `
}


type GasSettings struct {
		FeeCapMultiplier int64 ` + "`" + `toml:"fee_cap_multiplier"` + "`" + `
		TipCapMultiplier int64 ` + "`" + `toml:"tip_cap_multiplier"` + "`" + `
}

type Jobs struct {
		ConfigPollIntervalSeconds time.Duration ` + "`" + `toml:"config_poll_interval_sec"` + "`" + `
		MaxTaskDurationSec        time.Duration ` + "`" + `toml:"max_task_duration_sec"` + "`" + `
}

type ConfigPhase int

// deployLinkAndMint is a universal action that deploys link token and mints required amount of LINK token for all the nodes.
func deployLinkAndMint(ctx context.Context, in *Cfg, c *ethclient.Client, auth *bind.TransactOpts, rootAddr string, transmitters []common.Address) (*link_token.LinkToken, error) {
	addr, tx, lt, err := link_token.DeployLinkToken(auth, c)
	if err != nil {
		return nil, fmt.Errorf("could not create link token contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, tx)
	if err != nil {
		return nil, err
	}
	Plog.Info().Str("Address", addr.Hex()).Msg("Deployed link token contract")
	tx, err = lt.GrantMintRole(auth, common.HexToAddress(rootAddr))
	if err != nil {
		return nil, fmt.Errorf("could not grant mint role: %w", err)
	}
	_, err = bind.WaitMined(ctx, c, tx)
	if err != nil {
		return nil, err
	}
	// mint for public keys of nodes directly instead of transferring
	for _, transmitter := range transmitters {
		amount := new(big.Float).Mul(big.NewFloat(in.OnChain.CLNodesFundingLink), big.NewFloat(1e18))
		amountWei, _ := amount.Int(nil)
		Plog.Info().Msgf("Minting LINK for transmitter address: %s", transmitter.Hex())
		tx, err = lt.Mint(auth, transmitter, amountWei)
		if err != nil {
			return nil, fmt.Errorf("could not transfer link token contract: %w", err)
		}
		_, err = bind.WaitMined(ctx, c, tx)
		if err != nil {
			return nil, err
		}
	}
	return lt, nil
}


func configureContracts(in *Cfg, c *ethclient.Client, auth *bind.TransactOpts, cl []*clclient.ChainlinkClient, rootAddr string, transmitters []common.Address) (*DeployedContracts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), in.OnChain.ContractsConfigurationTimeoutSec*time.Second)
	defer cancel()
	Plog.Info().Msg("Deploying LINK token contract")
	_, err := deployLinkAndMint(ctx, in, c, auth, rootAddr, transmitters)
	if err != nil {
		return nil, fmt.Errorf("could not create link token contract and mint: %w", err)
	}
	// TODO: use client and deploy your contracts
	return &DeployedContracts{
		SomeContractAddr: "",
	}, nil
}

func configureJobs(in *Cfg, clNodes []*clclient.ChainlinkClient, contracts *DeployedContracts) error {
	bootstrapNode := clNodes[0]
	workerNodes := clNodes[1:]
	// TODO: define your jobs
	job := ""
	_, _, err := bootstrapNode.CreateJobRaw(job)
	if err != nil {
		return fmt.Errorf("creating bootstrap job have failed: %w", err)
	}

	for _, chainlinkNode := range workerNodes {
		// TODO: define your job for nodes here
		job := ""
		_, _, err = chainlinkNode.CreateJobRaw(job)
		if err != nil {
			return fmt.Errorf("creating job on node have failed: %w", err)
		}
	}
	return nil
}

// DefaultProductConfiguration is default product configuration that includes:
// - Deploying required prerequisites (LINK token, shared contracts)
// - Applying product-specific changesets
// - Creating cldf.Environment, connecting to components, see *Cfg fields
// - Generating CL nodes configs
// All the data can be added *Cfg struct like and is synced between local machine and remote environment
// so later both local and remote tests can use it.
func DefaultProductConfiguration(in *Cfg, phase ConfigPhase) error {
	pkey := GetNetworkPrivateKey()
	if pkey == "" {
		return fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	switch phase {
	case ConfigureNodesNetwork:
		Plog.Info().Msg("Applying default CL nodes configuration")
		node := in.Blockchains[0].Out.Nodes[0]
		chainID := in.Blockchains[0].ChainID
		// configure node set and generate CL nodes configs
		netConfig := fmt.Sprintf(` + "`" + `
    [[EVM]]
    LogPollInterval = '1s'
    BlockBackfillDepth = 100
    LinkContractAddress = '%s'
    ChainID = '%s'
    MinIncomingConfirmations = 1
    MinContractPayment = '0.0000001 link'
    FinalityDepth = %d

    [[EVM.Nodes]]
    Name = 'default'
    WsUrl = '%s'
    HttpUrl = '%s'

    [Feature]
    FeedsManager = true
    LogPoller = true
    UICSAKeys = true
    [OCR2]
    Enabled = true
    SimulateTransactions = false
    DefaultTransactionQueueDepth = 1
    [P2P.V2]
    Enabled = true
    ListenAddresses = ['0.0.0.0:6690']

	   [Log]
   JSONConsole = true
   Level = 'debug'
   [Pyroscope]
   ServerAddress = 'http://host.docker.internal:4040'
   Environment = 'local'
   [WebServer]
    SessionTimeout = '999h0m0s'
    HTTPWriteTimeout = '3m'
   SecureCookies = false
   HTTPPort = 6688
   [WebServer.TLS]
   HTTPSPort = 0
    [WebServer.RateLimit]
    Authenticated = 5000
    Unauthenticated = 5000
   [JobPipeline]
   [JobPipeline.HTTPRequest]
   DefaultTimeout = '1m'
    [Log.File]
    MaxSize = '0b'
` + "`" + `, in.OnChain.LinkContractAddress, chainID, 5, node.InternalWSUrl, node.InternalHTTPUrl)
		for _, nodeSpec := range in.NodeSets[0].NodeSpecs {
			nodeSpec.Node.TestConfigOverrides = netConfig
		}
		Plog.Info().Msg("Nodes network configuration is finished")
	case ConfigureProductContractsJobs:
		Plog.Info().Msg("Connecting to CL nodes")
		nodeClients, err := clclient.New(in.NodeSets[0].Out.CLNodes)
		if err != nil {
			return err
		}
		transmitters := make([]common.Address, 0)
		ethKeyAddresses := make([]string, 0)
		for i, nc := range nodeClients {
			addr, err := nc.ReadPrimaryETHKey(in.Blockchains[0].ChainID)
			if err != nil {
				return err
			}
			ethKeyAddresses = append(ethKeyAddresses, addr.Attributes.Address)
			transmitters = append(transmitters, common.HexToAddress(addr.Attributes.Address))
			Plog.Info().
				Int("Idx", i).
				Str("ETH", addr.Attributes.Address).
				Msg("Node info")
		}
		// ETH examples
		c, auth, rootAddr, err := ETHClient(in)
		if err != nil {
			return fmt.Errorf("could not create basic eth client: %w", err)
		}
		for _, addr := range ethKeyAddresses {
			if err := FundNodeEIP1559(c, pkey, addr, in.OnChain.CLNodesFundingETH); err != nil {
				return err
			}
		}
		contracts, err := configureContracts(in, c, auth, nodeClients, rootAddr, transmitters)
		if err != nil {
			return err
		}
		if err := configureJobs(in, nodeClients, contracts); err != nil {
			return err
		}
		Plog.Info().Str("BootstrapNode", in.NodeSets[0].Out.CLNodes[0].Node.ExternalURL).Send()
		for _, n := range in.NodeSets[0].Out.CLNodes[1:] {
			Plog.Info().Str("Node", n.Node.ExternalURL).Send()
		}
		in.OnChain.DeployedContracts = contracts
	}
	return nil
}
`
	// JustFileTmpl is a Justfile template used for building and publishing Docker images
	JustFileTmpl = `set fallback

# Default: show available recipes
default:
    just --list

clean:
    rm -rf compose/ blockscout/

build-fakes:
    just fakes/build

push-fakes:
    just fakes/push

# Rebuild CLI
cli:
    pushd cmd/{{ .CLIName }} > /dev/null && go install -ldflags="-X main.Version=1.0.0" . && popd > /dev/null`
)

/* Template params in heirarchical order, module -> file(table test) -> test */

// GoModParams params for generating go.mod file
type GoModParams struct {
	ModuleName     string
	RuntimeVersion string
}

// ConfigTOMLParams default env.toml params
type ConfigTOMLParams struct {
	PackageName string
	Nodes       int
	NodeIndices []int
}

// JustfileParams Justfile params
type JustfileParams struct {
	PackageName string
	CLIName     string
}

// CLICompletionParams cli.go file params
type CLICompletionParams struct {
	PackageName string
	CLIName     string
}

// CLIParams cli.go file params
type CLIParams struct {
	PackageName     string
	CLIName         string
	DevEnvPkgImport string
}

// CLDFParams cldf.go file params
type CLDFParams struct {
	PackageName string
}

// ToolsParams tools.go file params
type ToolsParams struct {
	PackageName string
}

// ConfigParams config.go file params
type ConfigParams struct {
	PackageName string
}

// EnvParams environment.go file params
type EnvParams struct {
	PackageName string
}

// ProductConfigurationSimple product_configuration.go file params
type ProductConfigurationSimple struct {
	PackageName string
}

// TableTestParams params for generating a table test
type TableTestParams struct {
	Package       string
	TableTestName string
	TestCases     []TestCaseParams
	WorkloadCode  string
	GunCode       string
}

// TestCaseParams params for generating a test case
type TestCaseParams struct {
	Name    string
	RunFunc string
}

/* Codegen logic */

// EnvBuilder builder for load test codegen
type EnvBuilder struct {
	nodes       int
	outputDir   string
	packageName string
	cliName     string
	moduleName  string
}

// EnvCodegen is a load test code generator that creates workload and chaos experiments
type EnvCodegen struct {
	cfg *EnvBuilder
}

// NewEnvBuilder creates a new Chainlink Cluster developer environment
func NewEnvBuilder(cliName string, nodes int) *EnvBuilder {
	return &EnvBuilder{
		cliName:     cliName,
		nodes:       nodes,
		packageName: "devenv",
		outputDir:   "devenv",
		moduleName:  fmt.Sprintf("github.com/smartcontractkit/%s/devenv", cliName),
	}
}

// OutputDir sets the output directory for generated files
func (g *EnvBuilder) OutputDir(dir string) *EnvBuilder {
	g.outputDir = dir
	return g
}

// GoModName sets the Go module name for the generated project
func (g *EnvBuilder) GoModName(name string) *EnvBuilder {
	g.moduleName = name
	return g
}

// validate verifier that we can build codegen with provided params
// empty for now, add validation here if later it'd become more complex
func (g *EnvBuilder) validate() error {
	return nil
}

// Validate validate generation params
// for now it's empty but for more complex mutually exclusive cases we should
// add validation here
func (g *EnvBuilder) Build() (*EnvCodegen, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}
	return &EnvCodegen{g}, nil
}

// Read read K8s namespace and find all the pods
// some pods may be crashing but it doesn't matter for code generation
func (g *EnvCodegen) Read() error {
	return nil
}

// Write generates a complete boilerplate, can be multiple files
func (g *EnvCodegen) Write() error {
	// Create output directory
	if err := os.MkdirAll( //nolint:gosec
		g.cfg.outputDir,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate go.mod file
	goModContent, err := g.GenerateGoMod()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "go.mod"),
		[]byte(goModContent),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Generate Justfile
	justContents, err := g.GenerateJustfile()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "Justfile"),
		[]byte(justContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CLI completion file: %w", err)
	}

	// Generate default env.toml file
	tomlContents, err := g.GenerateDefaultTOMLConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "env.toml"),
		[]byte(tomlContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write default env.toml file: %w", err)
	}

	cliDir := filepath.Join(g.cfg.outputDir, "cmd", g.cfg.cliName)
	if err := os.MkdirAll( //nolint:gosec
		cliDir,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create CLI directory: %w", err)
	}

	// Generate CLI file
	cliContents, err := g.GenerateCLI()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(cliDir, fmt.Sprintf("%s.go", g.cfg.cliName)),
		[]byte(cliContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CLI file: %w", err)
	}

	// Generate completion file
	completionContents, err := g.GenerateCLICompletion()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(cliDir, "completion.go"),
		[]byte(completionContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write CLI completion file: %w", err)
	}

	// Generate tools.go
	toolsContents, err := g.GenerateDebugTools()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "tools.go"),
		[]byte(toolsContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Generate config.go
	configFileContents, err := g.GenerateConfig()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "config.go"),
		[]byte(configFileContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Generate cldf.go
	cldfContents, err := g.GenerateCLDF()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "cldf.go"),
		[]byte(cldfContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Generate environment.go
	envFileContents, err := g.GenerateEnvironment()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "environment.go"),
		[]byte(envFileContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write environment file: %w", err)
	}

	// Generate product_configuration.go
	prodConfigFileContents, err := g.GenerateSingleNetworkProductConfiguration()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(g.cfg.outputDir, "product_configuration.go"),
		[]byte(prodConfigFileContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write product configuration file: %w", err)
	}

	// tidy and finalize
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// nolint
	defer os.Chdir(currentDir)
	if err := os.Chdir(g.cfg.outputDir); err != nil {
		return err
	}
	log.Info().Msg("Downloading dependencies and running 'go mod tidy' ..")
	_, err = exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to tidy generated module: %w", err)
	}
	log.Info().
		Str("OutputDir", g.cfg.outputDir).
		Str("Module", g.cfg.moduleName).
		Msg("Developer environment generated")
	return nil
}

// GenerateGoMod generates a go.mod file
func (g *EnvCodegen) GenerateGoMod() (string, error) {
	data := GoModParams{
		ModuleName:     g.cfg.moduleName,
		RuntimeVersion: strings.ReplaceAll(runtime.Version(), "go", ""),
	}
	return render(GoModTemplate, data)
}

// GenerateDefaultTOMLConfig generate default env.toml config
func (g *EnvCodegen) GenerateDefaultTOMLConfig() (string, error) {
	p := ConfigTOMLParams{
		PackageName: g.cfg.packageName,
		Nodes:       g.cfg.nodes,
		NodeIndices: make([]int, g.cfg.nodes),
	}
	return render(ConfigTOMLTmpl, p)
}

// GenerateJustfile generate Justfile to build and publish Docker images
func (g *EnvCodegen) GenerateJustfile() (string, error) {
	p := JustfileParams{
		PackageName: g.cfg.packageName,
		CLIName:     g.cfg.cliName,
	}
	return render(JustFileTmpl, p)
}

// GenerateCLICompletion generate CLI completion for "go-prompt" library
func (g *EnvCodegen) GenerateCLICompletion() (string, error) {
	p := CLICompletionParams{
		PackageName: g.cfg.packageName,
		CLIName:     g.cfg.cliName,
	}
	return render(CompletionTmpl, p)
}

// GenerateCLI generate Cobra CLI
func (g *EnvCodegen) GenerateCLI() (string, error) {
	p := CLIParams{
		PackageName:     g.cfg.packageName,
		CLIName:         g.cfg.cliName,
		DevEnvPkgImport: g.cfg.moduleName,
	}
	return render(CLITmpl, p)
}

// GenerateSingleNetworkProductConfiguration generate a single-network product configuration
func (g *EnvCodegen) GenerateSingleNetworkProductConfiguration() (string, error) {
	p := ProductConfigurationSimple{
		PackageName: g.cfg.packageName,
	}
	return render(SingleNetworkProductConfigurationTmpl, p)
}

// GenerateEnvironment generate environment.go, our environment composition function
func (g *EnvCodegen) GenerateEnvironment() (string, error) {
	p := EnvParams{
		PackageName: g.cfg.packageName,
	}
	return render(EnvironmentTmpl, p)
}

// GenerateCLDF generate CLDF helpers
func (g *EnvCodegen) GenerateCLDF() (string, error) {
	p := CLDFParams{
		PackageName: g.cfg.packageName,
	}
	return render(CLDFTmpl, p)
}

// GenerateDebugTools generate debug tools (tracing)
func (g *EnvCodegen) GenerateDebugTools() (string, error) {
	p := ToolsParams{
		PackageName: g.cfg.packageName,
	}
	return render(DebugToolsTmpl, p)
}

// GenerateConfig generate read/write utilities for TOML configs
func (g *EnvCodegen) GenerateConfig() (string, error) {
	p := ConfigParams{
		PackageName: g.cfg.packageName,
	}
	return render(ConfigTmpl, p)
}

// GenerateTableTest generates all possible experiments for a namespace
// first generate all small pieces then insert into a table test template
func (g *EnvCodegen) GenerateTableTest() (string, error) {
	// TODO: generate a table test when we'll have chain-specific interface solidified
	return "", nil
}

// GenerateTestCases generates table test cases
func (g *EnvCodegen) GenerateTestCases() ([]TestCaseParams, error) {
	// TODO: generate test cases when we'll have chain-specific interface solidified
	return []TestCaseParams{}, nil
}

// render is just an internal function to parse and render template
func render(tmpl string, data any) (string, error) {
	parsed, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse table test template: %w", err)
	}
	var buf bytes.Buffer
	if err := parsed.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to generate table test: %w", err)
	}
	return buf.String(), err
}
