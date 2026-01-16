package leak_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/leak"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRunHog(t *testing.T) {
	ctx := context.Background()
	hog, err := SetupResourceHog(
		ctx,
		"resource-hog:latest",
		map[string]string{
			"WORKERS": "1,2,3,2,1",
			"MEMORY":  "1,2,3,2,1",
			"REPEAT":  "1",
		},
	)
	require.NoError(t, err)
	time.Sleep(15 * time.Minute)
	t.Cleanup(func() {
		if err := hog.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	})
}

func TestVerifyHog(t *testing.T) {
	lc := leak.NewResourceLeakChecker()
	// cpu
	diff, err := lc.MeasureDelta(&leak.CheckConfig{
		Query: `sum(rate(container_cpu_usage_seconds_total{name="resource-hog"}[5m])) * 100`,
		Start: mustTime("2026-01-16T13:20:30Z"),
		End:   mustTime("2026-01-16T13:39:45Z"),
	})
	fmt.Println(diff)
	require.NoError(t, err)

	// mem
	diff, err = lc.MeasureDelta(&leak.CheckConfig{
		Query: `avg_over_time(container_memory_rss{name="resource-hog"}[5m]) / 1024 / 1024`,
		Start: mustTime("2026-01-16T13:20:30Z"),
		End:   mustTime("2026-01-16T13:38:25Z"),
	})
	fmt.Println(diff)
	require.NoError(t, err)
}

// ResourceHogContainer represents a container that hogs CPU and memory
type ResourceHogContainer struct {
	testcontainers.Container
	URI string
}

// SetupResourceHog starts a container that consumes CPU and memory
func SetupResourceHog(ctx context.Context, image string, env map[string]string) (*ResourceHogContainer, error) {
	// Build request for the container
	req := testcontainers.ContainerRequest{
		Name:         "resource-hog",
		Image:        image,
		ExposedPorts: []string{},
		Env:          env,
		WaitingFor:   wait.ForLog("Starting CPU and Memory hog"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	return &ResourceHogContainer{Container: container}, nil
}
