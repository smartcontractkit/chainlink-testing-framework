package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/spf13/cobra"
)

var buildContractsCmd = &cobra.Command{
	Use:   "build_contracts",
	Short: "Builds all contracts",
	Long:  `Builds all contracts for provided versions`,
	RunE:  buildContractsE,
}

func buildContractsE(cmd *cobra.Command, _ []string) error {
	log.Info().Msg("Building contracts")
	configPath, err := cmd.Flags().GetString(FlagConfig)
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}
	eb := contracts.NewEthereumContractBuilder(&cfg.Contracts.Ethereum)
	targets, err := eb.Targets()
	if err != nil {
		return err
	}
	log.Debug().Interface("Targets", targets).Msg("Building contracts")
	err = eb.UpdateExternalSources()
	if err != nil {
		return err
	}
	err = eb.GenerateBindings(targets)
	if err != nil {
		return err
	}
	return nil
}
