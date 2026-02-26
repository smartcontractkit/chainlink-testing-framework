package chaos

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	// Docker and docker-tc commands
	CmdPause     = "pause"
	CmdDelay     = "delay"
	CmdLoss      = "loss"
	CmdDuplicate = "duplicate"
	CmdCorrupt   = "corrupt"
)

const (
	// dockerTCContainerName default "docker-tc" container name
	dockerTCContainerName = "dtc"
	// dockerTCInternalSvc docker-tc internal service name
	dockerTCInternalSvc = "localhost:4080"
)

var (
	defaultCURLCMD = fmt.Sprintf("docker exec %s curl", dockerTCContainerName)
	tcCommands     = []string{CmdDelay, CmdLoss, CmdCorrupt, CmdDuplicate}
)

// DockerChaos is a chaos generator for Docker
type DockerChaos struct {
	Experiments map[string]string
}

// NewDockerChaos creates a new "docker-tc" instance or reuses existing one
func NewDockerChaos(ctx context.Context) (*DockerChaos, error) {
	framework.L.Info().
		Str("Container", dockerTCContainerName).
		Msg("Starting new docker-tc container")

	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	req := testcontainers.ContainerRequest{
		Image:      "lukaszlach/docker-tc",
		Name:       dockerTCContainerName,
		CapAdd:     []string{"NET_ADMIN"},
		WaitingFor: wait.ForLog("Starting Docker Traffic Control"),
		HostConfigModifier: func(h *container.HostConfig) {
			h.Privileged = true
			h.NetworkMode = "host"
			h.Mounts = []mount.Mount{
				{
					Type:     "bind",
					Source:   "/var/run/docker.sock",
					Target:   "/var/run/docker.sock",
					ReadOnly: true,
				},
			}
		},
	}
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start docker-tc container: %w", err)
	}
	return &DockerChaos{
		Experiments: make(map[string]string, 0),
	}, nil
}

// RemoveAll removes all the experiments
func (m *DockerChaos) RemoveAll() error {
	for exName, exCmd := range m.Experiments {
		if _, err := framework.ExecCmd(exCmd); err != nil {
			return fmt.Errorf("failed to remove chaos experiment: name: %s, command:%s, err: %w", exName, exCmd, err)
		}
	}
	return nil
}

// Chaos executes either Docker or "docker-tc" commands
func (m *DockerChaos) Chaos(containerName string, cmd, val string) error {
	if slices.Contains(tcCommands, cmd) {
		m.Experiments[containerName] = fmt.Sprintf("%s -X DELETE %s/%s", defaultCURLCMD, dockerTCInternalSvc, containerName)
		if _, err := framework.ExecCmd(fmt.Sprintf("%s -d %s=%s %s/%s", defaultCURLCMD, cmd, val, dockerTCInternalSvc, containerName)); err != nil {
			return err
		}
	} else {
		m.Experiments[containerName] = fmt.Sprintf("docker unpause %s", containerName)
		if _, err := framework.ExecCmd(fmt.Sprintf("docker pause %s", containerName)); err != nil {
			return err
		}
	}
	return nil
}

// DEPRECATED: Since Pumba has outdated Docker dependencies it may not work without additional
// setting to allow using Docker client which is out of client<>server compatibility range.
// Use NewDockerChaos for pause and network experiments!
//
// ExecPumba executes Pumba (https://github.com/alexei-led/pumba) command
// since handling various docker race conditions is hard and there is no easy API for that
// for now you can provide time to wait until chaos is applied
func ExecPumba(command string, wait time.Duration) (func(), error) {
	ctx := context.Background()
	cmd := strings.Split(command, " ")
	pumbaReq := testcontainers.ContainerRequest{
		Name:       fmt.Sprintf("chaos-%s", uuid.NewString()[0:5]),
		Image:      "gaiaadm/pumba",
		Labels:     framework.DefaultTCLabels(),
		Privileged: true,
		Cmd:        cmd,
		HostConfigModifier: func(h *container.HostConfig) {
			h.Mounts = []mount.Mount{
				{
					Type:     "bind",
					Source:   "/var/run/docker.sock",
					Target:   "/var/run/docker.sock",
					ReadOnly: true,
				},
			}
		},
	}
	pumbaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: pumbaReq,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start pumba chaos container: %w", err)
	}
	framework.L.Info().Str("Cmd", command).Msg("Pumba chaos has started")
	time.Sleep(wait)
	framework.L.Info().Msg("Pumba chaos has finished")
	return func() {
		_ = pumbaContainer.Terminate(ctx)
	}, nil
}
