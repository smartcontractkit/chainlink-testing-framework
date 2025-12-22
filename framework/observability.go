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

	CTFCacheDir = ".local/share/ctf"
)

// getObservabilityDir returns the fixed directory where observability files are extracted
func getObservabilityDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, CTFCacheDir), nil
}

// extractAllFiles goes through the embedded directory and extracts all files to the fixed observability directory
func extractAllFiles(embeddedDir string) error {
	// Get fixed observability directory
	obsDir, err := getObservabilityDir()
	if err != nil {
		return err
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
		content, err := fs.ReadFile(EmbeddedObservabilityFiles, path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Determine the target path (strip out the `embeddedDir` part)
		relativePath, err := filepath.Rel(embeddedDir, path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", path, err)
		}
		targetPath := filepath.Join(obsDir, relativePath)

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

func BlockScoutUp(url, chainID string) error {
	L.Info().Msg("Creating local Blockscout stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	obsDir, err := getObservabilityDir()
	if err != nil {
		return err
	}
	blockscoutDir := filepath.Join(obsDir, "blockscout")
	os.Setenv("BLOCKSCOUT_RPC_URL", url)
	os.Setenv("BLOCKSCOUT_CHAIN_ID", chainID)
	// old migrations for v15 is still applied somehow, cleaning up DB helps
	if err := RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		rm -rf blockscout-db-data && \
		rm -rf logs && \
		rm -rf redis-data && \
		rm -rf stats-db-data && \
		rm -rf dets
	`, filepath.Join(blockscoutDir, "services"))); err != nil {
		return err
	}
	err = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, blockscoutDir))
	if err != nil {
		return err
	}
	fmt.Println()
	L.Info().Msgf("Blockscout is up at: %s", "http://localhost")
	return nil
}

func BlockScoutDown(url string) error {
	L.Info().Msg("Removing local Blockscout stack")
	obsDir, err := getObservabilityDir()
	if err != nil {
		return err
	}
	blockscoutDir := filepath.Join(obsDir, "blockscout")
	os.Setenv("BLOCKSCOUT_RPC_URL", url)
	return RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v
	`, blockscoutDir))
}

// ObservabilityUpOnlyLoki slim stack with only Loki to verify specific logs of CL nodes or services in tests
func ObservabilityUpOnlyLoki() error {
	L.Info().Msg("Creating local observability stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	obsDir, err := getObservabilityDir()
	if err != nil {
		return err
	}
	composeDir := filepath.Join(obsDir, "compose")
	_ = DefaultNetwork(nil)
	if err := NewPromtail(); err != nil {
		return err
	}
	err = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d loki grafana
	`, composeDir))
	if err != nil {
		return err
	}
	fmt.Println()
	L.Info().Msgf("Loki: %s", LocalLogsURL)
	return nil
}

// ObservabilityUp standard stack with logs/metrics for load testing and observability
func ObservabilityUp() error {
	L.Info().Msg("Creating local observability stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	obsDir, err := getObservabilityDir()
	if err != nil {
		return err
	}
	composeDir := filepath.Join(obsDir, "compose")
	_ = DefaultNetwork(nil)
	if err := NewPromtail(); err != nil {
		return err
	}
	err = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d otel-collector prometheus loki grafana
	`, composeDir))
	if err != nil {
		return err
	}
	fmt.Println()
	L.Info().Msgf("Loki: %s", LocalLogsURL)
	L.Info().Msgf("Prometheus: %s", LocalPrometheusURL)
	L.Info().Msgf("CL Node Errors: %s", LocalCLNodeErrorsURL)
	L.Info().Msgf("Workflow Engine: %s", LocalWorkflowEngineURL)
	return nil
}

// ObservabilityUpFull full stack for load testing and performance investigations
func ObservabilityUpFull() error {
	L.Info().Msg("Creating full local observability stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	obsDir, err := getObservabilityDir()
	if err != nil {
		return err
	}
	composeDir := filepath.Join(obsDir, "compose")
	_ = DefaultNetwork(nil)
	if err := NewPromtail(); err != nil {
		return err
	}
	err = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d
	`, composeDir))
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
	obsDir, err := getObservabilityDir()
	if err != nil {
		return err
	}
	composeDir := filepath.Join(obsDir, "compose")
	_ = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v && docker rm -f promtail
	`, composeDir))
	return nil
}
