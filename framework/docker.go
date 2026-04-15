package framework

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/docker/docker/api/types/container"
	dfilter "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"
)

const (
	DefaultCTFLogsDir = "logs/docker"
)

func CTFContainersListOpts() container.ListOptions {
	return container.ListOptions{
		All: true,
		Filters: dfilter.NewArgs(dfilter.KeyValuePair{
			Key:   "label",
			Value: "framework=ctf",
		}),
	}
}

func CTFContainersLogsOpts() container.LogsOptions {
	return container.LogsOptions{ShowStdout: true, ShowStderr: true}
}

func IsDockerRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Ping(context.Background())
	return err == nil
}

func GetHost(container tc.Container) (string, error) {
	return GetHostWithContext(context.Background(), container)
}

func GetHostWithContext(ctx context.Context, container tc.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}
	// if localhost then force it to ipv4 localhost
	if host == "localhost" {
		host = "127.0.0.1"
	}
	return host, nil
}

func MapTheSamePort(ports ...string) nat.PortMap {
	portMap := nat.PortMap{}
	for _, port := range ports {
		// need to split off /tcp or /udp
		onlyPort := strings.SplitN(port, "/", 2)
		portMap[nat.Port(port)] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: onlyPort[0],
			},
		}
	}
	return portMap
}

func DefaultTCLabels() map[string]string {
	return map[string]string{
		"framework": "ctf",
		"logging":   "promtail",
	}
}

func DefaultTCName(name string) string {
	return fmt.Sprintf("%s-%s", name, uuid.NewString()[0:5])
}

// RunCommand executes a command and prints the output.
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandDir executes a command in some directory and prints the output
func RunCommandDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

// DockerClient wraps a Docker API client and provides convenience methods
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient creates a new instance of DockerClient
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return &DockerClient{cli: cli}, nil
}

// ExecContainer executes a command inside a running container by name and returns the combined stdout/stderr.
func (dc *DockerClient) ExecContainer(containerName string, command []string) (string, error) {
	return dc.ExecContainerWithContext(context.Background(), containerName, command)
}

// ExecContainerWithContext executes a command inside a running container by name and returns the combined stdout/stderr.
func (dc *DockerClient) ExecContainerWithContext(ctx context.Context, containerName string, command []string) (string, error) {
	execConfig := container.ExecOptions{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	}

	return dc.ExecContainerOptionsWithContext(ctx, containerName, execConfig)
}

// ExecContainer executes a command inside a running container by name and returns the combined stdout/stderr.
func (dc *DockerClient) ExecContainerOptions(containerName string, execConfig container.ExecOptions) (string, error) {
	return dc.ExecContainerOptionsWithContext(context.Background(), containerName, execConfig)
}

// ExecContainerOptionsWithContext executes a command inside a running container by name and returns the combined stdout/stderr.
func (dc *DockerClient) ExecContainerOptionsWithContext(ctx context.Context, containerName string, execConfig container.ExecOptions) (string, error) {
	L.Info().Strs("Command", execConfig.Cmd).Str("ContainerName", containerName).Msg("Executing command")
	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}
	var containerID string
	for _, cont := range containers {
		for _, name := range cont.Names {
			if name == "/"+containerName {
				containerID = cont.ID
				break
			}
		}
	}
	if containerID == "" {
		return "", fmt.Errorf("container with name '%s' not found", containerName)
	}

	execID, err := dc.cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %w", err)
	}
	resp, err := dc.cli.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer resp.Close()
	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %w", err)
	}
	L.Info().Str("Output", string(output)).Msg("Command output")
	return string(output), nil
}

// CopyFile copies a file into a container by name
func (dc *DockerClient) CopyFile(containerName, sourceFile, targetPath string) error {
	ctx := context.Background()
	containerID, err := dc.findContainerIDByName(ctx, containerName)
	if err != nil {
		return fmt.Errorf("failed to find container ID by name: %s", containerName)
	}
	return dc.copyToContainer(containerID, sourceFile, targetPath)
}

// findContainerIDByName finds a container ID by its name
func (dc *DockerClient) findContainerIDByName(ctx context.Context, containerName string) (string, error) {
	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				return c.ID, nil
			}
		}
	}
	return "", fmt.Errorf("container with name %s not found", containerName)
}

// copyToContainer copies a file into a container
func (dc *DockerClient) copyToContainer(containerID, sourceFile, targetPath string) error {
	ctx := context.Background()
	src, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer src.Close()

	// Create a tar archive containing the file
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("could not stat source file: %w", err)
	}

	// Add file to tar
	header := &tar.Header{
		Name: info.Name(),
		Size: info.Size(),
		Mode: int64(info.Mode()),
	}
	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("could not write tar header: %w", err)
	}
	if _, err := io.Copy(tw, src); err != nil {
		return fmt.Errorf("could not write file to tar archive: %w", err)
	}
	tw.Close()

	// Copy the tar archive to the container
	err = dc.cli.CopyToContainer(ctx, containerID, targetPath, &buf, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if err != nil {
		return fmt.Errorf("could not copy file to container: %w", err)
	}
	return nil
}

// SearchLogFile searches logfile using regex and return matches or error
func SearchLogFile(fp string, regex string) ([]string, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	matches := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			L.Info().Str("Regex", regex).Msg("Log match found")
			matches = append(matches, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return matches, err
	}
	return matches, nil
}

func SaveAndCheckLogs(t *testing.T) error {
	_, err := SaveContainerLogs(fmt.Sprintf("%s-%s", DefaultCTFLogsDir, t.Name()))
	if err != nil {
		return err
	}
	err = CheckCLNodeContainerErrors()
	if err != nil {
		return err
	}
	return nil
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

// PrintFailedContainerLogs writes exited/dead CTF containers' last log lines to stdout.
func PrintFailedContainerLogs(logLinesCount uint64) error {
	logStream, lErr := StreamContainerLogs(ExitedCtfContainersListOpts, container.LogsOptions{
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

	L.Error().Msgf("Containers that exited with non-zero codes: %s", strings.Join(slices.Collect(maps.Keys(logStream)), ", "))

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

var ExitedCtfContainersListOpts = container.ListOptions{
	All: true,
	Filters: dfilter.NewArgs(dfilter.KeyValuePair{
		Key:   "label",
		Value: "framework=ctf",
	},
		dfilter.KeyValuePair{
			Key:   "status",
			Value: "exited"},
		dfilter.KeyValuePair{
			Key:   "status",
			Value: "dead"}),
}

func StreamContainerLogs(listOptions container.ListOptions, logOptions container.LogsOptions) (map[string]io.ReadCloser, error) {
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

	for _, containerInfo := range containers {
		eg.Go(func() error {
			containerName := containerInfo.Names[0]
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

// LogStreamConsumer represents a log stream consumer that receives one stream per container.
type LogStreamConsumer struct {
	Name    string
	Consume func(map[string]io.ReadCloser) error
}

// StreamContainerLogsFanout fetches container logs once and fans out streams to all consumers.
func StreamContainerLogsFanout(listOptions container.ListOptions, logOptions container.LogsOptions, consumers ...LogStreamConsumer) error {
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
		i := i
		consumer := consumer
		consumerGroup.Go(func() error {
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

func BuildImageOnce(once *sync.Once, dctx, dfile, nameAndTag string, buildArgs map[string]string) error {
	var err error
	once.Do(func() {
		err = BuildImage(dctx, dfile, nameAndTag, buildArgs)
		if err != nil {
			err = fmt.Errorf("failed to build Docker image: %w", err)
		}
	})
	return err
}

func BuildImage(dctx, dfile, nameAndTag string, buildArgs map[string]string) error {
	dfilePath := filepath.Join(dctx, dfile)

	if os.Getenv("CTF_CLNODE_DLV") == "true" {
		commandParts := []string{"docker", "buildx", "build", "--build-arg", `GO_GCFLAGS=all=-N -l`, "--build-arg", "CHAINLINK_USER=chainlink"}
		for k, v := range buildArgs {
			commandParts = append(commandParts, "--build-arg", fmt.Sprintf("%s=%s", k, v))
		}
		if os.Getenv("GITHUB_TOKEN") != "" {
			commandParts = append(commandParts, "--secret", "id=GIT_AUTH_TOKEN,env=GITHUB_TOKEN")
		}
		commandParts = append(commandParts, "--load", "-t", nameAndTag, "-f", dfilePath, dctx)
		return RunCommand(commandParts[0], commandParts[1:]...)
	}
	commandParts := []string{"docker", "buildx", "build", "--build-arg", "CHAINLINK_USER=chainlink"}
	for k, v := range buildArgs {
		commandParts = append(commandParts, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}
	if os.Getenv("GITHUB_TOKEN") != "" {
		commandParts = append(commandParts, "--secret", "id=GIT_AUTH_TOKEN,env=GITHUB_TOKEN")
	}
	commandParts = append(commandParts, "--load", "-t", nameAndTag, "-f", dfilePath, dctx)
	return RunCommand(commandParts[0], commandParts[1:]...)
}

// RemoveTestContainers removes all test containers, volumes and CTF docker network
func RemoveTestContainers() error {
	L.Info().Str("label", "framework=ctf").Msg("Cleaning up docker containers")
	// Bash command for removing Docker containers and networks with "framework=ctf" label
	cmd := exec.Command("bash", "-c", `
		docker ps -aq --filter "label=framework=ctf" | xargs -r docker rm -f && \
		docker volume ls -q | xargs -r docker volume rm || true
	`)
	L.Debug().Msg("Running command")
	if L.GetLevel() == zerolog.DebugLevel {
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}

func RemoveTestStack(name string) error {
	L.Info().Str("stack name", name).Msg("Cleaning up docker containers")
	// Bash command for removing Docker containers and networks with "framework=ctf" label
	//nolint:gosec //ignoring G204
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`
		docker ps -a --filter "label=com.docker.compose.project" --format '{{.ID}} {{.Label "com.docker.compose.project"}}' \
		| awk '$2 ~ /^%s/ { print $2 }' | sort -u \
		| xargs -I{} docker compose -p {} down -v --remove-orphans
	`, name))
	L.Debug().Msg("Running command")
	if L.GetLevel() == zerolog.DebugLevel {
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running clean command: %s", string(output))
	}
	return nil
}

type ContainerResources struct {
	CPUs     float64 `toml:"cpus" validate:"gte=0" comment:"CPU shares, ex.: 2"`
	MemoryMb uint    `toml:"memory_mb" comment:"Memory in MegaBytes, ex.:\"200\""`
}

// ResourceLimitsFunc returns a function to configure container resources based on the human-readable CPUs and memory in Mb
func ResourceLimitsFunc(h *container.HostConfig, resources *ContainerResources) {
	if resources == nil {
		return
	}
	if resources.MemoryMb > 0 {
		//nolint:gosec
		h.Memory = int64(resources.MemoryMb) * 1024 * 1024 // Memory in Mb
		//nolint:gosec
		h.MemoryReservation = int64(resources.MemoryMb) * 1024 * 1024 // Total memory that can be reserved (soft) in Mb
		// https://docs.docker.com/engine/containers/resource_constraints/ if both values are equal swap is off, read the docs
		h.MemorySwap = h.Memory
	}
	if resources.CPUs > 0 {
		// Set CPU limits using CPUQuota and CPUPeriod
		// we don't use runtime.NumCPU or docker API to get CPUs because h.CPUShares is relative to amount of containers you run
		// CPUPeriod and CPUQuota are absolute and easier to control
		h.CPUPeriod = 100000                        // Default period (100ms)
		h.CPUQuota = int64(resources.CPUs * 100000) // Quota in microseconds (e.g., 0.5 CPUs = 50000)
	}
}

// GenerateCustomPortsData generate custom ports data: exposed and forwarded port map
func GenerateCustomPortsData(portsProvided []string) ([]string, nat.PortMap, error) {
	exposedPorts := make([]string, 0)
	portBindings := nat.PortMap{}
	customPorts := make([]string, 0)
	for _, p := range portsProvided {
		if !strings.Contains(p, ":") {
			return nil, nil, fmt.Errorf("custom ports must have format external_port:internal_port, you provided: %s", p)
		}
		pp := strings.Split(p, ":")
		if len(pp) != 2 {
			return nil, nil, fmt.Errorf("custom_ports has ':' but you must provide both ports, you provided: %s", pp)
		}
		customPorts = append(customPorts, fmt.Sprintf("%s/tcp", pp[1]))

		dockerPort := nat.Port(fmt.Sprintf("%s/tcp", pp[1]))
		hostPort := pp[0]
		portBindings[dockerPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		}
	}
	exposedPorts = append(exposedPorts, customPorts...)
	return exposedPorts, portBindings, nil
}

// NoDNS removes default DNS server and sets it to localhost
func NoDNS(noDNS bool, hc *container.HostConfig) {
	if noDNS {
		hc.DNS = []string{"127.0.0.1"}
	}
}
