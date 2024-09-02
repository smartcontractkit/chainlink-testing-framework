package examples_wasp

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/seth"
)

type ExampleGun struct {
	client *seth.Client
	Data   []string
}

func NewExampleHTTPGun(client *seth.Client) *ExampleGun {
	return &ExampleGun{
		client: client,
		Data:   make([]string, 0),
	}
}

func (m *ExampleGun) Call(_ *wasp.Generator) *wasp.Response {
	_, err := m.client.Decode(
		TestEnv.DebugContract.AddCounter(m.client.NewTXKeyOpts(m.client.AnySyncedKey()), big.NewInt(0), big.NewInt(1)),
	)
	if err != nil {
		return &wasp.Response{Error: errors.Join(err).Error()}
	}
	return &wasp.Response{}
}

func TestWithWasp(t *testing.T) {
	t.Setenv(seth.ROOT_PRIVATE_KEY_ENV_VAR, "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	t.Setenv(seth.CONFIG_FILE_ENV_VAR, "seth.toml")
	cfg, err := seth.ReadConfig()
	require.NoError(t, err, "failed to read config")
	c, err := seth.NewClientWithConfig(cfg)
	require.NoError(t, err, "failed to initialise seth")
	labels := map[string]string{
		"go_test_name": "TestWithWasp",
		"gen_name":     "TestWithWasp",
		"branch":       "TestWithWasp",
		"commit":       "TestWithWasp",
	}
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.RPS,
		Schedule: wasp.CombineAndRepeat(
			2,
			wasp.Plain(2, 30*time.Second),
			wasp.Plain(10, 30*time.Second),
			wasp.Plain(2, 30*time.Second),
		),
		Gun:        NewExampleHTTPGun(c),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)
	gen.Run(true)
}
