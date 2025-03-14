package chaos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

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
