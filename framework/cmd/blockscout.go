package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"path/filepath"
)

func blockscoutUp() error {
	framework.L.Info().Msg("Creating local Blockscout stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	err := runCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, "blockscout"))
	if err != nil {
		return err
	}
	fmt.Println()
	framework.L.Info().Msgf("Blockscout is up at: %s", "http://localhost")
	return nil
}

func blockscoutDown() error {
	framework.L.Info().Msg("Removing local Blockscout stack")
	err := runCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, "blockscout"))
	if err != nil {
		return err
	}
	err = runCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		rm -rf blockscout-db-data && \
		rm -rf logs && \
		rm -rf redis-data && \
		rm -rf stats-db-data
	`, filepath.Join("blockscout", "services")))
	return nil
}
