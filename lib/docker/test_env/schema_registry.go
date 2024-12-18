package test_env

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const defaultSchemaRegistryImage = "confluentinc/cp-schema-registry:7.4.0"

type SchemaRegistry struct {
	EnvComponent
	EnvVars     map[string]string
	InternalUrl string
	ExternalUrl string
	l           zerolog.Logger
	t           *testing.T
}

// NewSchemaRegistry initializes a new SchemaRegistry instance with a unique container name and specified networks.
// It sets default environment variables and a startup timeout, making it suitable for managing schema registries in decentralized applications.
func NewSchemaRegistry(networks []string) *SchemaRegistry {
	id, _ := uuid.NewRandom()
	defaultEnvVars := map[string]string{
		"SCHEMA_REGISTRY_DEBUG": "true",
	}
	return &SchemaRegistry{
		EnvComponent: EnvComponent{
			ContainerName:  fmt.Sprintf("schema-registry-%s", id.String()),
			Networks:       networks,
			StartupTimeout: 1 * time.Minute,
		},
		EnvVars: defaultEnvVars,
	}
}

// WithTestInstance sets up a SchemaRegistry instance for testing purposes.
// It initializes the logger with a test logger and associates the testing context,
// allowing for better logging and error tracking during tests.
func (r *SchemaRegistry) WithTestInstance(t *testing.T) *SchemaRegistry {
	r.l = logging.GetTestLogger(t)
	r.t = t
	return r
}

// WithContainerName sets the container name for the SchemaRegistry instance.
// This allows users to customize the naming of the container for better organization
// and identification within their application. It returns the updated SchemaRegistry.
func (r *SchemaRegistry) WithContainerName(name string) *SchemaRegistry {
	r.ContainerName = name
	return r
}

// WithKafka sets the Kafka bootstrap servers for the Schema Registry using the provided URL.
// It returns the updated SchemaRegistry instance with the new environment variables applied.
func (r *SchemaRegistry) WithKafka(kafkaUrl string) *SchemaRegistry {
	envVars := map[string]string{
		"SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS": kafkaUrl,
	}
	return r.WithEnvVars(envVars)
}

// WithEnvVars merges the provided environment variables into the SchemaRegistry's existing set.
// It allows for dynamic configuration of the Schema Registry based on the specified environment settings.
func (r *SchemaRegistry) WithEnvVars(envVars map[string]string) *SchemaRegistry {
	if err := mergo.Merge(&r.EnvVars, envVars, mergo.WithOverride); err != nil {
		r.l.Fatal().Err(err).Msg("Failed to merge env vars")
	}
	return r
}

// StartContainer initializes and starts a Schema Registry container.
// It sets up the necessary environment variables and logs the internal
// and external URLs for accessing the container. This function is
// essential for users needing a running instance of Schema Registry
// for testing or development purposes.
func (r *SchemaRegistry) StartContainer() error {
	r.InternalUrl = fmt.Sprintf("http://%s:%s", r.ContainerName, "8081")
	l := logging.GetTestContainersGoTestLogger(r.t)
	envVars := map[string]string{
		"SCHEMA_REGISTRY_HOST_NAME": r.ContainerName,
		"SCHEMA_REGISTRY_LISTENERS": r.InternalUrl,
	}
	r.WithEnvVars(envVars)
	cr, err := r.getContainerRequest()
	if err != nil {
		return err
	}
	req := tc.GenericContainerRequest{
		ContainerRequest: cr,
		Started:          true,
		Reuse:            true,
		Logger:           l,
	}
	c, err := tc.GenericContainer(testcontext.Get(r.t), req)
	if err != nil {
		return fmt.Errorf("cannot start Schema Registry container: %w", err)
	}
	host, err := GetHost(testcontext.Get(r.t), c)
	if err != nil {
		return err
	}
	port, err := c.MappedPort(testcontext.Get(r.t), "8081/tcp")
	if err != nil {
		return err
	}
	r.ExternalUrl = fmt.Sprintf("%s:%s", host, port.Port())
	r.l.Info().
		Str("containerName", r.ContainerName).
		Str("internalUrl", r.InternalUrl).
		Str("externalUrl", r.ExternalUrl).
		Msgf("Started Schema Registry container")

	r.Container = c

	return nil
}

func (r *SchemaRegistry) getContainerRequest() (tc.ContainerRequest, error) {
	schemaImage := mirror.AddMirrorToImageIfSet(defaultSchemaRegistryImage)
	return tc.ContainerRequest{
		Name:         r.ContainerName,
		Image:        schemaImage,
		ExposedPorts: []string{"8081/tcp"},
		Env:          r.EnvVars,
		Networks:     r.Networks,
		WaitingFor: tcwait.ForLog("INFO Server started, listening for requests").
			WithStartupTimeout(r.StartupTimeout).
			WithPollInterval(100 * time.Millisecond),
	}, nil
}
