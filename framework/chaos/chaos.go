package chaos

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	"strings"
)

func ExecPumba(command string) (func(), error) {
	ctx := context.Background()
	cmd := strings.Split(command, " ")
	pumbaReq := testcontainers.ContainerRequest{
		Name:       fmt.Sprintf("chaos-%s", uuid.NewString()[0:5]),
		Image:      "gaiaadm/pumba",
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
	framework.L.Info().Msg("Pumba chaos started")
	return func() {
		_ = pumbaContainer.Terminate(ctx)
	}, nil
}
