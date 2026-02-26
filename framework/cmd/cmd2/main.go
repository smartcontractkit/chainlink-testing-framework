package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
)

func main() {
	ctx := context.Background()

	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	// Define the container request
	req := testcontainers.ContainerRequest{
		Image:      "lukaszlach/docker-tc",
		Name:       "dtc",
		AutoRemove: false,
		CapAdd:     []string{"NET_ADMIN"},
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
				{
					Type:   "bind",
					Source: "/var/docker-tc",
					Target: "/var/docker-tc",
				},
			}
		},
	}

	// Create the container
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %s", err)
	}
	time.Sleep(15 * time.Second)
	if _, err := framework.ExecCmd("docker exec dtc curl -d delay=8000ms localhost:4080/blockchain-node-2baf2"); err != nil {
		panic(err)
	}
	time.Sleep(30 * time.Second)
	if _, err := framework.ExecCmd(`docker exec dtc curl -X DELETE localhost:4080/blockchain-node-2baf2`); err != nil {
		panic(err)
	}
}
