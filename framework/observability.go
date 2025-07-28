package framework

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed observability/*
var EmbeddedObservabilityFiles embed.FS

const (
	LocalGrafanaBaseURL    = "http://localhost:3000"
	LocalLokiBaseURL       = "http://localhost:3030"
	LocalPrometheusBaseURL = "http://localhost:9099"
	LocalCLNodeErrorsURL   = "http://localhost:3000/d/a7de535b-3e0f-4066-bed7-d505b6ec9ef1/cl-node-errors?orgId=1&refresh=5s"
	LocalWorkflowEngineURL = "http://localhost:3000/d/ce589a98-b4be-4f80-bed1-bc62f3e4414a/workflow-engine?orgId=1&refresh=5s&from=now-15m&to=now"
	LocalLogsURL           = "http://localhost:3000/explore?panes=%7B%22qZw%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D%5C%22ctf%5C%22%7D%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22code%22%7D%5D,%22range%22:%7B%22from%22:%22now-15m%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1"
	LocalPrometheusURL     = "http://localhost:3000/explore?panes=%7B%22qZw%22:%7B%22datasource%22:%22PBFA97CFB590B2093%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%22,%22range%22:true,%22datasource%22:%7B%22type%22:%22prometheus%22,%22uid%22:%22PBFA97CFB590B2093%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-15m%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1"
	LocalPostgresDebugURL  = "http://localhost:3000/d/000000039/postgresql-database?orgId=1&refresh=5s&var-DS_PROMETHEUS=PBFA97CFB590B2093&var-interval=$__auto_interval_interval&var-namespace=&var-release=&var-instance=postgres_exporter_0:9187&var-datname=All&var-mode=All&from=now-15m&to=now"
	LocalPyroScopeURL      = "http://localhost:4040/?query=process_cpu%3Acpu%3Ananoseconds%3Acpu%3Ananoseconds%7Bservice_name%3D%22chainlink-node%22%7D&from=now-15m"
)

// extractAllFiles goes through the embedded directory and extracts all files to the current directory
func extractAllFiles(embeddedDir string) error {
	// Get current working directory where CLI is running
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Walk through the embedded files
	err = fs.WalkDir(EmbeddedObservabilityFiles, embeddedDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the directory: %w", err)
		}
		if strings.Contains(path, "README.md") {
			return nil
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read file content from embedded file system
		content, err := EmbeddedObservabilityFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Determine the target path (strip out the `embeddedDir` part)
		relativePath, err := filepath.Rel(embeddedDir, path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", path, err)
		}
		targetPath := filepath.Join(currentDir, relativePath)

		// Create target directories if necessary
		targetDir := filepath.Dir(targetPath)
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		// Write the file content to the target path
		//nolint
		err = os.WriteFile(targetPath, content, 0o777)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
		return nil
	})

	return err
}

func BlockScoutUp(url string) error {
	L.Info().Msg("Creating local Blockscout stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	os.Setenv("BLOCKSCOUT_RPC_URL", url)
	err := RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, "blockscout"))
	if err != nil {
		return err
	}
	fmt.Println()
	L.Info().Msgf("Blockscout is up at: %s", "http://localhost")
	return nil
}

func BlockScoutDown(url string) error {
	L.Info().Msg("Removing local Blockscout stack")
	os.Setenv("BLOCKSCOUT_RPC_URL", url)
	err := RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, "blockscout"))
	if err != nil {
		return err
	}
	return RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		rm -rf blockscout-db-data && \
		rm -rf logs && \
		rm -rf redis-data && \
		rm -rf stats-db-data
	`, filepath.Join("blockscout", "services")))
}

func ObservabilityUp() error {
	L.Info().Msg("Creating local observability stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	_ = DefaultNetwork(nil)
	if err := NewPromtail(); err != nil {
		return err
	}
	err := RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, "compose"))
	if err != nil {
		return err
	}
	fmt.Println()
	L.Info().Msgf("Loki: %s", LocalLogsURL)
	L.Info().Msgf("Prometheus: %s", LocalPrometheusURL)
	L.Info().Msgf("PostgreSQL: %s", LocalPostgresDebugURL)
	L.Info().Msgf("Pyroscope: %s", LocalPyroScopeURL)
	L.Info().Msgf("CL Node Errors: %s", LocalCLNodeErrorsURL)
	L.Info().Msgf("Workflow Engine: %s", LocalWorkflowEngineURL)
	return nil
}

func ObservabilityDown() error {
	L.Info().Msg("Removing local observability stack")
	err := RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v && docker rm -f promtail
	`, "compose"))
	if err != nil {
		return err
	}
	return nil
}
