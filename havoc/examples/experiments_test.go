package havoc_example

import (
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	"github.com/stretchr/testify/require"
	"testing"
)

func createMonkey(t *testing.T, l zerolog.Logger, namespace string) *havoc.Controller {
	havoc.SetGlobalLogger(l)
	cfg, err := havoc.ReadConfig("config.toml")
	require.NoError(t, err)
	c, err := havoc.NewController(cfg)
	err = c.GenerateSpecs(namespace)
	require.NoError(t, err)
	return c
}

func TestMyLoad(t *testing.T) {
	/* my testing logger */
	l := havoc.L
	/* my load test preparation here */
	/* wrapping with chaos monkey */
	monkey := createMonkey(t, l, "my namespace, get it from config")
	go monkey.Run()
	/* my test runs and ends */
	errs := monkey.Stop()
	require.Len(t, errs, 0)
}

func TestCodeRun(t *testing.T) {
	cfg, err := havoc.ReadConfig("../havoc.toml")
	require.NoError(t, err)
	c, err := havoc.NewController(cfg)
	require.NoError(t, err)
	nexp, err := havoc.NewNamedExperiment("../experiments-crib-core/failure/failure-app-node-1-bootstrap-69fb558d9-s7npw.yaml")
	require.NoError(t, err)
	err = c.ApplyAndAnnotate(nexp)
	require.NoError(t, err)
}
