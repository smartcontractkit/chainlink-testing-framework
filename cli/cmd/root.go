package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flag options
const (
	FlagType      = "type"
	FlagNodeCount = "nodeCount"
	FlagNetwork   = "network"

	FlagConfig = "config"
)

var rootCmd = &cobra.Command{
	Use:   "ifcli",
	Short: "ifcli is a quick way of creating ephemeral Chainlink environments, and build contracts",
	Long: `By using the k8s test frameworks environment functionality built with k8s' client-go,
ifcli can create configurable and quick environment clusters that can be used for local testing.`,
}

func init() {
	rootCmd.AddCommand(buildContractsCmd)
	buildContractsCmd.Flags().StringP(FlagConfig,
		"c", "", "path to contracts config")
	_ = buildContractsCmd.MarkFlagRequired(FlagConfig)

	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP(FlagType, "t", "chainlink", "type of environment to deploy")
	createCmd.Flags().IntP(FlagNodeCount, "c", 3, "number of Chainlink nodes to deploy")
	createCmd.Flags().StringP(FlagNetwork,
		"n", "ethereum_hardhat", "the network to deploy the Chainlink cluster on")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
