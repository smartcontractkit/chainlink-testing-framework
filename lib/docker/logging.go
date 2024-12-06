package docker

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"
)

// WriteAllContainersLogs writes all Docker container logs to the default logs directory
func WriteAllContainersLogs(logger zerolog.Logger, directory string) error {
	logger.Info().Msgf("Writing Docker containers logs to %s", directory)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", directory, err)
		}
	}
	provider, err := tc.NewDockerProvider()
	if err != nil {
		return fmt.Errorf("failed to create Docker provider: %w", err)
	}
	containers, err := provider.Client().ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list Docker containers: %w", err)
	}

	eg := &errgroup.Group{}

	for _, containerInfo := range containers {
		eg.Go(func() error {
			containerName := containerInfo.Names[0]
			if shouldIgnore(logger, containerName) {
				return nil
			}
			logger.Debug().Str("Container", containerName).Msg("Collecting logs")
			logOptions := container.LogsOptions{ShowStdout: true, ShowStderr: true}
			logs, err := provider.Client().ContainerLogs(context.Background(), containerInfo.ID, logOptions)
			if err != nil {
				logger.Error().Err(err).Str("Container", containerName).Msg("failed to fetch logs for container")
				return err
			}
			logFilePath := filepath.Join(directory, fmt.Sprintf("%s.log", containerName))
			logFile, err := os.Create(logFilePath)
			if err != nil {
				logger.Error().Err(err).Str("Container", containerName).Msg("failed to create container log file")
				return err
			}
			// Parse and write logs
			header := make([]byte, 8) // Docker stream header is 8 bytes
			for {
				_, err := io.ReadFull(logs, header)
				if err == io.EOF {
					break
				}
				if err != nil {
					logger.Error().Err(err).Str("Container", containerName).Msg("failed to read log stream header")
					break
				}

				// Extract log message size
				msgSize := binary.BigEndian.Uint32(header[4:8])

				// Read the log message
				msg := make([]byte, msgSize)
				_, err = io.ReadFull(logs, msg)
				if err != nil {
					logger.Error().Err(err).Str("Container", containerName).Msg("failed to read log message")
					break
				}

				// Write the log message to the file
				if _, err := logFile.Write(msg); err != nil {
					logger.Error().Err(err).Str("Container", containerName).Msg("failed to write log message to file")
					break
				}
			}
			return nil
		})
	}
	return eg.Wait()
}

func shouldIgnore(logger zerolog.Logger, containerName string) bool {
	if in(containerName, []string{"/sig-provider", "/stats", "/stats-db", "/db", "/backend", "/promtail", "/compose", "/blockscout", "/frontend", "/user-ops-indexer", "/visualizer", "/redis-db", "/proxy"}) {
		logger.Debug().Str("Container", containerName).Msg("Ignoring local tool container output")
		return true
	}
	return false
}

func in(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
