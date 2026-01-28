package framework

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

/* Templates */

const (
	ProductComponentImpl = `
	package services

/*
 * A simple template to add your project services to devenv.
 * Each service should have a file with Input/Output struct and a function deploying it.
 */

import (
		"context"
		"fmt"
		"os"
		"strconv"

		"github.com/rs/zerolog"
		"github.com/rs/zerolog/log"

		"github.com/docker/docker/api/types/container"
		"github.com/docker/go-connections/nat"
		"github.com/smartcontractkit/chainlink-testing-framework/framework"
		"github.com/testcontainers/testcontainers-go"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "datastreams"}).Logger()

const (
		DefaultExampleSvcName  = "example-service"
		DefaultExampleSvcImage = "busybox:latest"
)

// ExampleSvcInput example service input
type ExampleSvcInput struct {
		Image         string            ` + "`" + `toml:"image"` + "`" + `
		Port          int               ` + "`" + `toml:"port"` + "`" + `
		ContainerName string            ` + "`" + `toml:"container_name"` + "`" + `
		EnvVars       map[string]string ` + "`" + `toml:"env_vars"` + "`" + `
		Out           *ExampleSvcOutput ` + "`" + `toml:"out"` + "`" + `
}

// Default is a default input configuration
func (s *ExampleSvcInput) Default() {
		if s.Image == "" {
			s.Image = DefaultExampleSvcImage
		}
		if s.Port == 0 {
			s.Port = 9501
		}
		if s.ContainerName == "" {
			s.ContainerName = DefaultExampleSvcName
		}
}

// ExampleSvcOutput represents service connection details which can be consumed by tests
// or environment.go, or any other metadata your need to expose to clients
type ExampleSvcOutput struct {
		ServiceURL string ` + "`" + `toml:"service_url"` + "`" + `
}

// NewService deploys an example service and populates output data
func NewService(in *ExampleSvcInput) error {
		if in == nil || in.Out != nil {
			// either service key is not present in configuration
			// or it is already deployed because we have an output, skipping
			L.Info().Str("ServiceName", DefaultExampleSvcName).Msg("service is skipped or already deployed")
			return nil
		}
		ctx := context.Background()
		// read and apply default inputs
		in.Default()
		// create your service container, use static ports
		req := testcontainers.ContainerRequest{
			Image:    in.Image,
			Name:     in.ContainerName,
			Labels:   framework.DefaultTCLabels(),
			Networks: []string{framework.DefaultNetworkName},
			NetworkAliases: map[string][]string{
				framework.DefaultNetworkName: {in.ContainerName},
			},
			Env: in.EnvVars,
			Cmd: []string{"sleep", "infinity"},
			// add more internal ports here with /tcp suffix, ex.: 9501/tcp
			ExposedPorts: []string{"9501/tcp"},
			HostConfigModifier: func(h *container.HostConfig) {
				h.PortBindings = nat.PortMap{
					// add more internal/external pairs here, ex.: 9501/tcp as a key and HostPort is the exposed port (no /tcp prefix!)
					"9501/tcp": []nat.PortBinding{
						{HostPort: strconv.Itoa(in.Port)},
					},
				}
			},
			// for complex services wait for specific log message
			// WaitingFor: wait.ForLog("Some log message").
			// 	WithStartupTimeout(120 * time.Second).
			// 	WithPollInterval(3 * time.Second),
		}
		_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			return fmt.Errorf("failed to setup %s service", in.ContainerName)
		}

		// write outputs on a shared struct
		in.Out = &ExampleSvcOutput{
			ServiceURL: "https://example.com",
		}
		return nil
}
	`
)

type SvcImplParams struct {
	ProductName string
}

func (g *EnvCodegen) GenerateServiceImpl() (string, error) {
	log.Info().Msg("Generating service implementation")
	p := SvcImplParams{}
	return render(ProductComponentImpl, p)
}

// WriteFakes writes all files related to fake services used in testing
func (g *EnvCodegen) WriteServices() error {
	servicesRoot := filepath.Join(g.cfg.outputDir, "services")
	if err := os.MkdirAll( //nolint:gosec
		servicesRoot,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to create services directory: %w", err)
	}

	// generate Docker wrapper for a service
	serviceImplContents, err := g.GenerateServiceImpl()
	if err != nil {
		return err
	}
	if err := os.WriteFile( //nolint:gosec
		filepath.Join(servicesRoot, "svc.go"),
		[]byte(serviceImplContents),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return nil
}
