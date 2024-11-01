package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func observabilityUp() error {
	framework.L.Info().Msg("Creating local observability stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	err := runCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, "compose"))
	if err != nil {
		return err
	}
	framework.L.Info().Msg("Done")
	fmt.Println()
	framework.L.Info().Msgf("Loki: %s", LocalLogsURL)
	framework.L.Info().Msgf("All logs: %s", "{job=\"ctf\"}")
	framework.L.Info().Msgf("By log level: %s", "{job=\"ctf\", container=~\"node.*\"} |= \"WARN|INFO|DEBUG\"")
	framework.L.Info().Msgf("Pyroscope: %s", LocalPyroScopeURL)
	return nil
}

func observabilityDown() error {
	framework.L.Info().Msg("Removing local observability stack")
	err := runCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, "compose"))
	if err != nil {
		return err
	}
	framework.L.Info().Msg("Done")
	return nil
}
