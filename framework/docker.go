package framework

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
)

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

func safeContainerName(info container.Summary) string {
	if len(info.Names) > 0 {
		name := strings.TrimPrefix(info.Names[0], "/")
		if name != "" {
			// defensive: docker names normally don't include "/" beyond prefix,
			// but this guarantees safe map keys and filenames.
			return strings.ReplaceAll(name, "/", "_")
		}
	}
	// fallback when Names is missing/unexpected
	if len(info.ID) >= 12 {
		return info.ID[:12]
	}
	return info.ID
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
