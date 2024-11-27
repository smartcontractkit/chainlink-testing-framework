package test_env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

func TestPostgresCustomImageVersionNotInMirror(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	t.Setenv(config.EnvVarInternalDockerRepo, "")
	pgOpt := WithPostgresImageVersion("16.1")
	pg, err := NewPostgresDb([]string{network.Name}, pgOpt)
	require.NoError(t, err)
	err = pg.StartContainer()
	require.NoError(t, err)
}

func TestPostgresLowestCompatibleVersion(t *testing.T) {
	// requires internal docker repo to be set
	if os.Getenv(config.EnvVarInternalDockerRepo) == "" {
		t.Skipf("Skipping test because %s is not set", config.EnvVarInternalDockerRepo)
	}
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	pg, err := NewPostgresDb([]string{network.Name})
	require.NoError(t, err)
	err = pg.StartContainer()
	require.NoError(t, err)
}

func TestPostgresModernVersion(t *testing.T) {
	// requires internal docker repo to be set
	if os.Getenv(config.EnvVarInternalDockerRepo) == "" {
		t.Skipf("Skipping test because %s is not set", config.EnvVarInternalDockerRepo)
	}
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	pgOpt := WithPostgresImageVersion("15.6")
	pg, err := NewPostgresDb([]string{network.Name}, pgOpt)
	require.NoError(t, err)
	err = pg.StartContainer()
	require.NoError(t, err)
}
