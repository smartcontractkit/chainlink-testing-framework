package soak_runner

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/stretchr/testify/require"
)

func TestSoakOCR(t *testing.T) {
	t.Parallel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	env, err := environment.DeploySoakEnvironment(
		environment.NewSoakChainlinkConfig(environment.ChainlinkReplicas(6, nil)),
		"@soak-ocr",
		tools.ChartsRoot,
	)
	require.NoError(t, err)
	require.NotNil(t, env)
	log.Info().Str("Namespace", env.Namespace).Msg("Soak Test Successfully Launched")
}
