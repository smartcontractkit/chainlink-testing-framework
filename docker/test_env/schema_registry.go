package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

type SchemaRegistry struct {
	EnvComponent
	InternalUrl string
	ExternalUrl string
	l           zerolog.Logger
	t           *testing.T
}

func NewSchemaRegistry(networks []string) *SchemaRegistry {
	id, _ := uuid.NewRandom()
	return &SchemaRegistry{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("schema-registry-%s", id.String()),
			Networks:      networks,
		},
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

func (r *SchemaRegistry) StartContainer(envVars map[string]string) error {
	r.InternalUrl = fmt.Sprintf("http://%s:%s", r.ContainerName, "8081")

	l := tc.Logger
	if r.t != nil {
		l = logging.CustomT{
			T: r.t,
			L: r.l,
		}
	}
	req := tc.GenericContainerRequest{
		ContainerRequest: r.getContainerRequest(envVars),
		Started:          true,
		Reuse:            true,
		Logger:           l,
	}
	c, err := tc.GenericContainer(context.Background(), req)
	if err != nil {
		return errors.Wrapf(err, "cannot start Schema Registry container")
	}

	host, err := c.Host(context.Background())
	if err != nil {
		return err
	}
	port, err := c.MappedPort(context.Background(), "8081/tcp")
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

func (r *SchemaRegistry) getContainerRequest(envVars map[string]string) tc.ContainerRequest {
	defaultValues := map[string]string{
		"SCHEMA_REGISTRY_HOST_NAME":                    r.ContainerName,
		"SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS": "kafka:9092",
		"SCHEMA_REGISTRY_DEBUG":                        "true",
		"SCHEMA_REGISTRY_LISTENERS":                    r.InternalUrl,
	}
	if err := mergo.Merge(&defaultValues, envVars, mergo.WithOverride); err != nil {
		r.l.Fatal().Err(err).Msgf("Cannot merge env vars")
	}
	return tc.ContainerRequest{
		Name:         r.ContainerName,
		Image:        "confluentinc/cp-schema-registry:7.4.0",
		ExposedPorts: []string{"8081/tcp"},
		Env:          defaultValues,
		Networks:     r.Networks,
		WaitingFor: tcwait.ForLog("INFO Server started, listening for requests").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}
}
