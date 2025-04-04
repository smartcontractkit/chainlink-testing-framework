package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func loadLogs(rawURL, dirPath string, rps, chunks int) error {
	if rawURL == "" && dirPath == "" {
		return fmt.Errorf("at least one source must be provided, either -u $url or -d $dir")
	}
	jobID := uuid.New().String()[0:5]
	framework.L.Info().Str("JobID", jobID).Msg("Loading logs into Loki")
	limiter := ratelimit.New(rps)
	if rawURL != "" {
		L.Info().Msg("Downloading raw logs from URL")
		//nolint:gosec
		resp, err := http.Get(rawURL)
		if err != nil {
			return errors.Wrap(err, "error downloading raw logs")
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			return fmt.Errorf("non-success response code when downloading raw logs: %s", resp.Status)
		}

		if err := processAndUploadLog(rawURL, resp.Body, limiter, chunks, jobID); err != nil {
			return errors.Wrap(err, "error processing raw logs")
		}
	} else if dirPath != "" {
		L.Info().Msgf("Processing directory: %s", dirPath)
		if err := processAndUploadDir(dirPath, limiter, chunks, jobID); err != nil {
			return errors.Wrapf(err, "error processing directory: %s", dirPath)
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
	err := framework.RunCommand("bash", "-c", fmt.Sprintf(`
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
	framework.L.Info().Msgf("CL Node Errors: %s", LocalCLNodeErrorsURL)
	framework.L.Info().Msgf("Workflow Engine: %s", LocalWorkflowEngineURL)
	return nil
}

func observabilityDown() error {
	framework.L.Info().Msg("Removing local observability stack")
	err := framework.RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v && docker rm -f promtail
	`, "compose"))
	if err != nil {
		return err
	}
	return nil
}
