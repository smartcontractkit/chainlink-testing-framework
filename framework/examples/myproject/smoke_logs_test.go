package examples

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type CfgLogs struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSets           []*ns.Input       `toml:"nodesets" validate:"required"`
}

func TestLogsSmoke(t *testing.T) {
	in, err := framework.Load[CfgLogs](t)
	require.NoError(t, err)
	// most simple checks, save all the logs and check (CRIT|PANIC|FATAL) log levels
	t.Cleanup(func() {
		err := framework.SaveAndCheckLogs(t)
		require.NoError(t, err)
	})

	re := regexp.MustCompile(`name=HeadReporter version=\d+`)
	t.Cleanup(func() {
		err := framework.StreamCTFContainerLogsFanout(
			framework.LogStreamConsumer{
				Name: "custom-regex-assert",
				Consume: func(logStreams map[string]io.ReadCloser) error {
					for name, stream := range logStreams {
						scanner := bufio.NewScanner(stream)
						found := false
						for scanner.Scan() {
							if re.MatchString(scanner.Text()) {
								found = true
								break
							}
						}
						if err := scanner.Err(); err != nil {
							return fmt.Errorf("scan %s: %w", name, err)
						}
						if !found {
							return fmt.Errorf("missing HeadReporter log in %s", name)
						}
					}
					return nil
				},
			},
		)
		require.NoError(t, err)
	})

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
	_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	require.NoError(t, err)
	out, err := ns.NewSharedDBNodeSet(in.NodeSets[0], bc)
	require.NoError(t, err)

	t.Run("test something", func(t *testing.T) {
		for _, n := range out.CLNodes {
			require.NotEmpty(t, n.Node.ExternalURL)
		}
	})
}
