package examples

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
	components "github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components"
)

type Config struct {
	LocalS3Config *s3provider.Input `toml:"local_s3" validate:"required"`
}

func TestLocalS3(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	output, err := s3provider.NewMinioFactory().NewFrom(in.LocalS3Config)
	require.NoError(t, err)

	t.Run("verify that container can be accessed from host", func(t *testing.T) {
		resp, err := http.Get(output.ConsoleURL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("verify that container can be accesses internally", func(t *testing.T) {
		err := components.NewDockerFakeTester(output.ConsoleBaseURL)
		require.NoError(t, err)
	})
}
