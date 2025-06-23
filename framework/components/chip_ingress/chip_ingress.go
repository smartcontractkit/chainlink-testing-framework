package chipingress

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

type Output struct {
	InternalURL string
	ExternalURL string
}

func New(ctx context.Context) (*Output, error) {
	composeFile := "./docker-compose.yml"
	stack, err := compose.NewDockerCompose([]string{composeFile},
		compose.WithStackName("chip_stack"),
		compose.WithEnv(map[string]string{
			// optional: override any env vars here
			// "SERVER_METRICS_OTEL_EXPORTER_GRPC_ENDPOINT": "otel-lgtm:4317",
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compose stack")
	}

	// Add to existing Docker network
	stack = stack.WithAdditionalOptions(func(s *testcontainers.DockerCompose) {
		s.WithCustomNetwork("my_shared_net") // your existing network name
	})

	fmt.Println("Bringing up stack...")
	if err := stack.Up(ctx, compose.Wait(true)); err != nil {
		log.Fatalf("Failed to start stack: %v", err)
	}
	defer func() {
		_ = stack.Down(ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
	}()

	// Optional: wait for chip-ingress to become healthy
	time.Sleep(5 * time.Second)

	// Retrieve internal and external ports
	ingress, err := stack.ServiceContainer(ctx, "chip-ingress")
	if err != nil {
		log.Fatalf("Failed to get chip-ingress container: %v", err)
	}

	internalHost := "chip-ingress"
	internalPort := "50051" // from service definition
	externalHost := "localhost"
	externalPort, err := ingress.MappedPort(ctx, "50051")
	if err != nil {
		log.Fatalf("Failed to get mapped port: %v", err)
	}
}
