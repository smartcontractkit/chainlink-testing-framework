package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	FlagType      = "type"
	FlagNodeCount = "nodeCount"
	FlagNetwork   = "network"
)

var rootCmd = &cobra.Command{
	Use:   "clEnv",
	Short: "clEnv is a quick way of creating ephemeral Chainlink environments",
	Long: `By using the k8s test frameworks environment functionality built with k8s' client-go,
clEnv can create configurable and quick environment clusters that can be used for local testing.`,
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.PersistentFlags().StringP(FlagType, "t", "chainlink", "type of environment to deploy")
	rootCmd.PersistentFlags().IntP(FlagNodeCount, "c", 3, "number of Chainlink nodes to deploy")
	rootCmd.PersistentFlags().StringP(FlagNetwork, "n", "ethereum_hardhat", "the network to deploy the Chainlink cluster on")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
