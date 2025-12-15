package framework

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v72/github"
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

	CTFCacheDir              = ".local/share/ctf"
	DefaultGitHubOwner       = "smartcontractkit"
	DefaultGitHubRepo        = "chainlink-testing-framework"
	DefaultObservabilityPath = "framework/observability"
)

// resolveObservabilitySource determines where to load observability files from based on the source parameter
// - Empty string: use embedded files
// - file:// prefix: use local filesystem
// - http(s):// prefix: download and cache from remote URL
func resolveObservabilitySource(source string) (fs.FS, string, error) {
	if source == "" {
		// Default: use embedded files
		return EmbeddedObservabilityFiles, "observability", nil
	}

	if strings.HasPrefix(source, "file://") {
		// Local filesystem path
		localPath := strings.TrimPrefix(source, "file://")
		if _, err := os.Stat(localPath); err != nil {
			return nil, "", fmt.Errorf("local observability path does not exist: %s: %w", localPath, err)
		}
		return os.DirFS(localPath), ".", nil
	}

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		// Remote URL: download and cache
		cachePath, err := downloadAndCacheObservabilityFiles(source)
		if err != nil {
			return nil, "", fmt.Errorf("failed to download observability files: %w", err)
		}
		return os.DirFS(cachePath), ".", nil
	}

	return nil, "", fmt.Errorf("invalid source format: %s (must be empty, file://, or http(s)://)", source)
}

// downloadAndCacheObservabilityFiles downloads observability files from a GitHub URL and caches them
func downloadAndCacheObservabilityFiles(url string) (string, error) {
	// Parse URL to extract repo info and create cache key
	// Expected format: https://github.com/owner/repo/tree/ref/path/to/observability
	owner, repo, ref, path, err := parseGitHubURL(url)
	if err != nil {
		return "", err
	}

	// Create cache directory using just the ref
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	cachedPath := filepath.Join(homeDir, CTFCacheDir, "observability", ref)

	// Check if already cached and has content
	if info, err := os.Stat(cachedPath); err == nil && info.IsDir() {
		// Verify the cache directory has files
		entries, err := os.ReadDir(cachedPath)
		if err == nil && len(entries) > 0 {
			L.Debug().Msgf("Using cached observability files from %s", cachedPath)
			return cachedPath, nil
		}
		L.Debug().Msg("Cache directory exists but is empty, re-downloading")
	}

	L.Info().Msgf("Downloading observability files from GitHub: %s/%s@%s (path: %s)", owner, repo, ref, path)

	// Create GitHub client with optional authentication and timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var client *github.Client
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		L.Debug().Msg("Using authenticated GitHub client")
		client = github.NewClient(nil).WithAuthToken(token)
	} else {
		L.Debug().Msg("Using unauthenticated GitHub client")
		client = github.NewClient(nil)
	}

	// Download directory contents recursively
	if err := downloadDirectoryRecursive(ctx, client, owner, repo, ref, path, cachedPath); err != nil {
		return "", fmt.Errorf("failed to download directory: %w", err)
	}

	L.Info().Msgf("Observability files cached at: %s", cachedPath)
	return cachedPath, nil
}

// parseGitHubURL parses a GitHub URL and extracts owner, repo, ref, and path
func parseGitHubURL(url string) (owner, repo, ref, path string, err error) {
	// Expected format: https://github.com/owner/repo/tree/ref/path/to/observability
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "http://github.com/")
	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return "", "", "", "", fmt.Errorf("invalid GitHub URL format (expected: https://github.com/owner/repo/tree|blob/ref/path)")
	}

	owner = parts[0]
	repo = parts[1]
	treeOrBlob := parts[2]
	ref = parts[3]
	path = strings.Join(parts[4:], "/")

	if treeOrBlob != "tree" && treeOrBlob != "blob" {
		return "", "", "", "", fmt.Errorf("unsupported GitHub URL type: %s (expected 'tree' or 'blob')", treeOrBlob)
	}

	return owner, repo, ref, path, nil
}

// downloadDirectoryRecursive recursively downloads a directory from GitHub
func downloadDirectoryRecursive(ctx context.Context, client *github.Client, owner, repo, ref, path, destPath string) error {
	// Get directory contents
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return fmt.Errorf("failed to get directory contents: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(destPath, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Process each item in the directory
	for _, item := range directoryContent {
		if item.GetName() == "README.md" {
			continue
		}

		itemPath := item.GetPath()
		itemName := item.GetName()
		targetPath := filepath.Join(destPath, itemName)

		switch item.GetType() {
		case "file":
			// Download file
			fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, itemPath, &github.RepositoryContentGetOptions{
				Ref: ref,
			})
			if err != nil {
				return fmt.Errorf("failed to get file %s: %w", itemPath, err)
			}

			content, err := fileContent.GetContent()
			if err != nil {
				return fmt.Errorf("failed to decode file %s: %w", itemPath, err)
			}

			if err := os.WriteFile(targetPath, []byte(content), 0o644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}

		case "dir":
			// Recursively download subdirectory
			if err := downloadDirectoryRecursive(ctx, client, owner, repo, ref, itemPath, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// extractAllFiles goes through the embedded directory and extracts all files to the current directory
func extractAllFiles(embeddedDir string) error {
	return extractAllFilesFromFS(EmbeddedObservabilityFiles, embeddedDir)
}

// extractAllFilesFromFS goes through a filesystem and extracts all files to the current directory
func extractAllFilesFromFS(fsys fs.FS, embeddedDir string) error {
	// Get current working directory where CLI is running
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Walk through the files
	err = fs.WalkDir(fsys, embeddedDir, func(path string, d fs.DirEntry, err error) error {
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

		// Read file content from file system
		content, err := fs.ReadFile(fsys, path)
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

func BlockScoutUp(url, chainID string) error {
	L.Info().Msg("Creating local Blockscout stack")
	if err := extractAllFiles("observability"); err != nil {
		return err
	}
	os.Setenv("BLOCKSCOUT_RPC_URL", url)
	os.Setenv("BLOCKSCOUT_CHAID_ID", chainID)
	// old migrations for v15 is still applied somehow, cleaning up DB helps
	if err := RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		rm -rf blockscout-db-data && \
		rm -rf logs && \
		rm -rf redis-data && \
		rm -rf stats-db-data && \
		rm -rf dets
	`, filepath.Join("blockscout", "services"))); err != nil {
		return err
	}
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
	return RunCommand("bash", "-c", "rm -rf blockscout/")
}

// ObservabilityUpOnlyLoki slim stack with only Loki to verify specific logs of CL nodes or services in tests
func ObservabilityUpOnlyLoki() error {
	return ObservabilityUpOnlyLokiWithSource("")
}

// ObservabilityUpOnlyLokiWithSource slim stack with only Loki using custom observability file source
// source can be:
// - "" (empty): use embedded files (default)
// - "file:///path/to/observability": use local filesystem
// - "https://github.com/owner/repo/tree/tag/framework/observability": download from GitHub
func ObservabilityUpOnlyLokiWithSource(source string) error {
	L.Info().Msg("Creating local observability stack")
	fsys, dir, err := resolveObservabilitySource(source)
	if err != nil {
		return err
	}
	if err := extractAllFilesFromFS(fsys, dir); err != nil {
		return err
	}
	_ = DefaultNetwork(nil)
	if err := NewPromtail(); err != nil {
		return err
	}
	err = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d loki grafana
	`, "compose"))
	if err != nil {
		return err
	}
	fmt.Println()
	L.Info().Msgf("Loki: %s", LocalLogsURL)
	return nil
}

// ObservabilityUp standard stack with logs/metrics for load testing and observability
func ObservabilityUp() error {
	return ObservabilityUpWithSource("")
}

// ObservabilityUpWithSource standard stack with logs/metrics using custom observability file source
// source can be:
// - "" (empty): use embedded files (default)
// - "file:///path/to/observability": use local filesystem
// - "https://github.com/owner/repo/tree/tag/framework/observability": download from GitHub
func ObservabilityUpWithSource(source string) error {
	L.Info().Msg("Creating local observability stack")
	fsys, dir, err := resolveObservabilitySource(source)
	if err != nil {
		return err
	}
	if err := extractAllFilesFromFS(fsys, dir); err != nil {
		return err
	}
	_ = DefaultNetwork(nil)
	if err := NewPromtail(); err != nil {
		return err
	}
	err = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose up -d otel-collector prometheus loki grafana
	`, "compose"))
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
	return ObservabilityUpFullWithSource("")
}

// ObservabilityUpFullWithSource full stack for load testing using custom observability file source
// source can be:
// - "" (empty): use embedded files (default)
// - "file:///path/to/observability": use local filesystem
// - "https://github.com/owner/repo/tree/tag/framework/observability": download from GitHub
func ObservabilityUpFullWithSource(source string) error {
	L.Info().Msg("Creating full local observability stack")
	fsys, dir, err := resolveObservabilitySource(source)
	if err != nil {
		return err
	}
	if err := extractAllFilesFromFS(fsys, dir); err != nil {
		return err
	}
	_ = DefaultNetwork(nil)
	if err := NewPromtail(); err != nil {
		return err
	}
	err = RunCommand("bash", "-c", fmt.Sprintf(`
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
	_ = RunCommand("bash", "-c", fmt.Sprintf(`
		cd %s && \
		docker compose down -v && docker rm -f promtail
	`, "compose"))
	_ = RunCommand("bash", "-c", "rm -rf compose/")
	return nil
}
