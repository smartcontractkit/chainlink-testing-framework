package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func blockscoutUp(url string) error {
	framework.L.Info().Msg("Creating local Blockscout stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	os.Setenv("BLOCKSCOUT_RPC_URL", url)
	err := framework.RunCommand("bash", "-c", fmt.Sprintf(`
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

func blockscoutDown(url string) error {
	framework.L.Info().Msg("Removing local Blockscout stack")
	os.Setenv("BLOCKSCOUT_RPC_URL", url)
	err := framework.RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, "blockscout"))
	if err != nil {
		return err
	}
	return framework.RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		rm -rf blockscout-db-data && \
		rm -rf logs && \
		rm -rf redis-data && \
		rm -rf stats-db-data
	`, filepath.Join("blockscout", "services")))
}
