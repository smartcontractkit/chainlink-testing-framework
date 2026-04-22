package framework

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/moby/moby/client"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"
)

const (
	EnvVarIgnoreCriticalLogs = "CTF_IGNORE_CRITICAL_LOGS"
	DefaultCTFLogsDir        = "logs/docker"
)

var criticalLogLevelRegex = regexp.MustCompile(`(CRIT|PANIC|FATAL)`)

func checkNodeLogStream(source string, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	// safer for long log lines
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if criticalLogLevelRegex.MatchString(line) {
			return fmt.Errorf("source %s contains matching log level at line %d: %s", source, lineNumber, line)
		}
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading source %s: %w", source, err)
	}
	return nil
}

// New stream-first API
func checkNodeLogErrorsFromStreams(streams map[string]io.ReadCloser) error {
	if os.Getenv(EnvVarIgnoreCriticalLogs) == "true" {
		L.Warn().Msg(`CTF_IGNORE_CRITICAL_LOGS is set to true, we ignore all CRIT|FATAL|PANIC errors in node logs!`)
		return nil
	}
	for name, rc := range streams {
		if err := checkNodeLogStream(name, rc); err != nil {
			_ = rc.Close()
			return err
		}
		_ = rc.Close()
	}
	return nil
}

func StreamContainerLogs(listOptions client.ContainerListOptions, logOptions client.ContainerLogsOptions) (map[string]io.ReadCloser, error) {
	L.Info().Msg("Streaming Docker containers logs")
	provider, err := tc.NewDockerProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker provider: %w", err)
	}
	containers, err := provider.Client().ContainerList(context.Background(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker containers: %w", err)
	}

	eg := &errgroup.Group{}
	logMap := make(map[string]io.ReadCloser)
	var mutex sync.Mutex

	for _, containerInfo := range containers.Items {
		eg.Go(func() error {
			containerName := safeContainerName(containerInfo)
			L.Debug().Str("Container", containerName).Msg("Collecting logs")
			logs, err := provider.Client().ContainerLogs(context.Background(), containerInfo.ID, logOptions)
			if err != nil {
				L.Error().Err(err).Str("Container", containerName).Msg("failed to fetch logs for container")
				return err
			}
			mutex.Lock()
			defer mutex.Unlock()
			logMap[containerName] = logs
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return logMap, nil
}

func CTFContainersListOpts() client.ContainerListOptions {
	return client.ContainerListOptions{
		All: true,
		Filters: make(client.Filters).Add("label", "framework=ctf"),
	}
}

func CTFContainersLogsOpts() client.ContainerLogsOptions {
	return client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
}

// LogStreamConsumer represents a log stream consumer that receives one stream per container.
type LogStreamConsumer struct {
	Name    string
	Consume func(map[string]io.ReadCloser) error
}

// StreamContainerLogsFanout fetches container logs once and fans out streams to all consumers.
func StreamContainerLogsFanout(listOptions client.ContainerListOptions, logOptions client.ContainerLogsOptions, consumers ...LogStreamConsumer) error {
	logStream, err := StreamContainerLogs(listOptions, logOptions)
	if err != nil {
		return err
	}

	return fanoutContainerLogs(logStream, consumers...)
}

// StreamCTFContainerLogsFanout fetches CTF logs once and fans out streams to all consumers.
func StreamCTFContainerLogsFanout(consumers ...LogStreamConsumer) error {
	return StreamContainerLogsFanout(CTFContainersListOpts(), CTFContainersLogsOpts(), consumers...)
}

func fanoutContainerLogs(logStream map[string]io.ReadCloser, consumers ...LogStreamConsumer) error {
	if len(consumers) == 0 {
		for _, reader := range logStream {
			_ = reader.Close()
		}
		return nil
	}

	consumerStreams := make([]map[string]io.ReadCloser, len(consumers))
	for i := range consumers {
		consumerStreams[i] = make(map[string]io.ReadCloser, len(logStream))
	}

	pumpGroup := &errgroup.Group{}
	for containerName, sourceReader := range logStream {
		writers := make([]*io.PipeWriter, len(consumers))
		for i := range consumers {
			reader, writer := io.Pipe()
			consumerStreams[i][containerName] = reader
			writers[i] = writer
		}

		pumpGroup.Go(func() error {
			defer func() {
				_ = sourceReader.Close()
			}()

			readBuf := make([]byte, 32*1024)
			for {
				n, readErr := sourceReader.Read(readBuf)
				if n > 0 {
					// NOTE: io.Pipe is unbuffered, so a slow consumer can still backpressure others.
					// Future improvement: decouple per-consumer delivery with bounded buffering.
					chunk := readBuf[:n]
					for i, writer := range writers {
						if writer == nil {
							continue
						}
						if writeErr := writeAll(writer, chunk); writeErr != nil {
							if errors.Is(writeErr, io.ErrClosedPipe) {
								writers[i] = nil
								continue
							}
							closeAllPipeWritersWithError(writers, fmt.Errorf("failed writing stream for container %s: %w", containerName, writeErr))
							return fmt.Errorf("failed writing stream for container %s: %w", containerName, writeErr)
						}
					}
				}

				if readErr == io.EOF {
					closeAllPipeWriters(writers)
					return nil
				}
				if readErr != nil {
					closeAllPipeWritersWithError(writers, readErr)
					return fmt.Errorf("failed reading stream for container %s: %w", containerName, readErr)
				}
			}
		})
	}

	consumerGroup := &errgroup.Group{}
	for i, consumer := range consumers {
		consumerGroup.Go(func() error {
			defer closeAllPipeReaders(consumerStreams[i])

			if consumer.Consume == nil {
				return fmt.Errorf("consumer %q has nil Consume function", consumer.Name)
			}
			if err := consumer.Consume(consumerStreams[i]); err != nil {
				return fmt.Errorf("consumer %q failed: %w", consumer.Name, err)
			}
			return nil
		})
	}

	pumpErr := pumpGroup.Wait()
	consumerErr := consumerGroup.Wait()
	return errors.Join(pumpErr, consumerErr)
}

func writeDockerLogPayload(dst io.Writer, reader io.Reader) error {
	header := make([]byte, 8)
	for {
		_, err := io.ReadFull(reader, header)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read log stream header: %w", err)
		}

		msgSize := binary.BigEndian.Uint32(header[4:8])
		msg := make([]byte, msgSize)
		if _, err = io.ReadFull(reader, msg); err != nil {
			return fmt.Errorf("failed to read log message: %w", err)
		}
		if _, err = dst.Write(msg); err != nil {
			return fmt.Errorf("failed to write log message: %w", err)
		}
	}
}

func writeAll(writer io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := writer.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}

func closeAllPipeWriters(writers []*io.PipeWriter) {
	for _, writer := range writers {
		if writer == nil {
			continue
		}
		_ = writer.Close()
	}
}

func closeAllPipeWritersWithError(writers []*io.PipeWriter, err error) {
	for _, writer := range writers {
		if writer == nil {
			continue
		}
		_ = writer.CloseWithError(err)
	}
}

func closeAllPipeReaders(readers map[string]io.ReadCloser) {
	for _, reader := range readers {
		if reader == nil {
			continue
		}
		_ = reader.Close()
	}
}

func SaveAndCheckLogs(t *testing.T) error {
	return StreamCTFContainerLogsFanout(
		LogStreamConsumer{
			Name: "save-container-logs",
			Consume: func(logStreams map[string]io.ReadCloser) error {
				_, logsErr := SaveContainerLogsFromStreams(fmt.Sprintf("%s-%s", DefaultCTFLogsDir, t.Name()), logStreams)
				return logsErr
			},
		},
		LogStreamConsumer{
			Name: "check-container-logs",
			Consume: func(logStreams map[string]io.ReadCloser) error {
				checkErr := checkNodeLogErrorsFromStreams(logStreams)
				return checkErr
			},
		},
	)
}

// SaveContainerLogs writes all Docker container logs to some directory
func SaveContainerLogs(dir string) ([]string, error) {
	logStream, lErr := StreamContainerLogs(CTFContainersListOpts(), CTFContainersLogsOpts())

	if lErr != nil {
		return nil, lErr
	}

	return SaveContainerLogsFromStreams(dir, logStream)
}

// SaveContainerLogsFromStreams writes all provided Docker log streams to files in some directory.
func SaveContainerLogsFromStreams(dir string, logStream map[string]io.ReadCloser) ([]string, error) {
	L.Info().Msg("Writing Docker containers logs")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	eg := &errgroup.Group{}
	logFilePaths := make([]string, 0)
	var logFilePathsMu sync.Mutex
	for containerName, reader := range logStream {
		eg.Go(func() error {
			defer func() {
				_ = reader.Close()
			}()

			logFilePath := filepath.Join(dir, fmt.Sprintf("%s.log", containerName))
			logFile, err := os.Create(logFilePath)
			if err != nil {
				L.Error().Err(err).Str("Container", containerName).Msg("failed to create container log file")
				return err
			}
			defer func() {
				_ = logFile.Close()
			}()

			logFilePathsMu.Lock()
			logFilePaths = append(logFilePaths, logFilePath)
			logFilePathsMu.Unlock()

			if err := writeDockerLogPayload(logFile, reader); err != nil {
				L.Error().Err(err).Str("Container", containerName).Msg("failed to write container logs")
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return logFilePaths, nil
}

var ExitedCtfContainersListOpts = client.ContainerListOptions{
	All: true,
	Filters: make(client.Filters).Add("label", "framework=ctf").Add("status", "exited", "dead"),
}

// PrintFailedContainerLogs writes exited/dead CTF containers' last log lines to stdout.
func PrintFailedContainerLogs(logLinesCount uint64) error {
	logStream, lErr := StreamContainerLogs(ExitedCtfContainersListOpts, client.ContainerLogsOptions{
		ShowStderr: true,
		Tail:       strconv.FormatUint(logLinesCount, 10),
	})
	if lErr != nil {
		return lErr
	}

	return PrintFailedContainerLogsFromStreams(logStream, logLinesCount)
}

// PrintFailedContainerLogsFromStreams prints all provided container streams as red text.
func PrintFailedContainerLogsFromStreams(logStream map[string]io.ReadCloser, logLinesCount uint64) error {
	if len(logStream) == 0 {
		L.Info().Msg("No failed Docker containers found")
		return nil
	}

	L.Error().Msgf("Exited/dead containers: %s", strings.Join(slices.Collect(maps.Keys(logStream)), ", "))

	eg := &errgroup.Group{}
	for cName, ioReader := range logStream {
		eg.Go(func() error {
			defer func() {
				_ = ioReader.Close()
			}()

			var content strings.Builder
			if err := writeDockerLogPayload(&content, ioReader); err != nil {
				return fmt.Errorf("failed to read logs for container %s: %w", cName, err)
			}

			trimmed := strings.TrimSpace(content.String())
			if len(trimmed) > 0 {
				L.Info().Str("Container", cName).Msgf("Last %d lines of logs", logLinesCount)
				fmt.Println(RedText("%s\n", trimmed))
			}

			return nil
		})
	}

	return eg.Wait()
}

func CheckContainersForPanics(maxLinesAfterPanic int) bool {
	logstream, err := StreamContainerLogs(CTFContainersListOpts(), CTFContainersLogsOpts())
	if err != nil {
		L.Error().Err(err).Msg("failed to stream container logs for panic check")
		return true
	}

	return CheckContainersForPanicsFromStreams(logstream, maxLinesAfterPanic)
}

// CheckContainersForPanicsFromStreams scans the provided stream (usually Docker container logs) for panic-related patterns.
// When a panic is detected, it displays the panic line and up to maxLinesAfterPanic lines following it.
//
// This is useful for debugging test failures where a Docker container may have panicked.
// The function searches for common panic patterns including:
//   - "panic:" - Standard Go panic
//   - "runtime error:" - Runtime errors (nil pointer, index out of bounds, etc.)
//   - "fatal error:" - Fatal errors
//   - "goroutine N [running]" - Stack trace indicators
//
// The function scans all containers in parallel and stops as soon as the first panic is found.
//
// Parameters:
//   - logStream: Map of container names to their log streams (io.ReadCloser).
//   - maxLinesAfterPanic: Maximum number of lines to show after the panic line (including stack trace).
//     Recommended: 50-200 depending on how much context you want.
//
// Returns:
//   - true if any panics were found in any container, false otherwise
func CheckContainersForPanicsFromStreams(logStream map[string]io.ReadCloser, maxLinesAfterPanic int) bool {
	// Panic patterns to search for
	panicPatterns := map[string]*regexp.Regexp{
		"panic":         regexp.MustCompile(`(?i)panic:`),                // Go panic
		"runtime error": regexp.MustCompile(`(?i)runtime error:`),        // Runtime errors
		"fatal error":   regexp.MustCompile(`(?i)fatal error:`),          // Fatal errors
		"stack trace":   regexp.MustCompile(`goroutine \d+ \[running\]`), // Stack trace indicator
	}

	// Create context for early cancellation when first panic is found
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to receive panic results
	panicFoundChan := make(chan bool, len(logStream))
	var wg sync.WaitGroup

	// Scan all containers in parallel
	for containerName, reader := range logStream {
		wg.Add(1)
		go func(name string, r io.ReadCloser) {
			defer wg.Done()
			defer r.Close()

			panicFound := scanContainerForPanics(ctx, L, name, r, panicPatterns, maxLinesAfterPanic)
			if panicFound {
				panicFoundChan <- true
				cancel() // Signal other goroutines to stop
			}
		}(containerName, reader)
	}

	// Wait for all goroutines to finish in a separate goroutine
	go func() {
		wg.Wait()
		close(panicFoundChan)
	}()

	// Check if any panic was found
	for range panicFoundChan {
		return true // Return as soon as first panic is found
	}

	L.Info().Msg("No panics detected in any container logs")

	return false
}

// scanContainerForPanics scans a single container's log stream for panic patterns
// It checks the context for cancellation to enable early termination when another goroutine finds a panic
func scanContainerForPanics(ctx context.Context, logger zerolog.Logger, containerName string, reader io.Reader, patterns map[string]*regexp.Regexp, maxLinesAfter int) bool {
	scanner := bufio.NewScanner(reader)
	// Increase buffer size to handle large log lines
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var logLines []string
	lineNum := 0
	panicLineNum := -1
	patternNameFound := ""

	// Read all lines and detect panic
	for scanner.Scan() {
		// Check if context is cancelled (another goroutine found a panic)
		select {
		case <-ctx.Done():
			return false // Stop scanning, another container already found a panic
		default:
		}

		line := scanner.Text()
		lineNum++

		// If we found a panic, collect remaining lines up to maxLinesAfter
		if panicLineNum >= 0 {
			logLines = append(logLines, line)
			// Stop reading once we have enough context after the panic
			if lineNum >= panicLineNum+maxLinesAfter+1 {
				break
			}
			continue
		}

		// Still searching for panic - store all lines
		logLines = append(logLines, line)

		// Check if this line contains a panic pattern
		for patternName, pattern := range patterns {
			if pattern.MatchString(line) {
				patternNameFound = patternName
				panicLineNum = lineNum - 1 // Store index (0-based)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error().Err(err).Str("Container", containerName).Msg("error reading container logs")
		return false
	}

	// If panic was found, display it with context
	if panicLineNum >= 0 {
		logger.Error().
			Str("Container", containerName).
			Int("PanicLineNumber", panicLineNum+1).
			Msgf("🔥 %s DETECTED in container logs", strings.ToUpper(patternNameFound))

		// Calculate range to display
		startLine := panicLineNum
		endLine := min(len(logLines), panicLineNum+maxLinesAfter+1)

		// Build the output
		var output strings.Builder
		fmt.Fprintf(&output, "\n%s\n", strings.Repeat("=", 80))
		fmt.Fprintf(&output, "%s FOUND IN CONTAINER: %s (showing %d lines from panic)\n", strings.ToUpper(patternNameFound), containerName, endLine-startLine)
		fmt.Fprintf(&output, "%s\n", strings.Repeat("=", 80))

		for i := startLine; i < endLine; i++ {
			fmt.Fprintf(&output, "%s\n", logLines[i])
		}

		fmt.Fprintf(&output, "%s\n", strings.Repeat("=", 80))

		fmt.Println(RedText("%s\n", output.String()))
		return true
	}

	return false
}
