package cmd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create an environment",
	Long:  `Deploy a configurable environment using k8s that can be used for local testing.`,
	RunE:  createRunE,
}

func createRunE(cmd *cobra.Command, _ []string) error {
	envType, err := cmd.Flags().GetString(FlagType)
	if err != nil {
		return err
	}
	nodes, err := cmd.Flags().GetInt(FlagNodeCount)
	if err != nil {
		return err
	}
	network, err := cmd.Flags().GetString(FlagNetwork)
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig(tools.ProjectRoot)
	if err != nil {
		return err
	}
	cfg.Network = network
	networkConfig, err := client.NewNetworkFromConfig(cfg)
	if err != nil {
		return err
	}

	var env environment.Environment
	switch envType {
	case "chainlink":
		envSpec := environment.NewChainlinkCluster(nodes)
		env, err = environment.NewK8sEnvironment(envSpec, cfg, networkConfig)
	default:
		return fmt.Errorf("invalid environment type '%s' specified", envType)
	}
	if err != nil {
		return err
	}
	log.Info().Str("Namespace", env.ID()).Msgf("Environment created")
	return nil
}
