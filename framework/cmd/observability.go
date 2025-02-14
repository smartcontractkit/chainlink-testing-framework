package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"go.uber.org/ratelimit"
	"net/http"
)

func loadLogs(rawURL, dirPath string, rps, chunks int) error {
	framework.L.Info().Msg("Loading logs into Loki")
	sources := 0
	if rawURL != "" {
		sources++
	}
	if dirPath != "" {
		sources++
	}
	if sources != 1 {
		L.Error().Msg("Usage: provide exactly one of -raw-url or -dir")
		return nil
	}
	jobID := uuid.New().String()[0:5]
	L.Info().Msgf("Using unique job identifier: %s", jobID)
	limiter := ratelimit.New(rps)
	if rawURL != "" {
		L.Info().Msg("Downloading raw logs from URL")
		resp, err := http.Get(rawURL)
		if err != nil {
			L.Error().Err(err).Msg("Error downloading raw logs")
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			L.Error().Msgf("Non-success response downloading raw logs: %s", resp.Status)
			return nil
		}

		if err := processAndUploadLog(rawURL, resp.Body, limiter, chunks, jobID); err != nil {
			L.Error().Err(err).Msg("Error processing raw logs")
			return nil
		}
	} else if dirPath != "" {
		L.Info().Msgf("Processing directory: %s", dirPath)
		if err := processAndUploadDir(dirPath, limiter, chunks, jobID); err != nil {
			L.Error().Err(err).Msg("Error processing directory")
			return nil
		}
	}
	framework.L.Info().Str("JobID", jobID).Str("URL", grafanaURL+jobID+grafanaURL2).Msg("Upload complete")
	return nil
}

func observabilityUp() error {
	framework.L.Info().Msg("Creating local observability stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	if err := framework.NewPromtail(); err != nil {
		return err
	}
	err := runCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, "compose"))
	if err != nil {
		return err
	}
	fmt.Println()
	framework.L.Info().Msgf("Loki: %s", LocalLogsURL)
	framework.L.Info().Msgf("Prometheus: %s", LocalPrometheusURL)
	framework.L.Info().Msgf("PostgreSQL: %s", LocalPostgresDebugURL)
	framework.L.Info().Msgf("Pyroscope: %s", LocalPyroScopeURL)
	return nil
}

func observabilityDown() error {
	framework.L.Info().Msg("Removing local observability stack")
	err := runCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v && docker rm -f promtail
	`, "compose"))
	if err != nil {
		return err
	}
	return nil
}
