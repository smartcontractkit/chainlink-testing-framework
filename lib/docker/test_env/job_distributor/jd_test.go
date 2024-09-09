package job_distributor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

func TestJDSpinUp(t *testing.T) {
	t.Skipf("TODO enable this when jd image is available in ci")
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	// create a postgres first
	pg, err := test_env.NewPostgresDb(
		[]string{network.Name},
		test_env.WithPostgresDbName("jd-db"),
		test_env.WithPostgresImageVersion("14.1"))
	require.NoError(t, err)
	err = pg.StartContainer()
	require.NoError(t, err)

	jd := New([]string{network.Name},
		//TODO: replace with actual image
		WithImage("localhost:5001/jd"),
		WithVersion("latest"),
		WithDBURL(pg.InternalURL.String()),
	)

	err = jd.StartContainer()
	require.NoError(t, err)
	// create a jd connection
	_, err = newConnection(jd.Grpc)
	require.NoError(t, err)
}

func newConnection(target string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to service at %s. Err: %w", target, err)
	}

	return conn, nil
}
