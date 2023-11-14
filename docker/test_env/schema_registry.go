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

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

type SchemaRegistry struct {
	EnvComponent
	EnvVars     map[string]string
	InternalUrl string
	ExternalUrl string
	l           zerolog.Logger
	t           *testing.T
}

func NewSchemaRegistry(networks []string) *SchemaRegistry {
	id, _ := uuid.NewRandom()
	defaultEnvVars := map[string]string{
		"SCHEMA_REGISTRY_DEBUG": "true",
	}
	return &SchemaRegistry{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("schema-registry-%s", id.String()),
			Networks:      networks,
		},
		EnvVars: defaultEnvVars,
	}
}

func (r *SchemaRegistry) WithTestLogger(t *testing.T) *SchemaRegistry {
	r.l = logging.GetTestLogger(t)
	r.t = t
	return r
}

func (r *SchemaRegistry) WithContainerName(name string) *SchemaRegistry {
	r.ContainerName = name
	return r
}

func (r *SchemaRegistry) WithKafka(kafkaUrl string) *SchemaRegistry {
	envVars := map[string]string{
		"SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS": kafkaUrl,
	}
	return r.WithEnvVars(envVars)
}

func (r *SchemaRegistry) WithEnvVars(envVars map[string]string) *SchemaRegistry {
	if err := mergo.Merge(&r.EnvVars, envVars, mergo.WithOverride); err != nil {
		r.l.Fatal().Err(err).Msg("Failed to merge env vars")
	}
	return r
}

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
	c, err := tc.GenericContainer(utils.TestContext(r.t), req)
	if err != nil {
		return fmt.Errorf("cannot start Schema Registry container: %w", err)
	}
	host, err := GetHost(utils.TestContext(r.t), c)
	if err != nil {
		return err
	}
	port, err := c.MappedPort(utils.TestContext(r.t), "8081/tcp")
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
	schemaImage, err := mirror.GetImage("confluentinc/cp-schema-registry")
	if err != nil {
		return tc.ContainerRequest{}, err
	}
	return tc.ContainerRequest{
		Name:         r.ContainerName,
		Image:        schemaImage,
		ExposedPorts: []string{"8081/tcp"},
		Env:          r.EnvVars,
		Networks:     r.Networks,
		WaitingFor: tcwait.ForLog("INFO Server started, listening for requests").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}, nil
}
